# deleteLocation Endpoint Documentation

## Overview
The `deleteLocation` endpoint removes a location from the system. When a location is deleted, any devices associated with that location may be reassigned or have their location information updated.

## API Details
- **Path**: `/cationcation/deletelocation`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `DeleteLocationEntry` object with the following field:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | Integer | Yes | ID of the location to delete |

The `DeleteLocationEntry` extends the `BaseEntry` class, which includes common fields like language and tenantId.

### Example Request Body
```json
{
  "id": 12345
}
```

## Response Structure
The endpoint returns a `DeleteLocationResponse` object which extends `BaseResponse`:

### Base Response Fields
| Field | Type | Description |
|-------|------|-------------|
| result | int | Status code (0 for success, negative values indicate errors) |
| msg | String | Status message or error description |

### DeleteLocationResponse Additional Fields
| Field | Type | Description |
|-------|------|-------------|
| data.list | Array[DeviceBean] | List of devices affected by the location deletion |

Each DeviceBean in the list contains device information, including ID, name, serial number, and other attributes.

### Example Success Response
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 67890,
        "serialNumber": "ABC123XYZ",
        "name": "Front Door Camera"
        // other device attributes...
      }
    ]
  }
}
```

### Example Error Response
```json
{
  "result": -2001,
  "msg": "Network error"
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -2001 | Network error |
| -2002 | Access denied |
| -3001 | Location not found |
| -3002 | Cannot delete location with assigned devices |

## Usage in Application
The endpoint is called through the `DeviceLocationCore` class through the following sequence:
1. User selects a location to delete in the location management UI
2. A confirmation dialog may be shown to confirm the deletion
3. The `deleteUserDeviceLocation` method is called on the `DeviceLocationCore` instance with the location ID
4. A `DeleteLocationEntry` object is created with the location ID
5. The request is made to the API
6. The response is processed through callbacks
7. On success, the UI updates to remove the deleted location and potentially update device assignments
8. On error, an appropriate error message is displayed to the user

## Constraints
- The user must have appropriate permissions to delete the location
- Some implementations may prevent deleting a location if it has devices assigned to it
- The location ID must exist in the system
- The default/primary location may not be deletable