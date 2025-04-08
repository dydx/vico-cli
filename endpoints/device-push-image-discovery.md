# Device Push Image Endpoint Discovery

## Endpoint Information
- **Path:** `/device/devicePushImage`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Retrieves thumbnail images for user's devices

## Request Parameters
The endpoint takes a standard `BaseEntry` object in the request body without additional parameters:

```json
{
  "app": {
    "type": "string",     // Application type identifier
    "version": "string"   // Application version
  },
  "countryNo": "string",  // User's country code
  "language": "string",   // User's language preference
  "tenantId": "string"    // User's tenant ID
}
```

## Response Format
The endpoint returns a `DeviceThumbImageResponse` which extends the `BaseResponse` object:

```json
{
  "result": 0,            // Result code (0 indicates success)
  "msg": "string",        // Response message
  "data": [
    {
      "serialNumber": "string",     // Device serial number
      "lastPushImageUrl": "string", // URL to the device's thumbnail image
      "lastPushTime": 1612345678    // Timestamp of when the image was last updated (Unix time)
    },
    // Additional device entries...
  ]
}
```

## Code Analysis
The endpoint is implemented in multiple interfaces:
1. `com.ai.addxnet.vicohome_1742553098674_00O0o0oOO`
2. `com.ai.settingcore.vicohome_1742553098674_00O0o0oOO`

Example call signature from decompiled code:
```java
@POST("device/devicePushImage")
Observable<DeviceThumbImageResponse> deviceThumbImage(@Body BaseEntry baseEntry);
```

The endpoint is called by `DeviceSettingCore.getRecentCoverImage()` which makes the API call and processes the response.

Example implementation:
```java
public final void getRecentCoverImage(List<String> serialNumbers, vicohome_1742553098674_0O0oO00 callback) {
    com.ai.settingcore.vicohome_1742553098674_0O0o0oo.a().deviceThumbImage(new BaseEntry())
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<DeviceThumbImageResponse>() {
            @Override
            public void onSubscribe(Disposable d) {}

            @Override
            public void onNext(DeviceThumbImageResponse response) {
                if (response.getResult() < 0 || response.getData() == null) {
                    if (callback != null) {
                        callback.onError(response.getResult(), response.getMsg());
                    }
                    return;
                }
                
                List<DeviceCoverImageModel> filteredDevices = response.getData();
                if (serialNumbers != null && !serialNumbers.isEmpty()) {
                    // Filter devices by serial number
                    filteredDevices = new ArrayList<>();
                    for (DeviceCoverImageModel device : response.getData()) {
                        if (serialNumbers.contains(device.getSerialNumber())) {
                            filteredDevices.add(device);
                        }
                    }
                }
                
                if (callback != null) {
                    callback.onSuccess(filteredDevices);
                }
            }

            @Override
            public void onError(Throwable e) {
                if (callback != null) {
                    callback.onError(-1, e.getMessage());
                }
            }

            @Override
            public void onComplete() {}
        });
}
```

## Usage Context
This endpoint is used in the following scenarios:
1. Displaying device thumbnails in the camera list on the home screen
2. Refreshing device thumbnails periodically to show current camera views
3. Displaying device preview images in the device settings page

## Error Handling
The response extends BaseResponse which contains:
- `result`: Result code (0 indicates success, negative values indicate errors)
- `msg`: A message describing the result or error

The application implements custom error handling through the callback interface:
- If result < 0 or data is null, it calls onError with the error code and message
- Network errors are caught in the onError method of the Observer

## Filtering
The endpoint itself doesn't accept device filters in the request. However, the application provides client-side filtering:
- `DeviceSettingCore.getRecentCoverImage()` accepts an optional list of serial numbers
- When this list is provided, the app filters the response to only include the specified devices
- If no list is provided, all device thumbnails are returned

## Notes
- This endpoint provides thumbnails for all devices at once, which is efficient for populating the main camera view
- The thumbnail URLs are stored in a MutableLiveData object that the UI observes
- The thumbnails appear to be periodically updated on the server by the devices, and this endpoint retrieves the latest versions