# Insert Location Endpoint Discovery

## Endpoint Information
- **Path:** `/location/insertlocation/`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Creates a new location for organizing devices

## Request Parameters
The endpoint takes a `LocationEntry` object in the request body:

```json
{
  "locationName": "string",     // Required: Name of the location (e.g., "Home", "Office")
  "country": "string",          // Country where the location is situated
  "state": "string",            // State/province
  "city": "string",             // City
  "district": "string",         // District/neighborhood
  "streetAddress1": "string",   // Primary street address
  "streetAddress2": "string",   // Secondary address information (apt, suite, etc.)
  "postalCode": "string",       // Postal/ZIP code
  "adminId": number,            // Admin ID (if applicable)
  "id": 0,                      // ID value (usually 0 for new locations)
  "insertTime": 0,              // Time of insertion (usually 0, server will assign)
  "app": {                      // Standard BaseEntry fields
    "type": "string",           // Application type identifier
    "version": "string"         // Application version
  },
  "countryNo": "string",        // User's country code
  "language": "string",         // User's language preference 
  "tenantId": "string"          // User's tenant ID
}
```

## Response Format
The endpoint returns a response object with the following structure:

```json
{
  "result": 0,               // Result code (0 indicates success)
  "msg": "string",           // Response message
  "data": {                  // Location data (only present on success)
    "id": number,            // Generated location ID
    "locationName": "string", // Name of the location
    "country": "string",     // Country
    "state": "string",       // State/province
    "city": "string",        // City
    "district": "string",    // District/neighborhood
    "streetAddress1": "string", // Primary street address
    "streetAddress2": "string", // Secondary address information
    "postalCode": "string",  // Postal/ZIP code
    "adminId": number,       // Admin ID
    "insertTime": number     // Server timestamp of when the location was created
  }
}
```

## Code Analysis
The endpoint is implemented in multiple API interfaces:

Example call signature from decompiled code:
```java
@POST("/location/insertlocation/")
Observable<InsertLocationResponse> createUserDeviceLocation(@Body LocationEntry locationEntry);
```

Implementation in the application:
```java
public final void createUserDeviceLocation(LocationBean location, final DeviceLocationCallback callback) {
    LocationEntry locationEntry = convertBeanToEntry(location);
    
    apiClient.createUserDeviceLocation(locationEntry)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<InsertLocationResponse>() {
            @Override
            public void onNext(InsertLocationResponse response) {
                if (response.getResult() < 0 || response.getData() == null) {
                    if (callback != null) {
                        callback.onError(response.getResult(), response.getMsg());
                    }
                    return;
                }
                
                LocationBean location = convertResponseToBean(response.getData());
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
1. When a user adds a new location in the device management section
2. During initial setup when a user wants to assign a location for their device(s)
3. When organizing multiple devices into different physical locations

The typical user flow:
1. User navigates to location management
2. User selects to add a new location
3. User enters location details (name, address, etc.)
4. User submits the form, triggering this endpoint call
5. On success, the new location appears in the locations list

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by non-zero result codes with error messages
- Network errors are handled through the RxJava Observable pattern

Common errors include:
- Required fields missing (especially locationName)
- Location already exists with the same name
- User has exceeded maximum number of allowed locations
- Authentication issues
- Server-side errors

## Related Endpoints
- `/location/listlocation` - Lists all locations
- `/location/updatelocationinfo` - Updates an existing location
- `/location/deletelocation` - Deletes a location
- `/device/binddevicelocation` - Associates a device with a location
- `/device/updatedevicelocation` - Moves a device to a different location

## Notes
- Locations are organizational units for grouping devices in specific physical places
- The locationName is the only required field; address fields are optional
- The server assigns a unique ID to each new location
- The application typically refreshes the location list after a successful creation
- Locations can be managed from the device settings area of the application
- The trailing slash in the endpoint path ("/location/insertlocation/") is intentional