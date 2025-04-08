# Dormancy Edit Endpoint Discovery

## Endpoint Information
- **Path:** `/device/dormancy/edit`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Edits an existing sleep/dormancy plan for a device

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
@POST("/device/dormancy/edit")
Observable<BaseResponse> editSleepPlan(@Body DeviceSleepPlanBean sleepPlanBean);
```

Implementation in the application:
```java
public void editSleepPlan(DeviceSleepPlanBean sleepPlanBean, final Callback<BaseResponse> callback) {
    apiClient.editSleepPlan(sleepPlanBean)
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
1. When a user modifies an existing sleep schedule for a device
2. In the device configuration settings screen
3. When adjusting automated sleep times for battery conservation

The typical user flow:
1. User navigates to device settings → Sleep schedule settings
2. User sees list of existing sleep plans (retrieved via /device/dormancy/list)
3. User selects a plan to edit
4. User modifies start time, end time, and/or days of the week
5. User submits the form, triggering this endpoint call
6. On success, the modified sleep plan appears in the updated list

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by non-zero result codes with error messages
- Network errors are handled through the RxJava Observable pattern

## Related Endpoints
- `/device/dormancy/list` - Lists all sleep plans for a device
- `/device/dormancy/create` - Creates a new sleep plan
- `/device/dormancy/delete` - Deletes a sleep plan
- `/device/dormancy/switch` - Toggles the sleep mode feature on/off

## Notes
- This endpoint requires all the same fields as the create endpoint
- The period identifier must match an existing sleep plan to edit it properly
- The application typically refreshes the sleep plan list after editing a plan
- Sleep/dormancy plans allow devices to enter a low-power mode during specified time periods
- This feature helps save battery power and reduce unnecessary recordings/alerts
- The edited plan replaces the previous version completely, there is no partial update
- The time format is in 24-hour format (0-23 hours)
- Any changes to the plan take effect immediately if the sleep feature is enabled