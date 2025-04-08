# Dormancy Switch Endpoint Discovery

## Endpoint Information
- **Path:** `/device/dormancy/switch`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Toggles the sleep/dormancy mode on or off for a specific device

## Request Parameters
The endpoint takes a `SleepPlanData` object in the request body:

```json
{
  "serialNumber": "string",    // Required: Device serial number
  "dormancySwitch": number,    // Required: Sleep mode status (1 = on, 0 = off)
  "app": {                     // Standard BaseEntry fields
    "type": "string",          // Application type identifier
    "version": "string"        // Application version
  },
  "countryNo": "string",       // User's country code
  "language": "string",        // User's language preference 
  "tenantId": "string"         // User's tenant ID
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
@POST("/device/dormancy/switch")
Observable<BaseResponse> sleepSwitchSetting(@Body SleepPlanData sleepPlanData);
```

Implementation in the application:
```java
public Subscription setSleep(String serialNumber, boolean enabled, final Callback<BaseResponse> callback) {
    SleepPlanData sleepPlanData = new SleepPlanData(serialNumber, enabled ? 1 : 0);
    
    Subscription subscription = apiClient.sleepSwitchSetting(sleepPlanData)
        .subscribeOn(Schedulers.io())
        .observeOn(AndroidSchedulers.mainThread())
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
    
    subscriptions.add(subscription);
    return subscription;
}
```

## Usage Context
This endpoint is used in the following scenarios:
1. When a user toggles the sleep mode switch in the device settings
2. To manually enable or disable the dormancy feature
3. When sleep mode needs to be remotely controlled

The typical user flow:
1. User navigates to device settings → Sleep schedule settings
2. User toggles the main sleep mode switch
3. App calls this endpoint to update the device's sleep state
4. On success, the device enters or exits sleep mode

## Effects of Sleep Mode
When a device is in sleep mode:
- It stops recording and detecting motion
- The device reduces power consumption (especially important for battery-powered devices)
- Certain features like the status light may be turned off
- The device may not respond to some commands or events

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by non-zero result codes with error messages
- Network errors are handled through the RxJava Observable pattern

When an error occurs:
- The UI typically shows an error message
- The toggle switch is reverted to its previous state
- The application may retry the operation

## Related Endpoints
- `/device/dormancy/list` - Lists all sleep plans for a device
- `/device/dormancy/create` - Creates a new sleep plan
- `/device/dormancy/edit` - Edits an existing sleep plan
- `/device/dormancy/delete` - Deletes a sleep plan

## Notes
- This endpoint controls the master switch for the dormancy feature
- Even if sleep plans exist, they only take effect if dormancy is enabled via this endpoint
- The feature is primarily for power saving and privacy
- The dormancySwitch parameter uses 1 for enabled and 0 for disabled
- Sleep mode can be toggled independently from the sleep plans themselves
- The application maintains a subscription to track the state change
- The UI provides immediate feedback but reverts if the server request fails