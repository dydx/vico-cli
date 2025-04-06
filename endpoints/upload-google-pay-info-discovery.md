# Upload Google Pay Info API

## Endpoint Details
- **Path**: `/pay/google/order/verify`
- **Method**: POST
- **Content-Type**: application/json
- **Description**: Processes Google Pay payments and activates subscriptions for selected devices

## Request Parameters

```json
{
  "productId": 123,
  "tradeType": 1,
  "outTradeNo": "ORD12345678",
  "purchaseToken": "google-pay-token-abcdef123456",
  "subscriptionGroupId": "sub_group_premium",
  "tierDeviceList": ["ABC123XYZ", "DEF456UVW"],
  "guidanceSource": 0,
  "app": {
    "appName": "vicohome",
    "appVersion": "1.2.3",
    "appBuild": "123",
    "channelId": 1
  },
  "countryNo": "US",
  "language": "en",
  "tenantId": "default"
}
```

### Required Fields
- `productId`: Integer - The product identifier
- `tradeType`: Integer - Type of transaction (subscription or one-time)
- `outTradeNo`: String - Order number for tracking the transaction
- `purchaseToken`: String - Token from Google Pay after payment
- `tierDeviceList`: Array of Strings - Device identifiers to activate subscription for
- Standard app and user context fields (app, countryNo, language, tenantId)

### Optional Fields
- `subscriptionGroupId`: String - Group identifier for subscription
- `guidanceSource`: Integer - Source of purchase referral

## Response Format

```json
{
  "result": 0,
  "msg": "success",
  "data": "transaction_id_12345"
}
```

### Success Response
- `result`: 0
- `msg`: "success"
- `data`: String - Transaction identifier or confirmation ID

### Error Responses
- `result`: Negative value
- `msg`: Error description message
- Common error codes:
  - `-1001`: Invalid product ID
  - `-1002`: Invalid purchase token
  - `-1003`: Transaction verification failed
  - `-2001`: Network error
  - `-4001`: Authentication error

## Usage Notes
- This endpoint is used to verify and process payments made through Google Pay
- The request contains the purchase token obtained from Google Pay after payment
- The transaction type (tradeType) determines whether this is a subscription or one-time purchase
- Device IDs are provided to link the subscription to specific devices
- After successful processing, the response includes a transaction ID for reference