# Threat Model — DevSecOps Reference Pipeline

## Methodology

This document applies the **STRIDE** methodology to the CI/CD chain itself — not to the application it deploys. A compromised delivery pipeline is a critical attack vector: it has access to secrets, the container registry, and the production cluster.

> *Source: the CERT 2025 report on cybersecurity explicitly cites CI/CD compromise as a privilege escalation vector, alongside ADCS and hypervisor tooling.*

### STRIDE Categories

| Letter | Threat | Description |
|---|---|---|
| **S** | Spoofing | Identity impersonation |
| **T** | Tampering | Data or code alteration |
| **R** | Repudiation | Inability to trace an action |
| **I** | Information Disclosure | Sensitive data leakage |
| **D** | Denial of Service | Service interruption |
| **E** | Elevation of Privilege | Privilege escalation |

---

## Scope

```
[Developer workstation] → [GitHub repository] → [GitHub Actions runners]
→ [Container Registry] → [ArgoCD] → [Production K8s cluster]
```

Each component is a potential attack surface.

---

## 1. Developer Workstation

### Threats

| ID | STRIDE | Threat | Severity |
|---|---|---|---|
| DEV-01 | I | Hardcoded secret in source code (API key, token, password) | 🔴 Critical |
| DEV-02 | T | Compromised dependency installed locally (npm/pip/go typosquatting) | 🔴 Critical |
| DEV-03 | S | Git identity spoofing (commit with a forged author) | 🟡 Medium |
| DEV-04 | T | Pre-commit hook disabled or intentionally bypassed | 🟡 Medium |

### Mitigations

| ID | Mitigation | Tool |
|---|---|---|
| DEV-01 | Automated secret scanning on commit | Gitleaks pre-commit |
| DEV-02 | Locked dependency versions + SCA scan | go.sum / package-lock.json + Trivy |
| DEV-03 | GPG commit signing (recommended, non-blocking) | Git GPG |
| DEV-04 | Independent secret re-scan in CI regardless of local hook | Gitleaks CI |

---

## 2. GitHub Repository

### Threats

| ID | STRIDE | Threat | Severity |
|---|---|---|---|
| REPO-01 | T | Direct push to `main` without review (CI bypass) | 🔴 Critical |
| REPO-02 | E | Overly permissive GitHub access for contributors | 🔴 Critical |
| REPO-03 | T | Malicious modification of a CI workflow file (`.github/workflows/`) | 🔴 Critical |
| REPO-04 | I | Secret exposure in GitHub Actions logs | 🔴 Critical |
| REPO-05 | T | Supply chain attack via a compromised GitHub Action (`uses: author/action@v1`) | 🔴 Critical |
| REPO-06 | R | No traceability on merges and deployments | 🟠 High |

### Mitigations

| ID | Mitigation | Tool / Config |
|---|---|---|
| REPO-01 | Branch protection on `main`: required PR + green CI | GitHub Branch Protection Rules |
| REPO-02 | Minimum permissions per role (read / write / admin) | GitHub Repository Settings |
| REPO-03 | Mandatory review of workflow file changes | CODEOWNERS on `.github/workflows/` |
| REPO-04 | Automatic secret masking in logs | GitHub Secrets (never plaintext in steps) |
| REPO-05 | Pin actions by commit SHA | `uses: actions/checkout@a5ac7e51b41094c92402da3b24376905380afc29` |
| REPO-06 | GitHub Audit Log enabled + alerts on sensitive actions | GitHub Audit Log |

> ⚠️ **REPO-05 is critical**: a GitHub Action referenced by tag (`@v3`) can be modified at any time by its author. SHA pinning guarantees immutability. See: Codecov incident 2021, tj-actions incident 2025.

---

## 3. GitHub Actions Runners

### Threats

| ID | STRIDE | Threat | Severity |
|---|---|---|---|
| RUNNER-01 | E | CI secret exfiltration from the runner (env variables, temp files) | 🔴 Critical |
| RUNNER-02 | T | Command injection via workflow inputs (script injection) | 🔴 Critical |
| RUNNER-03 | T | Persistent compromised runner between jobs (self-hosted runner) | 🔴 Critical |
| RUNNER-04 | D | Abusive CI minute consumption (DoS on the pipeline) | 🟡 Medium |
| RUNNER-05 | I | Access to secrets from other repositories via a shared runner | 🟠 High |

### Mitigations

| ID | Mitigation | Tool / Config |
|---|---|---|
| RUNNER-01 | Use only GitHub Secrets, never plaintext variables | GitHub Secrets + `${{ secrets.NAME }}` |
| RUNNER-02 | Never interpolate user inputs directly in `run:` | Pass through intermediate environment variables |
| RUNNER-03 | Prefer GitHub-hosted runners (ephemeral by nature) over self-hosted | GitHub-hosted runners |
| RUNNER-04 | Job timeout + concurrency limit | `timeout-minutes` + `concurrency` in workflow |
| RUNNER-05 | Isolate self-hosted runners per project if used | Runner groups |

> ℹ️ **RUNNER-02 concrete example**:
> ```yaml
> # ❌ Dangerous — injection possible if PR title contains backticks
> - run: echo "Title: ${{ github.event.pull_request.title }}"
>
> # ✅ Safe — passed through environment variable
> - env:
>     PR_TITLE: ${{ github.event.pull_request.title }}
>   run: echo "Title: $PR_TITLE"
> ```

---

## 4. Container Registry (GHCR)

### Threats

