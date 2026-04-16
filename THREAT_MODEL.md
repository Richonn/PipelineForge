# Threat Model — DevSecOps Reference Pipeline

## Méthodologie

Ce document applique la méthode **STRIDE** à la chaîne CI/CD elle-même — pas à l'application qu'elle déploie. Une chaîne de production compromise est un vecteur d'attaque critique : elle a accès aux secrets, au registry, et au cluster de production.

> *Source : le rapport CERT 2025 sur la cybersécurité cite explicitement la compromission de la CI/CD comme vecteur d'élévation de privilèges, aux côtés de l'ADCS et des outils d'hypervision.*

### Catégories STRIDE

| Lettre | Menace | Traduction |
|---|---|---|
| **S** | Spoofing | Usurpation d'identité |
| **T** | Tampering | Altération de données ou de code |
| **R** | Repudiation | Impossibilité de retracer une action |
| **I** | Information Disclosure | Fuite d'informations sensibles |
| **D** | Denial of Service | Interruption de service |
| **E** | Elevation of Privilege | Escalade de privilèges |

---

## Périmètre analysé

```
[Poste développeur] → [GitHub repo] → [GitHub Actions runners]
→ [Container Registry] → [ArgoCD] → [Cluster K8s production]
```

Chaque composant est une surface d'attaque potentielle.

---

## 1. Poste développeur

### Menaces

| ID | STRIDE | Menace | Criticité |
|---|---|---|---|
| DEV-01 | I | Secret hardcodé dans le code source (clé API, token, mot de passe) | 🔴 Critical |
| DEV-02 | T | Dépendance compromise installée localement (typosquatting npm/pip/go) | 🔴 Critical |
| DEV-03 | S | Usurpation d'identité Git (commit signé avec un faux auteur) | 🟡 Medium |
| DEV-04 | T | Hook pre-commit désactivé ou contourné intentionnellement | 🟡 Medium |

### Mitigations

| ID | Mitigation | Outil |
|---|---|---|
| DEV-01 | Scan automatique des secrets au commit | Gitleaks pre-commit |
| DEV-02 | Lock des versions de dépendances + scan SCA | go.sum / package-lock.json + Trivy |
| DEV-03 | Signature des commits GPG (recommandé, non bloquant) | Git GPG |
| DEV-04 | Re-scan des secrets en CI indépendamment du hook local | Gitleaks CI |

---

## 2. Dépôt GitHub

### Menaces

| ID | STRIDE | Menace | Criticité |
|---|---|---|---|
| REPO-01 | T | Push direct sur `main` sans review (contournement de la CI) | 🔴 Critical |
| REPO-02 | E | Permissions GitHub trop larges pour les contributors | 🔴 Critical |
| REPO-03 | T | Modification malveillante d'un fichier de workflow CI (`.github/workflows/`) | 🔴 Critical |
| REPO-04 | I | Exposition de secrets dans les logs GitHub Actions | 🔴 Critical |
| REPO-05 | T | Supply chain attack via une GitHub Action compromise (`uses: auteur/action@v1`) | 🔴 Critical |
| REPO-06 | R | Absence de traçabilité sur les merges et déploiements | 🟠 High |

### Mitigations

| ID | Mitigation | Outil / Config |
|---|---|---|
| REPO-01 | Branch protection sur `main` : PR obligatoire + CI verte | GitHub Branch Protection Rules |
| REPO-02 | Permissions minimales par rôle (read / write / admin) | GitHub Repository Settings |
| REPO-03 | Review obligatoire des modifications de workflows | CODEOWNERS sur `.github/workflows/` |
| REPO-04 | Masquage automatique des secrets dans les logs | GitHub Secrets (jamais en clair dans les steps) |
| REPO-05 | Pinning des actions par SHA de commit | `uses: actions/checkout@a5ac7e51b41094c92402da3b24376905380afc29` |
| REPO-06 | Audit log GitHub activé + alertes sur les actions sensibles | GitHub Audit Log |

