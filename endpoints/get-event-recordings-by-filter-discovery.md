# Library Event Recordings Endpoint Discovery

## Endpoint Information
- **Path:** `/library/newselectlibrary/newevent`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Get event recordings by applying filters

## Request Parameters
The endpoint takes a `FilterEntry` object which extends `BaseEntry`:

```json
{
  "startTimestamp": 1643673600000,  // Required: Start time for filtering (Unix timestamp in ms)
  "endTimestamp": 1643760000000,    // Required: End time for filtering (Unix timestamp in ms)
  "from": 0,                        // Optional: Pagination start index
  "to": 20,                         // Optional: Pagination end index
  "marked": 0,                      // Optional: Filter by marked status (0=all, 1=marked)
  "missing": 0,                     // Optional: Filter by read status (0=all, 1=unread)
  "serialNumber": ["ABCD1234"],     // Optional: List of device serial numbers to filter by
  "tags": ["motion", "person"],     // Optional: List of event tags to filter by
  "objectNames": ["person", "car"], // Optional: List of detected objects to filter by
  "serialNumberToActivityZone": {   // Optional: Filter by specific activity zones
    "ABCD1234": [1, 2]              // Maps serial number to list of activity zone IDs
  },
  "isFromSDCard": false,            // Optional: Whether to fetch from SD card (local) or cloud
  "videoEventKey": "string",        // Optional: Video event key for specific event retrieval
  "doorbellTags": ["pressed"],      // Optional: Doorbell specific event tags
  "deviceCallEventTag": "string",   // Optional: Device call event tag
  "deviceName": "Front Door",       // Optional: Device name for filtering
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
The endpoint returns a `LibraryRecordEventResponse` object:

```json
{
  "result": 0,                      // Result code (0 indicates success)
  "msg": "string",                  // Response message
  "data": {                         // Response data container
    "eventList": [                  // List of event recordings
      {
        "deviceSn": "ABCD1234",     // Device serial number
        "deviceName": "Front Door", // Device name
        "readStatus": 0,            // Read status (0=read, 1=unread)
        "traceId": "abc123def456",  // Unique trace ID for the event
        "thumbnail": "https://...", // Thumbnail image URL
        "timestamp": 1643673600000, // Event timestamp (Unix ms)
        "eventType": "motion",      // Event type
        "marked": 0,                // Marked status (0=unmarked, 1=marked)
        "tags": ["motion", "person"], // Event tags
        "objectNames": ["person"],  // Detected objects
        "alertZone": 1,             // Activity zone ID that triggered alert
        "recordStatus": 1,          // Recording status code
        "fileSize": 2048000,        // File size in bytes
        "recordUrl": "https://...", // Recording URL
        "recordLen": 15,            // Recording length in seconds
        "videoQuality": "HD"        // Video quality
      }
    ],
    "totalCount": 42,               // Total count of matching recordings
    "hasMore": true                 // Whether more recordings are available
  }
}
```

## Code Analysis
The endpoint is implemented in the API interfaces:

```java
@POST("library/newselectlibrary/newevent")
Observable<LibraryRecordEventResponse> getEventRecordByFilter(@Body FilterEntry filterEntry);
```

Implementation in LibraryApiClient:
```java
@Override
public Observable<LibraryRecordEventResponse> getEventRecordByFilter(FilterEntry entry) {
    checkNotNullParameter(entry, "entry");
    if (entry.isFromSDCard()) {
        return this.localApiService.getEventRecordByFilter(entry);
    }
    NetworkUtils.getInstance().wrapBaseEntry(entry);
    return this.cloudApiService.getEventRecordByFilter(entry);
}
```

Usage example in LibraryViewModel:
```java
public void getEventRecordByFilter(FilterEntry mFilterEntry, RecordFilterCondition recordFilterCondition) {
    checkNotNullParameter(mFilterEntry, "mFilterEntry");
    checkNotNullParameter(recordFilterCondition, "recordFilterCondition");
    
    // Cancel any existing subscription
    Subscription subscription = this.pendingSubscription;
    if (subscription != null) {
        subscription.unsubscribe();
    }
    
    // Make the API call
    this.pendingSubscription = LibraryCore.getInstance().getEventRecordByFilter(
        mFilterEntry, 
        new FilterRecordCallback(recordFilterCondition, mFilterEntry)
    );
}

