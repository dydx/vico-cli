# Update Bird AI Setting Endpoint Discovery

## Endpoint Information
- **Path:** `/birdLovers/updateBirdAiSetting`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Updates bird detection AI settings for a specific camera device

## Request Parameters
The endpoint takes a request object with the following structure:

```json
{
  "serialNumber": "string",     // Required: Device serial number
  "aiAnalyzeSwitch": boolean,   // Optional: Whether AI bird analysis is enabled
  "aiNotifySwitch": boolean,    // Optional: Whether bird detection notifications are enabled
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
The endpoint is implemented in the API interfaces:

Example call signature from decompiled code:
```java
@POST("/birdLovers/updateBirdAiSetting")
Observable<BaseResponse> updateBirdAiSetting(@Body UpdateBirdAiSettingEntry entry);
```

Implementation in the application:
```java
public void updateBirdAiSetting(String serialNumber, boolean aiAnalyzeSwitch, boolean aiNotifySwitch, final Callback<BaseResponse> callback) {
    UpdateBirdAiSettingEntry entry = new UpdateBirdAiSettingEntry();
    entry.setSerialNumber(serialNumber);
    entry.setAiAnalyzeSwitch(aiAnalyzeSwitch);
    entry.setAiNotifySwitch(aiNotifySwitch);
    
    apiClient.updateBirdAiSetting(entry)
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

## Settings Description
The endpoint allows updating two boolean settings:
1. `aiAnalyzeSwitch` - Controls whether the device's AI will analyze video footage for birds
2. `aiNotifySwitch` - Controls whether notifications will be sent when birds are detected

## Usage Context
This endpoint is used in the notification settings screens of the application. It allows users to configure whether their camera should:
1. Analyze video footage for birds (aiAnalyzeSwitch)
2. Send notifications when birds are detected (aiNotifySwitch)

The typical user flow:
1. User navigates to device settings → AI settings → Bird detection settings
2. App loads current settings using `/birdLovers/queryBirdAiSetting`
3. User toggles bird analysis or notification settings
4. App calls this endpoint to save the new settings
5. Settings are applied to the device

## UI Integration
In the Android application, this endpoint is used in connection with:
- A toggle switch for enabling/disabling bird AI analysis (`ai_bird_switch`)
- A toggle for enabling/disabling notifications for bird detection (`ai_bird_notification`)

The UI shows different states:
- A "normal state" view for users who haven't subscribed to the feature
- A "paid state" view for users who have access to the bird detection feature

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by negative result codes with error messages
- Network errors are handled through the RxJava Observable pattern

Common errors include:
- Invalid or missing serialNumber
- Unauthorized access (if feature requires subscription)
- Server-side processing errors

## Related Endpoints
- `/birdLovers/queryBirdAiSetting` - Gets the current bird AI settings for a device
- `/birdLovers/searchBirdName` - Searches for bird names (likely for identified birds)
- `/birdLovers/queryVideoPossibleSubcategory` - Gets potential bird types for video analysis
- `/birdLovers/feedbackBirdName` - Provides feedback about bird identification

## Notes
- This feature appears to be part of a paid/premium feature that may require subscription
- Bird detection is likely a specialized AI feature for outdoor cameras
- The settings are device-specific and are identified by the device's serial number
- Both analysis and notification settings can be updated in a single request
- If either setting is omitted from the request, the current value is maintained
- This is a specialized feature that likely appeals to bird watchers or nature enthusiasts