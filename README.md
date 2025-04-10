
# Testing

## Devices

### `./vico-cli devices list` [PASS]

```bash
➜  vicohome git:(main) ✗ ./vico-cli devices list              
Serial Number                        Model                Name                 Network         IP              Battery
----------------------------------------------------------------------------------------------------------------
854396ddc826ed6e3e4263fa067ee288     CG625-BD-TNBD-SS2    Birdy House          Rocinante       192.168.10.223  100%
378b660598295ceca8b20871991a0409     CG623G-ST1BQJ        Birdies              Rocinante       192.168.10.107  100%
```

### `./vico-cli devices list --format table` [PASS]

```bash
➜  vicohome git:(main) ✗ ./vico-cli devices list --format table
Serial Number                        Model                Name                 Network         IP              Battery
----------------------------------------------------------------------------------------------------------------
854396ddc826ed6e3e4263fa067ee288     CG625-BD-TNBD-SS2    Birdy House          Rocinante       192.168.10.223  100%
378b660598295ceca8b20871991a0409     CG623G-ST1BQJ        Birdies              Rocinante       192.168.10.107  100%
```

### `./vico-cli devices list --format json` [PASS]

```bash
➜  vicohome git:(main) ✗ ./vico-cli devices list --format json
[
  {
    "serialNumber": "854396ddc826ed6e3e4263fa067ee288",
    "modelNo": "CG625-BD-TNBD-SS2",
    "deviceName": "Birdy House",
    "networkName": "Rocinante",
    "ip": "192.168.10.223",
    "batteryLevel": 100,
    "locationName": "Garden",
    "signalStrength": -54,
    "wifiChannel": 6,
    "isCharging": 0,
    "chargingMode": 0,
    "macAddress": "b4:61:e9:72:a0:15"
  },
  {
    "serialNumber": "378b660598295ceca8b20871991a0409",
    "modelNo": "CG623G-ST1BQJ",
    "deviceName": "Birdies",
    "networkName": "Rocinante",
    "ip": "192.168.10.107",
    "batteryLevel": 100,
    "locationName": "Garden",
    "signalStrength": -51,
    "wifiChannel": 6,
    "isCharging": 0,
    "chargingMode": 0,
    "macAddress": "b4:61:e9:35:7d:af"
  }
]
```

## Events

### `./vico-cli events list` [FAIL]

```bash
➜  vicohome git:(main) ✗ ./vico-cli events list
Trace ID                             Timestamp            Device Name                    Serial Number                       
------------------------------------------------------------------------------------------------
018594221744243886k4jua3TyFQq                             Birdies                        378b660598295ceca8b20871991a0409    
018594221744242378uCKbHnbEp74                             Birdies                        378b660598295ceca8b20871991a0409    
018594221744242295eOqKEyzQb2f                             Birdies                        378b660598295ceca8b20871991a0409    
018594221744241531w5075vdOeR3                             Birdies                        378b660598295ceca8b20871991a0409    
018594221744241450fq1YMeKsWvi                             Birdies                        378b660598295ceca8b20871991a0409    
018594221744241368HcUvRcjMpz9                             Birdies                        378b660598295ceca8b20871991a0409
```

It is apparent that our transformation of the value for `timestamp` is not correct. All events have a `timestamp`, it is an instrinsic part of being an Event. In observing response information, this appears to just be a json field like `"timestamp": 1744156238,`. Consider identifying the cause of our issue in transforming this correctly?

### `./vico-cli events list --hours 1` [FAIL]

```bash
➜  vicohome git:(main) ✗ ./vico-cli events list --hours 2
Trace ID                             Timestamp            Device Name                    Serial Number                       
------------------------------------------------------------------------------------------------
018594221744243886k4jua3TyFQq                             Birdies                        378b660598295ceca8b20871991a0409    
018594221744242378uCKbHnbEp74                             Birdies                        378b660598295ceca8b20871991a0409    
018594221744242295eOqKEyzQb2f                             Birdies                        378b660598295ceca8b20871991a0409
```

Same failure cause as the prior test. See the recommended resolution.

### `./vico-cli events list --hours 1 --format json` [FAIL]

