# Update User Rolling Day API

## Endpoint Details
- **Path**: `/vip/user/rolling/day`
- **Method**: POST
- **Content-Type**: application/json
- **Description**: Updates the billing cycle day for a user's VIP subscription service

## Request Parameters

```json
{
  "rollingDay": 15,         // Day of month for billing cycle (1-28)
  "userVipId": 12345,       // User's VIP subscription ID
  "app": {                  // Standard app information
    "appName": "vicohome", 
    "appVersion": "1.2.3",
    "appBuild": "123",
    "channelId": 1
  },
  "countryNo": "US",        // Country code
  "language": "en",         // Language code
  "tenantId": "default"     // Tenant identifier
}
```

### Required Fields
- `rollingDay`: Integer between 1-28 representing the day of month for billing cycle
- `userVipId`: Integer representing the user's VIP subscription ID
- Standard app and user context fields (app, countryNo, language, tenantId)

## Response Format

```json
{
  "result": 0,              // Status code (0 for success, negative for errors)
  "msg": "success"          // Status message
}
```

### Success Response
- `result`: 0
- `msg`: "success" or other success message

### Error Responses
- `-1001`: Invalid subscription ID
- `-1002`: Invalid rolling day value
- `-2001`: Network error
- `-4001`: Authentication error

## Usage Notes
- The rolling day value is limited to 1-28 to avoid issues with months having fewer than 31 days
- The user must be authenticated to access this endpoint
- The userVipId must be valid and belong to the authenticated user
- This endpoint is typically used in subscription management flows when a user wants to change their billing cycle date