# getEventDetail Endpoint Documentation

## Overview
The `getEventDetail` endpoint retrieves detailed information about a specific event or recording in the user's library. It can fetch details from either cloud storage or local SD card storage based on the parameters provided.

## API Details
- **Path**: `/library/newselectlibrary`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `FilterEntry` object that extends `BaseEntry` with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| videoEventKey | String | Yes | The unique identifier for the video event to retrieve details for |
| isFromSDCard | Boolean | No | Flag indicating whether to fetch from local SD card (true) or cloud (false) |
| startTimestamp | Long | No | Start time for filtering in milliseconds |
| endTimestamp | Long | No | End time for filtering in milliseconds |
| marked | Integer | No | Filter by favorite status (0 or 1) |
| missing | Integer | No | Filter by read/unread status (0 or 1) |
| from | Integer | No | Pagination start parameter |
| to | Integer | No | Pagination end parameter |
| serialNumber | Array | No | List of device serial numbers to filter by |
| tags | Array | No | Event tags to filter by |
| serialNumberToActivityZone | Object | No | Maps serial numbers to activity zones |
| doorbellTags | Array | No | Doorbell-specific tags for filtering |
| deviceCallEventTag | String | No | Device call event tag for filtering |
| deviceName | String | No | Device name for filtering |
| objectNames | Array | No | List of detected object names to filter by (e.g., "person", "vehicle") |
| app | Object | Yes | Application information (inherited from BaseEntry) |
| countryNo | String | Yes | Country code (e.g., "US") (inherited from BaseEntry) |
| language | String | Yes | Language code (e.g., "en") (inherited from BaseEntry) |
| tenantId | String | Yes | Tenant identifier (inherited from BaseEntry) |

### Example Request Body
```json
{
  "videoEventKey": "event_123456789",
  "isFromSDCard": false,
  "serialNumber": ["ABC123XYZ"],
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
The endpoint returns a `LibraryRecordResponse` object that extends `BaseResponse`:

### BaseResponse Fields
| Field | Type | Description |
|-------|------|-------------|
| result | Integer | Status code (0 for success, negative values indicate errors) |
| msg | String | Status message or error description |

### LibraryRecordResponse Additional Fields
| Field | Type | Description |
|-------|------|-------------|
| data.list | Array | List of recording details (typically contains only one item when querying by videoEventKey) |
| data.total | Integer | Total count of matching records (typically 1 for this endpoint) |

### Recording Detail Object Structure
| Field | Type | Description |
|-------|------|-------------|
| traceId | String | Unique recording identifier |
| serialNumber | String | Device serial number |
| eventTime | Long | Timestamp of event in milliseconds |
| eventType | Integer | Type of event (numeric code) |
| thumbnailUrl | String | URL to thumbnail image |
| videoUrl | String | URL to video recording |
| duration | Integer | Recording length in seconds |
| marked | Integer | Favorite status (0=not marked, 1=marked as favorite) |
| read | Integer | Viewed status (0=unread, 1=read) |
| activityZone | String | Activity zone ID if applicable |
| detectedObjects | Array | List of detected objects in the recording |
| tags | Array | List of event tags |
| deviceName | String | Name of the device that recorded the event |

### Example Success Response
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "traceId": "abc123def456",
        "serialNumber": "ABC123XYZ",
        "eventTime": 1648756321000,
        "eventType": 1,
        "thumbnailUrl": "https://cdn.vicohome.io/thumbnails/abc123.jpg",
        "videoUrl": "https://cdn.vicohome.io/videos/abc123.mp4",
        "duration": 15,
        "marked": 0,
        "read": 0,
        "activityZone": "zone1",
        "detectedObjects": ["person"],
        "tags": ["motion"],
        "deviceName": "Front Door Camera"
      }
    ],
    "total": 1
  }
}
```

### Example Error Response
```json
{
  "result": -1001,
  "msg": "Recording not found"
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -1001 | Recording not found |
| -1002 | Invalid video event key |
| -2001 | Network error |
| -4001 | Authentication error |

## Usage in Application
The endpoint is called in these scenarios:
1. When a user taps on a recording in the library to view its details
2. When the app needs to fetch full information about a recording before playback
3. In the LibraryCore.getEventDetail method, which handles the appropriate routing to either cloud or local storage based on the isFromSDCard flag
4. Through the LibraryViewModel, which acts as an intermediary between the UI and the API

## Related Endpoints
- `/library/newselectlibrary/newevent` (getEventRecordByFilter) - Gets recordings filtered by criteria
- `/library/setmarkstatus` (setMarkStatus) - Marks recordings as favorites
- `/library/setreadstatus` (setReadStatus) - Sets read status for recordings

## Constraints
- User must be authenticated to access this endpoint
- The videoEventKey must be valid and associated with the user's account
- The appropriate storage location (cloud or SD card) must be specified correctly through the isFromSDCard parameter
- This endpoint uses the same base path as getEventRecordByFilter but retrieves a single specific recording's details