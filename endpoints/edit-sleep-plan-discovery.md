# editSleepPlan Endpoint

## Overview
The editSleepPlan endpoint allows clients to update existing sleep plans (dormancy settings) for a specific device. Sleep plans define when a device enters low-power or sleep mode to conserve energy.

## API Details
- **Path**: `/device/dormancy/edit`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Updates an existing sleep plan for a specific device based on the provided parameters.

## Request Parameters
The request body should contain a JSON object with the following properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | The serial number of the device for which to update the sleep plan |
| startHour | Integer | Yes | Start hour in 24-hour format (0-23) |
| startMinute | Integer | Yes | Start minute (0-59) |
| endHour | Integer | Yes | End hour in 24-hour format (0-23) |
| endMinute | Integer | Yes | End minute (0-59) |
| period | Integer | Yes | Period identifier (1 for weekly schedule, 0 for single day) |
| planDay | Integer | No | Single day selection (1-7, used when period=0) |
| planStartDay | Array of Integers | No | List of days of the week (1=Monday, 7=Sunday, used when period=1) |

Plus standard BaseEntry properties:
- app
- countryNo
- language
- tenantId

## Request Example
Weekly schedule (period=1):
```json
{
  "serialNumber": "ABC123456789",
  "startHour": 22,
  "startMinute": 0,
  "endHour": 6,
  "endMinute": 0,
  "period": 1,
  "planStartDay": [1, 2, 3, 4, 5],
  "app": "vicohome",
  "countryNo": "US",
  "language": "en",
  "tenantId": "AAA12345678"
}
```

Single day schedule (period=0):
```json
{
  "serialNumber": "ABC123456789",
  "startHour": 22,
  "startMinute": 0,
  "endHour": 6,
  "endMinute": 0,
  "period": 0,
  "planDay": 1,
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
- When a user wants to modify an existing sleep plan for a device
- When time settings or active days for a sleep plan need to be adjusted
- When changing a sleep plan from daily to weekly or vice versa

## Related Endpoints
- `createSleepPlan` - Creates a new sleep plan
- `deleteSleepPlan` - Deletes an existing sleep plan
- `listSleepPlan` - Lists all sleep plans for a device

## Implementation Notes
The endpoint is implemented in the DeviceSleepPlanCore and DeviceSettingApiClient classes. It uses the same DeviceSleepPlanBean class as the createSleepPlan endpoint. The API is used to modify existing sleep plans for devices, allowing users to change the times, days, or frequency of the sleep schedule.