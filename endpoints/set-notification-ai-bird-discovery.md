# setNotificationAiBird Endpoint Documentation

## Overview
The `setNotificationAiBird` endpoint allows users to update the bird detection analysis and notification settings for a specific device. This enables users to control whether bird detection is active and whether they receive notifications for detected birds.

## API Details
- **Path**: `/birdLovers/updateBirdAiSetting`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts an `AiBirdUpdateEntry` object with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | Device serial number identifier |
| aiAnalyzeSwitch | Boolean | Yes | Flag to enable/disable bird detection analysis |
| aiNotifySwitch | Boolean | Yes | Flag to enable/disable notifications for bird detection |
| requestId | String | No | Unique identifier for the request |
| language | String | No | Language parameter for localized responses |

### Example Request Body
```json
{
  "serialNumber": "ABC123XYZ",
  "aiAnalyzeSwitch": true,
  "aiNotifySwitch": false
}
```

## Response Structure
The endpoint returns a standard `BaseResponse` object:

### Base Response Fields
| Field | Type | Description |
|-------|------|-------------|
| result | int | Status code (0 for success, negative values indicate errors) |
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
1. User toggles bird detection settings in the UI
2. The `updateBirdAiSetting` method is called on the `DeviceAICore` instance
3. An `AiBirdUpdateEntry` object is created with the device's serial number and the desired settings
4. The request is made through RxJava using the `DeviceSettingApiClient`
5. The response is handled through a callback, providing user feedback about the result

## Constraints
- The device must be compatible with bird detection features
- The user must have appropriate access rights to the device
- Only certain premium camera models support bird detection
- If `aiAnalyzeSwitch` is set to false, `aiNotifySwitch` is automatically treated as false as well
- Bird detection is a subscription feature on some device models