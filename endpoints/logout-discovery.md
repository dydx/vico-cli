# logout Endpoint Documentation

## Overview
The `logout` endpoint terminates a user's session and invalidates their authentication token. It provides a secure way for users to sign out of the application.

## API Details
- **Path**: `/account/logout/`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a standard `BaseEntry` object with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| app | AppBean | No | Application information |
| countryNo | String | No | Country code (e.g., "US") |
| language | String | No | Language code (e.g., "en") |
| tenantId | String | No | Tenant identifier |

No specific additional parameters are required beyond these common fields, which are typically automatically added by the API client.

### Example Request Body
```json
{
  "language": "en",
  "countryNo": "US"
}
```

## Response Structure
The endpoint returns a standard `BaseResponse` object:

### Base Response Fields
| Field | Type | Description |
|-------|------|-------------|
| result | int | Status code (0 for success, negative values indicate errors) |
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
  "msg": "Authentication failed"
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -1001 | Authentication failed |
| -2001 | Network error |

## Usage in Application
The endpoint is called when a user explicitly logs out:
1. User taps the logout button in the UI
2. A standard `BaseEntry` object is created
3. The request is made through RxJava using the account API client
4. On success, the application clears all local user data and authentication tokens
5. On failure, an error message may be displayed, but local logout actions are typically still performed

## Constraints
- The user must be currently logged in with a valid session
- After successful logout, any previously issued authentication tokens become invalid
- A new login is required to access authenticated endpoints after logout