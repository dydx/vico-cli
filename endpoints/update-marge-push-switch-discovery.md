# updateMargePushSwitch Endpoint Documentation

## Overview
The `updateMargePushSwitch` endpoint updates the user's push notification merging preferences. This setting determines whether the system should combine multiple notifications into a single notification (merged) or send individual notifications for each event (non-merged).

## API Details
- **Path**: `/usersetting/switch`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `MegerPushEntry` object that extends `BaseEntry` with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| messageMergeSwitch | Integer | Yes | Notification merging preference (1 = merged, 0 = not merged) |
| app | Object | Yes | Application information (inherited from BaseEntry) |
| countryNo | String | Yes | Country code (e.g., "US") (inherited from BaseEntry) |
| language | String | Yes | Language code (e.g., "en") (inherited from BaseEntry) |
| tenantId | String | Yes | Tenant identifier (inherited from BaseEntry) |

### Example Request Body
```json
{
  "messageMergeSwitch": 1,
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
  "result": -2001,
  "msg": "Network error"
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -2001 | Network error |
| -3001 | Invalid parameters |
| -4001 | Authentication error |

## Usage in Application
The endpoint is called when the user updates their notification settings:
1. User toggles the notification merging setting in the notification preferences UI
2. The application creates a `MegerPushEntry` object with the selected preference
3. The request is made through the `NotificationViewModel.updateMegerPushSwitch()` method
4. On success, the updated preference is reflected in the UI
5. On failure, an appropriate error message is displayed to the user

## Related Endpoints
- `/usersetting/queryswitch` (queryMargePushSwitch) - Gets the current push notification merging preference

## Constraints
- The user must be authenticated to update preferences
- The messageMergeSwitch parameter must be either 0 or 1