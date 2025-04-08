# Modify Device Attributes Endpoint Discovery

## Endpoint Information
- **Path:** `/device/modifyDeviceAttributes`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Modifies attributes and settings of a device identified by its serial number

## Request Parameters
The endpoint takes an `AttributesEntry` object in the request body:

```json
{
  "serialNumber": "string",         // Required: Device serial number
  "modifiableAttributes": [         // Required: Array of attributes to modify
    {
      "name": "string",             // Required: Attribute name
      "value": object               // Required: New value (type depends on attribute)
    }
  ]
}
```

### Common Modifiable Attributes

| Attribute Name | Type | Description |
|----------------|------|-------------|
| `alarmDuration` | String | Duration of alarm in seconds |
| `pirCooldownTime` | String | Device cool down time between motion detections |
| `alarmFlashLightSwitch` | Boolean | Whether to flash light during alarms |
| `recordingAudioSwitch` | Boolean | Toggle for audio recording |
| `pirSensitivity` | String | Motion detection sensitivity level |
| `motionTrackingSwitch` | Boolean | Toggle motion tracking feature |
| `pirSwitch` | Boolean | Toggle motion detection |
| `nightVisionSensitivity` | String | Night vision sensitivity level |
| `nightVisionSwitch` | Boolean | Toggle night vision |
| `videoAntiFlickerFrequency` | String | Anti-flicker frequency setting |
| `nightVisionMode` | String | Night vision mode setting |
| `pirCooldownSwitch` | Boolean | Toggle motion detection cooldown |
| `videoAntiFlickerSwitch` | Boolean | Toggle anti-flicker feature |
| `voiceVolume` | Integer | Voice feedback volume |
| `alarmVolume` | Integer | Alarm volume |
| `motionAlertSwitch` | Boolean | Toggle motion alerts |
| `pirRecordTime` | String | Duration to record after motion detection |
| `recLampSwitch` | Boolean | Toggle recording indicator lamp |

## Response Format
The endpoint returns a `BaseResponse` object:

```json
{
  "result": 0,            // Result code (0 indicates success, negative values indicate errors)
  "msg": "string"         // Response message
}
```

## Code Analysis
The endpoint is implemented in several interfaces:
1. `com.ai.addxnet.vicohome_1742553098674_00O0o0oOO`
2. `com.ai.settingcore.vicohome_1742553098674_00O0o0oOO`

Example call signature from decompiled code:
```java
@POST("device/modifyDeviceAttributes")
Observable<BaseResponse> modifyDeviceAttributes(@Body AttributesEntry attributesEntry);
```

Example implementation for updating attributes:
```java
public final void updateAttribute(String serialNumber, List<AttributesEntry.vicohome_1742553098674_00O0o0oOO> attributes, final Callback callback) {
    AttributesEntry attributesEntry = new AttributesEntry();
    attributesEntry.setSerialNumber(serialNumber);
    attributesEntry.setModifiableAttributes(attributes);
    
    apiClient.modifyDeviceAttributes(attributesEntry)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<BaseResponse>() {
            @Override
            public void onSubscribe(Disposable d) {}

            @Override
            public void onNext(BaseResponse response) {
                if (response.getResult() < 0) {
                    if (callback != null) {
                        callback.onError(response.getResult(), response.getMsg());
                    }
                    return;
                }
                
                if (callback != null) {
                    callback.onSuccess();
                }
            }

            @Override
            public void onError(Throwable e) {
                if (callback != null) {
                    callback.onError(-1, e.getMessage());
                }
            }

            @Override
            public void onComplete() {}
        });
}
```

## Usage Context
This endpoint is used in several contexts:
1. Updating device settings from the settings UI
2. Toggling features like motion detection, night vision, etc.
3. Adjusting sensitivity levels and other parameters
4. Changing alert and notification settings
5. Modifying recording parameters

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by negative result codes with error messages
- Network errors are handled through the RxJava Observable pattern

Common error conditions include:
- Device is offline or unreachable
- Invalid attribute names or values
- Insufficient permissions
- Server-side validation failures

## Related Endpoints
1. `/device/getDeviceAttributes` - Used to retrieve current attributes before modification
2. `/device/listDeviceAttributes` - Lists devices with their attributes

## Notes
- This endpoint is essential for the device settings functionality
- Changes to device attributes are typically cached locally after successful updates
- For binary attributes (on/off switches), boolean values are used
- For enumerated attributes, string values representing the enum are used
- For numeric attributes, values must be within the range defined by the device capabilities
- The application typically validates input parameters before calling this endpoint
- Multiple attributes can be modified in a single request for efficiency