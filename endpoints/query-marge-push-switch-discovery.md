# queryMargePushSwitch Endpoint

## Overview
The queryMargePushSwitch endpoint retrieves the user's push notification merging preferences. This setting determines whether the system should combine multiple notifications into a single notification (merged) or send individual notifications for each event (non-merged).

> Note: The name "marge" appears to be a typo for "merge" based on the implementation details and response structure.

## API Details
- **Path**: `/usersetting/queryswitch`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Retrieves the user's push notification merging preferences.

## Request Parameters
The request body only requires standard BaseEntry properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| app | Object | Yes | Application information |
| countryNo | String | Yes | Country code (e.g., "US") |
| language | String | Yes | Language code (e.g., "en") |
| tenantId | String | Yes | Tenant identifier |

## Request Example
```json
{
  "app": "vicohome",
  "countryNo": "US",
  "language": "en",
  "tenantId": "AAA12345678"
}
```

## Response Structure
The response contains information about the user's push notification merging preference:

| Property | Type | Description |
|----------|------|-------------|
| result | Integer | Status code (0 for success, negative values for errors) |
| msg | String | Status message or error description |
| data | Object | Contains push notification settings |

### Data Object Structure
| Property | Type | Description |
|----------|------|-------------|
| messageMergeSwitch | String | Notification merging setting ("0" = not merged, "1" = merged) |
| userId | String | User identifier |

## Response Example
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "messageMergeSwitch": "1",
    "userId": "user123456"
  }
}
```

## Error Codes
| Error Code | Description |
|------------|-------------|
| -1 | General error |
| -2 | Authentication error |
| -5 | Network error |

## Usage Context
This endpoint is typically used in the following scenarios:
- When initializing notification settings UI components
- Before displaying notification preferences to the user
- When determining how to present notifications to users
- During application startup to configure notification handling

## Related Endpoints
- `updateMargePushSwitch` - Updates the push notification merging preference

## Implementation Notes
The endpoint retrieves a user-level setting that controls how push notifications are handled across all devices associated with the user's account. When the messageMergeSwitch value is "1", multiple notifications that occur close together will be combined into a single notification. When the value is "0", each event generates a separate notification.

This is a global user preference rather than a device-specific setting, which is why the endpoint doesn't require a device serial number. The endpoint is called by the DeviceAICore.getMergePushData() method in the application code.

Note that despite the name "marge" in the endpoint name (likely a typo), the actual response structure uses "merge" terminology (MergePushResponse) which accurately describes the functionality.