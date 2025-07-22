package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"scain-backend/database"
	"scain-backend/models"
	
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// DeviceService handles device management operations
type DeviceService struct{}

// NewDeviceService creates a new device service instance
func NewDeviceService() *DeviceService {
	return &DeviceService{}
}

// RegisterDevice registers a new device in the system
func (s *DeviceService) RegisterDevice(deviceInfo *models.DeviceInfo) (*database.Device, error) {
	// Check if device already exists
	existingDevice, err := database.GetDeviceByID(deviceInfo.DeviceID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing device: %w", err)
	}
	if existingDevice != nil {
		return nil, fmt.Errorf("device with ID %s already exists", deviceInfo.DeviceID)
	}

	// Create database device
	dbDevice := &database.Device{
		DeviceID:        deviceInfo.DeviceID,
		Type:            string(deviceInfo.Type),
		SecureBoot:      deviceInfo.SecureBoot,
		OTACapable:      deviceInfo.OTACapable,
		FirmwareVersion: deviceInfo.FirmwareVersion,
		CalibrationDate: deviceInfo.CalibrationDate,
		BatteryPct:      deviceInfo.BatteryPct,
		LastHeartbeat:   deviceInfo.LastHeartbeat,
		IsActive:        true,
	}

	// Set join status if provided
	if deviceInfo.JoinStatus != nil {
		status := string(*deviceInfo.JoinStatus)
		dbDevice.JoinStatus = &status
	}

	// Save to database
	if err := database.CreateDevice(dbDevice); err != nil {
		return nil, fmt.Errorf("failed to create device in database: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"deviceId": deviceInfo.DeviceID,
		"type":     deviceInfo.Type,
	}).Info("Device registered successfully")

	return dbDevice, nil
}

// GetDevice retrieves device information by ID
func (s *DeviceService) GetDevice(deviceID string) (*models.DeviceInfo, error) {
	dbDevice, err := database.GetDeviceByID(deviceID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("device not found: %s", deviceID)
		}
		return nil, fmt.Errorf("failed to get device from database: %w", err)
	}

	// Convert to models.DeviceInfo
	deviceInfo := &models.DeviceInfo{
		DeviceID:        dbDevice.DeviceID,
		Type:            models.DeviceType(dbDevice.Type),
		SecureBoot:      dbDevice.SecureBoot,
		OTACapable:      dbDevice.OTACapable,
		FirmwareVersion: dbDevice.FirmwareVersion,
		CalibrationDate: dbDevice.CalibrationDate,
		BatteryPct:      dbDevice.BatteryPct,
		LastHeartbeat:   dbDevice.LastHeartbeat,
	}

	// Set join status if available
	if dbDevice.JoinStatus != nil {
		status := models.JoinStatus(*dbDevice.JoinStatus)
		deviceInfo.JoinStatus = &status
	}

	return deviceInfo, nil
}

// ClaimDevice claims a device using a claim code
func (s *DeviceService) ClaimDevice(claimCode *models.ClaimCode) (*database.Device, error) {
	// Validate claim code
	if err := s.validateClaimCode(claimCode.ClaimCode); err != nil {
		return nil, fmt.Errorf("invalid claim code: %w", err)
	}

	// Check if claim code exists and is valid
	claimEntry, err := database.GetClaimCodeByCode(claimCode.ClaimCode)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid claim code: %s", claimCode.ClaimCode)
		}
		return nil, fmt.Errorf("failed to validate claim code: %w", err)
	}

	// Check if claim code is already used
	if claimEntry.IsUsed {
		return nil, fmt.Errorf("claim code %s has already been used", claimCode.ClaimCode)
	}

	// Check if claim code has expired
	if claimEntry.ExpiresAt != nil && time.Now().After(*claimEntry.ExpiresAt) {
		return nil, fmt.Errorf("claim code %s has expired", claimCode.ClaimCode)
	}

	// Check if device type matches
	if claimEntry.DeviceType != string(claimCode.Type) {
		return nil, fmt.Errorf("device type mismatch: expected %s, got %s", 
			claimEntry.DeviceType, claimCode.Type)
	}

	// Generate device ID based on claim code and type
	deviceID := s.generateDeviceID(claimCode.Type, claimCode.ClaimCode)

	// Create new device
	device := &database.Device{
		DeviceID:    deviceID,
		Type:        string(claimCode.Type),
		ClaimCode:   &claimCode.ClaimCode,
		ClaimedAt:   &time.Time{},
		IsActive:    true,
	}
	now := time.Now()
	device.ClaimedAt = &now

	// Save device to database
	if err := database.CreateDevice(device); err != nil {
		return nil, fmt.Errorf("failed to create claimed device: %w", err)
	}

	// Mark claim code as used
	if err := database.MarkClaimCodeAsUsed(claimCode.ClaimCode, deviceID); err != nil {
		return nil, fmt.Errorf("failed to mark claim code as used: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"deviceId":  deviceID,
		"claimCode": claimCode.ClaimCode,
		"type":      claimCode.Type,
	}).Info("Device claimed successfully")

	return device, nil
}

