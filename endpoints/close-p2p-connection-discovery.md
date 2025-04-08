# Close P2P Connection Discovery

## Endpoint Details
- **Path**: `/dc/rtcconnection/close`
- **Method**: POST
- **Description**: Close a peer-to-peer (P2P) connection for live streaming

## Request Parameters

The request expects a JSON object with the following structure:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | List\<String\> | Yes | List of device serial numbers to close P2P connections with |

### Request Example
```json
{
  "serialNumber": ["DEVICE_SN_1", "DEVICE_SN_2"]
}
```

## Response

The endpoint returns a simple Object response with no specific structure. The success is indicated by the HTTP status code.

### Response Example
```json
{}
```

## Error Handling

Errors are handled through standard error callbacks. The error information includes:
- Error code (integer)
- Error message (string)

## Implementation Details

1. The endpoint is called via the `closeP2PConnection` method in the `LibraryCore` class.
2. It expects a list of device serial numbers for which to close P2P connections.
3. The request is processed via a reactive Observable pattern.
4. Errors are reported through callback functions.
5. The implementation uses RxJava for asynchronous processing.

## Usage Context

This endpoint is used to close WebRTC connections established with cameras. It's typically called:
- When a user exits the live streaming view
- When the app is going to background
- When the streaming session needs to be terminated for any reason

The companion to this endpoint is `/dc/rtcconnection/open` which establishes the P2P connection.

## Internal Implementation

When called, the app performs the following actions:
1. Creates a `LibraryP2PConnectRequest` object with the provided serial numbers
2. Sends the request to the backend server
3. Upon receiving a successful response, the WebRTC data channels are closed
4. Resources associated with the connection are released

This endpoint is part of the peer connection management system that handles device-to-app streaming communication.