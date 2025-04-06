# getEventRecordByFilter Endpoint

## Overview
The getEventRecordByFilter endpoint retrieves video recordings and events from a device's library based on specified filters, including time periods. This endpoint allows clients to query recordings by date ranges, event types, and other criteria.

## API Details
- **Path**: `/library/newselectlibrary/newevent`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Retrieves recording events filtered by time period and other criteria.

## Request Parameters
The request body should contain a FilterEntry JSON object with the following properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| startTimestamp | Long | Yes | Start time for filtering events (milliseconds since epoch) |
| endTimestamp | Long | Yes | End time for filtering events (milliseconds since epoch) |
| serialNumber | Array | Yes | List of device serial numbers to filter by |
| from | Integer | No | Starting index for pagination |
| to | Integer | No | Ending index for pagination |
| tags | Array | No | List of event tags to filter by |
| objectNames | Array | No | List of detected object names to filter by |
| marked | Integer | No | Filter by marked/favorite status (0 or 1) |
| missing | Integer | No | Filter by read/unread status |
| serialNumberToActivityZone | Object | No | Map linking serial numbers to activity zones |
| videoEventKey | String | No | Key for specific video events |
| isFromSDCard | Boolean | No | Flag indicating if content should be fetched from local storage |

Plus standard BaseEntry properties:
- app
- countryNo
- language
- tenantId

## Request Example
```json
{
  "startTimestamp": 1649030400000,
  "endTimestamp": 1649116800000,
  "serialNumber": ["ABC123456789"],
  "from": 0,
  "to": 20,
  "tags": ["motion", "person"],
  "marked": 0,
  "isFromSDCard": false,
  "app": "vicohome",
  "countryNo": "US",
  "language": "en",
  "tenantId": "AAA12345678"
}
```

## Response Structure
The response contains a list of events matching the filter criteria:

| Property | Type | Description |
|----------|------|-------------|
| result | Integer | Status code (0 for success, negative values for errors) |
| msg | String | Status message or error description |
| data | Object | Contains the list of recordings and total count |

### Data Object Structure
| Property | Type | Description |
|----------|------|-------------|
| list | Array | List of RecordBean objects representing recordings |
| total | Integer | Total number of records matching the filter criteria |

### RecordBean Structure
Each recording in the list has the following properties:

| Property | Type | Description |
|----------|------|-------------|
| traceId | String | Unique identifier for the recording |
| serialNumber | String | Device serial number |
| eventTime | Long | Timestamp when the event occurred |
| eventType | Integer | Type of event (motion, person detection, etc.) |
| thumbnailUrl | String | URL to a thumbnail image for the event |
| videoUrl | String | URL to the video recording |
| duration | Integer | Length of the recording in seconds |
| marked | Integer | Whether the recording is marked/favorited (0 or 1) |
| read | Integer | Whether the recording has been viewed (0 or 1) |
| activityZone | Integer | ID of the activity zone that triggered the event |
| detectedObjects | Array | List of objects detected in the recording |

## Response Example
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "traceId": "ABC123-1649050200000",
        "serialNumber": "ABC123456789",
        "eventTime": 1649050200000,
        "eventType": 1,
        "thumbnailUrl": "https://cloud.vicohome.io/thumb/ABC123-1649050200000.jpg",
        "videoUrl": "https://cloud.vicohome.io/video/ABC123-1649050200000.mp4",
        "duration": 15,
        "marked": 0,
        "read": 1,
        "activityZone": 2,
        "detectedObjects": ["person"]
      },
      {
        "traceId": "ABC123-1649060100000",
        "serialNumber": "ABC123456789",
        "eventTime": 1649060100000,
        "eventType": 2,
        "thumbnailUrl": "https://cloud.vicohome.io/thumb/ABC123-1649060100000.jpg",
        "videoUrl": "https://cloud.vicohome.io/video/ABC123-1649060100000.mp4",
        "duration": 12,
        "marked": 0,
        "read": 0,
        "activityZone": 1,
        "detectedObjects": ["motion"]
      }
    ],
    "total": 2
  }
}
```

## Error Codes
| Error Code | Description |
|------------|-------------|
| -1 | General error |
| -2 | Invalid parameters |
| -100 | No recordings found |
| -101 | Unauthorized access |

## Usage Context
This endpoint is typically used in the following scenarios:
- Displaying the video library timeline filtered by date
- Searching for specific events within a time range
- Filtering recordings by event type (motion, person detection, etc.)
- Retrieving only marked/favorite recordings

## Related Endpoints
- `getLibraryStatus` - Gets library status including available dates with recordings
- `getLibraryByTraceId` - Gets specific recording by trace ID
- `deleteRecord` - Deletes selected recordings
- `setMarkStatus` - Marks recordings as important/favorite
- `setReadStatus` - Marks recordings as read/viewed

## Implementation Details
The endpoint supports retrieving recordings from both cloud storage and local SD card storage, depending on the `isFromSDCard` parameter. The time-based filtering is the primary mechanism for narrowing down the list of recordings, with additional filters available for more specific queries.

When using pagination, the response includes the total count of matching recordings, allowing the client application to implement proper pagination controls. The application typically calls this endpoint with incremental date ranges to populate a timeline view of recordings in the library section of the app.