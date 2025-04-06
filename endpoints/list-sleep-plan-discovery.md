# listSleepPlan Endpoint

## Overview
The listSleepPlan endpoint retrieves all sleep plans configured for a specific device. Sleep plans define when a device enters low-power or sleep mode to conserve energy.

## API Details
- **Path**: `/device/dormancy/list`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Retrieves all sleep plans for a specific device identified by its serial number.

## Request Parameters
The request body should contain a JSON object with the following properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | The serial number of the device for which to retrieve sleep plans |

Plus standard BaseEntry properties:
- app
- countryNo
- language
- tenantId

## Request Example
```json
{
  "serialNumber": "ABC123456789",
  "app": "vicohome",
  "countryNo": "US",
  "language": "en",
  "tenantId": "AAA12345678"
}
```

## Response Structure
The response contains a list of sleep plans for the specified device:

| Property | Type | Description |
|----------|------|-------------|
| result | Integer | Status code (0 for success, negative values for errors) |
| msg | String | Status message or error description |
| data | Array | List of DeviceSleepPlanBean objects |

### DeviceSleepPlanBean Structure
Each sleep plan in the list has the following properties:

| Property | Type | Description |
|----------|------|-------------|
| serialNumber | String | Device serial number identifier |
| period | Integer | Period identifier (1 for weekly schedule, 0 for single day) |
| startHour | Integer | Start hour in 24-hour format (0-23) |
| startMinute | Integer | Start minute (0-59) |
| endHour | Integer | End hour in 24-hour format (0-23) |
| endMinute | Integer | End minute (0-59) |
| planDay | Integer | Single day selection (1-7, used when period=0) |
| planStartDay | Array | List of days of the week (1=Monday, 7=Sunday, used when period=1) |

## Response Example
```json
{
  "result": 0,
  "msg": "success",
  "data": [
    {
      "serialNumber": "ABC123456789",
      "period": 1,
      "startHour": 22,
      "startMinute": 0,
      "endHour": 6,
      "endMinute": 0,
      "planStartDay": [1, 2, 3, 4, 5]
    },
    {
      "serialNumber": "ABC123456789",
      "period": 0,
      "startHour": 23,
      "startMinute": 30,
      "endHour": 7,
      "endMinute": 30,
      "planDay": 6
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
- When a user opens the sleep plan management screen
- When the application needs to display existing sleep plans
- Before creating or editing sleep plans to check for existing configurations
- When synchronizing device settings between app instances

## Related Endpoints
- `createSleepPlan` - Creates a new sleep plan
- `deleteSleepPlan` - Deletes an existing sleep plan
- `editSleepPlan` - Modifies an existing sleep plan

## Implementation Notes
The endpoint is implemented in the DeviceSleepPlanCore class with the getSleepPlanList method. It returns all sleep plans configured for the specified device, which might include both daily (period=0) and weekly (period=1) sleep schedules. The application then uses this data to populate the sleep plan management UI, allowing users to view, edit, or delete existing sleep plans.