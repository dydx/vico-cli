# getVipUserServiceInfo Endpoint Documentation

## Overview
The `getVipUserServiceInfo` endpoint retrieves detailed information about a user's current VIP subscription services. This includes information about active devices, service tiers, subscription durations, billing information, and supported features.

## API Details
- **Path**: `/vip/user/service/info`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts either a `BaseEntry` object or a `Vip4GRequest` object (which extends `BaseEntry`):

### When using BaseEntry:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| app | Object | Yes | Application information |
| countryNo | String | Yes | Country code (e.g., "US") |
| language | String | Yes | Language code (e.g., "en") |
| tenantId | String | Yes | Tenant identifier |

### When using Vip4GRequest:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| tierServiceType | Integer | Yes | Type of tier service to retrieve (typically set to 1) |
| app | Object | Yes | Application information (inherited from BaseEntry) |
| countryNo | String | Yes | Country code (e.g., "US") (inherited from BaseEntry) |
| language | String | Yes | Language code (e.g., "en") (inherited from BaseEntry) |
| tenantId | String | Yes | Tenant identifier (inherited from BaseEntry) |

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
The endpoint returns a `VipServiceResponse` object that extends `BaseResponse`:

### BaseResponse Fields
| Field | Type | Description |
|-------|------|-------------|
| result | Integer | Status code (0 for success, negative values indicate errors) |
| msg | String | Status message or error description |

### VipServiceResponse Additional Fields
| Field | Type | Description |
|-------|------|-------------|
| data.notInTierDeviceList | Array | List of user's devices not included in the current subscription tier |
| data.notInTierDeviceThirdPartyList | Array | List of third-party devices not included in the current tier |
| data.outOfPlanKey | String | Key indicating why a device is out of plan |
| data.rollingDay | Integer | Billing cycle day |
| data.serviceActiveDevice | Object | Information about the active subscription and associated devices |
| data.serviceActiveDevice.activeDeviceList | Array | List of devices activated under this subscription |
| data.serviceActiveDevice.tierNameKey | String | Key/name of the current subscription tier |
| data.serviceActiveDevice.endTime | Long | Timestamp when the subscription ends |
| data.serviceActiveDevice.maxDeviceNum | Integer | Maximum number of devices allowed in this tier |
| data.serviceActiveDevice.tierId | Integer | ID of the subscription tier |
| data.serviceActiveDevice.billingCycleType | Integer | Type of billing cycle |
| data.serviceActiveDevice.billingCycleDuration | Integer | Duration of the billing cycle |
| data.serviceAdditionList | Array | List of add-on services |
| data.supportRollingDay | Array | List of supported billing cycle days |
| data.tierDescribeList | Array | List of tier descriptions with titles and links |
| data.userVipId | Integer | Unique identifier for the user's VIP status |

### PaySupportDevice Object Structure
| Field | Type | Description |
|-------|------|-------------|
| deviceId | Integer | Unique identifier for the device |
| deviceInfo | Object | Device information |
| deviceInfo.deviceName | String | User-assigned device name |
| deviceInfo.serialNumber | String | Device serial number |
| deviceInfo.modelName | String | Model name |
| icon | String | URL to device icon |
| isSelected | Boolean | Whether the device is selected |
| order | Integer | Device order/priority |

### Example Success Response
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "rollingDay": 15,
    "serviceActiveDevice": {
      "activeDeviceList": [
        {
          "deviceId": 12345,
          "deviceInfo": {
            "deviceName": "Front Door Camera",
            "serialNumber": "ABC123XYZ",
            "modelName": "VH100"
          },
          "icon": "https://cdn.vicohome.io/icons/camera_doorbell.png",
          "isSelected": true,
          "order": 1
        }
      ],
      "tierNameKey": "premium",
      "endTime": 1701388800000,
      "maxDeviceNum": 5,
      "tierId": 2,
      "billingCycleType": 1,
      "billingCycleDuration": 30
    },
    "notInTierDeviceList": [
      {
        "deviceId": 67890,
        "deviceInfo": {
          "deviceName": "Backyard Camera",
          "serialNumber": "DEF456UVW",
          "modelName": "VH200"
        },
        "icon": "https://cdn.vicohome.io/icons/camera_outdoor.png",
        "isSelected": false,
        "order": 2
      }
    ],
    "supportRollingDay": [1, 15, 28],
    "tierDescribeList": [
      {
        "title": "Premium Features",
        "describe": "30-day video history, AI detection",
        "url": "https://vicohome.io/premium-features"
      }
    ],
    "userVipId": 98765
  }
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
1. When a user navigates to their VIP/subscription management screen
2. In the VipServiceActivity and Vip4GServiceActivity screens
3. To display current subscription information, active devices, and renewal dates
4. To allow users to manage which devices are included in their subscription
5. To provide information for modifying billing cycle dates

## Related Endpoints
- `/vip/tier/list/v4` (getTierListV4) - Gets subscription tier options
- `/vip/device/cloud/info` - Gets subscription information for a specific device

## Constraints
- User must be authenticated to access this endpoint
- The endpoint will only return information if the user has an active subscription
- The subscription details and available features may vary by region