> ⚠️ **REPO-05 est critique** : une action GitHub référencée par tag (`@v3`) peut être modifiée à tout moment par son auteur. Le pinning par SHA garantit l'immutabilité. Voir : incident Codecov 2021, incident tj-actions 2025.

---

## 3. GitHub Actions Runners

### Menaces

| ID | STRIDE | Menace | Criticité |
|---|---|---|---|
| RUNNER-01 | E | Exfiltration des secrets CI depuis le runner (env variables, fichiers temporaires) | 🔴 Critical |
| RUNNER-02 | T | Injection de commandes dans les inputs de workflow (script injection) | 🔴 Critical |
| RUNNER-03 | T | Runner compromis persistant entre deux jobs (self-hosted runner) | 🔴 Critical |
| RUNNER-04 | D | Consommation abusive des minutes CI (DoS sur le pipeline) | 🟡 Medium |
| RUNNER-05 | I | Accès aux secrets d'autres dépôts via un runner partagé | 🟠 High |

### Mitigations

| ID | Mitigation | Outil / Config |
|---|---|---|
| RUNNER-01 | Utiliser uniquement des secrets GitHub, jamais de variables en clair | GitHub Secrets + `${{ secrets.NOM }}` |
| RUNNER-02 | Ne jamais interpoler d'inputs utilisateur directement dans `run:` | Passer par des variables d'environnement intermédiaires |
| RUNNER-03 | Préférer les runners GitHub-hosted (éphémères par nature) aux self-hosted | GitHub-hosted runners |
| RUNNER-04 | Timeout sur les jobs CI + limite de concurrence | `timeout-minutes` + `concurrency` dans le workflow |
| RUNNER-05 | Isoler les runners self-hosted par projet si utilisés | Runner groups |

> ℹ️ **RUNNER-02 exemple concret** :
> ```yaml
> # ❌ Dangereux — injection possible si PR title contient des backticks
> - run: echo "Titre : ${{ github.event.pull_request.title }}"
>
> # ✅ Sûr — passage par variable d'environnement
> - env:
>     PR_TITLE: ${{ github.event.pull_request.title }}
>   run: echo "Titre : $PR_TITLE"
> ```

---

## 4. Container Registry (GHCR)

### Menaces

| ID | STRIDE | Menace | Criticité |
|---|---|---|---|
| REG-01 | T | Push d'une image malveillante dans le registry | 🔴 Critical |
| REG-02 | I | Registry public exposant des images contenant des secrets | 🔴 Critical |
| REG-03 | T | Écrasement d'un tag existant (`latest`) avec une image compromise | 🔴 Critical |
| REG-04 | S | Usurpation d'image via collision de nom (image pull sans digest) | 🟠 High |

### Mitigations

| ID | Mitigation | Outil / Config |
|---|---|---|
| REG-01 | Seul le runner CI peut pusher (token à durée limitée) | GitHub Actions OIDC + GHCR permissions |
| REG-02 | Scanner les images avant push + vérifier absence de secrets | Trivy + Gitleaks sur le filesystem de l'image |
| REG-03 | Tagging immutable par SHA de commit, pas uniquement par `latest` | `ghcr.io/user/app:sha-abc123` |
| REG-04 | Pull d'images par digest SHA256 en production | `image: ghcr.io/user/app@sha256:...` |

---

## 5. ArgoCD (CD / GitOps)

### Menaces

| ID | STRIDE | Menace | Criticité |
|---|---|---|---|
| ARGO-01 | E | Accès non autorisé à l'interface ArgoCD (UI ou API) | 🔴 Critical |
| ARGO-02 | T | Modification directe des manifests K8s en bypassant Git | 🟠 High |
| ARGO-03 | I | Secrets Kubernetes en clair dans le dépôt de config | 🔴 Critical |
| ARGO-04 | E | ArgoCD avec des droits cluster-admin inutilement larges | 🟠 High |

### Mitigations

| ID | Mitigation | Outil / Config |
|---|---|---|
| ARGO-01 | SSO + MFA sur l'accès ArgoCD | ArgoCD SSO (Dex) |
| ARGO-02 | `syncPolicy: automated` + drift detection active | ArgoCD sync policy |
| ARGO-03 | Chiffrement des secrets dans Git | Sealed Secrets (Bitnami) |
| ARGO-04 | RBAC ArgoCD : droits limités par namespace | ArgoCD RBAC config |

