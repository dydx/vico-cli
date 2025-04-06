# updateFloodlightSwitch Endpoint Documentation

## Overview
The `updateFloodlightSwitch` endpoint allows users to control the floodlight (turn on/off) for compatible devices. This provides remote control of the device's illumination feature.

## API Details
- **Path**: `/device/updateFloodlightSwitch`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts an `UpdateFloodlightSwitch` object with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | Device serial number identifier |
| switchOn | Boolean | Yes | Flag to turn the floodlight on (true) or off (false) |

The request object also includes standard authentication fields inherited from the `BaseEntry` class, which are automatically added by the API client.

### Example Request Body
```json
{
  "serialNumber": "ABC123XYZ",
  "switchOn": true
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
| -3001 | Device does not support floodlight control |

## Usage in Application
The endpoint is called from the `FloodLightView` class through the following sequence:
1. User taps the floodlight toggle button in the UI
2. An `UpdateFloodlightSwitch` object is created with the device's serial number and the desired state (on/off)
3. The request is made through RxJava using the `DeviceSettingApiClient`
4. The UI shows a loading indicator during the API call
5. On success, the UI updates to reflect the new floodlight state
6. On failure, an error message is displayed to the user

## Constraints
- The device must be online and accessible
- The device must have floodlight functionality
- The user must have appropriate access rights to control the device
- Battery-powered devices may have restrictions on how long the floodlight can remain on