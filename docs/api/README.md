# API Documentation

This document describes the REST API endpoints provided by the Scain **Go** backend service.

## ✅ Implementation Status

**Backend is PRODUCTION READY** - Complete EPCIS implementation with blockchain integration, device management, and comprehensive API.

### ✅ What's Working
- **Go backend** with Gin HTTP framework
- **SQLite database** with full CRUD operations
- **Hyperledger Fabric integration** for blockchain anchoring
- Health check endpoint (`/health`)  
- API information endpoint (`/api`)
- **EPCIS event creation** (`POST /api/events`) with blockchain anchoring
- **Event retrieval** (`GET /api/events/:id`)
- **Device registration** (`POST /api/devices`)
- **Device retrieval** (`GET /api/devices/:id`)
- **Raw data ingestion** (`POST /api/ingest`) with auto-transformation
- **Device claiming** (`POST /api/claim`)
- **Blockchain verification** (`GET /api/events/:id/verify`)
- **Transaction history** (`GET /api/events/:id/history`)
- **Comprehensive data models** with validation
- **Hash utilities** for data integrity
- **Canonical JSON** serialization

### 🆕 New Features
- **Blockchain Integration**: Events automatically anchored on Hyperledger Fabric
- **Event Verification**: Cryptographic verification against blockchain
- **Transaction History**: Complete audit trail from blockchain
- **Configurable Blockchain**: Enable/disable via environment variables

---

## Base URL

- **Development**: `http://localhost:8081`
- **Production**: `https://api.scain.com` (when deployed)

## Authentication

**Currently**: No authentication required (development mode)

**Planned**: JWT-based authentication with role-based access control

```bash
# Future authentication headers
Authorization: Bearer <jwt-token>
X-API-Key: <api-key>
```

## Response Format

All API responses follow this standard format:

### Success Response
```json
{
  "status": "created|ok",
  "data": { ... },
  "hash": "sha256-hash" // For EPCIS events
}
```

### Error Response
```json
{
  "error": "error-type",
  "message": "Human readable error message",
  "code": 400
  "data": {},
  "message": "Optional message",
  "timestamp": "2025-07-20T07:29:21.139Z"
}
```

---

## API Endpoints

### Health & System

#### GET `/health`

