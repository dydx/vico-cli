# uploadPersonDetectionDetailInfo Endpoint Documentation

## Overview
The `uploadPersonDetectionDetailInfo` endpoint updates AI detection notification settings for a specific device. These settings control which types of objects (person, package, pets, vehicles, etc.) will trigger notifications when detected by the device's AI system.

## API Details
- **Path**: `/device/updateMessageNotification/v1`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `PersonDetectEntry` object with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | Device serial number identifier |
| eventObjectType | Object | Yes | Notification configuration for different object types |
| app | Object | Yes | Application information (inherited from BaseEntry) |
| countryNo | String | Yes | Country code (e.g., "US") (inherited from BaseEntry) |
| language | String | Yes | Language code (e.g., "en") (inherited from BaseEntry) |
| tenantId | String | Yes | Tenant identifier (inherited from BaseEntry) |

### EventObjectType Structure
The `eventObjectType` object contains the following properties:

| Property | Type | Description |
|----------|------|-------------|
| person | Array | List of person detection notification settings |
| pet | Array | List of pet detection notification settings |
| vehicle | Array | List of vehicle detection notification settings |
| packageContent | Array | List of package detection notification settings |
| other | Array | List of other detection notification settings |

Each array item contains:

| Property | Type | Description |
|----------|------|-------------|
| name | String | Detection category name (e.g., "Person", "Familiar Face") |
| choice | Boolean | Whether notifications are enabled for this category |

### Example Request Body
```json
{
  "serialNumber": "ABC123XYZ",
  "eventObjectType": {
    "person": [
      {
        "name": "Person",
        "choice": true
      },
      {
        "name": "Familiar Face",
        "choice": false
      }
    ],
    "pet": [
      {
        "name": "Pet",
        "choice": true
      }
    ],
    "vehicle": [
      {
        "name": "Vehicle",
        "choice": false
      }
    ],
    "packageContent": [
      {
        "name": "Package",
        "choice": true
      }
    ],
    "other": []
  },
  "app": {
    "appName": "vicohome",
    "appVersion": "1.2.3",
    "appBuild": "123",
    "channelId": 1
  },
  "countryNo": "US",
  "language": "en",
  "tenantId": "default"
}
```

## Response Structure
The endpoint returns a `BaseResponse` object:

### BaseResponse Fields
| Field | Type | Description |
|-------|------|-------------|
| result | Integer | Status code (0 for success, negative values indicate errors) |
| msg | String | Status message or error description |

### Example Success Response
```json
{
  "result": 0,
  "msg": "success"
}
```

### Example Error Response
```json
{
  "result": -1001,
  "msg": "Device not found"
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -1001 | Device not found |
| -2001 | Network error |
| -3001 | Invalid parameters |
| -4001 | Authentication error |

## Usage in Application
The endpoint is called when the user updates notification settings for AI-detected objects:
1. User toggles detection notification settings in the device settings UI
2. The application constructs a `PersonDetectEntry` object with the updated settings
3. The request is made through the `DeviceAICore.updateMessageNotificationConfig()` method
4. On success, the updated settings are reflected in the UI
5. On failure, an appropriate error message is displayed to the user

## Related Endpoints
- `/device/queryMessageNotification/v1` (loadPersonDetectionDetailInfo) - Gets the current AI detection notification settings for a device

## Constraints
- The device must exist and be associated with the user's account
- The serialNumber parameter must not be empty
- The user must have permission to modify device settings