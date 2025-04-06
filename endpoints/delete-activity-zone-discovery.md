# deleteActivityZone Endpoint Documentation

## Overview
The `deleteActivityZone` endpoint allows users to delete custom activity zones for device monitoring. Users can remove previously created zones when they are no longer needed.

## API Details
- **Path**: `/device/deleteactivityzone`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `ZoneDeleteEntry` object with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | Integer | Yes | The ID of the activity zone to delete |
| serialNumber | String | Yes | Device serial number identifier |
| language | String | No | Language parameter for localized responses |
| requestId | String | No | Unique identifier for the request |

### Example Request Body
```json
{
  "id": 12345,
  "serialNumber": "ABC123XYZ"
}
```

## Response Structure
The endpoint returns a `ZoneDeleteResponse` object which extends `BaseResponse`:

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

## Usage in Application
The endpoint is called from the `DeviceActivityZoneCore` class through the following sequence:
1. User selects a zone for deletion in the UI
2. The `deleteActivityZone` method is called on the `DeviceActivityZoneCore` instance
3. The zone ID and device serial number are packaged into a `ZoneDeleteEntry` object
4. The request is made through a callback pattern using the `DeviceSettingApiClient`
5. The response is handled in callback methods that update the UI accordingly

## Constraints
- Only zones created by the user can be deleted
- The zone ID must be valid and associated with the specified device
- Deleting a zone cannot be undone