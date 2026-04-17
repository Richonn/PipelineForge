# DevSecOps Reference Guide

This document covers the principles and practices behind the PipelineForge pipeline. It is not a tutorial — it is a reference that explains the *why* behind each architectural decision.

---

## Shift Left

### Definition

"Shift Left" means moving security controls earlier in the development lifecycle. The term comes from reading a pipeline left to right: the further left a check occurs, the earlier it catches issues.

```
[Local commit] -> [CI pipeline] -> [Registry] -> [Runtime]
    faster                                           slower
    cheaper                                          costlier
    earlier                                           later
```

A secret caught at the pre-commit hook stage requires no action beyond removing it from the uncommitted change. The same secret caught in production requires credential rotation, incident response, audit review, and customer notification.

### How It Is Applied Here

| Stage | Tool | What gets caught early |
|---|---|---|
| Pre-commit | Gitleaks | Secrets before they enter the repository |
| Pre-commit | Linter | Code quality issues before they reach CI |
| CI — SAST | Semgrep, CodeQL | Vulnerability patterns before the code is deployed |
| CI — SCA | Trivy | Known CVEs before the image is built |
| CI — Image | Trivy | New CVEs introduced by the base image |
| CI — IaC | Checkov | Misconfigurations before manifests reach the cluster |

The further right a vulnerability is caught, the more expensive it is to fix. Shift Left is a cost optimization strategy as much as a security strategy.

---

## SAST vs DAST vs SCA

Three distinct families of security analysis tools are used in this pipeline. They are not interchangeable — each targets a different phase and a different class of vulnerability.

### SAST — Static Application Security Testing

**What it does**: Analyzes source code without executing it. Looks for patterns, data flows, and structures that indicate a vulnerability.

**When it runs**: During CI, before the build.

**What it catches**: SQL injections, command injections, hardcoded secrets, unsafe function calls, insecure cryptographic usage.

**What it misses**: Runtime behavior, environment-specific issues, authentication flaws that only manifest during execution.

**Tools used**: Semgrep (pattern matching), CodeQL (data flow analysis).

**Key distinction between Semgrep and CodeQL**:

- Semgrep matches code *patterns*. It can detect `exec(user_input)` by recognizing the structural pattern of a dangerous function call with a variable argument.
- CodeQL traces *data flows*. It can detect that a value entering via an HTTP parameter (`c.Query("name")`) eventually reaches a SQL query (`db.Query(...)`) three function calls later — even if no single line looks obviously dangerous.

Both are necessary. Pattern matching is fast and catches common mistakes. Data flow analysis catches complex vulnerabilities that span multiple files and abstraction layers.

### DAST — Dynamic Application Security Testing

**What it does**: Tests the running application by sending malformed HTTP requests and observing responses.

**When it runs**: Against a deployed staging environment, after the build.

**What it catches**: HTTP-level vulnerabilities (missing security headers, XSS reflected in responses, open redirects, basic injection points exposed via the API surface).

**What it misses**: Vulnerabilities in code paths that are not exposed via HTTP, logic flaws, authentication issues that require session state.

**Tool used**: OWASP ZAP (baseline scan mode).

**Limitation**: The ZAP baseline scan is a safety net, not a penetration test. It runs a fixed set of passive and active checks in under 5 minutes. It will catch obvious misconfigurations but will not find complex chained vulnerabilities. Its value is in catching regressions — a security header that was present last week but is now missing.

### SCA — Software Composition Analysis

**What it does**: Inventories third-party dependencies and checks them against CVE databases.

**When it runs**: Twice — once against dependency manifests (go.sum, package-lock.json) in CI, and once against the built Docker image.

**What it catches**: Known vulnerabilities in libraries and base OS packages, with CVE identifiers, CVSS scores, and fix versions.

**What it misses**: Zero-day vulnerabilities not yet in any database, vulnerabilities introduced by how a library is *used* (misuse of a cryptographic library, for example).

**Tool used**: Trivy.

**Why scan twice**: Dependency manifests list the application's declared dependencies. The Docker image contains the application's dependencies *plus* the OS packages installed in the base image (Alpine, distroless). An image can be vulnerable even if all application dependencies are clean.

---

## False Positive Management

An overly noisy pipeline is a broken pipeline. If every CI run produces dozens of warnings that developers learn to ignore, the pipeline loses its value. False positive management is not optional — it is a prerequisite for a pipeline that people actually respect.

### Severity Policy

| Level | Definition | Action |
|---|---|---|
| Critical | CVSS >= 9.0, or confirmed exploitable vulnerability | Blocks merge — immediate fix required |
| High | CVSS >= 7.0, or likely exploitable vulnerability | Blocks merge — fix or documented exception required |
| Medium | CVSS 4.0–6.9, or theoretical vulnerability | Warning on PR — does not block |
| Low / Info | CVSS < 4.0, or informational finding | Ignored in CI, visible in dashboard |

