# VIP User Device List Endpoint Discovery

## Endpoint Information
- **Path:** `/vip/user/device/list`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Lists devices with VIP status and their subscription tier positions

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

## Response Format
The endpoint returns a `VipManagerResponse` object:

```json
{
  "result": 0,                  // Result code (0 indicates success)
  "msg": "string",              // Response message
  "data": [                     // List of VIP manager tiers
    {
      "tierId": 3,              // Tier ID
      "tierName": "Premium",    // Tier name
      "tierNameKey": "premium", // Tier name localization key
      "maxDeviceNum": 5,        // Maximum devices allowed in tier
      "activeDeviceList": [     // Devices in this tier
        {
          "serialNumber": "string",   // Device serial number
          "deviceName": "string",     // Device name
          "deviceStatus": 1,          // Device status code
          "online": 1,                // Online status (1=online, 0=offline)
          "icon": "https://example.com/icon.png", // Device icon URL
          "tierName": "Premium",      // Tier name device belongs to
          "isAdmin": true,            // Whether user is admin for this device
          "deviceType": "string",     // Device type identifier
          "deviceStorageType": 1      // Storage type for device (cloud vs local)
        }
      ],
      "supportDeviceSnList": [  // List of supported device serial numbers
        "string"
      ],
      "customSupportAllDevice": false  // Whether all devices are supported
    }
  ]
}
```

## Code Analysis
The endpoint is implemented in the API interfaces:

```java
@POST("/vip/user/device/list")
Observable<VipManagerResponse> queryVipDevicePosition(@Body BaseEntry baseEntry);
```

Implementation in PayApiClient:
```java
@Override
public Observable<VipManagerResponse> queryVipDevicePosition(BaseEntry entry) {
    checkNotNullParameter(entry, "entry");
    return createApiService(entry).queryVipDevicePosition(entry);
}
```

Usage example:
```java
// In VipManagerFragment
private void loadDataAndSetToUI() {
    PayApiClient.getInstance().queryVipDevicePosition(new BaseEntry())
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<VipManagerResponse>() {
            @Override
            public void onNext(VipManagerResponse response) {
                if (response != null && response.getResult() == 0) {
                    // Process response data and set to UI
                    processManagerResponse(response);
                } else {
                    // Handle error
                    String errorMsg = response != null ? response.getMsg() : "Unknown error";
                    showToast(errorMsg);
                }
            }
            
            @Override
            public void onError(Throwable e) {
                // Handle network or other errors
                showToast("Open Failed, Retry");
            }
        });
}
```

## Usage Context
This endpoint is used in the following scenarios:

1. VIP Manager UI:
   - To display a list of devices organized by subscription tier
   - To show which devices are included in paid subscription tiers
   - To show which devices are not included in paid tiers (basic tier)

2. Device management:
   - To allow users to view all their devices' subscription status
   - To enable drag-and-drop reordering of devices between tiers
   - To enforce maximum device limits per tier

3. Subscription management:
   - To help users visualize which devices receive premium features
   - To support decision-making about upgrading tiers based on device counts

## UI Implementation
The app processes the response data in several key components:

1. `VipManagerFragment`:
   - Main UI component that displays device subscription tiers
   - Transforms the API response into UI models
   - Implements drag-and-drop functionality for device reordering

2. `VipManagerAdapter`:
   - RecyclerView adapter for rendering device items
   - Supports different view types (tier headers and device items)
   - Shows visual indicators for device status and tier membership

3. Drag-and-drop implementation:
   - Uses ItemTouchHelper.SimpleCallback to handle drag events
   - Tracks start and end positions to update device tier memberships
   - Validates moves against tier device limits
   - Updates server with new device positions after changes

## Error Handling
The application handles errors by:

1. Checking for null responses and non-zero result codes
2. Displaying appropriate messages to the user:
   - Error code -2006: Shows "Device Selection Reminder" dialog
   - Network errors: Shows "Open Failed, Retry" toast
   - Other errors: Shows the error message from the API response

3. Implementing recovery options:
   - Retry button for network failures
   - Cancel option to dismiss errors and continue

## Device Position Management
The response provides structured information about device tier positions:

1. Data organization:
   - Devices are grouped by subscription tiers
   - Each tier has a maximum device limit
   - Tiers are presented in order of feature level

2. Position changes:
   - Users can drag devices between tiers
   - Position changes update the local data model
   - Changes are saved to the server with the updateVipDevicePosition API

3. Constraints:
   - Tier device limits are enforced
   - Device moves are validated before being accepted
   - Admin-only operations are restricted based on user permissions

## Related Endpoints
This endpoint works closely with these related endpoints:

1. `/vip/user/service/info` - Gets VIP user service information
   - Provides more detailed subscription information
   - Used to initialize the tier structure

2. `/vip/device/update/position` - Updates device positions (not shown in endpoints.yml)
   - Called after drag-and-drop operations
   - Saves the new device tier positions to the server

3. `/vip/tier/list/v4` - Gets subscription tier list
   - Provides tier pricing and feature information
   - Complements the device positioning data

4. `/vip/device/cloud/info` - Query device cloud service VIP info
   - Gets device-specific subscription details
   - May be called before or after position changes

## Notes
- This endpoint is crucial for the VIP/subscription management UI
- The response structure supports a flexible tier-based subscription model
- The UI implements an intuitive drag-and-drop model for device management
- Error code -2006 appears to be specifically for device selection issues
- The endpoint works alongside other VIP endpoints to provide a complete subscription management experience
- The deviceStorageType field differentiates between cloud and local storage devices
- Online status allows the UI to show connectivity status alongside subscription status