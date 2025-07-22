# Scain Backend

A production-ready Go backend implementing EPCIS 2.0 standards with Hyperledger Fabric blockchain integration for supply chain traceability.

## 🚀 Features

- **EPCIS 2.0 Compliance**: Full implementation of Electronic Product Code Information Services
- **Blockchain Integration**: Hyperledger Fabric support for immutable event anchoring
- **Device Management**: Registration, claiming, and heartbeat monitoring
- **Data Transformation**: Automatic conversion of raw sensor data to EPCIS events
- **Cryptographic Integrity**: SHA-256 hashing for all events
- **RESTful API**: Comprehensive endpoints with validation middleware
- **Real-time Processing**: Automatic data ingestion and transformation

## 🏗️ Architecture

```
backend/
├── main.go                 # Server entry point and routing
├── services/               # Business logic layer
│   ├── epcis_service.go   # EPCIS event processing
│   ├── device_service.go  # Device management
│   └── blockchain_service.go # Fabric blockchain integration
├── models/                 # Data models and structures
│   └── epcis.go           # EPCIS event models
├── database/               # Data persistence layer
│   └── database.go        # SQLite operations
├── middleware/             # HTTP middleware
│   └── validation.go      # Request validation
├── utils/                  # Utility functions
│   ├── hash.go            # Cryptographic utilities
│   └── canonical.go       # JSON canonicalization
└── admin/                  # Administrative tools
    └── generate_claim_codes.go # Device claim code generation
```

## 🛠 Setup & Installation

### Prerequisites
- Go 1.21+
- SQLite3
- Docker (for blockchain)

### Quick Start

```bash
# Clone and navigate
cd backend

# Install dependencies
go mod tidy

# Copy environment configuration
cp .env.example .env

# Run the server
go run .

# Or build and run
go build -o scain-backend
./scain-backend
```

### Environment Configuration

```bash
# Server
PORT=8081
NODE_ENV=development

# Database
DATABASE_PATH=./scain.db

# Blockchain (Optional)
ENABLE_BLOCKCHAIN=false
FABRIC_CCP_PATH=../blockchain/network/connection-profile.yaml
FABRIC_WALLET_PATH=../blockchain/network/wallet
FABRIC_USER_ID=appUser
FABRIC_CHANNEL_NAME=mychannel
FABRIC_CHAINCODE_NAME=scain
```

## 🔗 API Endpoints

### Health & Info
- `GET /health` - Health check
- `GET /api` - API information

### EPCIS Events
- `POST /api/events` - Create EPCIS event (+ blockchain anchoring)
- `GET /api/events/:id` - Retrieve event by ID
- `POST /api/ingest` - Ingest raw sensor data

### Device Management
- `POST /api/devices` - Register device
- `GET /api/devices/:id` - Get device info
- `POST /api/claim` - Claim device with code

### Blockchain (when enabled)
- `GET /api/events/:id/verify` - Verify event on blockchain
- `GET /api/events/:id/history` - Get blockchain transaction history

## 🔄 Data Flow

1. **Device Registration**: Devices register via `/api/devices`
2. **Data Ingestion**: Raw data comes via `/api/ingest`
3. **EPCIS Transformation**: Raw data converted to EPCIS events
4. **Database Storage**: Events stored in SQLite with hash
5. **Blockchain Anchoring**: Events submitted to Fabric (if enabled)
6. **API Access**: Events accessible via REST endpoints

## ⛓️ Blockchain Integration

### How It Works
1. When `ENABLE_BLOCKCHAIN=true`, the backend initializes Fabric SDK connection
2. After storing events in database, they're automatically submitted to blockchain
3. Event hash and metadata stored on-chain for immutability
4. Blockchain transaction ID stored in database for traceability

### Blockchain Service API
```go
// Submit event to blockchain
record, err := blockchainService.SubmitEvent(eventData)

// Retrieve event from blockchain
record, err := blockchainService.GetEvent(eventID)

// Verify event integrity
isValid, err := blockchainService.VerifyEvent(eventID, eventData)
```

## 🧪 Testing

```bash
# Run tests
go test ./...

# Test with sample data
./test_api.sh

# Generate claim codes
go run admin/generate_claim_codes.go
```

## 📊 Database Schema

### Events Table
- `id` - Primary key
- `event_type` - EPCIS event type
- `event_time` - Event timestamp
- `hash` - SHA-256 hash of event data
- `raw_data` - JSON event data
- `blockchain_tx_id` - Fabric transaction ID (optional)
- `device_id` - Source device ID
- `lot_code` - Product lot identifier

### Devices Table
- `id` - Primary key
- `device_id` - Unique device identifier
- `device_type` - Device type (ESP32, etc.)
- `claim_code` - Device claiming code
- `is_claimed` - Claim status
- `last_heartbeat` - Last communication

## 🔐 Security Features

- **Input Validation**: All requests validated with go-playground/validator
- **Content-Type Enforcement**: Prevents content confusion attacks
- **Request Size Limiting**: Prevents DoS attacks
- **Cryptographic Hashing**: SHA-256 for data integrity
- **Blockchain Immutability**: Tamper-proof event records

## 📈 Performance

- **Concurrent Processing**: Goroutines for blockchain operations
- **Database Optimization**: Indexed queries and prepared statements
- **Error Handling**: Graceful degradation for blockchain failures
- **Logging**: Structured logging with logrus

## 🚀 Production Deployment

```bash
# Build optimized binary
CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o scain-backend

# Docker deployment
docker build -t scain-backend .
docker run -p 8081:8081 -v $(pwd)/data:/app/data scain-backend
```

## 🤝 Contributing

1. Follow Go conventions and formatting (`go fmt`)
2. Add tests for new features
3. Update API documentation
4. Ensure blockchain integration works with/without Fabric

## 📚 Dependencies

- **gin-gonic/gin**: HTTP web framework
- **go-playground/validator**: Request validation
- **hyperledger/fabric-sdk-go**: Blockchain integration
- **mattn/go-sqlite3**: Database driver
- **sirupsen/logrus**: Structured logging

---

**Status: ✅ Production Ready with Blockchain Integration** 