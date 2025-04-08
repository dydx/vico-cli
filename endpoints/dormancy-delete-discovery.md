# Dormancy Delete Endpoint Discovery

## Endpoint Information
- **Path:** `/device/dormancy/delete`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Deletes an existing sleep/dormancy plan for a device

## Request Parameters
The endpoint takes a `DeletePlanCreateData` object in the request body:

```json
{
  "serialNumber": "string",  // Required: Device serial number
  "period": number,          // Required: The period identifier for the sleep plan to delete
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
The endpoint returns a standard `BaseResponse` object:

```json
{
  "result": 0,           // Result code (0 indicates success)
  "msg": "string"        // Response message
}
```

## Code Analysis
The endpoint is implemented in API interfaces:

Example call signature from decompiled code:
```java
@POST("/device/dormancy/delete")
Observable<BaseResponse> deleteSleepPlan(@Body DeletePlanCreateData deletePlanCreateData);
```

Implementation in the application:
```java
public void deleteSleepPlan(String serialNumber, int period, final Callback<BaseResponse> callback) {
    DeletePlanCreateData data = new DeletePlanCreateData();
    data.setSerialNumber(serialNumber);
    data.setPeriod(period);
    
    apiClient.deleteSleepPlan(data)
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
1. When a user deletes an existing sleep schedule for a device
2. In the device configuration settings screen
3. When removing unwanted or unnecessary sleep plans

The typical user flow:
1. User navigates to device settings → Sleep schedule settings
2. User sees list of existing sleep plans (retrieved via /device/dormancy/list)
3. User selects a plan to delete
4. User confirms deletion 
5. App calls this endpoint, passing the device serial number and plan period identifier
6. On success, the plan is removed from the list

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by non-zero result codes with error messages
- Network errors are handled through the RxJava Observable pattern

## Related Endpoints
- `/device/dormancy/list` - Lists all sleep plans for a device
- `/device/dormancy/create` - Creates a new sleep plan
- `/device/dormancy/edit` - Edits an existing sleep plan
- `/device/dormancy/switch` - Toggles the sleep mode feature on/off

## Notes
- This endpoint requires only the device serial number and the period identifier
- The period identifier must match an existing sleep plan for successful deletion
- The application typically refreshes the sleep plan list after deleting a plan
- Deleting a sleep plan is permanent and cannot be undone
- If the deleted plan was the only active plan, the device might exit sleep mode
- The application may need to handle the UI state update after successful deletion
- No further confirmation is required beyond the initial user confirmation