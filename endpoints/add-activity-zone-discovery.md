# addActivityZone Endpoint Documentation

## Overview
The `addActivityZone` endpoint allows users to create custom activity zones for device monitoring. These zones define specific areas within a camera's view where motion detection and alerts should be active.

## API Details
- **Path**: `/device/insertactivityzone`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `ZoneBean` object with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | Device serial number identifier |
| zoneName | String | Yes | Name for the activity zone |
| vertices | String | Yes | Comma-separated list of coordinates that define the zone's shape (format: "x1,y1,x2,y2,...") |
| needRecord | int | Yes | Flag to record events in this zone (1 = enabled, 0 = disabled) |
| needAlarm | int | No | Flag to trigger alarm for events in this zone (1 = enabled, 0 = disabled) |
| needPush | int | Yes | Flag to send push notifications for events in this zone (1 = enabled, 0 = disabled) |

### Example Request Body
```json
{
  "serialNumber": "ABC123XYZ",
  "zoneName": "Zone 1",
  "vertices": "0.2,0.2,0.2,0.5,0.2,0.8,0.5,0.8,0.8,0.8,0.8,0.5,0.8,0.2,0.5,0.2",
  "needRecord": 1,
  "needPush": 1
}
```

## Response Structure
The endpoint returns a `ZoneAddResponse` object which extends `BaseResponse`:

### Base Response Fields
| Field | Type | Description |
|-------|------|-------------|
| result | int | Status code (negative values indicate errors) |
| msg | String | Status message or error description |

### ZoneAddResponse Additional Fields
| Field | Type | Description |
|-------|------|-------------|
| data.id | int | ID of the newly created activity zone |

### Example Success Response
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "id": 12345
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

## Usage in Application
The endpoint is called from the `ZoneActivity` class through the following sequence:
1. User configures a zone in the UI by drawing a zone shape
2. The `createActivityZone` method is called on the `DeviceActivityZoneCore` instance
3. The defined zone parameters are packaged into a `ZoneBean` object
4. The request is made through RxJava using the `DeviceSettingApiClient` which wraps the API call
5. The response is observed in the `ZoneViewModel` and updates the UI accordingly

## Constraints
- The app enforces a limit of 3 activity zones per device
- Zone coordinates are represented as relative values (0.0 to 1.0) rather than absolute pixel coordinates
- Zone shapes are polygons defined by their vertex coordinates