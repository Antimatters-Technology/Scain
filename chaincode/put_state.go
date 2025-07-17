/*
 * Scain Hyperledger Fabric Chaincode
 * Stores EPCIS 2.0 JSON events with SHA-256 hashing
 * Implements traceability record storage for food safety compliance
 */

package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ScainContract defines the smart contract
type ScainContract struct {
	contractapi.Contract
}

// EPCISEvent represents an EPCIS 2.0 event structure
type EPCISEvent struct {
	Context       string    `json:"@context"`
	Type          string    `json:"type"`
	SchemaVersion string    `json:"schemaVersion"`
	CreationDate  string    `json:"creationDate"`
	EPCISBody     EPCISBody `json:"epcisBody"`
}

// EPCISBody contains the event list
type EPCISBody struct {
	EventList []ObjectEvent `json:"eventList"`
}

// ObjectEvent represents a single EPCIS object event
type ObjectEvent struct {
	EventType         string              `json:"eventType"`
	EventTime         string              `json:"eventTime"`
	EventTimeZone     string              `json:"eventTimeZoneOffset"`
	RecordTime        string              `json:"recordTime"`
	EPCList           []string            `json:"epcList"`
	Action            string              `json:"action"`
	BizStep           string              `json:"bizStep"`
	Disposition       string              `json:"disposition"`
	SensorElementList []SensorElement     `json:"sensorElementList"`
	UserExtensions    map[string]interface{} `json:"userExtensions,omitempty"`
}

// SensorElement represents sensor data
type SensorElement struct {
	SensorMetadata SensorMetadata `json:"sensorMetadata"`
	SensorReport   []SensorReport `json:"sensorReport"`
}

// SensorMetadata contains sensor information
type SensorMetadata struct {
	Time           string `json:"time"`
	DeviceID       string `json:"deviceID"`
	DeviceMetadata string `json:"deviceMetadata"`
}

// SensorReport contains individual sensor readings
type SensorReport struct {
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
	UOM       string      `json:"uom"`
	Component string      `json:"component"`
}

// TraceabilityRecord represents a stored traceability record
type TraceabilityRecord struct {
	Hash         string    `json:"hash"`
	EPCISData    string    `json:"epcisData"`
	Timestamp    time.Time `json:"timestamp"`
	DeviceID     string    `json:"deviceId"`
	EPCList      []string  `json:"epcList"`
	EventType    string    `json:"eventType"`
	Temperature  float64   `json:"temperature,omitempty"`
	Humidity     float64   `json:"humidity,omitempty"`
	ProbeTemp    float64   `json:"probeTemp,omitempty"`
}

// QueryResult represents query response
type QueryResult struct {
	Key    string             `json:"key"`
	Record TraceabilityRecord `json:"record"`
}

// RecordEvent stores an EPCIS event in the ledger
func (c *ScainContract) RecordEvent(ctx contractapi.TransactionContextInterface, epcisJSON string) error {
	// Parse EPCIS JSON
	var event EPCISEvent
	if err := json.Unmarshal([]byte(epcisJSON), &event); err != nil {
		return fmt.Errorf("failed to parse EPCIS JSON: %v", err)
	}

	// Validate EPCIS structure
	if len(event.EPCISBody.EventList) == 0 {
		return fmt.Errorf("no events found in EPCIS body")
	}

	// Generate hash
	hash := sha256.Sum256([]byte(epcisJSON))
	hashString := fmt.Sprintf("%x", hash)

	// Extract sensor data for indexing
	objectEvent := event.EPCISBody.EventList[0]
	record := TraceabilityRecord{
		Hash:      hashString,
		EPCISData: epcisJSON,
		Timestamp: time.Now(),
		DeviceID:  extractDeviceID(objectEvent),
		EPCList:   objectEvent.EPCList,
		EventType: objectEvent.EventType,
	}

	// Extract sensor values for easier querying
	if len(objectEvent.SensorElementList) > 0 {
		for _, report := range objectEvent.SensorElementList[0].SensorReport {
			if val, ok := report.Value.(float64); ok {
				switch report.Component {
				case "air":
					if report.Type == "gs1:Temperature" {
						record.Temperature = val
					} else if report.Type == "gs1:RelativeHumidity" {
						record.Humidity = val
					}
				case "probe":
					if report.Type == "gs1:Temperature" {
						record.ProbeTemp = val
					}
				}
			}
		}
	}

	// Store record
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %v", err)
	}

	// Store with hash as key
	if err := ctx.GetStub().PutState(hashString, recordJSON); err != nil {
		return fmt.Errorf("failed to store record: %v", err)
	}

	// Create composite key for device-based queries
	deviceKey, err := ctx.GetStub().CreateCompositeKey("device", []string{record.DeviceID, record.Timestamp.Format(time.RFC3339)})
	if err != nil {
		return fmt.Errorf("failed to create device composite key: %v", err)
	}

	if err := ctx.GetStub().PutState(deviceKey, []byte(hashString)); err != nil {
		return fmt.Errorf("failed to store device index: %v", err)
	}

	// Create composite key for EPC-based queries
	for _, epc := range record.EPCList {
		epcKey, err := ctx.GetStub().CreateCompositeKey("epc", []string{epc, record.Timestamp.Format(time.RFC3339)})
		if err != nil {
			return fmt.Errorf("failed to create EPC composite key: %v", err)
		}

		if err := ctx.GetStub().PutState(epcKey, []byte(hashString)); err != nil {
			return fmt.Errorf("failed to store EPC index: %v", err)
		}
	}

	// Log successful storage
	log.Printf("Stored traceability record: %s for device: %s", hashString, record.DeviceID)

	return nil
}

