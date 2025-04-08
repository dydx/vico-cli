# Endpoint: /user/updatepassword

## Purpose
This endpoint allows users to change their account password. It requires both the old password and the new password to be provided, and it performs validation to ensure the new password meets security requirements.

## Request Method
POST

## Request Parameters
The request body is a JSON object containing:

### PasswordUpdateEntry object
- `oldPassword` (String, required): The user's current password
- `newPassword` (String, required): The new password the user wants to set
- Plus standard BaseEntry fields:
  - `app` (AppBean object): Application information
  - `countryNo` (String, optional): Country code
  - `language` (String, optional): Language preference
  - `tenantId` (String, optional): Tenant ID for the system

## Request Example
```json
{
  "oldPassword": "CurrentPassword123",
  "newPassword": "NewPassword456",
  "app": {
    "type": "android",
    "version": "1.2.3"
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
    "node": "us-east-1",
    "token": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "tokentype": "Bearer"
    }
  }
}
```

### Error Response (result < 0)
```json
{
  "result": -1021,
  "msg": "Incorrect old password"
}
```

## Error Codes
- `-1021`: Incorrect old password
- `-99999`: General error

## Password Requirements
The new password must match the following requirements:
- At least 8 characters long
- Contains at least one lowercase letter
- Contains at least one uppercase letter
- Regular expression pattern: `^(?=.*[a-z])(?=.*[A-Z]).{8,}$`

## Usage in App
This endpoint is used in the ChangePasswordActivity when a user wants to update their password. The typical flow is:

1. User navigates to security settings and selects to change password
2. User enters their current password, new password, and confirms the new password
3. Client-side validates that the new password and confirmation match
4. Client-side validates that the new password meets the pattern requirements
5. App calls this endpoint with the old and new passwords
6. On success, the app updates the local user data with the new authentication token
7. On error, appropriate error messages are displayed

## Client Implementation
The client creates a PasswordUpdateEntry with the old and new passwords and makes the API call:

```java
PasswordUpdateEntry entry = new PasswordUpdateEntry();
entry.setOldPassword(oldPassword);
entry.setNewPassword(newPassword);
apiClient.updatePassword(entry, new Subscriber<PasswordUpdateResponse>() {
    @Override
    public void doOnNext(PasswordUpdateResponse response) {
        if (response.getResult() < 0) {
            // Show specific error based on result code
            if (response.getResult() == -1021) {
                // Show incorrect old password error
            } else {
                // Show general error
            }
            return;
        }
        // Update local user data with new token
        UserBean user = response.getData();
        UserManager.getInstance().updateUser(user);
        // Show success message and close activity
    }
    
    @Override
    public void doOnError(Throwable e) {
        // Show network error
    }
});
```

## Error Handling
In the app, error responses are handled by:
- Checking result code and showing specific error messages
- Common errors include incorrect old password (-1021)
- Network errors are handled separately

## Security Considerations
- This endpoint handles user credentials and should only be used over HTTPS
- Requires authentication token
- Password strength requirements are enforced both on client and server side
- After successful password change, a new authentication token is issued