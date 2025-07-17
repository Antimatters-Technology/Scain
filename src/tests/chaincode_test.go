package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStub for testing
type MockStub struct {
	mock.Mock
	shim.ChaincodeStubInterface
}

func (m *MockStub) PutState(key string, value []byte) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockStub) GetState(key string) ([]byte, error) {
	args := m.Called(key)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockStub) CreateCompositeKey(objectType string, attributes []string) (string, error) {
	args := m.Called(objectType, attributes)
	return args.String(0), args.Error(1)
}

func (m *MockStub) GetStateByPartialCompositeKey(objectType string, attributes []string) (shim.StateQueryIteratorInterface, error) {
	args := m.Called(objectType, attributes)
	return args.Get(0).(shim.StateQueryIteratorInterface), args.Error(1)
}

func (m *MockStub) GetStateByRange(startKey, endKey string) (shim.StateQueryIteratorInterface, error) {
	args := m.Called(startKey, endKey)
	return args.Get(0).(shim.StateQueryIteratorInterface), args.Error(1)
}

// MockTransactionContext for testing
type MockTransactionContext struct {
	mock.Mock
	contractapi.TransactionContextInterface
	stub *MockStub
}

func (m *MockTransactionContext) GetStub() shim.ChaincodeStubInterface {
	return m.stub
}

// Test data
const testEPCISJSON = `{
  "@context": "https://ref.gs1.org/standards/epcis/2.0.0/epcis-context.jsonld",
  "type": "EPCISDocument",
  "schemaVersion": "2.0",
  "creationDate": "2024-01-15T12:00:00Z",
  "epcisBody": {
    "eventList": [
      {
        "eventType": "ObjectEvent",
        "eventTime": "2024-01-15T12:00:00Z",
        "eventTimeZoneOffset": "+00:00",
        "recordTime": "2024-01-15T12:00:00Z",
        "epcList": ["urn:epc:id:sgtin:0614141.KG-ESP32-001"],
        "action": "OBSERVE",
        "bizStep": "urn:epcglobal:cbv:bizstep:sensor_reporting",
        "disposition": "urn:epcglobal:cbv:disp:in_transit",
        "sensorElementList": [
          {
            "sensorMetadata": {
              "time": "2024-01-15T12:00:00Z",
              "deviceID": "KG-ESP32-001",
              "deviceMetadata": "KnowGraph ESP32 Node v1.0"
            },
            "sensorReport": [
              {
                "type": "gs1:Temperature",
                "value": 4.5,
                "uom": "CEL",
                "component": "probe"
              },
              {
                "type": "gs1:Temperature",
                "value": 7.2,
                "uom": "CEL",
                "component": "air"
              },
              {
                "type": "gs1:RelativeHumidity",
                "value": 65.3,
                "uom": "A93",
                "component": "air"
              }
            ]
          }
        ]
      }
    ]
  }
}`

func TestSHA256Hashing(t *testing.T) {
	// Test SHA-256 hash generation
	hash := sha256.Sum256([]byte(testEPCISJSON))
	hashString := fmt.Sprintf("%x", hash)

	assert.Equal(t, 64, len(hashString), "SHA-256 hash should be 64 characters long")
	assert.Regexp(t, "^[a-f0-9]{64}$", hashString, "Hash should be valid hex string")

	// Test hash consistency
	hash2 := sha256.Sum256([]byte(testEPCISJSON))
	hashString2 := fmt.Sprintf("%x", hash2)
	assert.Equal(t, hashString, hashString2, "Hash should be consistent")

	// Test different input produces different hash
	modifiedJSON := testEPCISJSON + " "
	hash3 := sha256.Sum256([]byte(modifiedJSON))
	hashString3 := fmt.Sprintf("%x", hash3)
	assert.NotEqual(t, hashString, hashString3, "Different inputs should produce different hashes")
}

func TestEPCISJSONParsing(t *testing.T) {
	var event EPCISEvent
	err := json.Unmarshal([]byte(testEPCISJSON), &event)

	assert.NoError(t, err, "Should parse valid EPCIS JSON")
	assert.Equal(t, "EPCISDocument", event.Type, "Document type should be EPCISDocument")
	assert.Equal(t, "2.0", event.SchemaVersion, "Schema version should be 2.0")
	assert.Equal(t, 1, len(event.EPCISBody.EventList), "Should have one event")

	objectEvent := event.EPCISBody.EventList[0]
	assert.Equal(t, "ObjectEvent", objectEvent.EventType, "Event type should be ObjectEvent")
	assert.Equal(t, "OBSERVE", objectEvent.Action, "Action should be OBSERVE")
	assert.Equal(t, 1, len(objectEvent.EPCList), "Should have one EPC")
	assert.Equal(t, "urn:epc:id:sgtin:0614141.KG-ESP32-001", objectEvent.EPCList[0], "EPC should match")
}

