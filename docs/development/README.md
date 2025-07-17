# 👩‍💻 Developer Guide

Welcome, contributor! This doc helps you set up a local Scain dev environment in **15 min or less**.

## 1. Prerequisites

| Tool | Version | Notes |
|------|---------|-------|
| Git | ≥ 2.34 | `brew install git` |
| Docker Desktop | ≥ 4.0 | Enable Kubernetes optional |
| Node.js | 18 LTS | Use `nvm install 18` |
| Go | 1.22 | `brew install go` |
| Python | 3.11 | `pyenv install 3.11` |
| ESP-IDF | v5.1 | `brew install --cask esp-idf` |

> **Tip:** macOS users can run `make setup` to auto-install missing deps via Homebrew.

## 2. Repository Workflow

```bash
# Fork & clone
 git clone git@github.com:Antimatters-Technology/Scain.git
 cd Scain

# Create feature branch
 git checkout -b feat/my-awesome-feature

# Run stack & dashboard
 make devnet        # starts docker-compose
 make dashboard     # hot-reload at localhost:3000

# Run unit tests
 make test

# Commit using Conventional Commits
 git commit -m "feat(api): add pagination to Fabric client"

# Push & open PR
 git push origin feat/my-awesome-feature
```

### Branch Rules
- `main` – protected, auto-deploys to staging.
- `dev` – integration branch (rebased weekly).
- `feat/*`, `fix/*`, `docs/*` – short-lived topic branches.

### Commit Emoji Cheatsheet
| Emoji | Type | Example |
|-------|------|---------|
| ✨ `feat` | New feature | `feat(dashboard): add dark mode` |
| 🐛 `fix` | Bug fix | `fix(chaincode): nil ptr deref on empty EPC` |
| 📝 `docs` | Docs | `docs(api): add GraphQL examples` |
| ♻️ `refactor` | Refactor | `refactor(firmware): extract sensor utils` |
| 🚀 `perf` | Performance | `perf(bridge): batch Fabric writes` |

## 3. Coding Standards

### Go (Chaincode)
- `go fmt ./...` before commit.
- Use dependency injection for external clients.

### C++ (Firmware)
- Follow **Google C++ Style Guide**.
- Prefer `std::chrono` over `delay()`.

### TypeScript (Dashboard)
- ESLint + Prettier enforced on commit.
- Use React hooks + Zustand for state.

### Python (Bridge)
- Black formatter + ruff linter.
- Retry network calls (`tenacity` or `retry` pkg).

## 4. Testing Strategy
- **Unit** – `go test`, `jest`, `pytest`.
- **Integration** – spin up `make devnet` and run Cypress (TODO).
- **Hardware-in-loop** – GitHub Actions self-hosted runner with ESP32.

## 5. Useful Make Targets
| Target | Purpose |
|--------|---------|
| `flash` | Build & flash firmware to `/dev/ttyUSB0` |
| `devnet-stop` | Bring down containers |
| `lint` | Run Go + TS linters |
| `simulate-sensors` | Publish fake MQTT data |

## 6. VS Code Tips
- Install **`ms-vscode.cpptools`**, **`esbenp.prettier-vscode`**, **`ms-python.python`**, **`golang.go`**.
- Add `.vscode/settings.json` (see `development/snippets`).

Happy hacking! 🚀 