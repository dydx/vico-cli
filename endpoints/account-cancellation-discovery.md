# Endpoint: /user/cancellation

## Purpose
This endpoint permanently deletes a user's account and all associated data. It is the final step in the account deletion process after the user has verified their identity and confirmed their intention to delete the account.

## Request Method
POST

## Request Parameters
The request body is a JSON object containing standard BaseEntry fields:

- `app` (AppBean object, optional): Application information
  - `appType` (String, optional): Type of application
  - `version` (String, optional): Application version
  - `timeZone` (String, optional): User's timezone
- `countryNo` (String, optional): User's country code
- `language` (String, optional): User's preferred language
- `tenantId` (String, optional): Tenant ID for the system

No additional parameters beyond the standard authentication token are required.

## Request Example
```json
{
  "app": {
    "appType": "Android",
    "version": "1.2.3",
    "timeZone": "America/New_York"
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
  "msg": "Error deleting account"
}
```

## Error Codes
- General error codes apply (network errors, authentication errors)
- No specific error codes were identified for this endpoint

## Usage in App
This endpoint is used as the final step in a multi-stage account deletion process:

1. User navigates to account settings and selects "Delete Account"
2. User is prompted to enter their password for verification
3. Password is verified using the `/account/passwordvalidation` endpoint
4. User is shown a confirmation screen with warnings about data loss
5. User must check a checkbox acknowledging they understand the consequences
6. User taps "Continue" to proceed with deletion
7. App unregisters the device from push notifications
8. App calls the `/user/cancellation` endpoint to delete the account
9. On success, user is redirected to the login screen and all activities are finished

## Client Implementation
The client creates a BaseEntry object and calls the deleteAccount method:

```java
// In AccountViewModel
public void deleteAccount() {
    BaseEntry baseEntry = new BaseEntry();
    // ... set any additional fields ...
    apiClient.deleteAccount(baseEntry)
        .subscribeOn(Schedulers.io())
        .observeOn(AndroidSchedulers.mainThread())
        .subscribe(new Observer<BaseResponse>() {
            @Override
            public void onNext(BaseResponse response) {
                if (response.getResult() >= 0) {
                    // Success handling: redirect to login
                    Intent intent = new Intent(context, LoginAndRegisterActivity.class);
                    intent.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
                    context.startActivity(intent);
                    Utils.finishActivityExcept(LoginAndRegisterActivity.class);
                } else {
                    // Error handling
                    Utils.showToastCenter(context, R.string.delete_account_failed);
                }
            }
            
            @Override
            public void onError(Throwable e) {
                // Network error handling
                Utils.showToastCenter(context, R.string.delete_account_failed);
            }
        });
}
```

## Error Handling
In the app, error responses are handled by:
- Checking if result < 0 to identify an error condition
- Showing a toast message "Delete account failed"
- Keeping the user on the current screen

## Security Considerations
- This endpoint performs a permanent, irreversible action
- Requires proper authentication token
- The app implements multiple confirmation steps and password verification before calling this endpoint
- All communication should be over HTTPS
- The UI provides warnings about the consequences of account deletion