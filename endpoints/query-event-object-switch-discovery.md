# Query Event Object Switch Endpoint Discovery

## Endpoint Information
- **Path:** `/aiAssist/queryEventObjectSwitch`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Queries AI event object detection settings for devices

## Request Parameters
The endpoint takes an `AiEventEntry` object in the request body:

```json
{
  "isAll": boolean,           // Whether to query all devices
  "serialNumbers": ["string"], // Array of device serial numbers to query (used if isAll is false)
  "app": {                    // Standard BaseEntry fields
    "type": "string",         // Application type identifier
    "version": "string"       // Application version
  },
  "countryNo": "string",      // User's country code
  "language": "string",       // User's language preference 
  "tenantId": "string"        // User's tenant ID
}
```

## Response Format
The endpoint returns an `AiEventResponse` object:

```json
{
  "code": 0,              // Response code (0 indicates success)
  "message": "string",    // Response message
  "data": [
    {
      "serialNumber": "string", // Device serial number
      "deviceName": "string",   // Name of the device
      "list": [
        {
          "eventObject": "string", // Type of event object (e.g., "person", "vehicle", "animal")
          "checked": boolean       // Whether detection for this object type is enabled
        },
        // Additional event objects...
      ]
    },
    // Additional devices...
  ]
}
```

## Code Analysis
The endpoint is implemented in multiple interfaces:
1. `com.ai.addxnet.vicohome_1742553098674_00O0o0oOO`
2. `com.ai.settingcore.vicohome_1742553098674_00O0o0oOO`
3. `com.ai.addxnet.vicohome_1742553098674_0oOOoOO`

Example call signature from decompiled code:
```java
@POST("/aiAssist/queryEventObjectSwitch")
Observable<AiEventResponse> queryEventObjectSwitch(@Body AiEventEntry aiEventEntry);
```

Implementation example:
```java
public void getEventObjectSwitch(List<String> serialNumbers, final Callback<List<DeviceAiEventBean>> callback) {
    AiEventEntry entry = new AiEventEntry();
    if (serialNumbers == null || serialNumbers.isEmpty()) {
        entry.setAll(true);
    } else {
        entry.setAll(false);
        entry.setSerialNumbers(serialNumbers);
    }
    
    apiClient.queryEventObjectSwitch(entry)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<AiEventResponse>() {
            @Override
            public void onNext(AiEventResponse response) {
                if (response.getCode() < 0 || response.getData() == null) {
                    if (callback != null) {
                        callback.onError(response.getCode(), response.getMessage());
                    }
                    return;
                }
                
                if (callback != null) {
                    callback.onSuccess(response.getData());
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
- "vehicle" - Cars and other vehicles
- "animal" - Animals (pets, wildlife)
- "package" - Package deliveries
- "face" - Face detection

The exact list may vary depending on the device model and capabilities.

## Usage Context
This endpoint is used in the following scenarios:
1. When loading the AI detection settings page
2. When viewing smart detection features for a device
3. Before configuring which object types should trigger alerts and recordings

The typical user flow:
1. User navigates to device settings → AI detection settings
2. App calls this endpoint to retrieve current settings
3. UI displays toggles for each object type the device can detect
4. User can view which detection types are enabled

## Related Endpoints
- `/aiAssist/updateEventObjectSwitch` - Updates AI event object detection settings

## Notes
- The isAll parameter provides flexibility to query either a specific set of devices or all devices
- Different device models may support different sets of detectable objects
- This endpoint focuses on what event types are detected, not notification settings
- The response contains both device information and its detection settings
- The checked boolean indicates whether that specific detection type is enabled
- To update these settings, the application uses the related endpoint `/aiAssist/updateEventObjectSwitch`