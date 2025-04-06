# modifyDeviceAttributes Endpoint

## Overview
The modifyDeviceAttributes endpoint allows clients to update configurable settings for a specific device. This endpoint provides a way to modify various device attributes such as motion detection settings, alarm configurations, night vision settings, and other adjustable parameters.

## API Details
- **Path**: `/device/modifyDeviceAttributes`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Updates one or more configurable attributes for a specific device.

## Request Parameters
The request body should contain an AttributesEntry JSON object with the following properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | The serial number of the device to modify |
| modifiableAttributes | Array | Yes | List of attribute objects to modify |

### Attribute Object Structure
Each object in the modifiableAttributes array has:

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| name | String | Yes | Name of the attribute to modify |
| value | Varies | Yes | New value for the attribute |

Plus standard BaseEntry properties:
- app
- countryNo
- language
- tenantId

## Request Example
```json
{
  "serialNumber": "ABC123456789",
  "modifiableAttributes": [
    {
      "name": "motionTrackingSwitch",
      "value": true
    },
    {
      "name": "alarmVolume",
      "value": 75
    },
    {
      "name": "nightVisionMode",
      "value": 2
    }
  ],
  "app": "vicohome",
  "countryNo": "US",
  "language": "en",
  "tenantId": "AAA12345678"
}
```

## Response Structure
The response is a standard BaseResponse structure:

| Property | Type | Description |
|----------|------|-------------|
| result | Integer | Status code (0 for success, negative values for errors) |
| msg | String | Status message or error description |

## Response Example
```json
{
  "result": 0,
  "msg": "success"
}
```

## Common Modifiable Attributes

| Attribute Name | Type | Description |
|----------------|------|-------------|
| alarmDuration | Integer | Alarm duration in seconds |
| pirCooldownTime | Integer | Motion detection cooldown time in seconds |
| alarmFlashLightSwitch | Boolean | Enable/disable alarm flash light |
| pirSensitivity | Integer | Motion detection sensitivity (typically 1-5) |
| motionTrackingSwitch | Boolean | Enable/disable motion tracking |
| pirSwitch | Boolean | Enable/disable motion detection (PIR) |
| nightVisionSensitivity | Integer | Night vision sensitivity |
| nightVisionSwitch | Boolean | Enable/disable night vision |
| videoAntiFlickerFrequency | Integer | Anti-flicker frequency setting (typically 50/60 Hz) |
| nightVisionMode | Integer | Night vision mode (0=auto, 1=on, 2=off) |
| pirCooldownSwitch | Boolean | Enable/disable motion cooldown period |
| videoAntiFlickerSwitch | Boolean | Enable/disable video anti-flicker |
| voiceVolume | Integer | Voice/audio volume (typically 0-100) |
| alarmVolume | Integer | Alarm volume (typically 0-100) |
| motionAlertSwitch | Boolean | Enable/disable motion alerts |
| pirRecordTime | Integer | Recording time for motion events in seconds |
| recLampSwitch | Boolean | Enable/disable recording indicator lamp |

## Error Codes
| Error Code | Description |
|------------|-------------|
| -1 | General error |
| -2 | Invalid parameters |
| -100 | Device not found |
| -101 | Device offline |
| -102 | Invalid attribute name |
| -103 | Invalid attribute value |

## Usage Context
This endpoint is typically used in the following scenarios:
- Adjusting device settings from the app settings screen
- Configuring detection sensitivity and alert preferences
- Enabling or disabling specific device features
- Adjusting audio volumes for alerts and notifications
- Configuring night vision and video recording settings

## Related Endpoints
- `getDeviceAttributes` - Gets current device attributes including modifiable ones
- `updateEventObjectSwitch` - Updates object detection switch settings
- `updateFloodlightSwitch` - Controls floodlight on/off state

## Implementation Notes
The endpoint accepts multiple attribute changes in a single request, allowing for efficient batch updates of device settings. The attributes that can be modified depend on the device model and its capabilities. Not all attributes are applicable to all device types, and attempting to modify an unsupported attribute may result in an error.

When updating attributes, the system performs validation to ensure the new values are within acceptable ranges. For example, volume settings typically must be between 0 and 100, while mode selectors must match predefined enum values. The implementation in the DeviceSettingCore class provides convenience methods for updating individual attributes using the updateAttribute method, which internally calls this endpoint.