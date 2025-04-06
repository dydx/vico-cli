# unbindDevice Endpoint

## Overview
The unbindDevice endpoint removes a device from a user's account. This process, also known as "unbinding," disconnects the association between the device and the user account, allowing the device to be reset or added to another account.

## API Details
- **Path**: `/device/deactivatedevice`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Removes a device from a user's account, breaking the account-device association.

## Request Parameters
The request body should contain a JSON object with the following properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | The serial number of the device to unbind |

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
| -102 | Permission denied (not device owner) |
| -103 | Device in use |

## Usage Context
This endpoint is typically used in the following scenarios:
- When a user wants to remove a device from their account
- When selling or giving away a device to another user
- When factory resetting a device
- During troubleshooting device connection issues
- When transferring device ownership

## Related Endpoints
- `listDevice` - Lists all devices associated with an account
- `getSingleDevice` - Gets detailed information about a specific device

## Implementation Notes
The endpoint is implemented in the DeviceManageCore class with the deleteDevice and unbindDevice methods. The implementation follows a standard pattern with asynchronous API calls and callbacks for handling successful unbinding or errors.

The unbinding process includes several steps:
1. Validating that the user has permission to unbind the device
2. Sending the unbind request to the server
3. Removing the device from the user's account
4. Updating the local device cache to reflect the removal

In the application code, callbacks are defined for:
- onStartUnBind - Called when the unbinding process starts
- onUnBindSuccess - Called when the device is successfully unbound
- onUnbindError - Called when an error occurs during the unbinding process

Once a device is unbound, it typically enters a state where it can be set up again and added to the same or a different user account.