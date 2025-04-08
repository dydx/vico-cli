# Reset Password Verification Code Endpoint Discovery

## Endpoint Information
- **Path:** `/account/resetconfirm`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Sends a verification code to the user's email or phone number as part of the password reset process. This is the first step in the password reset flow.

## Request Parameters
The endpoint takes a `VerificationCodeEntry` object in the request body which extends `BaseEntry`:

```json
{
  "email": "string",      // Required if resetting via email: User's email address
  "phone": "string",      // Required if resetting via phone: User's phone number
  "language": "string",   // Optional: User language preference
  "requestId": "string",  // Optional: Request identifier for tracking
  "app": {                // Standard BaseEntry fields
    "type": "string",     // Application type identifier
    "version": "string"   // Application version
  },
  "countryNo": "string",  // User's country code
  "tenantId": "string"    // User's tenant ID
}
```

**Note:** Either `email` or `phone` must be provided, but not both.

## Response Format
The endpoint returns a `BaseResponse` object:

```json
{
  "result": 0,     // Result code (0 indicates success)
  "msg": "string"  // Response message
}
```

## Code Analysis
The endpoint is implemented in the API interface:

```java
@POST("account/resetconfirm")
Observable<BaseResponse> sendForgotVerificationCode(@Body VerificationCodeEntry verificationCodeEntry);
```

In the `AccountViewModel` class, the method to call this endpoint is:

```java
public void getVCodeForForgotPassword(String str, Activity activity) {
    VerificationCodeEntry verificationCodeEntry = new VerificationCodeEntry();
    if (vicohome_1742553098674_0O0oo0.checkPhoneNumber(str)) {
        verificationCodeEntry.phone = str;
    } else {
        verificationCodeEntry.setEmail(str);
    }
    addSubscription(com.ai.addxnet.vicohome_1742553098674_0O0o0oo.getInstance().sendForgotVerificationCode(verificationCodeEntry, new vicohome_1742553098674_00O0o0oOO(activity, str)));
}
```

## Usage Context
This endpoint is used in the following scenarios:
1. When a user initiates the "Forgot Password" process
2. When a user wants to reset their password but doesn't remember the current one
3. As the first step in the password reset flow, before entering a new password

The typical user flow:
1. User navigates to the login screen and selects "Forgot Password"
2. User enters their email or phone number
3. App calls this endpoint to request a verification code
4. User receives the code via email or SMS
5. User enters the code in the next screen along with a new password
6. App calls `/account/resetpswd` to complete the password reset process

## Error Handling
The response follows standard error handling:
- Success is indicated by `result: 0`
- Errors are indicated by negative result codes with error messages:
  - `-1011`: Invalid email or phone format
  - `-1021`: Invalid verification code
  - `-1028`: Request too frequent
  - `-1033`: Request too frequent (alternative code)
  - `-1035`: Request too frequent (another alternative code)
  - `-1036`: Verification code expired
  - `-1001`: Email not registered
  - `-1002`: Account blocked or other security issue

Error handling in the client code:
```java
@Override
public void doOnNext(BaseResponse baseResponse) {
    int result = baseResponse.getResult();
    if (result < 0) {
        // Handle specific error codes
        if (result == -1011) {
            // Invalid email/phone format
        } else if (result == -1035 || result == -1033 || result == -1028) {
            // Too many requests
        } else if (result == -1036) {
            // Code expired
        } else {
            // Other errors
        }
    } else {
        // Success case
    }
}
```

## Related Endpoints
- `/account/resetpswd` - Completes the password reset process with the verification code and new password
- `/account/checkconfirm` - Validates the verification code
- `/account/registconfirm` - Sends a verification code for the registration process
- `/account/retryLoginCode` - Sends a verification code when a password login fails

## Notes
- The verification code is typically a 6-digit number
- The code is valid for a limited time (usually 5-10 minutes)
- There are rate limits on how frequently this endpoint can be called for security reasons
- The app should handle various error scenarios, especially for frequent requests
- This endpoint is part of a two-step password reset process:
  1. Send verification code (`/account/resetconfirm`)
  2. Verify code and set new password (`/account/resetpswd`)
- The email/phone provided must be associated with an existing account
- For security reasons, the API returns success even if the email/phone doesn't exist, but no code is sent in that case