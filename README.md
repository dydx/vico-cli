
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

### `./vico-cli events list --hours 1 --format json` [PASS]

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
    "birdName": "Eastern Phoebe",
    "birdLatin": "Sayornis phoebe",
    "birdConfidence": 0.996811,
    "keyShotUrl": "https://a4x-prod-us.s3.amazonaws.com/ai-saas-out-storage/keyshot_front_bird_018594221744243886k4jua3TyFQq_countryNo_US.jpg",
    "imageUrl": "https://a4x-prod-us.s3.amazonaws.com/ai-saas-out-storage/018594221744243886k4jua3TyFQq_gallery.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Date=20250410T020507Z\u0026X-Amz-SignedHeaders=host\u0026X-Amz-Expires=172799\u0026X-Amz-Credential=AKIAQBFG53LBAA5AEUVF%2F20250410%2Fus-east-1%2Fs3%2Faws4_request\u0026X-Amz-Signature=ff9fb6d04654ac8bdd460368bb5bb57745716bc4634310a068a70c33ba099a6d",
    "videoUrl": "https://api-us.vicohome.io/video/download/m3u8/018594221744243886k4jua3TyFQq.m3u8?token=eyJhbGciOiJIUzUxMiJ9.eyJ1c2VySWQiOjE4NTk0MjIsInRyYWNlSWQiOiIwMTg1OTQyMjE3NDQyNDM4ODZrNGp1YTNUeUZRcSIsImV4cCI6MTc0NDQyMzUwN30.sZkpWZnQ31UB2s6h7kTSredwFYdT8eCxBmR3F-xHK4CuqAwhW5aXBebykZ_sKfZtDVzYPeGSZLAgqrUY6OhhDQ"
  }
]
```

### `./vico-cli events get <> --format json` [PASS]

```bash
➜  vicohome git:(main) ✗ ./vico-cli events get 018594221744243886k4jua3TyFQq --format json
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
  "keyShotUrl": "https://a4x-prod-us.s3.amazonaws.com/ai-saas-out-storage/keyshot_front_bird_018594221744243886k4jua3TyFQq_countryNo_US.jpg",
  "imageUrl": "https://a4x-prod-us.s3.amazonaws.com/ai-saas-out-storage/018594221744243886k4jua3TyFQq_gallery.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Date=20250410T020548Z\u0026X-Amz-SignedHeaders=host\u0026X-Amz-Expires=172800\u0026X-Amz-Credential=AKIAQBFG53LBAA5AEUVF%2F20250410%2Fus-east-1%2Fs3%2Faws4_request\u0026X-Amz-Signature=48ff2901444648a2e2454f9b45fe97c5e0370baf92b2e179509ab33ca3a3d170",
  "videoUrl": "https://api-us.vicohome.io/video/download/m3u8/018594221744243886k4jua3TyFQq.m3u8?token=eyJhbGciOiJIUzUxMiJ9.eyJ1c2VySWQiOjE4NTk0MjIsInRyYWNlSWQiOiIwMTg1OTQyMjE3NDQyNDM4ODZrNGp1YTNUeUZRcSIsImV4cCI6MTc0NDQyMzU0OH0.0f43zpnBQ6Xydq5asTsF5D5oHB69HAnZj7_42aeJW_wxE04TW4XT_uvQckr8Q5jA3-d0anwaYbtqjlfDfVST-Q"
}
```