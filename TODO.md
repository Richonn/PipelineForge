# DevSecOps Reference Pipeline — TODO

## Project Description

**DevSecOps Reference Pipeline** is an open-source GitHub repository modeling a secure, agile software delivery pipeline from end to end. This is not just a CI/CD pipeline: it is a **documented, justified, and reproducible reference deliverable** — designed as a state-of-the-art example that any team can study, adapt, or deploy.

The project covers three application stacks (Go, Python, Node.js) and integrates security controls at every stage of the software lifecycle: from commit to production deployment. Every tool choice is documented and justified against the risks it addresses.

It targets both DevOps teams looking to secure their pipelines, and consultants or architects seeking a concrete, well-argued reference.

---

## Tech Stack

- **CI/CD**: GitHub Actions
- **Containerization**: Docker, multi-stage builds
- **Orchestration**: Kubernetes (optional, via KubeForge)
- **SAST**: Semgrep + CodeQL
- **Secrets scanning**: Gitleaks (pre-commit + CI)
- **SCA / CVE**: Trivy (images + dependencies)
- **DAST**: OWASP ZAP (smoke test on staging environment)
- **IaC scanning**: Checkov (Dockerfile, K8s manifests)
- **Threat modeling**: structured THREAT_MODEL.md (STRIDE)
- **Reporting**: GitHub Security Tab + automated summaries on PRs

---

## TODO

### Phase 1 — Repository Foundations
- [ ] Create repository structure with main README
- [ ] Write `ARCHITECTURE.md`: full pipeline diagram (Mermaid)
- [ ] Write `THREAT_MODEL.md`: threat modeling of the CI/CD chain itself (STRIDE)
  - What risks on GitHub Actions?
  - What risks on the Docker registry?
  - What risks around secrets in CI?
  - What risks related to third-party dependencies?
- [ ] Write `TOOL_CHOICES.md`: justification of each selected tool vs. alternatives (e.g. why Semgrep over SonarQube?)

### Phase 2 — Demo Applications
- [ ] Create a minimal Go app (simple REST API)
- [ ] Create a minimal Python app (Flask or FastAPI)
- [ ] Create a minimal Node.js app (Express)
- [ ] Each app must intentionally contain fixable vulnerabilities to demonstrate the scans

### Phase 3 — Secure CI/CD Pipeline
- [ ] Base pipeline: lint, build, unit tests
- [ ] Integrate Gitleaks as pre-commit hook (local) + CI stage
- [ ] Integrate Semgrep (SAST) for all three stacks
- [ ] Integrate CodeQL for Go and JavaScript
- [ ] Integrate Trivy for Docker image and dependency scanning
- [ ] Integrate Checkov for Dockerfiles and K8s manifests
- [ ] Integrate OWASP ZAP in baseline mode on a staging environment (DAST)
- [ ] Configure blocking thresholds: Critical/High block merge, Medium warns
- [ ] Publish results to GitHub Security Tab (SARIF)
- [ ] Generate an automated security summary on each PR

### Phase 4 — Secrets Management and Supply Chain
- [ ] Set up secrets management via GitHub Secrets + document best practices
- [ ] Pin GitHub Actions by SHA (`uses: actions/checkout@SHA` instead of `@v3`)
- [ ] Enable Dependabot for automated dependency updates
- [ ] Document supply chain risks in `THREAT_MODEL.md` (link to CERT-Wavestone 2025 report: 500+ compromised npm packages)

### Phase 5 — State-of-the-Art Documentation
- [ ] Write `DEVSECOPS_REFERENCE.md`: full pipeline guide
  - Shift Left: definition and practical implementation
  - SAST vs DAST vs SCA: when to use which
  - False positive management: strategy and configuration
  - Agile pipeline vs secure pipeline: how to reconcile both
- [ ] Add badges to README (pipeline status, Trivy, Semgrep, license)
- [ ] Create a visual diagram of the full chain (Mermaid in README)
- [ ] Write a blog post or LinkedIn article at publication (optional)

### Phase 6 — Polish & Publication
- [ ] Verify each tool is properly configured, not just present
- [ ] Ensure the pipeline runs without errors on all three apps
- [ ] Full documentation review
- [ ] Publish the repository publicly on GitHub (Richonn account)
- [ ] Submit to r/devops and r/netsec for feedback

---

## What This Project Demonstrates

| Skill | Evidence |
|---|---|
| Knowledge of CI/CD tools and their risks | `THREAT_MODEL.md` + `TOOL_CHOICES.md` |
| Integration of security tools in a delivery chain | Functional GitHub Actions pipeline |
| State-of-the-art pipeline representation | `ARCHITECTURE.md` + `DEVSECOPS_REFERENCE.md` |
| Concrete working implementation | All three apps + running pipeline |
| Organizational dimension | Process documentation and best practices |

---

## Related Projects

- **ShieldCI** → generates pipelines; this project is the "manual reference" version of what ShieldCI automates
- **KubeForge** → covers runtime/CD; this project covers CI and the build phase

Together, the three projects provide complete DevSecOps coverage from commit to cluster.
