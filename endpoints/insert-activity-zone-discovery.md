# Insert Activity Zone Endpoint Discovery

## Endpoint Information
- **Path:** `/device/insertactivityzone`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Adds a new activity zone for a device camera

## Request Parameters
The endpoint takes a `ZoneAddEntry` or `ZoneBean` object in the request body:

```json
{
  "serialNumber": "string",  // Required: Device serial number
  "zoneName": "string",      // Required: Name of the activity zone
  "vertices": "string",      // Required: Coordinates that define the zone boundary
  "needAlarm": 0|1,          // Required: Flag to enable/disable alarm (0=disabled, 1=enabled)
  "needPush": 0|1,           // Required: Flag to enable/disable push notifications
  "needRecord": 0|1,         // Required: Flag to enable/disable recording for events
  "language": "string",      // Optional: User's language preference
  "requestId": "string",     // Optional: Request identifier for tracking
  "app": {                   // Standard BaseEntry fields
    "type": "string",        // Application type identifier
    "version": "string"      // Application version
  },
  "countryNo": "string",     // User's country code
  "tenantId": "string"       // User's tenant ID
}
```

## Response Format
The endpoint returns a `ZoneAddResponse` object which extends the `BaseResponse`:

```json
{
  "result": 0,            // Result code (0 indicates success)
  "msg": "string",        // Response message
  "data": {
    "id": 12345           // Integer ID of the newly created activity zone
  }
}
```

## Code Analysis
The endpoint is implemented in at least two API interfaces:
1. `com.ai.addxnet.vicohome_1742553098674_00O0o0oOO`
2. `com.ai.settingcore.vicohome_1742553098674_00O0o0oOO`

Example call signature from decompiled code:
```java
@POST("/device/insertactivityzone")
Observable<ZoneAddResponse> addActivityZone(@Body ZoneAddEntry zoneAddEntry);
```

The `DeviceActivityZoneCore` class provides a higher-level implementation:
```java
public final void createActivityZone(ZoneBean zone, Callback<BaseResponse> callBack) {
    DeviceSettingApiClient.getInstance()
        .addActivityZone(zone)
        .subscribeOn(Schedulers.io())
        .observeOn(AndroidSchedulers.mainThread())
        .subscribe(new Subscriber<ZoneAddResponse>() {
            @Override
            public void onNext(ZoneAddResponse response) {
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

## Vertices Format
The `vertices` parameter is a string representation of coordinate points that define the polygon shape of the activity zone. The format appears to be a comma-separated list of x,y coordinates, where each point is separated from the next.

Example: `"10,20,150,20,150,100,10,100"` would represent a rectangle with corners at (10,20), (150,20), (150,100), and (10,100).

## Usage Context
This endpoint is used in the following scenarios:
1. When a user creates a new activity zone in the camera settings
2. During the configuration of motion detection settings for specific areas
3. As part of the camera setup or reconfiguration process

The typical user flow:
1. User navigates to camera settings → activity zones
2. User chooses to add a new zone
3. User draws the zone boundaries on the camera preview image
4. User provides a name and configures notification/recording settings
5. User submits the form, which triggers this endpoint call
6. On success, the new zone appears in the list of active zones

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by non-zero result codes with error messages
- Network errors are handled through the RxJava Observable pattern

## Related Endpoints
- `/device/getactivityzone` - Retrieves activity zones for a device
- `/device/updateactivityzone` - Updates an existing activity zone
- `/device/deleteactivityzone` - Deletes an activity zone

## Notes
- Activity zones are an important feature for reducing false alerts
- They allow users to specify regions of interest in the camera's field of view
- Each zone can have its own notification and recording settings
- The zone ID returned in the response is used for subsequent zone management operations
- The vertices string representation is a compact way to define polygon shapes
- Zone settings (needAlarm, needPush, needRecord) determine what actions are taken when motion is detected in that specific zone