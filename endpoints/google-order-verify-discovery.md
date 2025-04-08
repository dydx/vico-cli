# Google Order Verify Endpoint Discovery

## Endpoint Information
- **Path:** `/pay/google/order/verify`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Verifies Google Play subscription purchases and activates service

## Request Parameters
The endpoint takes an `OrderEntry` object which extends `BaseEntry`:

```json
{
  "productId": 12345,            // Required: Subscription product ID being purchased
  "purchaseToken": "string",     // Required: Token from Google Play for verification
  "subscriptionGroupId": "string", // Required: ID of the subscription group
  "outTradeNo": "string",        // Required: Google's order ID for the transaction
  "guidanceSource": 1,           // Required: Source of purchase flow (e.g., 1=device detail, 2=account)
  "tradeType": 2,                // Required: Type of transaction (1=one-time, 2=subscription)
  "tierDeviceList": [            // Required: List of device serial numbers selected for subscription
    "SN12345678",
    "SN87654321"
  ],
  "app": {                       // Standard BaseEntry fields
    "apiVersion": "string",      // API version
    "appName": "string",         // Application name
    "appType": "string",         // Application type
    "bundle": "string",          // Bundle identifier
    "countlyId": "string",       // Analytics ID
    "env": "string",             // Environment (e.g., "prod")
    "tenantId": "string",        // Tenant identifier
    "timeZone": "string",        // User's timezone
    "version": "string",         // App version number
    "versionName": "string"      // App version name
  },
  "countryNo": "string",         // Country code
  "language": "string",          // User's language preference
  "tenantId": "string"           // User's tenant ID
}
```

## Response Format
The endpoint returns a `PayResultResponse` object:

```json
{
  "result": 0,                  // Result code (0 indicates success)
  "msg": "string",              // Response message
  "data": "string"              // Additional response data (usually empty or confirmation token)
}
```

## Error Codes
The endpoint can return various error codes:

- `0`: Success - Payment verified and subscription activated
- `3010`: General payment verification failure
- `-1`: Generic error condition
- Error subcodes (1-4) may be present in the message for specific error types:
  - `1`: Network errors
  - `2`: Payment processing failures
  - `3`: Verification failures
  - `4`: Subscription activation failures

## Code Analysis
The endpoint is implemented in the API interfaces:

```java
@POST("/pay/google/order/verify")
Observable<PayResultResponse> uploadGooglePayInfo(@Body OrderEntry orderEntry);
```

Implementation in PayApiClient:
```java
@Override
public Observable<PayResultResponse> uploadGooglePayInfo(OrderEntry orderEntry) {
    checkNotNullParameter(orderEntry, "orderEntry");
    return createApiService(orderEntry).uploadGooglePayInfo(orderEntry);
}
```

Usage example in GoogleProcessManager:
```java
// In GoogleProcessManager.verifyOrder method
public void verifyOrder(Purchase purchase, PaymentProduct product) {
    OrderEntry orderEntry = new OrderEntry();
    orderEntry.setProductId(product.getProductId());
    orderEntry.setPurchaseToken(purchase.getPurchaseToken());
    orderEntry.setSubscriptionGroupId(product.getSubscriptionGroupId());
    orderEntry.setOutTradeNo(purchase.getOrderId());
    orderEntry.setGuidanceSource(mGuidanceSource);
    orderEntry.setTradeType(2); // Subscription
    orderEntry.setTierDeviceList(VIPManager.getInstance().getVIPSNList());
    
    payApiClient.uploadGooglePayInfo(orderEntry)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<PayResultResponse>() {
            @Override
            public void onNext(PayResultResponse response) {
                if (response != null && response.getResult() == 0) {
                    // Success handling
                    showSuccessDialog();
                    refreshVipStatus();
                    trackPaymentSuccess();
                } else {
                    // Error handling
                    String errorMsg = response != null ? response.getMsg() : "Unknown error";
                    showErrorDialog(errorMsg);
                    trackPaymentFailure(errorMsg);
                }
            }
            
            @Override
            public void onError(Throwable e) {
                // Network error handling
                showNetworkErrorDialog();
                trackPaymentFailure("Network error");
            }
        });
}
```

## Usage Context
This endpoint is used in the following scenario:

1. Google Play Subscription Flow:
   - After user selects a subscription tier
   - After selecting devices to include in subscription
   - After completing payment through Google Play Billing
   - Before activating subscription features in the app

2. The verification process:
   - Sends Google Play purchase details to VicoHome server
   - Server validates the purchase with Google's servers
   - Server activates the subscription for the user's account
   - Server associates selected devices with the subscription

3. Post-verification actions:
   - App shows success dialog
   - Account status is refreshed
   - Device lists are updated
   - Subscription features are enabled

## UI Implementation
The app presents this functionality in the following way:

1. During verification:
   - Shows a loading dialog ("Verifying payment...")
   - Prevents user interaction until verification completes

2. On success:
   - Dismisses loading dialog
   - Shows "Plan Activated" success dialog
   - Updates UI to reflect new subscription status
   - Enables premium features for selected devices

3. On failure:
   - Shows appropriate error dialog based on error type
   - Offers retry options for network errors
   - Provides clear error messages for payment failures
   - Tracks failure reasons for analytics

## Error Handling
The application handles errors by:

1. Checking response codes:
   - Success (result = 0): Shows success dialog
   - Failure (result != 0): Shows error dialog

2. Network error handling:
   - Catches and processes Throwable from RxJava
   - Shows network error dialog with retry option
   - Logs failure for analytics

3. Payment verification errors:
   - Shows specific error messages based on error code
   - Offers appropriate recovery options based on error type
   - Uses subcodes to determine specific failure reason

4. Recovery options:
   - "Try Again" button for network errors
   - Clear error messaging for payment failures
   - Option to contact support for persistent issues

## Payment Flow Integration
The endpoint is a critical part of the payment verification flow:

1. Pre-verification:
   - User selects subscription plan
   - User selects devices to include
   - Payment is processed through Google Play

2. Verification (this endpoint):
   - App sends purchase details to server
   - Server validates with Google Play
   - Server activates subscription

3. Post-verification:
   - App updates account status
   - UI reflects subscription changes
   - Premium features are enabled
   - Analytics events are logged

## Business Logic
The verification process implements important business rules:

1. Purchase validation:
   - Ensures purchase is legitimate
   - Confirms correct product was purchased
   - Validates purchase token with Google

2. Device association:
   - Links selected devices to subscription
   - Respects device limits for tier
   - Enables cloud storage for devices

3. Subscription activation:
   - Sets correct expiration date
   - Applies appropriate feature set
   - Records billing information

## Related Endpoints
This endpoint works with these related endpoints:

1. `/vip/tier/list/v4` - Gets subscription tiers
   - Provides subscription options before payment
   - Determines pricing and features

2. `/vip/available/purchase/device` - Gets purchasable devices
   - Used before payment to select devices
   - Determines which devices can be included

3. `/vip/user/service/info` - Gets VIP service information
   - Used after verification to confirm subscription status
   - Shows updated device allocations and features

## Notes
- This endpoint is specifically for Google Play subscription verification
- There would likely be similar endpoints for other payment platforms (Apple, etc.)
- The verification process is security-critical for preventing subscription fraud
- The `guidanceSource` parameter tracks where in the app the purchase was initiated
- The `tradeType` parameter distinguishes between one-time purchases and subscriptions
- The integration follows Google's recommended server-side verification pattern
- The process involves client-side, server-side, and Google Play server coordination
- Analytics tracking is important for monitoring conversion rates and failures