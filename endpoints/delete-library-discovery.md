# Delete Recording Discovery

## Endpoint
- **Path**: `/library/deletelibrary/`
- **Method**: POST
- **Description**: Delete one or more video recordings from cloud storage

## Request Parameters
The request uses a `DeleteRecordEntry` object with the following properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| traceIdList | List<String> | Yes | List of recording IDs (trace IDs) to delete |

Example request body:
```json
{
  "traceIdList": ["trace123", "trace456", "trace789"]
}
```

## Response Format
The response is a `LibraryDeleteRecordResponse` object which extends `BaseResponse`:

| Parameter | Type | Description |
|-----------|------|-------------|
| result | int | Result code (0 for success) |
| msg | String | Message describing the result |
| data | Object | Data object containing deletion details |

Example response body:
```json
{
  "result": 0,
  "msg": "Success",
  "data": {
    // Deletion details (not fully specified in the code)
  }
}
```

## Implementation
In the app, deletion follows this process:
1. The user selects one or more recordings to delete in the LibraryFragment
2. The `deleteLibraryVideo()` method is called from the HomeActivity
3. A confirmation dialog is shown to the user
4. If confirmed, the LibraryFragment collects the trace IDs of selected recordings
5. The LibraryViewModel's `deleteRecord()` method is called with the list of trace IDs
6. The request is processed by LibraryCore which calls the API endpoint through LibraryApiClient
7. The response is observed and handled through a callback system

## Error Handling
Error handling occurs at multiple levels:
1. API level errors are returned with non-zero result codes
2. Network errors are caught by the RxJava subscription system
3. Errors are propagated through the callback mechanism to the ViewModel
4. The ViewModel updates its LiveData which notifies observers in the UI

## Notes
- Trace IDs seem to be unique identifiers for recordings in the system
- Multiple recordings can be deleted in a single request
- The app will refresh the library list after deletion is complete
- This endpoint does not delete recordings from local storage, only from cloud

## Sample Code
```java
// From LibraryViewModel.java
public final void deleteRecord(List<String> entry) {
    vicohome_1742553098674_0O0oOoo.checkNotNullParameter(entry, "entry");
    LibraryCore.vicohome_1742553098674_o00OOoo.getSInstance().deleteRecord(entry, new vicohome_1742553098674_0O0o0oo());
}

// From LibraryCore.java
public final void deleteRecord(List<String> traceIdList, com.ai.addxbase.vicohome_1742553098674_0O0oO0O<vicohome_1742553098674_oO00ooOO.vicohome_1742553098674_00O0o0oOO> callBack) {
    vicohome_1742553098674_0O0oOoo.checkNotNullParameter(traceIdList, "traceIdList");
    vicohome_1742553098674_0O0oOoo.checkNotNullParameter(callBack, "callBack");
    LibraryApiClient sInstance = LibraryApiClient.vicohome_1742553098674_o00OOoo.getSInstance();
    DeleteRecordEntry deleteRecordEntry = new DeleteRecordEntry();
    deleteRecordEntry.setTraceIdList(traceIdList);
    this.vicohome_1742553098674_0o00OOoOo.add(sInstance.deleteRecord(deleteRecordEntry).subscribeOn(Schedulers.io()).observeOn(AndroidSchedulers.mainThread()).subscribe((Subscriber<? super LibraryDeleteRecordResponse>) new vicohome_1742553098674_0O0o0ooO(callBack)));
}
```