// Callback handling
class FilterRecordCallback implements ApiCallback<RecordEventResponse> {
    private final RecordFilterCondition condition;
    private final FilterEntry filterEntry;
    
    FilterRecordCallback(RecordFilterCondition condition, FilterEntry filterEntry) {
        this.condition = condition;
        this.filterEntry = filterEntry;
    }
    
    @Override
    public void onSuccess(int code, String msg, RecordEventResponse response) {
        // Update UI with response data
        viewModel.getFilterRecordLiveData().postValue(
            new Triple<>(response, this.condition, Integer.valueOf(code))
        );
        
        // Track SD card usage if applicable
        if (this.filterEntry.isFromSDCard()) {
            viewModel.reportSDReview(code);
        }
    }
    
    @Override
    public void onError(int code, String errorMsg) {
        // Handle error
        viewModel.getFilterRecordLiveData().postValue(
            new Triple<>(null, this.condition, Integer.valueOf(code))
        );
        
        // Track SD card errors if applicable
        if (this.filterEntry.isFromSDCard()) {
            viewModel.reportSDReview(code);
        }
    }
}
```

## Usage Context
This endpoint is used in the following scenarios:

1. In the events/recordings library view:
   - When user opens the events/recordings tab
   - When applying filters to the recordings list
   - When loading more recordings via pagination
   - When switching between cloud and local storage views

2. Filtering recordings:
   - By date range
   - By device
   - By event type (motion, person detection, etc.)
   - By marked/important status
   - By read/unread status
   - By activity zone

3. Display contexts:
   - Main camera feed timeline
   - Events library view
   - Calendar date selection view
   - Event type filtering

## UI Implementation
The app presents this functionality in the following way:

1. In the library fragment:
   - Displays a scrollable list of event recordings
   - Each item shows a thumbnail, timestamp, event type
   - Unread items have special visual indicators
   - Marked items show a star or similar indicator
   - Pagination loads more items when scrolling
   - Pull-to-refresh updates the list

2. Filtering options:
   - Date picker for time range selection
   - Device selector dropdown
   - Event type filter buttons/chips
   - Marked only toggle
   - Unread only toggle
   - Activity zone filter

3. Loading states:
   - Initial loading spinner
   - Pagination loading indicator
   - Error state with retry option
   - Empty state with explanatory text

## Error Handling
The application handles errors by:

1. Network errors:
   - Shows appropriate error UI
   - Logs error events
   - Provides retry option
   - Falls back to cached data when possible

2. Empty results:
   - Shows "No recordings found" UI
   - Suggests modifying filter criteria
   - May offer quick filters to try

3. Permission/access issues:
   - Handles authentication failures
   - Shows subscription upsell if appropriate
   - Explains retention policy limits

4. Source switching:
   - Gracefully handles switching between cloud and SD card
   - Shows appropriate loading indicators during switch
   - Displays different UI elements for different sources

## Data Flow
The data flow for this endpoint is:

1. Request creation:
   - User selects filters in UI
   - App constructs FilterEntry object
   - Pagination parameters are added
   - App determines isFromSDCard based on user selection

2. Request routing:
   - Based on isFromSDCard flag, request goes to:
     - Cloud API (https://api-us.vicohome.io)
     - Local API (https://127.0.0.1:port)

3. Response handling:
   - Data is parsed into model objects
   - UI is updated with results
   - Pagination state is updated
   - Analytics events are logged

## Related Endpoints
This endpoint works with these related endpoints:

1. `/library/newselectlibrary` - Get event details
   - Similar but with different focus
   - Less filtering options

2. `/library/newselectsinglelibrary` - Get library by trace ID
   - Used when selecting a specific recording
   - Gets detailed information about one recording

3. `/library/deletelibrary/` - Delete recording
   - Used to remove recordings from the library
   - Often called after viewing recording details

4. `/library/librarystatus` - Get library status
   - Gets overall status of recordings library
   - Used for showing counts, storage usage, etc.

## Notes
- This endpoint is part of the library management API group
- The isFromSDCard parameter controls whether to fetch from cloud or local storage
- Pagination is handled via from/to parameters
- Different device types may support different filtering options
- The response includes both metadata and content URLs
- Performance is optimized for quickly browsing recordings
- Both mobile and tablet UIs are supported
- Analytics tracks user interaction with the recordings
- Free-tier users may have limited retention periods