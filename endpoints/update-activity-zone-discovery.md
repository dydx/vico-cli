# Update Activity Zone Endpoint Discovery

## Endpoint Information
- **Path:** `/device/updateactivityzone`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Updates an existing activity zone for a camera device

## Request Parameters
The endpoint takes a `ZoneUpdateEntry` or `ZoneBean` object in the request body:

```json
{
  "id": 12345,              // Required: Zone ID to update
  "serialNumber": "string", // Required: Device serial number
  "zoneName": "string",     // Required: Name of the activity zone
  "vertices": "string",     // Required: Coordinates that define the zone boundary
  "needAlarm": 0|1,         // Required: Flag to enable/disable alarm (0=disabled, 1=enabled)
  "needPush": 0|1,          // Required: Flag to enable/disable push notifications
  "needRecord": 0|1,        // Required: Flag to enable/disable recording for events
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
The endpoint returns a `ZoneUpdateResponse` object which extends the `BaseResponse`:

```json
{
  "result": 0,            // Result code (0 indicates success)
  "msg": "string"         // Response message
}
```

## Vertices Format
The `vertices` parameter is a string representation of coordinate points that define the polygon shape of the activity zone. The format is a comma-separated list of x,y coordinates, where each point is separated from the next.

Example: `"0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8"` would define a polygon with four points at (0.1,0.2), (0.3,0.4), (0.5,0.6), and (0.7,0.8).

## Code Analysis
The endpoint is implemented in multiple interfaces:

Example call signature from decompiled code:
```java
@POST("/device/updateactivityzone")
Observable<ZoneUpdateResponse> updateActivityZone(@Body ZoneUpdateEntry zoneUpdateEntry);
```

The `DeviceActivityZoneCore` class provides a higher-level implementation:
```java
public final void updateActivityZone(ZoneBean zone, Callback<BaseResponse> callBack) {
    DeviceSettingApiClient.getInstance().updateActivityZone(zone)
        .subscribeOn(Schedulers.io())
        .observeOn(AndroidSchedulers.mainThread())
        .subscribe(new Subscriber<ZoneUpdateResponse>() {
            @Override
            public void onNext(ZoneUpdateResponse response) {
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
1. When a user modifies an existing activity zone in the camera settings
2. After retrieving zones with `/device/getactivityzone` and making changes
3. When the user changes the zone's boundaries, name, or notification settings

The typical user flow:
1. User navigates to camera settings → activity zones
2. App loads existing zones using `/device/getactivityzone`
3. User selects a zone to edit
4. User modifies the zone (boundaries, name, settings)
5. User saves changes, triggering this endpoint call
6. On success, the updated zone appears in the list

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by non-zero result codes with error messages
- Network errors are handled through the RxJava Observable pattern

## Related Endpoints
- `/device/insertactivityzone` - Creates a new activity zone
- `/device/getactivityzone` - Retrieves activity zones for a device
- `/device/deleteactivityzone` - Deletes an activity zone

## Notes
- Activity zones help users define specific areas for motion detection
- Each zone can have its own notification, alarm, and recording settings
- Updating a zone requires both the zone ID and device serial number
- The application typically updates the local cache after a successful update
- Zone settings (needAlarm, needPush, needRecord) determine the actions taken when motion is detected in that zone
- The vertices coordinates may be normalized values (0.0-1.0) relative to the camera view dimensions