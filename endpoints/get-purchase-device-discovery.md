# getPurchaseDevice Endpoint Documentation

## Overview
The `getPurchaseDevice` endpoint retrieves information about devices that are available for subscription purchase. It's used to determine which devices can be associated with a new subscription before proceeding with payment.

## API Details
- **Path**: `/vip/available/purchase/device`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `ProductEntry` object that extends `BaseEntry` with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| productId | Integer | Yes | The ID of the product/subscription plan to purchase |
| app | Object | Yes | Application information (inherited from BaseEntry) |
| countryNo | String | Yes | Country code (e.g., "US") (inherited from BaseEntry) |
| language | String | Yes | Language code (e.g., "en") (inherited from BaseEntry) |
| tenantId | String | Yes | Tenant identifier (inherited from BaseEntry) |

### Example Request Body
```json
{
  "productId": 123,
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
The endpoint returns a `PurchaseDevice` object that extends `BaseResponse`:

### BaseResponse Fields
| Field | Type | Description |
|-------|------|-------------|
| result | Integer | Status code (0 for success, negative values indicate errors) |
| msg | String | Status message or error description |

### PurchaseDevice Additional Fields
| Field | Type | Description |
|-------|------|-------------|
| data.availableForPurchase | Boolean | Whether purchase is available for any device |
| data.allSimThirdParty | Boolean | Flag related to SIM card/third-party connectivity |
| data.deviceList | Array | List of devices available for subscription |
| data.deviceList[].canChoose | Boolean | Whether the device can be selected for subscription |
| data.deviceList[].deviceName | String | Name of the device |
| data.deviceList[].icon | String | Device icon URL |
| data.deviceList[].serialNumber | String | Device serial number |
| data.deviceList[].isCheckout | Boolean | Whether the device is already checked out |

### Example Success Response
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "availableForPurchase": true,
    "allSimThirdParty": false,
    "deviceList": [
      {
        "canChoose": true,
        "deviceName": "Front Door Camera",
        "icon": "https://cdn.vicohome.io/icons/camera_doorbell.png",
        "serialNumber": "ABC123XYZ",
        "isCheckout": false
      },
      {
        "canChoose": true,
        "deviceName": "Backyard Camera",
        "icon": "https://cdn.vicohome.io/icons/camera_outdoor.png",
        "serialNumber": "DEF456UVW",
        "isCheckout": false
      }
    ]
  }
}
```

### Example Error Response
```json
{
  "result": -1001,
  "msg": "No available devices"
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -1001 | No available devices |
| -1002 | Purchase limit reached |
| -2001 | Network error |
| -4001 | Authentication error |

## Usage in Application
The endpoint is called during the subscription/payment flow:
1. User selects a subscription plan to purchase
2. Application calls the endpoint with the product ID
3. If multiple devices are available, the user is shown a device selection dialog
4. If only one device is available, it is automatically selected
5. If no devices are available, a "No Device" dialog is shown
6. If purchase limit is reached, a "Reached Limits" dialog is shown
7. After device selection, the user proceeds to the payment screen

## Related Endpoints
- `/vip/device/cloud/info` - Gets existing subscription information for a device
- `/vip/tier/list/v4` - Gets subscription tier options

## Constraints
- User must be authenticated to access this endpoint
- User must have at least one device registered to their account
- Some devices may not be eligible for certain subscription plans