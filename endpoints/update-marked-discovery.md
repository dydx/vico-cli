# `/library/updatemarked` Endpoint

## Overview
The `/library/updatemarked` endpoint is used to set the "marked" status of recordings in the video library. This allows users to mark or unmark recordings as favorites or important.

## Request Method
- **Method**: POST

## Request Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| marked | int | Yes | Status to set: 1 = mark, 0 = unmark |
| traceIds | String | Yes | Comma-separated list of recording trace IDs to update |

## Sample Request
```json
{
  "marked": 1,
  "traceIds": "trace123,trace456,trace789"
}
```

## Response Format
Returns a standard `BaseResponse`:

| Field | Type | Description |
|-------|------|-------------|
| result | int | 0 for success, error code otherwise |
| msg | String | Result message or error details |

## Sample Response
```json
{
  "result": 0,
  "msg": "Success"
}
```

## Error Handling
If an error occurs, the response will include a non-zero result code and an error message:

```json
{
  "result": 1001,
  "msg": "Operation failed"
}
```

## Implementation Notes
This endpoint is used in the recording library UI to allow users to mark important recordings. The app displays marked recordings with a star icon and allows filtering by marked status.

The implementation in the app:
1. Collects the trace IDs of selected recordings
2. Creates a SetMarkEntry object with:
   - marked = 0 or 1 (to unmark or mark)
   - traceIds = comma-separated list of recording trace IDs
3. Makes the API call
4. Updates the UI based on success/failure

## UI Implementation
In the LibraryEventPlayActivity, when a user selects recordings and taps the "mark" button, the app:
1. Collects the trace IDs of all selected recordings
2. Determines the marking action (mark or unmark) based on current states
3. Creates a SetMarkEntry with the appropriate parameters
4. Calls LibraryCore.setMarkStatus()
5. Shows a toast message "Marked successfully" or "Mark canceled" based on result
6. Updates the UI to show/hide star icons on the recordings