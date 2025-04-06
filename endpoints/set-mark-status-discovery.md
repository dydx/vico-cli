# setMarkStatus Endpoint Documentation

## Overview
The `setMarkStatus` endpoint marks or unmarks recordings in the user's library as favorites. This allows users to flag important recordings for easy access and to prevent their automatic deletion.

## API Details
- **Path**: `/library/updatemarked`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `SetMarkEntry` object that extends `BaseEntry` with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| marked | Integer | Yes | Status to set (0 = unmarked, 1 = marked as favorite) |
| traceId | String | No* | Single trace ID of the recording to mark/unmark |
| traceIds | String | No* | Multiple trace IDs as comma-separated string |
| markIds | Array | No* | List of IDs to mark/unmark |
| app | Object | Yes | Application information (inherited from BaseEntry) |
| countryNo | String | Yes | Country code (e.g., "US") (inherited from BaseEntry) |
| language | String | Yes | Language code (e.g., "en") (inherited from BaseEntry) |
| tenantId | String | Yes | Tenant identifier (inherited from BaseEntry) |

*At least one of these three parameters must be provided.

### Example Request Body (Single traceId)
```json
{
  "marked": 1,
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
  "marked": 0,
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
1. When a user taps the favorite/star icon on a recording
2. When a user selects multiple recordings and marks them as favorites
3. When a user removes the favorite status from recordings
4. In the LibraryEventPlayActivity when viewing a specific recording
5. Through the LibraryCore.setMarkStatus helper method which wraps the API call

## Implementation Details
1. The application creates a SetMarkEntry object with the appropriate parameters
2. The marked parameter is set to 1 (favorite) or 0 (not favorite)
3. The traceId or traceIds parameter is set with the recording identifiers
4. The API call is made through the LibraryApiClient
5. Upon success, the UI is updated to reflect the new favorite status

## Related Endpoints
- `/library/setreadstatus` (setReadStatus) - Sets read status for recordings
- `/library/newselectlibrary/newevent` (getEventRecordByFilter) - Gets recordings filtered by criteria

## Constraints
- User must be authenticated to access this endpoint
- The trace IDs must be valid and associated with the user's account
- When using traceIds for multiple recordings, they must be comma-separated
- Marked recordings are typically excluded from automatic deletion policies