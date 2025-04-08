# Delete Activity Zone Endpoint Discovery

## Endpoint Information
- **Path:** `/device/deleteactivityzone`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Deletes an existing activity zone for a camera device

## Request Parameters
The endpoint takes a `ZoneDeleteEntry` object in the request body:

```json
{
  "id": 12345,              // Required: Zone ID to delete
  "serialNumber": "string", // Required: Device serial number
  "language": "string",     // Optional: User's language preference
  "requestId": "string",    // Optional: Request identifier for tracking
  "app": {                  // Standard BaseEntry fields
    "type": "string",       // Application type identifier
    "version": "string"     // Application version
  },
  "countryNo": "string",    // User's country code
  "tenantId": "string"      // User's tenant ID
}
```

## Response Format
The endpoint returns a `ZoneDeleteResponse` object which extends the `BaseResponse`:

```json
{
  "result": 0,            // Result code (0 indicates success)
  "msg": "string"         // Response message
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -2002 | Device access error - user does not have permission to access this device |
| -2001 | Other device operation errors |
| -1 | General error - failed to delete activity zone |

## Code Analysis
The endpoint is implemented in multiple interfaces:

Example call signature from decompiled code:
```java
@POST("/device/deleteactivityzone")
Observable<ZoneDeleteResponse> deleteActivityZone(@Body ZoneDeleteEntry zoneDeleteEntry);
```

The `DeviceActivityZoneCore` class provides a higher-level implementation:
```java
public final void deleteActivityZone(int zoneId, String serialNumber, Callback<BaseResponse> callBack) {
    ZoneDeleteEntry zoneDeleteEntry = new ZoneDeleteEntry();
    zoneDeleteEntry.setId(zoneId);
    zoneDeleteEntry.setSerialNumber(serialNumber);
    
    DeviceSettingApiClient.getInstance().deleteActivityZone(zoneDeleteEntry)
        .subscribeOn(Schedulers.io())
        .observeOn(AndroidSchedulers.mainThread())
        .subscribe(new Subscriber<ZoneDeleteResponse>() {
            @Override
            public void onNext(ZoneDeleteResponse response) {
                if (response.getResult() < 0) {
                    if (callBack != null) {
                        callBack.onError(response.getResult(), response.getMsg());
                    }
                    return;
                }
                
                if (callBack != null) {
                    callBack.onSuccess(response);
                }
            }
            
            @Override
            public void onError(Throwable e) {
                if (callBack != null) {
                    callBack.onError(-1, e.getMessage());
                }
            }
            
            // Other subscriber methods
        });
}
```

## Usage Context
This endpoint is used in the following scenarios:
1. When a user selects to delete an activity zone in the camera settings
2. In the ZoneActivity UI where users can manage their activity zones
3. After retrieving zones with `/device/getactivityzone` and selecting one for removal

The typical user flow:
1. User navigates to camera settings → activity zones
2. App loads existing zones using `/device/getactivityzone`
3. User selects a zone to delete
4. User confirms deletion
5. App calls this endpoint with the zone ID and device serial number
6. On success, the zone is removed from the list

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by non-zero result codes with error messages
- Common errors include permission issues (-2002) and general failures (-1)
- Network errors are handled through the RxJava Observable pattern

## Related Endpoints
- `/device/insertactivityzone` - Creates a new activity zone
- `/device/getactivityzone` - Retrieves activity zones for a device
- `/device/updateactivityzone` - Updates an existing activity zone

## Notes
- Activity zones are an important feature for reducing false alerts
- Deleting a zone requires both the zone ID and device serial number
- The application typically updates the local cache after a successful deletion
- Access to modify zones is restricted - the endpoint returns error -2002 if the user doesn't have permission
- Upon successful deletion, the UI is updated to reflect the removal of the zone
- The endpoint is typically called from the ZoneViewModel when a user selects the delete option