# Scain Makefile
# Provides targets for firmware flashing, development network setup, and dashboard

.PHONY: help flash devnet dashboard dev build test clean install setup

# Default target
help:
	@echo "Scain Food Traceability MVP"
	@echo "Available targets:"
	@echo "  dev       - Start development server (frontend)"
	@echo "  build     - Build frontend for production"
	@echo "  flash     - Flash ESP32 firmware"
	@echo "  devnet    - Start development network (Docker)"
	@echo "  dashboard - Start Next.js dashboard"
	@echo "  test      - Run all tests"
	@echo "  test-go   - Run Go tests only"
	@echo "  clean     - Clean build artifacts"
	@echo "  install   - Install dependencies"
	@echo "  setup     - Initial project setup"

# Frontend targets
dev:
	@echo "Starting development server..."
	cd frontend && npm run dev

build:
	@echo "Building frontend for production..."
	cd frontend && npm run build

# ESP32 Firmware targets
flash:
	@echo "Flashing ESP32 firmware..."
	cd src/firmware && \
	idf.py set-target esp32 && \
	idf.py build && \
	idf.py -p /dev/ttyUSB0 flash monitor

flash-ota:
	@echo "OTA firmware update..."
	cd src/firmware && \
	idf.py build && \
	python scripts/ota_update.py --firmware build/main.bin

# Development network
devnet:
	@echo "Starting development network..."
	docker-compose down
	docker-compose up -d
	@echo "Waiting for services to start..."
	sleep 30
	@echo "Services available at:"
	@echo "  Dashboard: http://localhost:3000"
	@echo "  Fabric API: http://localhost:4000"
	@echo "  OpenEPCIS: http://localhost:8080"
	@echo "  CouchDB: http://localhost:5984"

devnet-logs:
	docker-compose logs -f

devnet-stop:
	docker-compose down

devnet-clean:
	docker-compose down -v
	docker system prune -f

# Dashboard targets
dashboard:
	@echo "Starting Next.js dashboard..."
	cd frontend && npm run dev

dashboard-build:
	@echo "Building dashboard for production..."
	cd frontend && npm run build

dashboard-test:
	@echo "Running dashboard tests..."
	cd frontend && npm test

# Chaincode targets
chaincode-package:
	@echo "Packaging chaincode..."
	cd chaincode && \
	GO111MODULE=on go mod vendor && \
	peer lifecycle chaincode package knowgraph-cc.tar.gz \
		--path . \
		--lang golang \
		--label knowgraph-cc_1.0

chaincode-install:
	@echo "Installing chaincode..."
	peer lifecycle chaincode install chaincode/knowgraph-cc.tar.gz

chaincode-deploy:
	@echo "Deploying chaincode..."
	scripts/deploy-chaincode.sh

# Test targets
test:
	@echo "Running all tests..."
	make test-chaincode
	make test-dashboard
	make test-firmware

test-go:
	@echo "Running Go tests..."
	cd chaincode && go test -v ./...
	cd src/tests && go test -v ./...

test-chaincode:
	@echo "Testing chaincode..."
	cd chaincode && go test -v ./...

test-dashboard:
	@echo "Testing dashboard..."
	cd frontend && npm test

test-firmware:
	@echo "Testing firmware..."
	cd src/firmware && idf.py build
	@echo "Firmware build successful"

# Setup and installation
setup:
	@echo "Setting up KnowGraph development environment..."
	make install
	make setup-fabric
	make setup-crypto

install:
	@echo "Installing dependencies..."
	# Install Node.js dependencies
	cd frontend && npm install
	# Install Go dependencies
	cd chaincode && go mod tidy
	# Install Python dependencies for scripts
	pip install -r requirements.txt

setup-fabric:
	@echo "Setting up Hyperledger Fabric..."
	mkdir -p config
	scripts/generate-crypto.sh
	scripts/generate-genesis.sh

setup-crypto:
	@echo "Generating crypto materials..."
	scripts/generate-crypto.sh

# Configuration targets
config-mosquitto:
	@echo "Configuring Mosquitto MQTT broker..."
	mkdir -p config
	cp config/mosquitto.conf.template config/mosquitto.conf
	mosquitto_passwd -c config/mosquitto.passwd knowgraph

config-fabric:
	@echo "Configuring Fabric network..."
	scripts/configure-fabric.sh

# Monitoring and maintenance
monitor:
	@echo "Starting monitoring dashboard..."
	docker-compose -f docker-compose.monitoring.yml up -d

logs:
	@echo "Showing system logs..."
	docker-compose logs -f --tail=100

status:
	@echo "System status:"
	docker-compose ps
	@echo ""
	@echo "Fabric peer status:"
	docker exec knowgraph-peer0 peer node status

# Clean targets
clean:
	@echo "Cleaning build artifacts..."
	cd src/firmware && idf.py clean
	cd frontend && rm -rf .next node_modules/.cache
	cd chaincode && go clean
	rm -rf build/
	rm -rf dist/

clean-docker:
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	docker system prune -f
	docker volume prune -f

# Production targets
deploy-prod:
	@echo "Deploying to production..."
	docker-compose -f docker-compose.prod.yml up -d

backup:
	@echo "Creating backup..."
	scripts/backup.sh

restore:
	@echo "Restoring from backup..."
	scripts/restore.sh

# Utility targets
generate-certs:
	@echo "Generating SSL certificates..."
	scripts/generate-certs.sh

update-deps:
	@echo "Updating dependencies..."
	cd frontend && npm update
	cd chaincode && go get -u ./...

lint:
	@echo "Running linters..."
	cd frontend && npm run lint
	cd chaincode && golangci-lint run

# Development utilities
dev-reset:
	@echo "Resetting development environment..."
	make devnet-clean
	make clean
	make setup
	make devnet

simulate-sensors:
	@echo "Starting sensor simulation..."
	python scripts/simulate_sensors.py

# Documentation
docs:
	@echo "Generating documentation..."
	cd frontend && npm run build-docs
	cd chaincode && godoc -http=:6060

# Variables
ESP32_PORT ?= /dev/ttyUSB0
FABRIC_VERSION ?= 2.5
NODE_VERSION ?= 18 