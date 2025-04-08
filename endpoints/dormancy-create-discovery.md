# Dormancy Create Endpoint Discovery

## Endpoint Information
- **Path:** `/device/dormancy/create`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Creates a new sleep/dormancy plan for a device

## Request Parameters
The endpoint takes a `DeviceSleepPlanBean` object in the request body:

```json
{
  "serialNumber": "string",        // Required: Device serial number
  "startHour": number,             // Required: Start hour of sleep period (0-23)
  "startMinute": number,           // Required: Start minute of sleep period (0-59)
  "endHour": number,               // Required: End hour of sleep period (0-23)
  "endMinute": number,             // Required: End minute of sleep period (0-59)
  "period": number,                // Required: Identifier for this sleep plan
  "planDay": number,               // Optional: Bitmask of days for the plan (127 = all days)
  "planStartDay": [number],        // Optional: List of start days for the plan (1=Monday, 7=Sunday)
  "app": {                         // Standard BaseEntry fields
    "type": "string",              // Application type identifier
    "version": "string"            // Application version
  },
  "countryNo": "string",           // User's country code
  "language": "string",            // User's language preference 
  "tenantId": "string"             // User's tenant ID
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

Alternatively, the `planStartDay` array can be used with day numbers (1-7).

## Code Analysis
The endpoint is implemented in API interfaces:

Example call signature from decompiled code:
```java
@POST("/device/dormancy/create")
Observable<BaseResponse> createSleepPlan(@Body DeviceSleepPlanBean sleepPlanBean);
```

Implementation in the application:
```java
public void createSleepPlan(DeviceSleepPlanBean sleepPlanBean, final Callback<BaseResponse> callback) {
    apiClient.createSleepPlan(sleepPlanBean)
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
1. When a user creates a new sleep schedule for a device
2. In the device settings page for power management
3. When setting up automated sleep times for battery conservation

The typical user flow:
1. User navigates to device settings → Sleep schedule settings
2. User selects to add a new sleep plan
3. User sets start time, end time, and days of the week
4. User submits the form, triggering this endpoint call
5. On success, the new sleep plan appears in the list

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by non-zero result codes with error messages
- Network errors are handled through the RxJava Observable pattern

## Related Endpoints
- `/device/dormancy/list` - Lists all sleep plans for a device
- `/device/dormancy/edit` - Edits an existing sleep plan
- `/device/dormancy/delete` - Deletes a sleep plan
- `/device/dormancy/switch` - Toggles the sleep mode feature on/off

## Notes
- Sleep/dormancy plans allow devices to enter a sleep mode during specified time periods
- This feature helps save battery power and reduce unnecessary recordings/alerts
- Multiple sleep plans can be configured for a device
- Plans can be set for specific days of the week
- The time format is in 24-hour format (0-23 hours)
- After creating a sleep plan, it's typically enabled automatically
- The period parameter appears to be a unique identifier for the sleep plan
- The app may need to refresh the sleep plan list after creating a new plan