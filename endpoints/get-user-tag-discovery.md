# API Endpoint Discovery: /api/userTag

## Overview
The `/api/userTag` endpoint is used to retrieve user-specific bucket parameters (tags) for A/B testing and feature flag functionality. It's implemented through the GrowthBook SDK for managing features and experiments.

## Request

### HTTP Method
POST

### Request Parameters
The endpoint accepts a `BaseEntry` object in the request body with the following structure:

```json
{
  "app": {
    "apiVersion": "string",
    "appName": "string", 
    "appType": "Android",
    "bundle": "string",
    "countlyId": "string",
    "env": "string",
    "tenantId": "string",
    "timeZone": "string", 
    "version": 0,
    "versionName": "string"
  },
  "countryNo": "string",
  "language": "string",
  "tenantId": "string"
}
```

## Response

### Success Response
The endpoint returns a `BucketResponseModel` which extends `BaseResponse`:

```json
{
  "result": 0,  // 0 indicates success
  "msg": "string",
  "data": {
    // Key-value pairs of user tags/attributes
    "key1": "value1",
    "key2": "value2"
    // Additional properties as needed
  }
}
```

The `data` field contains a HashMap of string keys to object values that represent user-specific attributes or tags that are used for A/B testing, feature flagging, and personalization.

### Error Response
In case of an error, the response maintains the BaseResponse structure:

```json
{
  "result": -1,  // Non-zero indicates error
  "msg": "Error message"
}
```

## Error Handling
Error handling is implemented in the client code through RxJava Observables. When an error occurs:

1. The error code and message are logged
2. The default `doOnError` method from `vicohome_1742553098674_0O0oO000` is called with the error code and message

## Usage Context
The endpoint is primarily used in the A/B testing implementation to:

1. Fetch user-specific attributes (tags) from the server
2. Save these attributes locally via `ABTestSharePreferencesManager`
3. Update the GrowthBook SDK with the retrieved attributes
4. Enable feature flag checking and A/B testing based on these attributes

The response data is used to determine which features or variations a user should see based on their assigned bucket or segment.

## Notes
- The response data is stored locally for offline access
- The attributes are merged with locally cached attributes when setting up GrowthBook SDK
- This endpoint is called during the initialization of the A/B testing module