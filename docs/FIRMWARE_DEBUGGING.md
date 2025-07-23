# Firmware-to-Backend Communication Debugging Guide

## Overview
This guide provides comprehensive debugging techniques for troubleshooting communication between IoT firmware devices and the Scain backend. It covers common issues, debugging tools, and step-by-step troubleshooting procedures.

## Communication Flow

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Firmware  │───▶│   Network   │───▶│   Backend   │
│   Device    │    │   (WiFi/    │    │   API       │
│             │    │   Cellular/ │    │             │
│             │    │   LoRaWAN)  │    │             │
└─────────────┘    └─────────────┘    └─────────────┘
```

## Debugging Tools

### 1. Network Tools
- **Wireshark:** Packet capture and analysis
- **tcpdump:** Command-line packet capture
- **curl:** Manual API testing
- **Postman:** API testing and debugging
- **netcat:** Network connectivity testing

### 2. Firmware Debugging
- **Serial Monitor:** Real-time firmware output
- **PlatformIO:** Integrated debugging
- **GDB:** Advanced debugging
- **Logic Analyzer:** Hardware signal analysis

### 3. Backend Debugging
- **Backend Logs:** Application-level debugging
- **Database Inspection:** Data verification
- **API Testing:** Endpoint validation

## Step-by-Step Debugging Process

### Step 1: Verify Firmware Configuration

#### Check WiFi Configuration (ESP32)
```cpp
// In firmware/config.h
#define WIFI_SSID "your_ssid"
#define WIFI_PASSWORD "your_password"
#define BACKEND_URL "http://192.168.1.100:8081"
#define DEVICE_ID "esp32-sensor-001"
```

#### Verify Network Connectivity
```cpp
void testWiFiConnection() {
  Serial.println("Testing WiFi connection...");
  
  if (WiFi.status() == WL_CONNECTED) {
    Serial.println("WiFi connected successfully");
    Serial.print("IP Address: ");
    Serial.println(WiFi.localIP());
    Serial.print("Signal Strength: ");
    Serial.println(WiFi.RSSI());
  } else {
    Serial.println("WiFi connection failed");
  }
}
```

### Step 2: Test Network Connectivity

#### Ping Test
```bash
# From firmware device or same network
ping 192.168.1.100  # Backend server IP
```

#### Port Connectivity Test
```bash
# Test if backend port is reachable
nc -zv 192.168.1.100 8081
```

#### Manual API Test
```bash
# Test backend health endpoint
curl -X GET http://192.168.1.100:8081/health

# Expected response:
# {"status":"ok","timestamp":"2024-07-21T12:00:00Z","version":"1.0.0"}
```

### Step 3: Debug Firmware HTTP Requests

#### Add Debug Logging to Firmware
```cpp
void sendDataToBackend(float temperature, float humidity, float pressure) {
  // Create JSON payload
  String jsonPayload = "{\"deviceType\":\"ESP32\",";
  jsonPayload += "\"deviceId\":\"" + String(DEVICE_ID) + "\",";
  jsonPayload += "\"timestamp\":\"" + getISOTimestamp() + "\",";
  jsonPayload += "\"lotCode\":\"LOT123456\",";
  jsonPayload += "\"data\":{";
  jsonPayload += "\"temperature\":" + String(temperature) + ",";
  jsonPayload += "\"humidity\":" + String(humidity) + ",";
  jsonPayload += "\"pressure\":" + String(pressure);
  jsonPayload += "}}";
  
  // Debug: Print payload
  Serial.println("Sending payload:");
  Serial.println(jsonPayload);
  
  // Create HTTP client
  HTTPClient http;
  http.begin(BACKEND_URL + "/api/ingest");
  http.addHeader("Content-Type", "application/json");
  
  // Debug: Print request details
  Serial.println("Request URL: " + BACKEND_URL + "/api/ingest");
  Serial.println("Content-Type: application/json");
  
  // Send POST request
  int httpResponseCode = http.POST(jsonPayload);
  
  // Debug: Print response
  Serial.print("HTTP Response code: ");
  Serial.println(httpResponseCode);
  
  if (httpResponseCode > 0) {
    String response = http.getString();
    Serial.println("Response: " + response);
  } else {
    Serial.print("Error code: ");
    Serial.println(httpResponseCode);
    Serial.println("Error: " + http.errorToString(httpResponseCode));
  }
  
  http.end();
}
```

### Step 4: Monitor Backend Logs

#### Enable Debug Logging
```go
// In backend/main.go
import (
    "github.com/sirupsen/logrus"
)

