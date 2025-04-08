# Update Floodlight Switch Endpoint Discovery

## Endpoint Information
- **Path:** `/device/updateFloodlightSwitch`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Updates the floodlight switch status (on/off) for a device with floodlight capability

## Request Parameters
The endpoint takes an `UpdateFloodlightSwitch` object in the request body which extends `BaseEntry`:

```json
{
  "serialNumber": "string",  // Required: Device serial number
  "switchOn": boolean,       // Required: True to turn on, false to turn off
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
  "result": 0,            // Result code (0 indicates success)
  "msg": "string"         // Response message
}
```

## Code Analysis
The endpoint is implemented in the API interfaces:

```java
@POST("/device/updateFloodlightSwitch")
Observable<BaseResponse> updateFloodlightSwitch(@Body UpdateFloodlightSwitch updateFloodlightSwitch);
```

Implementation example from the application:
```java
public void updateFloodlightSwitch(String serialNumber, boolean switchOn, final Callback callback) {
    UpdateFloodlightSwitch request = new UpdateFloodlightSwitch();
    request.setSerialNumber(serialNumber);
    request.setSwitchOn(switchOn);
    
    apiClient.updateFloodlightSwitch(request)
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
                    callback.onSuccess();
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
This endpoint is used in the following UI components:
1. `FloodLightView` - A dedicated UI control for toggling the floodlight
2. Device live view screens where floodlight control is available
3. Both normal view and full/split screen views of the camera feed

The typical user flow:
1. User navigates to a camera view with floodlight capability
2. The app displays a floodlight toggle icon/button
3. User taps the button to turn the floodlight on/off
4. App shows a loading animation during the request
5. UI updates to reflect the new state after successful response

## Error Handling
The application implements error handling through:
- RxJava subscription error callbacks
- Loading state indicators during the request
- Proper subscription lifecycle management to prevent memory leaks
- Null checks on parameters before making the API call

## Related Endpoints
- `/device/updateFloodlightLuminance` - Updates the brightness level of the floodlight
- `/device/getDeviceAttributes` - Gets device attributes including floodlight capabilities
- `/device/modifyDeviceAttributes` - For other device settings modifications

## Authorization
The UI shows different options based on the user's permission level:
- Admin users can see additional settings options
- Non-admin users have limited control over the floodlight

## Analytics Tracking
The application tracks floodlight usage with analytics events:
- Exposure events (EXP) when the floodlight control is displayed
- Click events (CLK) when the floodlight switch is toggled
- Tracking includes the current floodlight status and device context

## Notes
- This endpoint is only available for devices with floodlight capability
- The floodlight feature is typically found on security cameras with built-in lighting
- The application checks device capabilities before showing floodlight controls
- Floodlight controls are displayed in different UI layouts depending on the viewing mode
- Quick control of the floodlight is an important feature for security purposes