func TestRecordEvent(t *testing.T) {
	// Setup
	contract := &KnowGraphContract{}
	mockStub := &MockStub{}
	mockCtx := &MockTransactionContext{stub: mockStub}

	// Expected hash
	hash := sha256.Sum256([]byte(testEPCISJSON))
	expectedHash := fmt.Sprintf("%x", hash)

	// Mock expectations
	mockStub.On("PutState", expectedHash, mock.AnythingOfType("[]uint8")).Return(nil)
	mockStub.On("CreateCompositeKey", "device", []string{"KG-ESP32-001", mock.AnythingOfType("string")}).Return("device~KG-ESP32-001~timestamp", nil)
	mockStub.On("PutState", "device~KG-ESP32-001~timestamp", []byte(expectedHash)).Return(nil)
	mockStub.On("CreateCompositeKey", "epc", []string{"urn:epc:id:sgtin:0614141.KG-ESP32-001", mock.AnythingOfType("string")}).Return("epc~urn:epc:id:sgtin:0614141.KG-ESP32-001~timestamp", nil)
	mockStub.On("PutState", "epc~urn:epc:id:sgtin:0614141.KG-ESP32-001~timestamp", []byte(expectedHash)).Return(nil)

	// Test
	err := contract.RecordEvent(mockCtx, testEPCISJSON)

	// Assertions
	assert.NoError(t, err, "RecordEvent should succeed")
	mockStub.AssertExpectations(t)
}

func TestRecordEventInvalidJSON(t *testing.T) {
	contract := &KnowGraphContract{}
	mockStub := &MockStub{}
	mockCtx := &MockTransactionContext{stub: mockStub}

	invalidJSON := `{"invalid": json}`

	err := contract.RecordEvent(mockCtx, invalidJSON)
	assert.Error(t, err, "Should return error for invalid JSON")
	assert.Contains(t, err.Error(), "failed to parse EPCIS JSON", "Error should mention JSON parsing")
}

func TestRecordEventEmptyEventList(t *testing.T) {
	contract := &KnowGraphContract{}
	mockStub := &MockStub{}
	mockCtx := &MockTransactionContext{stub: mockStub}

	emptyEventJSON := `{
		"@context": "https://ref.gs1.org/standards/epcis/2.0.0/epcis-context.jsonld",
		"type": "EPCISDocument",
		"schemaVersion": "2.0",
		"creationDate": "2024-01-15T12:00:00Z",
		"epcisBody": {
			"eventList": []
		}
	}`

	err := contract.RecordEvent(mockCtx, emptyEventJSON)
	assert.Error(t, err, "Should return error for empty event list")
	assert.Contains(t, err.Error(), "no events found", "Error should mention empty event list")
}

func TestGetEvent(t *testing.T) {
	contract := &KnowGraphContract{}
	mockStub := &MockStub{}
	mockCtx := &MockTransactionContext{stub: mockStub}

	hash := "test-hash"
	expectedRecord := TraceabilityRecord{
		Hash:        hash,
		EPCISData:   testEPCISJSON,
		Timestamp:   time.Now(),
		DeviceID:    "KG-ESP32-001",
		EPCList:     []string{"urn:epc:id:sgtin:0614141.KG-ESP32-001"},
		EventType:   "ObjectEvent",
		Temperature: 7.2,
		Humidity:    65.3,
		ProbeTemp:   4.5,
	}

	recordJSON, _ := json.Marshal(expectedRecord)
	mockStub.On("GetState", hash).Return(recordJSON, nil)

	record, err := contract.GetEvent(mockCtx, hash)

	assert.NoError(t, err, "GetEvent should succeed")
	assert.Equal(t, expectedRecord.Hash, record.Hash, "Hash should match")
	assert.Equal(t, expectedRecord.DeviceID, record.DeviceID, "Device ID should match")
	assert.Equal(t, expectedRecord.EventType, record.EventType, "Event type should match")
	mockStub.AssertExpectations(t)
}

func TestGetEventNotFound(t *testing.T) {
	contract := &KnowGraphContract{}
	mockStub := &MockStub{}
	mockCtx := &MockTransactionContext{stub: mockStub}

	hash := "nonexistent-hash"
	mockStub.On("GetState", hash).Return([]byte(nil), nil)

	record, err := contract.GetEvent(mockCtx, hash)

	assert.Error(t, err, "Should return error for non-existent record")
	assert.Nil(t, record, "Record should be nil")
	assert.Contains(t, err.Error(), "record not found", "Error should mention record not found")
	mockStub.AssertExpectations(t)
}

