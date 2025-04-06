# updateUserRollingDay Endpoint Documentation

## Overview
The `updateUserRollingDay` endpoint updates the billing cycle day for a user's subscription. This setting determines which day of the month the user will be billed for their recurring subscription.

## API Details
- **Path**: `/vip/user/rolling/day`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `RollingDayEntry` object that extends `BaseEntry` with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| rollingDay | Integer | Yes | The day of the month for the billing cycle (1-28) |
| userVipId | Integer | Yes | The ID of the user's VIP subscription |
| app | Object | Yes | Application information (inherited from BaseEntry) |
| countryNo | String | Yes | Country code (e.g., "US") (inherited from BaseEntry) |
| language | String | Yes | Language code (e.g., "en") (inherited from BaseEntry) |
| tenantId | String | Yes | Tenant identifier (inherited from BaseEntry) |

### Example Request Body
```json
{
  "rollingDay": 15,
  "userVipId": 12345,
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
  "msg": "Invalid subscription ID"
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -1001 | Invalid subscription ID |
| -1002 | Invalid rolling day value |
| -2001 | Network error |
| -4001 | Authentication error |

## Usage in Application
The endpoint is called in the subscription settings flow:
1. When a user navigates to their VIP subscription management screen (in VipServiceActivity)
2. User selects a new billing cycle day from the available options
3. The application creates a RollingDayEntry with the selected day and the user's subscription ID
4. Upon successful update, the UI is refreshed to show the new billing cycle day
5. This setting affects when the user will be charged for their recurring subscription

## Related Endpoints
- `/vip/user/service/info` (getVipUserServiceInfo) - Gets VIP user service information, including current rolling day

## Constraints
- User must be authenticated to access this endpoint
- The userVipId must be valid and belong to the authenticated user
- The rollingDay value is typically limited to 1-28 to avoid issues with months that have fewer than 31 days
- This setting typically only applies to monthly subscriptions, not annual ones
- Changing the billing cycle day may affect the prorated charges in the next billing period