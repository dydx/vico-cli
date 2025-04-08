# OTA Status Endpoint Discovery

## Endpoint Information
- **Path:** `/device/otastatus/`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Retrieves the current status of an Over-The-Air (OTA) firmware update for a specific device

## Request Parameters
The endpoint takes a `SerialNoEntry` object in the request body:

```json
{
  "serialNumber": "string",  // Required: The serial number of the device
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
The endpoint returns a `GetOtaStatueResponse` object:

```json
{
  "result": 0,              // Integer result code (0 for success, negative values for errors)
  "msg": "string",          // Human-readable message
  "data": {
    "serialNumber": "string", // Device serial number
    "inProgress": 0,          // Flag indicating if update is in progress (0 or 1)
    "status": 0,              // General status code
    "otaStatus": 0,           // OTA specific status code (0-5)
    "transferredSize": 0,     // Bytes transferred so far
    "totalSize": 0,           // Total bytes to transfer
    "targetFirmware": "string", // Target firmware version
    "localDataProgress": 0,   // Local data transfer progress
    "localDataStartTime": 0   // Timestamp when local data transfer started
  }
}
```

### OTA Status Codes
The `otaStatus` field in the response can have the following values:
- `0`: Not started or waiting
- `1`: In progress (downloading firmware)
- `2`: Installing
- `3`: Completed successfully
- `4`: Failed with an error
- `5`: Failed due to timeout

### Error Codes
- `-2131`: General error that should stop the OTA process
- `-10401`: Timeout error
- `-10402`: Network error
- `-3111`: Low power error
- `-2133`: Device is in sleep mode

## Code Analysis
The endpoint is implemented in several interfaces and used by the OTA management modules:

Example call signature from decompiled code:
```java
@POST("device/otastatus/")
Observable<GetOtaStatueResponse> getOtaStatus(@Body SerialNoEntry serialNoEntry);
```

The application primarily uses this endpoint through the `DeviceOTACore.getOtaStatus()` method:
```java
public void getOtaStatus(String serialNumber, OtaStatusCallback callback) {
    SerialNoEntry entry = new SerialNoEntry();
    entry.setSerialNumber(serialNumber);
    
    apiClient.getOtaStatus(entry)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<GetOtaStatueResponse>() {
            @Override
            public void onNext(GetOtaStatueResponse response) {
                if (response.getResult() < 0 || response.getData() == null) {
                    callback.onError(response.getResult(), response.getMsg());
                    return;
                }
                
                callback.onSuccess(response.getData());
            }
            
            @Override
            public void onError(Throwable e) {
                callback.onError(-10402, e.getMessage());
            }
            
            // Other observer methods
        });
}
```

## Usage Context
This endpoint is used in the following scenarios:
1. To check the status of an ongoing OTA update
2. To track the progress of firmware downloads
3. To determine if an update has completed or failed
4. As part of a polling mechanism that regularly checks update status during the update process

The application implements a polling mechanism that checks every 5-6 seconds during an update:
- Polling continues until the update completes, fails, or times out
- The status values are used to calculate progress percentage for UI updates:
  - When `otaStatus` is 1, progress = (transferredSize * 90) / totalSize
  - When `otaStatus` is 2, progress = 95% (installing)
  - When `otaStatus` is 3, progress = 100% (complete)

## Error Handling
The response includes standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by negative result codes with error messages
- Network errors result in `-10402` error code
- Timeouts result in `-10401` error code
- Low battery conditions may cause `-3111` error code
- Device in sleep mode returns `-2133` error code

## Related Endpoints
- `/device/otastart/`: Initiates an OTA update for a device
- `/device/otaignore/`: Ignores a pending OTA update
- `/dev/otalist`: Lists available firmware versions
- `/device/otadetail/`: Provides detailed information about available updates

## Notes
- This endpoint is critical for the OTA update process user experience
- It provides both status information and progress details for UI feedback
- The application handles transitional states between different OTA phases
- Error codes provide specific information about why an update might fail