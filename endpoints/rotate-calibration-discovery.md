# Rotate Calibration Endpoint Discovery

## Endpoint Information
- **Path:** `/device/rotate-calibration`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Initiates or checks the status of device rotation calibration for PTZ (Pan-Tilt-Zoom) cameras

## Request Parameters
The endpoint takes a `RotateCalibrationRequest` object in the request body:

```json
{
  "serialNumber": "string",    // Required: Device serial number
  "needCalibration": boolean,  // Required: true to start calibration, false to check status
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
The endpoint returns a `RotateCalibrationResponse` which extends the `BaseResponse` object:

```json
{
  "result": 0,            // Result code (0 indicates success)
  "msg": "string",        // Response message
  "data": {
    "calibrationFinished": boolean  // Whether calibration has completed
  }
}
```

## Code Analysis
The endpoint is implemented in multiple interfaces and used by device configuration modules:

```java
@POST("device/rotate-calibration")
Observable<RotateCalibrationResponse> rotateCalibration(@Body RotateCalibrationRequest request);
```

The application uses this endpoint through the `DeviceSettingCore` class with these methods:
```java
// Check calibration status
public void getRotateCalibrationStatus(String serialNumber, RotateCalibrationCallback callback) {
    RotateCalibrationRequest request = new RotateCalibrationRequest();
    request.setSerialNumber(serialNumber);
    request.setNeedCalibration(false);
    
    apiClient.rotateCalibration(request)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<RotateCalibrationResponse>() {
            // Handle response
        });
}

// Start calibration
public void startRotateCalibration(String serialNumber, RotateCalibrationCallback callback) {
    RotateCalibrationRequest request = new RotateCalibrationRequest();
    request.setSerialNumber(serialNumber);
    request.setNeedCalibration(true);
    
    apiClient.rotateCalibration(request)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<RotateCalibrationResponse>() {
            // Handle response
        });
}
```

## Calibration Process Flow
The calibration process follows these steps:
1. User initiates calibration from device settings
2. App sends request with `needCalibration=true`
3. App starts polling with `needCalibration=false` every 5 seconds
4. The UI shows a loading state during calibration
5. When `calibrationFinished=true` is received, the process is complete
6. If timeout occurs (55 seconds), app calls the timeout endpoint

The polling implementation:
```java
private void monitorCalibrationStatus(String serialNumber, RotateCalibrationCallback callback) {
    calibrationSubscription = Observable.interval(5, 5, TimeUnit.SECONDS)
        .take(11) // Poll for up to 55 seconds
        .flatMap(tick -> {
            RotateCalibrationRequest request = new RotateCalibrationRequest();
            request.setSerialNumber(serialNumber);
            request.setNeedCalibration(false);
            return apiClient.rotateCalibration(request);
        })
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<RotateCalibrationResponse>() {
            // Process status updates
            // Check for completion
        });
}
```

## Related Endpoints
- **Path:** `/device/rotate-calibration/timeout`
- **Method:** POST
- **Description:** Reports a timeout in the calibration process
- **Request Body:** `CommonSNRequest` with just the device's serial number

```java
@POST("device/rotate-calibration/timeout")
Observable<BaseResponse> rotateCalibrationTimeout(@Body CommonSNRequest request);
```

## Usage Context
This endpoint is used in these scenarios:
1. When the user accesses pan-tilt settings in the device configuration
2. When a PTZ camera needs calibration after installation or reset
3. As part of the device setup process for calibration-capable devices

The UI flow typically involves:
1. User navigates to device settings → pan-tilt settings
2. User selects "Calibrate" option
3. Application shows loading UI and calls this endpoint
4. UI shows calibration progress until completion or timeout
5. User is notified of success or failure

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by non-zero result codes with error messages
- Timeouts are handled by the dedicated timeout endpoint
- Network errors trigger standard error callbacks

## Notes
- This endpoint is only applicable to PTZ (Pan-Tilt-Zoom) camera devices
- The calibration process runs on the camera hardware to properly align its movement mechanisms
- The process may take up to a minute to complete
- During calibration, the camera may physically move to test its rotation capabilities
- Calibration is needed to ensure accurate camera movement control
- The app tracks the calibration state locally to update the UI accordingly