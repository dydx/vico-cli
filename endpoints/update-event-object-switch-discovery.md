# updateEventObjectSwitch Endpoint

## Overview
The updateEventObjectSwitch endpoint allows clients to update AI object detection settings for a specific device. This endpoint controls which types of objects (such as people, packages, pets, vehicles) a device will detect and generate notifications for.

## API Details
- **Path**: `/aiAssist/updateEventObjectSwitch`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Updates AI object detection settings for a specific device.

## Request Parameters
The request body should contain a JSON object with the following properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | The serial number of the device |
| list | Array | Yes | List of event objects to update |

### Event Object Structure
Each object in the list array has the following properties:

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| eventObject | String | Yes | The type of object to detect ("person", "package", "pet", "vehicle") |
| checked | Boolean | Yes | Whether detection for this object type is enabled |

Plus standard BaseEntry properties:
- app
- countryNo
- language
- tenantId

## Request Example
```json
{
  "serialNumber": "ABC123456789",
  "list": [
    {
      "eventObject": "person",
      "checked": true
    },
    {
      "eventObject": "package",
      "checked": false
    },
    {
      "eventObject": "pet",
      "checked": true
    },
    {
      "eventObject": "vehicle",
      "checked": false
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

## Error Codes
| Error Code | Description |
|------------|-------------|
| -1 | General error |
| -2 | Invalid parameters |
| -100 | Device not found |
| -101 | Device offline |
| -102 | Permission denied |
| -103 | Feature not available (requires premium subscription) |

## Usage Context
This endpoint is typically used in the following scenarios:
- When a user enables or disables specific object detection in device settings
- During initial device setup to configure detection preferences
- When customizing notification settings to receive alerts for specific objects
- When modifying AI features based on user preference or subscription level

## Related Endpoints
- `queryAiEventInfo` - Gets current AI object detection settings
- `queryNotificationAiBird` - Gets bird detection notification settings
- `setNotificationAiBird` - Updates bird detection notification settings

## Implementation Notes
The endpoint is implemented in the DeviceAICore class with the updateAnalysisEventConfig method. In the mobile application, this endpoint is typically called from the AISettingActivity when users toggle object detection switches in the UI.

The implementation follows these patterns:
1. The mobile app creates a SwitchAiEventEntry with the device serial number and the list of objects to modify
2. Each object type (person, package, pet, vehicle) is updated with its enabled/disabled state
3. The server processes the request and applies the settings to the device
4. The server returns a success or error response

This feature often requires a premium/VIP subscription to access, as indicated by checks in the AISettingActivity before enabling these settings. The endpoint allows users to customize which object detection capabilities are active on their devices, tailoring their notification experience to their specific needs.