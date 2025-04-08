# Report Banner Shown Endpoint Discovery

## Endpoint Information
- **Path:** `/library/freeuser/banner/close`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Report when a payment promotion banner has been shown or closed by the user

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
The endpoint returns a standard `BaseResponse` object:

```json
{
  "result": 0,                // Result code (0 indicates success)
  "msg": "success"            // Response message
}
```

## Code Analysis
The endpoint is implemented in the API interface:

```java
@POST("/library/freeuser/banner/close")
Observable<BaseResponse> reportBannerShown(@Body BaseEntry baseEntry);
```

The endpoint is called in the `LibraryBannerHelper` class when hiding the payment banner:

```java
// Inside LibraryBannerHelper$hidePayBanner$1 method:
GlobalSwap globalSwap = GlobalSwap.vicohome_1742553098674_0o00OOoOo;
Observable<BaseResponse> reportBannerShown = com.ai.addxnet.vicohome_1742553098674_0O0o0oo.getInstance().reportBannerShown(new BaseEntry());
kotlin.jvm.internal.vicohome_1742553098674_0O0oOoo.checkNotNullExpressionValue(reportBannerShown, "getInstance().reportBannerShown(BaseEntry())");
this.label = 1;
if (globalSwap.await(reportBannerShown, this) == coroutine_suspended) {
    return coroutine_suspended;
}
```

## Usage Context
This endpoint is used in tandem with the `/library/freeuser/banner` endpoint as part of the user flow for payment promotion banners:

1. The app first calls `/library/freeuser/banner` to determine if a promotional banner should be shown
2. If a banner is displayed to the user, the app later calls `/library/freeuser/banner/close` when:
   - The user dismisses the banner by clicking the close button
   - The banner is hidden automatically (e.g., after timing out)
   - The user navigates away from the library section

## Implementation Details
The endpoint is called as part of an animation sequence that hides the banner with a smooth transition:

1. The banner's height is animated from its current height to zero
2. During this animation, the API call to `/library/freeuser/banner/close` is made
3. After the animation completes, the banner view is completely removed from the container

```java
// Animation code from LibraryBannerHelper$hidePayBanner$1:
ValueAnimator ofFloat = ValueAnimator.ofFloat(1.0f, 0.0f);
ofFloat.setDuration(300L);
ofFloat.addUpdateListener(new ValueAnimator.AnimatorUpdateListener() {
    @Override
    public final void onAnimationUpdate(ValueAnimator valueAnimator) {
        // Gradually reduce the banner height
    }
});
ofFloat.addListener(new AnimatorListenerAdapter() {
    @Override
    public void onAnimationEnd(Animator animation, boolean z) {
        // Remove the banner from view hierarchy
        container.removeView(banner);
    }
});
ofFloat.start();

// Make the API call
Observable<BaseResponse> reportBannerShown = ApiClient.getInstance().reportBannerShown(new BaseEntry());
// Wait for response using coroutine
```

## Error Handling
There is minimal error handling for this endpoint. The code makes the API call but does not specifically check the response result. This suggests that:

1. The operation is primarily for analytics/reporting purposes
2. The app does not depend on a successful response to continue functioning
3. Any network failures would be handled by global error handling mechanisms

## Business Logic
This endpoint serves several business purposes:

1. **Analytics tracking**: Records when users see promotion banners
2. **UX optimization**: Helps track banner effectiveness and user engagement
3. **Server-side logic**: May be used to limit the frequency of banner displays
4. **Conversion tracking**: Contributes to analyzing the payment conversion funnel

## Related Endpoints
This endpoint works together with the `/library/freeuser/banner` endpoint:

1. `/library/freeuser/banner` - Determines if a banner should be displayed
2. `/library/freeuser/banner/close` - Reports when the banner is shown/closed

These endpoints are part of the subscription promotion system targeted at free users.

## Notes
- The endpoint call is made regardless of whether the user manually closes the banner or it's hidden automatically
- The API call is fire-and-forget; the app doesn't wait for a successful response before continuing
- No additional parameters are sent beyond the standard BaseEntry, suggesting the server identifies the user and context from this information alone
- This is primarily an analytics endpoint rather than one that significantly affects app functionality