# 🚀 Scain Technical Action-Plan (Next 10 Days)

This living document tracks the short-term engineering tasks required to turn the current MVP into a **demo-ready, investor-friendly** product.

_Last updated: <!-- date placeholder -->_

---

## 0. Legend

| Symbol | Meaning |
|--------|---------|
| ✅ | Complete |
| ⏳ | In progress |
| 🔜 | Next up |
| 🛠  | Needs contributor |

---

## 1. Demo-Ready Stack

| Status | Task | Owner | Notes |
|--------|------|-------|-------|
| 🔜 | Enable TLS for Mosquitto, Fabric-REST, Dashboard | @devops | Use `scripts/generate-certs.sh`; mount into `docker-compose.yml`. |
| 🔜 | One-click bootstrap script (`scripts/bootstrap.sh`) | @devops | Detect deps, run `make setup devnet`. |
| 🔜 | Host live demo (Fly.io / Render) | @cloud | Build/push `scain/dashboard`, seed with `make simulate-sensors`. |

## 2. CI / CD

| Status | Task | Owner | Notes |
|--------|------|-------|-------|
| 🔜 | GitHub Action `ci.yml` – lint + `make test` | @maintainer | Block PR merge on failure. |
| 🔜 | Docker build-and-push on `main` | @maintainer | Tags: `dashboard:<sha>`, `mqtt-bridge:<sha>`. |

## 3. Pilot Hardware

| Status | Task | Owner | Notes |
|--------|------|-------|-------|
| 🔜 | Procure 5× BOM kits | @hardware | See `/docs/BOM.xlsx`. |
| 🔜 | Flash firmware (`make flash`) | @hardware | Verify publish to `scain/events`. |
| 🔜 | LOI with pilot customer | @bizdev | 30-day free trial, weekly feedback calls. |

## 4. Story & Visibility

| Status | Task | Owner | Notes |
|--------|------|-------|-------|
| 🔜 | Draft 10-slide deck (`/docs/pitch/`) | @founder | Problem → Solution → Demo → Ask. |
| 🔜 | Blog post: “Open-sourcing Scain” | @marketing | Cross-post LinkedIn + Hackaday. |
| 🔜 | Create 5 “good first issue” tickets | @maintainer | Label + small scope. |

---

## 5. Stretch (Day 30)

* LoRaWAN integration (Heltec) prototype.  
* Role-based auth via Keycloak.  
* Helm charts for EKS.

---

### How to Update This File
1. Edit via PR and update status symbols.  
2. Keep notes column short; link to issue numbers where relevant.  
3. Update the _Last updated_ line. 