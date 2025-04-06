# closeP2PConnection Endpoint Documentation

## Overview
The `closeP2PConnection` endpoint closes active peer-to-peer (P2P) connections with one or more Vicohome devices. These P2P connections are typically used for real-time video streaming and two-way audio communication with cameras and doorbells.

## API Details
- **Path**: `/dc/rtcconnection/close`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `LibraryP2PConnectRequest` object with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | Array | Yes | List of device serial numbers to close P2P connections for |
| app | Object | Yes | Application information (inherited from BaseEntry) |
| countryNo | String | Yes | Country code (e.g., "US") (inherited from BaseEntry) |
| language | String | Yes | Language code (e.g., "en") (inherited from BaseEntry) |
| tenantId | String | Yes | Tenant identifier (inherited from BaseEntry) |

### Example Request Body
```json
{
  "serialNumber": ["ABC123XYZ", "DEF456UVW"],
  "app": {
    "appName": "vicohome",
    "appVersion": "1.2.3",
    "appBuild": "123",
    "channelId": 1
  },
  "countryNo": "US",
  "language": "en",
  "tenantId": "default"
}
```

## Response Structure
The endpoint returns a generic Object response:

| Field | Type | Description |
|-------|------|-------------|
| result | Integer | Status code (0 for success, negative values indicate errors) |
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
  "result": -1001,
  "msg": "Device not found"
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -1001 | Device not found |
| -1002 | No active connection to close |
| -2001 | Network error |
| -4001 | Authentication error |

## Usage in Application
The endpoint is called in these scenarios:
1. When a user exits a live view or playback screen, closing the real-time connection
2. When the app detects network changes that require reconnection
3. When the app is backgrounded or closed to release resources
4. As part of error handling when P2P connections encounter issues
5. During device removal or account logout

The implementation is in the A4xPeerFetch class, which manages WebRTC connections in the app.

## Related Endpoints
- `/dc/rtcconnection/open` (openP2PConnection) - Opens peer-to-peer connections with devices

## Constraints
- User must be authenticated to access this endpoint
- The devices in the serialNumber list must belong to the authenticated user
- Only active connections can be closed
- This endpoint primarily affects server-side signaling for WebRTC connections