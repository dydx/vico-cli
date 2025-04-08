# Get Event Details API

## Endpoint

**Path**: `/library/newselectlibrary`
**Method**: POST

## Description

This endpoint retrieves detailed information about recorded events based on filtering criteria. It returns a list of recordings that match the specified filters.

## Request Parameters

The request body should contain a JSON object with the following fields:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| serialNumber | Array of Strings | No | List of device serial numbers to filter by |
| startTimestamp | Long | No | Start time in milliseconds since epoch |
| endTimestamp | Long | No | End time in milliseconds since epoch |
| from | Integer | No | Pagination start index |
| to | Integer | No | Pagination end index |
| tags | Array of Strings | No | Filter by specific event tags |
| objectNames | Array of Strings | No | Filter by detected object types |
| marked | Integer | No | Filter by marked status (0=unmarked, 1=marked) |
| missing | Integer | No | Filter by read status (0=read, 1=unread) |
| isFromSDCard | Boolean | No | Whether to fetch from SD card storage |
| serialNumberToActivityZone | Map<String, List<Integer>> | No | Filter by activity zones for specific devices |
| videoEventKey | String | No | Filter by specific video event key |
| doorbellTags | Array of Strings | No | Filter by doorbell specific tags |
| deviceCallEventTag | String | No | Filter by device call event tag |
| deviceName | String | No | Filter by device name |

## Response Parameters

| Field | Type | Description |
|-------|------|-------------|
| result | Integer | Result code (0 indicates success) |
| msg | String | Response message |
| data | Object | Response data containing the event details |
| data.total | Integer | Total number of records matching the filter criteria |
| data.list | Array | List of event records |

### Record Object Fields

| Field | Type | Description |
|-------|------|-------------|
| traceId | String | Unique identifier for the recording |
| serialNumber | String | Device serial number |
| deviceName | String | Name of the device |
| timestamp | Integer | Event timestamp in seconds since epoch |
| date | String | Formatted date string |
| videoUrl | String | URL to the video recording |
| imageUrl | String | URL to the event thumbnail image |
| fileSize | Integer | Size of the recording file in bytes |
| type | Integer | Event type code |
| period | Double | Recording duration in seconds |
| tags | String | Event tags |
| marked | Integer | Whether the event is marked/favorited (0=no, 1=yes) |
| missing | Integer | Whether the event is unread (0=read, 1=unread) |
| userId | Integer | User ID |
| userName | String | User name |
| adminId | Integer | Admin ID |
| adminName | String | Admin name |
| adminIsVip | Boolean | Whether the admin has VIP subscription |
| locationId | Integer | ID of the location |
| locationName | String | Name of the location |
| timeZone | Integer | Timezone offset in minutes |
| dst | Integer | Daylight saving time adjustment |
| timeFormat | Integer | Time format preference (12/24 hour) |
| activityZoneName | String | Name of the triggered activity zone |
| role | String | User role |
| mediaType | String | Type of media ("video" or "image") |
| imageOnly | Integer | Indicates if only image is available |
| videoEvent | String | Video event information |
| deviceAiEventList | Array of Strings | List of AI detection events |
| eventInfoList | Array of Strings | Additional event information |
| hasPossibleSubcategory | Boolean | Whether subcategories exist for this event |
| subcategoryInfoList | Array of Objects | List of object subcategory information |

### Subcategory Object Fields

| Field | Type | Description |
|-------|------|-------------|
| objectName | String | Name of the detected object |
| url | String | Image URL for the detected object |
| birdStdName | String | Standard name (for bird detection) |

## Examples

### Request

```json
{
  "serialNumber": ["ABC123456789"],
  "startTimestamp": 1649764800000,
  "endTimestamp": 1649851200000,
  "from": 0,
  "to": 20,
  "tags": ["motion"],
  "marked": 0
}
```

### Successful Response

```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "total": 42,
    "list": [
      {
        "traceId": "20220412150322_ABC123456789",
        "serialNumber": "ABC123456789",
        "deviceName": "Front Door Camera",
        "timestamp": 1649764800,
        "date": "2022-04-12 15:03:22",
        "videoUrl": "https://storage.vicohome.io/videos/20220412150322_ABC123456789.mp4",
        "imageUrl": "https://storage.vicohome.io/images/20220412150322_ABC123456789.jpg",
        "fileSize": 1245678,
        "type": 1,
        "period": 15.5,
        "tags": "motion",
        "marked": 0,
        "missing": 0,
        "userId": 12345,
        "userName": "user@example.com",
        "adminId": 12345,
        "adminName": "user@example.com",
        "adminIsVip": true,
        "locationId": 6789,
        "locationName": "Home",
        "timeZone": -300,
        "dst": 1,
        "timeFormat": 12,
        "activityZoneName": "Driveway",
        "role": "owner",
        "mediaType": "video",
        "imageOnly": 0,
        "deviceAiEventList": ["person", "vehicle"],
        "eventInfoList": ["Person detected"],
        "hasPossibleSubcategory": true,
        "subcategoryInfoList": [
          {
            "objectName": "person",
            "url": "https://storage.vicohome.io/ai/person/20220412150322_ABC123456789.jpg"
          }
        ]
      }
    ]
  }
}
```

### Error Response

```json
{
  "result": 10001,
  "msg": "invalid parameter"
}
```

## Error Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 10001 | Invalid parameter |
| 10010 | User not found |
| 11001 | Device not found |
| 20001 | Server error |

## Notes

- This endpoint supports pagination through the `from` and `to` parameters
- The `marked` parameter can be used to filter favorited events
- The `missing` parameter can be used to filter unread events
- For time-based filtering, use `startTimestamp` and `endTimestamp` in milliseconds
- When filtering by AI events, use the `objectNames` field