# Security Policy

## Purpose of This Repository

PipelineForge is a **DevSecOps reference pipeline** designed for educational and demonstration purposes. It intentionally contains vulnerable code in the `apps/` directory to demonstrate the detection capabilities of the integrated security scanning tools (Semgrep, CodeQL, Trivy, Gitleaks).

**The vulnerabilities in `apps/` are deliberate.** They are not oversights — they exist to produce findings in the CI/CD pipeline and serve as concrete examples of what each tool detects.

## What Is Intentionally Vulnerable

| Location | Vulnerability | Purpose |
|---|---|---|
| `apps/go-app/main.go` | SQL injection (`/user` endpoint) | Demonstrates CodeQL data flow analysis |
| `apps/go-app/main.go` | OS command injection (`/ping` endpoint) | Demonstrates Semgrep pattern detection |
| `apps/go-app/main.go` | Hardcoded API key | Demonstrates Gitleaks secret detection |
| `apps/go-app/go.mod` | Outdated dependency (gin v1.7.7) | Demonstrates Trivy SCA scanning |

These are never deployed to any production environment.

## What Is Designed to Be Secure

The CI/CD infrastructure itself follows security best practices:

- All GitHub Actions are pinned by commit SHA (no tag-based references)
- Workflows run with least-privilege `GITHUB_TOKEN` permissions
- Secrets are never interpolated directly in `run:` steps
- Docker images are built using multi-stage builds with minimal runtime layers
- Branch protection is enforced on `main`

## Known Tool Limitations

### Semgrep

Semgrep findings are visible in the CI job logs but do not appear in the GitHub Security Tab. This is a known integration limitation: Semgrep's SARIF output does not include a `level` field on individual results, which GitHub Code Scanning requires to create alerts. The tool itself functions correctly — 4 findings are detected on every run (SQL injection, OS command injection, hardcoded secret, missing Dockerfile user). A Semgrep Cloud account with `SEMGREP_APP_TOKEN` would resolve this by routing results through Semgrep's own platform.

## Reporting a Vulnerability

If you discover a genuine security issue in the **pipeline infrastructure** (workflows, Dockerfile patterns, CI configuration), please report it responsibly:

1. Do not open a public GitHub issue
2. Contact the maintainer directly via GitHub: [@Richonn](https://github.com/Richonn)
3. Include a description of the issue, the affected file, and potential impact

Vulnerabilities in `apps/` are known and intentional — no need to report those.
