# Device Ingestion Payloads

## Overview
This document provides examples of payloads that IoT devices send to the Scain backend via the `/api/ingest` endpoint. Each device type has specific data formats and requirements.

## API Endpoint
```
POST /api/ingest
Content-Type: application/json
```

## Device Types

### 1. ESP32 Device Payload

**Device Type:** `ESP32`

**Purpose:** Environmental monitoring with temperature, humidity, and pressure sensors.

**Example Payload:**
```json
{
  "deviceType": "ESP32",
  "deviceId": "esp32-sensor-001",
  "timestamp": "2024-07-21T12:00:00Z",
  "lotCode": "LOT123456",
  "data": {
    "temperature": 24.5,
    "humidity": 62.3,
    "pressure": 1013.25,
    "battery": 85.2,
    "rssi": -45
  },
  "metadata": {
    "location": "Warehouse A",
    "zone": "Cold Storage",
    "sensorCalibration": "2024-06-15"
  }
}
```

**Data Fields:**
- `temperature`: Float (Celsius)
- `humidity`: Float (Percentage)
- `pressure`: Float (hPa)
- `battery`: Float (Percentage, optional)
- `rssi`: Integer (Signal strength, optional)

**EPCIS Event Generated:** ObjectEvent with sensor data

---

### 2. ExpressLink Device Payload

**Device Type:** `ExpressLink`

**Purpose:** AWS IoT ExpressLink integration for cloud connectivity.

**Example Payload:**
```json
{
  "deviceType": "ExpressLink",
  "deviceId": "expresslink-001",
  "timestamp": "2024-07-21T12:00:00Z",
  "lotCode": "LOT789012",
  "data": {
    "temperature": 22.1,
    "humidity": 58.7,
    "awsStatus": "connected",
    "messageCount": 145,
    "lastSync": "2024-07-21T11:59:30Z"
  },
  "metadata": {
    "awsThingName": "scain-expresslink-001",
    "region": "us-east-1",
    "certificateExpiry": "2025-01-15"
  }
}
```

**Data Fields:**
- `temperature`: Float (Celsius)
- `humidity`: Float (Percentage)
- `awsStatus`: String (connected/disconnected)
- `messageCount`: Integer (Total messages sent)
- `lastSync`: String (ISO timestamp)

**EPCIS Event Generated:** ObjectEvent with AWS integration data

---

### 3. LoRaWAN Device Payload

**Device Type:** `LoRaWAN`

**Purpose:** Long-range wireless communication for remote monitoring.

**Example Payload:**
```json
{
  "deviceType": "LoRaWAN",
  "deviceId": "lorawan-001",
  "timestamp": "2024-07-21T12:00:00Z",
  "lotCode": "LOT345678",
  "data": {
    "temperature": 18.3,
    "humidity": 45.2,
    "spreadingFactor": 7,
    "frequency": 868.1,
    "batteryLevel": 92,
    "signalQuality": -65
  },
  "metadata": {
    "gatewayId": "gateway-eu-001",
    "networkId": "thethingsnetwork",
    "devEUI": "A8404194A1B1C1D1"
  }
}
```

**Data Fields:**
- `temperature`: Float (Celsius)
- `humidity`: Float (Percentage)
- `spreadingFactor`: Integer (LoRa spreading factor)
- `frequency`: Float (MHz)
- `batteryLevel`: Integer (Percentage)
- `signalQuality`: Integer (dBm)

**EPCIS Event Generated:** ObjectEvent with LoRaWAN network data

---

### 4. GPS Tracker Payload

**Device Type:** `Tracker`

**Purpose:** Location tracking for shipments and logistics.

**Example Payload:**
```json
{
  "deviceType": "Tracker",
  "deviceId": "gps-tracker-001",
  "timestamp": "2024-07-21T12:00:00Z",
  "lotCode": "LOT901234",
  "data": {
    "latitude": 40.7128,
    "longitude": -74.0060,
    "altitude": 10.5,
    "speed": 25.3,
    "heading": 180,
    "satellites": 8,
    "accuracy": 3.2
  },
  "metadata": {
    "vehicleId": "TRUCK-001",
    "driverId": "DRIVER-123",
    "routeId": "ROUTE-NYC-CHI"
  }
}
```

