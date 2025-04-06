# getOtaStatue Endpoint

## Overview
The getOtaStatue endpoint retrieves the current status of Over-The-Air (OTA) firmware updates for a specific device. This endpoint allows clients to monitor the progress of firmware updates, including download and installation status.

> Note: There appears to be a typo in the endpoint name ("statue" instead of "status"), but this is consistent throughout the codebase.

## API Details
- **Path**: `/device/otastatus/`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Retrieves the current firmware update status for a specific device.

## Request Parameters
The request body should contain a JSON object with the following properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | The serial number of the device to check update status for |

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
The response contains details about the current state of the device's firmware update:

| Property | Type | Description |
|----------|------|-------------|
| result | Integer | Status code (0 for success, negative values for errors) |
| msg | String | Status message or error description |
| data | Object | Contains OTA update status information |

### Data Object Structure
| Property | Type | Description |
|----------|------|-------------|
| inProgress | Boolean | Whether an update is currently in progress |
| localDataProgress | Integer | Current update progress percentage (0-100) |
| localDataStartTime | Long | Timestamp when the update was initiated |
| otaStatus | Integer | Update status code (see status codes below) |
| serialNumber | String | Device serial number identifier |
| status | Integer | General device status code |
| targetFirmware | String | Target firmware version being updated to |
| totalSize | Long | Total size of the firmware update in bytes |
| transferredSize | Long | Amount of data transferred in bytes |

### OTA Status Codes
| Code | Description |
|------|-------------|
| 0 | Not started/Initial state |
| 1 | In progress (downloading firmware) |
| 2 | Installing firmware |
| 3 | Complete |
| 4 | Error (failure) |
| 5 | Error (another type) |

## Response Example
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "inProgress": true,
    "localDataProgress": 65,
    "localDataStartTime": 1649876543210,
    "otaStatus": 1,
    "serialNumber": "ABC123456789",
    "status": 1,
    "targetFirmware": "1.2.345",
    "totalSize": 15728640,
    "transferredSize": 10223616
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
| -10401 | Update download failed |
| -10402 | Update installation failed |

## Usage Context
This endpoint is typically used in the following scenarios:
- After initiating a firmware update to monitor progress
- Displaying update status and progress bars in the app
- Checking if an interrupted update is still in progress
- Verifying if a device has the latest firmware

## Related Endpoints
- `startOta` - Initiates a firmware update for a device
- `getDeviceAttributes` - Gets device information including current firmware version

## Implementation Notes
The endpoint is designed to be polled at regular intervals (typically every 5 seconds) during an active firmware update. The client can calculate the progress percentage:
- When otaStatus=1 (downloading): progress = (transferredSize * 90) / totalSize
- When otaStatus=2 (installing): progress = 95%
- When otaStatus=3 (complete): progress = 100%

This polling approach allows the client application to display a real-time progress indicator and notify the user when the update is complete or if an error occurs.