| ID | STRIDE | Threat | Severity |
|---|---|---|---|
| REG-01 | T | Malicious image pushed to the registry | 🔴 Critical |
| REG-02 | I | Public registry exposing images containing secrets | 🔴 Critical |
| REG-03 | T | Overwriting an existing tag (`latest`) with a compromised image | 🔴 Critical |
| REG-04 | S | Image spoofing via name collision (image pull without digest) | 🟠 High |

### Mitigations

| ID | Mitigation | Tool / Config |
|---|---|---|
| REG-01 | Only the CI runner can push (short-lived token) | GitHub Actions OIDC + GHCR permissions |
| REG-02 | Scan images before push + verify absence of secrets | Trivy + Gitleaks on image filesystem |
| REG-03 | Immutable tagging by commit SHA, not only by `latest` | `ghcr.io/user/app:sha-abc123` |
| REG-04 | Pull images by SHA256 digest in production | `image: ghcr.io/user/app@sha256:...` |

---

## 5. ArgoCD (CD / GitOps)

### Threats

| ID | STRIDE | Threat | Severity |
|---|---|---|---|
| ARGO-01 | E | Unauthorized access to the ArgoCD UI or API | 🔴 Critical |
| ARGO-02 | T | Direct K8s manifest modification bypassing Git | 🟠 High |
| ARGO-03 | I | Kubernetes secrets in plaintext in the config repository | 🔴 Critical |
| ARGO-04 | E | ArgoCD with unnecessarily broad cluster-admin rights | 🟠 High |

### Mitigations

| ID | Mitigation | Tool / Config |
|---|---|---|
| ARGO-01 | SSO + MFA on ArgoCD access | ArgoCD SSO (Dex) |
| ARGO-02 | `syncPolicy: automated` + active drift detection | ArgoCD sync policy |
| ARGO-03 | Encrypt secrets in Git | Sealed Secrets (Bitnami) |
| ARGO-04 | ArgoCD RBAC: scoped permissions per namespace | ArgoCD RBAC config |

---

## 6. Kubernetes Cluster (Runtime)

### Threats

| ID | STRIDE | Threat | Severity |
|---|---|---|---|
| K8S-01 | E | Container running as root with node access | 🔴 Critical |
| K8S-02 | E | Container escape via misconfiguration (hostPID, hostNetwork, privileged) | 🔴 Critical |
| K8S-03 | T | Image with new CVEs published post-deployment | 🟠 High |
| K8S-04 | I | Kubernetes secret accessible by all pods in the namespace | 🟠 High |
| K8S-05 | D | Missing limits/requests → resource starvation | 🟡 Medium |

### Mitigations

| ID | Mitigation | Tool / Config |
|---|---|---|
| K8S-01 | `runAsNonRoot: true` + `readOnlyRootFilesystem: true` in SecurityContext | KubeForge SecurityContext |
| K8S-02 | PodSecurity Admission Controller in `restricted` mode | K8s PSA |
| K8S-03 | Continuous scanning of deployed images | Trivy Operator |
| K8S-04 | RBAC per service account + Sealed Secrets | KubeForge RBAC |
| K8S-05 | Resource limits/requests on all workloads | KubeForge manifests |

---

## 7. Supply Chain (Third-party Dependencies)

### Threats

| ID | STRIDE | Threat | Severity |
|---|---|---|---|
| SC-01 | T | npm/pip/go dependency compromised by an attacker (typosquatting, hijacking) | 🔴 Critical |
| SC-02 | T | Automatic dependency update introducing a vulnerability | 🟠 High |
| SC-03 | T | Malicious package published by a compromised maintainer | 🔴 Critical |

> *Source: CERT-Wavestone 2025 report — the Shai-Hulud worm compromised over 500 npm packages by stealing maintainer tokens and injecting malicious code.*

### Mitigations

| ID | Mitigation | Tool |
|---|---|---|
| SC-01 | Lock exact versions + checksum verification | go.sum / package-lock.json |
| SC-02 | Manual review of Dependabot PRs before merge | Dependabot + branch protection |
| SC-03 | SCA on every build + CVE alerts | Trivy + GitHub Advisory |

---

## Risk Summary by Severity

### 🔴 Critical (immediate fix)
- Hardcoded secrets in code (DEV-01)
- Direct push to main without CI (REPO-01)
- GitHub Actions not pinned by SHA (REPO-05)
- Command injection in workflows (RUNNER-02)
- Public registry with secrets in images (REG-02)
- Kubernetes secrets in plaintext in Git (ARGO-03)
- Containers running as root (K8S-01)

### 🟠 High (planned fix)
- Overpermissive GitHub access (REPO-02)
- Secret exfiltration from CI runner (RUNNER-01)
- Registry tag overwriting (REG-03)
- ArgoCD overly broad permissions (ARGO-04)
- Post-deployment CVEs (K8S-03)

### 🟡 Medium (continuous improvement)
- Bypassable pre-commit hook (DEV-04)
- DoS on the CI pipeline (RUNNER-04)
- Missing K8s resource limits (K8S-05)

---

## References

- [CERT-Wavestone Report 2025](https://www.wavestone.com/fr/insight/2025-wavestone-cert-report/)
- [OWASP Top 10 CI/CD Security Risks](https://owasp.org/www-project-top-10-ci-cd-security-risks/)
- [SLSA Framework (Supply chain Levels for Software Artifacts)](https://slsa.dev/)
- [GitHub Actions Security Hardening](https://docs.github.com/en/actions/security-guides/security-hardening-for-github-actions)
- [Incident tj-actions/changed-files 2025](https://www.wiz.io/blog/new-github-action-supply-chain-attack-reviewdog-tj-actions)
