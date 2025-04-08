# Query Free Tier Endpoint Discovery

## Endpoint Information
- **Path:** `/user/queryFreeTier`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Queries the user's free cloud storage tier status and limits

## Request Parameters
The endpoint accepts a standard `BaseEntry` object in the request body:

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

## Response Format
The endpoint returns a `FreeTierResponse` object:

```json
{
  "result": 0,                  // Result code (0 indicates success)
  "msg": "string",              // Response message
  "data": {
    "effectTime": 1639036800,   // Unix timestamp when free tier began
    "endTime": 1670572800,      // Unix timestamp when free tier expires
    "rollingDays": 7,           // Number of days recordings are kept
    "storage": 1073741824,      // Storage space in bytes (1GB in this example)
    "valid": true               // Whether the free tier is currently valid
  }
}
```

## Code Analysis
The endpoint is implemented in the API interfaces:

```java
@POST("/user/queryFreeTier")
Observable<FreeTierResponse> queryFreeTier(@Body BaseEntry baseEntry);
```

Implementation in DeviceSettingApiClient:
```java
@Override
public Observable<FreeTierResponse> queryFreeTier(BaseEntry baseEntry) {
    BaseEntry baseEntry2 = new BaseEntry();
    NetworkHelper.getInstance().wrapBaseEntry(baseEntry2);
    return apiService.queryFreeTier(baseEntry2);
}

public final Observable<FreeTierResponse> queryFreeTier() {
    BaseEntry baseEntry = new BaseEntry();
    NetworkHelper.getInstance().wrapBaseEntry(baseEntry);
    return queryFreeTier(baseEntry);
}
```

Usage example in DeviceManageCore:
```java
public final void queryFreeTier(Callback<FreeTierResponse> callback) {
    checkNotNullParameter(callback, "callback");
    DeviceSettingApiClient.getInstance().queryFreeTier()
        .subscribeOn(Schedulers.io())
        .observeOn(AndroidSchedulers.mainThread())
        .subscribe(new Subscriber<FreeTierResponse>() {
            @Override
            public void onNext(FreeTierResponse response) {
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
This endpoint is used to determine the free cloud storage tier status for a user. Based on the data structure, it's used to:

1. Check if the user has a valid free tier cloud storage allocation
2. Determine how much storage space is available (in bytes)
3. Find out how long recordings are kept (rolling days)
4. Check when the free tier began and when it expires

The information is likely used for:
- Determining if a user needs to upgrade to a paid subscription
- Calculating available cloud storage space
- Managing cloud recording retention policies
- Showing appropriate storage-related UI elements

## Free Tier Parameters
The response contains several important parameters:

1. `valid`: Boolean flag indicating if the free tier is currently active
2. `storage`: The amount of cloud storage allocated in bytes
3. `rollingDays`: How many days of footage are kept before being overwritten
4. `effectTime`: When the free tier became active (Unix timestamp)
5. `endTime`: When the free tier expires (Unix timestamp)

These parameters define the limitations of the free cloud storage offering.

## Business Logic
The free tier system implements several business rules:

1. Time-limited free trial:
   - Has specific start and end dates
   - May be provided to new users or devices
   - Has an expiration after which users must upgrade

2. Storage limitations:
   - Fixed storage amount (likely 1GB based on the code)
   - Limited retention period (rollingDays)
   - Older recordings are deleted when the limit is reached

3. Service differentiation:
   - Free tier offers limited functionality
   - Paid tiers offer more storage and longer retention
   - This endpoint helps determine which UI options to show

## Related Endpoints
This endpoint works with these related endpoints:

1. `/vip/tier/list/v4` - Gets available subscription tiers
   - Shows upgrade options beyond the free tier

2. `/vip/device/cloud/info` - Gets device VIP information
   - More detailed device-specific subscription info

3. `/vip/user/service/info` - Gets user VIP information
   - Detailed information about paid subscription status

## Notes
- This endpoint is part of the VIP/subscription management API group
- The free tier appears to be a limited cloud storage offering
- The API provides time-bound information about the free storage
- The rollingDays parameter indicates a circular buffer approach to storage
- Storage amount is likely represented in bytes
- This endpoint helps the app determine whether to show upgrade prompts
- The valid flag allows quick checking of free tier eligibility
- The time parameters (effectTime and endTime) could be used to show countdown timers for free trial expiration