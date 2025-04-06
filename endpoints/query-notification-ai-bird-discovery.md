# queryNotificationAiBird Endpoint Documentation

## Overview
The `queryNotificationAiBird` endpoint retrieves the current bird detection analysis and notification settings for a specific device. This allows users to check whether bird detection and related notifications are enabled for compatible cameras.

## API Details
- **Path**: `/birdLovers/queryBirdAiSetting`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts an `SNEntry` object with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | Device serial number identifier |
| reason | Integer | No | Optional reason code (defaults to 0) |

### Example Request Body
```json
{
  "serialNumber": "ABC123XYZ",
  "reason": 0
}
```

## Response Structure
The endpoint returns an `AiBirdResponse` object which extends `BaseResponse`:

### Base Response Fields
| Field | Type | Description |
|-------|------|-------------|
| result | int | Status code (0 for success, negative values indicate errors) |
| msg | String | Status message or error description |

### AiBirdResponse Additional Fields
| Field | Type | Description |
|-------|------|-------------|
| data.aiAnalyzeSwitch | boolean | Whether bird detection analysis is enabled on the device |
| data.aiNotifySwitch | boolean | Whether notifications for bird detection events are enabled |

### Example Success Response
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "aiAnalyzeSwitch": true,
    "aiNotifySwitch": false
  }
}
```

### Example Error Response
```json
{
  "result": -2001,
  "msg": "Network error"
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -2001 | Network error |
| -2002 | Access denied |
| -3001 | Device not compatible with bird detection |

## Usage in Application
The endpoint is called from the `DeviceAICore` class through the following sequence:
1. User navigates to the bird detection settings screen in the UI
2. The `getBirdAiSetting` method is called on the `DeviceAICore` instance
3. An `SNEntry` object is created with the device's serial number
4. The request is made through RxJava using the `DeviceSettingApiClient`
5. The response is handled through a callback, updating the UI toggle switches for bird detection analysis and notifications

## Constraints
- The device must be compatible with bird detection features
- The user must have appropriate access rights to the device
- Only certain premium camera models support bird detection
- Bird detection is a subscription feature on some device models