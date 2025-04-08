# Endpoint: /device/bindComplete

## Description
This endpoint is used to notify the server that a device binding step has been completed. It's part of the device setup process and helps track the progress of binding a device to a user account.

## Request Method
POST

## Authentication
Requires user authentication

## Request Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| operationId | String | Yes | The unique operation ID for the binding process. This is obtained during the device binding initialization. |
| bindStep | int | Yes | The binding step that has been completed (values typically range from 1-4). |

### Inherited from BaseEntry
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| app | Object | No | Application information object |
| countryNo | String | No | Country code |
| language | String | No | Language setting |
| tenantId | String | No | Tenant identifier |

### App Object (AppBean)
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| apiVersion | String | No | API version |
| appName | String | No | Application name |
| appType | String | No | Application type (default: "Android") |
| bundle | String | No | Bundle identifier |
| countlyId | String | No | Countly analytics ID |
| env | String | No | Environment |
| tenantId | String | No | Tenant identifier |
| timeZone | String | No | User's timezone |
| version | int | No | Application version code |
| versionName | String | No | Application version name |

## Response

### Success Response
```json
{
  "result": 0,
  "msg": "success"
}
```

### Error Response
```json
{
  "result": -XXXX,
  "msg": "Error message"
}
```

## Usage Examples

### Code Example
```java
BindCompleteEntry bindCompleteEntry = new BindCompleteEntry(operationId, stepNumber);
api.bindComplete(bindCompleteEntry)
   .subscribeOn(Schedulers.io())
   .observeOn(AndroidSchedulers.mainThread())
   .subscribe(new Subscriber<BaseResponse>() {
      @Override
      public void doOnError(Throwable e) {
         LogUtils.d(TAG, "bindComplete error: " + e);
      }
      
      @Override
      public void doOnNext(BaseResponse baseResponse) {
         LogUtils.d(TAG, "bindComplete success");
      }
   });
```

## Notes
- This endpoint is called during the device binding process to notify the server about the completion of specific binding steps.
- The binding process typically involves multiple steps:
  1. Step 1: WiFi connection established
  2. Step 2: Device registered with cloud
  3. Step 3: Device setup completed
  4. Step 4: Initialization completed
- The operation ID is a unique identifier for the binding session, obtained during binding initialization.
- The server uses this information to track the device binding progress and provide appropriate next steps.
- After completing all steps, the device should appear in the user's device list.