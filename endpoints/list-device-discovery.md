# listDevice Endpoint

## Overview
The listDevice endpoint retrieves a list of all devices associated with the user's account. This endpoint provides comprehensive information about each device, including status, capabilities, and configuration settings.

## API Details
- **Path**: `/device/listuserdevices`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Retrieves a list of all devices associated with the authenticated user.

## Request Parameters
The request body requires only standard BaseEntry properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| app | Object | Yes | Application information |
| countryNo | String | Yes | Country code (e.g., "US") |
| language | String | Yes | Language code (e.g., "en") |
| tenantId | String | Yes | Tenant identifier |

## Request Example
```json
{
  "app": {
    "appType": "Android",
    "appName": "Vicohome",
    "versionName": "1.0.0",
    "version": 100,
    "apiVersion": "1",
    "bundle": "io.vicohome.app",
    "timeZone": "America/Los_Angeles",
    "env": "production"
  },
  "countryNo": "US",
  "language": "en",
  "tenantId": "AAA12345678"
}
```

## Response Structure
The response contains an array of device information:

| Property | Type | Description |
|----------|------|-------------|
| result | Integer | Status code (0 for success, negative values for errors) |
| msg | String | Status message or error description |
| data | Object | Contains the device list and related data |

### Data Object Structure
| Property | Type | Description |
|----------|------|-------------|
| list | Array | Array of DeviceBean objects representing each device |

### DeviceBean Structure
Each device in the list has the following properties:

| Property | Type | Description |
|----------|------|-------------|
| serialNumber | String | Device serial number identifier |
| deviceName | String | User-assigned device name |
| online | Integer | Online status (1=online, 0=offline) |
| deviceModel | Object | Device model information |
| batteryLevel | Integer | Current battery level (0-100) |
| deviceStatus | Integer | Device status code (e.g., 3 for sleep mode) |
| adminId | Integer | ID of the device administrator |
| adminName | String | Name of the device administrator |
| role | Integer | User's role for this device |
| roleName | String | Description of the user's role |
| thumbImgUrl | String | URL to device thumbnail image |
| thumbImgTime | Long | Timestamp of the thumbnail image |
| sdCard | Object | SD card information (if applicable) |
| firmwareId | String | Current firmware version |
| needOta | Integer | OTA update status flag |
| locationId | Integer | ID of the location where device is installed |
| locationName | String | Name of the location |
| deviceSupport | Object | Device capabilities and supported features |

## Response Example
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "serialNumber": "ABC123456789",
        "deviceName": "Front Door Camera",
        "online": 1,
        "deviceModel": {
          "modelName": "G1",
          "canStandby": true,
          "audioCodecType": "aac",
          "isCanRotate": false,
          "supportMotionTrack": true,
          "deviceType": 1
        },
        "batteryLevel": 85,
        "deviceStatus": 1,
        "adminId": 12345,
        "adminName": "John Smith",
        "role": 1,
        "roleName": "Owner",
        "thumbImgUrl": "https://example.com/thumb1.jpg",
        "thumbImgTime": 1617293000,
        "sdCard": {
          "formatStatus": 1,
          "total": 32000,
          "used": 12500
        },
        "firmwareId": "1.2.345",
        "needOta": 0,
        "locationId": 789,
        "locationName": "Home",
        "deviceSupport": {
          "deviceSupportAlarm": true,
          "supportWebrtc": 1,
          "deviceDormancySupport": 1
        }
      },
      {
        "serialNumber": "XYZ987654321",
        "deviceName": "Backyard Camera",
        "online": 1,
        "deviceModel": {
          "modelName": "G2",
          "canStandby": true,
          "audioCodecType": "aac",
          "isCanRotate": true,
          "supportMotionTrack": true,
          "deviceType": 2
        },
        "batteryLevel": 72,
        "deviceStatus": 1,
        "adminId": 12345,
        "adminName": "John Smith",
        "role": 1,
        "roleName": "Owner",
        "thumbImgUrl": "https://example.com/thumb2.jpg",
        "thumbImgTime": 1617293100,
        "sdCard": {
          "formatStatus": 1,
          "total": 64000,
          "used": 18200
        },
        "firmwareId": "1.3.210",
        "needOta": 0,
        "locationId": 789,
        "locationName": "Home",
        "deviceSupport": {
          "deviceSupportAlarm": true,
          "supportWebrtc": 1,
          "deviceDormancySupport": 1
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
| -2 | Authentication error |
| -5 | Network error |

## Usage Context
This endpoint is typically used in the following scenarios:
- Loading the main device dashboard or list view
- Refreshing device list to show current status
- Initial app load to retrieve all devices
- Checking if any devices need firmware updates

## Related Endpoints
- `getSingleDevice` - Gets detailed information about a specific device
- `getDeviceAttributes` - Gets specific device attributes
- `deviceThumbImage` - Gets device thumbnail images

## Implementation Notes
The endpoint retrieves all devices associated with the authenticated user's account, regardless of status or type. The response includes comprehensive information about each device, allowing the client application to display a detailed device list with status indicators, thumbnail images, and other relevant information. The application typically caches this data locally and refreshes it periodically to maintain an up-to-date view of all devices.