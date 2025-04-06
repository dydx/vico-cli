# updateFloodlightLuminance Endpoint Documentation

## Overview
The `updateFloodlightLuminance` endpoint allows users to adjust the brightness level of a device's floodlight. This provides fine-grained control over the illumination intensity for compatible devices.

## API Details
- **Path**: `/device/updateFloodlightLuminance`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts an `UpdateFloodlightLuminanceRequest` object with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | Device serial number identifier |
| luminance | Integer | Yes | The brightness level for the floodlight (typically 1-100) |

The request object also includes standard authentication fields inherited from the `BaseEntry` class, which are automatically added by the API client.

### Example Request Body
```json
{
  "serialNumber": "ABC123XYZ",
  "luminance": 75
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
| -3001 | Device does not support floodlight brightness control |

## Usage in Application
The endpoint is called from the `FloodLightView` class through the following sequence:
1. User adjusts the brightness slider in the UI
2. When the user stops dragging the slider (onStopTrackingTouch), an `UpdateFloodlightLuminanceRequest` object is created with the device's serial number and the slider's current value
3. The request is made through RxJava using the `DeviceSettingApiClient`
4. On success, the new brightness value is saved in the application's state
5. On failure, the UI slider is reset to the previous value

## Constraints
- The device must be online and accessible
- The device must have a floodlight with adjustable brightness
- The user must have appropriate access rights to control the device
- This endpoint only adjusts brightness and does not turn the floodlight on or off (use updateFloodlightSwitch for that purpose)
- The valid range for the luminance parameter depends on the device model, but is typically between 1 and 100