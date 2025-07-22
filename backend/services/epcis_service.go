package services

import (
	"encoding/json"
	"fmt"
	"os"

	"scain-backend/database"
	"scain-backend/models"
	"scain-backend/utils"
	
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// EPCISService handles EPCIS event processing
type EPCISService struct{
	blockchainService *BlockchainService
}

// NewEPCISService creates a new EPCIS service instance
func NewEPCISService() *EPCISService {
	service := &EPCISService{}
	
	// Initialize blockchain service if enabled
	if os.Getenv("ENABLE_BLOCKCHAIN") == "true" {
		blockchainService, err := NewBlockchainService()
		if err != nil {
			logger.Warnf("Failed to initialize blockchain service: %v", err)
		} else {
			service.blockchainService = blockchainService
			logger.Info("Blockchain service initialized")
		}
	}
	
	return service
}

// CreateEvent processes and stores an EPCIS event
func (s *EPCISService) CreateEvent(event *models.EpcisEvent) (*database.Event, error) {
	// Compute hash for integrity
	hash, err := utils.ComputeSHA256(event)
	if err != nil {
		return nil, fmt.Errorf("failed to compute event hash: %w", err)
	}

	// Convert to JSON for storage
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event to JSON: %w", err)
	}

	// Create database event
	dbEvent := &database.Event{
		EventType:           string(event.EventType),
		EventTime:           event.EventTime,
		EventTimeZoneOffset: event.EventTimeZoneOffset,
		Hash:                hash,
		RawData:             string(eventJSON),
	}

	// Set optional fields
	if event.BizStep != nil {
		bizStep := string(*event.BizStep)
		dbEvent.BizStep = &bizStep
	}
	if event.Disposition != nil {
		dbEvent.Disposition = event.Disposition
	}
	if event.ReadPoint != nil {
		dbEvent.ReadPointID = &event.ReadPoint.ID
	}
	if event.BizLocation != nil {
		dbEvent.BizLocationID = &event.BizLocation.ID
	}
	if event.LotCode != nil {
		dbEvent.LotCode = event.LotCode
	}
	if event.DeviceID != nil {
		dbEvent.DeviceID = event.DeviceID
	}
	if event.DeviceTimestamp != nil {
		dbEvent.DeviceTimestamp = event.DeviceTimestamp
	}

	// Save to database
	if err := database.CreateEvent(dbEvent); err != nil {
		return nil, fmt.Errorf("failed to create event in database: %w", err)
	}

	// Submit to blockchain if enabled
	var blockchainTxID string
	if s.blockchainService != nil {
		record, err := s.blockchainService.SubmitEvent(event)
		if err != nil {
			logger.Warnf("Failed to submit event to blockchain: %v", err)
		} else {
			blockchainTxID = record.TxID
			// Update database with blockchain transaction ID
			dbEvent.BlockchainTxID = &blockchainTxID
			if err := database.UpdateEvent(dbEvent); err != nil {
				logger.Warnf("Failed to update event with blockchain TX ID: %v", err)
			}
		}
	}

	logger.WithFields(logrus.Fields{
		"eventId":       dbEvent.ID,
		"eventType":     dbEvent.EventType,
		"hash":          hash,
		"blockchainTx":  blockchainTxID,
	}).Info("EPCIS event created successfully")

	return dbEvent, nil
}

// GetEvent retrieves an EPCIS event by ID
func (s *EPCISService) GetEvent(id string) (*models.EpcisEvent, error) {
	dbEvent, err := database.GetEventByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get event from database: %w", err)
	}

	// Unmarshal the raw data back to EPCIS event
	var event models.EpcisEvent
	if err := json.Unmarshal([]byte(dbEvent.RawData), &event); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	return &event, nil
}

// TransformRawData transforms raw device data into EPCIS events
func (s *EPCISService) TransformRawData(payload *models.RawIngestPayload) ([]*models.EpcisEvent, error) {
	var events []*models.EpcisEvent

	switch payload.DeviceType {
	case models.ESP32DeviceType:
		events = s.transformESP32Data(payload)
	case models.ExpressLinkDeviceType:
		events = s.transformExpressLinkData(payload)
	case models.LoRaWANDeviceType:
		events = s.transformLoRaWANData(payload)
	case models.TrackerDeviceType:
		events = s.transformTrackerData(payload)
	case models.ERPDeviceType:
		events = s.transformERPData(payload)
	default:
		return nil, fmt.Errorf("unsupported device type: %s", payload.DeviceType)
	}

	logger.WithFields(logrus.Fields{
		"deviceType": payload.DeviceType,
		"deviceId":   payload.DeviceID,
		"eventCount": len(events),
	}).Info("Raw data transformed to EPCIS events")

	return events, nil
}

