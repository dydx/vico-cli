# uploadGooglePayInfo Endpoint Documentation

## Overview
The `uploadGooglePayInfo` endpoint processes and verifies Google Pay payments for subscription purchases within the Vicohome app. This endpoint is responsible for validating purchase tokens with Google, recording transactions, and activating subscriptions.

## API Details
- **Path**: `/pay/google/order/verify`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts an `OrderEntry` object that extends `BaseEntry` with the following fields:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| productId | Integer | Yes | The product identifier for the subscription plan |
| tradeType | Integer | Yes | Type of transaction (subscription or one-time) |
| outTradeNo | String | Yes | Order number for tracking the transaction |
| purchaseToken | String | Yes | Token received from Google Pay after user payment |
| subscriptionGroupId | String | No | Group identifier for subscription |
| tierDeviceList | Array | Yes | List of device identifiers to activate the subscription for |
| guidanceSource | Integer | No | Source of purchase guidance/referral |
| app | Object | Yes | Application information (inherited from BaseEntry) |
| countryNo | String | Yes | Country code (e.g., "US") (inherited from BaseEntry) |
| language | String | Yes | Language code (e.g., "en") (inherited from BaseEntry) |
| tenantId | String | Yes | Tenant identifier (inherited from BaseEntry) |

### Example Request Body
```json
{
  "productId": 123,
  "tradeType": 1,
  "outTradeNo": "ORD12345678",
  "purchaseToken": "google-pay-token-abcdef123456",
  "subscriptionGroupId": "sub_group_premium",
  "tierDeviceList": ["ABC123XYZ", "DEF456UVW"],
  "guidanceSource": 0,
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
The endpoint returns a `PayResultResponse` object that extends `BaseResponse`:

### BaseResponse Fields
| Field | Type | Description |
|-------|------|-------------|
| result | Integer | Status code (0 for success, negative values indicate errors) |
| msg | String | Status message or error description |

### PayResultResponse Additional Fields
| Field | Type | Description |
|-------|------|-------------|
| data | String | Additional response data, typically a confirmation identifier |

### Example Success Response
```json
{
  "result": 0,
  "msg": "success",
  "data": "transaction_id_12345"
}
```

### Example Error Response
```json
{
  "result": -1001,
  "msg": "Invalid payment token"
}
```

## Error Codes
| Code | Description |
|------|-------------|
| -1001 | Invalid payment token |
| -1002 | Product not found |
| -1003 | Transaction already processed |
| -2001 | Network error |
| -3001 | Server error during verification |
| -4001 | Authentication error |

## Usage in Application
The endpoint is called during the payment process:
1. User selects a subscription plan and proceeds to checkout
2. User completes payment through Google Pay
3. Google Pay returns a purchase token to the app
4. The app creates an OrderEntry with the token and device information
5. The endpoint verifies the token with Google's servers
6. Upon successful verification, the subscription is activated for the specified devices
7. The app updates the UI to reflect the active subscription

## Related Endpoints
- `/vip/tier/list/v4` (getTierListV4) - Gets subscription tier options
- `/vip/user/service/info` (getVipUserServiceInfo) - Gets VIP user service information after purchase

## Constraints
- User must be authenticated to access this endpoint
- The purchaseToken must be valid and recently generated from Google Pay
- The productId must correspond to a valid subscription product
- The devices in tierDeviceList must belong to the authenticated user
- The endpoint is specifically for Google Pay transactions on Android devices