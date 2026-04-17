# Tool Choices — DevSecOps Reference Pipeline

This document justifies every tool selected for this pipeline against its alternatives. Each choice reflects a specific threat, a specific constraint, or a deliberate trade-off. Tools were not selected for popularity — they were selected because they are the right fit for the job.

---

## Secrets Scanning — Gitleaks

**Selected over**: TruffleHog, detect-secrets, GitGuardian

| Criterion | Rationale |
|---|---|
| Pre-commit support | Runs as a local hook before any push — secrets are caught before they enter the repository |
| Git history scanning | Scans the full commit history, not just the current diff — past leaks are not missed |
| SARIF output | Native SARIF export for GitHub Security Tab integration |
| Zero configuration | Works out of the box with 150+ built-in rules covering all major secret patterns |
| Self-hosted | No data leaves the runner — no external service dependency |

**Why not TruffleHog**: TruffleHog is powerful for deep history scanning but slower and more complex to configure for CI integration. Gitleaks covers the same detection surface with less operational overhead.

**Why not GitGuardian**: GitGuardian requires a cloud account and sends code to an external service. Unacceptable for pipelines handling proprietary code.

---

## SAST — Semgrep

**Selected over**: SonarQube, Bandit, gosec, ESLint security plugins

| Criterion | Rationale |
|---|---|
| Multi-language | Single tool covers Go, Python, and Node.js — no need to maintain separate scanners per stack |
| AST-based matching | Rules match code structure, not text — immune to formatting differences and line breaks |
| Community rule sets | 1000+ production-ready rules maintained by the community (`p/default`, `p/golang`, `p/python`) |
| Custom rules | Simple YAML syntax to write project-specific rules when needed |
| CI performance | Typically completes in under 60 seconds on small codebases |

**Why not SonarQube**: SonarQube requires a server (SonarCloud or self-hosted), adds significant operational complexity, and is primarily designed for code quality — security rules are available but not the core focus. Semgrep is security-first and requires no infrastructure.

**Why not gosec / Bandit**: Language-specific tools. Using both gosec (Go) and Bandit (Python) means maintaining two separate configs, two rule sets, and two integrations. Semgrep replaces both with a single workflow step.

**Known limitation**: Semgrep's SARIF output omits the `level` field on individual results. GitHub Code Scanning requires this field to create alerts. Findings are visible in CI logs but do not appear in the GitHub Security Tab without a Semgrep Cloud account. See [SECURITY.md](./SECURITY.md) for details.

---

## SAST — CodeQL

**Selected alongside Semgrep, not instead of it**

| Criterion | Rationale |
|---|---|
| Data flow analysis | Traces tainted data across function boundaries — detects SQL injections and command injections that Semgrep's pattern matching would miss |
| GitHub-native | Built into GitHub Advanced Security — no additional infrastructure, native SARIF integration |
| Deep vulnerability detection | Detects complex vulnerability chains that require understanding program semantics |

**Why both Semgrep and CodeQL**: They are complementary, not redundant. Semgrep catches fast, pattern-based issues (unsafe function calls, missing security headers, bad crypto usage). CodeQL catches complex data flow vulnerabilities (tainted input reaching a sink). Running both maximizes coverage.

**Trade-off**: CodeQL adds 3–5 minutes to the pipeline. Acceptable given the depth of analysis it provides.

---

## SCA — Trivy

**Selected over**: Snyk, OWASP Dependency-Check, Grype

| Criterion | Rationale |
|---|---|
| Dual-mode scanning | Scans both dependency manifests (go.sum, package-lock.json) and Docker image layers in a single tool |
| CVE database coverage | Aggregates NVD, OSV, and GitHub Advisory — one of the broadest coverage sets available |
| SARIF output | Native integration with GitHub Security Tab |
| Speed | Significantly faster than OWASP Dependency-Check for equivalent coverage |
| Self-hosted | No API key or external service required |

**Why not Snyk**: Snyk requires an API token and sends dependency data to an external service. It also enforces rate limits on the free tier, which can block CI pipelines. Trivy is fully self-contained.

**Why not OWASP Dependency-Check**: Slow (typically 5–10 minutes), requires downloading a large CVE database on every run, and lacks native Docker image scanning. Trivy covers the same surface faster and more reliably.

---

## IaC Scanning — Checkov

**Selected over**: Trivy (IaC mode), tfsec, Terrascan, Hadolint

| Criterion | Rationale |
|---|---|
| Broad framework support | Covers Dockerfiles, Kubernetes manifests, Terraform, and Helm in a single tool |
| Docker-specific rules | Purpose-built checks for Dockerfile security (USER, HEALTHCHECK, COPY, etc.) |
| SARIF output | Native GitHub Security Tab integration |
| Rule depth | 1000+ built-in policies across all supported frameworks |

**Why not Hadolint**: Hadolint is excellent for Dockerfile linting but limited to Dockerfiles only. When Kubernetes manifests are added to the pipeline, a second tool would be needed. Checkov covers both.

**Why not tfsec**: tfsec is Terraform-specific. This pipeline does not use Terraform — Checkov's broader coverage is more appropriate.

---

## CI/CD — GitHub Actions

**Selected over**: GitLab CI, Jenkins, CircleCI, Drone CI

| Criterion | Rationale |
|---|---|
| Zero infrastructure | No server to maintain — compute is fully managed by GitHub |
| Native security integrations | GitHub Advanced Security, SARIF upload, Dependabot, and secret scanning are all first-party |
| Reusable workflows | Modular pipeline composition without external orchestration tools |
| OIDC support | Keyless authentication to cloud providers and registries without stored credentials |
| Ecosystem | The largest marketplace of pre-built actions for security tooling |

**Why not Jenkins**: Jenkins requires infrastructure management, plugin maintenance, and significant operational overhead. The security posture of Jenkins pipelines is harder to enforce consistently.

---

## Container Registry — GHCR

**Selected over**: Docker Hub, AWS ECR, Google Artifact Registry

| Criterion | Rationale |
|---|---|
| GitHub-native | Authentication via `GITHUB_TOKEN` — no additional credentials to manage |
| Free for public repositories | No cost for open-source projects |
| Package visibility tied to repo | Registry permissions mirror repository permissions automatically |

---

## CD — ArgoCD

**Selected over**: Flux, Spinnaker, direct `kubectl apply`

| Criterion | Rationale |
|---|---|
| GitOps model | Cluster state is always derived from Git — full auditability, no manual interventions |
| Drift detection | Automatically detects and alerts on out-of-band changes to the cluster |
| Separation of concerns | The CI pipeline never has direct cluster access — it only pushes to a registry |
| RBAC | Fine-grained access control per application and namespace |

**Why not Flux**: Flux and ArgoCD are functionally equivalent for this use case. ArgoCD was chosen for its UI, which makes the GitOps synchronization state visible and easier to demonstrate.

**Why not direct kubectl**: Direct deployment from CI gives the pipeline cluster-admin access, which is a critical security risk. The GitOps model ensures the cluster only pulls from the registry.

---

## Dependency Updates — Dependabot

**Selected over**: Renovate, manual updates

| Criterion | Rationale |
|---|---|
| GitHub-native | No additional configuration or external service |
| Automatic PRs | Opens pull requests for outdated dependencies — updates go through the full CI pipeline before merge |
| Security alerts | Integrates with GitHub Advisory Database to flag vulnerable dependencies proactively |

**Why not Renovate**: Renovate is more configurable but requires a separate bot account or GitHub App. Dependabot is zero-setup for GitHub repositories and sufficient for this pipeline's needs.
