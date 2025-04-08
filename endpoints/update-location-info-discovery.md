# Update Location Info Endpoint Discovery

## Endpoint Information
- **Path:** `/location/updatelocationinfo`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Updates an existing location's information

## Request Parameters
The endpoint takes an `UpdateLocationEntry` object in the request body:

```json
{
  "id": number,                // Required: ID of the location to update
  "locationName": "string",    // Required: Name of the location
  "country": "string",         // Country where the location is situated
  "state": "string",           // State/province
  "city": "string",            // City
  "district": "string",        // District/neighborhood
  "streetAddress1": "string",  // Primary street address
  "streetAddress2": "string",  // Secondary address information
  "postalCode": "string",      // Postal/ZIP code
  "adminId": number,           // Admin ID (if applicable)
  "insertTime": number,        // Original insertion time
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
The endpoint is implemented in multiple API interfaces:

Example call signature from decompiled code:
```java
@POST("/location/updatelocationinfo")
Observable<BaseResponse> updateUserDeviceLocation(@Body UpdateLocationEntry updateLocationEntry);
```

Implementation in the application:
```java
public final void updateUserDeviceLocation(LocationBean location, final DeviceLocationCallback callback) {
    UpdateLocationEntry updateLocationEntry = convertBeanToEntry(location);
    
    apiClient.updateUserDeviceLocation(updateLocationEntry)
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
                    callback.onSuccess(location);
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
1. When a user edits an existing location in the device management section
2. When updating address or name information for an existing location
3. When refining or correcting location details after initial creation

The typical user flow:
1. User navigates to location management and selects an existing location
2. User modifies location details (name, address, etc.)
3. User submits the form, triggering this endpoint call
4. On success, the updated location appears in the locations list

## Required Parameters
The essential parameters for this endpoint are:
- `id`: The unique identifier of the location to update
- `locationName`: The user-defined name for this location

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by non-zero result codes with error messages
- Network errors are handled through the RxJava Observable pattern

Common errors include:
- Invalid location ID
- Missing required fields
- Location name already exists
- Authentication issues
- Server-side errors

## Related Endpoints
- `/location/listlocation` - Lists all locations
- `/location/insertlocation/` - Creates a new location
- `/location/deletelocation` - Deletes a location
- `/device/updatedevicelocation` - Moves a device to a different location

## Notes
- This endpoint requires the ID of an existing location
- Unlike creation, the update endpoint returns a simple success/error response without the updated location data
- The application typically refreshes the location list after a successful update
- The postalCode field supports various formats to accommodate international addresses
- Location data is important for organizing devices geographically in the application