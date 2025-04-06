# queryUserActivityZone Endpoint Documentation

## Overview
The `queryUserActivityZone` endpoint retrieves all activity zones configured by the user across all their devices. Unlike the device-specific getActivityZone endpoint, this provides a global view of all zones in the user's account.

## API Details
- **Path**: `/device/queryUserActivityZone`
- **Method**: GET
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
This endpoint doesn't require any specific parameters in the request body beyond standard authentication tokens, as it retrieves all zones for the authenticated user.

## Response Structure
The endpoint returns a `UserAllAZResponse` object that extends `BaseResponse`:

### BaseResponse Fields
| Field | Type | Description |
|-------|------|-------------|
| result | Integer | Status code (0 for success, negative values indicate errors) |
| msg | String | Status message or error description |

### UserAllAZResponse Additional Fields
| Field | Type | Description |
|-------|------|-------------|
| data | Array | List of ZoneBean objects representing all activity zones |

### ZoneBean Object Structure
| Field | Type | Description |
|-------|------|-------------|
| id | Integer | Unique identifier for the zone |
| serialNumber | String | Device serial number |
| zoneName | String | Name of the activity zone |
| vertices | String | Comma-separated coordinates defining the zone shape (format: "x1,y1,x2,y2,...") |
| needRecord | Integer | Flag to record events in this zone (1=enabled, 0=disabled) |
| needAlarm | Integer | Flag to trigger alarm for events (1=enabled, 0=disabled) |
| needPush | Integer | Flag to send push notifications (1=enabled, 0=disabled) |
| zoneNameId | Integer | Identifier for the zone name |
| modelCategory | Integer | Model category identifier |
| deviceName | String | Name of the device |
| deviceBean | Object | Additional device information |
| isSelect | Boolean | Selection state (used in UI) |

### Example Success Response
```json
{
  "result": 0,
  "msg": "success",
  "data": [
    {
      "id": 12345,
      "serialNumber": "ABC123XYZ",
      "zoneName": "Front Yard",
      "vertices": "0.2,0.2,0.2,0.5,0.2,0.8,0.5,0.8,0.8,0.8,0.8,0.5,0.8,0.2,0.5,0.2",
      "needRecord": 1,
      "needAlarm": 0,
      "needPush": 1,
      "zoneNameId": 1,
      "modelCategory": 1,
      "deviceName": "Front Door Camera"
    },
    {
      "id": 67890,
      "serialNumber": "DEF456UVW",
      "zoneName": "Driveway",
      "vertices": "0.3,0.3,0.3,0.6,0.3,0.9,0.6,0.9,0.9,0.9,0.9,0.6,0.9,0.3,0.6,0.3",
      "needRecord": 1,
      "needAlarm": 1,
      "needPush": 1,
      "zoneNameId": 2,
      "modelCategory": 1,
      "deviceName": "Driveway Camera"
    }
  ]
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
| -1001 | No zones found |
| -2001 | Network error |
| -4001 | Authentication error |

## Usage in Application
The endpoint is called in these scenarios:
1. When loading the global activity zone management screen
2. In the ZoneViewModel class to populate a list of all zones
3. To provide filtering options for recordings based on zones across multiple devices
4. When the app needs to display a comprehensive view of all configured zones

## Implementation Details
1. The endpoint is called through LibraryApiClient's queryUserActivityZone method
2. LibraryCore processes the response and converts it to a list of ZoneBean objects
3. The ZoneViewModel receives the data and updates its observable state
4. UI components observe this state to display zone information

## Related Endpoints
- `/device/getactivityzone` (getActivityZone) - Gets activity zones for a specific device
- `/device/insertactivityzone` (addActivityZone) - Adds an activity zone to a device
- `/device/updateactivityzone` (updateActivityZone) - Updates an existing activity zone
- `/device/deleteactivityzone` (deleteActivityZone) - Deletes an activity zone

## Constraints
- User must be authenticated to access this endpoint
- The response includes zones from all devices the user has permission to access
- For performance reasons, bitmap representations of zones are not included in the response