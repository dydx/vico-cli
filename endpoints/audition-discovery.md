# audition Endpoint Documentation

## Overview
The `audition` endpoint triggers a doorbell's chime to play a test sound. This allows users to verify the doorbell's audio functionality and adjust volume settings accordingly.

## API Details
- **Path**: `/device/mechanical/dingdong/audition`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `CommonSNRequest` object with the following field:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | Device serial number identifier |

### Example Request Body
```json
{
  "serialNumber": "ABC123XYZ"
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
| -3001 | Device not compatible with audition feature |
| -3002 | Device offline |

## Usage in Application
The endpoint is called from the `RingtoneSettingFragment` through the following sequence:
1. User taps the play button in the chime settings screen
2. The `playAudition` method is called in the `RingtoneSettingFragment`
3. The fragment calls the `audition` method on the `DeviceConfigViewModel`
4. The view model creates a `CommonSNRequest` with the device's serial number
5. The request is sent through the `DeviceSettingCore` to the API
6. The UI shows a rotating animation while the audition is in progress
7. On completion or error, the UI reverts to showing the play button

## Constraints
- The device must be online and accessible
- The feature is only available on compatible doorbell models
- There may be a cooldown period between consecutive audition requests
- The audition plays the current doorbell chime at the current volume settings