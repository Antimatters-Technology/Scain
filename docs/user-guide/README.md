# User Guide

Welcome to the Scain user guide. This document provides instructions for using the food traceability platform.

## 🚀 Getting Started

### Accessing the Application

1. **Open your web browser**
2. **Navigate to**: `http://localhost:3000` (development) or `https://scain.com` (production)
3. **Login** with your credentials (if authentication is enabled)

### First Time Setup

1. **Complete your profile**
   - Update your contact information
   - Set your organization details
   - Configure notification preferences

2. **Explore the dashboard**
   - Familiarize yourself with the layout
   - Review the navigation menu
   - Check the status indicators

## 📊 Dashboard Overview

### Main Dashboard

The dashboard provides an overview of your food traceability system:

```
┌─────────────────────────────────────────────────────────────┐
│ Header: Logo, Navigation, User Menu                       │
├─────────────────────────────────────────────────────────────┤
│ Sidebar: Navigation Menu                                   │
│ ├── Dashboard                                             │
│ ├── Devices                                              │
│ ├── Events                                               │
│ ├── Traceability                                         │
│ ├── Reports                                              │
│ └── Settings                                             │
├─────────────────────────────────────────────────────────────┤
│ Main Content Area                                         │
│ ├── Quick Stats Cards                                     │
│ ├── Recent Activity                                       │
│ ├── System Status                                         │
│ └── Alerts & Notifications                               │
└─────────────────────────────────────────────────────────────┘
```

### Key Metrics

The dashboard displays important metrics:

- **Total Devices**: Number of connected IoT devices
- **Active Events**: Recent traceability events
- **System Health**: Overall system status
- **Alerts**: Critical notifications requiring attention

## 🔧 Device Management

### Adding a New Device

1. **Navigate to Devices**
   - Click "Devices" in the sidebar
   - Click "Add Device" button

2. **Enter Device Information**
   ```
   Device Name: ESP32-Sensor-001
   Device Type: Temperature Sensor
   Location: Cold Storage Room A
   EPC Code: urn:epc:id:sgtin:0614141.Scain-001
   ```

3. **Configure Settings**
   - Set temperature thresholds
   - Configure alert notifications
   - Set reporting frequency

4. **Save Device**
   - Click "Save" to register the device
   - The device will appear in your device list

### Device Status

Devices can have the following statuses:

- **🟢 Online**: Device is connected and reporting
- **🟡 Warning**: Device has issues but still functioning
- **🔴 Offline**: Device is not responding
- **⚪ Unknown**: Device status is unclear

### Device Actions

For each device, you can:

- **View Details**: See device information and history
- **Edit Settings**: Modify device configuration
- **View Data**: Access sensor readings and events
- **Generate Report**: Create device-specific reports
- **Remove Device**: Delete device from system

## 📈 Event Monitoring

### Viewing Events

1. **Navigate to Events**
   - Click "Events" in the sidebar
   - View recent traceability events

2. **Filter Events**
   - Use date range picker
   - Filter by device type
   - Search by EPC code
   - Filter by event type

### Event Types

The system tracks various event types:

- **ObjectEvent**: Product observation events
- **AggregationEvent**: Product grouping events
- **TransactionEvent**: Business transaction events
- **TransformationEvent**: Product transformation events

### Event Details

Each event contains:

- **Event ID**: Unique identifier
- **Event Time**: When the event occurred
- **Device ID**: Which device generated the event
- **EPC List**: Product identifiers involved
- **Sensor Data**: Environmental readings
- **Location**: Where the event occurred

## 🔍 Traceability Features

### Product Tracing

1. **Enter Product Code**
   - Use the search bar to enter an EPC code
   - Or scan a QR code/barcode

2. **View Trace History**
   - See the complete product journey
   - View all events associated with the product
   - Check environmental conditions

3. **Export Trace Data**
   - Download trace data as CSV
   - Generate compliance reports
   - Share with stakeholders

### Lot Tracking

1. **Select Lot**
   - Choose from your lot list
   - Or create a new lot

2. **View Lot Details**
   - See all products in the lot
   - View lot creation and processing events
   - Check lot status and location

3. **Lot Operations**
   - Split lots into smaller groups
   - Merge multiple lots
   - Transfer lot ownership

## 📋 Compliance Reporting

### Generate Reports

1. **Navigate to Reports**
   - Click "Reports" in the sidebar
   - Select report type

2. **Configure Report**
   - Set date range
   - Choose devices/products
   - Select report format