```bash
➜  vicohome git:(main) ✗ ./vico-cli events list --hours 2 --format json
[
  {
    "traceId": "018594221744243886k4jua3TyFQq",
    "timestamp": "",
    "deviceName": "Birdies",
    "serialNumber": "378b660598295ceca8b20871991a0409",
    "adminName": "jpsandlin",
    "period": "",
    "keyshots": [
      {
        "imageUrl": "https://a4x-prod-us.s3.amazonaws.com/ai-saas-out-storage/keyshot_front_bird_018594221744243886k4jua3TyFQq_countryNo_US.jpg",
        "message": "1.000000",
        "objectCategory": "bird",
        "subCategoryName": "Sayornis phoebe"
      }
    ],
    "imageUrl": "https://a4x-prod-us.s3.amazonaws.com/ai-saas-out-storage/018594221744243886k4jua3TyFQq_gallery.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Date=20250410T013547Z\u0026X-Amz-SignedHeaders=host\u0026X-Amz-Expires=172800\u0026X-Amz-Credential=AKIAQBFG53LBAA5AEUVF%2F20250410%2Fus-east-1%2Fs3%2Faws4_request\u0026X-Amz-Signature=087828a6b38e69ad2fce823f2f54d2aa27481214726a02f6e28c3f53dd9f1313",
    "videoUrl": "https://api-us.vicohome.io/video/download/m3u8/018594221744243886k4jua3TyFQq.m3u8?token=eyJhbGciOiJIUzUxMiJ9.eyJ1c2VySWQiOjE4NTk0MjIsInRyYWNlSWQiOiIwMTg1OTQyMjE3NDQyNDM4ODZrNGp1YTNUeUZRcSIsImV4cCI6MTc0NDQyMTc0N30.PXiy3XHUz3FkmC2kd0utQhnQLOU0APYrDaD0DD0YmX12i0euoUB7Cfo8qGVj40FzPhmwIGSj5MPVb2lKVDHSTw"
  },
  {
    "traceId": "018594221744242378uCKbHnbEp74",
    "timestamp": "",
    "deviceName": "Birdies",
    "serialNumber": "378b660598295ceca8b20871991a0409",
    "adminName": "jpsandlin",
    "period": "",
    "keyshots": [],
    "imageUrl": "https://a4x-prod-us-vip-pro.s3.amazonaws.com/device_video_slice/378b660598295ceca8b20871991a0409/018594221744242378uCKbHnbEp74/image.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Date=20250410T013547Z\u0026X-Amz-SignedHeaders=host\u0026X-Amz-Expires=172800\u0026X-Amz-Credential=AKIAQBFG53LBAA5AEUVF%2F20250410%2Fus-east-1%2Fs3%2Faws4_request\u0026X-Amz-Signature=0e32e5425763b984ab62ab521ddfc16b45fac96c6eb7ab9229484ebc8f5dcd09",
    "videoUrl": "https://api-us.vicohome.io/video/download/m3u8/018594221744242378uCKbHnbEp74.m3u8?token=eyJhbGciOiJIUzUxMiJ9.eyJ1c2VySWQiOjE4NTk0MjIsInRyYWNlSWQiOiIwMTg1OTQyMjE3NDQyNDIzNzh1Q0tiSG5iRXA3NCIsImV4cCI6MTc0NDQyMTc0N30.1phAHRD2e60g6v29s94YbzV_3i7d7f9TXQWAvnsjbDK2gluRzLqcJYP7p3ZAJGaFWD717b18SULXaLByv9OCow"
  },
  {
    "traceId": "018594221744242295eOqKEyzQb2f",
    "timestamp": "",
    "deviceName": "Birdies",
    "serialNumber": "378b660598295ceca8b20871991a0409",
    "adminName": "jpsandlin",
    "period": "",
    "keyshots": [],
    "imageUrl": "https://a4x-prod-us-vip-pro.s3.amazonaws.com/device_video_slice/378b660598295ceca8b20871991a0409/018594221744242295eOqKEyzQb2f/image.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Date=20250410T013547Z\u0026X-Amz-SignedHeaders=host\u0026X-Amz-Expires=172800\u0026X-Amz-Credential=AKIAQBFG53LBAA5AEUVF%2F20250410%2Fus-east-1%2Fs3%2Faws4_request\u0026X-Amz-Signature=cb9afcaa6fd0414207197252a4df258f6ad554437d41efcd0295eb2bcaf01da0",
    "videoUrl": "https://api-us.vicohome.io/video/download/m3u8/018594221744242295eOqKEyzQb2f.m3u8?token=eyJhbGciOiJIUzUxMiJ9.eyJ1c2VySWQiOjE4NTk0MjIsInRyYWNlSWQiOiIwMTg1OTQyMjE3NDQyNDIyOTVlT3FLRXl6UWIyZiIsImV4cCI6MTc0NDQyMTc0N30.neDr3EqP8DOErIKMN4RTcAbZFWiA00UIWWA-iF_n7tG9dg7RN6o9CDYQXgCtofsm8sNqx5xFgi8YCXoqC3PYRQ"
  }
]
```

