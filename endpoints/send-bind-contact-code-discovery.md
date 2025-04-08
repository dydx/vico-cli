# Send Bind Contact Code Endpoint Discovery

## Endpoint Information
- **Path:** `/account/sendbindcontactcode`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Sends a verification code to the provided phone number or email address for binding contact information to an account.

## Request Parameters
The endpoint takes a `VerificationCodeEntry` object in the request body which extends `BaseEntry`:

```json
{
  "phone": "string",        // Required if binding a phone number (mutually exclusive with email)
  "email": "string",        // Required if binding an email address (mutually exclusive with phone)
  "loginType": 1|2,         // Optional: 1 for email verification, 2 for phone verification
  "language": "string",     // Optional: User language preference
  "requestId": "string",    // Optional: Request identifier for tracking
  "app": {                  // Standard BaseEntry fields
    "type": "string",       // Application type identifier
    "version": "string"     // Application version
  },
  "countryNo": "string",    // User's country code
  "tenantId": "string"      // User's tenant ID
}
```

## Response Format
The endpoint returns a `BaseResponse` object:

```json
{
  "result": 0,              // Result code (0 or positive indicates success, negative indicates error)
  "msg": "string"           // Response message
}
```

## Code Analysis
The endpoint is implemented in interface `com.ai.addxnet.vicohome_1742553098674_00O0o0oOO`:

Example call signature from decompiled code:
```java
@POST("/account/sendbindcontactcode")
Observable<BaseResponse> getVCodeForBindPhoneEmail(@Body VerificationCodeEntry verificationCodeEntry);
```

The `AccountViewModel` class provides a higher-level implementation:
```java
public void getVCodeForPhoneEmail(String str, Activity activity) {
    VerificationCodeEntry verificationCodeEntry = new VerificationCodeEntry();
    if (vicohome_1742553098674_0O0oo0.checkPhoneNumber(str)) {
        verificationCodeEntry.phone = str;
    } else {
        verificationCodeEntry.setEmail(str);
    }
    addSubscription(com.ai.addxnet.vicohome_1742553098674_0O0o0oo.getInstance().getVCodeForBindPhoneEmail(verificationCodeEntry)
        .subscribeOn(Schedulers.io())
        .observeOn(AndroidSchedulers.mainThread())
        .subscribe((Subscriber<? super BaseResponse>) new vicohome_1742553098674_0O0o0oo(activity, str)));
}
```

## Usage Context
This endpoint is used in the following scenarios:
1. When a user wants to bind a phone number to their account
2. When a user wants to bind an email address to their account
3. When a user wants to change their current phone number or email address

The typical user flow:
1. User navigates to account settings → bind phone/email
2. User enters the phone number or email to bind
3. App calls this endpoint to send a verification code
4. User receives the code via SMS or email
5. User enters the code for verification
6. App validates the code using the `/account/checkbindcontact` endpoint

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0` or positive value
- Errors are indicated by negative result codes with corresponding messages:
  - `-1011`: Invalid phone number or email format
  - `-1021`: Password error
  - `-1036`: Verification code expired
  - `-1035`, `-1033`, `-1028`: Request too frequent
  - `-1034`, `-1032`, `-1031`: Invalid verification code
  - `-1027`: Verification code required
  - `-1026`: Password/code error

## Related Endpoints
- `/account/checkbindcontact` - Verifies the code for binding phone/email
- `/account/quickLoginCode` - Sends verification code for quick login
- `/account/retryLoginCode` - Resends verification code when login password is wrong
- `/account/registconfirm` - Sends verification code for registration
- `/account/resetconfirm` - Sends verification code for password reset

## Notes
- The endpoint determines whether to send an SMS or email based on whether `phone` or `email` is provided
- The verification code typically expires after a certain period (usually 5-10 minutes)
- There are rate limits for sending verification codes to prevent abuse
- When a verification code is successfully sent, users should be prompted to check their messages or email
- The application should handle potential delays in code delivery
- This verification process is a security measure to ensure that the user owns the phone number or email address they're trying to bind