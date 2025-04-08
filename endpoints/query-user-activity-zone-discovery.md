# Query User Activity Zones Endpoint Discovery

## Endpoint Information
- **Path:** `/device/queryUserActivityZone`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Retrieves all activity zones defined across all devices for the current user

## Request Parameters
The endpoint takes no specific request parameters. Authentication is handled via session token in the request headers.

## Response Format
The endpoint returns a `UserAllAZResponse` object which extends the `BaseResponse`:

```json
{
  "result": 0,           // Result code (0 indicates success)
  "msg": "string",       // Response message
  "data": [              // Array of ZoneBean objects
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
```

## Vertices Format
The `vertices` parameter in the response is a string representation of coordinate points that define the polygon shape of the activity zone. The format is a comma-separated list of x,y coordinates, where each point is separated from the next.

Example: `"10,20,150,20,150,100,10,100"` would represent a rectangle with corners at (10,20), (150,20), (150,100), and (10,100).

## Code Analysis
This endpoint is implemented in multiple interfaces:
1. `com.ai.addxnet.vicohome_1742553098674_00O0o0oOO`
2. `com.a4x.a4xlibrary.vicohome_1742553098674_00O0o0oOO`

Example call signature from decompiled code:
```java
@POST("/device/queryUserActivityZone")
Observable<UserAllAZResponse> queryUserActivityZone();
```

The `LibraryCore` class provides an implementation for calling this endpoint:
```java
public final void queryUserActivityZone(com.ai.addxbase.vicohome_1742553098674_0O0oO0O<List<ZoneBean>> callBack) {
    vicohome_1742553098674_0O0oOoo.checkNotNullParameter(callBack, "callBack");
    this.vicohome_1742553098674_0o00OOoOo.add(
        LibraryApiClient.vicohome_1742553098674_o00OOoo.getSInstance()
            .queryUserActivityZone()
            .subscribeOn(Schedulers.io())
            .observeOn(AndroidSchedulers.mainThread())
            .subscribe((Subscriber<? super UserAllAZResponse>) new vicohome_1742553098674_0O0oO0O0(callBack))
    );
}
```

The `ZoneViewModel` class has a method that calls this function:
```java
public void queryUserActivityZone() {
    this.vicohome_1742553098674_o0oOO.postValue(new Pair<>(RxViewModel.State.LOADING, null));
    LibraryCore.vicohome_1742553098674_o00OOoo.getSInstance().queryUserActivityZone(new vicohome_1742553098674_0O0o0oo());
}
```

## Usage Context
This endpoint is used in the following scenarios:
1. When a user wants to view all activity zones across all their devices
2. When searching or filtering activity zones across multiple cameras
3. For presenting a consolidated view of all zones in the system

The typical user flow:
1. User navigates to a cross-device activity zone management screen
2. App calls this endpoint to retrieve zones from all devices
3. UI displays the comprehensive list of all zones
4. User can select zones to view, edit, or interact with them further

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by non-zero result codes with error messages
- Network errors are handled through the RxJava Observable pattern
- The implementation in `LibraryCore` transforms network or parsing errors into callback error states

## Related Endpoints
- `/device/getactivityzone` - Gets activity zones for a specific device
- `/device/insertactivityzone` - Creates a new activity zone
- `/device/updateactivityzone` - Updates an existing activity zone
- `/device/deleteactivityzone` - Deletes an activity zone

## Notes
- Unlike the `/device/getactivityzone` endpoint which requires a specific device serial number, this endpoint retrieves zones across all devices
- This is useful for apps that need to consolidate zone information from multiple cameras
- The response includes the device name and serial number for each zone, allowing for proper organization in the UI
- The endpoint doesn't require parameters, making it simpler to use for global zone queries
- Each zone includes its own settings (needAlarm, needPush, needRecord) which determine notification behaviors when motion is detected