# startOta Endpoint

## Overview
The startOta endpoint initiates an Over-The-Air (OTA) firmware update for a specific device. This endpoint starts the update process, after which clients should use the getOtaStatue endpoint to monitor the update progress.

## API Details
- **Path**: `/device/otastart/`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Initiates an OTA firmware update for a specific device.

## Request Parameters
The request body should contain a JSON object with the following properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | The serial number of the device to update |

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
The response contains information about the initial state of the OTA update:

| Property | Type | Description |
|----------|------|-------------|
| result | Integer | Status code (0 for success, negative values for errors) |
| msg | String | Status message or error description |
| data | Object | Contains initial OTA status information |

### Data Object Structure
| Property | Type | Description |
|----------|------|-------------|
| inProgress | Boolean | Indicates if the update is in progress |
| otaStatus | Integer | Status code for the update (0=not started, 1=downloading, etc.) |
| serialNumber | String | Device serial number |
| status | Integer | General device status |
| totalSize | Long | Total size of the firmware update in bytes |
| transferredSize | Long | Amount of data transferred so far (initially 0) |

## Response Example
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "inProgress": true,
    "otaStatus": 0,
    "serialNumber": "ABC123456789",
    "status": 1,
    "totalSize": 15728640,
    "transferredSize": 0
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
| -102 | No update available |
| -103 | Insufficient battery level |
| -10401 | Update download failed |
| -10402 | Update installation failed |

## Usage Context
This endpoint is typically used in the following scenarios:
- When a new firmware version is available for a device
- During scheduled maintenance updates
- When a user manually initiates a firmware update
- When critical security patches need to be applied

## Related Endpoints
- `getOtaStatue` - Checks the status of an ongoing OTA update

## Implementation Notes
The endpoint initiates the OTA update process for a device, which typically consists of several phases:

1. Initial validation (device compatibility, battery level, etc.)
2. Downloading the firmware update package
3. Verifying the downloaded package
4. Installing the update
5. Rebooting the device with the new firmware

After calling startOta, the client application should poll the getOtaStatue endpoint at regular intervals (typically every 5 seconds) to monitor the progress of the update. The otaStatus field in the response indicates the current phase of the update:

- 0: Not started/Initial state
- 1: Downloading firmware
- 2: Installing firmware
- 3: Complete
- 4: Error (failure)
- 5: Error (another type)

When polling getOtaStatue, the client can calculate progress percentage based on the transferredSize and totalSize values when otaStatus=1 (downloading). When otaStatus=2 (installing), progress is typically shown as 95%, and when otaStatus=3 (complete), progress is 100%.