// UpdateDeviceHeartbeat updates the last heartbeat timestamp for a device
func (s *DeviceService) UpdateDeviceHeartbeat(deviceID string, batteryPct *int) error {
	now := time.Now()
	updates := map[string]interface{}{
		"last_heartbeat": now,
	}

	if batteryPct != nil {
		updates["battery_pct"] = *batteryPct
	}

	err := database.DB.Model(&database.Device{}).
		Where("device_id = ?", deviceID).
		Updates(updates).Error

	if err != nil {
		return fmt.Errorf("failed to update device heartbeat: %w", err)
	}

	return nil
}

// GenerateClaimCodes generates claim codes for testing/admin purposes
func (s *DeviceService) GenerateClaimCodes(deviceType models.DeviceType, count int, expiresInHours int) ([]*database.ClaimCodeEntry, error) {
	var claimCodes []*database.ClaimCodeEntry
	
	for i := 0; i < count; i++ {
		code := s.generateClaimCode()
		
		claimCode := &database.ClaimCodeEntry{
			ClaimCode:  code,
			DeviceType: string(deviceType),
			IsUsed:     false,
		}

		// Set expiration if specified
		if expiresInHours > 0 {
			expiry := time.Now().Add(time.Duration(expiresInHours) * time.Hour)
			claimCode.ExpiresAt = &expiry
		}

		if err := database.CreateClaimCode(claimCode); err != nil {
			return nil, fmt.Errorf("failed to create claim code: %w", err)
		}

		claimCodes = append(claimCodes, claimCode)
	}

	logger.WithFields(logrus.Fields{
		"deviceType": deviceType,
		"count":      count,
		"expires":    expiresInHours,
	}).Info("Claim codes generated")

	return claimCodes, nil
}

// ProcessRawDataIngestion stores raw device data for processing
func (s *DeviceService) ProcessRawDataIngestion(payload *models.RawIngestPayload) (*database.RawDataIngestion, error) {
	// Convert data and metadata to JSON
	dataJSON, err := json.Marshal(payload.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	metadataJSON, err := json.Marshal(payload.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata to JSON: %w", err)
	}

	// Create raw data ingestion record
	ingestion := &database.RawDataIngestion{
		DeviceType:       string(payload.DeviceType),
		DeviceID:         payload.DeviceID,
		Timestamp:        payload.Timestamp,
		LotCode:          payload.LotCode,
		RawData:          string(dataJSON),
		Metadata:         string(metadataJSON),
		ProcessingStatus: "pending",
	}

	// Save to database
	if err := database.CreateRawDataIngestion(ingestion); err != nil {
		return nil, fmt.Errorf("failed to create raw data ingestion: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"ingestionId": ingestion.ID,
		"deviceType":  payload.DeviceType,
		"deviceId":    payload.DeviceID,
	}).Info("Raw data ingestion created")

	// Update device heartbeat
	if payload.DeviceType != models.ERPDeviceType {
		s.UpdateDeviceHeartbeat(payload.DeviceID, nil)
	}

	return ingestion, nil
}

// validateClaimCode validates the format of a claim code
func (s *DeviceService) validateClaimCode(code string) error {
	if len(code) != 8 {
		return fmt.Errorf("claim code must be 8 characters long")
	}

	// Check if alphanumeric
	for _, char := range code {
		if !((char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return fmt.Errorf("claim code must contain only uppercase letters and numbers")
		}
	}

	return nil
}

// generateDeviceID generates a unique device ID based on type and claim code
func (s *DeviceService) generateDeviceID(deviceType models.DeviceType, claimCode string) string {
	prefix := strings.ToLower(string(deviceType))
	return fmt.Sprintf("%s-%s-%d", prefix, strings.ToLower(claimCode), time.Now().Unix())
}

// generateClaimCode generates an 8-character alphanumeric claim code
func (s *DeviceService) generateClaimCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8
	
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
	
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	
	return string(code)
} 