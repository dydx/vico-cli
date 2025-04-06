# openP2PConnection Endpoint Documentation

## Overview
The `openP2PConnection` endpoint establishes peer-to-peer (P2P) connections with Vicohome devices for real-time video streaming and two-way audio communication. This endpoint is primarily used for direct communication with devices, especially when accessing SD card recordings.

## API Details
- **Path**: `/dc/rtcconnection/open`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `LibraryP2PConnectRequest` object that extends `BaseEntry` with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | Array | Yes | List of device serial numbers to establish P2P connections with |
| app | Object | Yes | Application information (inherited from BaseEntry) |
| countryNo | String | Yes | Country code (e.g., "US") (inherited from BaseEntry) |
| language | String | Yes | Language code (e.g., "en") (inherited from BaseEntry) |
| tenantId | String | Yes | Tenant identifier (inherited from BaseEntry) |

### Example Request Body
```json
{
  "serialNumber": ["ABC123XYZ"],
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
| -1002 | Device offline |
| -1003 | Connection failed |
| -2001 | Network error |
| -4001 | Authentication error |

## Usage in Application
The endpoint is called in these scenarios:
1. When a user selects an SD card camera to view its recordings
2. Before establishing WebRTC connections for live view
3. When accessing the SD card review functionality from the Library view
4. For two-way audio communication with devices
5. Through the LibraryCore.openP2PConnection method

The implementation is in the A4xPeerFetch class, which manages WebRTC connections in the app.

## Connection Flow
1. User selects an SD card camera to view its recordings
2. The app calls openP2PConnection with the device's serial number
3. Server establishes a signaling channel for WebRTC
4. When the connection is successfully established, LibraryFragment.onSDCardReviewDone() handles the UI transition
5. When viewing is complete, closeP2PConnection is called to release resources

## Related Endpoints
- `/dc/rtcconnection/close` (closeP2PConnection) - Closes peer-to-peer connections with devices

## Constraints
- User must be authenticated to access this endpoint
- The devices in the serialNumber list must belong to the authenticated user
- Devices must be online for the connection to be established
- P2P connections require significant bandwidth and may affect network performance
- This endpoint primarily affects server-side signaling for WebRTC connections