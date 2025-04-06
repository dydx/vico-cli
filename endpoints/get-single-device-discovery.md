# getSingleDevice Endpoint

## Overview
The getSingleDevice endpoint provides comprehensive information about a specific device. It returns detailed device properties, capabilities, status, and configuration information for a single device identified by its serial number.

## API Details
- **Path**: `/device/selectsingledevice`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Retrieves complete information for a specific device.

## Request Parameters
The request body should contain a JSON object with the following properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | The serial number of the device to query |

Plus standard BaseEntry properties:
- app
- countryNo
- language
- tenantId

## Request Example
```json
{
  "serialNumber": "ABC123456789",
  "app": "vicohome",
  "countryNo": "US",
  "language": "en",
  "tenantId": "AAA12345678"
}
```

## Response Structure
The response contains comprehensive device information:

| Property | Type | Description |
|----------|------|-------------|
| result | Integer | Status code (0 for success, negative values for errors) |
| msg | String | Status message or error description |
| data | Object | DeviceBean object containing all device information |

### DeviceBean Structure
| Property | Type | Description |
|----------|------|-------------|
| serialNumber | String | Device serial number identifier |
| deviceName | String | User-assigned device name |
| deviceModel | Object | Device model information object |
| batteryLevel | Integer | Current battery level (0-100) |
| online | Integer | Online status (1=online, 0=offline) |
| awake | Integer | Awake status (1=awake, 0=asleep) |
| isCharging | Integer | Charging status (1=charging, 0=not charging) |
| locationId | Integer | ID of the location where device is installed |
| locationName | String | Name of the location |
| thumbImgTime | Long | Timestamp of the latest thumbnail image |
| thumbImgUrl | String | URL to the latest thumbnail image |
| sdCard | Object | SD card information |
| deviceSupport | Object | Device capabilities and supported features |

### DeviceModel Structure
| Property | Type | Description |
|----------|------|-------------|
| modelName | String | Model name identifier |
| canStandby | Boolean | Whether device supports standby mode |
| audioCodecType | String | Audio codec information |
| isCanRotate | Boolean | Whether device can rotate/pan |
| supportMotionTrack | Boolean | Whether device supports motion tracking |
| deviceType | Integer | Device type identifier |
| modelType | Integer | Model type identifier |
| productCategory | Integer | Product category identifier |

### SDCard Structure
| Property | Type | Description |
|----------|------|-------------|
| formatStatus | Integer | Format status code |
| total | Long | Total storage capacity in MB |
| used | Long | Used storage in MB |

### DeviceSupport Structure
| Property | Type | Description |
|----------|------|-------------|
| deviceSupportAlarm | Boolean | Whether device supports alarms |
| supportWebrtc | Integer | WebRTC support level |
| deviceDormancySupport | Integer | Sleep/dormancy support level |
| supportFloodlight | Boolean | Whether device has floodlight |
| supportFloodlightLuminance | Boolean | Whether floodlight brightness is adjustable |

## Response Example
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "serialNumber": "ABC123456789",
    "deviceName": "Front Door Camera",
    "deviceModel": {
      "modelName": "G1",
      "canStandby": true,
      "audioCodecType": "aac",
      "isCanRotate": false,
      "supportMotionTrack": true,
      "deviceType": 1,
      "modelType": 3,
      "productCategory": 2
    },
    "batteryLevel": 85,
    "online": 1,
    "awake": 1,
    "isCharging": 0,
    "locationId": 12345,
    "locationName": "Home",
    "thumbImgTime": 1617293000,
    "thumbImgUrl": "https://example.com/thumbnail.jpg",
    "sdCard": {
      "formatStatus": 1,
      "total": 32000,
      "used": 12500
    },
    "deviceSupport": {
      "deviceSupportAlarm": true,
      "supportWebrtc": 1,
      "deviceDormancySupport": 1,
      "supportFloodlight": true,
      "supportFloodlightLuminance": true
    }
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

## Usage Context
This endpoint is typically used in the following scenarios:
- Opening a device detail page in the app
- Retrieving full device information after selecting from a list
- Checking device status and capabilities
- Getting device information prior to configuration changes

## Related Endpoints
- `listDevice` - Lists all devices with basic information
- `getDeviceAttributes` - Gets detailed device attributes
- `modifyDeviceAttributes` - Updates device attributes

## Implementation Notes
The endpoint provides a comprehensive view of a single device with all its properties, capabilities, and status information. The response contains a large number of properties that describe various aspects of the device. The actual properties returned may vary based on device type and capabilities, with certain fields only populated for specific device models.