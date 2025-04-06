# getLibraryByTraceId Endpoint Documentation

## Overview
The `getLibraryByTraceId` endpoint retrieves a single library item (recording or image) by its unique trace ID. This endpoint provides comprehensive details about the specific media item, including URLs, metadata, and status information.

## API Details
- **Path**: `/library/newselectsinglelibrary`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `TraceIdEntry` object that extends `BaseEntry` with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| traceId | String | Yes | The unique identifier for the library item to retrieve |
| app | Object | Yes | Application information (inherited from BaseEntry) |
| countryNo | String | Yes | Country code (e.g., "US") (inherited from BaseEntry) |
| language | String | Yes | Language code (e.g., "en") (inherited from BaseEntry) |
| tenantId | String | Yes | Tenant identifier (inherited from BaseEntry) |

### Example Request Body
```json
{
  "traceId": "trace123abc",
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
The endpoint returns a `SingleLibraryResponse` object that extends `BaseResponse`:

### BaseResponse Fields
| Field | Type | Description |
|-------|------|-------------|
| result | Integer | Status code (0 for success, negative values indicate errors) |
| msg | String | Status message or error description |

### SingleLibraryResponse Additional Fields
| Field | Type | Description |
|-------|------|-------------|
| data | Object | Detailed RecordBean object with library item information |

### RecordBean Object Structure
| Field | Type | Description |
|-------|------|-------------|
| traceId | String | Unique identifier for the library item |
| imageUrl | String | URL to the image thumbnail |
| videoUrl | String | URL to the video recording |
| deviceName | String | Name of the device that recorded the media |
| serialNumber | String | Serial number of the device |
| timestamp | Integer | Timestamp of the recording in seconds since epoch |
| date | String | Formatted date string of the recording |
| fileSize | Integer | Size of the media file in bytes |
| mediaType | String | Type of media (e.g., "video" or "image") |
| marked | Integer | Whether the item is marked as favorite (0=no, 1=yes) |
| missing | Integer | Whether the item has been viewed (0=unread, 1=read) |
| type | Integer | Type code for the media item |
| tags | String | Comma-separated tags associated with the media |
| activityZoneName | String | Name of the activity zone if applicable |
| locationName | String | Name of the location where the device is installed |
| duration | Integer | Length of the video in seconds |
| eventType | Integer | Type of event that triggered the recording |
| objectName | String | Objects detected in the recording (comma-separated) |
| activityZone | String | ID of the activity zone |
| width | Integer | Width of the media in pixels |
| height | Integer | Height of the media in pixels |

### Example Success Response
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "traceId": "trace123abc",
    "imageUrl": "https://cdn.vicohome.io/thumbnails/trace123abc.jpg",
    "videoUrl": "https://cdn.vicohome.io/videos/trace123abc.mp4",
    "deviceName": "Front Door Camera",
    "serialNumber": "ABC123XYZ",
    "timestamp": 1648756321,
    "date": "2022-04-01 10:12:01",
    "fileSize": 5242880,
    "mediaType": "video",
    "marked": 0,
    "missing": 1,
    "type": 1,
    "tags": "motion",
    "activityZoneName": "Driveway",
    "locationName": "Home",
    "duration": 15,
    "eventType": 1,
    "objectName": "person",
    "activityZone": "zone1",
    "width": 1920,
    "height": 1080
  }
}
```

### Example Error Response
```json
{
  "result": -1001,
  "msg": "Library item not found"
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -1001 | Library item not found |
| -1002 | Invalid trace ID |
| -2001 | Network error |
| -4001 | Authentication error |

## Usage in Application
The endpoint is called in these scenarios:
1. When a user accesses a specific recording from a notification
2. When the app needs to retrieve detailed information about a specific recording
3. In the NotifyHandlerViewModel to handle media from notifications
4. When the app needs to display a single library item with all its metadata
5. Through the LibraryCore.getLibraryByTraceId method, which provides proper handling and callbacks

## Related Endpoints
- `/library/newselectlibrary` (getEventDetail) - Gets details for an event
- `/library/newselectlibrary/newevent` (getEventRecordByFilter) - Gets recordings filtered by criteria

## Constraints
- User must be authenticated to access this endpoint
- The traceId must be valid and associated with the user's account
- This endpoint retrieves a single specific library item rather than filtering multiple items