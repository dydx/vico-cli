# Query Message Notification Endpoint Discovery

## Endpoint Information
- **Path:** `/device/queryMessageNotification/v1`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Retrieves message notification settings and configuration for a specific device

## Request Parameters
The endpoint takes a request object with the following structure:

```json
{
  "serialNumber": "string",      // Required: Device serial number
  "filterByAiAnalyze": boolean,  // Optional: Whether to filter results by AI analysis capability
  "userId": number,              // Optional: User ID
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
The endpoint returns a response object with the following structure:

```json
{
  "result": 0,                // Result code (0 indicates success)
  "msg": "string",            // Response message
  "data": {
    "list": [
      {
        "name": "string",     // Event type name (e.g., "person", "vehicle", "pet", "package", "other")
        "choice": boolean,    // Whether this event type is selected for notification
        "subEvent": [         // Optional: Sub-events for this event type
          {
            "name": "string", // Sub-event name
            "choice": boolean // Whether this sub-event is selected for notification
          }
        ]
      }
    ]
  }
}
```

## Common Event Types
Based on the code analysis, the following event types are supported:
- Person (`person`)
- Vehicle (`vehicle`)
- Pet (`pet`)
- Package (`package`)
- Other (`other`)

These event types may have sub-events that provide more granular notification control.

## Code Analysis
The endpoint is implemented in multiple API interfaces and is typically called using the `loadPersonDetectionDetailInfo` method name:

```java
@POST("/device/queryMessageNotification/v1")
Observable<PersonDetectionResponse> loadPersonDetectionDetailInfo(@Body PersonDetectEntry personDetectEntry);
```

Implementation example:
```java
public void getMessageNotificationConfig(String serialNumber, boolean filterByAiAnalyze, final Callback<List<PersonDetectionBean>> callback) {
    PersonDetectEntry entry = new PersonDetectEntry();
    entry.setSerialNumber(serialNumber);
    entry.setFilterByAiAnalyze(filterByAiAnalyze);
    
    apiClient.loadPersonDetectionDetailInfo(entry)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<PersonDetectionResponse>() {
            @Override
            public void onNext(PersonDetectionResponse response) {
                if (response.getResult() < 0 || response.getData() == null) {
                    if (callback != null) {
                        callback.onError(response.getResult(), response.getMsg());
                    }
                    return;
                }
                
                List<PersonDetectionBean> beans = convertResponseToBean(response.getData().getList());
                if (callback != null) {
                    callback.onSuccess(beans);
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
1. When loading notification settings for a device
2. In device settings screens to display notification options
3. When preparing to update notification preferences
4. Before configuring which event types should trigger notifications

The typical user flow:
1. User navigates to device settings → Notification settings
2. App calls this endpoint to retrieve current notification configuration
3. UI displays toggles for different event types (person, vehicle, etc.)
4. User can view which event types will trigger notifications

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by negative result codes with error messages
- Network errors are handled through the RxJava Observable pattern

## Related Endpoints
- `/device/updateMessageNotification/v1` - Updates the message notification settings for a device

## Notes
- This endpoint focuses on notification preferences for different AI detection events
- The `filterByAiAnalyze` parameter can be used to show only event types supported by the device's AI capabilities
- Different device models may support different sets of detectable objects
- The response provides both the main event types and their sub-events (if any)
- Each event type has a `choice` boolean indicating whether it's enabled for notification
- These settings determine which AI detection events will trigger notifications
- The endpoint uses RxJava observables for asynchronous processing
- Response data is typically converted to internal model objects for use in the application
- The v1 suffix suggests this may be the first version of this API