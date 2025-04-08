# Quick Login Code API

Gets a verification code for quick login via SMS.

**URL** : `/account/quickLoginCode`

**Method** : `POST`

**Auth required** : No

**Permissions required** : None

## Request Parameters

```json
{
  "phone": "[user phone number]",
  "loginType": 2
}
```

### Required Fields

* `phone` - User's phone number

### Optional Fields

* `loginType` - Type of login (2 represents phone login)

## Success Response

**Code** : `200 OK`

**Content example**

```json
{
  "result": 0,
  "msg": "success"
}
```

## Error Responses

**Condition** : If phone number is invalid

**Code** : `200 OK`

**Content** :

```json
{
  "result": -1011,
  "msg": "Invalid phone number"
}
```

**Condition** : If request is too frequent

**Content** :

```json
{
  "result": -1035,
  "msg": "Request too frequent"
}
```

**Condition** : Server error or network issue

**Content** :

```json
{
  "result": -1,
  "msg": "Error"
}
```

## Notes

* This endpoint is used when a user opts to login with a verification code instead of password
* Once the code is received by SMS, it should be used with the `/account/login` endpoint to complete authentication
* The verification code has a limited validity period (typically 120 seconds)
* After requesting a code, there is a cooldown period before another code can be requested for the same phone number

## Client Implementation

1. Client calls this endpoint with the user's phone number
2. Server validates the phone number and sends a verification code via SMS
3. User enters the received code in the app
4. Client uses the code in a subsequent login request (`/account/login`)

## Error Codes

* `-1011`: Invalid phone number
* `-1035`: Request too frequent
* `-1033`: Request too frequent (alternative)
* `-1028`: Request too frequent (alternative)
* `-1026`: Invalid verification code
* `-1036`: Verification code expired