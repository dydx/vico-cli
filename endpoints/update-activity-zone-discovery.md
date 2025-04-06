# updateActivityZone Endpoint Documentation

## Overview
The `updateActivityZone` endpoint allows users to modify existing activity zones for device monitoring. These zones define specific areas within a camera's view where motion detection and alerts should be active.

## API Details
- **Path**: `/device/updateactivityzone`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `ZoneBean` object with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | Integer | Yes | Unique identifier for the zone to update |
| serialNumber | String | Yes | Device serial number identifier |
| zoneName | String | Yes | Name for the activity zone |
| vertices | String | Yes | Comma-separated list of coordinates that define the zone's shape (format: "x1,y1,x2,y2,...") |
| needRecord | Integer | Yes | Flag to record events in this zone (1 = enabled, 0 = disabled) |
| needPush | Integer | Yes | Flag to send push notifications for events in this zone (1 = enabled, 0 = disabled) |
| needAlarm | Integer | No | Flag to trigger alarm for events in this zone (1 = enabled, 0 = disabled) |
| zoneNameId | Integer | No | Identifier for the zone name |
| modelCategory | Integer | No | Model category identifier |
| deviceName | String | No | Name of the device |

### Example Request Body
```json
{
  "id": 12345,
  "serialNumber": "ABC123XYZ",
  "zoneName": "Front Yard",
  "vertices": "0.2,0.2,0.2,0.5,0.2,0.8,0.5,0.8,0.8,0.8,0.8,0.5,0.8,0.2,0.5,0.2",
  "needRecord": 1,
  "needPush": 1,
  "needAlarm": 0
}
```

## Response Structure
The endpoint returns a `ZoneUpdateResponse` object which extends `BaseResponse`:

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
The endpoint is called from the `ZoneViewModel` class through the following sequence:
1. User modifies a zone in the UI
2. The modified zone properties are packaged into a `ZoneBean` object
3. The `updateActivityZone` method is called on the `ZoneViewModel` instance
4. The `ZoneBean` is passed to the `DeviceActivityZoneCore` which makes the API call
5. The request is made through RxJava using the `DeviceSettingApiClient`
6. The response is observed through LiveData in the `ZoneViewModel`
7. The UI is updated based on the success or failure of the operation

## Constraints
- The app enforces a limit of 3 activity zones per device
- Zone coordinates are represented as relative values (0.0 to 1.0) rather than absolute pixel coordinates
- Zone shapes are polygons defined by their vertex coordinates 
- The zone ID must match an existing zone for the specified device
- The user must have appropriate permissions to modify the zone