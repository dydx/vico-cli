# Delete Location Endpoint Discovery

## Endpoint Information
- **Path:** `/location/deletelocation`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Deletes a location and returns a list of affected devices

## Request Parameters
The endpoint takes a `DeleteLocationEntry` object in the request body:

```json
{
  "id": number,               // Required: ID of the location to delete
  "app": {                    // Standard BaseEntry fields
    "type": "string",         // Application type identifier
    "version": "string"       // Application version
  },
  "countryNo": "string",      // User's country code
  "language": "string",       // User's language preference 
  "tenantId": "string"        // User's tenant ID
}
```

## Response Format
The endpoint returns a `DeleteLocationResponse` which extends the `BaseResponse`:

```json
{
  "result": 0,             // Result code (0 indicates success)
  "msg": "string",         // Response message
  "data": {
    "list": [              // Array of devices that were affected by the location deletion
      {
        "serialNumber": "string",  // Device serial number
        "deviceName": "string",    // Name of the device
        "online": 0|1,             // Device online status
        "deviceModel": "string",   // Model information
        "batteryLevel": number,    // Battery level percentage
        // Other DeviceBean fields...
      }
    ]
  }
}
```

## Code Analysis
The endpoint is implemented in multiple API interfaces:

Example call signature from decompiled code:
```java
@POST("/location/deletelocation")
Observable<DeleteLocationResponse> deleteLocation(@Body DeleteLocationEntry deleteLocationEntry);
```

Implementation in the application:
```java
public final void deleteUserDeviceLocation(int locationId, final Callback<List<DeviceBean>> callback) {
    DeleteLocationEntry deleteLocationEntry = new DeleteLocationEntry();
    deleteLocationEntry.setId(locationId);
    
    apiClient.deleteLocation(deleteLocationEntry)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<DeleteLocationResponse>() {
            @Override
            public void onNext(DeleteLocationResponse response) {
                if (response.getResult() < 0 || response.getData() == null) {
                    if (callback != null) {
                        callback.onError(response.getResult(), response.getMsg());
                    }
                    return;
                }
                
                List<DeviceBean> devices = response.getData().getList();
                if (callback != null) {
                    callback.onSuccess(devices);
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
1. When a user wants to delete a location from their account
2. During cleanup of unused locations
3. When reorganizing device management structure

The typical user flow:
1. User navigates to location management
2. User selects a location to delete
3. User confirms the deletion
4. App calls this endpoint with the location ID
5. App processes the response which includes affected devices
6. UI updates to remove the location and possibly reassign devices

## Important Considerations
1. The only required parameter is the location ID
2. The endpoint returns affected devices, which suggests the location is deleted even if devices are assigned to it
3. The response data includes a list of devices that were previously associated with the deleted location
4. No explicit confirmation or validation parameters are required in the request

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by non-zero result codes with error messages
- Network errors are handled through the RxJava Observable pattern

Common errors include:
- Invalid location ID
- Attempt to delete a protected system location
- Authentication issues
- Server-side errors

## Related Endpoints
- `/location/listlocation` - Lists all locations
- `/location/insertlocation/` - Creates a new location
- `/location/updatelocationinfo` - Updates an existing location
- `/device/updatedevicelocation` - Might be used to reassign devices after location deletion

## Notes
- This endpoint requires only the ID of the location to delete
- The response includes devices that were affected by the deletion, which allows the application to handle device reassignment
- When a location is deleted, the devices previously assigned to it need to be reassigned to another location
- The application likely handles reassignment of devices shown in the response
- Special system locations (All Devices, Shared Devices) cannot be deleted as they are client-side only