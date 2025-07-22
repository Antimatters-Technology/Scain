package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ScainChaincode provides functions for managing EPCIS events
type ScainChaincode struct {
	contractapi.Contract
}

// EventRecord represents an EPCIS event stored on the blockchain
type EventRecord struct {
	EventID     string    `json:"eventId"`
	EventHash   string    `json:"eventHash"`
	Timestamp   time.Time `json:"timestamp"`
	EventType   string    `json:"eventType"`
	DeviceID    string    `json:"deviceId,omitempty"`
	Data        string    `json:"data"`
	TxID        string    `json:"txId,omitempty"`
	Owner       string    `json:"owner,omitempty"`
}

// StoreEvent stores an EPCIS event on the blockchain
func (s *ScainChaincode) StoreEvent(ctx contractapi.TransactionContextInterface, 
	eventID string, eventHash string, timestamp string, eventType string, data string) error {
	
	// Check if event already exists
	exists, err := s.EventExists(ctx, eventID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("event %s already exists", eventID)
	}

	// Parse timestamp
	ts, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return fmt.Errorf("invalid timestamp format: %v", err)
	}

	// Get transaction ID and submitter
	txID := ctx.GetStub().GetTxID()
	submitter, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get submitter identity: %v", err)
	}

	// Create event record
	event := EventRecord{
		EventID:   eventID,
		EventHash: eventHash,
		Timestamp: ts,
		EventType: eventType,
		Data:      data,
		TxID:      txID,
		Owner:     submitter,
	}

	// Marshal to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	// Store on ledger
	err = ctx.GetStub().PutState(eventID, eventJSON)
	if err != nil {
		return fmt.Errorf("failed to put state: %v", err)
	}

	// Create composite key for querying by type
	indexName := "eventType~eventId"
	eventTypeIndexKey, err := ctx.GetStub().CreateCompositeKey(indexName, []string{eventType, eventID})
	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}

	// Store index
	value := []byte{0x00}
	err = ctx.GetStub().PutState(eventTypeIndexKey, value)
	if err != nil {
		return fmt.Errorf("failed to put index state: %v", err)
	}

	log.Printf("Event stored: %s", eventID)
	return nil
}

// GetEvent retrieves an event by ID
func (s *ScainChaincode) GetEvent(ctx contractapi.TransactionContextInterface, eventID string) (*EventRecord, error) {
	eventJSON, err := ctx.GetStub().GetState(eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if eventJSON == nil {
		return nil, fmt.Errorf("event %s does not exist", eventID)
	}

	var event EventRecord
	err = json.Unmarshal(eventJSON, &event)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal event: %v", err)
	}

	return &event, nil
}

// EventExists checks if an event exists
func (s *ScainChaincode) EventExists(ctx contractapi.TransactionContextInterface, eventID string) (bool, error) {
	eventJSON, err := ctx.GetStub().GetState(eventID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return eventJSON != nil, nil
}

// GetEventsByType retrieves all events of a specific type
func (s *ScainChaincode) GetEventsByType(ctx contractapi.TransactionContextInterface, eventType string) ([]*EventRecord, error) {
	// Query by composite key
	indexName := "eventType~eventId"
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(indexName, []string{eventType})
	if err != nil {
		return nil, fmt.Errorf("failed to get state by partial composite key: %v", err)
	}
	defer resultsIterator.Close()

	var events []*EventRecord
	for resultsIterator.HasNext() {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next result: %v", err)
		}

		// Extract event ID from composite key
		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to split composite key: %v", err)
		}

		if len(compositeKeyParts) > 1 {
			eventID := compositeKeyParts[1]
			event, err := s.GetEvent(ctx, eventID)
			if err != nil {
				return nil, err
			}
			events = append(events, event)
		}
	}

	return events, nil
}

// GetEventHistory retrieves the transaction history for an event
func (s *ScainChaincode) GetEventHistory(ctx contractapi.TransactionContextInterface, eventID string) ([]map[string]interface{}, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for key %s: %v", eventID, err)
	}
	defer resultsIterator.Close()

	var history []map[string]interface{}
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next result: %v", err)
		}

		var event EventRecord
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &event)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal event: %v", err)
			}
		}

		record := map[string]interface{}{
			"txId":      response.TxId,
			"timestamp": response.Timestamp,
			"isDelete":  response.IsDelete,
			"value":     event,
		}
		history = append(history, record)
	}

	return history, nil
}

// VerifyEventHash verifies the integrity of an event by comparing hashes
func (s *ScainChaincode) VerifyEventHash(ctx contractapi.TransactionContextInterface, 
	eventID string, expectedHash string) (bool, error) {
	
	event, err := s.GetEvent(ctx, eventID)
	if err != nil {
		return false, err
	}

	return event.EventHash == expectedHash, nil
}

// InitLedger initializes the ledger with sample data (for testing)
func (s *ScainChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	log.Println("Scain chaincode initialized")
	return nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&ScainChaincode{})
	if err != nil {
		log.Panicf("Error creating Scain chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting Scain chaincode: %v", err)
	}
} 