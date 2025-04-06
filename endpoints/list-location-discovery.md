# listLocation Endpoint

## Overview
The listLocation endpoint retrieves a list of all locations associated with the user's account. Locations in the Vicohome ecosystem represent physical places where devices are installed, such as homes, offices, or vacation properties.

## API Details
- **Path**: `/location/listlocation`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Retrieves a list of all locations associated with the authenticated user.

## Request Parameters
The request body requires only standard BaseEntry properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| app | Object | Yes | Application information |
| countryNo | String | Yes | Country code (e.g., "US") |
| language | String | Yes | Language code (e.g., "en") |
| tenantId | String | Yes | Tenant identifier |

## Request Example
```json
{
  "app": {
    "appType": "Android",
    "appName": "Vicohome",
    "versionName": "1.0.0",
    "version": 100,
    "apiVersion": "1",
    "bundle": "io.vicohome.app",
    "timeZone": "America/Los_Angeles",
    "env": "production"
  },
  "countryNo": "US",
  "language": "en",
  "tenantId": "AAA12345678"
}
```

## Response Structure
The response contains an array of location information:

| Property | Type | Description |
|----------|------|-------------|
| result | Integer | Status code (0 for success, negative values for errors) |
| msg | String | Status message or error description |
| data | Array | List of LocationBean objects |

### LocationBean Structure
Each location in the list has the following properties:

| Property | Type | Description |
|----------|------|-------------|
| id | Integer | Unique identifier for the location |
| locationName | String | User-assigned name for the location |
| adminId | Integer | User ID of the location administrator |
| country | String | Country name |
| state | String | State/province name |
| city | String | City name |
| district | String | District/area name |
| streetAddress1 | String | Primary street address |
| streetAddress2 | String | Secondary street address (apt, unit, etc.) |
| postalCode | String | Postal/ZIP code |
| insertTime | Long | Timestamp of when the location was created |
| isLocalData | Boolean | Flag indicating if this is local data |
| isSelected | Boolean | Flag indicating if this location is currently selected |

## Response Example
```json
{
  "result": 0,
  "msg": "success",
  "data": [
    {
      "id": 12345,
      "locationName": "Home",
      "adminId": 67890,
      "country": "US",
      "state": "California",
      "city": "San Francisco",
      "district": "Mission District",
      "streetAddress1": "123 Main St",
      "streetAddress2": "Apt 4B",
      "postalCode": "94107",
      "insertTime": 1617293000,
      "isLocalData": false,
      "isSelected": true
    },
    {
      "id": 12346,
      "locationName": "Vacation Home",
      "adminId": 67890,
      "country": "US",
      "state": "Colorado",
      "city": "Aspen",
      "district": "Downtown",
      "streetAddress1": "456 Mountain Ave",
      "streetAddress2": "",
      "postalCode": "81611",
      "insertTime": 1620148200,
      "isLocalData": false,
      "isSelected": false
    }
  ]
}
```

## Error Codes
| Error Code | Description |
|------------|-------------|
| -1 | General error |
| -2 | Authentication error |
| -5 | Network error |

## Usage Context
This endpoint is typically used in the following scenarios:
- Displaying a list of locations in the app's settings
- When selecting a location to filter devices by location
- When configuring location-specific settings
- During the initial app setup to determine available locations

## Related Endpoints
- `addLocation` - Adds a new location
- `deleteLocation` - Deletes an existing location
- `updateLocation` - Updates location details
- `listDevice` - Lists all devices (which may be associated with locations)

## Implementation Notes
The endpoint retrieves all locations associated with the authenticated user's account. The response includes detailed information about each location, allowing the client application to display a comprehensive location list. The application typically uses this data to organize devices by location, enabling users to manage multiple installation sites within a single account.