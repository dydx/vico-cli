# VIP User Rolling Day Endpoint Discovery

## Endpoint Information
- **Path:** `/vip/user/rolling/day`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Updates a user's subscription billing cycle day

## Request Parameters
The endpoint takes a `RollingDayEntry` object which extends `BaseEntry`:

```json
{
  "rollingDay": 15,            // Required: The day of month to set as billing day (1-31)
  "userVipId": 12345,          // Required: The user's VIP subscription ID
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
The endpoint returns a standard `BaseResponse` object:

```json
{
  "result": 0,                  // Result code (0 indicates success)
  "msg": "string"               // Response message
}
```

## Code Analysis
The endpoint is implemented in the API interfaces:

```java
@POST("/vip/user/rolling/day")
Observable<BaseResponse> updateUserRollingDay(@Body RollingDayEntry rollingDayEntry);
```

Implementation in PayApiClient:
```java
@Override
public Observable<BaseResponse> updateUserRollingDay(RollingDayEntry rollingDayEntry) {
    checkNotNullParameter(rollingDayEntry, "rollingDayEntry");
    return createApiService(rollingDayEntry).updateUserRollingDay(rollingDayEntry);
}
```

Usage example:
```java
public void updateRollingDay(int day, int userVipId, final Callback<BaseResponse> callback) {
    RollingDayEntry entry = new RollingDayEntry();
    entry.setRollingDay(day);
    entry.setUserVipId(userVipId);
    
    payApiClient.updateUserRollingDay(entry)
        .compose(ApplySchedulers.io_main())
        .subscribe(new Observer<BaseResponse>() {
            @Override
            public void onNext(BaseResponse response) {
                if (response != null && response.getResult() == 0) {
                    callback.onSuccess(response);
                } else {
                    callback.onFailure(response != null ? response.getMsg() : "Unknown error");
                }
            }
            
            @Override
            public void onError(Throwable e) {
                callback.onFailure("open_fail_retry");
            }
        });
}
```

## Usage Context
This endpoint is used in the following scenario:

1. In the VIP service management screen:
   - To allow users to change when their subscription is billed each month
   - To select from a set of predefined days (typically 1, 5, 10, 15, 20, 25)
   - To align billing with user's preferred financial schedule (e.g., payday)

2. The rolling day is important because:
   - It determines when subscription charges occur
   - It may affect the subscription renewal date
   - It helps users manage their monthly expenses
   - It can reduce the chance of declined payments

## UI Implementation
The app presents this functionality in the following way:

1. In VipServiceActivity:
   - The current rolling day is displayed in a setting item
   - Labeled as "Video History" with the current day value 
   - Clicking this item opens a bottom dialog

2. Dialog implementation:
   - Shows a list of available days from the `supportRollingDay` response field
   - Presents each option as "Video History X Days"
   - Highlights the currently selected option
   - Allows selecting a new option

3. After selection:
   - The API call is made with the new rolling day value
   - On success, the UI updates to show the new value
   - On failure, shows an error message

## Error Handling
The application handles errors by:

1. Checking for null responses or non-zero result codes
2. Displaying a toast message "open_fail_retry" on failures
3. Not updating the UI if the server update fails
4. Using standard RxJava error handling patterns

## Available Rolling Days
The available billing cycle days are:

1. Provided by the `/vip/user/service/info` endpoint
2. Typically include days 1, 5, 10, 15, 20, and 25 of each month
3. The current rolling day is also returned in that response
4. Only days returned in the `supportRollingDay` list can be selected

## Business Logic
The rolling day functionality serves several purposes:

1. Customer flexibility:
   - Allows users to align subscription billing with their income schedule
   - Helps users manage monthly expenses by selecting preferred payment dates

2. Payment success optimization:
   - Can reduce failed payments by letting users choose when they have funds available
   - May improve retention by reducing involuntary churn

3. User convenience:
   - Single interface to change billing date without contacting customer service
   - Immediate feedback on whether the change was successful

## Related Endpoints
This endpoint works with these related endpoints:

1. `/vip/user/service/info` - Gets current rolling day and available options
   - Returns the `rollingDay` field with current value
   - Returns the `supportRollingDay` list of available options

2. `/vip/tier/list/v4` - Gets subscription tier information
   - Provides context about the subscription being modified

3. `/vip/device/cloud/info` - Gets device VIP information
   - Provides subscription details that are affected by billing date

## Notes
- This endpoint is part of the VIP/subscription management API group
- The API allows only selecting from predefined days, not arbitrary dates
- The change is immediate and affects the next billing cycle
- The userVipId is required to identify which subscription to modify
- The endpoint only updates the billing day, not other subscription parameters
- This allows a level of customization in the subscription management
- There is no validation in the API that the selected day exists in all months (e.g., 31)
- The typical offered days (1, 5, 10, 15, 20, 25) avoid this issue by being available in all months