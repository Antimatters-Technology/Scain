# Scain - Supply Chain Traceability Platform

A comprehensive supply chain traceability platform implementing EPCIS 2.0 standards with **blockchain integration** and **IoT device support** for immutable audit trails.

## 🚀 Features

### Backend (Go) ✅ **FULLY FUNCTIONAL**
- **EPCIS 2.0 Compliance**: Full implementation of Electronic Product Code Information Services standard
- **Multi-Device Support**: ESP32, AWS IoT ExpressLink, LoRaWAN, GPS Trackers, ERP integration
- **Automatic Data Transformation**: Raw sensor data → EPCIS events
- **Device Management**: Registration, claiming, and heartbeat monitoring  
- **SQLite Database**: Persistent storage with automatic migrations
- **Cryptographic Integrity**: SHA-256 hashing for event verification
- **RESTful API**: Comprehensive endpoints with validation middleware
- **Real-time Processing**: Automatic raw data ingestion and transformation
- **🆕 Blockchain Integration**: Hyperledger Fabric support for immutable event anchoring

### Frontend (Next.js)
- **Modern Dashboard**: Supply chain visibility and analytics
- **Device Monitoring**: Real-time sensor data visualization  
- **Lot Traceability**: Complete product journey tracking
- **Responsive Design**: Mobile-first approach with Tailwind CSS
- **Admin Interface**: Device and user management
- **🆕 Blockchain Verification**: View blockchain proof for events

### IoT Integration ✅ **ESP32 FIRMWARE INCLUDED**
- **Multi-Protocol Support**: MQTT, HTTP, LoRaWAN, Cellular
- **Edge Processing**: Local data processing capabilities
- **Secure Boot**: Hardware-based security foundation
- **OTA Updates**: Remote firmware management
- **🆕 Reference Firmware**: Ready-to-use ESP32 code with sensor support

### Blockchain (Hyperledger Fabric) 🆕 **READY FOR DEPLOYMENT**
- **Immutable Audit Trail**: All events anchored on blockchain
- **Smart Contracts**: Custom chaincode for EPCIS event management
- **Event Verification**: Cryptographic proof of data integrity
- **Transaction History**: Complete audit trail for compliance

## 📦 Project Structure

```
Scain/
├── backend/                 # Go backend server (FULLY FUNCTIONAL + BLOCKCHAIN)
│   ├── main.go             # Server entry point
│   ├── database/           # SQLite models and migrations
│   ├── services/           # Business logic (EPCIS, devices, blockchain)
│   ├── middleware/         # HTTP validation and security
│   ├── models/             # EPCIS data structures
│   └── utils/              # Cryptographic utilities
├── frontend/               # Next.js web application
│   ├── app/                # Next.js 14 App Router
│   ├── components/         # Reusable UI components
│   └── types/              # TypeScript definitions
├── firmware/               # 🆕 ESP32 firmware sources
│   ├── main.cpp            # Main firmware code
│   ├── config.h            # Configuration header
│   ├── platformio.ini      # PlatformIO project config
│   └── README.md           # Build and flash instructions
├── blockchain/             # 🆕 Hyperledger Fabric integration
│   ├── chaincode/          # Smart contracts
│   │   ├── scain_chaincode.go  # Main chaincode
│   │   └── go.mod          # Chaincode dependencies
│   ├── network/            # Network setup scripts
│   └── README.md           # Deployment guide
├── docs/                   # Documentation
└── package.json            # Root package configuration
```

## 🛠 Quick Start

### Prerequisites
- **Go 1.21+** (for backend)
- **Node.js 18+** (for frontend)
- **npm/yarn** (package manager)
- **Docker** (for Hyperledger Fabric - optional)
- **PlatformIO** (for ESP32 firmware - optional)

### Backend Setup (Ready to Use!)

```bash
# Navigate to backend
cd backend

# Install dependencies
go mod tidy

# Start the server
go run .

# Or using npm from root
npm run dev:backend
```

The Go backend will start at `http://localhost:8081` with full EPCIS functionality.

### Frontend Setup

```bash
# Navigate to frontend  
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

The frontend will be available at `http://localhost:3000`.

### 🆕 ESP32 Firmware Setup

```bash
# Navigate to firmware directory
cd firmware

# Install PlatformIO (if not installed)
pip install platformio

# Configure WiFi and API endpoint in config.h
# Build and upload to ESP32
pio run --target upload
```

### 🆕 Blockchain Setup (Optional)

```bash
# Navigate to blockchain network directory
cd blockchain/network

# Set up Fabric network
./setup.sh

# Start the network
./start.sh

# Deploy chaincode
./deploy-chaincode.sh

# Enable blockchain in backend
export ENABLE_BLOCKCHAIN=true
```

### Full Stack Development

```bash
# Install all dependencies
npm install

# Start both backend and frontend
npm run dev

# With blockchain enabled
export ENABLE_BLOCKCHAIN=true && npm run dev
```

