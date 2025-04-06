# getipinfo Endpoint Documentation

## Overview
The `getipinfo` endpoint retrieves information about the user's IP address and geographic location. This information is used for region-specific settings, support routing, and determining the appropriate API servers for the user's geographical location.

## API Details
- **Path**: `/account/getipinfo`
- **Method**: GET
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
This endpoint does not require any parameters. It determines the IP address from the incoming request.

## Response Structure
The endpoint returns a `NodeInfoResponse` object which extends `BaseResponse`:

### Base Response Fields
| Field | Type | Description |
|-------|------|-------------|
| result | int | Status code (0 for success, negative values indicate errors) |
| msg | String | Status message or error description |

### NodeInfoResponse Additional Fields
| Field | Type | Description |
|-------|------|-------------|
| data.ip | String | The user's IP address |
| data.countryCode | String | Two-letter country code based on IP geolocation (e.g., "US", "UK") |

### Example Success Response
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "ip": "203.0.113.1",
    "countryCode": "US"
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

## Usage in Application
The endpoint is used by the application in the following scenarios:
1. During startup to determine the user's region for proper API server routing
2. In the `AccountManager` class to configure region-specific settings:
   - Setting the appropriate support host URL for customer service
   - Updating API base URLs based on the user's region
   - Configuring region-specific tracking parameters
3. For determining which content delivery networks (CDNs) to use for media content

## Constraints
- No authentication is required for this endpoint
- Rate limiting may apply to prevent abuse
- IP address information is approximate and may not always reflect the user's exact location