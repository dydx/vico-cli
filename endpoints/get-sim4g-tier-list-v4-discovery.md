# getSim4gTierListV4 Endpoint Documentation

## Overview
The `getSim4gTierListV4` endpoint retrieves a list of available 4G SIM subscription tiers and plans for Vicohome devices with cellular connectivity. This endpoint provides comprehensive information about different subscription options including monthly and yearly plans.

## API Details
- **Path**: `/vip/tier/list/v4`
- **Method**: POST
- **Base URL**: https://api-us.vicohome.io

## Request Parameters
The endpoint accepts a `Vip4GRequest` object that extends `BaseEntry` with the following fields:

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
  "tierServiceType": 1,
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
    "commendProduct4G": [
      {
        "id": 123,
        "name": "4G Basic Monthly",
        "price": 9.99,
        "currency": "USD",
        "periodType": 1,
        "features": ["Unlimited Data", "30-Day Video History"],
        "tier": 1,
        "discountPercentage": 0,
        "originalPrice": 9.99
      }
    ],
    "commendProduct4GYearList": [
      {
        "id": 456,
        "name": "4G Premium Annual",
        "price": 99.99,
        "currency": "USD",
        "periodType": 12,
        "features": ["Unlimited Data", "30-Day Video History", "AI Detection"],
        "tier": 2,
        "discountPercentage": 16.67,
        "originalPrice": 119.99
      }
    ],
    "hasBirdDevice": false,
    "tierReceive": true,
    "tierReceive4G": true
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
2. When a user selects a 4G-enabled device for subscription
3. The PayServerHelper uses this endpoint to fetch available subscription options
4. The response is used to populate the subscription selection screen
5. The user can then select a preferred subscription plan for purchase

## Related Endpoints
- `/vip/tier/list/v4` (getTierListV4) - Gets standard subscription tier options
- `/vip/available/purchase/device` (getPurchaseDevice) - Gets devices available for subscription

## Constraints
- User must be authenticated to access this endpoint
- The endpoint will only return relevant subscription options based on the user's country