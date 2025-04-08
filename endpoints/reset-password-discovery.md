# Reset Password Endpoint Discovery

## Endpoint Information
- **URL**: `/account/resetpswd`
- **Method**: POST
- **Description**: Resets a user's password using a verification code

## Request Parameters

The endpoint accepts a JSON object with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| code | string | Yes | Verification code previously sent to the user's email or phone |
| password | string | Yes | New password (must match the pattern: at least 8 characters, containing at least one uppercase and one lowercase letter) |
| email | string | Conditional* | User's email address |
| phone | string | Conditional* | User's phone number |

*Either email or phone must be provided, depending on which contact method was used to send the verification code.

## Request Example

```json
{
  "email": "user@example.com",
  "code": "123456",
  "password": "NewPassword123"
}
```

Or with phone:

```json
{
  "phone": "+15551234567",
  "code": "123456", 
  "password": "NewPassword123"
}
```

## Response

### Success Response
- **Status Code**: 200 OK
- **Content Type**: application/json

```json
{
  "code": 0,
  "message": "Success"
}
```

### Error Responses

#### Invalid Request
- **Status Code**: 400 Bad Request
- **Content Type**: application/json

```json
{
  "code": 400,
  "message": "Invalid request parameters"
}
```

#### Invalid Verification Code
- **Status Code**: 400 Bad Request
- **Content Type**: application/json

```json
{
  "code": -1,
  "message": "Invalid verification code"
}
```

#### Code Expired
- **Status Code**: 410 Gone
- **Content Type**: application/json

```json
{
  "code": -2,
  "message": "Verification code has expired"
}
```

#### Password Requirements Not Met
- **Status Code**: 400 Bad Request
- **Content Type**: application/json

```json
{
  "code": -3,
  "message": "Password does not meet requirements"
}
```

#### Server Error
- **Status Code**: 500 Internal Server Error
- **Content Type**: application/json

```json
{
  "code": 500,
  "message": "Server error"
}
```

## Notes
- Before calling this endpoint, the client should request a verification code using the `/account/resetconfirm` endpoint
- Password must match the pattern: ^(?=.*[a-z])(?=.*[A-Z]).{8,}$ (at least 8 characters with at least one uppercase and one lowercase letter)
- After successful password reset, the user is redirected to the login screen

## Implementation Details
- The client constructs a `ResetPasswordEntry` object with the verification code, new password, and either email or phone
- The request is sent to `/account/resetpswd` 
- On success, the user is notified and redirected to the login screen
- On failure, an appropriate error message is displayed

## Related Endpoints
- `/account/resetconfirm` - Used to request a verification code for password reset
- `/account/login` - Used to authenticate with the new password after reset