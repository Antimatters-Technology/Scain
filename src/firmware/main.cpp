/*
 * Scain Food Traceability MVP - ESP32 Firmware
 * Reads DHT11, DS18B20, and RFID sensors every 5 minutes
 * Publishes EPCIS-compliant JSON over MQTT with deep sleep
 */

#include <WiFi.h>
#include <PubSubClient.h>
#include <DHT.h>
#include <OneWire.h>
#include <DallasTemperature.h>
#include <ArduinoJson.h>
#include <esp_sleep.h>
#include <esp_wifi.h>
#include "time.h"
#include "awslink.h"

// Hardware pin definitions
#define DHT_PIN 4
#define DHT_TYPE DHT11
#define DS18B20_PIN 2
#define RFID_CS_PIN 5
#define LED_PIN 2
#define SLEEP_DURATION 300 // 5 minutes in seconds

// Network credentials
const char* ssid = "YOUR_WIFI_SSID";
const char* password = "YOUR_WIFI_PASSWORD";
const char* mqtt_server = "localhost";
const int mqtt_port = 1883;
const char* mqtt_topic = "scain/events";

// Sensor instances
DHT dht(DHT_PIN, DHT_TYPE);
OneWire oneWire(DS18B20_PIN);
DallasTemperature ds18b20(&oneWire);

// MQTT client
WiFiClient espClient;
PubSubClient client(espClient);

// AWS ExpressLink instance
AWSLink awsLink;

// Global variables
String deviceEPC = "KG-ESP32-001"; // Default EPC, will be read from RFID
unsigned long lastSensorRead = 0;
bool useAWSLink = false;

void setup() {
  Serial.begin(115200);
  
  // Initialize sensors
  dht.begin();
  ds18b20.begin();
  
  // Initialize LED
  pinMode(LED_PIN, OUTPUT);
  digitalWrite(LED_PIN, LOW);
  
  // Check if we should use AWS ExpressLink
  if (awsLink.init()) {
    useAWSLink = true;
    Serial.println("AWS ExpressLink initialized");
  } else {
    Serial.println("Using local MQTT");
  }
  
  // Setup WiFi and MQTT
  setupWiFi();
  if (!useAWSLink) {
    client.setServer(mqtt_server, mqtt_port);
  }
  
  // Configure NTP for timestamps
  configTime(0, 0, "pool.ntp.org");
  
  Serial.println("Scain sensor node started");
  blinkLED(3);
}

void loop() {
  if (millis() - lastSensorRead > 5000) { // Read every 5 seconds in active mode
    readAndPublishSensors();
    lastSensorRead = millis();
  }
  
  // Enter deep sleep after 30 seconds of operation
  if (millis() > 30000) {
    enterDeepSleep();
  }
  
  delay(1000);
}

void setupWiFi() {
  delay(10);
  Serial.println();
  Serial.print("Connecting to ");
  Serial.println(ssid);
  
  WiFi.begin(ssid, password);
  
  int attempts = 0;
  while (WiFi.status() != WL_CONNECTED && attempts < 20) {
    delay(500);
    Serial.print(".");
    attempts++;
  }
  
  if (WiFi.status() == WL_CONNECTED) {
    Serial.println("");
    Serial.println("WiFi connected");
    Serial.println("IP address: ");
    Serial.println(WiFi.localIP());
  } else {
    Serial.println("WiFi connection failed, using fallback mode");
  }
}

