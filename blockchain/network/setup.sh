#!/bin/bash

# Scain Hyperledger Fabric Network Setup Script
# This script sets up a basic Fabric network for development and testing

set -e

# Configuration
FABRIC_VERSION="2.5.0"
CA_VERSION="1.5.0"
NETWORK_NAME="scain-network"
CHANNEL_NAME="mychannel"
CHAINCODE_NAME="scain"

echo "🚀 Setting up Scain Hyperledger Fabric Network..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Download Fabric binaries if not present
if [ ! -d "bin" ]; then
    echo "📥 Downloading Hyperledger Fabric binaries..."
    curl -sSL https://bit.ly/2ysbOFE | bash -s -- ${FABRIC_VERSION} ${CA_VERSION}
fi

# Create network directory structure
mkdir -p organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com
mkdir -p organizations/ordererOrganizations/example.com/orderers/orderer.example.com
mkdir -p channel-artifacts
mkdir -p wallet

# Generate crypto material
echo "🔐 Generating crypto material..."
./bin/cryptogen generate --config=./crypto-config.yaml --output="organizations"

# Generate genesis block
echo "🏗️ Generating genesis block..."
export FABRIC_CFG_PATH=$PWD
./bin/configtxgen -profile TwoOrgsOrdererGenesis -channelID system-channel -outputBlock ./channel-artifacts/genesis.block

# Generate channel configuration transaction
echo "📋 Generating channel configuration..."
./bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/${CHANNEL_NAME}.tx -channelID ${CHANNEL_NAME}

# Generate anchor peer transactions
./bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID ${CHANNEL_NAME} -asOrg Org1MSP

echo "✅ Fabric network setup completed!"
echo "📝 Next steps:"
echo "   1. Run './start.sh' to start the network"
echo "   2. Run './deploy-chaincode.sh' to deploy the Scain chaincode"
echo "   3. Configure your backend with the generated connection profile" 