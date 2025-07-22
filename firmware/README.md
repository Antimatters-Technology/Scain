# Scain ESP32 Firmware

This directory contains firmware for ESP32-based sensor nodes that send EPCIS-compliant events to the Scain backend.

## Features
- Reads DHT11, DS18B20, and RFID sensors
- Publishes EPCIS 2.0 events via MQTT or AWS IoT ExpressLink
- Deep sleep for power saving

## Build & Flash Instructions
1. Open the project in [PlatformIO](https://platformio.org/) or Arduino IDE.
2. Configure your WiFi, MQTT, and API endpoint in `main.cpp` or via a config header.
3. Connect your ESP32 to your computer.
4. Build and upload the firmware.

## Configuration
- **WiFi:** Set your SSID and password in the config section.
- **MQTT:** Set the broker address and topic.
- **API Endpoint:** If using HTTP, set the backend URL.

## Usage
- The device will read sensors every 5 minutes and publish data.
- Data is formatted as EPCIS 2.0 events for ingestion by the Scain backend.

## Troubleshooting
- Check serial output for errors.
- Ensure your backend is reachable from the ESP32.

## License
MIT 