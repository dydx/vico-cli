# Send Registration Verification Code

Sends a verification code to the provided email or phone number for user registration.

## Endpoint

`POST /account/registconfirm`

## Request Parameters

| Parameter   | Type   | Required | Description                             |
|-------------|--------|----------|-----------------------------------------|
| email       | String | No*      | User's email address (if registering with email) |
| phone       | String | No*      | User's phone number (if registering with phone)  |
| loginType   | Integer| Yes      | Login type (1 for email, 2 for phone)    |
| app         | Object | Yes      | Application information                  |
| - appName   | String | Yes      | Name of the application                  |
| - appType   | String | Yes      | Type of application (default: "Android") |
| - version   | Integer| Yes      | Version code                             |
| - versionName| String| Yes      | Version name                             |
| - apiVersion| String | No       | API version                              |
| - bundle    | String | No       | Bundle identifier                        |
| - countlyId | String | No       | Analytics ID                             |
| - env       | String | No       | Environment                              |
| - tenantId  | String | No       | Tenant identifier                        |
| - timeZone  | String | No       | User's timezone                          |
| countryNo   | String | No       | Country code                             |
| language    | String | No       | Language setting                         |
| tenantId    | String | No       | Tenant identifier                        |

*Either email or phone must be provided

## Response

| Field  | Type    | Description                                      |
|--------|---------|--------------------------------------------------|
| result | Integer | Status code (negative values indicate errors)    |
| msg    | String  | Response message                                 |

## Success Response

```json
{
  "result": 0,
  "msg": "Success"
}
```

## Error Codes

| Error Code | Description                                      |
|------------|--------------------------------------------------|
| -1036      | Verification code expired                        |
| -1035      | Request too frequent                             |
| -1034      | Invalid verification code                        |
| -1033      | Request too frequent                             |
| -1032      | Invalid verification code                        |
| -1031      | Invalid verification code                        |
| -1028      | Request too frequent                             |
| -1027      | Verification code required                       |
| -1026      | Password/code error                              |
| -1021      | Password error                                   |
| -1012      | Password or code invalid                         |
| -1011      | Phone number or email invalid                    |
| -1002      | Account-related error                            |
| -1001      | Email not registered                             |

## Sample Request

```json
{
  "email": "user@example.com",
  "loginType": 1,
  "app": {
    "appName": "VicoHome",
    "appType": "Android",
    "version": 123,
    "versionName": "1.2.3"
  },
  "countryNo": "US",
  "language": "en"
}
```

## Notes

- This endpoint triggers the sending of a verification code to the provided email or phone number
- The verification code is then used in subsequent registration steps
- Either an email or phone number must be provided, but not both
- Request frequency is limited to prevent abuse