# Endpoint: /user/accountlinkedplatforms

## Purpose
This endpoint retrieves a list of third-party platforms (like voice assistants or smart home systems) that have been linked to the user's account. It's used to determine which platform-specific features should be displayed to the user and what integrations are available.

## Request Method
POST

## Request Parameters
The request body is a JSON object containing standard BaseEntry fields:

- `app` (AppBean object, optional): Application information
  - `appType` (String, optional): Type of application
  - `version` (String, optional): Application version
  - `timeZone` (String, optional): User's timezone
- `countryNo` (String, optional): User's country code
- `language` (String, optional): User's preferred language
- `tenantId` (String, optional): Tenant ID for the system

No additional parameters beyond the standard authentication token are required.

## Request Example
```json
{
  "app": {
    "appType": "Android",
    "version": "1.2.3",
    "timeZone": "America/New_York"
  },
  "countryNo": "US",
  "language": "en"
}
```

## Response
The response extends the BaseResponse object and includes a data field with a list of linked platforms:

### Success Response (result = 0)
```json
{
  "result": 0,
  "msg": "Success",
  "data": {
    "linkedPlatforms": ["Alexa", "Google Assistant"]
  }
}
```

### Error Response (result < 0)
```json
{
  "result": -1,
  "msg": "Error retrieving linked platforms"
}
```

### Empty Response (no linked platforms)
```json
{
  "result": 0,
  "msg": "Success",
  "data": {
    "linkedPlatforms": []
  }
}
```

## Error Codes
- General error codes apply (network errors, authentication errors)
- No specific error codes were identified for this endpoint

## Usage in App
This endpoint is used to determine what platform-specific features and integrations should be shown to the user:

1. The app calls the endpoint to check what platforms are linked to the user's account
2. Based on the response, the app may show platform-specific features (e.g., Alexa skills, Google Assistant actions)
3. The app uses this information to customize the user interface and available options

Specific usage examples include:
- `checkShowAlexa()` method in `AppConfigViewModel` checks for Alexa integration
- The method determines if Alexa-specific features should be displayed in the app

## Client Implementation
The client creates a BaseEntry object and calls the getAccountLinkedPlatforms method:

```java
BaseEntry baseEntry = new BaseEntry();
apiClient.getAccountLinkedPlatforms(baseEntry)
    .subscribeOn(Schedulers.io())
    .observeOn(AndroidSchedulers.mainThread())
    .subscribe(new Observer<LinkedPlatformsResponse>() {
        @Override
        public void onNext(LinkedPlatformsResponse response) {
            if (response.getResult() >= 0 && response.getData() != null) {
                List<String> platforms = response.getData().getLinkedPlatforms();
                // Check if any platforms are linked
                boolean hasLinkedPlatforms = platforms != null && platforms.size() > 0;
                // Update UI based on the result
                // ...
            }
        }
        
        @Override
        public void onError(Throwable e) {
            // Handle error
        }
    });
```

## Error Handling
In the app, error responses are handled by:
- Checking if result < 0 to identify an error condition
- The app typically proceeds with default behavior (not showing platform-specific features) in case of error

## Security Considerations
- This endpoint requires authentication
- All communication should be over HTTPS
- The information is used to determine what features are available to the user