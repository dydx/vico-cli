# addLocation Endpoint Documentation

## Overview
The `addLocation` endpoint allows users to create a new location in the system. Locations are used to organize devices geographically, helping users manage devices across multiple sites.

## API Details
- **Path**: `/location/insertlocation/`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `LocationEntry` object with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| adminId | Integer | No | User ID of the admin |
| locationName | String | Yes | Name of the location |
| country | String | Yes | Country name |
| state | String | Yes | State/province name |
| city | String | Yes | City name |
| district | String | No | District/area name |
| streetAddress1 | String | No | Primary street address |
| streetAddress2 | String | No | Secondary street address (e.g., apartment number) |
| postalCode | String | No | Postal/ZIP code |
| insertTime | Integer | No | Timestamp for insertion |
| id | Integer | No | Location ID (used for updates rather than creation) |

The `LocationEntry` extends the `BaseEntry` class, which includes common fields like language and tenantId.

### Example Request Body
```json
{
  "locationName": "Home",
  "country": "US",
  "state": "California",
  "city": "San Francisco",
  "streetAddress1": "123 Main St",
  "postalCode": "94107"
}
```

## Response Structure
The endpoint returns an `AddLocationResponse` object which extends `BaseResponse`:

### Base Response Fields
| Field | Type | Description |
|-------|------|-------------|
| result | int | Status code (0 for success, negative values indicate errors) |
| msg | String | Status message or error description |

### AddLocationResponse Additional Fields
| Field | Type | Description |
|-------|------|-------------|
| data.id | int | Unique identifier for the created location |
| data.adminId | int | User ID of the admin |
| data.locationName | String | Name of the location |
| data.country | String | Country name |
| data.state | String | State/province name |
| data.city | String | City name |
| data.district | String | District/area name |
| data.streetAddress1 | String | Primary street address |
| data.streetAddress2 | String | Secondary street address |
| data.postalCode | String | Postal/ZIP code |
| data.insertTime | int | Timestamp of creation |
| data.isLocalData | boolean | Flag indicating if this is local data |
| data.isSelected | boolean | Flag indicating if this location is selected |

### Example Success Response
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "id": 12345,
    "adminId": 67890,
    "locationName": "Home",
    "country": "US",
    "state": "California",
    "city": "San Francisco",
    "streetAddress1": "123 Main St",
    "postalCode": "94107",
    "insertTime": 1617293000
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
| -3001 | Invalid location data |

## Usage in Application
The endpoint is called through the `DeviceLocationCore` class:
1. User creates a new location through the location management UI
2. A `LocationEntry` object is created with the location details
3. The `createUserDeviceLocation` method is called on the `DeviceLocationCore` instance
4. The request is processed and a response is returned
5. On success, the new location is displayed in the location list
6. On error, an appropriate error message is shown to the user

## Constraints
- A user may have a limited number of locations
- The location name must be unique for the user
- Required fields must be provided (locationName, country, state, city)
- The ID field should not be set when creating a new location