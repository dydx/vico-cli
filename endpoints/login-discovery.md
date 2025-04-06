# login Endpoint Documentation

## Overview
The `login` endpoint authenticates users and provides access tokens for the API. It allows users to log in using either email/password or phone/password combinations.

## API Details
- **Path**: `/account/login`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts an `AccountEntry` object with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| email | String | Yes (if phone not provided) | User's email address |
| password | String | Yes | User's password |
| phone | String | Yes (if email not provided) | User's phone number |
| code | String | No | Verification code (used for code-based login) |
| loginType | Integer | No | Login type (defaults to 0) |

### Example Request Body
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Or alternatively:

```json
{
  "phone": "1234567890",
  "password": "password123"
}
```

## Response Structure
The endpoint returns a `LoginResponse` object:

### Base Response Fields
| Field | Type | Description |
|-------|------|-------------|
| result | Integer | Status code (0 for success, negative values indicate errors) |
| msg | String | Status message or error description |

### LoginResponse Additional Fields
| Field | Type | Description |
|-------|------|-------------|
| data.id | Integer | User ID |
| data.name | String | User name |
| data.email | String | User email address |
| data.phone | String | User phone number |
| data.node | String | Node identifier |
| data.trackerToken | String | Tracking token |
| data.token.token | String | Authentication token |
| data.token.tokentype | String | Token type (usually "Bearer") |

### Example Success Response
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "id": 12345,
    "name": "User Name",
    "email": "user@example.com",
    "phone": "1234567890",
    "node": "node-identifier",
    "trackerToken": "tracking-token-string",
    "token": {
      "token": "authentication-token-string",
      "tokentype": "Bearer"
    }
  }
}
```

### Example Error Response
```json
{
  "result": -1001,
  "msg": "Invalid credentials"
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -1001 | Invalid credentials |
| -1002 | Account locked due to too many attempts |
| -1003 | Account not found |
| -1004 | Account requires verification |

## Usage in Application
The endpoint is called when the user attempts to log in through the login screen:
1. User enters email/phone and password in the login form
2. The application creates an `AccountEntry` object with the provided credentials
3. The request is made using RxJava through the `AccountApiClient`
4. On success, the authentication token is stored for future API calls
5. On failure, an appropriate error message is displayed to the user

## Constraints
- The user must provide either an email or phone number
- Password must meet the application's security requirements
- After multiple failed attempts, account may be temporarily locked
- For code-based login, a valid verification code must be provided