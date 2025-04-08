# Endpoint: `/device/bindcheck/step`

## Description
This endpoint is used to check the status of device binding process. The client periodically polls this endpoint to determine the current step of the device binding workflow.

## Request

### HTTP Method
POST

### Headers
- Content-Type: application/json

### Body Parameters
```json
{
  "operationId": "string",  // Required: The ID of the binding operation
  "app": {                  // Basic app information
    "appType": "Android",   // Default: "Android"
    "apiVersion": "string", // API version
    "appName": "string",    // Application name
    "bundle": "string",     // Bundle identifier
    "countlyId": "string",  // Analytics ID
    "env": "string",        // Environment
    "tenantId": "string",   // Tenant ID
    "timeZone": "string",   // Device timezone
    "version": 0,           // App version code
    "versionName": "string" // App version name
  },
  "countryNo": "string",    // Optional: Country code
  "language": "string",     // Optional: Language code
  "tenantId": "string"      // Optional: Tenant ID
}
```

## Response

### Success Response (200 OK)
```json
{
  "result": 0,               // Result code (0 = success)
  "msg": "string",           // Message description
  "data": {
    "deviceBindStep": "string", // Current binding step (1-4)
    "opretionId": "string",     // Operation ID (same as request)
    "serialNumber": "string"    // Device serial number
  }
}
```

#### Device Binding Steps
- `"1"`: WIFI_OK - Device connected to WiFi
- `"2"`: DEVICE_OK - Device operational
- `"3"`: REGISTER_CLOUD_OK - Device registered with cloud
- `"4"`: DEVICE_INITED_OK - Device initialization complete

### Error Response
```json
{
  "result": <error_code>,
  "msg": "error description"
}
```

#### Error Codes
- `-3012`: Not ready for binding
- `-2111`: Not ready for binding
- `-9011`: Binding limited (too many attempts)
- `-9`: Binding timeout
- `-1`: Cannot bind

## Usage
- The client should periodically call this endpoint (typically every 6 seconds) to check the binding progress
- The binding process is complete when `deviceBindStep` reaches `"4"` (DEVICE_INITED_OK)
- If the response contains an error code, the client should handle accordingly (retry, show error to user, etc.)

## Examples

### Request Example
```json
{
  "operationId": "abc123xyz",
  "app": {
    "appType": "Android",
    "apiVersion": "1.0",
    "appName": "VicoHome",
    "bundle": "com.ai.guard.vicohome",
    "countlyId": "user123",
    "env": "production",
    "tenantId": "default",
    "timeZone": "America/Los_Angeles",
    "version": 124,
    "versionName": "2.3.1"
  },
  "language": "en"
}
```

### Response Example - In Progress
```json
{
  "result": 0,
  "msg": "Success",
  "data": {
    "deviceBindStep": "2",
    "opretionId": "abc123xyz",
    "serialNumber": "VICO12345678"
  }
}
```

### Response Example - Complete
```json
{
  "result": 0,
  "msg": "Success",
  "data": {
    "deviceBindStep": "4",
    "opretionId": "abc123xyz",
    "serialNumber": "VICO12345678"
  }
}
```

### Response Example - Error
```json
{
  "result": -9,
  "msg": "Binding timeout"
}
```

## Notes
- This endpoint is part of the device binding workflow
- The client must first obtain an `operationId` before using this endpoint
- The client should implement a timeout mechanism to avoid infinite polling
- Typically used in conjunction with other binding endpoints like `/device/bindComplete` which is called after successful binding