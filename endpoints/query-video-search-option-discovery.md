# queryVideoSearchOption Endpoint Documentation

## Overview
The `queryVideoSearchOption` endpoint retrieves available video search and filter options for recordings in the user's library. It returns categorized tags, events, and device options that can be used to search and filter video content.

## API Details
- **Path**: `/library/queryVideoSearchOption`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `LibraryOptionRequest` object that extends `BaseEntry` with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| SN | String | No | Serial number of a specific device to filter options for |
| isFromSDCard | Boolean | No | Flag indicating whether to fetch options from SD card (true) or cloud (false) |
| app | Object | Yes | Application information (inherited from BaseEntry) |
| countryNo | String | Yes | Country code (e.g., "US") (inherited from BaseEntry) |
| language | String | Yes | Language code (e.g., "en") (inherited from BaseEntry) |
| tenantId | String | Yes | Tenant identifier (inherited from BaseEntry) |

### Example Request Body
```json
{
  "SN": "ABC123XYZ",
  "isFromSDCard": false,
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
| data.aiEventTags | Array | List of AI event tags (e.g., person, vehicle, package detection) |
| data.deviceEventTags | Array | List of device event tags (e.g., motion, sound, doorbell press) |
| data.devices | Array | List of devices available for filtering |
| data.operateOptions | Array | List of operation options (e.g., favorites, unread) |

### OptionTag Object Structure
| Field | Type | Description |
|-------|------|-------------|
| name | String | Tag name or identifier |
| displayName | String | Localized display name for the tag |
| count | Integer | Count of recordings with this tag (if available) |
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
| activityZone | Array | List of activity zones for this device (if available) |

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
      },
      {
        "name": "package",
        "displayName": "Package",
        "count": 8,
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
  "result": -2001,
  "msg": "Network error"
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -1001 | Invalid parameters |
| -2001 | Network error |
| -4001 | Authentication error |

## Usage in Application
The endpoint is called in these scenarios:
1. When the user opens the Library tab to get initial filter options
2. When the user selects a different storage location (cloud or SD card)
3. In the LibraryViewModel's getFilterOption method
4. Through the LibraryCore's queryVideoSearchOption method
5. To populate filter controls in the library interface

## Related Endpoints
- `/library/queryVideoSearchOptionBySn` (queryVideoSearchOptionBySn) - Gets video search options for a specific device
- `/library/newselectlibrary/newevent` (getEventRecordByFilter) - Gets recordings filtered by criteria

## Constraints
- User must be authenticated to access this endpoint
- The available options may differ between cloud storage and SD card storage
- If a specific device serial number is provided, the options will be tailored to that device
- The response may be cached to improve performance when frequently accessing filter options