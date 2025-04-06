# queryDeviceCloudVipInfo Endpoint

## Overview
The queryDeviceCloudVipInfo endpoint retrieves information about premium cloud service subscriptions for a specific device. This endpoint is used to determine if a device has VIP/premium services enabled, what tier of service it has, and when the subscription expires.

## API Details
- **Path**: `/vip/device/cloud/info`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Retrieves cloud service subscription information for a specific device.

## Request Parameters
The request body should contain a JSON object with the following properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| serialNumber | String | Yes | The serial number of the device to query |

Plus standard BaseEntry properties:
- app
- countryNo
- language
- tenantId

## Request Example
```json
{
  "serialNumber": "ABC123456789",
  "app": "vicohome",
  "countryNo": "US",
  "language": "en",
  "tenantId": "AAA12345678"
}
```

## Response Structure
The response contains information about the device's cloud service subscription:

| Property | Type | Description |
|----------|------|-------------|
| result | Integer | Status code (0 for success, negative values for errors) |
| msg | String | Status message or error description |
| data | Object | Contains cloud service subscription information |

### Data Object Structure
| Property | Type | Description |
|----------|------|-------------|
| hasVip | Boolean | Whether the device has an active premium subscription |
| vipTag | Integer | VIP level or tag identifier |
| tierNameKey | String | Identifier for the subscription tier |
| endTime | Long | Timestamp for when the subscription ends |
| freeLicenseId | Integer | ID for any free license applied to the device |
| vipNotifyShow | Boolean | Whether VIP notifications should be shown |
| tierList | Array | List of available subscription tiers/options |

### Tier Object Structure (in tierList)
| Property | Type | Description |
|----------|------|-------------|
| tierId | Integer | Unique identifier for the tier |
| tierName | String | Display name for the tier |
| tierDescription | String | Description of the tier's features |
| price | Number | Price of the subscription |
| currency | String | Currency code |
| interval | String | Billing interval (e.g., "month", "year") |

## Response Example
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "hasVip": true,
    "vipTag": 2,
    "tierNameKey": "premium_plus",
    "endTime": 1686009600000,
    "freeLicenseId": 0,
    "vipNotifyShow": false,
    "tierList": [
      {
        "tierId": 1,
        "tierName": "Basic",
        "tierDescription": "7-day cloud storage",
        "price": 2.99,
        "currency": "USD",
        "interval": "month"
      },
      {
        "tierId": 2,
        "tierName": "Premium Plus",
        "tierDescription": "30-day cloud storage with additional features",
        "price": 9.99,
        "currency": "USD",
        "interval": "month"
      }
    ]
  }
}
```

## Error Codes
| Error Code | Description |
|------------|-------------|
| -1 | General error |
| -2 | Invalid parameters |
| -100 | Device not found |
| -101 | Device offline |
| -200 | Authentication error |

## Usage Context
This endpoint is typically used in the following scenarios:
- When the app needs to determine what cloud features are available to a device
- When displaying subscription status in device settings
- Before attempting to use premium features that require a subscription
- When showing subscription upgrade options to users
- During initial device setup to determine available cloud services

## Related Endpoints
- `queryFreeTier` - Queries free tier information
- `getVipUserServiceInfo` - Gets VIP service information for a user
- `getTierListV4` - Gets detailed subscription tier options

## Implementation Notes
The endpoint is only called when the device is online and the user has administrative privileges for the device. The response includes comprehensive information about the device's current cloud service subscription status, including whether it has premium features enabled, when the subscription expires, and what subscription tiers are available.

The application uses this information to determine what cloud-based features (like extended cloud storage, advanced AI detection, etc.) are available for the device and to adjust the user interface accordingly. The tierList in the response can be used to display upgrade options for users without an active premium subscription.