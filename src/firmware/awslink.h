/*
 * Scain AWS ExpressLink Header
 * Class definition for ESP32-C3 ExpressLink AT command interface
 */

#ifndef AWSLINK_H
#define AWSLINK_H

#include <Arduino.h>
#include <HardwareSerial.h>

class AWSLink {
private:
  HardwareSerial serial;
  bool initialized;
  bool connected;
  bool debugEnabled = false;
  
  bool sendATCommand(const String& command, const String& expectedResponse, unsigned long timeout);
  bool sendATCommand(const String& command, String& response, unsigned long timeout);
  void debugPrint(const String& message);

public:
  AWSLink();
  
  // Initialization and connection
  bool init();
  bool connect();
  bool disconnect();
  bool isConnected();
  String getStatus();
  
  // MQTT operations
  bool publish(const String& topic, const String& payload);
  bool subscribe(const String& topic);
  String receive();
  
  // Debug control
  void enableDebug(bool enable);
};

#endif // AWSLINK_H 