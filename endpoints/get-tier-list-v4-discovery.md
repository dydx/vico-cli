# getTierListV4 Endpoint Documentation

## Overview
The `getTierListV4` endpoint retrieves a list of available subscription tiers and plans for Vicohome devices. It provides comprehensive information about different subscription options including monthly and yearly plans, features, and pricing.

## API Details
- **Path**: `/vip/tier/list/v4`
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
The endpoint returns a `TierListResponseV4` object that extends `BaseResponse`:

### BaseResponse Fields
| Field | Type | Description |
|-------|------|-------------|
| result | Integer | Status code (0 for success, negative values indicate errors) |
| msg | String | Status message or error description |

### TierListResponseV4 Additional Fields
| Field | Type | Description |
|-------|------|-------------|
| data.cloudFreeTrialDay | Integer | Number of days in the free trial period |
| data.commendProduct | Array | List of recommended standard subscription products |
| data.commendProduct4G | Array | List of recommended 4G subscription products |
| data.commendProduct4GHalfYearList | Array | List of recommended 6-month 4G subscription products |
| data.commendProduct4GYearList | Array | List of recommended annual 4G subscription products |
| data.copywriteDiff | Object | Contains subscription information and feature differences |
| data.hasBirdDevice | Boolean | Whether the user has a bird detection compatible device |
| data.productExplain | Array | List of product explanations/descriptions |
| data.tierList | Object | Contains tierV1List and tierV2List, each with monthly and yearly product lists |
| data.tierReceive | Boolean | Whether the user can receive standard subscription offers |
| data.tierReceive4G | Boolean | Whether the user can receive 4G subscription offers |

### ProductList Object Structure
| Field | Type | Description |
|-------|------|-------------|
| id | Integer | Unique identifier for the product |
| name | String | Product name |
| price | Float | Product price |
| currency | String | Currency code (e.g., "USD") |
| periodType | Integer | Subscription period type (1=Monthly, 12=Yearly) |
| features | Array | List of features included in this product |
| tier | Integer | Tier level |
| discountPercentage | Float | Discount percentage if applicable |
| originalPrice | Float | Original price before discount |

### Example Success Response
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "cloudFreeTrialDay": 30,
    "commendProduct": [
      {
        "id": 123,
        "name": "Basic Monthly",
        "price": 4.99,
        "currency": "USD",
        "periodType": 1,
        "features": ["Motion Detection", "7-Day Video History"],
        "tier": 1,
        "discountPercentage": 0,
        "originalPrice": 4.99
      },
      {
        "id": 456,
        "name": "Premium Monthly",
        "price": 9.99,
        "currency": "USD",
        "periodType": 1,
        "features": ["Motion Detection", "30-Day Video History", "AI Detection"],
        "tier": 2,
        "discountPercentage": 0,
        "originalPrice": 9.99
      }
    ],
    "tierList": {
      "tierV1List": {
        "monthlyList": [
          {
            "id": 123,
            "name": "Basic Monthly",
            "price": 4.99,
            "currency": "USD",
            "periodType": 1,
            "features": ["Motion Detection", "7-Day Video History"],
            "tier": 1
          }
        ],
        "yearlyList": [
          {
            "id": 789,
            "name": "Basic Annual",
            "price": 49.99,
            "currency": "USD",
            "periodType": 12,
            "features": ["Motion Detection", "7-Day Video History"],
            "tier": 1,
            "discountPercentage": 16.67,
            "originalPrice": 59.88
          }
        ]
      },
      "tierV2List": {
        "monthlyList": [
          {
            "id": 456,
            "name": "Premium Monthly",
            "price": 9.99,
            "currency": "USD",
            "periodType": 1,
            "features": ["Motion Detection", "30-Day Video History", "AI Detection"],
            "tier": 2
          }
        ],
        "yearlyList": [
          {
            "id": 1011,
            "name": "Premium Annual",
            "price": 99.99,
            "currency": "USD",
            "periodType": 12,
            "features": ["Motion Detection", "30-Day Video History", "AI Detection"],
            "tier": 2,
            "discountPercentage": 16.67,
            "originalPrice": 119.88
          }
        ]
      }
    },
    "hasBirdDevice": false,
    "tierReceive": true,
    "tierReceive4G": false
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
| -3001 | Invalid parameters |
| -4001 | Authentication error |

## Usage in Application
The endpoint is called in the subscription management flow:
1. When a user navigates to the subscription management screen
2. When a user selects a device for subscription
3. The response is used to populate the subscription selection screen with various tiers and plans
4. The user can then select a preferred subscription plan for purchase
5. The application may highlight recommended plans based on the "commendProduct" fields

## Related Endpoints
- `/vip/tier/list/v4` (getSim4gTierListV4) - Gets 4G subscription tier options with a specific service type parameter
- `/vip/available/purchase/device` (getPurchaseDevice) - Gets devices available for subscription

## Constraints
- User must be authenticated to access this endpoint
- The endpoint will return subscription options based on the user's country and location
- Pricing and available plans may vary by region