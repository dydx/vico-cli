# Endpoint: /user/setusertoken

## Purpose
This endpoint registers or updates a device's push notification token for a user's account. It associates the token with the user so that the server can send push notifications to their device for events like motion detection, person detection, and other camera alerts.

## Request Method
POST

## Request Parameters
The request body is a JSON object containing:

### TokenEntry object
- `msgType` (String, required): The type of push notification service (e.g., "FCM" for Firebase Cloud Messaging, "APNS" for Apple Push Notification Service)
- `msgToken` (String, required): The actual push notification token obtained from the device's push notification service
- Plus standard BaseEntry fields:
  - `app` (AppBean object): Application information
  - `countryNo` (String, optional): Country code
  - `language` (String, optional): Language preference
  - `tenantId` (String, optional): Tenant ID for the system

## Request Example
```json
{
  "msgType": "FCM",
  "msgToken": "cNR9qIzSSomwn8HJLZlmN9:APA91bEYH6DSFvmZ5XGDvD1MEGNs...",
  "app": {
    "type": "android",
    "version": "1.2.3"
  },
  "countryNo": "US",
  "language": "en"
}
```

## Response
The response is a standard BaseResponse object:

### Success Response (result = 0)
```json
{
  "result": 0,
  "msg": "Success"
}
```

### Error Response (result < 0)
```json
{
  "result": -1,
  "msg": "Error setting push token"
}
```

## Error Codes
- General error codes apply (network errors, authentication errors)
- No specific error codes were identified for this endpoint

## Usage in App
This endpoint is used in the push notification system:

1. When the device receives a new push notification token from FCM/APNS
2. During the initial device setup process
3. When the token needs to be refreshed or updated
4. After app reinstallation or token invalidation

The typical flow is:
1. Device obtains a push notification token from Firebase Cloud Messaging (Android) or APNS (iOS)
2. App calls this endpoint to associate the token with the user's account on the server
3. Server stores the token and uses it to send notifications for camera events

## Client Implementation
The client implementation is primarily in the `A4xPushBaseManager` class:

```java
// In A4xPushBaseManager
private void pushTokenToServer(String token) {
    TokenEntry tokenEntry = new TokenEntry();
    tokenEntry.setMsgType("FCM"); // or "APNS" for iOS
    tokenEntry.setMsgToken(token);
    
    apiClient.setPushToken(tokenEntry)
        .subscribeOn(Schedulers.io())
        .observeOn(AndroidSchedulers.mainThread())
        .subscribe(new Observer<BaseResponse>() {
            @Override
            public void onNext(BaseResponse response) {
                if (response.getResult() >= 0) {
                    // Success handling
                    updateTokenStatus(Step.FINISH);
                    onPushTokenToServerSuccess();
                } else {
                    // Error handling
                    updateTokenStatus(Step.FAILED);
                    onPushTokenToServerError();
                }
            }
            
            @Override
            public void onError(Throwable e) {
                // Network error handling
                updateTokenStatus(Step.FAILED);
                onPushTokenToServerError();
            }
        });
}
```

## Error Handling
In the app, error responses are handled by:
- Tracking the push token registration status with an enumeration (`Step`)
- Implementing a retry mechanism for errors
- Logging errors for debugging purposes
- Notifying the application about token registration status

## Security Considerations
- This endpoint handles device notification tokens
- Requires proper authentication token
- All communication should be over HTTPS
- The push token allows the server to send notifications to the user's device, so it's important for security and privacy