#ifndef CONFIG_H
#define CONFIG_H

// WiFi Configuration
#define WIFI_SSID "YOUR_WIFI_SSID"
#define WIFI_PASSWORD "YOUR_WIFI_PASSWORD"
#define WIFI_TIMEOUT_MS 20000

// MQTT Configuration
#define MQTT_SERVER "localhost"
#define MQTT_PORT 1883
#define MQTT_TOPIC "scain/events"
#define MQTT_CLIENT_ID_PREFIX "Scain-"

// HTTP API Configuration
#define API_SERVER "http://localhost:8081"
#define API_ENDPOINT "/api/ingest"
#define USE_HTTPS false

// Device Configuration
#define DEVICE_EPC "KG-ESP32-001"
#define DEVICE_MODEL "Scain ESP32 Node v1.0"
#define FIRMWARE_VERSION "1.0.0"

// Sensor Configuration
#define DHT_PIN 4
#define DHT_TYPE DHT11
#define DS18B20_PIN 2
#define RFID_CS_PIN 5
#define LED_PIN 2

// Timing Configuration
#define SENSOR_READ_INTERVAL_MS 5000      // 5 seconds in active mode
#define ACTIVE_MODE_DURATION_MS 30000     // 30 seconds active
#define SLEEP_DURATION_SECONDS 300        // 5 minutes sleep

// NTP Configuration
#define NTP_SERVER "pool.ntp.org"
#define GMT_OFFSET_SEC 0
#define DAYLIGHT_OFFSET_SEC 0

// Debug Configuration
#define DEBUG_SERIAL true
#define DEBUG_BAUD_RATE 115200

#endif // CONFIG_H 