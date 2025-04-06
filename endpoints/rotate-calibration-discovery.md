# rotateCalibration Endpoint

## Overview
The rotateCalibration endpoint handles the calibration of a device's physical orientation or rotation. This is particularly important for cameras that need to have a proper understanding of their orientation to display footage correctly and to make accurate motion detection judgments.

## API Details
- **Path**: `/device/rotate-calibration`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Initiates or checks the status of device rotation calibration.

## Request Parameters
The request body should contain a JSON object with the following properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | The serial number of the device to calibrate |
| needCalibration | Boolean | Yes | When true, initiates calibration; when false, checks status |

Plus standard BaseEntry properties:
- app
- countryNo
- language
- tenantId

## Request Example
Starting calibration:
```json
{
  "serialNumber": "ABC123456789",
  "needCalibration": true,
  "app": "vicohome",
  "countryNo": "US",
  "language": "en",
  "tenantId": "AAA12345678"
}
```

Checking calibration status:
```json
{
  "serialNumber": "ABC123456789",
  "needCalibration": false,
  "app": "vicohome",
  "countryNo": "US",
  "language": "en",
  "tenantId": "AAA12345678"
}
```

## Response Structure
The response contains information about the calibration status:

| Property | Type | Description |
|----------|------|-------------|
| result | Integer | Status code (0 for success, negative values for errors) |
| msg | String | Status message or error description |
| data | Object | Contains calibration information |

### Data Object Structure
| Property | Type | Description |
|----------|------|-------------|
| calibrationFinished | Boolean | Whether the calibration process is complete |

## Response Example
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "calibrationFinished": false
  }
}
```

## Error Codes
| Error Code | Description |
|------------|-------------|
| -1 | General error |
| -2 | Invalid parameters |
| -100 | Device not found |
| -101 | Device offline |
| -102 | Device does not support rotation calibration |

## Usage Context
This endpoint is typically used in the following scenarios:
- During initial device setup to ensure proper orientation
- After physically relocating or adjusting a camera's position
- When video footage appears to have incorrect orientation
- As part of troubleshooting camera display or motion detection issues

## Related Endpoints
- `getDeviceAttributes` - Gets device capabilities including rotation support
- `modifyDeviceAttributes` - Updates various device settings

## Implementation Notes
The endpoint is implemented in the DeviceSettingCore class with two separate methods:
- startRotateCalibration: Initiates the calibration process when needCalibration is true
- getRotateCalibrationStatus: Checks if calibration is complete when needCalibration is false

The application includes a dedicated UI component (RotateCalibrationFragment) that guides the user through the calibration process. The typical workflow involves:

1. User initiates calibration through the UI
2. App calls the endpoint with needCalibration=true to start the process
3. The device performs internal orientation calibration
4. App periodically polls with needCalibration=false to check completion status
5. When the response returns calibrationFinished=true, the process is complete

Not all device models support rotation calibration. The application should check device capabilities before offering this feature to users.