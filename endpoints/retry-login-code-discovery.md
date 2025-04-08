# Endpoint: /account/retryLoginCode

## Overview
This endpoint is used to request a verification code when a user has entered the wrong password during login. It provides an alternative login method via verification code.

## Request

### HTTP Method
POST

### Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| loginType | int | Yes | Type of login: 1 for email, 2 for phone |
| email | string | Conditional | User's email address (required if loginType=1) |
| phone | string | Conditional | User's phone number (required if loginType=2) |

### Sample Request
```json
{
  "loginType": 1,
  "email": "user@example.com"
}
```

OR

```json
{
  "loginType": 2,
  "phone": "+1234567890"
}
```

## Response

### Success Response
```json
{
  "result": 0,
  "msg": "Success"
}
```

### Error Response
```json
{
  "result": -1034,
  "msg": "Invalid verification code"
}
```

### Common Error Codes
| Error Code | Description |
|------------|-------------|
| -1011 | Invalid email or phone format |
| -1021 | Invalid verification code |
| -1028 | Request too frequent |
| -1034 | Invalid verification code |
| -1035 | Request too frequent |
| -1036 | Verification code expired |

## Usage Notes
1. This endpoint is specifically used when a user enters an incorrect password during login
2. After a successful call, a verification code will be sent to the user's email or phone
3. The user can then use this code with the login endpoint to authenticate
4. The verification code has a limited validity period
5. There are rate limits on how frequently this API can be called for the same user

## Implementation Details
1. This endpoint is called from the `getCodeWhenPasswordWrong` method in the AccountViewModel
2. The response is handled in the `vicohome_1742553098674_0O0o0ooO` class which updates the UI state based on success or error
3. On success (result >= 0), the app proceeds to show a verification code input UI
4. On error, the app shows appropriate error messages based on the error code

## Error Handling
The app handles various error scenarios:
- Network connectivity issues
- Rate limiting errors
- Invalid input format errors
- Expired or invalid verification codes

The error messages are displayed to the user as toast notifications or text field errors depending on the error type.