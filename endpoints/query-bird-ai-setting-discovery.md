# Query Bird AI Setting Endpoint Discovery

## Endpoint Information
- **Path:** `/birdLovers/queryBirdAiSetting`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Queries the bird detection AI settings for a specific camera device

## Request Parameters
The endpoint takes a request object with the following structure:

```json
{
  "serialNumber": "string",  // Required: Device serial number
  "reason": 0,               // Optional: Reason for the query (default: 0)
  "app": {                   // Standard BaseEntry fields
    "type": "string",        // Application type identifier
    "version": "string"      // Application version
  },
  "countryNo": "string",     // User's country code
  "language": "string",      // User's language preference 
  "tenantId": "string"       // User's tenant ID
}
```

## Response Format
The endpoint returns a response object with the following structure:

```json
{
  "result": 0,           // Result code (0 indicates success)
  "msg": "string",       // Response message
  "data": {
    "aiAnalyzeSwitch": boolean,  // Whether AI bird analysis is enabled (default: true)
    "aiNotifySwitch": boolean    // Whether AI bird notifications are enabled
  }
}
```

## Code Analysis
The endpoint is implemented in the API interfaces:

Example call signature from decompiled code:
```java
@POST("/birdLovers/queryBirdAiSetting")
Observable<BirdAiSettingResponse> queryBirdAiSetting(@Body BirdAiSettingEntry entry);
```

Implementation in the application:
```java
public void getBirdAiSetting(String serialNumber, final Callback<BirdAiSettingBean> callback) {
    BirdAiSettingEntry entry = new BirdAiSettingEntry();
    entry.setSerialNumber(serialNumber);
    
    apiClient.queryBirdAiSetting(entry)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<BirdAiSettingResponse>() {
            @Override
            public void onNext(BirdAiSettingResponse response) {
                if (response.getResult() < 0 || response.getData() == null) {
                    if (callback != null) {
                        callback.onError(response.getResult(), response.getMsg());
                    }
                    return;
                }
                
                BirdAiSettingBean bean = new BirdAiSettingBean();
                bean.setAiAnalyzeSwitch(response.getData().isAiAnalyzeSwitch());
                bean.setAiNotifySwitch(response.getData().isAiNotifySwitch());
                
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

## Settings Description
The endpoint returns two boolean settings:
1. `aiAnalyzeSwitch` - Controls whether the device's AI will analyze video footage for birds
2. `aiNotifySwitch` - Controls whether notifications will be sent when birds are detected

## Usage Context
This endpoint is used in the notification settings screens of the application. It allows users to configure whether their camera should:
1. Analyze video footage for birds (aiAnalyzeSwitch)
2. Send notifications when birds are detected (aiNotifySwitch)

The typical user flow:
1. User navigates to device settings → AI settings → Bird detection settings
2. App calls this endpoint to retrieve current bird detection settings
3. UI displays toggles for bird analysis and notifications
4. User can view and optionally modify these settings

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by negative result codes with error messages
- Network errors are handled through the RxJava Observable pattern

## Related Endpoints
- `/birdLovers/updateBirdAiSetting` - Updates the bird detection AI settings

## Notes
- This feature appears to be specific to camera devices that support bird detection
- The settings are device-specific and are identified by the device's serial number
- Bird detection is likely a specialized AI feature for outdoor cameras
- The application typically queries this endpoint when loading the notification settings page
- Default setting for analysis appears to be enabled (true)
- This is a specialized feature that may appeal to bird watchers or nature enthusiasts
- The UI components that control these settings are labeled "ai_bird_switch" and "ai_bird_notification"