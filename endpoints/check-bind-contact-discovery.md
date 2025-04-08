# Endpoint: /account/checkbindcontact

## Overview
This endpoint is used to verify verification codes (VCodes) specifically when binding or changing a phone number or email address to a user account. It validates whether the code provided by the user matches the one sent to their email or phone number during the contact binding process.

## Request

### HTTP Method
POST

### URL
`/account/checkbindcontact`

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
| -1011      | Phone number or email format is invalid |
| -1012      | Password or code is invalid |
| -1026      | Password/code error |
| -1027      | Verification code required |
| -1028      | Request too frequent |
| -1031      | Invalid verification code |
| -1032      | Invalid verification code |
| -1033      | Request frequency limit |
| -1034      | Invalid verification code |
| -1035      | Request too frequent |
| -1036      | Verification code expired |

## Usage Example
This endpoint is typically called after a verification code has been sent to a user's email or phone number through the `/account/sendbindcontactcode` endpoint. The client app collects the code from the user and submits it to this endpoint for verification. Upon successful verification, the phone number or email address is bound to the user's account.

**Typical flow:**
1. User initiates binding/changing a phone number or email address to their account
2. App calls `/account/sendbindcontactcode` to send a verification code
3. User receives and enters the code
4. App calls `/account/checkbindcontact` with the code and contact info
5. If successful, the contact method is bound to the user's account

## Implementation Notes
- The endpoint returns a success response with `result >= 0` if the verification code is valid.
- If the verification code is invalid or expired, appropriate error codes are returned.
- The client handles error responses by displaying appropriate error messages based on the error code.
- After successful verification, the user's profile is updated with the new contact information.
- For security reasons, the contact binding process requires verification to prevent unauthorized changes.

## Related Endpoints
- `/account/sendbindcontactcode` - Sends a verification code to the email/phone to be bound
- `/account/checkconfirm` - Similar endpoint used for general verification code validation
- `/account/getaccountinfo` - Gets account information including bound contact methods