### Exception Management

When a finding is a confirmed false positive, it must be suppressed with a documented justification — not silently ignored.

**Semgrep**: Add a comment on the flagged line:
```go
result := db.Query(query) // nosemgrep: go.lang.security.audit.sqli
```

Or use `.semgrepignore`:
```
# False positive: query is built from a validated allowlist, not user input
apps/go-app/internal/reports/query_builder.go
```

**Trivy**: Use `.trivyignore`:
```
# CVE-2023-XXXX: not exploitable in our deployment context (no network exposure)
# Reviewed: 2026-04-16, owner: @Richonn
CVE-2023-XXXX
```

**Rule**: Every suppression requires a justification comment. "We'll fix it later" is not a justification. Suppressions without justification are treated as technical debt and flagged in code review.

### The Cost of Over-Blocking

Blocking on Medium and Low findings creates alert fatigue and incentivizes engineers to disable scans entirely. The Critical/High blocking threshold is calibrated to catch real risk while keeping the pipeline navigable. A finding that cannot be exploited in the current deployment context is noise — document it and move on.

---

## Agile Pipeline vs Secure Pipeline

The common assumption is that security and velocity are in tension. This is true only when security is bolted on after the fact. A pipeline designed with security from the start can be both fast and secure.

### Where Security Adds Latency

| Stage | Typical duration | Mitigation |
|---|---|---|
| Gitleaks pre-commit | 1–3 seconds | Negligible |
| Semgrep SAST | 30–90 seconds | Runs in parallel with other jobs |
| CodeQL analysis | 3–8 minutes | Runs in parallel, not on the critical path for small PRs |
| Trivy SCA | 15–30 seconds | Fast — local CVE database cache |
| Trivy image scan | 30–60 seconds | Runs after build, parallel to other checks |
| Checkov IaC | 10–20 seconds | Fast |

Total additional latency over a basic build-and-test pipeline: approximately 3–8 minutes, depending on codebase size. This is acceptable. A developer who introduces a Critical vulnerability and has it caught in 5 minutes saves hours of incident response.

### Concurrency as the Key Enabler

The pipeline is structured so that independent security jobs run in parallel:

```
build
  |
  +-- gitleaks
  +-- semgrep       (parallel)
  +-- codeql        (parallel)
  +-- trivy SCA     (parallel)
  |
docker build
  |
  +-- trivy image   (parallel)
  +-- checkov       (parallel)
```

The wall-clock time of the security pipeline is determined by the slowest parallel job (CodeQL), not the sum of all jobs.

### Developer Experience Considerations

Security tooling that creates friction without feedback will be circumvented. The pipeline is designed to give actionable feedback at the right level:

- **Pre-commit hooks** fail fast with a specific message and line number
- **CI findings** are annotated directly on the PR diff — the developer sees the issue on the exact line, not in a separate report
- **SARIF integration** routes all findings to a single Security Tab — no need to check multiple dashboards
- **Severity thresholds** ensure that only actionable findings block the merge — informational noise is filtered

A developer who never gets blocked by a false positive, and always gets a clear explanation when they are blocked by a real finding, will trust the pipeline.

---

## Supply Chain Security

The CI/CD pipeline itself is a supply chain attack surface. Every external dependency introduced into a workflow is a potential vector.

### GitHub Actions Pinning

Every action in this pipeline is pinned by commit SHA, not by tag:

```yaml
# Vulnerable — the tag can be moved at any time
- uses: actions/checkout@v4

# Secure — the SHA is immutable
- uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd
```

A tag like `@v4` is a mutable pointer. The action author (or an attacker who has compromised the author's account) can push a new commit to that tag at any time. SHA pinning ensures that the exact code that was reviewed is the code that runs — always.

**Reference**: The tj-actions/changed-files incident (2025) demonstrated this attack at scale. A compromised GitHub token was used to move action tags to malicious commits, exfiltrating CI secrets from thousands of repositories.

### Least-Privilege Tokens

Every workflow job declares only the permissions it needs:

```yaml
permissions:
  contents: read      # checkout only
  security-events: write  # SARIF upload only
  packages: write     # GHCR push only
```

The default `GITHUB_TOKEN` has broad permissions. Explicit permission scoping ensures that a compromised job cannot write to the repository, create releases, or perform other privileged operations.

### Dependency Locking

All dependency versions are pinned and checksummed:

- **Go**: `go.sum` contains cryptographic hashes of every dependency
- **Node.js**: `package-lock.json` locks the full dependency tree
- **Python**: `requirements.txt` with pinned versions

`go mod tidy` without a lock file is equivalent to running untrusted code from the internet on every build.
