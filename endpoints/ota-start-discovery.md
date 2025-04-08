# OTA Start Endpoint Discovery

## Endpoint Information
- **Path:** `/device/otastart/`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Initiates the Over-The-Air (OTA) firmware update process for a specified device

## Request Parameters
The endpoint takes a `SerialNoEntry` object in the request body:

```json
{
  "serialNumber": "string",  // Required: The device's serial number
  "app": {
    "type": "string",        // Application type identifier
    "version": "string"      // Application version
  },
  "countryNo": "string",     // User's country code
  "language": "string",      // User's language preference
  "tenantId": "string"       // User's tenant ID
}
```

## Response Format
The endpoint returns a `GetOtaStatueResponse` object:

```json
{
  "result": 0,               // Result code (0 for success, negative for errors)
  "msg": "string",           // Response message
  "data": {
    "serialNumber": "string", // Device serial number
    "inProgress": 0,          // Flag indicating if update is in progress (0 or 1)
    "status": 0,              // General status code
    "otaStatus": 0,           // OTA specific status code (0-5)
    "transferredSize": 0,     // Bytes transferred so far
    "totalSize": 0,           // Total firmware size in bytes 
    "targetFirmware": "string", // Target firmware version
    "localDataProgress": 0,   // Local data transfer progress
    "localDataStartTime": 0   // Timestamp when local data transfer started
  }
}
```

## OTA Status Codes
Based on the code analysis, the following status codes are used:
- `0`: Idle/Waiting
- `1`: Transferring (In Progress)
- `2`: Installing
- `3`: Completed (Success)
- `4`: Failed
- `5`: Network Error

## Error Codes
Some identified error codes:
- `-10401`: Timeout error
- `-10402`: Network transfer failure
- `-2131`: General error
- `-3111`: Low power error
- `-2133`: Device in sleep mode

## Code Analysis
The endpoint is implemented in at least these interfaces:
1. `com.ai.addxnet.vicohome_1742553098674_00O0o0oOO`
2. `com.ai.settingcore.vicohome_1742553098674_00O0o0oOO`

Example call signature from decompiled code:
```java
@POST("device/otastart/")
Observable<GetOtaStatueResponse> otaStart(@Body SerialNoEntry serialNoEntry);
```

The application implements the OTA start process through the `DeviceOTACore` class:
```java
public void otaStart(String serialNumber, OtaProgressCallback callback) {
    // Create request
    SerialNoEntry entry = new SerialNoEntry();
    entry.setSerialNumber(serialNumber);
    
    // Stop any previous OTA process
    stopOta();
    
    // Make API call
    apiClient.otaStart(entry)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<GetOtaStatueResponse>() {
            @Override
            public void onNext(GetOtaStatueResponse response) {
                if (response.getResult() < 0 || response.getData() == null) {
                    callback.onError(response.getResult(), response.getMsg());
                    return;
                }
                
                // Start polling for status
                startOtaStatusPolling(serialNumber, callback);
            }
            
            @Override
            public void onError(Throwable e) {
                callback.onError(-10402, e.getMessage());
            }
            
            // Other observer methods
        });
}
```

## Progress Monitoring
The OTA process progress is monitored through the `onProcess` callback:
- For status 0: 0% progress
- For status 1: Progress calculated as `(transferredSize * 90) / totalSize`
- For status 2: 95% progress
- For status 3: 100% progress

The polling mechanism checks status every 5 seconds:
```java
private void startOtaStatusPolling(String serialNumber, OtaProgressCallback callback) {
    otaStatusSubscription = Observable.interval(5, 5, TimeUnit.SECONDS)
        .take(109) // Maximum polling attempts
        .flatMap(tick -> {
            SerialNoEntry entry = new SerialNoEntry();
            entry.setSerialNumber(serialNumber);
            return apiClient.getOtaStatus(entry);
        })
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<GetOtaStatueResponse>() {
            @Override
            public void onNext(GetOtaStatueResponse response) {
                // Process status updates and calculate progress
                // Update callback.onProcess with percentage
                // Check for completion or error states
            }
            
            // Error handling and other methods
        });
}
```

## Usage Context
This endpoint is used in these scenarios:
1. When a user manually initiates a firmware update from device settings
2. When the device requires a mandatory firmware update
3. When an automatic firmware update is scheduled (if otaAutoUpgrade is enabled)

The user flow typically involves:
1. User is notified of available update
2. User clicks to start the update
3. Application calls this endpoint to begin update
4. UI shows progress as polling continues
5. Update completes or fails with appropriate notification

## Error Handling
The application handles several error conditions:
- Network connectivity issues result in `-10402` error
- Server-side errors are indicated by negative result codes
- Timeout conditions trigger appropriate error messages
- Missing or incomplete update data is caught and reported

When errors occur, the OTA process is stopped using the `stopOta()` method, which clears all active subscriptions.

## Related Endpoints
- `/device/otastatus/`: Checks the current status of an OTA update
- `/device/otadetail/`: Retrieves details about available firmware updates
- `/device/otaignore/`: Ignores/postpones an available update
- `/dev/otalist`: Lists available firmware versions

## Notes
- The OTA process is critical for device maintenance and security
- The endpoint initiates the process but monitoring happens through the status endpoint
- The application implements a robust polling mechanism with timeout handling
- Progress reporting is calculated based on different stages of the update process
- For successful updates, the device may reboot and be temporarily offline