# Available Purchase Device Endpoint Discovery

## Endpoint Information
- **Path:** `/vip/available/purchase/device`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Gets a list of devices available for purchase with a specific subscription plan

## Request Parameters
The endpoint takes a `ProductEntry` object which extends `BaseEntry`:

```json
{
  "productId": 12345,          // Required: The subscription product ID to query
  "app": {                     // Standard BaseEntry fields
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
The endpoint returns a `PurchaseDevice` object:

```json
{
  "result": 0,                  // Result code (0 indicates success)
  "msg": "string",              // Response message
  "data": {
    "allSimThirdParty": false,  // Whether all devices are sim third-party devices
    "availableForPurchase": true, // Whether any devices are available for purchase
    "deviceList": [             // List of devices available for selection
      {
        "canChoose": true,      // Whether the device can be selected
        "deviceName": "Front Door Camera", // Device name
        "icon": "https://example.com/icon.png", // Device icon URL 
        "isCheckout": true,     // Whether device is pre-selected
        "serialNumber": "ABCD1234", // Device serial number
        "deviceType": "doorbell" // Device type identifier
      }
    ]
  }
}
```

## Code Analysis
The endpoint is implemented in the API interfaces:

```java
@POST("/vip/available/purchase/device")
Observable<PurchaseDevice> getPurchaseDevice(@Body ProductEntry productEntry);
```

Implementation in PayApiClient:
```java
@Override
public Observable<PurchaseDevice> getPurchaseDevice(ProductEntry productEntry) {
    checkNotNullParameter(productEntry, "productEntry");
    return createApiService(productEntry).getPurchaseDevice(productEntry);
}
```

Usage example:
```java
public void getPurchaseDevices(int productId, final Callback<PurchaseDevice> callback) {
    ProductEntry entry = new ProductEntry();
    entry.setProductId(productId);
    
    payApiClient.getPurchaseDevice(entry)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<PurchaseDevice>() {
            @Override
            public void onNext(PurchaseDevice response) {
                if (response != null && response.getResult() == 0) {
                    callback.onSuccess(response);
                } else {
                    callback.onFailure(response != null ? response.getMsg() : "Unknown error");
                }
            }
            
            @Override
            public void onError(Throwable e) {
                callback.onFailure("Network low");
            }
        });
}
```

## Usage Context
This endpoint is used in the following scenarios:

1. During the subscription purchase flow:
   - When a user selects a subscription plan
   - Before proceeding to payment
   - To allow selection of devices to include in the subscription

2. For subscription management:
   - When a user wants to upgrade their subscription
   - When adding devices to an existing subscription
   - When changing which devices are covered by a subscription

3. UI implementation:
   - The response populates a device selection dialog
   - Users can select/deselect devices via checkboxes
   - Device thumbnails and names are displayed
   - Selection affects the final purchase price

## UI Implementation
The app presents this functionality in the following way:

1. In PayDialogActivity:
   - Displays a list of devices in a RecyclerView
   - Each device has a checkbox for selection
   - Shows device name and thumbnail image
   - Enforces selection limits based on subscription tier
   - Updates the total device count in real-time

2. Selection process:
   - Tapping a device toggles its selection state
   - A running count shows selected/total devices
   - A "Select All" checkbox can select all available devices at once
   - Continue button becomes enabled when valid selection is made

3. Interface elements:
   - Devices have checkboxes and thumbnail images
   - Disabled devices (canChoose=false) are shown but cannot be selected
   - Selected device count is displayed (e.g., "2/5 devices")
   - Continue button proceeds to payment

## Error Handling
The application handles errors by:

1. Network errors:
   - Shows "Network low" toast message
   - Logs error events for analytics
   - Allows retry operation

2. Device availability errors:
   - If no devices are available, shows a special dialog
   - Dialog explains why no devices can be selected
   - Provides options to add devices or cancel

3. Selection validation:
   - Ensures at least one device is selected
   - Enforces maximum device limits per subscription tier
   - Shows appropriate error messages for invalid selections

## Subscription Flow Integration
This endpoint is a key part of the subscription purchase flow:

1. Flow sequence:
   - User selects a subscription tier and billing period
   - This endpoint is called with the selected productId
   - User selects devices from the returned list
   - Selected device list is passed to payment processing

2. Device limits by tier:
   - Basic tier: Limited device count (e.g., 1-2 devices)
   - Plus tier: Higher device limit (e.g., 5 devices)
   - Premium tier: Maximum device support (e.g., 10+ devices)

3. Payment integration:
   - Selected devices are stored via VIPManager.setVIPSNList()
   - Device selection affects final price calculation
   - Payment method (Google Pay, etc.) is initiated after selection

## Business Logic
The endpoint implements several business rules:

1. Device eligibility:
   - The `canChoose` flag indicates device eligibility
   - Some devices may not be eligible due to device type, firmware, or other factors
   - Pre-existing subscription coverage may make devices non-selectable

2. Pre-selection logic:
   - Some devices may be pre-selected (`isCheckout=true`)
   - These are typically the user's primary or most-used devices
   - The UI respects these pre-selections

3. Availability check:
   - The `availableForPurchase` flag indicates if any valid devices exist
   - If false, the UI shows a "no devices" dialog
   - This prevents empty subscription purchases

## Related Endpoints
This endpoint works with these related endpoints:

1. `/vip/tier/list/v4` - Gets subscription tiers
   - Provides the productId needed for this endpoint
   - Defines the tier's device limits

2. `/vip/user/device/list` - Lists devices with VIP status
   - Shows current subscription allocation
   - Used before or after this endpoint

3. `/pay/google/order/verify` - Payment verification
   - Used after device selection to process payment
   - Includes the selected device list

## Notes
- This endpoint is part of the VIP/subscription management API group
- The response enables a flexible device selection UI
- Maximum device counts depend on the subscription tier
- Both new subscriptions and upgrades use this endpoint
- Some devices may be pre-selected based on usage patterns
- The `allSimThirdParty` flag appears to indicate if all devices are SIM-enabled third-party devices
- Device selection must include at least one device
- The selection affects subscription pricing and features
- User decisions during this flow are tracked for analytics purposes