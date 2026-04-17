# DevSecOps Reference Pipeline — TODO

## Project Description

**DevSecOps Reference Pipeline** is an open-source GitHub repository modeling a secure, agile software delivery pipeline from end to end. This is not just a CI/CD pipeline: it is a **documented, justified, and reproducible reference deliverable** — designed as a state-of-the-art example that any team can study, adapt, or deploy.

The project covers a Go application stack and integrates security controls at every stage of the software lifecycle: from commit to production deployment. Every tool choice is documented and justified against the risks it addresses.

It targets both DevOps teams looking to secure their pipelines, and consultants or architects seeking a concrete, well-argued reference.

---

## Tech Stack

- **CI/CD**: GitHub Actions (reusable workflows)
- **Containerization**: Docker, multi-stage builds
- **Orchestration**: Kubernetes (optional, via KubeForge)
- **SAST**: Semgrep + CodeQL
- **Secrets scanning**: Gitleaks (pre-commit + CI)
- **SCA / CVE**: Trivy (images + dependencies)
- **DAST**: OWASP ZAP (smoke test on staging environment)
- **IaC scanning**: Checkov (Dockerfile, K8s manifests)
- **Threat modeling**: THREAT_MODEL.md (STRIDE)
- **Reporting**: GitHub Security Tab (SARIF) + job summaries

---

## TODO

### Phase 1 — Repository Foundations
- [x] Create repository structure with main README
- [x] Write `ARCHITECTURE.md`: full pipeline diagram (Mermaid)
- [x] Write `THREAT_MODEL.md`: threat modeling of the CI/CD chain itself (STRIDE)
- [x] Write `TOOL_CHOICES.md`: justification of each selected tool vs. alternatives

### Phase 2 — Demo Application
- [x] Create a minimal Go app (REST API)
- [x] Intentional vulnerabilities: SQL injection, OS command injection, hardcoded secret, outdated dependency
- [ ] Create a minimal Python app (FastAPI) — optional
- [ ] Create a minimal Node.js app (Express) — optional

### Phase 3 — Secure CI/CD Pipeline
- [x] Base pipeline: build + unit tests
- [x] Integrate Gitleaks (secrets scan)
- [x] Integrate Semgrep (SAST — findings in CI logs)
- [x] Integrate CodeQL (SAST — data flow analysis, findings in Security Tab)
- [x] Integrate Trivy SCA (dependency CVE scan, findings in Security Tab)
- [x] Integrate Trivy image scan (Docker image CVE scan)
- [x] Integrate Checkov (Dockerfile IaC scan)
- [x] Configure blocking thresholds: Critical/High continue-on-error, findings reported
- [x] Publish results to GitHub Security Tab (SARIF)
- [x] Generate job summary on each run
- [ ] Integrate OWASP ZAP (DAST — requires staging environment, out of scope for this reference)

### Phase 4 — Secrets Management and Supply Chain
- [x] GitHub Secrets for all sensitive values
- [x] All GitHub Actions pinned by commit SHA
- [ ] Enable Dependabot for automated dependency updates
- [x] Supply chain risks documented in `THREAT_MODEL.md`

### Phase 5 — State-of-the-Art Documentation
- [x] Write `DEVSECOPS_REFERENCE.md`: Shift Left, SAST/DAST/SCA, false positives, supply chain
- [x] Write `TOOL_CHOICES.md`: tool justification vs alternatives
- [x] Write `SECURITY.md`: known limitations and vulnerability reporting policy
- [ ] Add CI pipeline status badge to README

### Phase 6 — Polish & Publication
- [x] Repository is public on GitHub
- [ ] Full documentation review
- [ ] Submit to r/devops and r/netsec for feedback

---

## What This Project Demonstrates

| Skill | Evidence |
|---|---|
| Knowledge of CI/CD tools and their risks | `THREAT_MODEL.md` + `TOOL_CHOICES.md` |
| Integration of security tools in a delivery chain | Functional GitHub Actions pipeline (4 reusable workflows) |
| State-of-the-art pipeline design | `ARCHITECTURE.md` + `DEVSECOPS_REFERENCE.md` |
| Concrete working implementation | Go app + pipeline with active findings in Security Tab |
| Supply chain hardening | SHA-pinned actions, least-privilege tokens, locked dependencies |
| Organizational dimension | False positive policy, exception management, severity thresholds |

---

## Related Projects

- **ShieldCI** — generates pipelines automatically; PipelineForge is the manually documented reference version
- **KubeForge** — covers Kubernetes runtime security; PipelineForge covers CI and the build phase

Together, the three projects provide complete DevSecOps coverage from commit to cluster.
