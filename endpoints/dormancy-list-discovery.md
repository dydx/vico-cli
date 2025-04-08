# Dormancy List Endpoint Discovery

## Endpoint Information
- **Path:** `/device/dormancy/list`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Retrieves the list of sleep/dormancy plans for a device

## Request Parameters
The endpoint takes a `CommonSNRequest` object in the request body:

```json
{
  "serialNumber": "string",  // Required: Device serial number
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
The endpoint returns a `SleepPlanResponse` object which extends the `BaseResponse`:

```json
{
  "result": 0,            // Result code (0 indicates success)
  "msg": "string",        // Response message
  "data": [
    {
      "serialNumber": "string", // Device serial number
      "startHour": number,      // Hour when sleep mode starts (0-23)
      "startMinute": number,    // Minute when sleep mode starts (0-59)
      "endHour": number,        // Hour when sleep mode ends (0-23)
      "endMinute": number,      // Minute when sleep mode ends (0-59)
      "period": number,         // Period value (e.g., 1 for daily)
      "planDay": number,        // Bitmask of days for the plan (127 = all days)
      "planStartDay": [number]  // List of start days for the plan (1=Monday, 7=Sunday)
    },
    // Additional sleep plans...
  ]
}
```

## Plan Day Representation
The `planDay` field is a bitmask representing days of the week:
- 1 = Monday
- 2 = Tuesday
- 4 = Wednesday
- 8 = Thursday
- 16 = Friday
- 32 = Saturday
- 64 = Sunday
- 127 = All days (1+2+4+8+16+32+64)

## Code Analysis
The endpoint is implemented in multiple API interfaces:

Example call signature from decompiled code:
```java
@POST("/device/dormancy/list")
Observable<SleepPlanResponse> getSleepPlanList(@Body CommonSNRequest request);
```

Implementation in the application:
```java
public void getSleepPlanList(String serialNumber, final Callback<List<DeviceSleepPlanBean>> callback) {
    CommonSNRequest request = new CommonSNRequest();
    request.setSerialNumber(serialNumber);
    
    apiClient.getSleepPlanList(request)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<SleepPlanResponse>() {
            @Override
            public void onNext(SleepPlanResponse response) {
                if (response.getResult() < 0 || response.getData() == null) {
                    if (callback != null) {
                        callback.onError(response.getResult(), response.getMsg());
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

## Usage Context
This endpoint is used in the following scenarios:
1. When loading the device's dormancy/sleep settings page
2. Before creating or editing sleep plans to view existing ones
3. To check if a device has any sleep schedules configured

The typical user flow:
1. User navigates to device settings → Sleep schedule settings
2. App calls this endpoint to retrieve existing sleep plans
3. UI displays the list of sleep plans
4. User can view, add, edit, or delete sleep plans

## Sleep Plan Details
Each sleep plan contains:
- Start time (hour and minute)
- End time (hour and minute)
- Days of the week when the plan applies
- Device serial number

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by negative result codes with error messages
- Network errors are handled through the RxJava Observable pattern

## Related Endpoints
- `/device/dormancy/create` - Creates a new sleep plan
- `/device/dormancy/edit` - Edits an existing sleep plan
- `/device/dormancy/delete` - Deletes a sleep plan
- `/device/dormancy/switch` - Toggles the sleep mode feature on/off

## Notes
- Sleep/dormancy plans allow devices to enter a sleep mode during specified time periods
- This feature likely helps save power and reduce notifications during certain hours
- Multiple sleep plans can be configured for a device
- Plans can be set for specific days of the week
- The time format is in 24-hour format (0-23 hours)
- The response contains all sleep plans configured for the specified device