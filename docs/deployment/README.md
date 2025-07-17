# 🚀 Deployment Guide

This guide explains how to deploy **Scain** in three scenarios:

1. **Local Laptop** – Developer workstation using Docker Desktop.
2. **Edge Gateway** – Fanless ARM/Linux box at the facility.
3. **Cloud Cluster** – Production Kubernetes or Swarm.

## 1. Local Laptop (Docker Compose)

```bash
# Clone repository
 git clone https://github.com/Antimatters-Technology/Scain.git && cd Scain

# Build images (first run)
make install           # installs Node, Go, Python deps

# Spin up the stack
make devnet            # docker-compose up -d

# Tail logs
make logs

# Tear down
make devnet-stop       # or make devnet-clean
```

### Ports
| Service | Port |
|---------|------|
| Dashboard | `http://localhost:3000` |
| Fabric REST | `http://localhost:4000` |
| OpenEPCIS | `http://localhost:8080` |
| Mosquitto MQTT | `localhost:1883` (TCP) / `ws://localhost:9001` (WebSockets) |
| CouchDB | `http://localhost:5984` |

## 2. Edge Gateway (ARM64)

```bash
# Build multi-arch images
DOCKER_BUILDKIT=1 docker buildx bake --set *.platform=linux/arm64 .

# Copy docker-compose.yml and edit volumes for persistent media
scp -r Scain/ edge-gw:/opt/scain
ssh edge-gw
cd /opt/scain && docker-compose up -d
```

- Use external USB SSD for `/opt/scain/data` to avoid SD-card wear.
- Enable TLS on Mosquitto & Fabric REST (see `config/certs/`).

## 3. Cloud (Kubernetes)

Helm charts are under `deployment/charts/` (TODO).

```bash
helm install scain deployment/charts/scain \
  --set dashboard.image.tag=$(git rev-parse --short HEAD)
```

### Storage Classes
- **CouchDB / Peer** – SSD or GP3
- **OpenEPCIS** – PostgreSQL operator (Crunchy / Bitnami)

## Environment Variables
| Name | Default | Description |
|------|---------|-------------|
| `FABRIC_ENDPOINT` | `fabric-rest-api:4000` | REST gateway for chaincode calls |
| `MQTT_BROKER` | `mosquitto:1883` | MQTT hostname:port |
| `OPENEPCIS_ENDPOINT` | `http://openepcis:8080` | EPCIS capture endpoint |

## Makefile Targets Reference
| Target | Function |
|--------|----------|
| `make setup` | Install dev dependencies |
| `make devnet` | Start Docker Compose stack |
| `make dashboard` | Run Next.js in hot-reload mode |
| `make flash` | Build & flash ESP32 firmware |
| `make test` | Run all unit tests |

---

> **Tip:** Always run `docker system prune -f` after heavy dev sessions to reclaim disk space. 