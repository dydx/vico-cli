# VIP User Service Info Endpoint Discovery

## Endpoint Information
- **Path:** `/vip/user/service/info`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Gets VIP user service subscription information and device allocation

## Request Parameters
The endpoint accepts a `BaseEntry` object in the request body:

```json
{
  "app": {                     // Standard app information
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

There's also a variant that accepts `Vip4GRequest` which extends `BaseEntry`:

```json
{
  // All BaseEntry fields, plus:
  "tierServiceType": 1         // Service type (1 = 4G service tiers)
}
```

## Response Format
The endpoint returns a `VipServiceResponse` object:

```json
{
  "result": 0,                  // Result code (0 indicates success)
  "msg": "string",              // Response message
  "data": {
    "notInTierDeviceList": [    // Devices not included in the tier
      {
        "deviceId": "string",   // Device identifier
        "serialNumber": "string", // Device serial number
        "deviceName": "string", // Device name
        "deviceType": "string", // Device type
        "imageUrl": "string",   // Device image URL
        "isAdmin": true         // Whether user is admin
      }
    ],
    "notInTierDeviceThirdPartyList": [  // Third-party devices not in tier
      // Same structure as notInTierDeviceList
    ],
    "outOfPlanKey": "string",   // Out-of-plan message key
    "rollingDay": 1,            // Billing day of month
    "serviceActiveDevice": {    // Currently active VIP service
      "activeDeviceList": [     // Devices active on this tier
        // Same structure as notInTierDeviceList
      ],
      "tierNameKey": "string",  // Key for tier name localization
      "endTime": 1648825200000, // Subscription end time (milliseconds)
      "maxDeviceNum": 5,        // Maximum devices allowed
      "tierId": 3,              // Tier ID
      "billingCycleType": 2,    // Billing cycle type (1=day, 2=month, 3=year)
      "billingCycleDuration": 1 // Billing cycle duration
    },
    "serviceAdditionList": [    // Additional services on subscription
      {
        "name": "string",       // Service name
        "description": "string" // Service description
      }
    ],
    "supportRollingDay": [1, 5, 10, 15, 20, 25], // Available billing day options
    "tierDescribeList": [       // Feature descriptions for tier
      {
        "title": "string",      // Feature title
        "type": "string",       // Feature type
        "value": "string"       // Feature value
      }
    ],
    "userVipId": 12345          // User VIP identifier
  }
}
```

## Code Analysis
The endpoint is implemented in the API interfaces:

```java
@POST("/vip/user/service/info")
Observable<VipServiceResponse> getVipUserServiceInfo(@Body BaseEntry baseEntry);
```

Implementation in PayApiClient:
```java
@Override
public Observable<VipServiceResponse> getVipUserServiceInfo(BaseEntry baseEntry) {
    return createApiService(baseEntry).getVipUserServiceInfo(baseEntry);
}
```

Usage example:
```java
public void getVipUserServiceInfo(final Callback<VipServiceResponse> callback) {
    BaseEntry baseEntry = new BaseEntry();
    
    payApiClient.getVipUserServiceInfo(baseEntry)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<VipServiceResponse>() {
            @Override
            public void onNext(VipServiceResponse response) {
                if (response != null && response.getResult() == 0) {
                    callback.onSuccess(response);
                } else {
                    callback.onFailure(response != null ? response.getMsg() : "Unknown error");
                }
            }
            
            @Override
            public void onError(Throwable e) {
                callback.onFailure(e.getMessage());
            }
        });
}
```

## Usage Context
This endpoint is used in the following scenarios:

1. In the VIP/subscription management screens:
   - To display current subscription tier and benefits
   - To show which devices are included in the subscription
   - To show which devices are not included in the subscription

2. For subscription management:
   - To allow users to change which devices are covered by their subscription
   - To enable drag-and-drop reordering of devices between "included" and "excluded" lists
   - To show available billing day options and current selection

3. For displaying subscription details:
   - Tier name and features
   - Expiration date
   - Maximum allowed devices
   - Current billing cycle information

## UI Implementation
The app processes the response data in several key components:

1. `VipServiceActivity`:
   - `setSupportDeviceList()`: Displays devices in the subscription
   - `setNameAndFunctionCard()`: Shows tier name and features
   - `addFunctionDescription()`: Adds feature descriptions to UI cards
   - `showVideoHistory()`: Displays video history duration

2. `VipManagerFragment`:
   - Handles drag-and-drop device reordering
   - Updates server when device allocation changes
   - Shows separate lists for included and excluded devices

## Error Handling
The application handles errors by:

1. Checking if the response object is null
2. Verifying that response.getResult() == 0 (success code)
3. On error, extracting the error message from response.getMsg() or exception
4. Displaying appropriate UI feedback (error dialogs or toasts)
5. When device reordering fails, showing a specific error dialog

## Device Management
The response provides crucial information about device subscription status:

1. `serviceActiveDevice.activeDeviceList`: Devices currently covered by the subscription
   - These devices receive premium features based on the subscription tier
   - Limited by `maxDeviceNum` in the subscription

2. `notInTierDeviceList`: User's devices not covered by the subscription
   - These devices only receive basic features
   - Can be added to the subscription if slots are available

3. Device allocation management:
   - Users can drag-and-drop devices between lists
   - Changes are saved back to the server with another API call
   - The app enforces the maximum device limit

## Subscription Information
The response contains comprehensive details about the subscription:

1. `serviceActiveDevice`: Core subscription details
   - `tierNameKey`: ID for localized tier name (e.g., "standard", "premium")
   - `endTime`: When the subscription expires
   - `maxDeviceNum`: Maximum devices allowed in this tier
   - `tierId`: Unique identifier for the subscription tier
   - `billingCycleType` and `billingCycleDuration`: Subscription period info

2. `tierDescribeList`: Feature descriptions for marketing and UI display
   - Lists benefits included in the subscription
   - Used to populate feature lists in the UI

3. `rollingDay` and `supportRollingDay`: Billing day settings
   - Current billing day of the month
   - List of available options for changing billing day

## Related Endpoints
This endpoint works with these related endpoints:
- `/vip/tier/list/v4`: Gets available subscription tiers
- `/vip/device/cloud/info`: Gets device-specific subscription info
- `/vip/user/device/list`: Lists devices with VIP status
- `/vip/user/rolling/day`: Updates the user's billing day

## Notes
- This endpoint is essential for the VIP/subscription management UI
- The endpoint allows for flexible device allocation within subscription limits
- The response structure supports both standard cloud subscriptions and 4G service plans
- The API design separates device management from subscription tier selection
- Third-party devices appear in a separate list (`notInTierDeviceThirdPartyList`)
- The `userVipId` uniquely identifies the user's subscription record
- Billing day management is included directly in this API response