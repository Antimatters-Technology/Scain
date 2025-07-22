# Scain Hyperledger Fabric Chaincode

This directory contains chaincode (smart contracts) for anchoring EPCIS events on a Hyperledger Fabric blockchain.

## Features
- Store and query EPCIS event hashes or full events
- Provide auditability and immutability for supply chain data

## Deployment
1. Set up a local Fabric network (see Fabric docs or use test network scripts).
2. Deploy this chaincode using the Fabric CLI or scripts.
3. Update the backend connection profile to point to your Fabric network.

## Testing
- Use Fabric CLI or SDK to invoke and query chaincode functions.
- Add unit tests for chaincode logic.

## Integration
- The backend submits event hashes or events to Fabric after ingesting from devices.
- Store Fabric transaction IDs in the backend database for traceability.

## License
MIT 