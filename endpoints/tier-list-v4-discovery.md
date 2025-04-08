# Tier List V4 Endpoint Discovery

## Endpoint Information
- **Path:** `/vip/tier/list/v4`
- **Method:** GET
- **Base URL:** https://api-us.vicohome.io
- **Description:** Retrieves available subscription tier options and product information

## Request Parameters
The endpoint can be called with two different request types:

1. Standard request using `BaseEntry`:
```json
{
  "app": {
    "apiVersion": "string",    // API version
    "appName": "string",       // Application name
    "appType": "string",       // Application type
    "bundle": "string",        // Bundle identifier
    "countlyId": "string",     // Analytics ID
    "env": "string",           // Environment (e.g., "prod")
    "tenantId": "string",      // Tenant identifier
    "timeZone": "string",      // User's timezone
    "version": "string",       // App version number
    "versionName": "string"    // App version name
  },
  "countryNo": "string",       // Country code
  "language": "string",        // User's language preference
  "tenantId": "string"         // User's tenant ID
}
```

2. 4G-specific request using `Vip4GRequest` (extends BaseEntry):
```json
{
  // All BaseEntry fields, plus:
  "tierServiceType": 1         // Service type (1 = 4G service tiers)
}
```

## Response Format
The endpoint returns a `TierListResponseV4` object with the following structure:

```json
{
  "result": 0,                  // Result code (0 indicates success)
  "msg": "string",              // Response message
  "data": {
    "cloudFreeTrialDay": 30,    // Free trial duration in days
    "commendProduct": [         // Recommended regular products
      {
        "dayLookBack": 30,      // Days of recording history available
        "deviceNum": 1,         // Number of devices included
        "localBillingCycleDuration": 1,  // Billing cycle duration
        "localBillingCycleType": 2,      // Billing cycle type (1=day, 2=month, 3=year)
        "maxDeviceNum": 1,      // Maximum devices allowed
        "productId": 12345,     // Product ID
        "showInTier": true,     // Whether shown in tier selection
        "subscriptionGroupId": "group_id",  // Subscription group
        "tierId": 2,            // Tier ID
        "currency": "$",        // Currency symbol
        "currentPrice": "9.99", // Current price as string
        "priceMicro": 9990000,  // Price in micro units (millionths)
        "productName": "Cloud Standard",    // Product name
        "savePercentOfMontly": 0,           // Percent saved vs monthly
        "tierLevel": "standard"             // Tier level
      }
    ],
    "commendProduct4G": [       // Recommended 4G products (similar structure)
    ],
    "commendProduct4GHalfYearList": [  // 6-month 4G products
    ],
    "commendProduct4GYearList": [      // Annual 4G products
    ],
    "copywriteDiff": {          // Feature comparison info
      "diffList": [             // Feature differences between tiers
        {
          "title": "Video History",     // Feature name
          "type": "days",               // Feature type
          "value": "1 day",             // Free tier value
          "vipValue": "30 days"         // Paid tier value
        }
      ],
      "isSubscribed": false,    // Whether user is subscribed
      "storageDayCount": 1,     // Storage days in free tier
      "vipType": "none"         // Current VIP level
    },
    "hasBirdDevice": false,     // Whether user has bird devices
    "productExplain": [         // Explanatory text for products
      "All prices include applicable taxes"
    ],
    "tierList": {               // All tier options
      "tierV1List": {           // Old tier structure (for backward compatibility)
        "monthlyProductList": [],
        "yearlyProductList": []
      },
      "tierV2List": {           // Current tier structure
        "monthlyProductList": [  // Monthly subscription options
          // Similar to commendProduct structure
        ],
        "yearlyProductList": [   // Yearly subscription options
          // Similar to commendProduct structure with savings percentage
        ]
      }
    },
    "tierReceive": true,        // Whether tier data was received
    "tierReceive4G": true       // Whether 4G tier data was received
  }
}
```

## Code Analysis
The endpoint is implemented in the API interfaces:

```java
@GET("/vip/tier/list/v4")
Observable<TierListResponseV4> getTierListV4(@Query BaseEntry baseEntry);

@GET("/vip/tier/list/v4")
Observable<TierListResponseV4> getSim4gTierListV4(@Query Vip4GRequest vip4GRequest);
```

Implementation in PayApiClient:
```java
public Observable<TierListResponseV4> getTierListV4(BaseEntry baseEntry) {
    return createApiService(baseEntry).getTierListV4(baseEntry);
}

public Observable<TierListResponseV4> getSim4gTierListV4(Vip4GRequest vip4GRequest) {
    return createApiService(vip4GRequest).getSim4gTierListV4(vip4GRequest);
}
```

PayServerHelper calls:
```java
public void queryV4Products(final ProductsCallback callback) {
    BaseEntry baseEntry = new BaseEntry();
    
    payApiClient.getTierListV4(baseEntry)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<TierListResponseV4>() {
            @Override
            public void onNext(TierListResponseV4 response) {
                if (response != null && response.getResult() == 0) {
                    callback.onSuccess(response);
                } else {
                    callback.onFailure();
                }
            }
            
            @Override
            public void onError(Throwable e) {
                callback.onFailure();
            }
        });
}

public void querySim4gV4Products(final ProductsCallback callback) {
    Vip4GRequest request = new Vip4GRequest(1); // tierServiceType=1
    
    payApiClient.getSim4gTierListV4(request)
        .compose(ApplySchedulers.io_main())
        .subscribe(/* Similar observer pattern */);
}
```

## Usage Context
This endpoint is used in the following scenarios:

1. When loading the subscription/VIP management screen
2. To display available subscription tiers with pricing information
3. To show feature comparisons between free and paid tiers
4. To populate both regular cloud storage tiers and 4G-specific data plans
5. To get updated pricing information based on user's location/region
6. For calculating the savings percentage for annual plans compared to monthly

The UI typically displays tier options in two tabs:
- Monthly subscription options
- Annual subscription options (highlighting savings percentage)

## Error Handling
The application handles errors by:

1. Checking if the response object is null
2. Verifying that response.getResult() == 0 (success code)
3. If either check fails, the onFailure() callback is triggered
4. There's no retry mechanism implemented; failures simply show an error UI

## Related Endpoints
This endpoint works with these related endpoints:
- `/vip/device/cloud/info` - Gets device-specific subscription information
- `/vip/user/service/info` - Gets user-level subscription details
- `/vip/user/device/list` - Lists devices with VIP status

## Subscription Information
The response contains comprehensive details about subscription options:

1. Regular Cloud Storage Tiers:
   - Multiple tiers (e.g., Basic, Standard, Premium)
   - Different device count options (1, 3, 5, 10 devices)
   - Monthly and annual billing options
   - Video history duration (e.g., 7, 30, 60 days)

2. 4G Data Plans (if applicable):
   - Data volume options
   - Monthly, semi-annual, and annual plans
   - Device-specific 4G connectivity services

Each product includes:
- Price information (currency, formatted price, price in micro units)
- Billing cycle (duration and type)
- Features (video history days, device counts)
- Savings percentage for annual plans

## Notes
- This endpoint is part of the VIP/subscription management API group
- The tiered design allows users to select plans based on:
  - Number of devices (1, 3, 5, 10)
  - Billing cycle (monthly, yearly)
  - Feature level (basic, standard, premium)
- The copywriteDiff section provides clear feature comparisons for marketing purposes
- The response structure shows that the API evolved over time (tierV1List vs tierV2List)
- Currency and pricing are localized based on the user's countryNo
- The productExplain field provides important disclaimers about pricing