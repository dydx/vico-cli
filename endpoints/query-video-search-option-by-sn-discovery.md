# queryVideoSearchOptionBySn Endpoint Documentation

## Overview
The `queryVideoSearchOptionBySn` endpoint retrieves available video search and filter options for recordings from a specific device identified by its serial number. It returns categorized tags, events, and operation options that can be used to search and filter video content for that particular device.

## API Details
- **Path**: `/library/queryVideoSearchOptionBySn`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `LibraryOptionRequestBySn` object that extends `BaseEntry` with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| sn | String | Yes | Serial number of the device to get search options for |
| app | Object | Yes | Application information (inherited from BaseEntry) |
| countryNo | String | Yes | Country code (e.g., "US") (inherited from BaseEntry) |
| language | String | Yes | Language code (e.g., "en") (inherited from BaseEntry) |
| tenantId | String | Yes | Tenant identifier (inherited from BaseEntry) |

### Example Request Body
```json
{
  "sn": "ABC123XYZ",
  "app": {
    "appName": "vicohome",
    "appVersion": "1.2.3",
    "appBuild": "123",
    "channelId": 1
  },
  "countryNo": "US",
  "language": "en",
  "tenantId": "default"
}
```

## Response Structure
The endpoint returns a `VideoSearchOptionResponse` object that extends `BaseResponse`:

### BaseResponse Fields
| Field | Type | Description |
|-------|------|-------------|
| result | Integer | Status code (0 for success, negative values indicate errors) |
| msg | String | Status message or error description |

### VideoSearchOptionResponse Additional Fields
| Field | Type | Description |
|-------|------|-------------|
| data | Object | TagBean object containing search and filter options |
| data.aiEventTags | Array | List of AI event tags for this device (e.g., person, vehicle, package detection) |
| data.deviceEventTags | Array | List of device event tags for this device (e.g., motion, sound, doorbell press) |
| data.devices | Array | Information about the requested device (typically only contains one device) |
| data.operateOptions | Array | List of operation options (e.g., favorites, unread) |

### OptionTag Object Structure
| Field | Type | Description |
|-------|------|-------------|
| name | String | Tag name or identifier |
| displayName | String | Localized display name for the tag |
| count | Integer | Count of recordings with this tag for this device |
| checked | Boolean | Whether the tag is selected by default |
| subTags | Array | List of nested sub-tags with the same structure |

### OptionDevice Object Structure
| Field | Type | Description |
|-------|------|-------------|
| serialNumber | String | Device serial number |
| deviceName | String | User-assigned device name |
| modelCategory | Integer | Device model category |
| roleId | Integer | User role ID (1 for admin) |
| isBind | Boolean | Whether the device is bound to the user |
| online | Integer | Online status (1=online, 0=offline) |
| activityZone | Array | List of activity zones for this device |

### Example Success Response
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "aiEventTags": [
      {
        "name": "person",
        "displayName": "Person",
        "count": 42,
        "checked": false,
        "subTags": []
      },
      {
        "name": "vehicle",
        "displayName": "Vehicle",
        "count": 15,
        "checked": false,
        "subTags": []
      }
    ],
    "deviceEventTags": [
      {
        "name": "motion",
        "displayName": "Motion",
        "count": 87,
        "checked": false,
        "subTags": []
      },
      {
        "name": "doorbell",
        "displayName": "Doorbell Press",
        "count": 23,
        "checked": false,
        "subTags": []
      }
    ],
    "devices": [
      {
        "serialNumber": "ABC123XYZ",
        "deviceName": "Front Door Camera",
        "modelCategory": 1,
        "roleId": 1,
        "isBind": true,
        "online": 1,
        "activityZone": [
          {
            "name": "driveway",
            "displayName": "Driveway",
            "count": 35,
            "checked": false,
            "subTags": []
          }
        ]
      }
    ],
    "operateOptions": [
      {
        "name": "favorite",
        "displayName": "Favorites",
        "count": 12,
        "checked": false,
        "subTags": []
      },
      {
        "name": "unread",
        "displayName": "Unread",
        "count": 5,
        "checked": false,
        "subTags": []
      }
    ]
  }
}
```

### Example Error Response
```json
{
  "result": -1001,
  "msg": "Device not found"
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -1001 | Device not found |
| -1002 | Invalid serial number |
| -2001 | Network error |
| -4001 | Authentication error |

## Usage in Application
The endpoint is called in these scenarios:
1. When the user selects a specific device in the library view
2. When filtering video content for a particular device
3. In the LibraryCore.queryVideoTagsByCloud method to retrieve search options from the cloud
4. To populate device-specific filter dropdowns and options
5. Before displaying recordings from a specific device to get relevant filter categories

## Related Endpoints
- `/library/queryVideoSearchOption` (queryVideoSearchOption) - Gets video search options for all devices
- `/library/newselectlibrary/newevent` (getEventRecordByFilter) - Gets recordings filtered by criteria

## Constraints
- User must be authenticated to access this endpoint
- The serial number must be valid and associated with the user's account
- This endpoint retrieves options only from cloud storage (not SD card)
- The response is specific to the requested device rather than providing options for all devices