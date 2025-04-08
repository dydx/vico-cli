# Get WebRTC Ticket

Get a WebRTC ticket for establishing a P2P connection with a device for live video streaming.

**URL:** `/device/getWebrtcTicket`
**Method:** `POST`
**Auth Required:** Yes

## Request Body

```json
{
  "serialNumber": "string"  // Required. The serial number of the device to connect to
}
```

## Response

### Success Response

**Code:** `200 OK`
**Content:**

```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "traceId": "string",          // Unique identifier for the connection session
    "groupId": "string",          // Group identifier for the connection
    "role": "string",             // Role in the connection (e.g., "publisher", "subscriber")
    "id": "string",               // Unique identifier
    "signalServer": "string",     // WebRTC signaling server URL
    "sign": "string",             // Authentication signature
    "time": "string",             // Timestamp
    "signalPingInterval": 0,      // Interval for ping messages to keep connection alive
    "appStopLiveTimeout": "string", // Timeout value for the app to stop the connection
    "expirationTime": 0,          // When the ticket expires
    "iceServer": [                // Array of ICE servers for WebRTC connection
      {
        "url": "string",         // STUN/TURN server URL
        "username": "string",    // Authentication username for TURN
        "credential": "string"   // Authentication credential for TURN
      }
    ]
  }
}
```

### Error Responses

**Condition:** Invalid device serial number or device offline
**Code:** `400 BAD REQUEST`
**Content:**

```json
{
  "code": 101,
  "message": "Invalid device serial number"
}
```

**Condition:** Device is offline or unavailable for P2P connection
**Code:** `400 BAD REQUEST`
**Content:**

```json
{
  "code": 102,
  "message": "Device is offline or unavailable"
}
```

**Condition:** Authentication error
**Code:** `401 UNAUTHORIZED`
**Content:**

```json
{
  "code": 401,
  "message": "Invalid authentication token"
}
```

## Usage Notes

- This endpoint is used to establish a WebRTC connection for live video streaming from a device
- The returned ticket includes ICE servers for establishing P2P connections
- The client should use the returned signaling server to establish the WebRTC connection
- The connection will expire after the specified expiration time
- Keep-alive pings should be sent at the specified interval to maintain the connection
- When implementing WebRTC, proper error handling should account for connection failures
- The traceId should be included in any connection report for debugging purposes
- VPN usage may affect connection quality, so the client often reports VPN status to the server
- The returned WebRTC ticket is used by the client's WebRTC SDK to establish the connection

## Sample Usage

### Request

```http
POST /device/getWebrtcTicket HTTP/1.1
Host: api-us.vicohome.io
Content-Type: application/json
Authorization: Bearer {access_token}

{
  "serialNumber": "VDCA1P0123456789"
}
```

### Response

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "code": 0,
  "message": "Success",
  "data": {
    "traceId": "webrtc-vdca1p-1234567890",
    "groupId": "grp-vdca1p-1234567890",
    "role": "subscriber",
    "id": "usr-1234567890",
    "signalServer": "wss://webrtc-signal.vicohome.io",
    "sign": "a1b2c3d4e5f6g7h8i9j0",
    "time": "2023-04-15T12:34:56.789Z",
    "signalPingInterval": 10000,
    "appStopLiveTimeout": "30000",
    "expirationTime": 1681564496789,
    "iceServer": [
      {
        "url": "stun:stun.vicohome.io:3478",
        "username": "",
        "credential": ""
      },
      {
        "url": "turn:turn.vicohome.io:3478?transport=udp",
        "username": "vicohome-turn",
        "credential": "turn-credential-123"
      },
      {
        "url": "turn:turn.vicohome.io:3478?transport=tcp",
        "username": "vicohome-turn",
        "credential": "turn-credential-123"
      }
    ]
  }
}
```