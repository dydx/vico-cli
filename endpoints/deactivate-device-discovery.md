# Deactivate Device Endpoint Discovery

## Endpoint Information
- **Path:** `/device/deactivatedevice`
- **Method:** POST
- **Base URL:** https://api-us.vicohome.io
- **Description:** Unbinds (deactivates) a device from a user's account

## Request Parameters
The endpoint takes a `SerialNoEntry` object in the request body, which extends `BaseEntry`:

```json
{
  "serialNumber": "string",  // The serial number of the device to unbind
  "app": {
    "type": "string",      // Application type identifier
    "version": "string"    // Application version
  },
  "countryNo": "string",   // User's country code
  "language": "string",    // User's language preference
  "tenantId": "string"     // User's tenant ID
}
```

## Response Format
The endpoint returns a `BaseResponse` object:

```json
{
  "result": 0,             // Result code (0 indicates success, other values indicate errors)
  "msg": "string"          // Response message
}
```

## Code Analysis
The endpoint is implemented in two different interfaces:
1. `com.ai.addxnet.vicohome_1742553098674_00O0o0oOO` (line 765)
2. `com.ai.settingcore.vicohome_1742553098674_00O0o0oOO` (line 134)

Example call signature from decompiled code:
```java
@POST("device/deactivatedevice")
Observable<BaseResponse> unbindDevice(@Body SerialNoEntry serialNoEntry);
```

The main implementation is found in `DeviceSettingApiClient.java`:
```java
@Override
public Observable<BaseResponse> unbindDevice(SerialNoEntry serialNoEntry) {
    com.ai.addxnet.vicohome_1742553098674_0O0o0oo.getInstance().wrapBaseEntry(serialNoEntry);
    return this.vicohome_1742553098674_0o00OOoOo.unbindDevice(serialNoEntry);
}
```

## Usage Context
This endpoint is called when:
1. A user wants to remove a device from their account
2. During device transfer between accounts
3. When troubleshooting device issues that require rebinding

Main client calls are found in:
1. `DeviceManageCore.java` - Handles the API call and processes responses 
2. `ADDXBind.java` - Provides a high-level unBindDevice method with callback handling

Error handling example from `ADDXBind.java`:
```java
public final void unBindDevice(String sn, vicohome_1742553098674_0oOOoOO callback) {
    if (callback != null) {
        callback.onStartUnBind(sn);
    }
    // API call with callbacks for success/failure
    com.ai.addxnet.vicohome_1742553098674_0O0o0oo.getInstance().unbindDevice(sn, new vicohome_1742553098674_00O0o0oOO(sn, callback));
}
```

## Error Handling
The response extends BaseResponse which contains:
- `result`: Result code (0 indicates success)
- `msg`: A message describing the result or error

The app handles these errors by:
1. Reporting unbind errors through the callback interface
2. Displaying appropriate UI messages to the user
3. Possibly attempting to retry the operation

Error conditions might include:
- Authentication failures
- Network connectivity issues
- Device is offline
- Server-side errors
- User does not have permission to unbind the device

## Notes
- This endpoint permanently removes the device association from a user's account
- After a successful call, the device will need to be set up again to be used
- The application appears to handle device state cleanup locally after a successful API call
- The implementation includes callback interfaces for tracking the unbinding process:
  - `onStartUnBind(String sn)` - Called when the unbind process starts
  - `onUnBindSuccess(String sn)` - Called on successful unbinding
  - `onUnbindError(String sn)` - Called when unbinding fails