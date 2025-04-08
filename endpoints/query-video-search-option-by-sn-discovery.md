# Query Video Search Options By Serial Number

This endpoint allows retrieving video search filtering options for a specific device identified by its serial number.

## Endpoint

- **URL**: `/library/queryVideoSearchOptionBySn`
- **Method**: `POST`
- **State**: Unverified
- **Group**: library

## Description

This endpoint retrieves available filtering options for video/event records specific to a camera device identified by its serial number. It returns categories of filtering options including AI event tags, device event tags, and operation options that can be used to filter library recordings.

## Request Parameters

The request is a JSON object with the following properties:

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| sn | String | Yes | The serial number of the device |

Plus standard `BaseEntry` properties that are handled by the library.

## Response

The endpoint returns a JSON response with the following structure:

```json
{
  "result": 0,          // 0 indicates success
  "msg": "success",     // Status message
  "data": {             // TagBean object containing filtering options
    "aiEventTags": [    // AI-detected event tags
      {
        "name": "string",
        "subTags": [] 
      }
    ],
    "deviceEventTags": [ // Device-triggered event tags
      {
        "name": "string",
        "subTags": []
      }
    ],
    "operateOptions": [  // Operation filtering options
      {
        "name": "string",
        "subTags": []
      }
    ],
    "devices": [         // Available devices
      {
        "serialNumber": "string",
        "deviceName": "string",
        "modelCategory": 0,
        "roleId": 0,
        "isBind": false
      }
    ]
  }
}
```

### Data Objects

#### TagBean
- `aiEventTags`: List of AI detection event tags (person, vehicle, etc.)
- `deviceEventTags`: List of device-triggered event tags (motion, etc.)
- `operateOptions`: List of user operation options
- `devices`: List of available devices

#### OptionTag
- `name`: The display name of the tag
- `subTags`: Optional list of sub-category tags

#### OptionDevice
- `serialNumber`: Device serial number
- `deviceName`: User-assigned device name
- `modelCategory`: Device model category identifier
- `roleId`: User role ID (1 = admin)
- `isBind`: Whether the device is bound

## Example Usage

From the app code, this endpoint is used when retrieving filtering options for the video library screen. It's specifically called in the `queryVideoTagsByCloud` method of `LibraryCore.java`:

```java
public final void queryVideoTagsByCloud(String serialNumber, com.ai.addxbase.vicohome_1742553098674_0O0oO0O<TagBean> callBack) {
    vicohome_1742553098674_0O0oOoo.checkNotNullParameter(serialNumber, "serialNumber");
    vicohome_1742553098674_0O0oOoo.checkNotNullParameter(callBack, "callBack");
    this.vicohome_1742553098674_0o00OOoOo.add(LibraryApiClient.vicohome_1742553098674_o00OOoo.getSInstance().queryVideoSearchOptionBySn(new LibraryOptionRequestBySn(serialNumber)).subscribeOn(Schedulers.io()).observeOn(AndroidSchedulers.mainThread()).subscribe((Subscriber<? super VideoSearchOptionResponse>) new vicohome_1742553098674_0O0oO(callBack)));
}
```

## Error Handling

Errors are handled by the `doOnError` method in the callback classes. In `LibraryCore.java`, the method `queryVideoTagsByCloud` handles errors by:

1. Checking if parameters are null
2. Processing network errors in the response handler
3. Passing any error code and message to the callback through the `onError` method

The client is expected to handle error responses where the `result` field is non-zero.

## Notes

This endpoint is specifically used to retrieve filtering options for a single device by serial number, while the related endpoint `/library/queryVideoSearchOption` can accept additional parameters.