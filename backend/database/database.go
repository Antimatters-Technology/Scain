package database

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Event represents EPCIS events in the database
type Event struct {
	ID                  string    `gorm:"primaryKey" json:"id"`
	EventType           string    `json:"eventType"`
	EventTime           time.Time `json:"eventTime"`
	EventTimeZoneOffset string    `json:"eventTimeZoneOffset"`
	BizStep             *string   `json:"bizStep"`
	Disposition         *string   `json:"disposition"`
	ReadPointID         *string   `json:"readPointId"`
	BizLocationID       *string   `json:"bizLocationId"`
	LotCode             *string   `json:"lotCode"`
	DeviceID            *string   `json:"deviceId"`
	DeviceTimestamp     *time.Time `json:"deviceTimestamp"`
	Hash                string    `json:"hash"`
	RawData             string    `gorm:"type:text" json:"rawData"` // Store full EPCIS event as JSON
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

// Device represents devices in the database
type Device struct {
	DeviceID         string     `gorm:"primaryKey" json:"deviceId"`
	Type             string     `json:"type"`
	SecureBoot       *bool      `json:"secureBoot"`
	OTACapable       *bool      `json:"otaCapable"`
	FirmwareVersion  *string    `json:"firmwareVersion"`
	CalibrationDate  *string    `json:"calibrationDate"`
	BatteryPct       *int       `json:"batteryPct"`
	LastHeartbeat    *time.Time `json:"lastHeartbeat"`
	JoinStatus       *string    `json:"joinStatus"`
	ClaimedAt        *time.Time `json:"claimedAt"`
	ClaimCode        *string    `json:"claimCode"`
	IsActive         bool       `gorm:"default:true" json:"isActive"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

// RawDataIngestion represents raw data before processing
type RawDataIngestion struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	DeviceType   string    `json:"deviceType"`
	DeviceID     string    `json:"deviceId"`
	Timestamp    time.Time `json:"timestamp"`
	LotCode      *string   `json:"lotCode"`
	RawData      string    `gorm:"type:text" json:"rawData"` // Store raw payload as JSON
	Metadata     string    `gorm:"type:text" json:"metadata"` // Store metadata as JSON
	ProcessedAt  *time.Time `json:"processedAt"`
	ProcessingStatus string `gorm:"default:'pending'" json:"processingStatus"` // pending, processed, failed
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// ClaimCode represents device claim codes
type ClaimCodeEntry struct {
	ID          string     `gorm:"primaryKey" json:"id"`
	ClaimCode   string     `gorm:"uniqueIndex" json:"claimCode"`
	DeviceType  string     `json:"deviceType"`
	DeviceID    *string    `json:"deviceId"` // Set when claimed
	IsUsed      bool       `gorm:"default:false" json:"isUsed"`
	UsedAt      *time.Time `json:"usedAt"`
	ExpiresAt   *time.Time `json:"expiresAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// InitDatabase initializes the database connection and runs migrations
func InitDatabase() error {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./scain.db"
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(
		&Event{},
		&Device{},
		&RawDataIngestion{},
		&ClaimCodeEntry{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database schema: %w", err)
	}

	return nil
}

// CreateEvent creates a new event in the database
func CreateEvent(event *Event) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	return DB.Create(event).Error
}

// GetEventByID retrieves an event by ID
func GetEventByID(id string) (*Event, error) {
	var event Event
	err := DB.First(&event, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// CreateDevice creates a new device in the database
func CreateDevice(device *Device) error {
	return DB.Create(device).Error
}

// GetDeviceByID retrieves a device by ID
func GetDeviceByID(deviceID string) (*Device, error) {
	var device Device
	err := DB.First(&device, "device_id = ?", deviceID).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

// CreateRawDataIngestion creates a new raw data ingestion record
func CreateRawDataIngestion(ingestion *RawDataIngestion) error {
	if ingestion.ID == "" {
		ingestion.ID = uuid.New().String()
	}
	return DB.Create(ingestion).Error
}

// GetClaimCodeByCode retrieves a claim code entry by code
func GetClaimCodeByCode(code string) (*ClaimCodeEntry, error) {
	var claimCode ClaimCodeEntry
	err := DB.First(&claimCode, "claim_code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &claimCode, nil
}

// MarkClaimCodeAsUsed marks a claim code as used and associates it with a device
func MarkClaimCodeAsUsed(code string, deviceID string) error {
	now := time.Now()
	return DB.Model(&ClaimCodeEntry{}).Where("claim_code = ?", code).Updates(map[string]interface{}{
		"is_used":   true,
		"device_id": deviceID,
		"used_at":   now,
	}).Error
}

// CreateClaimCode creates a new claim code (for testing/admin purposes)
func CreateClaimCode(claimCode *ClaimCodeEntry) error {
	if claimCode.ID == "" {
		claimCode.ID = uuid.New().String()
	}
	return DB.Create(claimCode).Error
} 