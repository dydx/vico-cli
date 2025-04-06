# deviceThumbImage Endpoint

## Overview
The deviceThumbImage endpoint allows clients to retrieve thumbnail images for camera devices associated with the user's account. These thumbnail images represent the most recent snapshot from each camera and are typically used in device list or dashboard views.

## API Details
- **Path**: `/device/devicePushImage`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Retrieves thumbnail images for all devices in the user's account or a specified subset of devices.

## Request Parameters
The request body should contain a JSON object with standard BaseEntry properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| app | Object | Yes | Application information |
| countryNo | String | Yes | Country code (e.g., "US") |
| language | String | Yes | Language code (e.g., "en") |
| tenantId | String | Yes | Tenant identifier |

### App Object Structure
| Parameter | Type | Description |
|-----------|------|-------------|
| appType | String | Application type (default "Android") |
| appName | String | Application name |
| versionName | String | App version name |
| version | Integer | App version number |
| apiVersion | String | API version |
| bundle | String | Bundle identifier |
| timeZone | String | Device timezone |
| env | String | Environment (e.g., "production") |

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
    "env": "production",
    "tenantId": "US"
  },
  "countryNo": "US",
  "language": "en",
  "tenantId": "US"
}
```

## Response Structure
The response contains a list of device thumbnail image information:

| Property | Type | Description |
|----------|------|-------------|
| result | Integer | Status code (0 for success, negative values for errors) |
| msg | String | Status message or error description |
| data | Array | List of DeviceCoverImageModel objects |

### DeviceCoverImageModel Structure
| Property | Type | Description |
|----------|------|-------------|
| serialNumber | String | Device serial number |
| lastPushImageUrl | String | URL to the thumbnail/preview image |
| lastPushTime | Long | Timestamp of when the image was last updated |

## Response Example
```json
{
  "result": 0,
  "msg": "success",
  "data": [
    {
      "serialNumber": "ABC123456789",
      "lastPushImageUrl": "https://images.vicohome.io/abc123/thumb.jpg",
      "lastPushTime": 1649876543210
    },
    {
      "serialNumber": "XYZ987654321",
      "lastPushImageUrl": "https://images.vicohome.io/xyz987/thumb.jpg",
      "lastPushTime": 1649876542105
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
- Loading the main device dashboard to show preview images for each camera
- Refreshing device thumbnails periodically to show recent images
- Displaying device cards in the app with current preview images

## Related Endpoints
- None directly related, as this is a specific thumbnail retrieval endpoint

## Implementation Notes
The endpoint is implemented in the DeviceSettingCore class with the getRecentCoverImage method. It can optionally filter results by specific device serial numbers, though in typical usage it retrieves thumbnails for all devices associated with the user's account. The implementation uses RxJava Observables for asynchronous API calls.