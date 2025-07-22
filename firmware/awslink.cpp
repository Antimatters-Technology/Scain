/*
 * Scain AWS ExpressLink Wrapper
 * Handles AT commands for ESP32-C3 ExpressLink module
 * Provides secure TLS connection to AWS IoT Core
 */

#include "awslink.h"
#include <HardwareSerial.h>

AWSLink::AWSLink() : serial(2), initialized(false), connected(false) {
  // Use Serial2 for AT commands (GPIO 16/17)
}

bool AWSLink::init() {
  serial.begin(115200, SERIAL_8N1, 16, 17); // RX=16, TX=17
  delay(1000);
  
  // Test basic AT command
  if (!sendATCommand("AT", "OK", 1000)) {
    Serial.println("ExpressLink module not responding");
    return false;
  }
  
  // Get module information
  String response;
  if (sendATCommand("AT+CONF? About", response, 2000)) {
    Serial.println("ExpressLink Info: " + response);
  }
  
  // Check connection status
  if (sendATCommand("AT+CONNECT?", response, 2000)) {
    if (response.indexOf("1") >= 0) {
      connected = true;
      Serial.println("ExpressLink already connected");
    }
  }
  
  initialized = true;
  return true;
}

bool AWSLink::connect() {
  if (!initialized) {
    return false;
  }
  
  if (connected) {
    return true;
  }
  
  Serial.println("Connecting to AWS IoT Core...");
  
  // Start connection
  if (!sendATCommand("AT+CONNECT", "OK", 2000)) {
    Serial.println("Failed to initiate connection");
    return false;
  }
  
  // Wait for connection establishment (up to 30 seconds)
  for (int i = 0; i < 30; i++) {
    String response;
    if (sendATCommand("AT+CONNECT?", response, 1000)) {
      if (response.indexOf("1") >= 0) {
        connected = true;
        Serial.println("Connected to AWS IoT Core");
        return true;
      }
    }
    delay(1000);
  }
  
  Serial.println("Connection timeout");
  return false;
}

bool AWSLink::publish(const String& topic, const String& payload) {
  if (!initialized || !connected) {
    if (!connect()) {
      return false;
    }
  }
  
  // Escape quotes in payload
  String escapedPayload = payload;
  escapedPayload.replace("\"", "\\\"");
  
  // Construct publish command
  String command = "AT+SEND " + topic + " " + escapedPayload;
  
  // Send publish command
  if (sendATCommand(command, "OK", 5000)) {
    Serial.println("Message published to: " + topic);
    return true;
  } else {
    Serial.println("Failed to publish message");
    connected = false; // Reset connection status
    return false;
  }
}

bool AWSLink::subscribe(const String& topic) {
  if (!initialized || !connected) {
    if (!connect()) {
      return false;
    }
  }
  
  String command = "AT+SUBSCRIBE " + topic;
  
  if (sendATCommand(command, "OK", 2000)) {
    Serial.println("Subscribed to: " + topic);
    return true;
  } else {
    Serial.println("Failed to subscribe to: " + topic);
    return false;
  }
}

String AWSLink::receive() {
  if (!initialized || !connected) {
    return "";
  }
  
  String response;
  if (sendATCommand("AT+GET", response, 1000)) {
    // Parse received message
    int startIndex = response.indexOf(" ");
    if (startIndex > 0) {
      return response.substring(startIndex + 1);
    }
  }
  
  return "";
}

bool AWSLink::disconnect() {
  if (!initialized) {
    return false;
  }
  
  if (sendATCommand("AT+DISCONNECT", "OK", 2000)) {
    connected = false;
    Serial.println("Disconnected from AWS IoT Core");
    return true;
  }
  
  return false;
}

bool AWSLink::sendATCommand(const String& command, const String& expectedResponse, unsigned long timeout) {
  String response;
  return sendATCommand(command, response, timeout) && response.indexOf(expectedResponse) >= 0;
}

bool AWSLink::sendATCommand(const String& command, String& response, unsigned long timeout) {
  // Clear serial buffer
  while (serial.available()) {
    serial.read();
  }
  
  // Send command
  serial.println(command);
  
  // Wait for response
  unsigned long startTime = millis();
  response = "";
  
  while (millis() - startTime < timeout) {
    if (serial.available()) {
      char c = serial.read();
      response += c;
      
      // Check for command completion
      if (response.endsWith("\r\n")) {
        response.trim();
        if (response.length() > 0) {
          return true;
        }
      }
    }
    delay(10);
  }
  
  return false;
}

bool AWSLink::isConnected() {
  return connected;
}

String AWSLink::getStatus() {
  if (!initialized) {
    return "Not initialized";
  }
  
  String response;
  if (sendATCommand("AT+CONNECT?", response, 1000)) {
    if (response.indexOf("1") >= 0) {
      return "Connected";
    } else {
      return "Disconnected";
    }
  }
  
  return "Status unknown";
}

void AWSLink::enableDebug(bool enable) {
  debugEnabled = enable;
}

void AWSLink::debugPrint(const String& message) {
  if (debugEnabled) {
    Serial.println("[AWSLink] " + message);
  }
} 