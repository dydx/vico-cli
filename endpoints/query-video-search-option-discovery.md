# Library Query Video Search Options Discovery

## Endpoint Information
- **Path:** `/library/queryVideoSearchOption`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Query available video search filter options for library recordings

## Request Parameters
The endpoint accepts a `LibraryOptionRequest` object in the request body:

```json
{
  "SN": "string",            // Device serial number
  "isFromSDCard": boolean,   // Whether to query from SD card or cloud
  "app": {                   // Standard app information
    "apiVersion": "string",
    "appName": "string",
    "appType": "string",
    "bundle": "string",
    "countlyId": "string", 
    "env": "string",
    "tenantId": "string",
    "timeZone": "string",
    "version": "string",
    "versionName": "string"
  },
  "countryNo": "string",     // Country code
  "language": "string",      // User's language preference
  "tenantId": "string"       // User's tenant ID
}
```

## Response Format
The endpoint returns a `VideoSearchOptionResponse` which contains a `TagBean` object:

```json
{
  "result": 0,               // Result code (0 indicates success)
  "msg": "string",           // Response message
  "data": {                 // TagBean object
    "aiEventTags": [        // List of AI event filter tags
      {
        "name": "string",   // Tag name
        "subTags": [        // Optional sub-tags
          {
            "name": "string",
            "subTags": []
          }
        ]
      }
    ],
    "deviceEventTags": [    // List of device event filter tags
      {
        "name": "string", 
        "subTags": []
      }
    ],
    "operateOptions": [     // List of operation filter tags
      {
        "name": "string",
        "subTags": []
      }
    ],
    "devices": [            // List of available devices
      {
        "serialNumber": "string",  // Device serial number
        "deviceName": "string",    // Device name
        "isBind": boolean,         // Whether device is bound to user
        "modelCategory": integer,  // Device model category
        "roleId": integer          // User role for device (1 = admin)
      }
    ]
  }
}
```

## Code Analysis
The endpoint is implemented in the API interface:

```java
@POST("/library/queryVideoSearchOption")
Observable<VideoSearchOptionResponse> queryVideoSearchOption(@Body LibraryOptionRequest libraryOptionRequest);
```

Implementation in LibraryApiClient:
```java
@Override
public Observable<VideoSearchOptionResponse> queryVideoSearchOption(LibraryOptionRequest entry) {
    checkNotNullParameter(entry, "entry");
    if (entry.isFromSDCard) {
        return this.localApiService.queryVideoSearchOption(entry);
    }
    NetworkManager.getInstance().wrapBaseEntry(entry);
    return this.cloudApiService.queryVideoSearchOption(entry);
}
```

Usage in LibraryCore:
```java
public final void queryVideoSearchOption(boolean fromSDCard, String serialNumber, Callback<TagBean> callBack) {
    checkNotNullParameter(serialNumber, "serialNumber");
    checkNotNullParameter(callBack, "callBack");
    
    this.subscriptions.add(
        LibraryApiClient.Companion.getSInstance()
            .queryVideoSearchOption(new LibraryOptionRequest(fromSDCard, serialNumber))
            .subscribeOn(Schedulers.io())
            .observeOn(AndroidSchedulers.mainThread())
            .subscribe(new VideoSearchOptionSubscriber(callBack))
    );
}
```

## Usage Context
This endpoint is used in the following scenarios:

1. Library/Recording Filtering:
   - To populate filter options for video recordings in the library view
   - To get available device list for filtering recordings across multiple devices
   - To get AI event tags for filtering recordings by event type (person, animal, vehicle, etc.)

2. Filter UI:
   - Called when the user navigates to the filter screen for recording library
   - Used to build the filter UI with checkboxes/options for each filter category
   - Supports both local (SD card) and cloud recordings based on the isFromSDCard parameter

3. Error handling in LibraryViewModel:
   - Results are processed and made available through LiveData
   - Error handling provides fallback behavior when filter options can't be retrieved

## UI Implementation
The app processes the response data in these components:

1. `LibraryViewModel`:
   - Calls `getFilterOption(boolean isFromSDCard, String serialNumber)`
   - Transforms response into UI data structures through LiveData
   - Handles errors when filter options can't be retrieved

2. `FilterActivity`:
   - Renders filter UI based on the TagBean response
   - Displays checkboxes/selectors for each category of filters
   - Allows users to select multiple filter options

## Error Handling
The application handles errors by:

1. Using a subscriber pattern to capture and process errors:
```java
@Override
public void doOnError(int code, String msg) {
    callBack.onError(code, msg);
}
```

2. The ViewModel provides error state management:
```java
@Override
public void onError(int code, String errorMsg) {
    getMFilterOptionData().setValue(null);
}
```

3. The UI handles missing data gracefully by:
   - Showing default/basic filters when custom options aren't available
   - Providing visual feedback when filter options can't be loaded
   - Allowing retry functionality

## Filter Option Structure
The response provides structured information about filter categories:

1. AI Event Tags:
   - Smart detection events like "Person", "Vehicle", "Animal", etc.
   - May include subtags for more specific filtering

2. Device Event Tags:
   - Device-specific events like "Motion Detection", "Sound Detection"
   - Hardware-related events that trigger recordings

3. Operation Options:
   - Additional filtering options like "Marked", "Unread", etc.
   - Used for user-managed statuses of recordings

4. Devices:
   - List of available devices for multi-device filtering
   - Includes device information and user's role for the device

## Related Endpoints
This endpoint works closely with these related endpoints:

1. `/library/queryVideoSearchOptionBySn` - Alternative endpoint that accepts only the serial number
   - Used when the app only needs to query by SN without SD card option
   - Returns the same response format

2. `/library/newselectlibrary/newevent` - Gets recordings based on filter criteria
   - Uses the filter options selected from the UI built with this endpoint's data
   - Returns the actual recording list based on selected filters

3. `/library/librarystatus` - Gets library recording status
   - Provides information about recording counts and storage
   - Often called before or alongside the filter options endpoint

## Notes
- The endpoint supports both local (SD card) and cloud recordings
- The response structure allows for hierarchical filtering with main categories and subcategories
- Filter options may vary based on device capabilities and firmware
- The FilterActivity implementation allows for saving filter preferences
- Admin status (roleId = 1) may affect available filter options for certain devices