func TestValidateEvent(t *testing.T) {
	contract := &KnowGraphContract{}
	mockStub := &MockStub{}
	mockCtx := &MockTransactionContext{stub: mockStub}

	// Test valid EPCIS
	valid, err := contract.ValidateEvent(mockCtx, testEPCISJSON)
	assert.NoError(t, err, "Should validate valid EPCIS")
	assert.True(t, valid, "Valid EPCIS should return true")

	// Test invalid JSON
	valid, err = contract.ValidateEvent(mockCtx, `{"invalid": json}`)
	assert.Error(t, err, "Should return error for invalid JSON")
	assert.False(t, valid, "Invalid JSON should return false")

	// Test missing context
	invalidEPCIS := `{
		"type": "EPCISDocument",
		"schemaVersion": "2.0",
		"epcisBody": {
			"eventList": [{}]
		}
	}`
	valid, err = contract.ValidateEvent(mockCtx, invalidEPCIS)
	assert.Error(t, err, "Should return error for missing context")
	assert.False(t, valid, "Invalid EPCIS should return false")
	assert.Contains(t, err.Error(), "missing @context", "Error should mention missing context")

	// Test wrong document type
	wrongTypeEPCIS := `{
		"@context": "https://ref.gs1.org/standards/epcis/2.0.0/epcis-context.jsonld",
		"type": "WrongType",
		"schemaVersion": "2.0",
		"epcisBody": {
			"eventList": [{}]
		}
	}`
	valid, err = contract.ValidateEvent(mockCtx, wrongTypeEPCIS)
	assert.Error(t, err, "Should return error for wrong document type")
	assert.False(t, valid, "Wrong type should return false")
	assert.Contains(t, err.Error(), "invalid document type", "Error should mention invalid document type")

	// Test empty event list
	emptyEventEPCIS := `{
		"@context": "https://ref.gs1.org/standards/epcis/2.0.0/epcis-context.jsonld",
		"type": "EPCISDocument",
		"schemaVersion": "2.0",
		"epcisBody": {
			"eventList": []
		}
	}`
	valid, err = contract.ValidateEvent(mockCtx, emptyEventEPCIS)
	assert.Error(t, err, "Should return error for empty event list")
	assert.False(t, valid, "Empty event list should return false")
	assert.Contains(t, err.Error(), "no events in event list", "Error should mention empty event list")
}

func TestExtractDeviceID(t *testing.T) {
	// Test with sensor data
	eventWithSensor := ObjectEvent{
		SensorElementList: []SensorElement{
			{
				SensorMetadata: SensorMetadata{
					DeviceID: "KG-ESP32-001",
				},
			},
		},
	}
	deviceID := extractDeviceID(eventWithSensor)
	assert.Equal(t, "KG-ESP32-001", deviceID, "Should extract device ID from sensor metadata")

	// Test without sensor data
	eventWithoutSensor := ObjectEvent{
		SensorElementList: []SensorElement{},
	}
	deviceID = extractDeviceID(eventWithoutSensor)
	assert.Equal(t, "unknown", deviceID, "Should return 'unknown' when no sensor data")
}

func TestSensorValueExtraction(t *testing.T) {
	var event EPCISEvent
	err := json.Unmarshal([]byte(testEPCISJSON), &event)
	assert.NoError(t, err, "Should parse test JSON")

	objectEvent := event.EPCISBody.EventList[0]
	sensorReports := objectEvent.SensorElementList[0].SensorReport

	// Test temperature extraction
	var airTemp, probeTemp, humidity float64
	for _, report := range sensorReports {
		if val, ok := report.Value.(float64); ok {
			switch report.Component {
			case "air":
				if report.Type == "gs1:Temperature" {
					airTemp = val
				} else if report.Type == "gs1:RelativeHumidity" {
					humidity = val
				}
			case "probe":
				if report.Type == "gs1:Temperature" {
					probeTemp = val
				}
			}
		}
	}

	assert.Equal(t, 7.2, airTemp, "Air temperature should be extracted correctly")
	assert.Equal(t, 4.5, probeTemp, "Probe temperature should be extracted correctly")
	assert.Equal(t, 65.3, humidity, "Humidity should be extracted correctly")
}

func BenchmarkSHA256Hashing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash := sha256.Sum256([]byte(testEPCISJSON))
		_ = fmt.Sprintf("%x", hash)
	}
}

func BenchmarkJSONParsing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var event EPCISEvent
		json.Unmarshal([]byte(testEPCISJSON), &event)
	}
} 