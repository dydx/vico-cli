# createSleepPlan Endpoint Documentation

## Overview
The `createSleepPlan` endpoint creates a new sleep schedule for a device. Sleep plans allow users to define specific time periods during which the device will enter a low-power state, reducing notifications and conserving energy.

## API Details
- **Path**: `/device/dormancy/create`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `DeviceSleepPlanBean` object with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | Device serial number identifier |
| startHour | Integer | Yes | Start hour in 24-hour format (0-23) |
| startMinute | Integer | Yes | Start minute (0-59) |
| endHour | Integer | Yes | End hour in 24-hour format (0-23) |
| endMinute | Integer | Yes | End minute (0-59) |
| period | Integer | Yes | Period identifier (1 for weekly schedule, 0 for single day) |
| planDay | Integer | No | Single day selection (used when period=0) |
| planStartDay | Array[Integer] | No | List of days of the week (1=Monday, 7=Sunday), used when period=1 |

### Example Request Body
```json
{
  "serialNumber": "ABC123XYZ",
  "startHour": 22,
  "startMinute": 0,
  "endHour": 6,
  "endMinute": 0,
  "period": 1,
  "planStartDay": [1, 2, 3, 4, 5]
}
```

## Response Structure
The endpoint returns a standard `BaseResponse` object:

### Base Response Fields
| Field | Type | Description |
|-------|------|-------------|
| result | int | Status code (0 for success, negative values indicate errors) |
| msg | String | Status message or error description |

### Example Success Response
```json
{
  "result": 0,
  "msg": "success"
}
```

### Example Error Response
```json
{
  "result": -2001,
  "msg": "Network error"
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -2001 | Network error |
| -2002 | Access denied |
| -3001 | Invalid sleep plan parameters |
| -3002 | Device does not support sleep plans |

## Usage in Application
The endpoint is called through the `DeviceSleepPlanCore` class through the following sequence:
1. User configures a sleep schedule in the sleep settings UI
2. A `DeviceSleepPlanBean` object is created with the schedule details
3. The `createSleepPlan` method is called on the `DeviceSleepPlanCore` instance
4. The request is processed and a response is returned through callbacks
5. On success, the UI updates to show the new sleep plan
6. On error, an appropriate error message is displayed to the user

## Constraints
- Start time must be before end time unless the schedule spans midnight
- The device must support sleep mode functionality
- Some devices may have limits on the number of sleep plans that can be created
- Days of the week are represented as integers (1=Monday, 7=Sunday)
- Time is specified in the device's local timezone