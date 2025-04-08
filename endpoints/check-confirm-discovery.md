# Endpoint: /account/checkconfirm

## Overview
This endpoint is used to verify verification codes (VCodes) sent to users during registration, password reset, or account verification processes. It validates whether the code provided by the user matches the one sent to their email or phone number.

## Request

### HTTP Method
POST

### URL
`/account/checkconfirm`

### Request Parameters
| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| code      | String | Yes      | The verification code sent to the user's email or phone |
| email     | String | Conditional | The user's email address (required if phone is not provided) |
| phone     | String | Conditional | The user's phone number (required if email is not provided) |

### Request Example
```json
{
  "code": "123456",
  "email": "user@example.com"
}
```

OR

```json
{
  "code": "123456",
  "phone": "+1234567890"
}
```

## Response

### Response Parameters
| Parameter | Type   | Description |
|-----------|--------|-------------|
| result    | Integer| Result code. Values >= 0 indicate success, values < 0 indicate failure |
| msg       | String | A message describing the result, especially useful in error cases |

### Success Response Example
```json
{
  "result": 0,
  "msg": "Success"
}
```

### Error Response Examples
```json
{
  "result": -1021,
  "msg": "Invalid verification code"
}
```

```json
{
  "result": -1036,
  "msg": "Verification code expired"
}
```

## Error Codes
| Error Code | Description |
|------------|-------------|
| -1021      | Invalid verification code or password |
| -1028      | Request too frequent |
| -1034      | Invalid verification code |
| -1035      | Request too frequent |
| -1036      | Verification code expired |

## Usage Example
This endpoint is typically called after a verification code has been sent to a user's email or phone number through another endpoint. The client app collects the code from the user and submits it to this endpoint for verification before allowing the user to complete their intended action (e.g., registration, password reset).

## Implementation Notes
- The endpoint returns a simple success response with `result >= 0` if the verification code is valid.
- If the verification code is invalid or expired, appropriate error codes are returned.
- The client typically handles these error responses by displaying appropriate error messages to the user.
- The authentication state is tracked via a token that is managed separately from this verification process.

## Related Endpoints
- `/account/registconfirm` - Sends verification code for registration
- `/account/resetconfirm` - Sends verification code for password reset
- `/account/checkbindcontact` - Similar endpoint for verifying codes when binding a phone or email to an account