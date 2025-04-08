# Update Read Status (/library/updatemissing)

## Endpoint Details

**URL**: `https://api-us.vicohome.io/library/updatemissing`
**Method**: POST
**Description**: Sets the read status for one or more event recordings in the user's library.

## Request Parameters

### SetReadStatueEntry

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| missing | int | Yes | The read status to set (1 = read, 0 = unread) |
| traceId | String | No* | Single trace ID of the recording to update |
| traceIds | String | No* | Multiple trace IDs of recordings to update, comma-separated |

*Either traceId or traceIds must be provided

### BaseEntry Fields (Inherited)

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| app | AppBean | Yes | Application information |
| countryNo | String | No | Country code |
| language | String | No | UI language |
| tenantId | String | No | Tenant ID |

### AppBean Fields

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| apiVersion | String | No | API version |
| appName | String | No | Application name |
| appType | String | No | Application type, defaults to "Android" |
| bundle | String | No | Bundle identifier |
| countlyId | String | No | Analytics ID |
| env | String | No | Environment |
| tenantId | String | No | Tenant ID |
| timeZone | String | No | User's timezone |
| version | int | No | App version code |
| versionName | String | No | App version name |

## Response

### BaseResponse

| Parameter | Type | Description |
|-----------|------|-------------|
| result | int | Result code (0 = success, non-zero = error) |
| msg | String | Result message or error description |

## Usage Example

```java
// Example from LibraryCore.java
public final void setReadStatus(int i, String str, com.ai.addxbase.vicohome_1742553098674_0O0oO0O<String> callBack) {
    vicohome_1742553098674_0O0oOoo.checkNotNullParameter(callBack, "callBack");
    LibraryApiClient sInstance = LibraryApiClient.vicohome_1742553098674_o00OOoo.getSInstance();
    SetReadStatueEntry setReadStatueEntry = new SetReadStatueEntry();
    setReadStatueEntry.setMissing(i);
    setReadStatueEntry.setTraceIds(str);
    this.vicohome_1742553098674_0o00OOoOo.add(sInstance.setReadStatus(setReadStatueEntry).subscribeOn(Schedulers.io()).observeOn(AndroidSchedulers.mainThread()).subscribe((Subscriber<? super BaseResponse>) new vicohome_1742553098674_0O0oOO0(callBack)));
}
```

## Error Handling

When an error occurs, the response includes a non-zero result code and an error message in the `msg` field. The app handles errors by passing them to the callback's `onError` method:

```java
@Override // com.ai.addxnet.vicohome_1742553098674_0O0oO000
public void doOnError(int i, String msg) {
    vicohome_1742553098674_0O0oOoo.checkNotNullParameter(msg, "msg");
    this.vicohome_1742553098674_o00Oo0o0.onError(i, msg);
}
```

On success, the result code will be 0 and the callback's `onSuccess` method is called:

```java
@Override // vicohome_1742553098674_oO0O00o.vicohome_1742553098674_0O0o0oo
public void doOnNext(BaseResponse response) {
    vicohome_1742553098674_0O0oOoo.checkNotNullParameter(response, "response");
    int result = response.getResult();
    String msg = response.getMsg();
    vicohome_1742553098674_0O0oOoo.checkNotNullExpressionValue(msg, "response.msg");
    this.vicohome_1742553098674_o00Oo0o0.onSuccess(result, msg, response.getMsg());
}
```