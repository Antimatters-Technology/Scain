package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	
	"scain-backend/database"
	"scain-backend/middleware"
	"scain-backend/models"
	"scain-backend/services"
)

var logger = logrus.New()
var validate = validator.New()

// Service instances
var epcisService *services.EPCISService
var deviceService *services.DeviceService

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// APIResponse represents the API info response
type APIResponse struct {
	Message   string   `json:"message"`
	Version   string   `json:"version"`
	Endpoints []string `json:"endpoints"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found")
	}

	// Configure logger
	logger.SetFormatter(&logrus.JSONFormatter{})
	if os.Getenv("NODE_ENV") == "development" {
		logger.SetLevel(logrus.DebugLevel)
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	// Initialize database
	if err := database.InitDatabase(); err != nil {
		logger.WithError(err).Fatal("Failed to initialize database")
	}

	// Initialize services
	epcisService = services.NewEPCISService()
	deviceService = services.NewDeviceService()
}

// healthHandler handles the health check endpoint
func healthHandler(c *gin.Context) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   "1.0.0",
	}
	c.JSON(http.StatusOK, response)
}

// apiHandler handles the API info endpoint
func apiHandler(c *gin.Context) {
	response := APIResponse{
		Message: "Scain Backend API - EPCIS Supply Chain Traceability",
		Version: "1.0.0",
		Endpoints: []string{
			"GET /health - Health check",
			"GET /api - API information",
			"POST /api/events - Create EPCIS event",
			"GET /api/events/{id} - Get EPCIS event",
			"POST /api/devices - Register device",
			"GET /api/devices/{deviceId} - Get device info",
			"POST /api/ingest - Raw device data ingestion",
			"POST /api/claim - Claim device with code",
		},
	}
	c.JSON(http.StatusOK, response)
}

// createEventHandler handles EPCIS event creation
func createEventHandler(c *gin.Context) {
	var event models.EpcisEvent
	
	if err := c.ShouldBindJSON(&event); err != nil {
		validationError := middleware.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, validationError)
		return
	}
	
	// Validate the event
	if err := validate.Struct(&event); err != nil {
		validationError := middleware.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, validationError)
		return
	}
	
	// Create event using service
	dbEvent, err := epcisService.CreateEvent(&event)
	if err != nil {
		logger.WithError(err).Error("Failed to create EPCIS event")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Message: "Failed to process event",
			Code:    500,
		})
		return
	}
	
	logger.WithFields(logrus.Fields{
		"eventId":   dbEvent.ID,
		"eventType": dbEvent.EventType,
		"hash":      dbEvent.Hash,
	}).Info("EPCIS event created")
	
	response := map[string]interface{}{
		"status":  "created",
		"eventId": dbEvent.ID,
		"hash":    dbEvent.Hash,
		"event":   event,
	}
	
	c.JSON(http.StatusCreated, response)
}

// getEventHandler handles EPCIS event retrieval
func getEventHandler(c *gin.Context) {
	eventId := c.Param("id")
	
	if eventId == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Missing event ID",
			Message: "Event ID is required",
			Code:    400,
		})
		return
	}
	
	// Retrieve event using service
	event, err := epcisService.GetEvent(eventId)
	if err != nil {
		logger.WithError(err).WithField("eventId", eventId).Error("Failed to retrieve event")
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Event not found",
			Message: err.Error(),
			Code:    404,
		})
		return
	}
	
	response := map[string]interface{}{
		"status": "found",
		"event":  event,
	}
	
	c.JSON(http.StatusOK, response)
}

// registerDeviceHandler handles device registration
func registerDeviceHandler(c *gin.Context) {
	var device models.DeviceInfo
	
	if err := c.ShouldBindJSON(&device); err != nil {
		validationError := middleware.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, validationError)
		return
	}
	
	// Validate the device
	if err := validate.Struct(&device); err != nil {
		validationError := middleware.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, validationError)
		return
	}
	
	// Register device using service
	dbDevice, err := deviceService.RegisterDevice(&device)
	if err != nil {
		logger.WithError(err).Error("Failed to register device")
		
		// Check if device already exists
		if err.Error() == "device with ID "+device.DeviceID+" already exists" {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "Device already exists",
				Message: err.Error(),
				Code:    409,
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Message: "Failed to register device",
			Code:    500,
		})
		return
	}
	
	logger.WithFields(logrus.Fields{
		"deviceId": device.DeviceID,
		"type":     device.Type,
	}).Info("Device registered")
	
	response := map[string]interface{}{
		"status": "registered",
		"device": dbDevice,
	}
	
	c.JSON(http.StatusCreated, response)
}

// getDeviceHandler handles device info retrieval
func getDeviceHandler(c *gin.Context) {
	deviceId := c.Param("deviceId")
	
	if deviceId == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Missing device ID",
			Message: "Device ID is required",
			Code:    400,
		})
		return
	}
	
	// Retrieve device using service
	device, err := deviceService.GetDevice(deviceId)
	if err != nil {
		logger.WithError(err).WithField("deviceId", deviceId).Error("Failed to retrieve device")
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Device not found",
			Message: err.Error(),
			Code:    404,
		})
		return
	}
	
	response := map[string]interface{}{
		"status": "found",
		"device": device,
	}
	
	c.JSON(http.StatusOK, response)
}

// ingestRawDataHandler handles raw device data ingestion
func ingestRawDataHandler(c *gin.Context) {
	var payload models.RawIngestPayload
	
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationError := middleware.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, validationError)
		return
	}
	
	// Validate the payload
	if err := validate.Struct(&payload); err != nil {
		validationError := middleware.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, validationError)
		return
	}
	
	// Store raw data using service
	ingestion, err := deviceService.ProcessRawDataIngestion(&payload)
	if err != nil {
		logger.WithError(err).Error("Failed to process raw data ingestion")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Message: "Failed to ingest raw data",
			Code:    500,
		})
		return
	}
	
	// Transform raw data to EPCIS events
	events, err := epcisService.TransformRawData(&payload)
	if err != nil {
		logger.WithError(err).Error("Failed to transform raw data to EPCIS events")
		// Don't fail the request, just log the error
	} else {
		// Create EPCIS events
		var createdEvents []string
		for _, event := range events {
			dbEvent, eventErr := epcisService.CreateEvent(event)
			if eventErr != nil {
				logger.WithError(eventErr).Error("Failed to create EPCIS event from raw data")
				continue
			}
			createdEvents = append(createdEvents, dbEvent.ID)
		}
		
		logger.WithFields(logrus.Fields{
			"ingestionId":  ingestion.ID,
			"deviceType":   payload.DeviceType,
			"deviceId":     payload.DeviceID,
			"eventCount":   len(createdEvents),
			"createdEvents": createdEvents,
		}).Info("Raw data processed and EPCIS events created")
	}
	
	logger.WithFields(logrus.Fields{
		"deviceType": payload.DeviceType,
		"deviceId":   payload.DeviceID,
		"timestamp":  payload.Timestamp,
	}).Info("Raw data ingested")
	
	response := map[string]interface{}{
		"status":      "ingested",
		"ingestionId": ingestion.ID,
		"payload":     payload,
	}
	
	c.JSON(http.StatusAccepted, response)
}

// claimDeviceHandler handles device claiming with claim codes
func claimDeviceHandler(c *gin.Context) {
	var claimCode models.ClaimCode
	
	if err := c.ShouldBindJSON(&claimCode); err != nil {
		validationError := middleware.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, validationError)
		return
	}
	
	// Validate the claim code
	if err := validate.Struct(&claimCode); err != nil {
		validationError := middleware.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, validationError)
		return
	}
	
	// Claim device using service
	device, err := deviceService.ClaimDevice(&claimCode)
	if err != nil {
		logger.WithError(err).Error("Failed to claim device")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Claim failed",
			Message: err.Error(),
			Code:    400,
		})
		return
	}
	
	logger.WithFields(logrus.Fields{
		"deviceId":  device.DeviceID,
		"claimCode": claimCode.ClaimCode,
		"type":      claimCode.Type,
	}).Info("Device claim attempted")
	
	response := map[string]interface{}{
		"status":   "claimed",
		"deviceId": device.DeviceID,
		"device":   device,
	}
	
	c.JSON(http.StatusOK, response)
}

func setupRoutes(r *gin.Engine) {
	// Add global middleware
	r.Use(middleware.ContentTypeMiddleware())
	r.Use(middleware.RequestSizeLimitMiddleware(10 * 1024 * 1024)) // 10MB limit
	
	// Health check endpoint
	r.GET("/health", healthHandler)
	
	// API info endpoint
	r.GET("/api", apiHandler)

	// API routes group
	api := r.Group("/api")
	{
		// EPCIS Events
		api.POST("/events", createEventHandler)
		api.GET("/events/:id", getEventHandler)
		
		// Device Management
		api.POST("/devices", registerDeviceHandler)
		api.GET("/devices/:deviceId", getDeviceHandler)
		
		// Data Ingestion
		api.POST("/ingest", ingestRawDataHandler)
		
		// Device Claiming
		api.POST("/claim", claimDeviceHandler)
	}
}

func main() {
	// Set Gin mode
	if os.Getenv("NODE_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.New()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	})

	// Setup routes
	setupRoutes(r)

	// Get port and host from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	// Start server
	logger.WithFields(logrus.Fields{
		"port":    port,
		"host":    host,
		"nodeEnv": os.Getenv("NODE_ENV"),
	}).Info("Scain EPCIS backend server started")

	// Graceful shutdown
	go func() {
		if err := r.Run(host + ":" + port); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	logger.Info("Shutting down server...")
	// Add graceful shutdown logic here if needed
	logger.Info("Server shutdown complete")
} 