Returns the health status of the API server.

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2025-07-20T07:29:21.139Z",
  "version": "1.0.0",
  "blockchain": "enabled" // or "disabled"
}
```

**Status Codes:**
- `200` - Server is healthy

#### GET `/api`

Returns API information and available endpoints.

**Response:**
```json
{
  "name": "Scain EPCIS API",
  "version": "1.0.0",
  "endpoints": [...],
  "blockchain_enabled": true
}
```

### EPCIS Events

#### POST `/api/events`

Create a new EPCIS event. Event will be stored in database and optionally anchored on blockchain.

**Request Body:**
```json
{
  "eventType": "ObjectEvent",
  "eventTime": "2024-01-01T12:00:00Z",
  "eventTimeZoneOffset": "+00:00",
  "action": "OBSERVE",
  "epcList": ["urn:epc:id:sgtin:0614141.107346.2018"],
  "bizStep": "urn:epcglobal:cbv:bizstep:harvesting",
  "disposition": "urn:epcglobal:cbv:disp:in_progress"
}
```

**Response:**
```json
{
  "status": "created",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "eventType": "ObjectEvent",
    "hash": "abc123...",
    "blockchainTxId": "tx123..." // if blockchain enabled
  },
  "hash": "abc123..."
}
```

**Status Codes:**
- `201` - Event created successfully
- `400` - Invalid request data
- `500` - Server error

#### GET `/api/events/:id`

Retrieve an EPCIS event by ID.

**Response:**
```json
{
  "status": "ok",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "eventType": "ObjectEvent",
    "eventTime": "2024-01-01T12:00:00Z",
    "hash": "abc123...",
    "blockchainTxId": "tx123...",
    "epcList": ["urn:epc:id:sgtin:0614141.107346.2018"]
  }
}
```

**Status Codes:**
- `200` - Event retrieved successfully
- `404` - Event not found

### Blockchain Operations

#### GET `/api/events/:id/verify`

Verify an event's integrity against the blockchain.

**Response:**
```json
{
  "status": "ok",
  "data": {
    "eventId": "123e4567-e89b-12d3-a456-426614174000",
    "isValid": true,
    "blockchainTxId": "tx123...",
    "verificationTime": "2024-01-01T12:05:00Z"
  }
}
```

#### GET `/api/events/:id/history`

Get the blockchain transaction history for an event.

**Response:**
```json
{
  "status": "ok",
  "data": [
    {
      "txId": "tx123...",
      "timestamp": "2024-01-01T12:00:00Z",
      "isDelete": false,
      "value": {...}
    }
  ]
}
```

### Device Management

#### POST `/api/devices`

Register a new device.

**Request Body:**
```json
{
  "deviceId": "esp32-001",
  "deviceType": "ESP32",
  "metadata": {
    "location": "Warehouse A",
    "sensors": ["temperature", "humidity"]
  }
}
```

**Response:**
```json
{
  "status": "created",
  "data": {
    "id": "device-uuid-123",
    "deviceId": "esp32-001",
    "deviceType": "ESP32",
    "isActive": true,
    "registeredAt": "2024-01-01T12:00:00Z"
  }
}
```

#### GET `/api/devices/:id`

Retrieve device information.

**Response:**
```json
{
  "status": "ok",
  "data": {
    "id": "device-uuid-123",
    "deviceId": "esp32-001",
    "deviceType": "ESP32",
    "lastHeartbeat": "2024-01-01T12:00:00Z",
    "isActive": true
  }
}
```

### Data Ingestion

#### POST `/api/ingest`

Ingest raw sensor data and automatically transform to EPCIS events.

**Request Body:**
```json
{
  "deviceType": "ESP32",
  "deviceId": "esp32-001",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "temperature": 25.5,
    "humidity": 60.2,
    "location": {
      "lat": 40.7128,
      "lng": -74.0060
    }
  }
}
```

**Response:**
```json
{
  "status": "created",
  "data": {
    "eventsCreated": 1,
    "eventIds": ["event-uuid-123"],
    "transformationType": "ESP32_TO_OBJECT_EVENT"
  }
}
```

#### POST `/api/claim`

Claim a device using a validation code.

**Request Body:**
```json
{
  "deviceId": "esp32-001",
  "claimCode": "ABC12345"
}
```

**Response:**
```json
{
  "status": "ok",
  "data": {
    "deviceId": "esp32-001",
    "claimed": true,
    "claimedAt": "2024-01-01T12:00:00Z"
  }
}
    "GET /api/devices/{deviceId} - Get device info",
    "POST /api/ingest - Raw device data ingestion",
    "POST /api/claim - Claim device with code"
  ]
}
```

**Status Codes:**
- `200` - API information retrieved successfully

---

## EPCIS & Device Management Endpoints (Go Implementation)

### Device Management

#### POST `/api/devices`
**Status**: ✅ Implemented (Go)

Register a new IoT device with comprehensive metadata.

**Request Body:**
```json
{
  "deviceId": "ESP32-001",
  "type": "ESP32",
  "secureBoot": true,
  "otaCapable": true,
  "firmwareVersion": "1.2.0",
  "calibrationDate": "2025-01-15T00:00:00Z",
  "regulatoryCerts": ["FCC", "CE"],
  "batteryPct": 95,
  "joinStatus": "joined"
}
```

**Response:**
```json
{
  "status": "registered",
  "device": {
    "deviceId": "ESP32-001",
    "type": "ESP32",
    "secureBoot": true,
    "otaCapable": true,
    "firmwareVersion": "1.2.0",
    "calibrationDate": "2025-01-15T00:00:00Z",
    "regulatoryCerts": ["FCC", "CE"],
    "batteryPct": 95,
    "joinStatus": "joined"
  }
}
```

**Status Codes:**
- `201` - Device registered successfully
- `400` - Invalid request body or validation error
- `500` - Server error

#### GET `/api/devices/:deviceId`
**Status**: ❌ Not implemented (database needed)

Get information about a specific device.

**Parameters:**
- `deviceId` - The unique device identifier

**Response:**
```json
{
  "error": "Not implemented",
  "message": "Device retrieval not yet implemented",
  "code": 501
}
```

**Status Codes:**
- `501` - Not yet implemented
- `404` - Device not found (when implemented)
- `200` - Device retrieved successfully (when implemented)

### Device Claiming

#### POST `/api/claim`
**Status**: ✅ Implemented (Go)

Claim a device using a claim code.

**Request Body:**
```json
{
  "claimCode": "ABCD1234",
  "type": "ESP32",
  "identifiers": {
    "macAddress": "AA:BB:CC:DD:EE:FF"
  }
}
```

**Response:**
```json
{
  "status": "claimed",
  "claimCode": {
    "claimCode": "ABCD1234",
    "type": "ESP32",
    "identifiers": {
      "macAddress": "AA:BB:CC:DD:EE:FF"
    }
  }
}
```

**Status Codes:**
- `200` - Device claimed successfully
- `400` - Invalid claim code or validation error
- `404` - Claim code not found (when validation implemented)

### EPCIS Event Management

#### POST `/api/events`
**Status**: ✅ Implemented (Go)

Create a new EPCIS event for supply chain traceability.

**Request Body:**
```json
{
  "eventType": "ObjectEvent",
  "eventTime": "2025-07-20T07:29:21.139Z",
  "eventTimeZoneOffset": "-06:00",
  "bizStep": "shipping",
  "disposition": "in_transit",
  "readPoint": {
    "id": "urn:epc:id:sgln:0614141.00777.0"
  },
  "bizLocation": {
    "id": "urn:epc:id:sgln:0614141.00888.0"
  },
  "epcList": [
    "urn:epc:id:sgtin:0614141.812345.400"
  ],
  "sensorElementList": [
    {
      "sensorMetaData": {
        "deviceId": "ESP32-001",
        "deviceMetadata": {
          "type": "ESP32",
          "firmwareVersion": "1.2.0"
        }
      },
      "sensorReport": [
        {
          "type": "temperature",
          "value": 2.5,
          "uom": "CEL",
          "time": "2025-07-20T07:29:21.139Z"
        }
      ]
    }
  ],
  "lotCode": "LOT123456",
  "deviceId": "ESP32-001"
}
```

**Response:**
```json
{
  "status": "created",
  "hash": "a1b2c3d4e5f6...",
  "event": {
    "eventType": "ObjectEvent",
    "eventTime": "2025-07-20T07:29:21.139Z",
    "eventTimeZoneOffset": "-06:00",
    "bizStep": "shipping",
    "disposition": "in_transit",
    "readPoint": {
      "id": "urn:epc:id:sgln:0614141.00777.0"
    },
    "bizLocation": {
      "id": "urn:epc:id:sgln:0614141.00888.0"
    },
    "epcList": [
      "urn:epc:id:sgtin:0614141.812345.400"
    ],
    "sensorElementList": [
      {
        "sensorMetaData": {
          "deviceId": "ESP32-001",
          "deviceMetadata": {
            "type": "ESP32",
            "firmwareVersion": "1.2.0"
          }
        },
        "sensorReport": [
          {
            "type": "temperature",
            "value": 2.5,
            "uom": "CEL",
            "time": "2025-07-20T07:29:21.139Z"
          }
        ]
      }
    ],
    "lotCode": "LOT123456",
    "deviceId": "ESP32-001"
  }
}
```

**Status Codes:**
- `201` - Event created successfully
- `400` - Invalid request body or validation error
- `500` - Server error

#### GET `/api/events/:id`
**Status**: ❌ Not implemented (database needed)

Retrieve an EPCIS event by ID.

**Parameters:**
- `id` - The unique event identifier

**Response:**
```json
{
  "error": "Not implemented", 
  "message": "Event retrieval not yet implemented",
  "code": 501
}
```

**Status Codes:**
- `501` - Not yet implemented
- `404` - Event not found (when implemented)
- `200` - Event retrieved successfully (when implemented)

### Raw Data Ingestion

#### POST `/api/ingest`
**Status**: ✅ Implemented (Go)

Ingest raw device data for processing into EPCIS events.

**Request Body:**
```json
{
  "deviceType": "ESP32",
  "deviceId": "ESP32-001",
  "timestamp": "2025-07-20T07:29:21.139Z",
  "lotCode": "LOT123456",
  "data": {
    "temperature": 2.5,
    "humidity": 65.2,
    "location": {
      "lat": 40.7128,
      "lon": -74.0060
    }
  },
  "metadata": {
    "firmwareVersion": "1.2.0",
    "batteryLevel": 95
  }
}
```

**Response:**
```json
{
  "status": "ingested",
  "payload": {
    "deviceType": "ESP32",
    "deviceId": "ESP32-001", 
    "timestamp": "2025-07-20T07:29:21.139Z",
    "lotCode": "LOT123456",
    "data": {
      "temperature": 2.5,
      "humidity": 65.2,
      "location": {
        "lat": 40.7128,
        "lon": -74.0060
      }
    },
    "metadata": {
      "firmwareVersion": "1.2.0",
      "batteryLevel": 95
    }
  }
}
```

**Status Codes:**
- `202` - Data ingested successfully
- `400` - Invalid request body or validation error
- `500` - Server error

---

## Data Models

### EPCIS Event Types

The Go backend supports the following EPCIS event types:

**Request Body:**
```json
{
  "eventType": "ObjectEvent",
  "eventTime": "2025-07-20T07:29:21.139Z",
  "epcList": ["urn:epc:id:sgtin:0614141.Scain-001"],
  "bizStep": "cooling",
  "sensorElementList": [
    {
      "sensorMetaData": {
        "deviceId": "ESP32-001"
      },
      "sensorReport": [
        {
          "type": "gs1:Temperature",
          "value": 4.2,
          "uom": "CEL",
          "time": "2025-07-20T07:29:21.139Z"
        }
      ]
    }
  ]
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "eventId": "evt_123456",
    "hash": "sha256:abc123...",
    "blockchainTx": "tx_789def...",
    "timestamp": "2025-07-20T07:29:21.139Z"
  }
}
```

#### GET `/api/events`
**Status**: ❌ Not implemented

Query events with filters.

**Query Parameters:**
- `deviceId` - Filter by device
- `eventType` - Filter by event type
- `startDate` - Start date (ISO 8601)
- `endDate` - End date (ISO 8601)
- `limit` - Number of results
- `offset` - Pagination offset

**Response:**
```json
{
  "status": "success",
  "data": {
    "events": [
      {
        "id": "evt_123456",
        "eventType": "ObjectEvent",
        "eventTime": "2025-07-20T07:29:21.139Z",
        "deviceId": "ESP32-001",
        "hash": "sha256:abc123..."
      }
    ],
    "pagination": {
      "total": 1,
      "limit": 50,
      "offset": 0
    }
  }
}
```

### Traceability

#### GET `/api/trace/:epc`
**Status**: ❌ Not implemented

Trace a product through the supply chain.

**Path Parameters:**
- `epc` - Electronic Product Code

**Response:**
```json
{
  "status": "success",
  "data": {
    "epc": "urn:epc:id:sgtin:0614141.Scain-001",
    "trace": [
      {
        "eventId": "evt_123456",
        "eventType": "ObjectEvent",
        "eventTime": "2025-07-20T07:29:21.139Z",
        "bizStep": "cooling",
        "location": "Cold Storage Room A",
        "deviceId": "ESP32-001"
      }
    ],
    "summary": {
      "totalEvents": 1,
      "firstSeen": "2025-07-20T07:29:21.139Z",
      "lastSeen": "2025-07-20T07:29:21.139Z"
    }
  }
}
```

#### GET `/api/trace/lot/:lotCode`
**Status**: ❌ Not implemented

Get traceability data for a specific lot.

**Path Parameters:**
- `lotCode` - Lot identification code

**Response:**
```json
{
  "status": "success",
  "data": {
    "lotCode": "LOT-2024-001",
    "products": [
      {
        "epc": "urn:epc:id:sgtin:0614141.Scain-001",
        "status": "active"
      }
    ],
    "events": [
      {
        "eventId": "evt_123456",
        "eventType": "ObjectEvent",
        "eventTime": "2025-07-20T07:29:21.139Z"
      }
    ]
  }
}
```

### Analytics

#### GET `/api/analytics/dashboard`
**Status**: ❌ Not implemented

Get dashboard analytics data.

**Response:**
```json
{
  "status": "success",
  "data": {
    "kpis": {
      "totalDevices": 5,
      "activeDevices": 4,
      "totalEvents": 1250,
      "eventsToday": 45
    },
    "temperature": {
      "average": 4.2,
      "min": 2.1,
      "max": 8.5,
      "alerts": 2
    },
    "recentActivity": [
      {
        "deviceId": "ESP32-001",
        "eventType": "ObjectEvent",
        "timestamp": "2025-07-20T07:29:21.139Z"
      }
    ]
  }
}
```

#### GET `/api/analytics/devices/:deviceId`
**Status**: ❌ Not implemented

Get analytics for a specific device.

**Path Parameters:**
- `deviceId` - Device identifier

**Query Parameters:**
- `period` - Time period (1h, 24h, 7d, 30d)

**Response:**
```json
{
  "status": "success",
  "data": {
    "deviceId": "ESP32-001",
    "status": "active",
    "uptime": 99.5,
    "temperature": {
      "current": 4.2,
      "average": 4.1,
      "trend": "stable"
    },
    "events": {
      "total": 250,
      "last24h": 45
    }
  }
}
```

---

## Error Handling

### Standard Error Response

```json
{
  "status": "error",
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request data",
    "details": {}
  },
  "timestamp": "2025-07-20T07:29:21.139Z"
}
```

### Common Error Codes

| Code | Description | HTTP Status |
|------|-------------|-------------|
| `VALIDATION_ERROR` | Request data validation failed | 400 |
| `NOT_FOUND` | Resource not found | 404 |
| `UNAUTHORIZED` | Authentication required | 401 |
| `FORBIDDEN` | Insufficient permissions | 403 |
| `INTERNAL_ERROR` | Server internal error | 500 |
| `SERVICE_UNAVAILABLE` | Service temporarily unavailable | 503 |

---

## WebSocket Support

### Planned WebSocket Endpoints

**Status**: ❌ Not implemented

- `ws://localhost:8081/ws/events` - Real-time event streaming
- `ws://localhost:8081/ws/devices` - Device status updates
- `ws://localhost:8081/ws/notifications` - Alert notifications