// GetEvent retrieves an event by its hash
func (c *ScainContract) GetEvent(ctx contractapi.TransactionContextInterface, hash string) (*TraceabilityRecord, error) {
	recordJSON, err := ctx.GetStub().GetState(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to read record: %v", err)
	}

	if recordJSON == nil {
		return nil, fmt.Errorf("record not found: %s", hash)
	}

	var record TraceabilityRecord
	if err := json.Unmarshal(recordJSON, &record); err != nil {
		return nil, fmt.Errorf("failed to unmarshal record: %v", err)
	}

	return &record, nil
}

// GetEventsByDevice retrieves all events for a specific device
func (c *ScainContract) GetEventsByDevice(ctx contractapi.TransactionContextInterface, deviceID string) ([]*QueryResult, error) {
	return c.getEventsByCompositeKey(ctx, "device", deviceID)
}

// GetEventsByEPC retrieves all events for a specific EPC
func (c *ScainContract) GetEventsByEPC(ctx contractapi.TransactionContextInterface, epc string) ([]*QueryResult, error) {
	return c.getEventsByCompositeKey(ctx, "epc", epc)
}

// GetRecentEvents retrieves the most recent events (up to 100)
func (c *ScainContract) GetRecentEvents(ctx contractapi.TransactionContextInterface, limit int) ([]*QueryResult, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	iterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get state iterator: %v", err)
	}
	defer iterator.Close()

	var results []*QueryResult
	count := 0

	for iterator.HasNext() && count < limit {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next result: %v", err)
		}

		// Skip composite keys
		if len(queryResponse.Key) == 64 { // SHA-256 hash length
			var record TraceabilityRecord
			if err := json.Unmarshal(queryResponse.Value, &record); err == nil {
				results = append(results, &QueryResult{
					Key:    queryResponse.Key,
					Record: record,
				})
				count++
			}
		}
	}

	return results, nil
}

// ValidateEvent validates EPCIS event structure
func (c *ScainContract) ValidateEvent(ctx contractapi.TransactionContextInterface, epcisJSON string) (bool, error) {
	var event EPCISEvent
	if err := json.Unmarshal([]byte(epcisJSON), &event); err != nil {
		return false, fmt.Errorf("invalid JSON: %v", err)
	}

	// Basic validation
	if event.Context == "" {
		return false, fmt.Errorf("missing @context")
	}

	if event.Type != "EPCISDocument" {
		return false, fmt.Errorf("invalid document type: %s", event.Type)
	}

	if len(event.EPCISBody.EventList) == 0 {
		return false, fmt.Errorf("no events in event list")
	}

	return true, nil
}

// Helper function to extract device ID from object event
func extractDeviceID(event ObjectEvent) string {
	if len(event.SensorElementList) > 0 {
		return event.SensorElementList[0].SensorMetadata.DeviceID
	}
	return "unknown"
}

// Helper function to query by composite key
func (c *ScainContract) getEventsByCompositeKey(ctx contractapi.TransactionContextInterface, objectType, key string) ([]*QueryResult, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(objectType, []string{key})
	if err != nil {
		return nil, fmt.Errorf("failed to get state by composite key: %v", err)
	}
	defer iterator.Close()

	var results []*QueryResult

	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next result: %v", err)
		}

		// Get the actual record using the hash
		hash := string(queryResponse.Value)
		record, err := c.GetEvent(ctx, hash)
		if err != nil {
			continue // Skip invalid records
		}

		results = append(results, &QueryResult{
			Key:    hash,
			Record: *record,
		})
	}

	return results, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&ScainContract{})
	if err != nil {
		log.Panicf("Error creating Scain chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting Scain chaincode: %v", err)
	}
} 