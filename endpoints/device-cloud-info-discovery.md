# Device Cloud Info Endpoint Discovery

## Endpoint Information
- **Path:** `/vip/device/cloud/info`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Queries VIP/subscription information for a specific device

## Request Parameters
The endpoint takes a `SerialNoEntry` object in the request body:

```json
{
  "serialNumber": "string",    // Required: Device serial number
  "app": {                     // Standard BaseEntry fields
    "type": "string",          // Application type identifier
    "version": "string"        // Application version
  },
  "countryNo": "string",       // User's country code
  "language": "string",        // User's language preference 
  "tenantId": "string"         // User's tenant ID
}
```

## Response Format
The endpoint returns a `DeviceCloudServiceResponse` object which extends the `BaseResponse`:

```json
{
  "result": 0,            // Result code (0 indicates success)
  "msg": "string",        // Response message
  "data": {
    "endTime": number,         // Subscription end time (Unix timestamp in seconds)
    "freeLicenseId": number,   // Free license ID (if applicable)
    "hasVip": boolean,         // Whether the device has an active VIP subscription
    "tierList": [              // List of tiers/subscription options
      {
        // Tier details
      }
    ],
    "tierNameKey": "string",   // Name key of the subscription tier
    "vipNotifyShow": boolean,  // Whether to show VIP notifications
    "vipTag": number           // VIP tag value used to identify subscription type
  }
}
```

## Code Analysis
The endpoint is implemented in API interfaces:

Example call signature from decompiled code:
```java
@POST("/vip/device/cloud/info")
Observable<DeviceCloudServiceResponse> getDeviceCloudServiceVIPInfo(@Body SerialNoEntry serialNoEntry);
```

Implementation in the application:
```java
public void getDeviceCloudServiceVIPInfo(String serialNumber, final Callback<DeviceCloudServiceResponse> callback) {
    SerialNoEntry entry = new SerialNoEntry();
    entry.setSerialNumber(serialNumber);
    
    apiClient.getDeviceCloudServiceVIPInfo(entry)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<DeviceCloudServiceResponse>() {
            @Override
            public void onNext(DeviceCloudServiceResponse response) {
                if (callback != null) {
                    callback.onSuccess(response);
                }
            }
            
            @Override
            public void onError(Throwable e) {
                if (callback != null) {
                    // Return default response with null data on error
                    DeviceCloudServiceResponse defaultResponse = new DeviceCloudServiceResponse();
                    defaultResponse.setResult(0);
                    defaultResponse.setData(null);
                    callback.onSuccess(defaultResponse);
                }
            }
            
            // Other observer methods
        });
}
```

## Usage Context
This endpoint is used in the following scenarios:
1. In the `DeviceConfigViewModel` when loading device settings
2. To determine if a device has an active cloud subscription
3. To display subscription tier information to the user
4. To check when a subscription expires
5. To decide whether to show VIP notifications or promotions

The endpoint is only called when:
- The user is an admin of the device
- The device is not being accessed via AP (direct) connection
- The device is online

## Error Handling
If any of the call conditions aren't met, the app creates a default response with null data:
```json
{
  "result": 0,
  "data": null
}
```

This default response is also used when there are network errors or other failures.

## Subscription Information
The response contains several important fields:
- `hasVip`: Indicates whether the device has an active subscription
- `endTime`: Shows when the current subscription expires
- `tierNameKey`: Identifies the subscription tier/level
- `vipTag`: A numeric identifier for the subscription type
- `freeLicenseId`: Used if the device has a free license rather than a paid subscription

## Notes
- This endpoint is part of the VIP/subscription management API group
- The API supports both paid subscriptions and free licenses
- The endpoint is only called for devices where the user has admin rights
- The response helps the application determine what features should be available based on subscription status
- The VIP notifications can be toggled on/off with the vipNotifyShow property