func init() {
    // Set log level to debug
    logrus.SetLevel(logrus.DebugLevel)
    logrus.SetFormatter(&logrus.JSONFormatter{})
}
```

#### Check Backend Logs
```bash
# Start backend with debug logging
cd backend
go run main.go

# Look for incoming requests in logs
# Example log output:
# {"level":"info","msg":"Incoming request","method":"POST","path":"/api/ingest","time":"2024-07-21T12:00:00Z"}
```

### Step 5: Database Verification

#### Check Raw Data Ingestion
```sql
-- Connect to SQLite database
sqlite3 scain.db

-- Check raw data table
SELECT * FROM raw_data_ingestions ORDER BY created_at DESC LIMIT 5;

-- Check events table
SELECT id, event_type, device_id, created_at FROM events ORDER BY created_at DESC LIMIT 5;
```

#### Verify Device Registration
```sql
-- Check if device exists
SELECT device_id, type, is_active FROM devices WHERE device_id = 'esp32-sensor-001';
```

## Common Issues and Solutions

### Issue 1: WiFi Connection Fails

**Symptoms:**
- Firmware stuck in WiFi connection loop
- Serial output shows "WiFi connection failed"

**Debugging Steps:**
```cpp
void debugWiFi() {
  Serial.println("=== WiFi Debug Info ===");
  Serial.print("SSID: ");
  Serial.println(WIFI_SSID);
  Serial.print("Password length: ");
  Serial.println(strlen(WIFI_PASSWORD));
  
  // Scan available networks
  int n = WiFi.scanNetworks();
  Serial.print("Found ");
  Serial.print(n);
  Serial.println(" networks:");
  
  for (int i = 0; i < n; ++i) {
    Serial.print(i + 1);
    Serial.print(": ");
    Serial.print(WiFi.SSID(i));
    Serial.print(" (");
    Serial.print(WiFi.RSSI(i));
    Serial.print(")");
    Serial.println((WiFi.encryptionType(i) == WIFI_AUTH_OPEN)?" ":"*");
  }
}
```

**Solutions:**
- Verify SSID and password
- Check WiFi signal strength
- Ensure 2.4GHz compatibility
- Try different WiFi channels

### Issue 2: HTTP Request Fails

**Symptoms:**
- HTTP response code 0 or negative
- "Connection refused" errors
- Timeout errors

**Debugging Steps:**
```cpp
void debugHTTPRequest() {
  HTTPClient http;
  
  // Test with simple GET request first
  http.begin("http://192.168.1.100:8081/health");
  int responseCode = http.GET();
  
  Serial.print("Health check response: ");
  Serial.println(responseCode);
  
  if (responseCode > 0) {
    String response = http.getString();
    Serial.println("Health response: " + response);
  }
  
  http.end();
}
```

**Solutions:**
- Verify backend server is running
- Check firewall settings
- Ensure correct IP address and port
- Test with curl from same network

### Issue 3: JSON Parsing Errors

**Symptoms:**
- Backend returns 400 Bad Request
- "Invalid JSON" errors
- Malformed payload errors

**Debugging Steps:**
```cpp
void validateJSON(String json) {
  // Use ArduinoJson to validate
  DynamicJsonDocument doc(1024);
  DeserializationError error = deserializeJson(doc, json);
  
  if (error) {
    Serial.print("JSON parsing failed: ");
    Serial.println(error.c_str());
    Serial.println("JSON string: " + json);
  } else {
    Serial.println("JSON is valid");
  }
}
```

**Solutions:**
- Validate JSON structure
- Check for special characters
- Ensure proper escaping
- Use JSON validation tools

### Issue 4: Device Not Registered

**Symptoms:**
- Backend returns "Device not found"
- 404 errors on device endpoints

**Debugging Steps:**
```bash
# Check if device exists in database
sqlite3 scain.db "SELECT * FROM devices WHERE device_id = 'esp32-sensor-001';"

# Register device manually if needed
curl -X POST http://localhost:8081/api/devices \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "esp32-sensor-001",
    "type": "ESP32",
    "firmwareVersion": "1.0.0"
  }'
