# getLibraryStatus Endpoint Documentation

## Overview
The `getLibraryStatus` endpoint retrieves status information about the user's library, including counts, availability, and other metadata about recordings. It can query both cloud storage and local SD card storage based on the provided parameters.

## API Details
- **Path**: `/library/librarystatus`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `FilterEntry` object that extends `BaseEntry` with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| isFromSDCard | Boolean | No | Whether to fetch from local SD card (true) or cloud (false) |
| serialNumber | Array | No | List of device serial numbers to query |
| startTimestamp | Long | No | Start time for filtering recordings in milliseconds |
| endTimestamp | Long | No | End time for filtering recordings in milliseconds |
| marked | Integer | No | Filter by favorite status (0 or 1) |
| missing | Integer | No | Filter by read/unread status (0 or 1) |
| from | Integer | No | Pagination start parameter |
| to | Integer | No | Pagination end parameter |
| tags | Array | No | Filter by specific tags |
| objectNames | Array | No | Filter by detected objects |
| serialNumberToActivityZone | Object | No | Maps serial numbers to activity zones |
| videoEventKey | String | No | Specific event key to query |
| doorbellTags | Array | No | Doorbell-specific tags to filter by |
| deviceCallEventTag | String | No | Device call event tag to filter by |
| deviceName | String | No | Filter by specific device name |
| app | Object | Yes | Application information (inherited from BaseEntry) |
| countryNo | String | Yes | Country code (e.g., "US") (inherited from BaseEntry) |
| language | String | Yes | Language code (e.g., "en") (inherited from BaseEntry) |
| tenantId | String | Yes | Tenant identifier (inherited from BaseEntry) |

### Example Request Body
```json
{
  "isFromSDCard": false,
  "serialNumber": ["ABC123XYZ", "DEF456UVW"],
  "startTimestamp": 1648756321000,
  "endTimestamp": 1648842721000,
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
The endpoint returns a `LibraryStatusResponse` object that extends `BaseResponse`:

### BaseResponse Fields
| Field | Type | Description |
|-------|------|-------------|
| result | Integer | Status code (0 for success, negative values indicate errors) |
| msg | String | Status message or error description |

### LibraryStatusResponse Additional Fields
| Field | Type | Description |
|-------|------|-------------|
| data.totalCount | Integer | Total number of recordings matching the criteria |
| data.totalUnreadCount | Integer | Number of unread recordings matching the criteria |
| data.hasMoreData | Boolean | Whether more data is available beyond the requested pagination |
| data.deviceStatusList | Array | Status information for each queried device |
| data.cloudStatus | Object | Information about cloud storage status |
| data.sdCardStatus | Object | Information about SD card storage status if available |

### DeviceStatus Object Structure
| Field | Type | Description |
|-------|------|-------------|
| serialNumber | String | Device serial number |
| recordCount | Integer | Count of recordings for this device |
| unreadCount | Integer | Count of unread recordings for this device |
| lastRecordTime | Long | Timestamp of the most recent recording |
| recordingAvailable | Boolean | Whether recording is available for this device |

### CloudStatus Object Structure
| Field | Type | Description |
|-------|------|-------------|
| enabled | Boolean | Whether cloud storage is enabled |
| usedSpace | Long | Amount of used cloud storage space in bytes |
| totalSpace | Long | Total available cloud storage space in bytes |
| expiryDate | Long | Expiration date of cloud storage subscription |

### Example Success Response
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "totalCount": 128,
    "totalUnreadCount": 5,
    "hasMoreData": false,
    "deviceStatusList": [
      {
        "serialNumber": "ABC123XYZ",
        "recordCount": 85,
        "unreadCount": 3,
        "lastRecordTime": 1648842721000,
        "recordingAvailable": true
      },
      {
        "serialNumber": "DEF456UVW",
        "recordCount": 43,
        "unreadCount": 2,
        "lastRecordTime": 1648841721000,
        "recordingAvailable": true
      }
    ],
    "cloudStatus": {
      "enabled": true,
      "usedSpace": 1073741824,
      "totalSpace": 5368709120,
      "expiryDate": 1680378321000
    }
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
1. When the Library tab is opened to get initial status information
2. When the user applies filters to library content
3. To check availability of recordings before attempting to fetch them
4. To display unread counts and notification badges
5. Through the LibraryViewModel.getLibraryStatue method (note the typo in the method name)
6. To determine if cloud storage or local storage is available for a device

## Related Endpoints
- `/library/newselectlibrary/newevent` (getEventRecordByFilter) - Gets recordings filtered by criteria
- `/library/newselectsinglelibrary` (getLibraryByTraceId) - Gets library data by trace ID

## Constraints
- User must be authenticated to access this endpoint
- The response may vary depending on whether cloud or local storage is being queried
- For performance reasons, the response may not include actual recording content, just status information