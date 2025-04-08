# Update Event Object Switch Endpoint Discovery

## Endpoint Information
- **Path:** `/aiAssist/updateEventObjectSwitch`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Updates AI event detection settings for devices

## Request Parameters
The endpoint takes a `SwitchAiEventEntry` object in the request body:

```json
{
  "serialNumber": "string",      // Required: Device serial number to update
  "list": [                      // Required: List of event objects to update
    {
      "eventObject": "string",   // Type of event object (e.g., "person", "vehicle", "pet", "package")
      "checked": boolean         // Whether detection for this object type should be enabled
    }
  ],
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
  "result": 0,           // Response code (0 indicates success)
  "msg": "string"        // Response message
}
```

## Event Object Types
Common AI detectable event objects include:
- "person" - Human detection
- "vehicle" - Cars and other vehicles
- "pet" - Pet detection (animals)
- "package" - Package deliveries
- "face" - Face detection

The exact list may vary depending on the device model and capabilities.

## Code Analysis
The endpoint is implemented in multiple interfaces:
1. `com.ai.addxnet.vicohome_1742553098674_00O0o0oOO`
2. `com.ai.settingcore.vicohome_1742553098674_00O0o0oOO`

Example call signature from decompiled code:
```java
@POST("/aiAssist/updateEventObjectSwitch")
Observable<BaseResponse> updateEventObjectSwitch(@Body SwitchAiEventEntry switchAiEventEntry);
```

Implementation example:
```java
public void updateEventObjectSwitch(String serialNumber, List<AiEventBean> eventList, final Callback<BaseResponse> callback) {
    SwitchAiEventEntry entry = new SwitchAiEventEntry();
    entry.setSerialNumber(serialNumber);
    entry.setList(eventList);
    
    apiClient.updateEventObjectSwitch(entry)
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

## Usage Context
This endpoint is used in the following scenarios:
1. When updating AI detection settings for a specific device
2. When enabling or disabling specific object detection types
3. When configuring which objects should trigger alerts and recordings

The typical user flow:
1. User navigates to device settings → AI detection settings
2. App retrieves current settings using `/aiAssist/queryEventObjectSwitch`
3. User toggles on/off specific detection types (person, vehicle, etc.)
4. App calls this endpoint to update the settings on the server
5. Device applies the new detection settings

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by negative result codes with error messages
- Network errors are handled through the RxJava Observable pattern

## Related Endpoints
- `/aiAssist/queryEventObjectSwitch` - Queries current AI event object detection settings

## Notes
- This endpoint can update multiple event object types in a single request
- The request specifically targets a single device (by serial number)
- Changes take effect on the device after server processing
- The response does not include the updated configuration - to verify changes, a follow-up call to query the settings is needed
- Some detection capabilities may be dependent on specific device hardware and not all detection types may be supported by all devices
- These AI detection settings determine what types of objects can trigger events, notifications, and recordings
- Different device models may have different AI capabilities and supported detection types