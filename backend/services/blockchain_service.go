package services

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type BlockchainService struct {
	gateway  *gateway.Gateway
	network  *gateway.Network
	contract *gateway.Contract
}

type EventRecord struct {
	EventID     string    `json:"eventId"`
	EventHash   string    `json:"eventHash"`
	Timestamp   time.Time `json:"timestamp"`
	EventType   string    `json:"eventType"`
	DeviceID    string    `json:"deviceId"`
	Data        string    `json:"data"`
	TxID        string    `json:"txId,omitempty"`
}

// NewBlockchainService initializes connection to Hyperledger Fabric
func NewBlockchainService() (*BlockchainService, error) {
	// Load connection profile from environment or config
	ccpPath := os.Getenv("FABRIC_CCP_PATH")
	if ccpPath == "" {
		ccpPath = "./blockchain/network/connection-profile.yaml"
	}

	// Create gateway connection
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(ccpPath)),
		gateway.WithIdentity(os.Getenv("FABRIC_WALLET_PATH"), os.Getenv("FABRIC_USER_ID")),
		gateway.WithDiscovery(true),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gateway: %w", err)
	}

	// Get network and contract
	network, err := gw.GetNetwork(os.Getenv("FABRIC_CHANNEL_NAME"))
	if err != nil {
		return nil, fmt.Errorf("failed to get network: %w", err)
	}

	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	return &BlockchainService{
		gateway:  gw,
		network:  network,
		contract: contract,
	}, nil
}

// SubmitEvent submits an EPCIS event to the blockchain
func (bs *BlockchainService) SubmitEvent(eventData interface{}) (*EventRecord, error) {
	// Serialize event data
	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event data: %w", err)
	}

	// Create hash of the event
	hash := sha256.Sum256(eventJSON)
	eventHash := hex.EncodeToString(hash[:])

	// Create event record
	record := &EventRecord{
		EventID:   generateEventID(),
		EventHash: eventHash,
		Timestamp: time.Now(),
		EventType: "EPCIS_EVENT",
		Data:      string(eventJSON),
	}

	// Submit transaction to blockchain
	result, err := bs.contract.SubmitTransaction("StoreEvent", 
		record.EventID, 
		record.EventHash, 
		record.Timestamp.Format(time.RFC3339),
		record.EventType,
		record.Data,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to submit transaction: %w", err)
	}

	// Get transaction ID from result
	record.TxID = string(result)
	
	log.Printf("Event submitted to blockchain: EventID=%s, TxID=%s", record.EventID, record.TxID)
	return record, nil
}

// GetEvent retrieves an event from the blockchain by ID
func (bs *BlockchainService) GetEvent(eventID string) (*EventRecord, error) {
	result, err := bs.contract.EvaluateTransaction("GetEvent", eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var record EventRecord
	if err := json.Unmarshal(result, &record); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event record: %w", err)
	}

	return &record, nil
}

// VerifyEvent verifies the integrity of an event by comparing hashes
func (bs *BlockchainService) VerifyEvent(eventID string, eventData interface{}) (bool, error) {
	// Get event from blockchain
	record, err := bs.GetEvent(eventID)
	if err != nil {
		return false, err
	}

	// Calculate hash of provided data
	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		return false, fmt.Errorf("failed to marshal event data: %w", err)
	}

	hash := sha256.Sum256(eventJSON)
	calculatedHash := hex.EncodeToString(hash[:])

	// Compare hashes
	return record.EventHash == calculatedHash, nil
}

// Close closes the gateway connection
func (bs *BlockchainService) Close() error {
	if bs.gateway != nil {
		return bs.gateway.Close()
	}
	return nil
}

// generateEventID creates a unique event ID
func generateEventID() string {
	return fmt.Sprintf("evt_%d_%s", time.Now().Unix(), 
		hex.EncodeToString([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))[:8])
} 