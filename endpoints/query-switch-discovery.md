# Query Switch Endpoint Discovery

## Endpoint Information
- **Path:** `/usersetting/queryswitch`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Retrieves the user's notification merge settings

## Request Parameters
The endpoint takes a standard `BaseEntry` object in the request body without additional parameters:

```json
{
  "app": {                  // Standard BaseEntry fields
    "type": "string",       // Application type identifier
    "version": "string"     // Application version
  },
  "countryNo": "string",    // User's country code
  "language": "string",     // User's language preference 
  "tenantId": "string"      // User's tenant ID
}
```

## Response Format
The endpoint returns a `MergePushResponse` object which extends the `BaseResponse`:

```json
{
  "result": 0,              // Result code (0 indicates success)
  "msg": "string",          // Response message
  "data": {
    "messageMergeSwitch": "string",  // "1" for enabled, "0" for disabled
    "userId": "string"               // The user's ID
  }
}
```

## Code Analysis
The endpoint is implemented in API interfaces:

Example call signature from decompiled code:
```java
@POST("/usersetting/queryswitch")
Observable<MergePushResponse> getMergePushData(@Body BaseEntry baseEntry);
```

Implementation in the application:
```java
public void getMergePushData(final Callback<MergePushBean> callback) {
    BaseEntry baseEntry = new BaseEntry();
    
    apiClient.getMergePushData(baseEntry)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<MergePushResponse>() {
            @Override
            public void onNext(MergePushResponse response) {
                if (response.getResult() < 0 || response.getData() == null) {
                    if (callback != null) {
                        callback.onError(response.getResult(), response.getMsg());
                    }
                    return;
                }
                
                MergePushBean bean = new MergePushBean();
                bean.setMessageMergeSwitch("1".equals(response.getData().getMessageMergeSwitch()));
                bean.setUserId(response.getData().getUserId());
                
                if (callback != null) {
                    callback.onSuccess(bean);
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
- Value "1" (enabled): Notifications from multiple devices are consolidated/merged
- Value "0" (disabled): Each device's notifications are displayed individually

## Usage Context
This endpoint is used in the following scenarios:
1. When loading notification settings in the app
2. In the `NotificationViewModel` class to configure how alerts are displayed
3. Before making changes to notification merge preferences

The typical user flow:
1. User navigates to notification settings
2. App calls this endpoint to retrieve current merge setting
3. UI displays a toggle switch based on the current setting
4. User can change the setting using the `/usersetting/switch` endpoint

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by negative result codes with error messages
- Network errors are handled through the RxJava Observable pattern

## Related Endpoints
- `/usersetting/switch` - Updates the notification merge setting

## Notes
- This endpoint requires no specific parameters beyond authentication details
- The setting affects the user experience by controlling how device notifications are grouped
- Client-side, the setting is represented as a boolean, but is stored as a string on the server ("1" or "0")
- This is a user-level setting, not a device-specific setting
- The application converts between string ("1"/"0") and boolean (true/false) representations
- Merging notifications can help reduce notification clutter for users with many devices