# 🔌 API Reference

This document covers the primary interfaces exposed by Scain.

## 1. MQTT Topics

| Topic | Direction | QoS | Payload |
|-------|-----------|-----|---------|
| `scain/events` | publish | 1 | EPCIS 2.0 JSON (sensor readings) |
| `scain/commands` | subscribe | 0 | `{ "action": "reboot" }` |

### Example Payload
```json
{
  "@context": "https://ref.gs1.org/standards/epcis/2.0.0/epcis-context.jsonld",
  "type": "EPCISDocument",
  "epcisBody": {
    "eventList": [{
      "eventType": "ObjectEvent",
      "epcList": ["urn:epc:id:sgtin:0614141.Scain-001"],
      ...
    }]
  }
}
```

---

## 2. Fabric REST Gateway

Base URL: `http://localhost:4000/api/v1`

### Endpoints
| Method | Path | Function |
|--------|------|----------|
| POST | `/channels/scain-channel/chaincodes/scain-cc` | Invoke chaincode function |
| GET | `/transactions/{txid}` | Query transaction details |

#### `RecordEvent` Example
```bash
curl -X POST \
  $BASE/channels/scain-channel/chaincodes/scain-cc \
  -H 'Content-Type: application/json' \
  -d '{
        "fcn": "RecordEvent",
        "args": ["<EPCIS_JSON_STRING>"]
      }'
```

---

## 3. OpenEPCIS Repository

Capture endpoint:
```bash
POST http://localhost:8080/api/capture
Headers:
  Content-Type: application/json
  GS1-EPCIS-Version: 2.0.0
```

GraphQL playground: `http://localhost:8080/graphql`

Sample query:
```graphql
query {
  events(first: 10) {
    edges {
      node {
        eventTime
        epcList
        sensorElementList {
          sensorReport {
            type
            value
          }
        }
      }
    }
  }
}
```

---

## 4. Dashboard Hooks

| Hook | Description |
|------|-------------|
| `useTraceabilityData()` | SWR hook to fetch recent Fabric events. |
| `useMQTTConnection()` | WebSocket hook for live MQTT updates. |

---

> **Note:** For a full OpenAPI (Swagger) spec of the Fabric REST proxy, see `api/fabric-openapi.yaml` (TODO). 