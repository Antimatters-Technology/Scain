package models

import (
	"time"
)

// EventType represents the EPCIS event types
type EventType string

const (
	ObjectEventType       EventType = "ObjectEvent"
	TransformationEventType EventType = "TransformationEvent"
	AggregationEventType   EventType = "AggregationEvent"
	TransactionEventType   EventType = "TransactionEvent"
)

// BusinessStep represents the business step enumeration
type BusinessStep string

const (
	Harvesting        BusinessStep = "harvesting"
	Cooling          BusinessStep = "cooling"
	InitialPacking   BusinessStep = "initialPacking"
	FirstLandReceiving BusinessStep = "firstLandReceiving"
	Shipping         BusinessStep = "shipping"
	Receiving        BusinessStep = "receiving"
	Transformation   BusinessStep = "transformation"
)

// DeviceType represents the device type enumeration
type DeviceType string

const (
	ESP32DeviceType      DeviceType = "ESP32"
	ExpressLinkDeviceType DeviceType = "ExpressLink"
	LoRaWANDeviceType    DeviceType = "LoRaWAN"
	TrackerDeviceType    DeviceType = "Tracker"
	ERPDeviceType        DeviceType = "ERP"
)

// JoinStatus represents the join status for devices
type JoinStatus string

const (
	JoinPending JoinStatus = "pending"
	JoinJoined  JoinStatus = "joined"
	JoinFailed  JoinStatus = "failed"
)

// SensorReport represents individual sensor readings
type SensorReport struct {
	Type  string      `json:"type" validate:"required"`
	Value interface{} `json:"value" validate:"required"`
	UOM   *string     `json:"uom,omitempty"`
	Time  time.Time   `json:"time" validate:"required"`
}

// DeviceMetadata represents metadata for sensor devices
type DeviceMetadata struct {
	Type             DeviceType `json:"type" validate:"required,oneof=ESP32 ExpressLink LoRaWAN Tracker ERP"`
	FirmwareVersion  *string    `json:"firmwareVersion,omitempty"`
	CalibrationDate  *string    `json:"calibrationDate,omitempty"`
}

// SensorMetadata represents sensor metadata
type SensorMetadata struct {
	DeviceID       string          `json:"deviceId" validate:"required"`
	DeviceMetadata *DeviceMetadata `json:"deviceMetadata,omitempty"`
}

// SensorElement represents a sensor element with metadata and reports
type SensorElement struct {
	SensorMetaData SensorMetadata  `json:"sensorMetaData" validate:"required"`
	SensorReport   []SensorReport  `json:"sensorReport" validate:"required,min=1"`
}

// ReadPoint represents a read point location
type ReadPoint struct {
	ID string `json:"id" validate:"required"`
}

// BizLocation represents a business location
type BizLocation struct {
	ID string `json:"id" validate:"required"`
}

// QuantityElement represents quantity information
type QuantityElement struct {
	EPCClass string  `json:"epcClass" validate:"required"`
	Quantity float64 `json:"quantity" validate:"required,min=0"`
	UOM      *string `json:"uom,omitempty"`
}

// BizTransaction represents a business transaction
type BizTransaction struct {
	Type           string `json:"type" validate:"required"`
	BizTransaction string `json:"bizTransaction" validate:"required"`
}

// EpcisEvent represents the main EPCIS event structure
type EpcisEvent struct {
	EventType           EventType             `json:"eventType" validate:"required,oneof=ObjectEvent TransformationEvent AggregationEvent TransactionEvent"`
	EventTime           time.Time             `json:"eventTime" validate:"required"`
	EventTimeZoneOffset string                `json:"eventTimeZoneOffset" validate:"required"`
	BizStep             *BusinessStep         `json:"bizStep,omitempty"`
	Disposition         *string               `json:"disposition,omitempty"`
	ReadPoint           *ReadPoint            `json:"readPoint,omitempty"`
	BizLocation         *BizLocation          `json:"bizLocation,omitempty"`
	EPCList             []string              `json:"epcList,omitempty"`
	QuantityList        []QuantityElement     `json:"quantityList,omitempty"`
	BizTransactionList  []BizTransaction      `json:"bizTransactionList,omitempty"`
	SensorElementList   []SensorElement       `json:"sensorElementList,omitempty"`
	
	// Custom extensions
	LotCode         *string    `json:"lotCode,omitempty"`
	DeviceID        *string    `json:"deviceId,omitempty"`
	DeviceTimestamp *time.Time `json:"deviceTimestamp,omitempty"`
}

// DeviceInfo represents comprehensive device information
type DeviceInfo struct {
	DeviceID         string      `json:"deviceId" validate:"required"`
	Type             DeviceType  `json:"type" validate:"required,oneof=ESP32 ExpressLink LoRaWAN Tracker ERP"`
	SecureBoot       *bool       `json:"secureBoot,omitempty"`
	OTACapable       *bool       `json:"otaCapable,omitempty"`
	FirmwareVersion  *string     `json:"firmwareVersion,omitempty"`
	CalibrationDate  *string     `json:"calibrationDate,omitempty"`
	RegulatoryCerts  []string    `json:"regulatoryCerts,omitempty"`
	BatteryPct       *int        `json:"batteryPct,omitempty" validate:"omitempty,min=0,max=100"`
	LastHeartbeat    *time.Time  `json:"lastHeartbeat,omitempty"`
	JoinStatus       *JoinStatus `json:"joinStatus,omitempty"`
}

// RawIngestPayload represents raw device data before normalization
type RawIngestPayload struct {
	DeviceType DeviceType             `json:"deviceType" validate:"required,oneof=ESP32 ExpressLink LoRaWAN Tracker ERP"`
	DeviceID   string                 `json:"deviceId" validate:"required"`
	Timestamp  time.Time              `json:"timestamp" validate:"required"`
	LotCode    *string                `json:"lotCode,omitempty"`
	Data       map[string]interface{} `json:"data" validate:"required"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// ClaimCode represents device claim codes with validation
type ClaimCode struct {
	ClaimCode   string            `json:"claimCode" validate:"required,len=8,alphanum"`
	Type        DeviceType        `json:"type" validate:"required,oneof=ESP32 ExpressLink LoRaWAN Tracker"`
	Identifiers map[string]string `json:"identifiers,omitempty"`
} 