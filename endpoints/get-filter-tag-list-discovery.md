# Get Filter Tag List API Discovery

## Endpoint Information
- **Path**: `/library/othertagnamelist`
- **Method**: POST
- **Group**: library
- **Description**: Get filter tag list for library events

## Request Parameters
The endpoint uses a BaseEntry request object that contains basic authentication information. No additional parameters specific to this endpoint are required.

```java
// Example from code
Observable<FilterTagListResponse> getFilterTagList(@vicohome_1742553098674_00O0ooOO0 BaseEntry baseEntry);
```

## Response Structure
The response is a `FilterTagListResponse` that extends `BaseResponse` and contains a list of tag names as strings.

```java
public final class FilterTagListResponse extends BaseResponse {
    private vicohome_1742553098674_00O0o0oOO data;

    public static final class vicohome_1742553098674_00O0o0oOO implements JsonBean {
        private List<String> otherTagNameList;

        public final List<String> getOtherTagNameList() {
            return this.otherTagNameList;
        }

        public final void setOtherTagNameList(List<String> list) {
            this.otherTagNameList = list;
        }
    }

    public final vicohome_1742553098674_00O0o0oOO getData() {
        return this.data;
    }

    public final void setData(vicohome_1742553098674_00O0o0oOO data) {
        this.data = data;
    }
}
```

## Usage Context
This endpoint is used in the library/recordings filtering functionality of the app. It retrieves a list of tag names that can be used for filtering recorded events. The UI implementation shows these tags are displayed in the FilterActivity, where users can select different filters to narrow down their recordings.

The tags are used in combination with device filters, AI detection filters, and user-set filters to create a comprehensive filtering system for recorded events in the library section of the app.

## Error Handling
The endpoint inherits error handling from BaseResponse, which likely includes standard response fields like a status code and message. No specific error handling unique to this endpoint was observed in the code.

## Additional Notes
- The tags returned by this endpoint include various categories like person detection, package detection, vehicle detection, doorbell events, bird detection, etc.
- These tags are displayed in the UI as selectable filter options that users can apply to filter their video recordings library.
- The implementation in FilterActivity shows how these tags are organized into different groupings like "Video Tag", "Package", "Vehicle" and other categories.