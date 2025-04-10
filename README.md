
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

### `./vico-cli events list` [PASS]

```bash
➜  vicohome git:(main) ✗ ./vico-cli events list --hours 2
Trace ID                             Timestamp            Device Name               Bird Name                 Bird Latin               
--------------------------------------------------------------------------------------------------
018594221744243886k4jua3TyFQq        2025-04-09 20:11:24  Birdies                   Eastern Phoebe            Sayornis phoebe 
```

### `./vico-cli events list --hours 2` [PASS]

```bash
➜  vicohome git:(main) ✗ ./vico-cli events list --hours 2
Trace ID                             Timestamp            Device Name               Bird Name                 Bird Latin               
--------------------------------------------------------------------------------------------------
018594221744243886k4jua3TyFQq        2025-04-09 20:11:24  Birdies                   Eastern Phoebe            Sayornis phoebe 
```

### `./vico-cli events list --hours 1 --format json` [FAIL]

```bash
➜  vicohome git:(main) ✗ ./vico-cli events list --hours 2 --format json
[
  {
    "traceId": "018594221744243886k4jua3TyFQq",
    "timestamp": "2025-04-09 20:11:24",
    "deviceName": "Birdies",
    "serialNumber": "378b660598295ceca8b20871991a0409",
    "adminName": "jpsandlin",
    "period": "19.66s",
    "keyshots": [
      {
        "imageUrl": "https://a4x-prod-us.s3.amazonaws.com/ai-saas-out-storage/keyshot_front_bird_018594221744243886k4jua3TyFQq_countryNo_US.jpg",
        "message": "1.000000",
        "objectCategory": "bird",
        "subCategoryName": "Sayornis phoebe"
      }
    ],
    "birdName": "Eastern Phoebe",
    "birdLatin": "Sayornis phoebe",
    "birdConfidence": 0.996811,
    "imageUrl": "https://a4x-prod-us.s3.amazonaws.com/ai-saas-out-storage/018594221744243886k4jua3TyFQq_gallery.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Date=20250410T015408Z\u0026X-Amz-SignedHeaders=host\u0026X-Amz-Expires=172800\u0026X-Amz-Credential=AKIAQBFG53LBAA5AEUVF%2F20250410%2Fus-east-1%2Fs3%2Faws4_request\u0026X-Amz-Signature=d316ba785c46ba12b810c373c3c606b238ee4c0c49b9e0e400cfb01d784b844d",
    "videoUrl": "https://api-us.vicohome.io/video/download/m3u8/018594221744243886k4jua3TyFQq.m3u8?token=eyJhbGciOiJIUzUxMiJ9.eyJ1c2VySWQiOjE4NTk0MjIsInRyYWNlSWQiOiIwMTg1OTQyMjE3NDQyNDM4ODZrNGp1YTNUeUZRcSIsImV4cCI6MTc0NDQyMjg0OH0.RfeVYvRGRYXEl8RY7CUVQyMowaSH7wwyXFK-ad_2aHyMCVf5am1dQ2BKW_Qr3yCP-TXOgNnU-oHzhpxEbGavnw"
  }
]
```

What we have now is almost perfect, but lets clean up the response event further so its totally flat. Pull `keyshots.imageUrl` up as `keyShotUrl` and omit the `keyshots` field in the response JSON.

```json
[
  {
    "traceId": "018594221744243886k4jua3TyFQq",
    "timestamp": "2025-04-09 20:11:24",
    "deviceName": "Birdies",
    "serialNumber": "378b660598295ceca8b20871991a0409",
    "adminName": "jpsandlin",
    "period": "19.66s",
    "birdName": "Eastern Phoebe",
    "birdLatin": "Sayornis phoebe",
    "birdConfidence": 0.996811,
    "keyShotUrl": "https://a4x-prod-us.s3.amazonaws.com/ai-saas-out-storage/keyshot_front_bird_018594221744243886k4jua3TyFQq_countryNo_US.jpg",
    "imageUrl": "https://a4x-prod-us.s3.amazonaws.com/ai-saas-out-storage/018594221744243886k4jua3TyFQq_gallery.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Date=20250410T015408Z\u0026X-Amz-SignedHeaders=host\u0026X-Amz-Expires=172800\u0026X-Amz-Credential=AKIAQBFG53LBAA5AEUVF%2F20250410%2Fus-east-1%2Fs3%2Faws4_request\u0026X-Amz-Signature=d316ba785c46ba12b810c373c3c606b238ee4c0c49b9e0e400cfb01d784b844d",
    "videoUrl": "https://api-us.vicohome.io/video/download/m3u8/018594221744243886k4jua3TyFQq.m3u8?token=eyJhbGciOiJIUzUxMiJ9.eyJ1c2VySWQiOjE4NTk0MjIsInRyYWNlSWQiOiIwMTg1OTQyMjE3NDQyNDM4ODZrNGp1YTNUeUZRcSIsImV4cCI6MTc0NDQyMjg0OH0.RfeVYvRGRYXEl8RY7CUVQyMowaSH7wwyXFK-ad_2aHyMCVf5am1dQ2BKW_Qr3yCP-TXOgNnU-oHzhpxEbGavnw"
  }
]
```