### WebSocket Message Format

```json
{
  "type": "event|device|notification",
  "data": {},
  "timestamp": "2025-07-20T07:29:21.139Z"
}
```

---

## Implementation Roadmap

### Phase 1: Core APIs (Week 1-2)
- [ ] Database setup and models
- [ ] Device management endpoints
- [ ] Basic event ingestion
- [ ] Authentication system

### Phase 2: Business Logic (Week 3-4)
- [ ] EPCIS event processing
- [ ] Traceability queries
- [ ] Analytics endpoints
- [ ] Real-time WebSocket support

### Phase 3: Advanced Features (Week 5-6)
- [ ] Background job processing
- [ ] Blockchain integration
- [ ] Advanced analytics
- [ ] Export functionality

---

## Testing

### Current Testing Status
- **Tests**: ❌ No tests implemented
- **Coverage**: 0%
- **Test Framework**: Jest (configured but unused)

### Planned Test Structure
```
backend/src/
├── api/routes/__tests__/
│   ├── devices.test.ts
│   ├── events.test.ts
│   └── trace.test.ts
├── services/__tests__/
│   ├── deviceService.test.ts
│   └── eventService.test.ts
└── utils/__tests__/
    ├── hash.test.ts
    └── canonical.test.ts
```

---

## SDKs and Libraries

### JavaScript/TypeScript

```bash
npm install @scain/api-client
```

```typescript
import { ScainAPI } from '@scain/api-client';

const api = new ScainAPI({
  baseURL: 'http://localhost:8081',
  apiKey: 'your-api-key'
});

// Health check
const health = await api.health.check();

// Device management (when implemented)
const devices = await api.devices.list();
const device = await api.devices.create({
  deviceId: 'ESP32-001',
  type: 'ESP32'
});
```

---

## Support

For API support and questions:

- **Documentation**: [GitHub Wiki](https://github.com/your-org/scain/wiki)
- **Issues**: [GitHub Issues](https://github.com/your-org/scain/issues)
- **Email**: api-support@scain.com

---

**Last Updated**: July 2025  
**API Version**: 1.0.0  
**Implementation Status**: 5% Complete  
**Priority**: Core API endpoints needed 