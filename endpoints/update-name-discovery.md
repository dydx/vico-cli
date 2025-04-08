# Endpoint: /user/updatename

## Purpose
This endpoint updates the user's display name/nickname in their profile. It allows users to change how their name appears in the application and to other users.

## Request Method
POST

## Request Parameters
The request body is a JSON object containing:

### NickNameEntry object
- `name` (String, required): The new display name for the user
- Plus standard BaseEntry fields:
  - `app` (AppBean object): Application information
  - `countryNo` (String, optional): Country code
  - `language` (String, optional): Language preference
  - `tenantId` (String, optional): Tenant ID for the system

## Request Example
```json
{
  "name": "John Smith",
  "app": {
    "type": "android",
    "version": "1.2.3"
  },
  "countryNo": "US",
  "language": "en"
}
```

## Response
The response is a standard BaseResponse object:

### Success Response (result = 0)
```json
{
  "result": 0,
  "msg": "Success"
}
```

### Error Response (result < 0)
```json
{
  "result": -1,
  "msg": "Error updating name"
}
```

## Error Codes
- General error codes apply (network errors, authentication errors)
- No specific error codes were identified for this endpoint

## Usage in App
This endpoint is used in the RenameNickNameActivity when a user edits their display name. The typical flow is:

1. User navigates to the profile screen and selects to edit their name
2. RenameNickNameActivity shows the current name in an editable field
3. User enters a new name and taps "Save"
4. App validates input (not empty and different from current name)
5. App calls this endpoint with the new name
6. On success, the local user profile is updated and the activity closes
7. On error, a toast message is shown

## Client Implementation
The client creates a NickNameEntry, sets the name field, and makes the API call:

```java
NickNameEntry nickNameEntry = new NickNameEntry();
nickNameEntry.setName(newName);
apiClient.setUserNiceName(nickNameEntry, new Subscriber<BaseResponse>() {
    @Override
    public void doOnNext(BaseResponse response) {
        if (response.getResult() < 0) {
            // Show error
            return;
        }
        // Update local user object
        UserBean user = UserManager.getInstance().getUser();
        user.setName(newName);
        UserManager.getInstance().updateUser(user);
        // Close activity
    }
    
    @Override
    public void doOnError(Throwable e) {
        // Show error
    }
});
```

## Error Handling
In the app, error responses are handled by:
- Checking if result < 0 to identify an error condition
- Showing appropriate error messages based on the error code
- Network errors are handled through the RxJava Observable pattern

## Security Considerations
- This endpoint handles user profile information
- Requires authentication token (handled at the HTTP client level)
- All communication should be over HTTPS