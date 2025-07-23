# Hyperledger Fabric Blockchain Setup Guide

## Overview
This guide walks you through setting up a local Hyperledger Fabric network for the Scain supply chain traceability system. The blockchain provides immutable storage of EPCIS events and ensures data integrity.

## Prerequisites

### 1. System Requirements
- **OS:** Linux, macOS, or Windows (WSL2)
- **RAM:** Minimum 8GB, Recommended 16GB
- **Storage:** 10GB free space
- **Docker:** Version 20.10 or later
- **Docker Compose:** Version 2.0 or later
- **Go:** Version 1.19 or later
- **Node.js:** Version 16 or later (for Fabric tools)

### 2. Install Dependencies

**macOS:**
```bash
# Install Homebrew if not installed
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install dependencies
brew install docker docker-compose go node

# Start Docker Desktop
open /Applications/Docker.app
```

**Ubuntu/Debian:**
```bash
# Update package list
sudo apt update

# Install Docker
sudo apt install docker.io docker-compose

# Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# Add user to docker group
sudo usermod -aG docker $USER
newgrp docker
```

## Step 1: Download Fabric Binaries

```bash
# Create a directory for Fabric
mkdir -p ~/fabric-samples
cd ~/fabric-samples

# Download Fabric scripts
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.4.0 1.5.0

# Add Fabric binaries to PATH
export PATH=$PATH:~/fabric-samples/bin
echo 'export PATH=$PATH:~/fabric-samples/bin' >> ~/.bashrc
```

## Step 2: Set Up Network Configuration

```bash
# Navigate to Scain blockchain directory
cd /path/to/Scain/blockchain

# Create network configuration
mkdir -p network/config
```

Create `network/config/connection-profile.yaml`:
```yaml
name: "scain-network"
version: "1.0.0"
client:
  organization: Org1
  connection:
    timeout:
      peer:
        endorser: '300'
      orderer: '300'
organizations:
  Org1:
    mspid: Org1MSP
    peers:
      - peer0.org1.example.com
    certificateAuthorities:
      - ca.org1.example.com
peers:
  peer0.org1.example.com:
    url: grpcs://localhost:7051
    tlsCACerts:
      path: network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    grpcOptions:
      ssl-target-name-override: peer0.org1.example.com
      hostnameOverride: peer0.org1.example.com
certificateAuthorities:
  ca.org1.example.com:
    url: https://localhost:7054
    caName: ca-org1
    tlsCACerts:
      path: network/organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem
    registrar:
      - enrollId: admin
        enrollSecret: adminpw
orderers:
  orderer.example.com:
    url: grpcs://localhost:7050
    tlsCACerts:
      path: network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
    grpcOptions:
      ssl-target-name-override: orderer.example.com
      hostnameOverride: orderer.example.com
channels:
  scainchannel:
    orderers:
      - orderer.example.com
    peers:
      peer0.org1.example.com:
        endorsingPeer: true
        chaincodeQuery: true
        ledgerQuery: true
        eventSource: true
```

## Step 3: Start the Fabric Network

```bash
# Navigate to network directory
cd network

# Make setup script executable
chmod +x setup.sh

# Start the network
./setup.sh
```

The setup script will:
1. Generate cryptographic materials
2. Create the network configuration
3. Start Docker containers
4. Create and join the channel
5. Install and instantiate the chaincode

## Step 4: Create Wallet and Identity

```bash
# Create wallet directory
mkdir -p wallet

# Generate user identity
fabric-ca-client enroll -u https://admin:adminpw@localhost:7054 --caname ca-org1 -M wallet/admin

# Register a new user
fabric-ca-client register --caname ca-org1 --id.name appUser --id.secret appUserPw --id.type client --tls.certfiles network/organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem

# Enroll the user
fabric-ca-client enroll -u https://appUser:appUserPw@localhost:7054 --caname ca-org1 -M wallet/appUser
```

## Step 5: Configure Backend Environment

