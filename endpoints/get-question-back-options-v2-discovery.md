# Get Library Feedback Options

## Endpoint Information
- **Path**: `/library/getQuestionBackOptionsV2`
- **Method**: POST
- **Description**: Retrieves feedback options for videos/events in the library, allowing users to provide feedback on detection accuracy and other issues.

## Request Parameters

The request takes a `LibraryFeedbackEntry` object with the following properties:

| Field     | Type         | Required | Description |
|-----------|--------------|----------|-------------|
| traceId   | String       | Yes      | The trace ID of the event/recording |
| userId    | Integer      | Yes      | The user's ID |
| codes     | List<Integer>| No       | Pre-selected option codes (when retrieving existing feedback) |
| remark    | String       | No       | Custom feedback text (defaults to empty string) |

## Response Format

The response is a `LibraryFeedbackResponse` object with the following structure:

| Field     | Type          | Description |
|-----------|---------------|-------------|
| result    | Integer       | Result code (0 for success, negative for errors) |
| msg       | String        | Status message |
| data      | FeedbackData  | The feedback data object |

### FeedbackData Object

| Field        | Type                 | Description |
|--------------|----------------------|-------------|
| userId       | Integer              | The user's ID |
| traceId      | String               | The trace ID of the event/recording |
| libraryId    | String               | The library ID |
| serialNumber | String               | The device serial number |
| isFirst      | Boolean              | Whether this is the first feedback for this event |
| checkedCodes | List<Integer>        | Previously selected feedback option codes |
| options      | List<FeedbackOption> | Available feedback options |
| remark       | String               | Custom feedback text |

### FeedbackOption Object

| Field        | Type                 | Description |
|--------------|----------------------|-------------|
| code         | Integer              | Unique identifier for this option |
| title        | String               | Display text for the option |
| checked      | Boolean              | Whether this option is selected |
| hasTag       | Boolean              | Whether this option has an additional tag |
| tag          | String               | Optional tag text |
| childOptions | List<FeedbackOption> | Nested options (for hierarchical feedback) |

## How It's Called in the App

The endpoint is accessed through `VideoAiFeedBackCore` class, which provides two main methods:

1. `fetchVideoFeedbackList(String traceId, callback)` - Retrieves feedback options for a specific event
2. `commitVideoFeedback(String traceId, List<Integer> reasonCode, String remark, callback)` - Submits feedback for an event

### Getting Feedback Options
```java
// Example usage for fetching feedback options
LibraryFeedbackEntry libraryFeedbackEntry = new LibraryFeedbackEntry();
libraryFeedbackEntry.setTraceId(traceId);
// API call
getInstance().getLibraryFeedback(libraryFeedbackEntry)
    .subscribeOn(Schedulers.io())
    .observeOn(AndroidSchedulers.mainThread())
    .subscribe(new Subscriber<LibraryFeedbackResponse>() {
        @Override
        public void onNext(LibraryFeedbackResponse response) {
            // Process response
            if (response.getResult() < 0) {
                // Error handling
                callBack.onError(response.getResult(), response.getMsg());
            } else if (response.getData() != null) {
                // Success handling
                callBack.onSuccess(response.getResult(), response.getMsg(), response.getData());
            }
        }
        
        @Override
        public void onError(Throwable e) {
            callBack.onError(-1, e.toString());
        }
    });
```

### Submitting Feedback
```java
// Example usage for committing feedback
LibraryFeedbackEntry libraryFeedbackEntry = new LibraryFeedbackEntry();
libraryFeedbackEntry.setTraceId(traceId);
libraryFeedbackEntry.setCodes(reasonCodes);
libraryFeedbackEntry.setRemark(remark);
libraryFeedbackEntry.setUserId(Integer.valueOf(getInstance().getUserId()));
// API call
getInstance().uploadFeedbackInfo(libraryFeedbackEntry)
    .subscribeOn(Schedulers.io())
    .observeOn(AndroidSchedulers.mainThread())
    .subscribe(new Subscriber<BaseResponse>() {
        @Override
        public void onNext(BaseResponse response) {
            // Process response
            if (response.getResult() < 0) {
                callBack.onError(response.getResult(), response.getMsg());
            } else {
                callBack.onSuccess(response.getResult(), response.getMsg(), "");
            }
        }
        
        @Override
        public void onError(int code, String msg) {
            callBack.onError(code, msg);
        }
    });
```

## Error Handling

Error handling is implemented through the callback pattern:

1. If the HTTP request fails, the error is passed to the callback with code -1 and the error message.
2. If the server returns a result code less than 0, it's considered an error and passed to the error callback with the server's error message.
3. Success responses include the full data object (for retrieving options) or just a success message (for submitting feedback).

## Sample Request/Response

### Request (Retrieve Options)
```json
{
  "traceId": "cba3d567890",
  "userId": 12345
}
```

### Response (Retrieve Options)
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "userId": 12345,
    "traceId": "cba3d567890",
    "libraryId": "lib123456",
    "serialNumber": "SN12345678",
    "isFirst": true,
    "checkedCodes": [],
    "options": [
      {
        "code": 1,
        "title": "Detection inaccurate",
        "checked": false,
        "hasTag": false,
        "tag": "",
        "childOptions": [
          {
            "code": 11,
            "title": "False alarm",
            "checked": false,
            "hasTag": false,
            "tag": "",
            "childOptions": []
          },
          {
            "code": 12,
            "title": "Missed detection",
            "checked": false,
            "hasTag": false,
            "tag": "",
            "childOptions": []
          }
        ]
      },
      {
        "code": 2,
        "title": "Video quality issue",
        "checked": false,
        "hasTag": false,
        "tag": "",
        "childOptions": []
      }
    ],
    "remark": ""
  }
}
```

### Request (Submit Feedback)
```json
{
  "traceId": "cba3d567890",
  "userId": 12345,
  "codes": [1, 11],
  "remark": "Camera detected motion but it was just a shadow"
}
```

### Response (Submit Feedback)
```json
{
  "result": 0,
  "msg": "Feedback submitted successfully"
}
```