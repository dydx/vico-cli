# setReadStatus Endpoint Documentation

## Overview
The `setReadStatus` endpoint updates the read/unread status of recordings in the user's library. This allows users and the application to track which recordings have been viewed and which are still new or unread.

## API Details
- **Path**: `/library/updatemissing`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `SetReadStatueEntry` object that extends `BaseEntry` with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| missing | Integer | Yes | Status to set (0 = read, 1 = unread) |
| traceId | String | No* | Single trace ID of the recording to update |
| traceIds | String | No* | Multiple trace IDs as comma-separated string |
| app | Object | Yes | Application information (inherited from BaseEntry) |
| countryNo | String | Yes | Country code (e.g., "US") (inherited from BaseEntry) |
| language | String | Yes | Language code (e.g., "en") (inherited from BaseEntry) |
| tenantId | String | Yes | Tenant identifier (inherited from BaseEntry) |

*At least one of these two parameters must be provided.

### Example Request Body (Single traceId)
```json
{
  "missing": 0,
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

### Example Request Body (Multiple traceIds)
```json
{
  "missing": 0,
  "traceIds": "trace123abc,trace456def,trace789ghi",
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
The endpoint returns a `BaseResponse` object:

### BaseResponse Fields
| Field | Type | Description |
|-------|------|-------------|
| result | Integer | Status code (0 for success, negative values indicate errors) |
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
  "result": -1001,
  "msg": "Invalid trace ID"
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -1001 | Invalid trace ID |
| -1002 | No recordings to update |
| -2001 | Network error |
| -4001 | Authentication error |

## Usage in Application
The endpoint is called in these scenarios:
1. When a user views a recording for the first time (auto-marks as read)
2. When a user marks recordings as read or unread manually
3. When a user selects multiple recordings and marks them as read/unread
4. Through the LibraryCore.setReadStatus helper method which wraps the API call
5. To reset notification badges for new recordings

## Implementation Details
1. The application creates a SetReadStatueEntry object with the appropriate parameters
2. The missing parameter is set to 0 (read) or 1 (unread)
3. The traceId or traceIds parameter is set with the recording identifiers
4. The API call is made through the LibraryApiClient
5. Upon success, the UI is updated to reflect the new read status

## Related Endpoints
- `/library/updatemarked` (setMarkStatus) - Sets mark/favorite status for recordings
- `/library/newselectlibrary/newevent` (getEventRecordByFilter) - Gets recordings filtered by criteria

## Constraints
- User must be authenticated to access this endpoint
- The trace IDs must be valid and associated with the user's account
- When using traceIds for multiple recordings, they must be comma-separated
- Read status may affect notification badge counts and filters in the library view