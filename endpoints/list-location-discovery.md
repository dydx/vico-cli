# List Location Endpoint Discovery

## Endpoint Information
- **Path:** `/location/listlocation`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Retrieves a list of all locations created by the user

## Request Parameters
The endpoint takes a standard `BaseEntry` object in the request body:

```json
{
  "app": {                  // Standard BaseEntry fields
    "type": "string",       // Application type identifier
    "version": "string"     // Application version
  },
  "countryNo": "string",    // User's country code
  "language": "string",     // User's language preference 
  "tenantId": "string"      // User's tenant ID
}
```

## Response Format
The endpoint returns an `AllLocationResponse` object which extends the `BaseResponse`:

```json
{
  "result": 0,              // Result code (0 indicates success)
  "msg": "string",          // Response message
  "data": [
    {
      "id": 0,                 // Location ID
      "adminId": 0,            // User ID of the admin
      "locationName": "string", // Name of the location
      "country": "string",     // Country
      "state": "string",       // State/province
      "city": "string",        // City
      "district": "string",    // District
      "streetAddress1": "string", // Street address line 1
      "streetAddress2": "string", // Street address line 2
      "postalCode": "string",  // Postal/ZIP code
      "insertTime": 0          // Creation timestamp
    },
    // Additional location entries...
  ]
}
```

## Special Location ID Values
The application defines two special location ID constants:
- `ID_ALL_DEVICE = -666`: Represents all devices regardless of location
- `ID_SHAREID_DEVICE = -999`: Represents shared devices

## Code Analysis
The endpoint is implemented in multiple API interfaces:

Example call signature from decompiled code:
```java
@POST("/location/listlocation")
Observable<AllLocationResponse> getUserDeviceLocationList(@Body BaseEntry baseEntry);
```

Implementation in the application:
```java
public final void getUserDeviceLocationList(final DeviceLocationCallback callback) {
    BaseEntry baseEntry = new BaseEntry();
    
    apiClient.getUserDeviceLocationList(baseEntry)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<AllLocationResponse>() {
            @Override
            public void onNext(AllLocationResponse response) {
                if (response.getResult() < 0 || response.getData() == null) {
                    if (callback != null) {
                        callback.onError(response.getResult(), response.getMsg());
                    }
                    return;
                }
                
                List<LocationBean> locations = convertResponseToBean(response.getData());
                
                // Add special locations
                LocationBean allDevices = new LocationBean();
                allDevices.setId(ID_ALL_DEVICE);
                allDevices.setLocationName("All Devices");
                allDevices.setLocalData(true);
                
                LocationBean sharedDevices = new LocationBean();
                sharedDevices.setId(ID_SHAREID_DEVICE);
                sharedDevices.setLocationName("Shared Devices");
                sharedDevices.setLocalData(true);
                
                locations.add(0, allDevices);
                locations.add(sharedDevices);
                
                if (callback != null) {
                    callback.onSuccess(locations);
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
1. When loading the home screen to organize devices by location
2. When viewing the location management screen
3. When assigning a device to a location
4. When filtering devices by location

The typical user flow:
1. User opens the app or navigates to a device listing screen
2. App calls this endpoint to retrieve all locations
3. UI displays locations as filters or organization groups
4. User can select a location to view devices in that location

## Client-Side Location Handling
The client application adds two special locations to the response list:
1. "All Devices" (ID: -666) - Shows all devices regardless of location
2. "Shared Devices" (ID: -999) - Shows devices shared with the user

These special locations are marked as `isLocalData = true` to indicate they are client-side constructs.

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by negative result codes with error messages
- Network errors are handled through the RxJava Observable pattern

## Related Endpoints
- `/location/insertlocation/` - Creates a new location
- `/location/deletelocation` - Deletes an existing location
- `/location/updatelocationinfo` - Updates an existing location
- `/device/updatedevicelocation` - Reassigns a device to a different location

## Notes
- Locations are used to organize and group devices in the application
- The endpoint returns both standard and address information for each location
- The client adds special locations not stored on the server
- Location IDs are used to filter devices in the device list
- The distinction between server locations and local-only locations is tracked through the `isLocalData` flag
- The response is sorted by insertion time, with special locations added at specific positions