### `./vico-cli events get <> --format json` [FAIl]

```bash
➜  vicohome git:(main) ✗ ./vico-cli events get 018594221744243886k4jua3TyFQq --format json
{
  "traceId": "018594221744243886k4jua3TyFQq",
  "timestamp": "2025-04-09 20:11:24",
  "deviceName": "Birdies",
  "serialNumber": "378b660598295ceca8b20871991a0409",
  "adminName": "jpsandlin",
  "period": "19.66s",
  "keyshots": [
    {
      "imageUrl": "https://a4x-prod-us.s3.amazonaws.com/ai-saas-out-storage/keyshot_front_bird_018594221744243886k4jua3TyFQq_countryNo_US.jpg",
      "message": "1.000000",
      "objectCategory": "bird",
      "subCategoryName": "Sayornis phoebe"
    }
  ],
  "birdName": "Unidentified",
  "birdLatin": "",
  "birdConfidence": 0,
  "imageUrl": "https://a4x-prod-us.s3.amazonaws.com/ai-saas-out-storage/018594221744243886k4jua3TyFQq_gallery.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Date=20250410T015747Z\u0026X-Amz-SignedHeaders=host\u0026X-Amz-Expires=172800\u0026X-Amz-Credential=AKIAQBFG53LBAA5AEUVF%2F20250410%2Fus-east-1%2Fs3%2Faws4_request\u0026X-Amz-Signature=26198ef1e0f17a0505801f7582787fb2c854c7d0e556ec6fa49dfd878192a716",
  "videoUrl": "https://api-us.vicohome.io/video/download/m3u8/018594221744243886k4jua3TyFQq.m3u8?token=eyJhbGciOiJIUzUxMiJ9.eyJ1c2VySWQiOjE4NTk0MjIsInRyYWNlSWQiOiIwMTg1OTQyMjE3NDQyNDM4ODZrNGp1YTNUeUZRcSIsImV4cCI6MTc0NDQyMzA2N30.DkhgzGQdxPPmdIluYpcEfzTw6Uob92ijDO0sTZJfe_46AGiOuZfLDp4Dd-W6BH6JkcaVy4_HboF-RstyfftISw"
}
```

```json
{
  "traceId": "018594221744243886k4jua3TyFQq",
  "timestamp": "2025-04-09 20:11:24",
  "deviceName": "Birdies",
  "serialNumber": "378b660598295ceca8b20871991a0409",
  "adminName": "jpsandlin",
  "period": "19.66s",
  "birdName": "Unidentified",
  "birdLatin": "",
  "birdConfidence": 0,
  "keyShotUrl": "https://a4x-prod-us.s3.amazonaws.com/ai-saas-out-storage/keyshot_front_bird_018594221744243886k4jua3TyFQq_countryNo_US.jpg"
  "imageUrl": "https://a4x-prod-us.s3.amazonaws.com/ai-saas-out-storage/018594221744243886k4jua3TyFQq_gallery.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Date=20250410T015747Z\u0026X-Amz-SignedHeaders=host\u0026X-Amz-Expires=172800\u0026X-Amz-Credential=AKIAQBFG53LBAA5AEUVF%2F20250410%2Fus-east-1%2Fs3%2Faws4_request\u0026X-Amz-Signature=26198ef1e0f17a0505801f7582787fb2c854c7d0e556ec6fa49dfd878192a716",
  "videoUrl": "https://api-us.vicohome.io/video/download/m3u8/018594221744243886k4jua3TyFQq.m3u8?token=eyJhbGciOiJIUzUxMiJ9.eyJ1c2VySWQiOjE4NTk0MjIsInRyYWNlSWQiOiIwMTg1OTQyMjE3NDQyNDM4ODZrNGp1YTNUeUZRcSIsImV4cCI6MTc0NDQyMzA2N30.DkhgzGQdxPPmdIluYpcEfzTw6Uob92ijDO0sTZJfe_46AGiOuZfLDp4Dd-W6BH6JkcaVy4_HboF-RstyfftISw"
}
```

Same failure cause as the prior test. See the recommended resolution. 