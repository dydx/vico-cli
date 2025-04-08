# Get Device Attributes Endpoint Discovery

## Endpoint Information
- **Path:** `/device/getDeviceAttributes`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Retrieves various attributes and properties of a device identified by its serial number

## Request Parameters
The endpoint takes a device attributes request object:

```json
{
  "serialNumber": "string",          // Required: Device serial number
  "returnFixedAttributes": boolean,  // Optional: Whether to include fixed attributes (default: true)
  "returnRealTimeAttributes": boolean, // Optional: Whether to include real-time attributes (default: true)
  "returnModifiableAttributes": boolean // Optional: Whether to include modifiable attributes (default: true)
}
```

## Response Format
The endpoint returns a comprehensive set of device attributes based on the request flags:

```json
{
  "result": 0,              // Result code (0 indicates success)
  "msg": "string",          // Response message
  "data": {
    "serialNumber": "string",
    "fixedAttributes": {    // Only included if returnFixedAttributes is true
      "canStandby": number,
      "displayModelNo": "string",
      "floodlightLuminanceRange": {
        "min": number,
        "max": number,
        "interval": number
      },
      "iccid": "string",
      "icon": "string",
      "macAddress": "string",
      "modelCategory": number,
      "modelNo": "string",
      "quantityCharge": boolean,
      "roleName": "string",
      "smallIcon": "string",
      "support12Hours": boolean,
      "supportIndoor": boolean,
      "supportJson": {
        // JSON of supported features
      },
      "supportManualFloodlightLuminance": boolean,
      "supportManualFloodlightSwitch": boolean,
      "supportOtaAutoUpgrade": boolean,
      "supportPirAi": boolean,
      "supportPlanFloodlightLuminance": boolean,
      "supportPlanFloodlightSwitch": boolean,
      "supportStarlightSensor": boolean,
      "supportTriggerFloodlightSwitch": boolean,
      "supportWhiteLight": boolean,
      "triggerFloodlightCooldownTimeOptions": ["string"],
      "userSn": "string",
      "wiredMacAddress": "string"
    },
    "realTimeAttributes": { // Only included if returnRealTimeAttributes is true
      "batteryLevel": number,
      "isCharging": number,
      "chargingMode": number,
      "deviceNetType": number,
      "deviceStatus": number,
      "displayGitSha": "string",
      "firmwareId": "string",
      "firmwareStatus": number,
      "iccid": "string",
      "imei": "string",
      "ip": "string",
      "mcuNumber": "string",
      "networkName": "string",
      "newestFirmwareId": "string",
      "offlineTime": number,
      "online": number,
      "sdCard": {
        "formatStatus": number,
        "free": number,
        "total": number,
        "used": number
      },
      "signalLevel": number,
      "signalStrength": number,
      "simStatus": number,
      "simThirdParty": number,
      "whiteLight": number,
      "wifiChannel": number
    },
    "modifiableAttributes": [ // Only included if returnModifiableAttributes is true
      {
        "name": "string",      // Attribute name (e.g., "pirSwitch")
        "type": "string",      // Attribute type (e.g., "boolean", "int", "enum")
        "value": object,       // Current value (type depends on "type" field)
        "disabled": boolean,   // Whether this attribute is currently disabled
        "disabledOptions": ["string"], // Options that are disabled if this is an enum
        "intRange": {          // Range info if this is a numeric attribute
          "min": number,
          "max": number,
          "interval": number
        },
        "options": [object]    // Available options if this is an enum type
      }
      // Additional modifiable attributes...
    ]
  }
}
```

## Code Analysis
The endpoint is implemented in several interfaces:
1. `com.ai.addxnet.vicohome_1742553098674_00O0o0oOO`
2. `com.ai.settingcore.vicohome_1742553098674_00O0o0oOO`

Example call signature from decompiled code:
```java
@POST("device/getDeviceAttributes")
Observable<DeviceAttributesResponse> getDeviceAttributes(@Body DeviceAttributesRequest request);
```

The application retrieves device attributes using the `DeviceSettingCore.getDeviceAttributes()` method:
```java
public final void getDeviceAttributes(String serialNumber, Boolean returnFixed, Boolean returnRealTime, Boolean returnModifiable, DeviceAttributesCallback callback) {
    DeviceAttributesRequest request = new DeviceAttributesRequest();
    request.setSerialNumber(serialNumber);
    request.setReturnFixedAttributes(returnFixed);
    request.setReturnRealTimeAttributes(returnRealTime);
    request.setReturnModifiableAttributes(returnModifiable);
    
    apiClient.getDeviceAttributes(request)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<DeviceAttributesResponse>() {
            // Handle response and callbacks
        });
}
```

## Common Modifiable Attributes
The device attributes endpoint returns a list of attributes that can be modified. Common modifiable attributes include:

- `pirSwitch` - Motion detection on/off
- `motionTrackingSwitch` - Motion tracking on/off
- `pirSensitivity` - Motion sensitivity level
- `alarmFlashLightSwitch` - Whether to flash light during alarm
- `nightVisionMode` - Night vision operation mode
- `nightVisionSwitch` - Night vision on/off
- `nightVisionSensitivity` - Sensitivity for night vision
- `videoAntiFlickerFrequency` - Anti-flicker frequency (50Hz/60Hz)
- `videoAntiFlickerSwitch` - Anti-flicker on/off
- `pirCooldownTime` - Time between motion detections
- `pirCooldownSwitch` - Cool down period on/off
- `alarmDuration` - Duration of alarm
- `pirRecordTime` - Duration of motion recording
- `voiceVolume` - Volume level for voice
- `alarmVolume` - Volume level for alarm
- `recLampSwitch` - Recording indicator light on/off
- `motionAlertSwitch` - Motion alert notifications on/off

## Usage Context
This endpoint is used in several contexts:
1. When loading the device settings page to show all available settings
2. Before modifying device attributes to determine current values
3. When checking device capabilities to enable/disable features in the UI
4. Retrieving real-time status information like battery level and connectivity

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by non-zero result codes with error messages
- Network errors are handled through the RxJava Observable pattern

## Related Endpoints
1. `/device/modifyDeviceAttributes` - Used to update modifiable attributes after retrieving them
2. `/device/listDeviceAttributes` - Lists all devices with their attributes in a more summary format

## Notes
- This endpoint provides comprehensive device information and is essential for the device settings UI
- The flexible request format allows for optimizing payloads by only requesting needed attribute types
- Fixed attributes rarely change and could be cached locally
- Real-time attributes should be refreshed periodically for up-to-date status
- Modifiable attributes form the basis of the device settings UI