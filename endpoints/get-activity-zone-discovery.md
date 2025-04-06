# getActivityZone Endpoint Documentation

## Overview
The `getActivityZone` endpoint retrieves all configured activity zones for a specific device. Activity zones are user-defined areas within a camera's field of view where motion detection and alerts are active.

## API Details
- **Path**: `/device/getactivityzone`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `ZoneGetEntry` object with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | Device serial number identifier |
| requestId | String | No | Unique identifier for the request |
| language | String | No | Language parameter for localized responses |

### Example Request Body
```json
{
  "serialNumber": "ABC123XYZ"
}
```

## Response Structure
The endpoint returns a `ZoneGetResponse` object which extends `BaseResponse`:

### Base Response Fields
| Field | Type | Description |
|-------|------|-------------|
| result | int | Status code (0 for success, negative values indicate errors) |
| msg | String | Status message or error description |

### ZoneGetResponse Additional Fields
| Field | Type | Description |
|-------|------|-------------|
| data.list | List<ZoneBean> | List of activity zones for the specified device |

### ZoneBean Structure
| Field | Type | Description |
|-------|------|-------------|
| id | int | Unique identifier for the zone |
| serialNumber | String | Device serial number |
| zoneName | String | Name of the activity zone |
| vertices | String | Comma-separated coordinates (x1,y1,x2,y2,...) defining the zone shape |
| needRecord | int | Flag to record events in this zone (1=enabled, 0=disabled) |
| needAlarm | int | Flag to trigger alarm for events (1=enabled, 0=disabled) |
| needPush | int | Flag to send push notifications (1=enabled, 0=disabled) |
| zoneNameId | int | Identifier for the zone name |
| modelCategory | int | Model category identifier |
| deviceName | String | Name of the device |

### Example Success Response
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 12345,
        "serialNumber": "ABC123XYZ",
        "zoneName": "Front Door",
        "vertices": "0.2,0.2,0.2,0.5,0.5,0.5,0.5,0.2",
        "needRecord": 1,
        "needAlarm": 0,
        "needPush": 1,
        "zoneNameId": 1,
        "modelCategory": 1,
        "deviceName": "Front Camera"
      },
      {
        "id": 12346,
        "serialNumber": "ABC123XYZ",
        "zoneName": "Driveway",
        "vertices": "0.6,0.6,0.6,0.8,0.8,0.8,0.8,0.6",
        "needRecord": 1,
        "needAlarm": 1,
        "needPush": 1,
        "zoneNameId": 2,
        "modelCategory": 1,
        "deviceName": "Front Camera"
      }
    ]
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
The endpoint is called from the `DeviceActivityZoneCore` class through the following sequence:
1. User navigates to the zone management screen in the UI
2. The `getActivityZone` method is called on the `DeviceActivityZoneCore` instance
3. The request is made through RxJava using the `DeviceSettingApiClient`
4. The response is processed in the `ZoneViewModel` and displayed in the `ZoneActivity`

## Constraints
- The device must be online and accessible to retrieve its zone configuration
- The user must have appropriate access rights to the device
- Only zones created for the specified device will be returned
- The response will be an empty list if no zones are configured for the device