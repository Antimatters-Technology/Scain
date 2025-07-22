#!/bin/bash

# Scain Backend API Test Script
# This script tests all endpoints of the Scain backend

set -e

BASE_URL="http://localhost:8081"
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Scain Backend API Test Suite${NC}"
echo "=================================="

# Function to make HTTP requests and show results
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo -e "\n${YELLOW}Testing:${NC} $description"
    echo -e "${BLUE}$method $endpoint${NC}"
    
    if [ -n "$data" ]; then
        echo -e "${BLUE}Data:${NC} $data"
        response=$(curl -s -X $method "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    else
        response=$(curl -s -X $method "$BASE_URL$endpoint")
    fi
    
    echo -e "${GREEN}Response:${NC}"
    echo "$response" | jq . 2>/dev/null || echo "$response"
    
    return 0
}

# Test 1: Health Check
test_endpoint "GET" "/health" "" "Health Check"

# Test 2: API Info  
test_endpoint "GET" "/api" "" "API Information"

# Test 3: Register Device
DEVICE_DATA='{"deviceId": "test-esp32-001", "type": "ESP32", "firmwareVersion": "1.0.0", "secureBoot": true}'
test_endpoint "POST" "/api/devices" "$DEVICE_DATA" "Device Registration"

# Test 4: Get Device
test_endpoint "GET" "/api/devices/test-esp32-001" "" "Device Retrieval"

# Test 5: Create EPCIS Event
EVENT_DATA='{
  "eventType": "ObjectEvent",
  "eventTime": "2024-07-21T12:00:00Z", 
  "eventTimeZoneOffset": "+00:00",
  "bizStep": "harvesting",
  "epcList": ["urn:epc:id:sgtin:123456.789012.001"],
  "deviceId": "test-esp32-001"
}'
RESPONSE=$(test_endpoint "POST" "/api/events" "$EVENT_DATA" "EPCIS Event Creation")
EVENT_ID=$(echo "$RESPONSE" | jq -r '.eventId' 2>/dev/null || echo "")

# Test 6: Get EPCIS Event
if [ -n "$EVENT_ID" ] && [ "$EVENT_ID" != "null" ]; then
    test_endpoint "GET" "/api/events/$EVENT_ID" "" "EPCIS Event Retrieval"
fi

# Test 7: Raw Data Ingestion
RAW_DATA='{
  "deviceType": "ESP32",
  "deviceId": "test-esp32-001", 
  "timestamp": "2024-07-21T12:00:00Z",
  "lotCode": "LOT123456",
  "data": {
    "temperature": 24.5,
    "humidity": 62.3,
    "pressure": 1013.25
  }
}'
test_endpoint "POST" "/api/ingest" "$RAW_DATA" "Raw Data Ingestion"

# Test 8: Invalid Device Claim (should fail)
CLAIM_DATA='{"claimCode": "INVALID1", "type": "ESP32"}'
test_endpoint "POST" "/api/claim" "$CLAIM_DATA" "Invalid Device Claim (should fail)"

# Test 9: Non-existent Device (should fail)
test_endpoint "GET" "/api/devices/nonexistent" "" "Non-existent Device (should fail)"

# Test 10: Invalid Event ID (should fail)
test_endpoint "GET" "/api/events/invalid-id" "" "Invalid Event ID (should fail)"

echo -e "\n${GREEN}✅ All tests completed!${NC}"
echo -e "${BLUE}💡 Check the responses above for any errors${NC}"

# Show database statistics
echo -e "\n${YELLOW}📊 Database Statistics:${NC}"
if [ -f "scain.db" ]; then
    echo "Devices: $(sqlite3 scain.db 'SELECT COUNT(*) FROM devices;')"
    echo "Events: $(sqlite3 scain.db 'SELECT COUNT(*) FROM events;')"  
    echo "Raw Ingestions: $(sqlite3 scain.db 'SELECT COUNT(*) FROM raw_data_ingestions;')"
    echo "Claim Codes: $(sqlite3 scain.db 'SELECT COUNT(*) FROM claim_code_entries;')"
else
    echo "Database file not found"
fi 