void readAndPublishSensors() {
  // Read DHT11 sensor
  float humidity = dht.readHumidity();
  float airTemp = dht.readTemperature();
  
  // Read DS18B20 probe
  ds18b20.requestTemperatures();
  float probeTemp = ds18b20.getTempCByIndex(0);
  
  // Check for valid readings
  if (isnan(humidity) || isnan(airTemp)) {
    Serial.println("Failed to read from DHT sensor!");
    return;
  }
  
  if (probeTemp == DEVICE_DISCONNECTED_C) {
    Serial.println("DS18B20 sensor disconnected!");
    probeTemp = -999.0; // Error value
  }
  
  // Get current timestamp
  time_t now;
  struct tm timeinfo;
  time(&now);
  localtime_r(&now, &timeinfo);
  
  char timestamp[32];
  strftime(timestamp, sizeof(timestamp), "%Y-%m-%dT%H:%M:%SZ", &timeinfo);
  
  // Create EPCIS 2.0 compliant JSON payload
  DynamicJsonDocument doc(1024);
  
  // EPCIS Event structure
  doc["@context"] = "https://ref.gs1.org/standards/epcis/2.0.0/epcis-context.jsonld";
  doc["type"] = "EPCISDocument";
  doc["schemaVersion"] = "2.0";
  doc["creationDate"] = timestamp;
  
  JsonObject event = doc.createNestedObject("epcisBody").createNestedArray("eventList").createNestedObject();
  event["eventType"] = "ObjectEvent";
  event["eventTime"] = timestamp;
  event["eventTimeZoneOffset"] = "+00:00";
  event["recordTime"] = timestamp;
  
  // EPC and business context
  JsonArray epcList = event.createNestedArray("epcList");
  epcList.add("urn:epc:id:sgtin:0614141." + deviceEPC);
  
  event["action"] = "OBSERVE";
  event["bizStep"] = "urn:epcglobal:cbv:bizstep:sensor_reporting";
  event["disposition"] = "urn:epcglobal:cbv:disp:in_transit";
  
  // Sensor data extensions
  JsonObject sensorData = event.createNestedObject("sensorElementList").createNestedObject();
  sensorData["sensorMetadata"]["time"] = timestamp;
  sensorData["sensorMetadata"]["deviceID"] = deviceEPC;
      sensorData["sensorMetadata"]["deviceMetadata"] = "Scain ESP32 Node v1.0";
  
  JsonArray sensorReports = sensorData.createNestedArray("sensorReport");
  
  // Air temperature report
  JsonObject airTempReport = sensorReports.createNestedObject();
  airTempReport["type"] = "gs1:Temperature";
  airTempReport["value"] = airTemp;
  airTempReport["uom"] = "CEL";
  airTempReport["component"] = "air";
  
  // Probe temperature report
  JsonObject probeTempReport = sensorReports.createNestedObject();
  probeTempReport["type"] = "gs1:Temperature";
  probeTempReport["value"] = probeTemp;
  probeTempReport["uom"] = "CEL";
  probeTempReport["component"] = "probe";
  
  // Humidity report
  JsonObject humidityReport = sensorReports.createNestedObject();
  humidityReport["type"] = "gs1:RelativeHumidity";
  humidityReport["value"] = humidity;
  humidityReport["uom"] = "A93";
  humidityReport["component"] = "air";
  
  // Serialize JSON
  String jsonString;
  serializeJson(doc, jsonString);
  
  // Publish via MQTT or AWS ExpressLink
  if (useAWSLink) {
    publishViaAWS(jsonString);
  } else {
    publishViaMQTT(jsonString);
  }
  
  // Debug output
  Serial.println("Sensor data published:");
  Serial.printf("Air Temp: %.1f°C, Probe Temp: %.1f°C, Humidity: %.1f%%\n", 
                airTemp, probeTemp, humidity);
  
  blinkLED(1);
}

void publishViaMQTT(const String& payload) {
  if (!client.connected()) {
    reconnectMQTT();
  }
  
  if (client.connected()) {
    client.publish(mqtt_topic, payload.c_str());
    Serial.println("Published to MQTT");
  } else {
    Serial.println("MQTT publish failed");
  }
}

void publishViaAWS(const String& payload) {
  if (awsLink.publish("knowgraph/events", payload)) {
    Serial.println("Published to AWS IoT");
  } else {
    Serial.println("AWS publish failed, falling back to MQTT");
    publishViaMQTT(payload);
  }
}

void reconnectMQTT() {
  while (!client.connected()) {
    Serial.print("Attempting MQTT connection...");
    
    String clientId = "Scain-" + String(random(0xffff), HEX);
    
    if (client.connect(clientId.c_str())) {
      Serial.println("connected");
    } else {
      Serial.print("failed, rc=");
      Serial.print(client.state());
      Serial.println(" try again in 5 seconds");
      delay(5000);
    }
  }
}

void enterDeepSleep() {
  Serial.println("Entering deep sleep for 5 minutes...");
  
  // Disconnect WiFi to save power
  WiFi.disconnect(true);
  WiFi.mode(WIFI_OFF);
  
  // Configure wake up timer
  esp_sleep_enable_timer_wakeup(SLEEP_DURATION * 1000000ULL);
  
  // Enter deep sleep
  esp_deep_sleep_start();
}

void blinkLED(int times) {
  for (int i = 0; i < times; i++) {
    digitalWrite(LED_PIN, HIGH);
    delay(200);
    digitalWrite(LED_PIN, LOW);
    delay(200);
  }
} 