Create `.env` file in the backend directory:
```bash
# Blockchain Configuration
ENABLE_BLOCKCHAIN=true
FABRIC_CCP_PATH=./blockchain/network/config/connection-profile.yaml
FABRIC_WALLET_PATH=./wallet
FABRIC_USER_ID=appUser
FABRIC_CHANNEL_NAME=scainchannel
FABRIC_CHAINCODE_NAME=scain-chaincode

# Database Configuration
DATABASE_PATH=./scain.db

# Server Configuration
PORT=8081
NODE_ENV=development
```

## Step 6: Test Blockchain Integration

### 1. Start the Backend
```bash
cd backend
go run main.go
```

### 2. Test Event Creation with Blockchain
```bash
# Create a test event
curl -X POST http://localhost:8081/api/ingest \
  -H "Content-Type: application/json" \
  -d '{
    "deviceType": "ESP32",
    "deviceId": "test-device-001",
    "timestamp": "2024-07-21T12:00:00Z",
    "lotCode": "LOT123456",
    "data": {
      "temperature": 24.5,
      "humidity": 62.3
    }
  }'
```

### 3. Verify Blockchain Transaction
```bash
# Query the chaincode
peer chaincode query -C scainchannel -n scain-chaincode -c '{"Args":["GetEvent","EVENT_ID"]}'

# Check transaction history
peer chaincode query -C scainchannel -n scain-chaincode -c '{"Args":["GetAllEvents"]}'
```

## Step 7: Monitor the Network

### 1. Check Container Status
```bash
docker ps
```

Expected containers:
- `peer0.org1.example.com`
- `orderer.example.com`
- `ca.org1.example.com`
- `couchdb0`

### 2. View Logs
```bash
# Peer logs
docker logs peer0.org1.example.com

# Orderer logs
docker logs orderer.example.com

# CA logs
docker logs ca.org1.example.com
```

### 3. Access CouchDB (Optional)
```bash
# CouchDB is available at http://localhost:5984
# Username: admin
# Password: adminpw
```

## Step 8: Network Management

### Start Network
```bash
cd network
./start.sh
```

### Stop Network
```bash
cd network
./stop.sh
```

### Reset Network
```bash
cd network
./reset.sh
```

## Troubleshooting

### Common Issues

1. **Port Already in Use**
   ```bash
   # Check what's using the port
   lsof -i :7050
   lsof -i :7051
   lsof -i :7054
   
   # Kill the process
   kill -9 <PID>
   ```

2. **Docker Permission Issues**
   ```bash
   # Add user to docker group
   sudo usermod -aG docker $USER
   newgrp docker
   ```

3. **Certificate Issues**
   ```bash
   # Regenerate certificates
   cd network
   ./reset.sh
   ./setup.sh
   ```

4. **Chaincode Installation Fails**
   ```bash
   # Check chaincode logs
   docker logs dev-peer0.org1.example.com-scain-chaincode-1.0
   
   # Reinstall chaincode
   peer chaincode install -n scain-chaincode -v 1.0 -p github.com/scain-chaincode
   ```

### Debug Commands

```bash
# Check network status
peer channel list

# Check installed chaincodes
peer chaincode list --installed

# Check instantiated chaincodes
peer chaincode list --instantiated -C scainchannel

# Query channel info
peer channel getinfo -c scainchannel
```

## Production Considerations

### 1. Security
- Use proper TLS certificates
- Implement proper access control
- Secure wallet storage
- Use production-grade databases

### 2. Performance
- Configure proper resource limits
- Use production Docker images
- Implement connection pooling
- Monitor resource usage

### 3. Monitoring
- Set up logging aggregation
- Monitor blockchain metrics
- Implement health checks
- Set up alerting

### 4. Backup
- Regular wallet backups
- Database backups
- Configuration backups
- Disaster recovery plan

## Next Steps

1. **Test with Real Data:** Send data from actual IoT devices
2. **Add More Organizations:** Set up multi-org network
3. **Implement Private Data:** Use Fabric private data collections
4. **Add Monitoring:** Set up Prometheus/Grafana
5. **Security Hardening:** Implement proper authentication/authorization

## Resources

- [Hyperledger Fabric Documentation](https://hyperledger-fabric.readthedocs.io/)
- [Fabric Samples](https://github.com/hyperledger/fabric-samples)
- [Fabric Gateway Documentation](https://hyperledger.github.io/fabric-gateway/)
- [Scain Project Documentation](./README.md) 