3. **Generate and Download**
   - Click "Generate Report"
   - Download when ready
   - Share via email or cloud storage

### Report Types

Available report types:

- **Daily Summary**: Daily activity overview
- **Device Performance**: Device health and data quality
- **Compliance Report**: Regulatory compliance data
- **Traceability Report**: Product trace data
- **Alert Summary**: System alerts and notifications

## ⚙️ System Settings

### User Preferences

1. **Access Settings**
   - Click your user icon
   - Select "Settings"

2. **Configure Preferences**
   - **Notifications**: Email and SMS alerts
   - **Language**: Interface language
   - **Timezone**: Your local timezone
   - **Theme**: Light or dark mode

### Organization Settings

1. **Company Information**
   - Update company details
   - Set contact information
   - Configure business hours

2. **System Configuration**
   - Set default temperature units
   - Configure alert thresholds
   - Set data retention policies

## 🔔 Alerts and Notifications

### Alert Types

The system generates various alerts:

- **Temperature Alerts**: When temperature exceeds thresholds
- **Device Offline**: When devices stop reporting
- **Data Quality**: When sensor data is questionable
- **Compliance**: When compliance requirements are not met

### Managing Alerts

1. **View Alerts**
   - Check the alerts panel on the dashboard
   - Navigate to "Alerts" in the sidebar

2. **Respond to Alerts**
   - Acknowledge alerts
   - Take corrective action
   - Document responses

3. **Configure Alert Rules**
   - Set temperature thresholds
   - Configure notification preferences
   - Define escalation procedures

## 📱 Mobile Access

### Mobile Dashboard

The dashboard is responsive and works on mobile devices:

- **Touch-friendly interface**
- **Optimized for small screens**
- **Quick access to key features**

### Mobile Features

- **Real-time notifications**
- **Quick device status check**
- **Emergency alert response**
- **Photo capture for events**

## 🔐 Security Features

### User Authentication

- **Secure login**: Username and password
- **Two-factor authentication**: Optional 2FA
- **Session management**: Automatic logout
- **Password policies**: Strong password requirements

### Data Protection

- **Encrypted data transmission**: HTTPS/TLS
- **Secure data storage**: Encrypted databases
- **Access controls**: Role-based permissions
- **Audit trails**: Complete activity logging

## 🆘 Troubleshooting

### Common Issues

1. **Can't Access Dashboard**
   - Check internet connection
   - Verify URL is correct
   - Clear browser cache
   - Try different browser

2. **Device Not Reporting**
   - Check device power
   - Verify network connection
   - Check device configuration
   - Contact support if needed

3. **Data Not Loading**
   - Refresh the page
   - Check date range settings
   - Verify permissions
   - Contact administrator

### Getting Help

1. **Documentation**
   - Check this user guide
   - Review FAQ section
   - Search knowledge base

2. **Support Contact**
   - Email: support@scain.com
   - Phone: +1-555-SCAIN-1
   - Chat: Available in the app

3. **Emergency Support**
   - For critical issues affecting food safety
   - 24/7 emergency hotline
   - Immediate response guaranteed

## 📚 Keyboard Shortcuts

### Navigation Shortcuts

- **Ctrl/Cmd + D**: Go to Dashboard
- **Ctrl/Cmd + E**: Go to Events
- **Ctrl/Cmd + V**: Go to Devices
- **Ctrl/Cmd + T**: Go to Traceability
- **Ctrl/Cmd + R**: Go to Reports
- **Ctrl/Cmd + S**: Go to Settings

### Action Shortcuts

- **Ctrl/Cmd + N**: New device/event
- **Ctrl/Cmd + F**: Search
- **Ctrl/Cmd + P**: Print current view
- **Ctrl/Cmd + E**: Export data
- **F5**: Refresh page
- **Esc**: Close modal/dialog

## 📖 Glossary

### Technical Terms

- **EPC**: Electronic Product Code - unique product identifier
- **EPCIS**: Electronic Product Code Information Services
- **IoT**: Internet of Things - connected devices
- **Lot**: Group of products processed together
- **Traceability**: Ability to track product through supply chain

### Business Terms

- **FSMA**: Food Safety Modernization Act
- **SFCR**: Safe Food for Canadians Regulations
- **HACCP**: Hazard Analysis and Critical Control Points
- **GMP**: Good Manufacturing Practices

---

**Last Updated**: July 2025  
**User Guide Version**: 1.0.0 