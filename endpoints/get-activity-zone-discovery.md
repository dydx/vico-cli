# Get Activity Zone Endpoint Discovery

## Endpoint Information
- **Path:** `/device/getactivityzone`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Retrieves all activity zones defined for a specific device

## Request Parameters
The endpoint takes a `ZoneGetEntry` object in the request body which extends `BaseEntry`:

```json
{
  "serialNumber": "string",  // Required: Device serial number
  "language": "string",      // Optional: User language preference
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
The endpoint returns a `ZoneGetResponse` object which extends the `BaseResponse`:

```json
{
  "result": 0,            // Result code (0 indicates success)
  "msg": "string",        // Response message
  "data": {
    "list": [             // Array of activity zones
      {
        "id": 12345,             // Zone ID
        "serialNumber": "string", // Device serial number
        "zoneName": "string",     // Custom name of the zone
        "zoneNameId": 0,          // Zone name identifier
        "vertices": "string",     // Coordinates defining the zone boundaries
        "needPush": 0|1,          // Push notification setting (0=disabled, 1=enabled)
        "needAlarm": 0|1,         // Alarm setting (0=disabled, 1=enabled)
        "needRecord": 0|1,        // Recording setting (0=disabled, 1=enabled)
        "modelCategory": 0,       // Device model category
        "deviceName": "string"    // Name of the device
      },
      // Additional zone entries...
    ]
  }
}
```

## Vertices Format
The `vertices` parameter in the response is a string representation of coordinate points that define the polygon shape of the activity zone. The format is a comma-separated list of x,y coordinates, where each point is separated from the next.

Example: `"10,20,150,20,150,100,10,100"` would represent a rectangle with corners at (10,20), (150,20), (150,100), and (10,100).

## Code Analysis
The endpoint is implemented in multiple interfaces:
1. `com.ai.addxnet.vicohome_1742553098674_00O0o0oOO`
2. `com.ai.settingcore.vicohome_1742553098674_00O0o0oOO`
3. `com.ai.addxnet.vicohome_1742553098674_0oOOoOO`

Example call signature from decompiled code:
```java
@POST("/device/getactivityzone")
Observable<ZoneGetResponse> getActivityZone(@Body ZoneGetEntry zoneGetEntry);
```

The `DeviceActivityZoneCore` class provides a higher-level implementation:
```java
public final void getActivityZone(String serialNumber, Callback<List<ZoneBean>> callBack) {
    ZoneGetEntry zoneGetEntry = new ZoneGetEntry();
    zoneGetEntry.setSerialNumber(serialNumber);
    
    DeviceSettingApiClient.getInstance().getActivityZone(zoneGetEntry)
        .subscribeOn(Schedulers.io())
        .observeOn(AndroidSchedulers.mainThread())
        .subscribe(new Subscriber<ZoneGetResponse>() {
            @Override
            public void onNext(ZoneGetResponse response) {
                if (response.getResult() < 0 || response.getData() == null) {
                    if (callBack != null) {
                        callBack.onError(response.getResult(), response.getMsg());
                    }
                    return;
                }
                
                if (callBack != null) {
                    callBack.onSuccess(response.getData().getList());
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
1. When a user navigates to the activity zone settings page for a camera
2. Before editing or deleting an existing activity zone
3. When displaying the camera view with activity zone overlays
4. When initializing the activity zone management UI

The typical user flow:
1. User navigates to camera settings → activity zones
2. App calls this endpoint to retrieve existing zones
3. UI displays the list of zones with their names and settings
4. User can view, edit, or delete these zones

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by non-zero result codes with error messages
- Network errors are handled through the RxJava Observable pattern

## Related Endpoints
- `/device/insertactivityzone` - Creates a new activity zone
- `/device/updateactivityzone` - Updates an existing activity zone
- `/device/deleteactivityzone` - Deletes an activity zone
- `/device/queryUserActivityZone` - Gets all activity zones for the user

## Notes
- Activity zones are displayed as overlays on the camera view
- Each zone can have its own notification, alarm, and recording settings
- These zones help users focus monitoring on specific areas and reduce false alerts
- The zone ID returned in the response is used for subsequent zone management operations
- The application typically caches zone information for performance
- Zone settings (needAlarm, needPush, needRecord) determine what actions are taken when motion is detected in that specific zone