## 🔗 API Endpoints (Backend)

### Core EPCIS Operations
- `POST /api/events` - Create EPCIS events (+ blockchain anchoring)
- `GET /api/events/:id` - Retrieve events by ID
- `POST /api/ingest` - Ingest raw sensor data (auto-transforms to EPCIS)

### Device Management  
- `POST /api/devices` - Register new devices
- `GET /api/devices/:id` - Get device information
- `POST /api/claim` - Claim devices with validation codes

### 🆕 Blockchain Operations
- `GET /api/events/:id/verify` - Verify event on blockchain
- `GET /api/events/:id/history` - Get blockchain transaction history

### System Health
- `GET /health` - Health check endpoint
- `GET /api` - API documentation and endpoints

## 📊 Data Flow

1. **IoT Devices** (ESP32) → Raw sensor data via HTTP/MQTT
2. **Backend Ingestion** → `/api/ingest` endpoint  
3. **Automatic Transformation** → Raw data → EPCIS events
4. **Database Storage** → SQLite with integrity hashing
5. **🆕 Blockchain Anchoring** → Event hash/data stored on Fabric
6. **Frontend Visualization** → Real-time dashboard updates
7. **Traceability Queries** → Complete supply chain visibility + blockchain proof

## 🔐 Security Features

- **Input Validation**: Comprehensive request validation with detailed error messages
- **Content-Type Enforcement**: API security middleware
- **Request Size Limiting**: Prevention of abuse attacks
- **Cryptographic Integrity**: SHA-256 hashing for all events
- **Secure Device Claiming**: Validation code system
- **🆕 Blockchain Immutability**: Tamper-proof event records

## 🎯 Supported Device Types

- **ESP32**: General IoT sensor platforms (firmware included)
- **AWS IoT ExpressLink**: Cloud-connected modules  
- **LoRaWAN**: Long-range, low-power devices
- **GPS Trackers**: Location-based tracking
- **ERP Systems**: Enterprise integration

## 📋 EPCIS Event Types

- **ObjectEvent**: Individual product tracking
- **TransformationEvent**: Manufacturing processes
- **AggregationEvent**: Packaging and grouping
- **TransactionEvent**: Ownership transfers

## 🚀 Deployment

### Docker Support
```bash
# Backend
cd backend
docker build -t scain-backend .
docker run -p 8081:8081 scain-backend

# Frontend
cd frontend  
docker build -t scain-frontend .
docker run -p 3000:3000 scain-frontend

# Fabric Network
cd blockchain/network
docker-compose up -d
```

### Environment Configuration
```bash
# Copy example environment file
cp .env.example .env

# Configure for your environment
# Set ENABLE_BLOCKCHAIN=true to enable Fabric integration
```

## 📚 Documentation

- [Backend API Documentation](./backend/README.md)
- [Frontend Development Guide](./frontend/README.md)
- [🆕 ESP32 Firmware Guide](./firmware/README.md)
- [🆕 Blockchain Integration Guide](./blockchain/README.md)
- [Architecture Overview](./docs/architecture/README.md)
- [Deployment Guide](./docs/deployment/README.md)
- [User Guide](./docs/user-guide/README.md)

## 🧪 Testing

```bash
# Backend tests
cd backend && go test ./...

# Frontend tests  
cd frontend && npm test

# End-to-end tests
npm run test:e2e

# ESP32 firmware (hardware required)
cd firmware && pio test

# Blockchain integration tests
cd blockchain && npm test
```

## 🔄 Development Status

| Component | Status | Features |
|-----------|--------|----------|
| **Go Backend** | ✅ **PRODUCTION READY** | Full EPCIS implementation, device management, data transformation, comprehensive testing, **blockchain integration** |
| **Database Layer** | ✅ **PRODUCTION READY** | SQLite with automatic migrations, all CRUD operations working |
| **API Endpoints** | ✅ **PRODUCTION READY** | RESTful API with validation, error handling, security middleware |
| **Admin Tools** | ✅ **Complete** | Claim code generation, comprehensive test suite |
| **🆕 ESP32 Firmware** | ✅ **READY** | Complete sensor integration, EPCIS event generation, configurable |
| **🆕 Blockchain (Fabric)** | ✅ **READY** | Chaincode deployed, backend integration, event anchoring |
| **Frontend Dashboard** | 🚧 In Progress | Basic structure implemented, blockchain verification pending |
| **Device Integration** | ✅ **Complete** | ESP32 firmware, protocol support framework |

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see [LICENSE](LICENSE) for details.

## 🆘 Support

- [Issue Tracker](https://github.com/your-org/scain/issues)
- [Discussions](https://github.com/your-org/scain/discussions)
- Email: support@scain.io

---

**Built with ❤️ for supply chain transparency and traceability**
**🆕 Now with end-to-end blockchain integration and IoT device support!**
