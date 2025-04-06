# loadPersonDetectionDetailInfo Endpoint

## Overview
The loadPersonDetectionDetailInfo endpoint retrieves detailed information about person detection and other AI detection notification settings for a specific device. This endpoint allows clients to query which types of AI detections (people, packages, pets, vehicles, etc.) are configured to trigger notifications.

## API Details
- **Path**: `/device/queryMessageNotification/v1`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Retrieves AI detection notification settings for a specific device.

## Request Parameters
The request body should contain a JSON object with the following properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | The serial number of the device to query |
| filterByAiAnalyze | Boolean | No | Flag to filter events by AI analysis capability |
| personDetect | Integer | No | Person detection setting value |
| userId | Integer | Yes | User identifier |

Plus standard BaseEntry properties:
- app
- countryNo
- language
- tenantId

## Request Example
```json
{
  "serialNumber": "ABC123456789",
  "filterByAiAnalyze": true,
  "personDetect": 1,
  "userId": 12345,
  "app": "vicohome",
  "countryNo": "US",
  "language": "en",
  "tenantId": "AAA12345678"
}
```

## Response Structure
The response contains information about detection categories and their notification settings:

| Property | Type | Description |
|----------|------|-------------|
| result | Integer | Status code (0 for success, negative values for errors) |
| msg | String | Status message or error description |
| data | Array | List of detection category objects |

### Detection Category Object Structure
Each detection category in the list has the following properties:

| Property | Type | Description |
|----------|------|-------------|
| name | String | Name of the detection category (e.g., "Person", "Package", "Pet") |
| choice | Boolean | Flag indicating if this detection is enabled for notifications |
| subEvent | Array | List of sub-event objects for more granular detection categories |

### Sub-Event Object Structure
Each sub-event has the following properties:

| Property | Type | Description |
|----------|------|-------------|
| name | String | Name of the sub-event detection category |
| choice | Boolean | Flag indicating if this sub-event detection is enabled |

## Response Example
```json
{
  "result": 0,
  "msg": "success",
  "data": [
    {
      "name": "Person",
      "choice": true,
      "subEvent": [
        {
          "name": "Familiar Face",
          "choice": true
        },
        {
          "name": "Stranger",
          "choice": true
        }
      ]
    },
    {
      "name": "Package",
      "choice": true,
      "subEvent": []
    },
    {
      "name": "Vehicle",
      "choice": false,
      "subEvent": []
    },
    {
      "name": "Pet",
      "choice": true,
      "subEvent": []
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

## Usage Context
This endpoint is typically used in the following scenarios:
- When configuring notification settings for a camera
- When displaying the current AI detection configuration in settings screens
- When initializing the device settings UI to show which detections are enabled
- Before updating detection settings to get the current configuration first

## Related Endpoints
- `uploadPersonDetectionDetailInfo` - Updates person detection notification settings
- `updateEventObjectSwitch` - Updates object detection switch settings
- `queryAiEventInfo` - Queries AI event information

## Implementation Notes
The endpoint is implemented in the DeviceAICore class with the getMessageNotificationConfig method. It returns a hierarchical structure of detection categories and their enabled/disabled state, allowing for both top-level categories (like "Person") and more specific sub-categories ("Familiar Face" vs "Stranger"). The system supports various AI detection types that can be individually configured to trigger notifications based on user preferences.