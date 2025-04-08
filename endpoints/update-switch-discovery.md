# Update Switch Endpoint Discovery

## Endpoint Information
- **Path:** `/usersetting/switch`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Updates the user's notification merge settings

## Request Parameters
The endpoint takes a `MegerPushEntry` object in the request body which extends `BaseEntry`:

```json
{
  "messageMergeSwitch": number,  // Required: 1 to enable message merging, 0 to disable
  "app": {                      // Standard BaseEntry fields
    "type": "string",           // Application type identifier
    "version": "string"         // Application version
  },
  "countryNo": "string",        // User's country code
  "language": "string",         // User's language preference 
  "tenantId": "string"          // User's tenant ID
}
```

## Response Format
The endpoint returns a standard `BaseResponse` object:

```json
{
  "result": 0,           // Result code (0 indicates success)
  "msg": "string"        // Response message
}
```

## Code Analysis
The endpoint is implemented in API interfaces:

Example call signature from decompiled code:
```java
@POST("/usersetting/switch")
Observable<BaseResponse> updateMargePushSwitch(@Body MegerPushEntry megerPushEntry);
```

Implementation in the application:
```java
public void updateMergePushData(int switchValue, final Callback<BaseResponse> callback) {
    MegerPushEntry entry = new MegerPushEntry();
    entry.setMessageMergeSwitch(switchValue);
    
    apiClient.updateMargePushSwitch(entry)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<BaseResponse>() {
            @Override
            public void onNext(BaseResponse response) {
                if (response.getResult() < 0) {
                    if (callback != null) {
                        callback.onError(response.getResult(), response.getMsg());
                    }
                    return;
                }
                
                if (callback != null) {
                    callback.onSuccess(response);
                }
            }
            
            @Override
            public void onError(Throwable e) {
                if (callback != null) {
                    callback.onError(-1, e.getMessage());
                }
            }
            
            // Other observer methods
        });
}
```

## Notification Merge Setting
The `messageMergeSwitch` parameter determines how notifications from devices are displayed:
- Value 1 (enabled): Notifications from multiple devices are consolidated/merged
- Value 0 (disabled): Each device's notifications are displayed individually

## Usage Context
This endpoint is used in the following scenarios:
1. When a user changes notification settings in the app
2. In the `NotificationViewModel` class to update how alerts are displayed
3. After retrieving current settings with `/usersetting/queryswitch`

The typical user flow:
1. User navigates to notification settings
2. App calls `/usersetting/queryswitch` to retrieve current merge setting
3. UI displays a toggle switch based on the current setting
4. User changes the setting by toggling the switch
5. App calls this endpoint to update the server with the new preference

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by negative result codes with error messages
- Network errors are handled through the RxJava Observable pattern

During the update process, the application typically:
1. Sets the UI state to LOADING
2. Makes the API call
3. Updates to SUCCESS or ERROR state based on the response
4. Displays appropriate feedback to the user

## Related Endpoints
- `/usersetting/queryswitch` - Retrieves the current notification merge setting

## Notes
- This endpoint affects how notifications are grouped and displayed to the user
- The setting is account-wide and applies to all devices
- Client-side, the setting is represented as an integer (1 for enabled, 0 for disabled)
- Server-side, the setting appears to be stored as a string ("1" or "0")
- Merging notifications can help reduce notification clutter for users with many devices
- The functionality is part of the broader notification management features in the application