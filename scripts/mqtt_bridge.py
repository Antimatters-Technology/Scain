#!/usr/bin/env python3
"""
KnowGraph MQTT to Fabric Bridge
Forwards EPCIS sensor data from MQTT to Hyperledger Fabric chaincode
"""

import json
import logging
import os
import time
from typing import Dict, Any

import paho.mqtt.client as mqtt
import requests
from retry import retry
from pythonjsonlogger import jsonlogger

# Configure logging
logHandler = logging.StreamHandler()
formatter = jsonlogger.JsonFormatter()
logHandler.setFormatter(formatter)
logger = logging.getLogger()
logger.addHandler(logHandler)
logger.setLevel(logging.INFO)

# Configuration
MQTT_BROKER = os.getenv('MQTT_BROKER', 'mosquitto:1883')
FABRIC_ENDPOINT = os.getenv('FABRIC_ENDPOINT', 'fabric-rest-api:4000')
OPENEPCIS_ENDPOINT = os.getenv('OPENEPCIS_ENDPOINT', 'http://openepcis:8080')
MQTT_TOPIC = 'knowgraph/events'
FABRIC_CHANNEL = 'knowgraph-channel'
FABRIC_CHAINCODE = 'knowgraph-cc'

class MQTTBridge:
    def __init__(self):
        self.mqtt_client = None
        self.connected = False
        self.setup_mqtt()
        
    def setup_mqtt(self):
        """Initialize MQTT client"""
        self.mqtt_client = mqtt.Client()
        self.mqtt_client.on_connect = self.on_connect
        self.mqtt_client.on_message = self.on_message
        self.mqtt_client.on_disconnect = self.on_disconnect
        
        # Parse broker address
        if ':' in MQTT_BROKER:
            host, port = MQTT_BROKER.split(':')
            port = int(port)
        else:
            host = MQTT_BROKER
            port = 1883
            
        logger.info(f"Connecting to MQTT broker: {host}:{port}")
        self.mqtt_client.connect(host, port, 60)
        
    def on_connect(self, client, userdata, flags, rc):
        """Callback for MQTT connection"""
        if rc == 0:
            self.connected = True
            logger.info("Connected to MQTT broker")
            client.subscribe(MQTT_TOPIC)
            logger.info(f"Subscribed to topic: {MQTT_TOPIC}")
        else:
            logger.error(f"Failed to connect to MQTT broker: {rc}")
            
    def on_disconnect(self, client, userdata, rc):
        """Callback for MQTT disconnection"""
        self.connected = False
        logger.warning("Disconnected from MQTT broker")
        
    def on_message(self, client, userdata, msg):
        """Handle incoming MQTT messages"""
        try:
            payload = msg.payload.decode('utf-8')
            topic = msg.topic
            
            logger.info(f"Received message on topic {topic}")
            
            # Parse EPCIS JSON
            epcis_data = json.loads(payload)
            
            # Validate EPCIS structure
            if not self.validate_epcis(epcis_data):
                logger.error("Invalid EPCIS data received")
                return
                
            # Forward to Fabric
            self.forward_to_fabric(payload)
            
            # Forward to OpenEPCIS
            self.forward_to_openepcis(epcis_data)
            
        except Exception as e:
            logger.error(f"Error processing message: {e}")
            
    def validate_epcis(self, data: Dict[str, Any]) -> bool:
        """Validate EPCIS 2.0 structure"""
        required_fields = ['@context', 'type', 'epcisBody']
        
        for field in required_fields:
            if field not in data:
                logger.error(f"Missing required field: {field}")
                return False
                
        if data.get('type') != 'EPCISDocument':
            logger.error(f"Invalid document type: {data.get('type')}")
            return False
            
        if not data.get('epcisBody', {}).get('eventList'):
            logger.error("No events in eventList")
            return False
            
        return True
        
    @retry(tries=3, delay=1, backoff=2)
    def forward_to_fabric(self, epcis_json: str):
        """Forward EPCIS data to Hyperledger Fabric"""
        try:
            url = f"http://{FABRIC_ENDPOINT}/api/v1/channels/{FABRIC_CHANNEL}/chaincodes/{FABRIC_CHAINCODE}"
            
            payload = {
                "fcn": "RecordEvent",
                "args": [epcis_json]
            }
            
            headers = {
                'Content-Type': 'application/json',
                'X-Fabric-User': 'admin'
            }
            
            response = requests.post(url, json=payload, headers=headers, timeout=30)
            
            if response.status_code == 200:
                result = response.json()
                logger.info(f"Successfully recorded event in Fabric: {result}")
            else:
                logger.error(f"Fabric API error: {response.status_code} - {response.text}")
                raise Exception(f"Fabric API error: {response.status_code}")
                
        except Exception as e:
            logger.error(f"Failed to forward to Fabric: {e}")
            raise
            
    @retry(tries=3, delay=1, backoff=2)
    def forward_to_openepcis(self, epcis_data: Dict[str, Any]):
        """Forward EPCIS data to OpenEPCIS repository"""
        try:
            url = f"{OPENEPCIS_ENDPOINT}/api/capture"
            
            headers = {
                'Content-Type': 'application/json',
                'GS1-EPCIS-Version': '2.0.0'
            }
            
            response = requests.post(url, json=epcis_data, headers=headers, timeout=30)
            
            if response.status_code in [200, 201, 202]:
                logger.info("Successfully stored event in OpenEPCIS")
            else:
                logger.error(f"OpenEPCIS API error: {response.status_code} - {response.text}")
                raise Exception(f"OpenEPCIS API error: {response.status_code}")
                
        except Exception as e:
            logger.error(f"Failed to forward to OpenEPCIS: {e}")
            raise
            
    def health_check(self):
        """Perform health check"""
        try:
            # Check MQTT connection
            if not self.connected:
                logger.error("MQTT not connected")
                return False
                
            # Check Fabric API
            fabric_url = f"http://{FABRIC_ENDPOINT}/api/v1/health"
            fabric_response = requests.get(fabric_url, timeout=10)
            if fabric_response.status_code != 200:
                logger.error(f"Fabric API unhealthy: {fabric_response.status_code}")
                return False
                
            # Check OpenEPCIS
            openepcis_url = f"{OPENEPCIS_ENDPOINT}/api/health"
            openepcis_response = requests.get(openepcis_url, timeout=10)
            if openepcis_response.status_code != 200:
                logger.error(f"OpenEPCIS unhealthy: {openepcis_response.status_code}")
                return False
                
            logger.info("Health check passed")
            return True
            
        except Exception as e:
            logger.error(f"Health check failed: {e}")
            return False
            
    def run(self):
        """Start the bridge service"""
        logger.info("Starting KnowGraph MQTT Bridge")
        
        # Start MQTT client
        self.mqtt_client.loop_start()
        
        # Health check loop
        last_health_check = 0
        health_check_interval = 60  # seconds
        
        try:
            while True:
                current_time = time.time()
                
                # Periodic health check
                if current_time - last_health_check > health_check_interval:
                    self.health_check()
                    last_health_check = current_time
                    
                time.sleep(1)
                
        except KeyboardInterrupt:
            logger.info("Shutting down bridge")
            self.mqtt_client.loop_stop()
            self.mqtt_client.disconnect()

def main():
    """Main entry point"""
    bridge = MQTTBridge()
    bridge.run()

if __name__ == "__main__":
    main() 