Same failure cause as the prior test. See the recommended resolution. Additional failure causes:

* `period` is incorrect. All events have a period. It is an intrinsic part of them. A value, `"period": "",`, is sort of impossible. This indicates incorrect mapping. Correct value is found in event as like `"period": 19.659,`
* usage of `"keyshots"` instead of `"subcategoryInfoList"` when showing bird data. This results in not obtaining nested name value like `"objectName": "House Finch",` and not obtaining correct `url` and `confidence`.

Ideally, what we could actually do is have a placeholder in the cli response for `birdName`. If we have an event with no detected bird and no `"subcategoryInfoList"`, then `birdName` would default to "unidentified". If we do have a bird detected and `"subcategoryInfoList"` is present, we should look for `objectName` and map this to our cli response for `birdName`. This will be a value like `House Finch`. We can map `birdStdName` to `birdLatin` in our CLI response. This should lead to a more concise and consistent CLI output. For sake of reference, here is that `"subcategoryInfoList"` object from a real response:

```json
  "subcategoryInfoList": [
      {
          "objectType": "bird",
          "objectName": "House Finch",
          "birdStdName": "Haemorhous mexicanus",
          "url": "https://a4x-prod-us.s3.amazonaws.com/ai-saas-out-storage/0185942217441562400ny2nw9yIw0_countryNo_US_bird_gallery.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20250409T235501Z&X-Amz-SignedHeaders=host&X-Amz-Expires=172800&X-Amz-Credential=AKIAQBFG53LBAA5AEUVF%2F20250409%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Signature=45f98c501df00715a304f85b70b9b89a96c13103700ed8dd4f42da1b507310b1",
          "confidence": 0.979434
      }
  ],
```

### `./vico-cli events get <> --format json` [FAIl]

```bash
➜  vicohome git:(main) ✗ ./vico-cli events get 018594221744242295eOqKEyzQb2f --format json
{
  "traceId": "018594221744242295eOqKEyzQb2f",
  "timestamp": "",
  "deviceName": "Birdies",
  "serialNumber": "378b660598295ceca8b20871991a0409",
  "adminName": "jpsandlin",
  "period": "",
  "keyshots": [],
  "imageUrl": "https://a4x-prod-us-vip-pro.s3.amazonaws.com/device_video_slice/378b660598295ceca8b20871991a0409/018594221744242295eOqKEyzQb2f/image.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Date=20250410T014410Z\u0026X-Amz-SignedHeaders=host\u0026X-Amz-Expires=172800\u0026X-Amz-Credential=AKIAQBFG53LBAA5AEUVF%2F20250410%2Fus-east-1%2Fs3%2Faws4_request\u0026X-Amz-Signature=040fdd1c01351ee50a2e650b6b0a624a6af7a1557b8cbb9461a9a1bb01857101",
  "videoUrl": "https://api-us.vicohome.io/video/download/m3u8/018594221744242295eOqKEyzQb2f.m3u8?token=eyJhbGciOiJIUzUxMiJ9.eyJ1c2VySWQiOjE4NTk0MjIsInRyYWNlSWQiOiIwMTg1OTQyMjE3NDQyNDIyOTVlT3FLRXl6UWIyZiIsImV4cCI6MTc0NDQyMjI1MH0.tjlOVAyEuoKJqP47sVP5PjREcivD8vJS_ChxtVcBQaurB-kDX3Qxj-Ds5lBZzjVH3PytHOLb3TbuDNHkBVUUAQ"
}
```

Same failure cause as the prior test. See the recommended resolution. 