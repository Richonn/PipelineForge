# DevSecOps Reference Pipeline — TODO

## Description du projet

**DevSecOps Reference Pipeline** est un dépôt GitHub open source qui modélise une chaîne de production applicative sécurisée et agile, de bout en bout. Il ne s'agit pas d'un simple pipeline CI/CD : c'est un **livrable de référence** — documenté, justifié, et reproductible — pensé comme un état de l'art que n'importe quelle équipe peut étudier, adapter, ou déployer.

Le projet couvre trois stacks applicatives (Go, Python, Node.js) et intègre des contrôles de sécurité à chaque étape du cycle de vie logiciel : du commit au déploiement en production. Chaque choix d'outil est documenté et justifié face aux risques qu'il couvre.

Il s'adresse autant aux équipes DevOps qui veulent sécuriser leurs pipelines, qu'aux consultants ou architectes qui cherchent une référence concrète et argumentée.

---

## Stack technique

- **CI/CD** : GitHub Actions
- **Conteneurisation** : Docker, multi-stage builds
- **Orchestration** : Kubernetes (optionnel, via KubeForge)
- **SAST** : Semgrep + CodeQL
- **Secrets scanning** : Gitleaks (pre-commit + CI)
- **SCA / CVE** : Trivy (images + dépendances)
- **DAST** : OWASP ZAP (smoke test sur environnement de staging)
- **IaC scanning** : Checkov (Dockerfile, manifests K8s)
- **Threat modeling** : fichier THREAT_MODEL.md structuré (STRIDE)
- **Reporting** : GitHub Security Tab + résumés automatiques en PR

---

## TODO

### Phase 1 — Fondations du dépôt
- [ ] Créer la structure du dépôt avec README principal
- [ ] Rédiger le `ARCHITECTURE.md` : schéma de la chaîne complète (draw.io ou Mermaid)
- [ ] Rédiger le `THREAT_MODEL.md` : modélisation des menaces sur la chaîne CI/CD elle-même (STRIDE)
  - Quels risques sur GitHub Actions ?
  - Quels risques sur le registry Docker ?
  - Quels risques sur les secrets dans la CI ?
  - Quels risques liés aux dépendances tierces ?
- [ ] Rédiger le `TOOL_CHOICES.md` : justification de chaque outil retenu vs alternatives (ex : pourquoi Semgrep plutôt que SonarQube ?)

### Phase 2 — Applications de démonstration
- [ ] Créer une app Go minimaliste (API REST simple)
- [ ] Créer une app Python minimaliste (Flask ou FastAPI)
- [ ] Créer une app Node.js minimaliste (Express)
- [ ] Chaque app doit contenir volontairement des vulnérabilités corrigibles pour démontrer les scans

### Phase 3 — Pipeline CI/CD sécurisé
- [ ] Pipeline de base : lint, build, test unitaires
- [ ] Intégrer Gitleaks en pre-commit (hook local) + étape CI
- [ ] Intégrer Semgrep (SAST) pour les trois stacks
- [ ] Intégrer CodeQL pour Go et JavaScript
- [ ] Intégrer Trivy pour le scan d'image Docker et des dépendances
- [ ] Intégrer Checkov pour les Dockerfiles et manifests K8s
- [ ] Intégrer OWASP ZAP en mode baseline sur un environnement de staging (DAST)
- [ ] Configurer les seuils de blocage : Critical/High bloquent le merge, Medium warn
- [ ] Publier les résultats dans le GitHub Security Tab (SARIF)
- [ ] Générer un résumé automatique de sécurité dans chaque PR

### Phase 4 — Gestion des secrets et supply chain
- [ ] Mettre en place la gestion des secrets via GitHub Secrets + documentation des bonnes pratiques
- [ ] Pinning des actions GitHub (`uses: actions/checkout@SHA` plutôt que `@v3`)
- [ ] Activer Dependabot pour les mises à jour de dépendances automatiques
- [ ] Documenter les risques supply chain dans le `THREAT_MODEL.md` (lien avec rapport CERT-Wavestone 2025 : 500+ paquets npm compromis)

### Phase 5 — Documentation état de l'art
- [ ] Rédiger `DEVSECOPS_REFERENCE.md` : guide complet de la chaîne
  - Shift Left : définition et mise en pratique
  - SAST vs DAST vs SCA : quand utiliser quoi
  - Gestion des faux positifs : stratégie et configuration
  - Pipeline agile vs pipeline sécurisé : comment concilier les deux
- [ ] Ajouter des badges dans le README (pipeline status, Trivy, Semgrep, licence)
- [ ] Créer un schéma visuel de la chaîne complète (Mermaid dans le README)
- [ ] Écrire un article de blog ou post LinkedIn à la publication (optionnel)

### Phase 6 — Polish & publication
- [ ] Vérifier que chaque outil est bien configuré et pas juste présent
- [ ] S'assurer que le pipeline tourne sans erreur sur les trois apps
- [ ] Relecture complète de la documentation
- [ ] Publier le dépôt public sur GitHub (compte Richonn)
- [ ] Soumettre sur r/devops et r/netsec pour feedback

---

## Ce que ce projet démontre

| Compétence | Preuve |
|---|---|
| Connaissance des outils CI/CD et de leurs risques | `THREAT_MODEL.md` + `TOOL_CHOICES.md` |
| Intégration d'outils de sécurité dans une chaîne | Pipeline GitHub Actions fonctionnel |
| Représentation d'une chaîne à l'état de l'art | `ARCHITECTURE.md` + `DEVSECOPS_REFERENCE.md` |
| Réalisation d'une maquette concrète | Les trois apps + pipeline qui tourne |
| Volet organisationnel | Documentation des processus et bonnes pratiques |

---

## Liens avec les projets existants

- **ShieldCI** → génère des pipelines ; ce projet est la version "référence manuelle et documentée" de ce que ShieldCI automatise
- **KubeForge** → couvre le runtime/CD ; ce projet couvre la CI et le build

Les trois ensemble forment une couverture DevSecOps complète du commit au cluster.
