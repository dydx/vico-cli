# Endpoint: /account/register/

## Purpose
This endpoint registers a new user account in the VicoHome system. It requires a verification code that must be obtained through the `/account/registconfirm` endpoint first. The request accepts either email or phone-based registration with a password and nickname.

## Request Method
POST

## Request Parameters
The request body is a JSON object containing:

### RegisterEntry object
- `code` (String, required): Verification code received via email or SMS
- `email` (String, conditional): Email address (required if registering with email)
- `phone` (String, conditional): Phone number (required if registering with phone)
- `password` (String, required): Account password
- `name` (String, required): User's display name/nickname
- `supportFreeLicense` (Boolean, optional): Whether to enable free trial features
- Plus standard BaseEntry fields:
  - `app` (AppBean object): Application information
  - `countryNo` (String, optional): Country code
  - `language` (String, optional): Language preference
  - `tenantId` (String, optional): Tenant ID for the system

## Request Examples

### Email Registration
```json
{
  "code": "123456",
  "email": "user@example.com",
  "password": "securePassword123",
  "name": "John Doe",
  "supportFreeLicense": true,
  "app": {
    "type": "android",
    "version": "1.2.3"
  },
  "countryNo": "US",
  "language": "en"
}
```

### Phone Registration
```json
{
  "code": "123456",
  "phone": "+15551234567",
  "password": "securePassword123",
  "name": "John Doe",
  "supportFreeLicense": true,
  "app": {
    "type": "ios",
    "version": "1.2.3"
  },
  "countryNo": "US",
  "language": "en"
}
```

## Response

### Success Response (code = 0)
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "id": 123456,
    "name": "John Doe",
    "email": "user@example.com",
    "phone": "+15551234567",
    "node": "us-east-1",
    "token": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "tokentype": "Bearer"
    }
  }
}
```

### Error Response (code ≠ 0)
```json
{
  "code": 400,
  "message": "Invalid verification code"
}
```

## Error Codes
- `400`: Invalid request parameters
- `401`: Invalid verification code
- `409`: Email/phone already registered
- `410`: Verification code expired
- `500`: Server error

## Registration Process Flow
1. User enters email/phone and password
2. App calls `/account/registconfirm` to send verification code
3. User enters verification code
4. User enters desired display name
5. App calls `/account/register/` with all information
6. On success, user is logged in automatically

## Usage in App
This endpoint is used during the user registration process. The app first collects the user's email or phone, then requests a verification code through the `/account/registconfirm` endpoint. Once the user enters the verification code and completes the registration form, this endpoint is called to create the account.

Upon successful registration, the app:
1. Stores the authentication token for subsequent API calls
2. Automatically logs the user in
3. Navigates to the main application interface

## Error Handling
In the app, error responses are handled by checking the response's code:
- If code = 0, the registration was successful
- Otherwise, appropriate error messages are displayed based on the code and message

## Security Considerations
- This endpoint handles sensitive user information
- Passwords should meet security requirements (minimum length, complexity)
- The verification code has a limited validity period
- All communication should be over HTTPS