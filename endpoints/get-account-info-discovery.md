# Endpoint: /user/getaccountinfo

## Purpose
This endpoint retrieves user account information including profile details, contact information (email/phone), and authentication token information. It's used to get the current user's account details after authentication.

## Request Method
POST

## Request Parameters
The request body is a JSON object containing standard BaseEntry fields:

- `app` (AppBean object, optional): Application information
  - `appType` (String, optional): Type of application (defaults to "Android")
  - `version` (Integer, optional): Application version number
  - `versionName` (String, optional): Application version name
  - `timeZone` (String, optional): User's timezone
- `countryNo` (String, optional): User's country code
- `language` (String, optional): User's preferred language
- `tenantId` (String, optional): Tenant ID for the system

## Request Example
```json
{
  "app": {
    "appType": "Android",
    "version": 123,
    "versionName": "1.2.3",
    "timeZone": "America/New_York"
  },
  "countryNo": "US",
  "language": "en"
}
```

## Response
The response is a JSON object containing:

### Success Response (result = 0)
```json
{
  "result": 0,
  "msg": "Success",
  "data": {
    "id": 123456,
    "name": "John Doe",
    "email": "user@example.com",
    "phone": null,
    "node": "us-east-1",
    "token": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "tokentype": "Bearer"
    },
    "trackerToken": "abc123def456"
  }
}
```

### Error Response (result < 0)
```json
{
  "result": -1001,
  "msg": "User not authenticated"
}
```

## Error Codes
- `-1001`: User not authenticated
- `-1002`: Session expired
- `-1003`: Invalid token
- `-1004`: User account not found

## Usage in App
This endpoint is typically called in the following scenarios:

1. After a successful login to retrieve the user's account information
2. When the user navigates to their profile screen
3. After account updates to refresh local user data
4. When refreshing the authentication token

In the app, the endpoint is called by the `UserProfileViewModel` which then updates the UI with the user's profile information such as name, email, and phone number.

Upon successful response, the app:
1. Updates the user interface with account details
2. Stores relevant user information locally
3. May refresh the authentication token if a new one is provided

## Error Handling
In the app, error responses are handled by:
- Checking if result < 0 to identify an error condition
- Showing appropriate error messages based on the error code
- Potentially triggering a re-login flow if authentication errors are encountered

## Security Considerations
- This endpoint handles sensitive user information
- Requires authentication token (handled at the HTTP client level)
- All communication should be over HTTPS
- The response contains authentication tokens which should be securely stored on the device