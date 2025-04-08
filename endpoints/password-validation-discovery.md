# Endpoint: /account/passwordvalidation

## Purpose
This endpoint is used to verify a user's password without performing a full login action. It provides a way to confirm the current user's password is correct before proceeding with sensitive operations like changing account settings, deleting accounts, or updating credentials.

## Request Method
POST

## Request Parameters
The request body is a JSON object containing:

### VerifyEntry object
- `oldPassword` (String, required): The password to verify
- Plus standard BaseEntry fields:
  - `app` (AppBean object): Application information
  - `countryNo` (String, optional): Country code
  - `language` (String, optional): Language preference
  - `tenantId` (String, optional): Tenant ID for the system

## Request Example
```json
{
  "oldPassword": "user_password_here",
  "app": {
    // App information fields
  },
  "countryNo": "US",
  "language": "en"
}
```

## Response
The response is a BaseResponse object:

### Success Response (result >= 0)
```json
{
  "result": 0,
  "msg": "Success"
}
```

### Error Response (result < 0)
```json
{
  "result": -1021, // Example error code
  "msg": "Error message"
}
```

## Error Codes
- `-1021`: Invalid password

## Usage in App
This endpoint is primarily used in the app for:
1. Password verification before sensitive operations
2. Identity confirmation in security-critical functions

The request is made using the `verifyPassword` method in the `AccountViewModel` class, which creates a `VerifyEntry` object, populates it with the provided password, and sends it to the API. The response is processed to determine if the password is valid or not.

## Error Handling
In the app, error responses are handled by checking the response's result code:
- If result >= 0, the password is considered valid and the success state is set
- If result < 0, appropriate error messages are displayed based on the error code
- For invalid passwords, a specific error message is shown to the user

When a user enters an incorrect password, the error message "Password is incorrect" is displayed in the UI, typically in a text field error state.

## Security Considerations
- This endpoint handles sensitive user credentials
- Requires authentication token (handled at the HTTP client level)
- Should only be used over encrypted HTTPS connections