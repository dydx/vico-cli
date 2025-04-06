# deleteSleepPlan Endpoint

## Overview
The deleteSleepPlan endpoint allows clients to delete sleep plans (dormancy settings) for a specific device. Sleep plans define when a device enters low-power or sleep mode to conserve energy.

## API Details
- **Path**: `/device/dormancy/delete`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Deletes a sleep plan for a specific device based on the provided serial number and period.

## Request Parameters
The request body should contain a JSON object with the following properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | The serial number of the device for which to delete the sleep plan |
| period | Integer | Yes | Indicates the type of sleep plan to delete (0 or 1) |

Plus standard BaseEntry properties:
- app
- countryNo
- language
- tenantId

## Request Example
```json
{
  "serialNumber": "ABC123456789",
  "period": 0,
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

## Usage Context
This endpoint is typically used in the following scenarios:
- When a user wants to remove an existing sleep plan for a device
- When a device is being reconfigured with new settings
- As part of a device reset process

## Related Endpoints
- `createSleepPlan` - Creates a new sleep plan
- `editSleepPlan` - Modifies an existing sleep plan
- `listSleepPlan` - Lists all sleep plans for a device

## Implementation Notes
The endpoint is implemented in the DeviceSleepPlanCore and DeviceSettingApiClient classes. It sends a request with the serial number and period to delete the corresponding sleep plan. The server processes the deletion and returns a standard response indicating success or failure.