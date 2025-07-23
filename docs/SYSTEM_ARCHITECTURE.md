# Scain System Architecture

## Overview
Scain is a supply chain traceability system that combines IoT devices, blockchain technology, and EPCIS standards to provide end-to-end visibility of products through the supply chain.

## System Architecture Diagram

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   IoT Devices   │    │   Edge Devices  │    │   ERP Systems   │
│                 │    │                 │    │                 │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │   ESP32     │ │    │ │ ExpressLink │ │    │ │ SAP/ERP     │ │
│ │ Sensors     │ │    │ │ LoRaWAN     │ │    │ │ Integration │ │
│ │ GPS Tracker │ │    │ │ Gateway     │ │    │ │             │ │
│ └─────────────┘ │    │ └─────────────┘ │    │ └─────────────┘ │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────▼─────────────┐
                    │     Backend API           │
                    │   (Go/Gin Server)         │
                    │                           │
                    │ ┌───────────────────────┐ │
                    │ │   REST API Layer      │ │
                    │ │ • /api/events         │ │
                    │ │ • /api/devices        │ │
                    │ │ • /api/ingest         │ │
                    │ │ • /api/claim          │ │
                    │ └───────────────────────┘ │
                    │                           │
                    │ ┌───────────────────────┐ │
                    │ │   Business Logic      │ │
                    │ │ • EPCIS Processing    │ │
                    │ │ • Device Management   │ │
                    │ │ • Data Transformation │ │
                    │ └───────────────────────┘ │
                    │                           │
                    │ ┌───────────────────────┐ │
                    │ │   Blockchain Service  │ │
                    │ │ • Hyperledger Fabric  │ │
                    │ │ • Event Hashing       │ │
                    │ │ • Transaction Mgmt    │ │
                    │ └───────────────────────┘ │
                    └─────────────┬─────────────┘
                                  │
                    ┌─────────────▼─────────────┐
                    │      Database Layer       │
                    │                           │
                    │ ┌───────────────────────┐ │
                    │ │   SQLite Database     │ │
                    │ │ • Events              │ │
                    │ │ • Devices             │ │
                    │ │ • Raw Data            │ │
                    │ │ • Claim Codes         │ │
                    │ └───────────────────────┘ │
                    └─────────────┬─────────────┘
                                  │
                    ┌─────────────▼─────────────┐
                    │   Blockchain Network      │
                    │                           │
                    │ ┌───────────────────────┐ │
                    │ │ Hyperledger Fabric    │ │
                    │ │ • Peers               │ │
                    │ │ • Orderers            │ │
                    │ │ • Chaincode           │ │
                    │ └───────────────────────┘ │
                    └───────────────────────────┘
```

## Component Details

### 1. IoT Devices Layer
- **ESP32 Devices:** Temperature, humidity, pressure sensors
- **GPS Trackers:** Location tracking for shipments
- **ExpressLink:** AWS IoT integration
- **LoRaWAN:** Long-range wireless communication

### 2. Edge Devices Layer
- **ExpressLink Gateways:** AWS IoT connectivity
- **LoRaWAN Gateways:** Long-range data collection
- **Local Processing:** Data aggregation and filtering

### 3. Backend API Layer
- **Framework:** Go with Gin web framework
- **Port:** 8081 (configurable)
- **Features:**
  - RESTful API endpoints
  - EPCIS event processing
  - Device management
  - Blockchain integration
  - Data validation and transformation

### 4. Database Layer
- **Type:** SQLite (production-ready alternatives: PostgreSQL, MySQL)
- **Tables:**
  - `events`: EPCIS events with blockchain transaction IDs
  - `devices`: IoT device registry
  - `raw_data_ingestions`: Raw sensor data
  - `claim_code_entries`: Device claim codes

### 5. Blockchain Layer
- **Platform:** Hyperledger Fabric
- **Purpose:** Immutable event storage and audit trail
- **Components:**
  - Fabric peers for endorsement
  - Ordering service for consensus
  - Chaincode for smart contracts
  - Wallet for identity management

## Data Flow

### 1. Device Registration
```
Device → POST /api/devices → Database → Response
```

### 2. Data Ingestion
```
Device → POST /api/ingest → Raw Data Storage → EPCIS Transformation → Event Storage → Blockchain (optional)
```

### 3. Event Creation
```
Raw Data → EPCIS Processing → Database Storage → Blockchain Submission → Transaction ID Storage
```

### 4. Device Claiming
```
User → POST /api/claim → Claim Code Validation → Device Association → Database Update
```

## Security Features

### 1. Device Security
- Secure boot verification
- Firmware version tracking
- Device authentication via claim codes

### 2. Data Security
- SHA-256 hashing for data integrity
- Blockchain immutability
- Database transaction safety

### 3. API Security
- Input validation
- Error handling
- Rate limiting (recommended)

## Scalability Considerations

### 1. Horizontal Scaling
- Stateless API design
- Database connection pooling
- Load balancer ready

### 2. Vertical Scaling
- Efficient Go runtime
- Optimized database queries
- Memory-efficient data structures

### 3. Blockchain Scaling
- Fabric network expansion
- Channel partitioning
- Private data collections

## Monitoring and Observability

### 1. Health Checks
- `GET /health` endpoint
- Database connectivity
- Blockchain service status

### 2. Logging
- Structured logging with logrus
- Request/response tracking
- Error monitoring

### 3. Metrics
- API response times
- Database performance
- Blockchain transaction rates

## Deployment Options

### 1. Development
- Local SQLite database
- Mock blockchain service
- Single API instance

### 2. Production
- PostgreSQL/MySQL database
- Full Fabric network
- Multiple API instances
- Load balancer
- Monitoring stack 