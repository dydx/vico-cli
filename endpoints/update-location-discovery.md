# updateLocation Endpoint Documentation

## Overview
The `updateLocation` endpoint updates information about an existing location in the system. Locations are used to organize devices geographically, helping users manage devices across multiple sites.

## API Details
- **Path**: `/location/updatelocationinfo`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts an `UpdateLocationEntry` object with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | Integer | Yes | Unique identifier for the location to update |
| adminId | Integer | Yes | User ID of the admin |
| locationName | String | Yes | Name of the location |
| country | String | Yes | Country name |
| state | String | Yes | State/province name |
| city | String | Yes | City name |
| district | String | No | District/area name |
| streetAddress1 | String | No | Primary street address |
| streetAddress2 | String | No | Secondary street address |
| postalCode | String | No | Postal/ZIP code |
| insertTime | Integer | Yes | Timestamp for insertion |
| app | Object | Yes | Application information |
| countryNo | String | Yes | Country code |
| language | String | Yes | Language code |
| tenantId | String | Yes | Tenant identifier |

### Example Request Body
```json
{
  "id": 12345,
  "adminId": 67890,
  "locationName": "Home",
  "country": "US",
  "state": "California",
  "city": "San Francisco",
  "district": "Mission District",
  "streetAddress1": "123 Main St",
  "streetAddress2": "Apt 4B",
  "postalCode": "94107",
  "insertTime": 1617293000,
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
The endpoint returns a `BaseResponse` object:

### BaseResponse Fields
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
  "msg": "Location not found"
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -1001 | Location not found |
| -2001 | Network error |
| -3001 | Invalid parameters |

## Usage in Application
The endpoint is called when the user updates a location:
1. User modifies the location information in a form
2. The application creates an `UpdateLocationEntry` object with the provided data
3. The request is made through the `DeviceSettingApiClient`
4. On success, the updated location information is reflected in the UI
5. On failure, an appropriate error message is displayed to the user

## Constraints
- The location ID must exist in the system
- The user must have permission to update the location
- Required fields must not be empty or null