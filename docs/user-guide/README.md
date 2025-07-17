# 🧑‍🏭 User Guide

Welcome to **Scain**! This guide walks operators through daily tasks.

## 1. Hardware Setup

1. Affix the ESP32 sensor node to the pallet.
2. Attach the DS18B20 probe inside product core.
3. Stick the RFID inlay or QR label on the outer wrap.
4. Power up the node (18650 cell or USB-C PD bank).
5. Ensure the blue LED blinks **3×** (Wi-Fi connected).

> **TIP:** For cold-chain, position the DHT11 away from evaporator airflow.

## 2. Dashboard Overview

| Section | Purpose |
|---------|---------|
| Status Cards | Real-time temp/humidity and system state. |
| Alarm Panel | Highlights breaches (>8 °C probe, >85 % RH, etc.). |
| Trend Chart | 24 h temperature history (air + probe). |
| Events Table | Last 100 EPCIS events with hash link. |
| Compliance Blocks | FSMA & SFCR checklist status. |

### Access
Open `http://<gateway-ip>:3000` in Chrome/Edge.  
Default credentials: `admin / scain123` (change under **Settings → Users**).

## 3. Running a Recall Drill

1. Navigate to **Trace → Search**.
2. Scan or enter **Lot Code** (e.g., `SGT-20240715-001`).
3. Click **Retrieve Chain of Custody**.
4. CSV export downloads within 30 s (<24 h rule).

## 4. Firmware OTA Update

```bash
make flash-ota FW=firmware_v1.1.bin
```
Nodes will auto-reboot into new firmware and report version under **Devices**.

## 5. Troubleshooting

| Symptom | Likely Cause | Fix |
|---------|--------------|-----|
| Red LED solid | Wi-Fi not configured | Hold BOOT 5 s to start AP `SCAIN_SETUP`. |
| No data in dashboard | MQTT down | `docker restart scain-mosquitto`. |
| Chaincode error 500 | Peer unhealthy | `docker-compose restart peer0 couchdb`. |

## 6. Support

Email `support@scain.io` or open a GitHub issue with logs (Settings → Download Diagnostics). 