FROM python:3.11-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements
COPY requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY scripts/mqtt_bridge.py .

# Set environment variables
ENV PYTHONUNBUFFERED=1
ENV MQTT_BROKER=mosquitto:1883
ENV FABRIC_ENDPOINT=fabric-rest-api:4000
ENV OPENEPCIS_ENDPOINT=http://openepcis:8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD python -c "import paho.mqtt.client as mqtt; client = mqtt.Client(); client.connect('$MQTT_BROKER'.split(':')[0], int('$MQTT_BROKER'.split(':')[1]), 60)" || exit 1

# Run the bridge
CMD ["python", "mqtt_bridge.py"] 