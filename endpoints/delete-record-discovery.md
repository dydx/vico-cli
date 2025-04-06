# deleteRecord Endpoint Documentation

## Overview
The `deleteRecord` endpoint deletes one or more video recordings from the user's library. This endpoint permanently removes the specified recordings based on their unique trace IDs.

## API Details
- **Path**: `/library/deletelibrary/`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `DeleteRecordEntry` object that extends `BaseEntry` with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| traceIdList | Array | Yes | List of trace IDs for the recordings to delete |
| app | Object | Yes | Application information (inherited from BaseEntry) |
| countryNo | String | Yes | Country code (e.g., "US") (inherited from BaseEntry) |
| language | String | Yes | Language code (e.g., "en") (inherited from BaseEntry) |
| tenantId | String | Yes | Tenant identifier (inherited from BaseEntry) |

### Example Request Body
```json
{
  "traceIdList": ["trace123abc", "trace456def"],
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
The endpoint returns a `LibraryDeleteRecordResponse` object that extends `BaseResponse`:

### BaseResponse Fields
| Field | Type | Description |
|-------|------|-------------|
| result | Integer | Status code (0 for success, negative values indicate errors) |
| msg | String | Status message or error description |

### LibraryDeleteRecordResponse Additional Fields
| Field | Type | Description |
|-------|------|-------------|
| data | Object | Information about the deletion operation |
| data.deletedCount | Integer | Number of recordings successfully deleted |
| data.failedCount | Integer | Number of recordings that failed to delete |
| data.failedTraceIds | Array | List of trace IDs that could not be deleted |

### Example Success Response
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "deletedCount": 2,
    "failedCount": 0,
    "failedTraceIds": []
  }
}
```

### Example Partial Success Response
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "deletedCount": 1,
    "failedCount": 1,
    "failedTraceIds": ["trace456def"]
  }
}
```

### Example Error Response
```json
{
  "result": -1001,
  "msg": "Invalid trace IDs"
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -1001 | Invalid trace IDs |
| -1002 | No permission to delete recordings |
| -2001 | Network error |
| -4001 | Authentication error |

## Usage in Application
The endpoint is called in these scenarios:
1. When a user selects one or more recordings in the library and taps the delete button
2. When a user swipes to delete a recording in the library view
3. When automated cleanup of old recordings is triggered
4. Through the LibraryViewModel's deleteRecord method, which is the entry point for the UI

The implementation relies on:
1. The DeleteRecordEntry class to format the request
2. The LibraryApiClient to make the actual API call
3. The LibraryCore class to handle proper callbacks and error handling

## Related Endpoints
- `/library/newselectlibrary/newevent` (getEventRecordByFilter) - Gets recordings filtered by criteria
- `/library/setmarkstatus` (setMarkStatus) - Marks recordings as favorites

## Constraints
- User must be authenticated to access this endpoint
- The trace IDs must be valid and associated with the user's account
- Deletion is permanent and cannot be undone
- Recordings that are marked as favorites may require an additional confirmation step in the UI before deletion