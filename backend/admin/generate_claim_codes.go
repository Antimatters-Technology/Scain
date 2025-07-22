package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"scain-backend/database"
	"scain-backend/models"
	"scain-backend/services"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run admin/generate_claim_codes.go <device_type> <count> [expires_hours]")
		fmt.Println("Example: go run admin/generate_claim_codes.go ESP32 5 24")
		os.Exit(1)
	}

	deviceType := os.Args[1]
	count, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal("Invalid count:", err)
	}

	expiresHours := 0
	if len(os.Args) > 3 {
		expiresHours, err = strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal("Invalid expires_hours:", err)
		}
	}

	// Initialize database
	if err := database.InitDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Create device service
	deviceService := services.NewDeviceService()

	// Generate claim codes
	claimCodes, err := deviceService.GenerateClaimCodes(models.DeviceType(deviceType), count, expiresHours)
	if err != nil {
		log.Fatal("Failed to generate claim codes:", err)
	}

	fmt.Printf("\n✅ Generated %d claim codes for device type '%s':\n\n", len(claimCodes), deviceType)
	
	for i, code := range claimCodes {
		expiry := "Never"
		if code.ExpiresAt != nil {
			expiry = code.ExpiresAt.Format("2006-01-02 15:04:05")
		}
		fmt.Printf("%d. %s (expires: %s)\n", i+1, code.ClaimCode, expiry)
	}
	
	fmt.Printf("\n💡 Use these codes with: curl -X POST /api/claim -d '{\"claimCode\":\"<CODE>\",\"type\":\"%s\"}'\n", deviceType)
} 