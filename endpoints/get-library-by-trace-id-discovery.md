# Select Single Library Endpoint Discovery

## Endpoint Information
- **Path:** `/library/newselectsinglelibrary`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Retrieves detailed information about a single library record by its trace ID

## Request Parameters
The endpoint takes a `TraceIdEntry` object in the request body which extends the `BaseEntry`:

```json
{
  "traceId": "string",       // Required: The library record's trace ID
  "app": {                   // From BaseEntry
    "type": "string",        // Application type identifier
    "version": "string"      // Application version
  },
  "countryNo": "string",     // User's country code
  "language": "string",      // User's language preference
  "tenantId": "string"       // User's tenant ID
}
```

## Response Format
The endpoint returns a `SingleLibraryResponse` containing detailed recording information:

```json
{
  "result": 0,               // 0 indicates success
  "msg": "string",           // Response message
  "data": {                  // RecordBean object
    "traceId": "string",     // Unique identifier for this record
    "serialNumber": "string", // Device serial number
    "deviceName": "string",   // Device name
    "locationId": integer,    // Location ID
    "locationName": "string", // Location name
    "timestamp": integer,     // Unix timestamp for the recording
    "date": "string",         // Formatted date string
    "type": integer,          // Recording type
    "imageUrl": "string",     // Thumbnail image URL
    "videoUrl": "string",     // Video recording URL
    "period": number,         // Duration of recording in seconds
    "fileSize": integer,      // File size in bytes
    "imageOnly": integer,     // Whether it's only an image (1) or has video (0)
    "mediaType": "string",    // "image" or "video"
    "videoEvent": "string",   // Event type identifier
    "tags": "string",         // Tags associated with the recording
    "marked": integer,        // Whether the recording is marked/favorited (0=no, 1=yes)
    "missing": integer,       // Whether the recording has been viewed (0=viewed, 1=unviewed)
    "userId": integer,        // User ID
    "userName": "string",     // User name
    "adminId": integer,       // Admin ID
    "adminName": "string",    // Admin name
    "adminIsVip": boolean,    // Whether the admin has VIP status
    "role": "string",         // User role
    "timeZone": integer,      // Time zone offset
    "timeFormat": integer,    // Time format preference
    "dst": integer,           // Daylight saving time flag
    "activityZoneName": "string", // Activity zone name if applicable
    "hasPossibleSubcategory": boolean, // Whether subcategories are possible
    "deviceAiEventList": [     // List of AI event types detected
      "string"
    ],
    "eventInfoList": [         // List of event information
      "string"
    ],
    "subcategoryInfoList": [   // List of subcategory information
      {
        // Subcategory info structure
      }
    ]
  }
}
```

## Code Analysis
The endpoint is implemented in the `vicohome_1742553098674_00O0o0oOO` interface:

```java
@vicohome_1742553098674_0OO00OOo("/library/newselectsinglelibrary")
Observable<SingleLibraryResponse> getLibraryByTraceId(@vicohome_1742553098674_00O0ooOO0 TraceIdEntry traceIdEntry);
```

The `LibraryCore` class provides the implementation that calls this endpoint:

```java
public final void getLibraryByTraceId(String traceId, com.ai.addxbase.vicohome_1742553098674_0O0oO0O<RecordBean> callBack) {
    vicohome_1742553098674_0O0oOoo.checkNotNullParameter(traceId, "traceId");
    vicohome_1742553098674_0O0oOoo.checkNotNullParameter(callBack, "callBack");
    this.vicohome_1742553098674_0o00OOoOo.add(LibraryApiClient.vicohome_1742553098674_o00OOoo.getSInstance().getLibraryByTraceId(new TraceIdEntry(traceId)).subscribeOn(Schedulers.io()).observeOn(AndroidSchedulers.mainThread()).subscribe((Subscriber<? super SingleLibraryResponse>) new vicohome_1742553098674_0O0oO00(callBack)));
}
```

## Usage Context
This endpoint is commonly used in the following scenarios:

1. When viewing recording details after selecting a specific event notification
2. When sharing a specific recording with others
3. When downloading a single recording for offline viewing
4. When displaying detailed information about an event captured

The `NotifyHandlerViewModel` class shows one specific usage of this endpoint to fetch recording details when handling notifications:

```java
public final void getSingleLibraryByLibraryId(String traceId, Bundle bundle) {
    vicohome_1742553098674_0O0oOoo.checkNotNullParameter(traceId, "traceId");
    vicohome_1742553098674_0O0oOoo.checkNotNullParameter(bundle, "bundle");
    LibraryCore.vicohome_1742553098674_o00OOoo.getSInstance().getLibraryByTraceId(traceId, new vicohome_1742553098674_00O0o0oOO(bundle));
}
```

## Error Handling
The endpoint follows standard error handling practices:
- Success is indicated by `result: 0`
- Errors are indicated by non-zero result codes with error messages
- Common error conditions include record not found, authentication issues, and network failures

Error handling in the `NotifyHandlerViewModel` shows how application errors are processed:

```java
public void onError(int i, String str) {
    NotifyHandlerViewModel.this.getNotifyState().postValue(new Pair<>(State.ERROR_INTERNAL, new Pair(null, this.vicohome_1742553098674_0o00OOoo0)));
}
```

## Notes
- This endpoint is critical for displaying detailed information about a specific recording
- The `marked` and `missing` fields are used to manage user interaction states with recordings
- The response includes multiple URL formats for thumbnails and videos that can be used for displaying or downloading content
- AI event detection information is included when the camera has detected specific objects or events
- The response data can be used for sharing recordings through the app's sharing functionality