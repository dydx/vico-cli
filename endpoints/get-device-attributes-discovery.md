# getDeviceAttributes Endpoint

## Overview
The getDeviceAttributes endpoint provides comprehensive information about a device's capabilities, current status, and configurable settings. This endpoint retrieves static hardware specifications, real-time state information, and user-configurable attributes for a specific device.

## API Details
- **Path**: `/device/getDeviceAttributes`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Retrieves detailed attribute information for a specific device identified by its serial number.

## Request Parameters
The request body should contain a JSON object with the following properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | The serial number of the device to query |
| returnFixedAttributes | Boolean | No | Flag to include fixed/hardware attributes in the response (default: true) |
| returnRealTimeAttributes | Boolean | No | Flag to include real-time status attributes in the response (default: true) |
| returnModifiableAttributes | Boolean | No | Flag to include user-configurable attributes in the response (default: true) |

Plus standard BaseEntry properties:
- app
- countryNo
- language
- tenantId

## Request Example
```json
{
  "serialNumber": "ABC123456789",
  "returnFixedAttributes": true,
  "returnRealTimeAttributes": true,
  "returnModifiableAttributes": true,
  "app": "vicohome",
  "countryNo": "US",
  "language": "en",
  "tenantId": "AAA12345678"
}
```

## Response Structure
The response contains comprehensive device information divided into three main categories:

| Property | Type | Description |
|----------|------|-------------|
| result | Integer | Status code (0 for success, negative values for errors) |
| msg | String | Status message or error description |
| data | Object | Contains device attribute information |

### Data Object Structure
| Property | Type | Description |
|----------|------|-------------|
| serialNumber | String | Device serial number identifier |
| fixedAttributes | Object | Hardware and non-configurable attributes |
| realTimeAttributes | Object | Current device status information |
| modifiableAttributes | Array | User-configurable settings |

### Fixed Attributes Include
| Property | Type | Description |
|----------|------|-------------|
| modelNo | String | Device model number |
| displayModelNo | String | User-friendly model name |
| macAddress | String | Device MAC address |
| wiredMacAddress | String | Ethernet MAC address (if applicable) |
| icon | String | URL to device icon image |
| smallIcon | String | URL to small device icon image |
| supportWhiteLight | Boolean | Whether device supports white light |
| supportManualFloodlightSwitch | Boolean | Whether device has manual floodlight control |
| supportFloodlightLuminance | Boolean | Whether device supports floodlight brightness adjustment |
| floodlightLuminanceRange | Object | Min/max brightness values for floodlight |

### Real-Time Attributes Include
| Property | Type | Description |
|----------|------|-------------|
| batteryLevel | Integer | Current battery percentage (0-100) |
| online | Boolean | Whether device is currently connected |
| chargingMode | Integer | Charging mode identifier |
| charging | Boolean | Whether device is currently charging |
| firmwareId | String | Current firmware version |
| newestFirmwareId | String | Latest available firmware version |
| sdCard | Object | Storage information (capacity, used space) |
| signalStrength | Integer | WiFi signal strength |
| wifiChannel | Integer | Current WiFi channel |

### Modifiable Attributes Structure
Each modifiable attribute is an object with:

| Property | Type | Description |
|----------|------|-------------|
| name | String | Setting name |
| type | String | Data type (boolean, integer, string, etc.) |
| value | Varies | Current value of the setting |
| options | Array | Available options (for selection-type settings) |
| intRange | Object | Min/max/interval values (for numeric settings) |

## Response Example
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "serialNumber": "ABC123456789",
    "fixedAttributes": {
      "modelNo": "A8C1",
      "displayModelNo": "Doorbell Pro",
      "macAddress": "AA:BB:CC:DD:EE:FF",
      "icon": "https://resources.vicohome.io/icons/doorbell_pro.png",
      "supportWhiteLight": true,
      "supportManualFloodlightSwitch": true,
      "supportFloodlightLuminance": true,
      "floodlightLuminanceRange": {
        "min": 1,
        "max": 100
      }
    },
    "realTimeAttributes": {
      "batteryLevel": 85,
      "online": true,
      "chargingMode": 0,
      "charging": false,
      "firmwareId": "1.2.345",
      "newestFirmwareId": "1.2.345",
      "sdCard": {
        "capacity": 32000000000,
        "used": 12500000000
      },
      "signalStrength": 75,
      "wifiChannel": 6
    },
    "modifiableAttributes": [
      {
        "name": "motionDetection",
        "type": "boolean",
        "value": true
      },
      {
        "name": "nightVision",
        "type": "integer",
        "value": 1,
        "options": ["Auto", "On", "Off"],
        "intRange": {
          "min": 0,
          "max": 2,
          "interval": 1
        }
      }
    ]
  }
}
```

## Error Codes
| Error Code | Description |
|------------|-------------|
| -1 | General error |
| -2 | Invalid parameters |
| -100 | Device not found |
| -101 | Device offline |
| -102 | Access denied |

## Usage Context
This endpoint is typically used in the following scenarios:
- Initial device setup and configuration
- Device detail screens in the application
- Determining available device capabilities
- Checking device status (battery level, online status, etc.)
- Retrieving current settings before modifying them

## Related Endpoints
- `modifyDeviceAttributes` - Updates user-configurable device attributes
- `getSingleDevice` - Gets basic device information
- `listDevice` - Lists all devices with basic information

## Implementation Notes
The endpoint is designed to be flexible, allowing clients to request only the specific attribute categories they need (fixed, real-time, or modifiable) by setting the corresponding flags. This can reduce response size and processing time when only certain information is required. The implementation allows for a comprehensive view of device capabilities, current status, and configurable settings in a single request.