**Data Fields:**
- `latitude`: Float (Decimal degrees)
- `longitude`: Float (Decimal degrees)
- `altitude`: Float (Meters, optional)
- `speed`: Float (km/h, optional)
- `heading`: Integer (Degrees, optional)
- `satellites`: Integer (Number of GPS satellites)
- `accuracy`: Float (Meters, optional)

**EPCIS Event Generated:** ObjectEvent with location data and ReadPoint

---

### 5. ERP System Payload

**Device Type:** `ERP`

**Purpose:** Business process integration from ERP systems.

**Example Payload:**
```json
{
  "deviceType": "ERP",
  "deviceId": "erp-system-001",
  "timestamp": "2024-07-21T12:00:00Z",
  "lotCode": "LOT567890",
  "data": {
    "businessStep": "shipping",
    "disposition": "in_transit",
    "orderNumber": "ORD-2024-001",
    "customerId": "CUST-123",
    "quantity": 100,
    "unitPrice": 15.50,
    "totalValue": 1550.00
  },
  "metadata": {
    "erpSystem": "SAP",
    "module": "MM",
    "user": "operator-001",
    "sessionId": "SESS-456"
  }
}
```

**Data Fields:**
- `businessStep`: String (EPCIS business step)
- `disposition`: String (EPCIS disposition)
- `orderNumber`: String (Business order reference)
- `customerId`: String (Customer identifier)
- `quantity`: Integer (Product quantity)
- `unitPrice`: Float (Price per unit)
- `totalValue`: Float (Total order value)

**EPCIS Event Generated:** TransactionEvent with business transactions

---

## Response Format

**Success Response (200 OK):**
```json
{
  "ingestionId": "7a582534-56dc-4fc2-9565-ce61a3fe3215",
  "payload": {
    "deviceType": "ESP32",
    "deviceId": "esp32-sensor-001",
    "timestamp": "2024-07-21T12:00:00Z",
    "lotCode": "LOT123456",
    "data": {
      "temperature": 24.5,
      "humidity": 62.3,
      "pressure": 1013.25
    }
  },
  "status": "ingested"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Validation failed",
  "message": "Invalid device type: UNKNOWN",
  "code": 400
}
```

## Validation Rules

### Required Fields
- `deviceType`: Must be one of: ESP32, ExpressLink, LoRaWAN, Tracker, ERP
- `deviceId`: String, must be registered device
- `timestamp`: ISO 8601 format
- `data`: Object with device-specific fields

### Optional Fields
- `lotCode`: String (for lot tracking)
- `metadata`: Object (additional context)

### Data Type Validation
- Numeric fields must be numbers
- Timestamps must be valid ISO 8601 format
- Device IDs must exist in the device registry

## EPCIS Transformation

Each payload is automatically transformed into EPCIS events:

1. **ESP32/ExpressLink/LoRaWAN:** ObjectEvent with sensor data
2. **Tracker:** ObjectEvent with location data and ReadPoint
3. **ERP:** TransactionEvent with business transactions

## Testing Examples

### cURL Commands

**ESP32 Device:**
```bash
curl -X POST http://localhost:8081/api/ingest \
  -H "Content-Type: application/json" \
  -d '{
    "deviceType": "ESP32",
    "deviceId": "esp32-sensor-001",
    "timestamp": "2024-07-21T12:00:00Z",
    "lotCode": "LOT123456",
    "data": {
      "temperature": 24.5,
      "humidity": 62.3,
      "pressure": 1013.25
    }
  }'
```

**GPS Tracker:**
```bash
curl -X POST http://localhost:8081/api/ingest \
  -H "Content-Type: application/json" \
  -d '{
    "deviceType": "Tracker",
    "deviceId": "gps-tracker-001",
    "timestamp": "2024-07-21T12:00:00Z",
    "lotCode": "LOT901234",
    "data": {
      "latitude": 40.7128,
      "longitude": -74.0060,
      "speed": 25.3
    }
  }'
```

## Best Practices

1. **Timestamp Accuracy:** Use precise timestamps for accurate event tracking
2. **Data Validation:** Validate data ranges before sending (e.g., temperature -40 to 80°C)
3. **Error Handling:** Implement retry logic for failed requests
4. **Batch Processing:** Consider batching multiple readings for efficiency
5. **Security:** Use HTTPS in production environments
6. **Monitoring:** Track ingestion success rates and response times 