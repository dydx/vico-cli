# List User Devices Endpoint Discovery

## Endpoint Information
- **Path:** `/device/listuserdevices`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Lists all devices associated with a user's account

## Request Parameters
The endpoint takes a `BaseEntry` object in the request body which contains:

```json
{
  "app": {
    "type": "string",     // Application type identifier
    "version": "string"   // Application version
  },
  "countryNo": "string",  // User's country code
  "language": "string",   // User's language preference
  "tenantId": "string"    // User's tenant ID
}
```

## Response Format
The endpoint returns an `AllDeviceResponse` which extends the `BaseResponse` object:

```json
{
  "code": "string",       // Response code (typically "0" for success)
  "msg": "string",        // Response message
  "data": {
    "devices": [
      {
        "serialNumber": "string",    // Device serial number
        "deviceName": "string",      // Name of the device
        "online": 0,                 // Online status (1 = online, 0 = offline)
        "deviceModel": "string",     // Model information
        "batteryLevel": 0,           // Battery level (percentage)
        "sdCard": {                  // SD card information
          "available": true,
          "totalSize": 0,
          "usedSize": 0
        },
        "thumbImgUrl": "string",     // Device thumbnail image URL
        "bindTime": 0,               // Time when device was bound to account
        "location": {                // Location information
          "locationId": "string",
          "locationName": "string"
        },
        "capabilities": [            // Device capabilities
          "string"
        ],
        "settings": {                // Device settings
          "setting1": "value1",
          "setting2": "value2"
        }
        // Additional device attributes...
      }
    ]
  }
}
```

## Code Analysis
The endpoint is called in two different interfaces:
1. `com.ai.addxnet.vicohome_1742553098674_00O0o0oOO` (line 467)
2. `com.ai.settingcore.vicohome_1742553098674_00O0o0oOO` (line 92)

Example call signature from decompiled code:
```java
@POST("/device/listuserdevices")
Observable<AllDeviceResponse> listUserDevices(@Body BaseEntry baseEntry);
```

## Usage Context
This is a core API endpoint used by the application to retrieve a list of all devices associated with the user's account. It's called when:
- The user views their device list in the app
- The application needs to update device information
- When checking device status during app startup

## Error Handling
The response extends BaseResponse which contains:
- `code`: Response code string ("0" indicates success)
- `msg`: A message describing the result or error

Error conditions might include:
- Authentication failures
- Network connectivity issues
- Server-side errors

## Notes
- This is a fundamental endpoint for the VicoHome application as it provides the list of all devices and their current status
- The endpoint requires authentication (likely via a token or cookie that's added by the HTTP client)
- The response data is used to populate the main device list view in the application