# Library Status Endpoint Discovery

## Endpoint Information
- **Path:** `/library/librarystatus`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Get library status information including event count by date

## Request Parameters
The endpoint takes a `FilterEntry` object which extends `BaseEntry`:

```json
{
  "startTimestamp": 1643673600000,  // Optional: Start time for filtering (Unix timestamp in ms)
  "endTimestamp": 1643760000000,    // Optional: End time for filtering (Unix timestamp in ms)
  "from": 0,                        // Optional: Pagination start index
  "to": 20,                         // Optional: Pagination end index
  "serialNumber": ["ABCD1234"],     // Optional: List of device serial numbers to filter by
  "isFromSDCard": false,            // Optional: Whether to fetch from SD card (local) or cloud
  "app": {                          // Standard BaseEntry fields
    "apiVersion": "string",         // API version
    "appName": "string",            // Application name
    "appType": "Android",           // Application type
    "bundle": "string",             // Bundle identifier
    "countlyId": "string",          // Analytics ID
    "env": "string",                // Environment (e.g., "prod")
    "tenantId": "string",           // Tenant identifier
    "timeZone": "string",           // User's timezone
    "version": 0,                   // App version number
    "versionName": "string"         // App version name
  },
  "countryNo": "string",            // Country code
  "language": "string",             // User's language preference
  "tenantId": "string"              // User's tenant ID
}
```

## Response Format
The endpoint returns a `LibraryStatusResponse` object:

```json
{
  "result": 0,                      // Result code (0 indicates success)
  "msg": "string",                  // Response message
  "data": {                         // Response data container
    "list": [                       // List of status counts by date
      {
        "count": 5,                 // Number of recordings for this date
        "dateTimestamp": 1643587200 // Date timestamp (Unix seconds, not milliseconds)
      },
      {
        "count": 12,
        "dateTimestamp": 1643673600
      }
    ]
  }
}
```

## Code Analysis
The endpoint is implemented in the API interfaces:

```java
@POST("/library/librarystatus")
Observable<LibraryStatusResponse> getLibraryStatus(@Body FilterEntry filterEntry);
```

Implementation in LibraryApiClient:
```java
@Override
public Observable<LibraryStatusResponse> getLibraryStatus(FilterEntry entry) {
    checkNotNullParameter(entry, "entry");
    if (entry.isFromSDCard()) {
        return this.localApiService.getLibraryStatus(entry);
    }
    NetworkUtils.getInstance().wrapBaseEntry(entry);
    return this.cloudApiService.getLibraryStatus(entry);
}
```

Usage example in LibraryViewModel:
```java
public void getLibraryStatus(FilterEntry filterEntry) {
    checkNotNullParameter(filterEntry, "filterEntry");
    LibraryCore.getInstance().getLibraryStatus(filterEntry, new LibraryStatusCallback(filterEntry));
}

class LibraryStatusCallback implements ApiCallback<LibraryStatusResponse> {
    private final FilterEntry filterEntry;
    
    LibraryStatusCallback(FilterEntry filterEntry) {
        this.filterEntry = filterEntry;
    }
    
    @Override
    public void onSuccess(int code, String msg, LibraryStatusResponse response) {
        // Update UI with response data
        viewModel.getLibraryStatusLiveData().postValue(
            new Pair<>(this.filterEntry, response)
        );
    }
    
    @Override
    public void onError(int code, String errorMsg) {
        // Handle error
        viewModel.getLibraryStatusLiveData().postValue(
            new Pair<>(this.filterEntry, null)
        );
    }
}
```

## Usage Context
This endpoint is used in the following scenarios:

1. In the events/recordings library view:
   - To display recording counts per day (calendar view)
   - To indicate which days have recordings
   - To show recording density by date

2. Library Analytics:
   - To track total recordings in different time periods
   - To generate usage statistics

3. Display contexts:
   - Calendar date selection view
   - Usage statistics
   - Storage status

## UI Implementation
The app presents this functionality in the following way:

1. In the library calendar view:
   - Dates with recordings are highlighted
   - Dates with more recordings may have stronger highlighting
   - The current date range selection shows recording counts

2. In reporting/analytics views:
   - Event count over time graphs
   - Storage usage indicators

## Error Handling
The application handles errors by:

1. Network errors:
   - Shows appropriate error UI
   - Logs error events
   - Provides retry option
   - Falls back to cached data when possible

2. Empty results:
   - Shows "No recordings found" UI
   - Suggests expanding date range

3. Source switching:
   - Gracefully handles switching between cloud and SD card
   - Shows appropriate loading indicators during switch
   - Displays different UI elements for different sources

## Data Flow
The data flow for this endpoint is:

1. Request creation:
   - App determines time range for status request
   - App constructs FilterEntry object with date range
   - App determines isFromSDCard based on user selection

2. Request routing:
   - Based on isFromSDCard flag, request goes to:
     - Cloud API (https://api-us.vicohome.io)
     - Local API (https://127.0.0.1:port)

3. Response handling:
   - Data is parsed into model objects
   - UI calendar is updated with recording indicators
   - Analytics data is updated

## Related Endpoints
This endpoint works with these related endpoints:

1. `/library/newselectlibrary/newevent` - Get event recordings by filter
   - Used to fetch the actual recordings once the user selects a date
   - Provides detailed event information

2. `/library/newselectlibrary` - Get event details
   - Used to get more information about specific events
   - Called after user selects an event from the list

3. `/library/queryVideoSearchOption` - Query video search options
   - Used to get available filtering options
   - Often called before setting up filter criteria

## Notes
- This endpoint provides date-based statistics for library recordings
- It supports both cloud and local storage sources
- The response includes counts by date to populate calendar views
- May be used to determine date range limits for event queries
- Simple response format focuses on count per date
- Performance is optimized for quick status overview
- Used primarily for calendar views in the library interface