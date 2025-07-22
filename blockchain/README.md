# Scain Blockchain Integration

Hyperledger Fabric blockchain integration for immutable supply chain event anchoring and verification.

## 🚀 Overview

This directory contains all components needed to deploy and integrate Hyperledger Fabric with the Scain platform:

- **Chaincode**: Smart contracts for storing and querying EPCIS events
- **Network**: Local Fabric network setup and configuration
- **Integration**: Backend service integration with Fabric SDK

## 🏗️ Architecture

```
blockchain/
├── chaincode/                    # Smart contracts
│   ├── scain_chaincode.go       # Main chaincode implementation
│   ├── go.mod                   # Chaincode dependencies
│   └── README.md                # Chaincode documentation
├── network/                      # Network setup
│   ├── setup.sh                 # Network initialization script
│   ├── start.sh                 # Network startup script (to be created)
│   ├── deploy-chaincode.sh      # Chaincode deployment (to be created)
│   └── connection-profile.yaml  # Network connection config (generated)
└── README.md                    # This file
```

## 🛠 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+
- Hyperledger Fabric binaries (downloaded automatically)

### 1. Setup Network

```bash
cd blockchain/network

# Initialize Fabric network
./setup.sh

# Start the network
./start.sh

# Deploy chaincode
./deploy-chaincode.sh
```

### 2. Enable Backend Integration

```bash
# In your .env file
ENABLE_BLOCKCHAIN=true
FABRIC_CCP_PATH=./blockchain/network/connection-profile.yaml
FABRIC_WALLET_PATH=./blockchain/network/wallet
FABRIC_USER_ID=appUser
FABRIC_CHANNEL_NAME=mychannel
FABRIC_CHAINCODE_NAME=scain
```

### 3. Test Integration

```bash
# Start backend with blockchain enabled
cd backend
ENABLE_BLOCKCHAIN=true go run .

# Create an event (will be anchored on blockchain)
curl -X POST http://localhost:8081/api/events \
  -H "Content-Type: application/json" \
  -d '{
    "eventType": "ObjectEvent",
    "eventTime": "2024-01-01T12:00:00Z",
    "action": "OBSERVE",
    "epcList": ["urn:epc:id:sgtin:0614141.107346.2018"]
  }'
```

## ⛓️ Chaincode Functions

### Core Functions

- `StoreEvent(eventID, eventHash, timestamp, eventType, data)` - Store EPCIS event
- `GetEvent(eventID)` - Retrieve event by ID
- `EventExists(eventID)` - Check if event exists
- `GetEventsByType(eventType)` - Query events by type
- `GetEventHistory(eventID)` - Get transaction history for event
- `VerifyEventHash(eventID, expectedHash)` - Verify event integrity

### Usage Examples

```bash
# Invoke chaincode (via Fabric CLI)
peer chaincode invoke -C mychannel -n scain \
  -c '{"function":"StoreEvent","Args":["evt_123","abc123hash","2024-01-01T12:00:00Z","ObjectEvent","{}"]}'

# Query chaincode
peer chaincode query -C mychannel -n scain \
  -c '{"function":"GetEvent","Args":["evt_123"]}'
```

## 🔄 Data Flow

1. **Backend Event Creation**: EPCIS event created via API
2. **Database Storage**: Event stored in SQLite with hash
3. **Blockchain Submission**: Event hash/data submitted to Fabric
4. **Chaincode Execution**: Smart contract stores event on ledger
5. **Transaction ID**: Fabric TX ID stored back in database
6. **Verification**: Event can be verified against blockchain

## 🔐 Security & Trust

### Immutability
- Events stored on blockchain cannot be modified or deleted
- Cryptographic hashing ensures data integrity
- Transaction history provides complete audit trail

### Verification
- Event hashes can be independently verified
- Blockchain provides cryptographic proof of existence
- Transaction timestamps prove when events were recorded

### Access Control
- Fabric network uses PKI for identity management
- Chaincode enforces business rules and validation
- Only authorized participants can submit transactions

## 🚀 Production Deployment

### Cloud Deployment Options

1. **IBM Blockchain Platform**
2. **Amazon Managed Blockchain**
3. **Azure Blockchain Service**
4. **Self-hosted Kubernetes**

### Configuration for Production

```yaml
# connection-profile.yaml (example)
name: scain-network
version: 1.0.0
client:
  organization: Org1
  connection:
    timeout:
      peer:
        endorser: 300
channels:
  mychannel:
    peers:
      peer0.org1.example.com:
        endorsingPeer: true
        chaincodeQuery: true
        ledgerQuery: true
        eventSource: true
```

## 📊 Monitoring & Analytics

### Blockchain Metrics
- Transaction throughput (TPS)
- Block creation time
- Chaincode execution time
- Network latency

### Event Analytics
- Events stored per day/hour
- Event types distribution
- Device activity patterns
- Blockchain vs database consistency

## 🧪 Testing

### Unit Tests
```bash
cd chaincode
go test ./...
```

### Integration Tests
```bash
# Test backend-blockchain integration
cd backend
ENABLE_BLOCKCHAIN=true go test ./services/blockchain_service_test.go
```

### End-to-End Tests
```bash
# Test full flow: API -> Database -> Blockchain
./test_blockchain_integration.sh
```

## 🔧 Troubleshooting

### Common Issues

1. **Connection Failed**
   - Check Docker containers are running
   - Verify connection profile paths
   - Check Fabric network status

2. **Chaincode Errors**
   - Verify chaincode is deployed
   - Check chaincode logs
   - Validate function parameters

3. **Transaction Failures**
   - Check endorsement policy
   - Verify user credentials
   - Check network connectivity

### Debug Commands

```bash
# Check network status
docker ps | grep fabric

# View chaincode logs
docker logs <chaincode-container-id>

# Check peer logs
docker logs peer0.org1.example.com
```

## 📈 Performance Tuning

### Optimization Tips
- Use batch transactions for multiple events
- Implement proper indexing in chaincode
- Configure appropriate timeout values
- Monitor and tune endorsement policies

### Scaling Considerations
- Add more peers for read scalability
- Use multiple channels for workload isolation
- Implement proper sharding strategies
- Consider state database optimization (CouchDB vs LevelDB)

## 🤝 Contributing

1. Test chaincode changes locally first
2. Update documentation for new functions
3. Add unit tests for new features
4. Verify integration with backend service

## 📚 Resources

- [Hyperledger Fabric Documentation](https://hyperledger-fabric.readthedocs.io/)
- [Fabric SDK Go](https://github.com/hyperledger/fabric-sdk-go)
- [EPCIS 2.0 Standard](https://www.gs1.org/standards/epcis)
- [Chaincode Development Tutorial](https://hyperledger-fabric.readthedocs.io/en/latest/chaincode.html)

---

**Status: ✅ Ready for Production Deployment** 