```

**Solutions:**
- Register device before sending data
- Check device ID consistency
- Verify device registration API

### Issue 5: Data Not Appearing in Database

**Symptoms:**
- HTTP 200 response but no data in database
- Missing events or raw data

**Debugging Steps:**
```bash
# Check backend logs for database errors
tail -f backend.log | grep -i "database\|error"

# Verify database file permissions
ls -la scain.db

# Check database integrity
sqlite3 scain.db "PRAGMA integrity_check;"
```

**Solutions:**
- Check database file permissions
- Verify database connection
- Check for transaction rollbacks
- Monitor database logs

## Advanced Debugging Techniques

### 1. Packet Capture with Wireshark

```bash
# Capture HTTP traffic
sudo wireshark -i eth0 -k -f "port 8081"

# Filter for specific device
# http.request.method == "POST" && http.host contains "192.168.1.100"
```

### 2. Firmware Memory Debugging

```cpp
void debugMemory() {
  Serial.print("Free heap: ");
  Serial.println(ESP.getFreeHeap());
  Serial.print("Largest free block: ");
  Serial.println(ESP.getMaxAllocHeap());
  Serial.print("Free stack: ");
  Serial.println(uxTaskGetStackHighWaterMark(NULL));
}
```

### 3. Backend Performance Monitoring

```go
// Add timing middleware
func timingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        duration := time.Since(start)
        logrus.WithFields(logrus.Fields{
            "method":     c.Request.Method,
            "path":       c.Request.URL.Path,
            "duration":   duration,
            "status":     c.Writer.Status(),
        }).Info("Request processed")
    }
}
```

## Testing Framework

### 1. Automated Testing Script

```bash
#!/bin/bash
# test_firmware_communication.sh

echo "Testing firmware-to-backend communication..."

# Test 1: Health check
echo "Test 1: Backend health check"
curl -s http://localhost:8081/health | jq .

# Test 2: Device registration
echo "Test 2: Device registration"
curl -s -X POST http://localhost:8081/api/devices \
  -H "Content-Type: application/json" \
  -d '{"deviceId":"test-device","type":"ESP32"}' | jq .

# Test 3: Data ingestion
echo "Test 3: Data ingestion"
curl -s -X POST http://localhost:8081/api/ingest \
  -H "Content-Type: application/json" \
  -d '{
    "deviceType": "ESP32",
    "deviceId": "test-device",
    "timestamp": "2024-07-21T12:00:00Z",
    "data": {"temperature": 25.0}
  }' | jq .

# Test 4: Verify data in database
echo "Test 4: Database verification"
sqlite3 scain.db "SELECT COUNT(*) as event_count FROM events;"
```

### 2. Firmware Test Mode

```cpp
void runDiagnostics() {
  Serial.println("=== Firmware Diagnostics ===");
  
  // Test WiFi
  testWiFiConnection();
  
  // Test sensors
  testSensors();
  
  // Test HTTP communication
  testHTTPCommunication();
  
  // Test JSON generation
  testJSONGeneration();
  
  Serial.println("=== Diagnostics Complete ===");
}
```

## Best Practices

### 1. Error Handling
- Implement retry logic with exponential backoff
- Log all errors with context
- Provide meaningful error messages
- Handle network timeouts gracefully

### 2. Data Validation
- Validate sensor data ranges
- Check for NaN or infinite values
- Verify timestamp accuracy
- Ensure JSON structure integrity

### 3. Monitoring
- Monitor device connectivity
- Track API response times
- Monitor database performance
- Set up alerting for failures

### 4. Security
- Use HTTPS in production
- Implement device authentication
- Validate input data
- Monitor for suspicious activity

## Troubleshooting Checklist

- [ ] Firmware compiles without errors
- [ ] WiFi credentials are correct
- [ ] Backend server is running
- [ ] Network connectivity is established
- [ ] Device is registered in backend
- [ ] JSON payload is valid
- [ ] HTTP request reaches backend
- [ ] Database operations succeed
- [ ] Response is received by firmware
- [ ] Error handling works correctly

## Resources

- [ESP32 Arduino Core Documentation](https://github.com/espressif/arduino-esp32)
- [PlatformIO Documentation](https://docs.platformio.org/)
- [HTTPClient Reference](https://docs.espressif.com/projects/arduino-esp32/en/latest/api/httpclient.html)
- [ArduinoJson Library](https://arduinojson.org/)
- [Scain Backend API Documentation](./API_TEST_REPORT.md) 