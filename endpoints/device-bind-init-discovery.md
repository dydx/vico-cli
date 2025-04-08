# Device Bind Initialization

Retrieves initialization information for a device during the binding process, including device name, firmware status, and available locations.

## Endpoint

`/device/deviceBindInit`

## Method

POST

## Request

The request should contain:

| Field | Type | Description |
|-------|------|-------------|
| serialNumber | String | The serial number of the device to bind |
| bindType | Integer | The type of binding (0 or 1) |

## Response

The response includes:

| Field | Type | Description |
|-------|------|-------------|
| result | Integer | Result code (0 for success, negative values for errors) |
| data | Object | The data object containing device information |
| data.deviceName | String | The device name |
| data.firmwareStatus | Integer | Firmware status code (bit 0 indicates if OTA is needed, bit 3 indicates if force OTA is needed) |
| data.locationDOList | Array | List of available locations |
| data.locationId | Integer | Default location ID |
| data.newestFirmwareId | String | The ID of the newest available firmware |
| data.serialNumber | String | The serial number of the device |

## Example Request

```json
{
  "serialNumber": "ABCD1234",
  "bindType": 0
}
```

## Example Response

```json
{
  "result": 0,
  "message": "Success",
  "data": {
    "deviceName": "My Camera",
    "firmwareStatus": 0,
    "locationDOList": [
      {
        "id": 123,
        "locationName": "Living Room",
        "isSelected": false,
        "isLocalData": false
      }
    ],
    "locationId": 123,
    "newestFirmwareId": "firmware-v2.0",
    "serialNumber": "ABCD1234"
  }
}
```

## Error Handling

The API returns a result code in the response. A value of 0 indicates success, while negative values indicate different errors.

Common error codes:
- Error code < 0: Generic error