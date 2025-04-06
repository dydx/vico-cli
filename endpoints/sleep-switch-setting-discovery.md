# sleepSwitchSetting Endpoint

## Overview
The sleepSwitchSetting endpoint allows clients to enable or disable sleep mode for a specific device. Unlike the sleep plan endpoints that manage schedules, this endpoint provides immediate control over a device's sleep/dormancy state.

## API Details
- **Path**: `/device/dormancy/switch`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Toggles sleep mode on or off for a specific device.

## Request Parameters
The request body should contain a JSON object with the following properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | The serial number of the device |
| dormancySwitch | Integer | Yes | Sleep mode toggle (0 = off, 1 = on) |

Plus standard BaseEntry properties:
- app
- countryNo
- language
- tenantId

## Request Example
```json
{
  "serialNumber": "ABC123456789",
  "dormancySwitch": 1,
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
| -102 | Device does not support sleep mode |

## Usage Context
This endpoint is typically used in the following scenarios:
- When a user wants to manually enable sleep mode to conserve battery
- When a user wants to disable sleep mode to ensure the device remains active
- During device setup or reconfiguration to set initial sleep settings
- When testing device functionality with sleep mode on or off
- To override scheduled sleep plans temporarily without changing the configured schedules

## Related Endpoints
- `createSleepPlan` - Creates a new sleep schedule
- `editSleepPlan` - Modifies an existing sleep schedule
- `deleteSleepPlan` - Deletes a sleep schedule
- `listSleepPlan` - Lists all sleep schedules for a device

## Implementation Notes
The endpoint is implemented in the DeviceSleepPlanCore class with the setSleep method. It provides a direct way to toggle a device's sleep state, which is separate from the sleep plan scheduling system. When sleep mode is enabled (dormancySwitch = 1), the device may enter a low-power state, potentially reducing its responsiveness or available features to conserve battery. When disabled (dormancySwitch = 0), the device remains fully active.

This endpoint affects the device's immediate state, while the sleep plan endpoints manage scheduled sleep periods. Both systems work together to provide comprehensive power management for devices.