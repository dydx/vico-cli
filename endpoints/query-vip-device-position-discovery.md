# queryVipDevicePosition Endpoint Documentation

## Overview
The `queryVipDevicePosition` endpoint retrieves information about the user's VIP subscription devices and their positions/statuses in the system. This endpoint provides details about which devices are activated under which subscription tiers and their current status.

## API Details
- **Path**: `/vip/user/device/list`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `BaseEntry` object with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| app | Object | Yes | Application information |
| countryNo | String | Yes | Country code (e.g., "US") |
| language | String | Yes | Language code (e.g., "en") |
| tenantId | String | Yes | Tenant identifier |

### Example Request Body
```json
{
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
The endpoint returns a `VipManagerResponse` object that extends `BaseResponse`:

### BaseResponse Fields
| Field | Type | Description |
|-------|------|-------------|
| result | Integer | Status code (0 for success, negative values indicate errors) |
| msg | String | Status message or error description |

### VipManagerResponse Additional Fields
| Field | Type | Description |
|-------|------|-------------|
| data | Array | List of VimManagerBean objects with subscription tier information |

### VimManagerBean Object Structure
| Field | Type | Description |
|-------|------|-------------|
| tierId | Integer | ID of the subscription tier |
| tierName | String | Display name of the subscription tier |
| tierNameKey | String | Key/identifier for the tier name |
| maxDeviceNum | Integer | Maximum number of devices allowed in this tier |
| customSupportAllDevice | Boolean | Whether all device types are supported in this tier |
| supportDeviceSnList | Array | List of serial numbers for supported devices |
| activeDeviceList | Array | List of devices currently active in this tier |

### PaySupportDevice Object Structure (in activeDeviceList)
| Field | Type | Description |
|-------|------|-------------|
| serialNumber | String | Device serial number |
| deviceName | String | User-assigned device name |
| icon | String | URL to device icon |
| smallIcon | String | URL to device small icon |
| tierList | String | List of tiers this device supports |
| tierName | String | Name of the tier the device is in |
| userSn | String | User's serial number |
| iccid | String | Integrated Circuit Card Identifier (for SIM devices) |
| deviceStatus | Integer | Current status code of the device |
| online | Integer | Online status (1 = online, 0 = offline) |
| simThirdParty | Integer | Flag for third-party SIM status |
| deviceStorageType | Integer | Type of storage used by the device |
| viewType | Integer | UI view type for this device |

### Example Success Response
```json
{
  "result": 0,
  "msg": "success",
  "data": [
    {
      "tierId": 2,
      "tierName": "Premium",
      "tierNameKey": "premium_tier",
      "maxDeviceNum": 5,
      "customSupportAllDevice": false,
      "supportDeviceSnList": ["VH100", "VH200", "VH300"],
      "activeDeviceList": [
        {
          "serialNumber": "ABC123XYZ",
          "deviceName": "Front Door Camera",
          "icon": "https://cdn.vicohome.io/icons/camera_doorbell.png",
          "smallIcon": "https://cdn.vicohome.io/icons/small/camera_doorbell.png",
          "tierList": "1,2,3",
          "tierName": "Premium",
          "userSn": "USR123456",
          "iccid": "",
          "deviceStatus": 1,
          "online": 1,
          "simThirdParty": 0,
          "deviceStorageType": 2,
          "viewType": 0
        },
        {
          "serialNumber": "DEF456UVW",
          "deviceName": "Backyard Camera",
          "icon": "https://cdn.vicohome.io/icons/camera_outdoor.png",
          "smallIcon": "https://cdn.vicohome.io/icons/small/camera_outdoor.png",
          "tierList": "1,2,3",
          "tierName": "Premium",
          "userSn": "USR123456",
          "iccid": "",
          "deviceStatus": 1,
          "online": 1,
          "simThirdParty": 0,
          "deviceStorageType": 2,
          "viewType": 0
        }
      ]
    },
    {
      "tierId": 3,
      "tierName": "Premium Plus",
      "tierNameKey": "premium_plus_tier",
      "maxDeviceNum": 10,
      "customSupportAllDevice": true,
      "supportDeviceSnList": [],
      "activeDeviceList": []
    }
  ]
}
```

### Example Error Response
```json
{
  "result": -1001,
  "msg": "No subscription found"
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -1001 | No subscription found |
| -2001 | Network error |
| -4001 | Authentication error |

## Usage in Application
The endpoint is called in the subscription management flow:
1. When a user navigates to the VIP Manager screen (VipManagerFragment)
2. To display a list of subscription tiers and devices activated under each tier
3. To show device status (online/offline) within the subscription interface
4. To provide information needed for adding or removing devices from subscriptions

## Related Endpoints
- `/vip/user/service/info` (getVipUserServiceInfo) - Gets VIP user service information
- `/vip/device/cloud/info` - Gets subscription information for a specific device

## Constraints
- User must be authenticated to access this endpoint
- The endpoint will only return information if the user has at least one active subscription
- The response will include all tiers the user has access to, even if there are no active devices in some tiers