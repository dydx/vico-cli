# Update Floodlight Luminance Endpoint Discovery

## Endpoint Information
- **Path:** `/device/updateFloodlightLuminance`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Updates the brightness level of a device's floodlight

## Request Parameters
The endpoint takes an `UpdateFloodlightLuminanceRequest` object in the request body which extends `BaseEntry`:

```json
{
  "serialNumber": "string",  // Required: Device serial number
  "luminance": integer,      // Required: Brightness level (1-100)
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
@POST("/device/updateFloodlightLuminance")
Observable<BaseResponse> updateFloodlightLuminance(@Body UpdateFloodlightLuminanceRequest request);
```

Implementation example from the application:
```java
public void updateFloodlightLuminance(String serialNumber, int luminance, final Callback callback) {
    UpdateFloodlightLuminanceRequest request = new UpdateFloodlightLuminanceRequest();
    request.setSerialNumber(serialNumber);
    request.setLuminance(luminance);
    
    apiClient.updateFloodlightLuminance(request)
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
1. `FloodLightView` - The compact view with basic brightness controls
2. `FloodLightFullView` - The expanded view with more detailed controls
3. Camera live view screens where floodlight control is available

The typical user flow:
1. User navigates to a camera view with floodlight capability
2. User interacts with a brightness slider control
3. When the user releases the slider, the app calls this endpoint
4. App shows a loading animation during the request
5. UI updates to reflect the new brightness after successful response
6. If an error occurs, the UI reverts to the previous brightness value

## UI Components
The brightness control is implemented as:
- A seekbar/slider component for adjusting values
- Numeric display showing the current brightness value
- Minimum brightness enforced at value 1
- Maximum brightness at 100

## Error Handling
The application implements error handling through:
- RxJava subscription error callbacks
- Loading state indicators during the request
- Reverting to previous brightness value on error
- Proper subscription lifecycle management to prevent memory leaks

## Related Endpoints
- `/device/updateFloodlightSwitch` - Toggles the floodlight on/off
- `/device/getDeviceAttributes` - Gets device attributes including floodlight capabilities and ranges
- `/device/modifyDeviceAttributes` - For other device settings modifications

## Notes
- This endpoint is only available for devices with floodlight capability
- The brightness level is an integer value, typically in the range 1-100
- The minimum brightness level is enforced at 1 (not 0)
- The application cancels any pending brightness adjustment requests before sending a new one
- Changes to brightness are made asynchronously on a background thread
- UI updates happen on the main thread for smooth user experience
- The floodlight brightness is an important feature for security cameras with built-in lighting
- Users can fine-tune the lighting level based on environmental conditions and preferences