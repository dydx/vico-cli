# Library Payment Banner Endpoint Discovery

## Endpoint Information
- **Path:** `/library/freeuser/banner`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Query payment promotion banner information for free users in the library section

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
The endpoint returns a `LibraryPayBannerResponse` object:

```json
{
  "result": 0,                  // Result code (0 indicates success)
  "msg": "string",              // Response message
  "data": {
    "notify": true,             // Whether to show the banner
    "notifyCount": 30,          // Value related to notification (days/GB)
    "notifyMessage": "string",  // Banner message text to display
    "notifyType": 0,            // Type of notification (0=days, 1=GB)
    "slotName": "string",       // Banner slot identifier
    "upgradeFreeBtn": "string", // Text for the banner button
    "recommendProduct": {       // Recommended product information
      "productId": 123,         // Unique product ID for subscription
      "productName": null,      // Product name (may be null)
      "maxDeviceNum": null,     // Maximum number of devices supported
      "subscriptionGroupId": "" // Subscription group identifier
    },
    "recommendProductDO": {     // More detailed product recommendations
      "defaultSelectedProductId": 123,   // Default selected product
      "buttonFeatureGroup": "string",    // Feature grouping for button
      "yearFeatureGroup": "string",      // Feature grouping for yearly plan
      "freeTrialDay": 7,                 // Free trial period in days
      "recommendProductList": [],        // Monthly subscription options
      "recommendProductYearList": []     // Yearly subscription options
    },
    "recommendProductV1": []    // List of recommended product options
  }
}
```

## Related Endpoint
The `/library/freeuser/banner/close` endpoint is called when the banner is closed or when a user interacts with it:

```json
// Request: BaseEntry (same as above)
// Response: BaseResponse with standard result code and message
{
  "result": 0,
  "msg": "success"
}
```

## Code Analysis
The endpoint is implemented in the API interfaces:

```java
@POST("/library/freeuser/banner")
Observable<LibraryPayBannerResponse> queryPaymentBanner(@Body BaseEntry baseEntry);

@POST("/library/freeuser/banner/close")
Observable<BaseResponse> reportBannerShown(@Body BaseEntry baseEntry);
```

The banner information is processed in the `LibraryBannerHelper` class:

```java
public final void initPayBanner() {
    FragmentActivity requireActivity = this.fragment.requireActivity();
    launch(LifecycleOwnerKt.getLifecycleScope(this.fragment), new LibraryBannerHelper$initPayBanner$1(this, requireActivity, null));
}

// Called inside the coroutine:
Observable<LibraryPayBannerResponse> queryPaymentBanner = ApiClient.getInstance().queryPaymentBanner(new BaseEntry());
LibraryPayBannerResponse response = await(queryPaymentBanner);

if (checkPaybannerResponse(response)) {
    // Display the banner
    this.notifyType = response.getData().getNotifyType() != null ? response.getData().getNotifyType() : -1;
    // Set up banner UI with response data
    view.findViewById(R.id.pay_banner_title).setText(response.getData().getNotifyMessage());
    ((TextView) view.findViewById(R.id.pay_banner_btn)).setText(response.getData().getUpgradeFreeBtn());
    // Set up click handlers
}
```

When the banner is closed or hidden:
```java
// In hidePayBanner method:
Observable<BaseResponse> reportBannerShown = ApiClient.getInstance().reportBannerShown(new BaseEntry());
await(reportBannerShown);
```

## Usage Context
This endpoint is used to display payment promotion banners to free users in the Library/Recordings section. The system:

1. Calls the endpoint when a free user accesses the library section
2. Displays a promotion banner if `notify` is true
3. Shows the notification message from `notifyMessage`
4. Uses the button text from `upgradeFreeBtn`
5. Reports when the banner is shown or interacted with

## Banner Types
The response includes a `notifyType` field which appears to distinguish different banner messages:
- Type 0: Storage capacity warning in days (e.g., "X days of recordings")
- Type 1: Storage capacity warning in GB (e.g., "X GB of storage space")

The `notifyCount` field provides the corresponding numeric value.

## Recommended Products
The banner includes subscription product recommendations:
1. `recommendProduct`: The primary recommended subscription
2. `recommendProductDO`: Detailed subscription options with:
   - Monthly subscription options (`recommendProductList`)
   - Yearly subscription options (`recommendProductYearList`)
   - Free trial period information (`freeTrialDay`)
   - Default selected product (`defaultSelectedProductId`)

## User Flow
The typical user flow for this feature is:
1. Free user opens the library/recordings section
2. App calls `/library/freeuser/banner` endpoint
3. If `notify` is true, banner is displayed with the message and button
4. User can:
   - Click the close button (reports using `/library/freeuser/banner/close`)
   - Click the upgrade button (opens subscription flow)
5. If user subscribes, the banner is no longer shown

## Error Handling
The code includes validation to ensure the banner is shown only when necessary:
```java
if (!LibraryBannerHelper.access$checkPaybannerResponse(this.this$0, libraryPayBannerResponse)) {
    return Unit.INSTANCE;
}
```

The validation checks:
1. Response is not null and result code is 0
2. Data object is not null
3. Notify flag is true
4. For non-China regions, product information must be valid

## Business Logic
This endpoint supports the conversion of free users to paid subscribers by:
1. Targeting users in the library section where storage limitations are most apparent
2. Providing contextual messages about storage limitations
3. Offering appropriate subscription options
4. Making subscription sign-up accessible with one click
5. Tracking banner impressions and interactions

## Notes
- The banner is specifically targeted at free users in the library section
- The implementation includes animation when showing/hiding the banner
- The endpoint supports regional differences (special handling for China payments)
- Banner content appears to be dynamic and controlled by the server
- The banner can promote both monthly and yearly subscription options
- The system tracks banner impressions and user interactions