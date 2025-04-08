# Update Message Notification Endpoint Discovery

## Endpoint Information
- **Path:** `/device/updateMessageNotification/v1`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Updates message notification settings for a specific device

## Request Parameters
The endpoint takes a request object with the following structure:

```json
{
  "serialNumber": "string",      // Required: Device serial number
  "userId": number,              // Optional: User ID
  "eventObjectType": {           // Required: Object containing notification settings
    "person": ["string"],        // Person detection notification types to enable
    "pet": ["string"],           // Pet detection notification types to enable
    "vehicle": ["string"],       // Vehicle detection notification types to enable
    "other": ["string"],         // Other detection notification types to enable
    "package": ["string"]        // Package detection notification types to enable
  },
  "app": {                       // Standard BaseEntry fields
    "type": "string",            // Application type identifier
    "version": "string"          // Application version
  },
  "countryNo": "string",         // User's country code
  "language": "string",          // User's language preference 
  "tenantId": "string"           // User's tenant ID
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
The endpoint is implemented in multiple API interfaces and is typically called using the `updatePersonDetectionConfig` method name:

```java
@POST("/device/updateMessageNotification/v1")
Observable<BaseResponse> updatePersonDetectionConfig(@Body PersonDetectEntry.vicohome_1742553098674_00O0o0oOO entry);
```

Implementation example:
```java
public void updateMessageNotificationConfig(String serialNumber, Map<String, List<String>> eventObjectType, final Callback<BaseResponse> callback) {
    PersonDetectEntry.vicohome_1742553098674_00O0o0oOO entry = new PersonDetectEntry.vicohome_1742553098674_00O0o0oOO();
    entry.setSerialNumber(serialNumber);
    entry.setEventObjectType(eventObjectType);
    
    apiClient.updatePersonDetectionConfig(entry)
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

## Event Object Types
Common AI detectable event objects include:
- "person" - Human detection
- "pet" - Animal detection
- "vehicle" - Vehicle detection
- "package" - Package delivery detection
- "other" - Other detection types

Each category can have multiple sub-types that can be individually enabled or disabled.

## Usage Context
This endpoint is used in the following scenarios:
1. When updating notification settings for a device
2. After a user toggles different notification options in the settings UI
3. To customize which AI-detected events trigger push notifications

The typical user flow:
1. User navigates to device settings → Notification settings
2. App loads current settings using `/device/queryMessageNotification/v1`
3. User toggles different event notification types (person, vehicle, etc.)
4. User saves changes, which triggers this endpoint call
5. New notification settings are applied to the device

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by negative result codes with error messages
- Network errors are handled through the RxJava Observable pattern

## Related Endpoints
- `/device/queryMessageNotification/v1` - Retrieves current message notification settings

## Notes
- This endpoint allows fine-grained control over which AI detection events trigger notifications
- The eventObjectType parameter contains categories and their enabled notification types
- Empty arrays for a category indicate notifications are disabled for that type
- Non-empty arrays indicate which specific subtypes have notifications enabled
- The endpoint uses RxJava observables for asynchronous processing
- These settings are device-specific and identified by the device's serial number
- The v1 suffix suggests this may be the first version of this API
- After a successful update, the application typically refreshes its local configuration