---

## 6. Cluster Kubernetes (Runtime)

### Menaces

| ID | STRIDE | Menace | Criticité |
|---|---|---|---|
| K8S-01 | E | Container s'exécutant en root avec accès au node | 🔴 Critical |
| K8S-02 | E | Escape de container via misconfiguration (hostPID, hostNetwork, privileged) | 🔴 Critical |
| K8S-03 | T | Image avec de nouvelles CVE publiées post-déploiement | 🟠 High |
| K8S-04 | I | Secret Kubernetes accessible par tous les pods du namespace | 🟠 High |
| K8S-05 | D | Absence de limits/requests → resource starvation | 🟡 Medium |

### Mitigations

| ID | Mitigation | Outil / Config |
|---|---|---|
| K8S-01 | `runAsNonRoot: true` + `readOnlyRootFilesystem: true` dans le SecurityContext | KubeForge SecurityContext |
| K8S-02 | PodSecurity Admission Controller en mode `restricted` | K8s PSA |
| K8S-03 | Scan continu des images déployées | Trivy Operator |
| K8S-04 | RBAC par service account + Sealed Secrets | KubeForge RBAC |
| K8S-05 | Resource limits/requests sur tous les workloads | KubeForge manifests |

---

## 7. Supply Chain (dépendances tierces)

### Menaces

| ID | STRIDE | Menace | Criticité |
|---|---|---|---|
| SC-01 | T | Dépendance npm/pip/go compromise par un attaquant (typosquatting, hijacking) | 🔴 Critical |
| SC-02 | T | Mise à jour automatique d'une dépendance introduisant une vulnérabilité | 🟠 High |
| SC-03 | T | Package malveillant publié par un mainteneur compromis | 🔴 Critical |

> *Source : rapport CERT-Wavestone 2025 — le ver Shai-Hulud a compromis plus de 500 paquets npm en volant les tokens des mainteneurs et en injectant du code malveillant.*

### Mitigations

| ID | Mitigation | Outil |
|---|---|---|
| SC-01 | Lock des versions exactes + vérification des checksums | go.sum / package-lock.json |
| SC-02 | Review manuelle des PRs Dependabot avant merge | Dependabot + branch protection |
| SC-03 | SCA à chaque build + alertes CVE | Trivy + GitHub Advisory |

---

## Résumé des risques par criticité

### 🔴 Critical (correction immédiate)
- Secrets hardcodés dans le code (DEV-01)
- Push direct sur main sans CI (REPO-01)
- Actions GitHub non pinnées par SHA (REPO-05)
- Injection de commandes dans les workflows (RUNNER-02)
- Registry public avec secrets dans les images (REG-02)
- Secrets K8s en clair dans Git (ARGO-03)
- Containers en root (K8S-01)

### 🟠 High (correction planifiée)
- Permissions GitHub trop larges (REPO-02)
- Exfiltration depuis runner CI (RUNNER-01)
- Écrasement de tags registry (REG-03)
- ArgoCD droits trop larges (ARGO-04)
- CVE post-déploiement (K8S-03)

### 🟡 Medium (amélioration continue)
- Hook pre-commit contournable (DEV-04)
- DoS sur le pipeline CI (RUNNER-04)
- Absence de resource limits K8s (K8S-05)

---

## Références

- [CERT-Wavestone Report 2025](https://www.wavestone.com/fr/insight/2025-wavestone-cert-report/)
- [OWASP Top 10 CI/CD Security Risks](https://owasp.org/www-project-top-10-ci-cd-security-risks/)
- [SLSA Framework (Supply chain Levels for Software Artifacts)](https://slsa.dev/)
- [GitHub Actions Security Hardening](https://docs.github.com/en/actions/security-guides/security-hardening-for-github-actions)
- [Incident tj-actions/changed-files 2025](https://www.wiz.io/blog/new-github-action-supply-chain-attack-reviewdog-tj-actions)
