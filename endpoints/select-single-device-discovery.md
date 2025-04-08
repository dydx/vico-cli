# Select Single Device Endpoint Discovery

## Endpoint Information
- **Path:** `/device/selectsingledevice`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Retrieves detailed information about a single device by its serial number

## Request Parameters
The endpoint takes a `SerialNoEntry` object in the request body which extends the `BaseEntry`:

```json
{
  "serialNumber": "string",  // Required: The device's serial number
  "app": {                   // From BaseEntry
    "type": "string",        // Application type identifier
    "version": "string"      // Application version
  },
  "countryNo": "string",     // User's country code
  "language": "string",      // User's language preference
  "tenantId": "string",      // User's tenant ID
  "voiceReminder": boolean   // Optional parameter
}
```

## Response Format
The endpoint returns a `GetSingleDeviceResponse` containing detailed device information:

```json
{
  "result": 0,               // 0 indicates success
  "message": "string",       // Response message
  "data": {
    "serialNumber": "string",
    "deviceName": "string",
    "online": 0|1,           // Device online status (0=offline, 1=online)
    "awake": 0|1,            // Device awake status (0=sleep, 1=awake)
    "deviceStatus": integer, // Status code (3=sleep, 11=shutdown-low-power, 12=shutdown-press-key, 13=shutdown-solar-low-power)
    "batteryLevel": integer, // Battery level percentage
    "adminName": "string",
    "adminPhone": "string",
    "adminId": integer,
    "deviceModel": {         // Device model information
      "modelName": "string",
      "canStandby": boolean,
      "isCanRotate": boolean,
      "supportMotionTrack": boolean,
      "whiteLight": boolean,
      "streamProtocol": "string",
      "audioCodecType": "string"
    },
    "sdCard": {              // SD card information if available
      "formatStatus": integer,
      "total": integer,
      "used": integer
    },
    "signalStrength": integer, // WiFi signal strength
    "deviceSupport": {       // Device capabilities
      "deviceSupportAlarm": boolean,
      "deviceSupportMirrorFlip": boolean,
      "deviceDormancySupport": integer,
      "supportWebrtc": integer,
      "supportRecordingAudioToggle": integer,
      "supportLiveAudioToggle": integer
    },
    "thumbImgUrl": "string", // Thumbnail image URL
    "thumbImgTime": long     // Timestamp of thumbnail image
  }
}
```

## Code Analysis
The endpoint is implemented in multiple interfaces:

1. `com.ai.addxnet.vicohome_1742553098674_00O0o0oOO` - Main API interface using RxJava Observable
2. `com.ai.settingcore.vicohome_1742553098674_00O0o0oOO` - Settings-related API interface
3. `com.ai.addxnet.vicohome_1742553098674_0oOOoOO` - Another API interface using Retrofit

Example call signature from decompiled code:
```java
@POST("/device/selectsingledevice")
Observable<GetSingleDeviceResponse> getSingleDevice(@Body SerialNoEntry serialNoEntry);
```

The `com.a4x.player.internal.vicohome_1742553098674_0ooO0Ooo.getDeviceInfo()` method calls this endpoint to retrieve device information for video playback.

## Usage Context
This endpoint is commonly used in the following scenarios:

1. When viewing device details in the device settings page
2. Before initiating a live stream to check device status
3. When checking device connectivity status
4. Before or after device operations to verify state changes

## Error Handling
The endpoint follows standard error handling practices:
- Success is indicated by `result: 0`
- Errors are indicated by non-zero result codes with error messages
- Common error conditions include device not found, authentication issues, and network failures

The implementation includes error handling for:
- Network failures
- Invalid response formats
- Server errors (non-zero result codes)

## Device Status Interpretation
Key device status values include:
- `online`: 0=offline, 1=online
- `awake`: 0=sleep mode, 1=awake
- `deviceStatus`: 
  - 3: device is in sleep mode
  - 11: shutdown due to low power
  - 12: shutdown by pressing power key
  - 13: shutdown due to solar low power

## Notes
- This endpoint is critical for the device details page and before starting any device-specific operations
- The response provides comprehensive device status information including battery level, connectivity, and capabilities
- Device capabilities are used by the app to conditionally enable/disable features in the UI