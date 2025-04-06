# queryAiEventInfo Endpoint

## Overview
The queryAiEventInfo endpoint retrieves information about AI event detection settings for one or more devices. This endpoint allows clients to query which types of AI events (such as person, package, pet, vehicle) are enabled or disabled for detection on specific devices.

## API Details
- **Path**: `/aiAssist/queryEventObjectSwitch`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Retrieves AI event detection configuration for specified devices.

## Request Parameters
The request body should contain a JSON object with the following properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumbers | Array | Yes* | List of device serial numbers to query |
| isAll | Boolean | No | Flag to query all devices (if true, serialNumbers can be omitted) |

*Required unless isAll is true

Plus standard BaseEntry properties:
- app
- countryNo
- language
- tenantId

## Request Example
```json
{
  "serialNumbers": ["ABC123456789", "XYZ987654321"],
  "isAll": false,
  "app": "vicohome",
  "countryNo": "US",
  "language": "en",
  "tenantId": "AAA12345678"
}
```

## Response Structure
The response contains information about AI event detection settings for each device:

| Property | Type | Description |
|----------|------|-------------|
| result | Integer | Status code (0 for success, negative values for errors) |
| msg | String | Status message or error description |
| data | Array | List of device objects with their AI event settings |

### Device Object Structure
Each device in the data array has the following properties:

| Property | Type | Description |
|----------|------|-------------|
| serialNumber | String | The device's serial number |
| deviceName | String | Name of the device |
| list | Array | List of AI event objects |

### AI Event Object Structure
Each AI event object in the list array has the following properties:

| Property | Type | Description |
|----------|------|-------------|
| eventObject | String | Type of AI event ("person", "package", "pet", "vehicle", etc.) |
| checked | Boolean | Whether detection for this event type is enabled |

## Response Example
```json
{
  "result": 0,
  "msg": "success",
  "data": [
    {
      "serialNumber": "ABC123456789",
      "deviceName": "Front Door Camera",
      "list": [
        {
          "eventObject": "person",
          "checked": true
        },
        {
          "eventObject": "package",
          "checked": true
        },
        {
          "eventObject": "pet",
          "checked": false
        },
        {
          "eventObject": "vehicle",
          "checked": true
        }
      ]
    },
    {
      "serialNumber": "XYZ987654321",
      "deviceName": "Backyard Camera",
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
      ]
    }
  ]
}
```

## Error Codes
| Error Code | Description |
|------------|-------------|
| -1 | General error |
| -2 | Invalid parameters |
| -100 | Device not found |
| -101 | Device offline |

## Usage Context
This endpoint is typically used in the following scenarios:
- Displaying AI detection settings in device configuration screens
- Initializing the AI settings UI to show which detections are enabled
- Before updating AI detection settings to get the current configuration
- When checking compatibility of a device with various AI detection types

## Related Endpoints
- `updateEventObjectSwitch` - Updates AI event detection settings
- `loadPersonDetectionDetailInfo` - Gets more detailed person detection settings
- `uploadPersonDetectionDetailInfo` - Updates person detection settings

## Implementation Notes
The endpoint is implemented in the DeviceAICore class with the getAnalysisEventConfig method. It interacts with the API service via the DeviceSettingApiClient class. The endpoint returns a list of AI event detections that are enabled or disabled for each requested device. This allows the application to display the current state of AI event detection settings and provide UI controls for modifying these settings.

Different device models may support different types of AI event detection capabilities. The response will only include the event types that are supported by each device.