// transformESP32Data transforms ESP32 sensor data to EPCIS events
func (s *EPCISService) transformESP32Data(payload *models.RawIngestPayload) []*models.EpcisEvent {
	var events []*models.EpcisEvent

	// Create sensor reports from raw data
	var sensorReports []models.SensorReport
	for key, value := range payload.Data {
		sensorReports = append(sensorReports, models.SensorReport{
			Type:  key,
			Value: value,
			Time:  payload.Timestamp,
		})
	}

	// Create EPCIS ObjectEvent with sensor data
	event := &models.EpcisEvent{
		EventType:           models.ObjectEventType,
		EventTime:           payload.Timestamp,
		EventTimeZoneOffset: "+00:00", // Default to UTC
		DeviceID:            &payload.DeviceID,
		DeviceTimestamp:     &payload.Timestamp,
		LotCode:             payload.LotCode,
		SensorElementList: []models.SensorElement{
			{
				SensorMetaData: models.SensorMetadata{
					DeviceID: payload.DeviceID,
					DeviceMetadata: &models.DeviceMetadata{
						Type: payload.DeviceType,
					},
				},
				SensorReport: sensorReports,
			},
		},
	}

	events = append(events, event)
	return events
}

// transformExpressLinkData transforms AWS IoT ExpressLink data to EPCIS events
func (s *EPCISService) transformExpressLinkData(payload *models.RawIngestPayload) []*models.EpcisEvent {
	// Similar to ESP32 but with ExpressLink specific handling
	return s.transformESP32Data(payload) // For now, use same logic
}

// transformLoRaWANData transforms LoRaWAN data to EPCIS events
func (s *EPCISService) transformLoRaWANData(payload *models.RawIngestPayload) []*models.EpcisEvent {
	// Similar to ESP32 but with LoRaWAN specific handling
	return s.transformESP32Data(payload) // For now, use same logic
}

// transformTrackerData transforms GPS tracker data to EPCIS events
func (s *EPCISService) transformTrackerData(payload *models.RawIngestPayload) []*models.EpcisEvent {
	var events []*models.EpcisEvent

	// Extract location data
	lat, hasLat := payload.Data["latitude"]
	lng, hasLng := payload.Data["longitude"]

	if hasLat && hasLng {
		// Create location-based event
		readPoint := &models.ReadPoint{
			ID: fmt.Sprintf("geo:%v,%v", lat, lng),
		}

		event := &models.EpcisEvent{
			EventType:           models.ObjectEventType,
			EventTime:           payload.Timestamp,
			EventTimeZoneOffset: "+00:00",
			BizStep:             nil, // Could be set based on business logic
			ReadPoint:           readPoint,
			DeviceID:            &payload.DeviceID,
			DeviceTimestamp:     &payload.Timestamp,
			LotCode:             payload.LotCode,
		}

		events = append(events, event)
	}

	return events
}

// transformERPData transforms ERP system data to EPCIS events
func (s *EPCISService) transformERPData(payload *models.RawIngestPayload) []*models.EpcisEvent {
	var events []*models.EpcisEvent

	// Extract business step from ERP data
	bizStepStr, hasBizStep := payload.Data["businessStep"].(string)
	var bizStep *models.BusinessStep
	if hasBizStep {
		bs := models.BusinessStep(bizStepStr)
		bizStep = &bs
	}

	// Create business event
	event := &models.EpcisEvent{
		EventType:           models.TransactionEventType,
		EventTime:           payload.Timestamp,
		EventTimeZoneOffset: "+00:00",
		BizStep:             bizStep,
		DeviceID:            &payload.DeviceID,
		DeviceTimestamp:     &payload.Timestamp,
		LotCode:             payload.LotCode,
	}

	// Add business transactions if present
	if txns, hasTxns := payload.Data["transactions"].([]interface{}); hasTxns {
		var bizTransactions []models.BizTransaction
		for _, txn := range txns {
			if txnMap, ok := txn.(map[string]interface{}); ok {
				bizTransactions = append(bizTransactions, models.BizTransaction{
					Type:           txnMap["type"].(string),
					BizTransaction: txnMap["id"].(string),
				})
			}
		}
		event.BizTransactionList = bizTransactions
	}

	events = append(events, event)
	return events
}

// ProcessRawDataIngestion processes a raw data ingestion record
func (s *EPCISService) ProcessRawDataIngestion(ingestionID string) error {
	// This would typically be called by a background worker
	// For now, we'll implement basic processing logic
	
	logger.WithField("ingestionId", ingestionID).Info("Processing raw data ingestion")
	
	// TODO: Implement background processing of raw data
	// 1. Retrieve raw data from database
	// 2. Transform to EPCIS events
	// 3. Store EPCIS events
	// 4. Mark ingestion as processed
	
	return nil
} 