# Upload Feedback Information Endpoint Discovery

## Endpoint Information
- **Path:** `/library/commitQuestionBack`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Submit user feedback about video events and detection quality

## Request Parameters
The endpoint accepts a `LibraryFeedbackEntry` object in the request body:

```json
{
  "traceId": "string",       // Required: The trace ID of the recording being given feedback on
  "userId": integer,         // Required: User ID submitting the feedback
  "codes": [                 // Required: Array of feedback reason codes selected by the user
    integer, integer, ...
  ],
  "remark": "string",        // Optional: Additional user comments/feedback
  "app": {                   // Standard app information (from BaseEntry)
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
The endpoint returns a standard `BaseResponse` object:

```json
{
  "result": 0,               // Result code (0 indicates success)
  "msg": "success"           // Response message
}
```

## Code Analysis
The endpoint is implemented in the API interface:

```java
@POST("/library/commitQuestionBack")
Observable<BaseResponse> commitFeedback(@Body LibraryFeedbackEntry libraryFeedbackEntry);
```

Called through the `VideoAiFeedBackCore` class:

```java
public final void commitVideoFeedback(String traceId, List<Integer> reasonCode, String remark, vicohome_1742553098674_0O0oO0O<String> callBack) {
    vicohome_1742553098674_0O0oOoo.checkNotNullParameter(reasonCode, "reasonCode");
    vicohome_1742553098674_0O0oOoo.checkNotNullParameter(callBack, "callBack");
    
    LibraryFeedbackEntry libraryFeedbackEntry = new LibraryFeedbackEntry();
    libraryFeedbackEntry.setTraceId(traceId);
    libraryFeedbackEntry.setCodes(reasonCode); 
    libraryFeedbackEntry.setRemark(remark);
    libraryFeedbackEntry.setUserId(Integer.valueOf(getInstance().getUserId()));
    
    this.vicohome_1742553098674_0o00OOoOo.add(
        vicohome_1742553098674_0O0o0oo.getInstance()
            .commitFeedback(libraryFeedbackEntry)
            .subscribeOn(Schedulers.io())
            .observeOn(AndroidSchedulers.mainThread())
            .subscribe(new vicohome_1742553098674_0O0oO0o(callBack))
    );
}
```

## User Flow
This endpoint is part of a complete feedback flow:

1. User views a recording in the library and wants to report an issue
2. User taps a "Feedback" or "Report Issue" button
3. App calls `/library/getQuestionBackOptionsV2` to get available feedback options
4. User selects from options (e.g., "False detection", "Missed event", etc.)
5. User optionally adds comments in a text field
6. User submits feedback, and app calls `/library/commitQuestionBack` with selections

The UI implementation is in `LibraryFeedBackActivity`:

```java
private void submitFeedBack() {
    if (this.mTraceId == null) {
        return;
    }
    
    List<Integer> selectedCodes = new ArrayList<>();
    // Collect all selected feedback options from UI
    for (FeedbackOption option : options) {
        if (option.isChecked()) {
            selectedCodes.add(option.getCode());
            // Also add any selected child options
            for (FeedbackOption childOption : option.getChildOptions()) {
                if (childOption.isChecked()) {
                    selectedCodes.add(childOption.getCode());
                }
            }
        }
    }
    
    // Get remark text from input field
    String remark = remarkEditText.getText().toString();
    
    // Submit feedback
    VideoAiFeedBackCore.getInstance().commitVideoFeedback(
        mTraceId, selectedCodes, remark, 
        new FeedbackSubmitCallback(this)
    );
}
```

## Error Handling
Error handling uses the callback pattern:

```java
public void onNext(BaseResponse response) {
    if (response.getResult() < 0) {
        // Error handling
        callBack.onError(response.getResult(), response.getMsg());
    } else {
        // Success handling
        callBack.onSuccess(response.getResult(), response.getMsg(), "");
    }
}

public void onError(Throwable e) {
    callBack.onError(-1, e.toString());
}
```

UI error handling:
```java
@Override
public void onError(int code, String errorMsg) {
    hideLoading();
    showToast(getString(R.string.feedback_failed));
}

@Override
public void onSuccess(int code, String msg, String data) {
    hideLoading();
    showToast(getString(R.string.feedback_success));
    setResult(Activity.RESULT_OK);
    finish();
}
```

## Related Endpoints
This endpoint works together with:
- `/library/getQuestionBackOptionsV2` - Gets available feedback options

## Purpose and Business Logic
This feedback system serves several purposes:

1. **Product Improvement**: Collecting data about detection accuracy
2. **Customer Support**: Identifying issues with specific devices/recordings
3. **Machine Learning**: Gathering data to improve AI detection algorithms
4. **User Satisfaction**: Giving users a channel to report issues

The feedback options are structured hierarchically:
- Top-level categories (e.g., "Detection Issue", "Video Quality")
- Sub-categories under each main category
- Free-form comments for additional context

## Notes
- The `codes` parameter contains integer IDs representing the feedback options selected by the user
- These codes correspond to the options retrieved from `/library/getQuestionBackOptionsV2`
- The `remark` field allows users to provide additional details beyond the pre-defined options
- The submission contains the recording's traceId to associate feedback with a specific event
- This feedback system is specifically for recordings/events, not general app feedback