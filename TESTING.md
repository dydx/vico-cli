
# Testing

Testing so far is really just done by my describing the interface I want in as much detail as I can, and then having the AI iterate to meet those requirements. I manually set each new test to `[FAIL]` and post the CLI usage stating "feature missing" or whatever. Claude Code then reads this test file, seeks for `[FAIL]` and addresses it.

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
➜  vicohome git:(main) ✗ ./vico-cli events list --startTime "2025-05-17 18:00:00" --endTime "2025-05-17 20:00:00"
Trace ID                             Timestamp            Device Name               Bird Name                 Bird Latin               
--------------------------------------------------------------------------------------------------
018594221744243886k4jua3TyFQq        2025-04-09 20:11:24  Birdies                   Eastern Phoebe            Sayornis phoebe 
```

### `./vico-cli events list --startTime <> --endTime <>` [PASS]

```bash
➜  vicohome git:(main) ✗ ./vico-cli events list --startTime "2025-05-17 18:00:00" --endTime "2025-05-17 20:00:00"
Trace ID                             Timestamp            Device Name               Bird Name                 Bird Latin               
--------------------------------------------------------------------------------------------------
018594221744243886k4jua3TyFQq        2025-04-09 20:11:24  Birdies                   Eastern Phoebe            Sayornis phoebe 
```

### `./vico-cli events list --startTime <> --endTime <> --format json` [PASS]

```bash
➜  vicohome git:(main) ✗ ./vico-cli events list --startTime "2025-05-17 18:00:00" --endTime "2025-05-17 20:00:00" --format json
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

### `./vico-cli events get <>` [PASS]

```bash
➜  vicohome git:(main) ✗ ./vico-cli events get 018594221744243886k4jua3TyFQq
Event Details:
------------------------------
Trace ID:       018594221744243886k4jua3TyFQq
Timestamp:      2025-04-09 20:11:24
Device Name:    Birdies
Serial Number:  378b660598295ceca8b20871991a0409
Admin Name:     jpsandlin
Period:         19.66s
Bird Name:      Unidentified
Bird Latin:     
KeyShot URL:    https://a4x-prod-us.s3.amazonaws.com/ai-saas-out-storage/keyshot_front_bird_018594221744243886k4jua3TyFQq_countryNo_US.jpg
Image URL:      https://a4x-prod-us.s3.amazonaws.com/ai-saas-out-storage/018594221744243886k4jua3TyFQq_gallery.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20250410T021654Z&X-Amz-SignedHeaders=host&X-Amz-Expires=172800&X-Amz-Credential=AKIAQBFG53LBAA5AEUVF%2F20250410%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Signature=bda0dc1a76f22ba659050ba5ff7882206e7e0ad9dd84b4154fd6c72d0984af7a
Video URL:      https://api-us.vicohome.io/video/download/m3u8/018594221744243886k4jua3TyFQq.m3u8?token=eyJhbGciOiJIUzUxMiJ9.eyJ1c2VySWQiOjE4NTk0MjIsInRyYWNlSWQiOiIwMTg1OTQyMjE3NDQyNDM4ODZrNGp1YTNUeUZRcSIsImV4cCI6MTc0NDQyNDIxNH0.H44qdq7n7s6jDfeM_3k7f00DwEpFSoiMXdFyXtnPZ1dr6zDBgUXzbuS758RYBdQ1XixzE9MjvSmFvQ-jyZpUgA
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

## Events Search

### `./vico-cli events search --field serialNumber <> --startTime <> --endTime <>` [PASS]

```bash
➜  vicohome git:(main) ✗ ./vico-cli events search --field serialNumber 378b660598295ceca8b20871991a0409 --startTime "2025-05-17 13:00:00" --endTime "2025-05-17 15:00:00"
Trace ID                             Timestamp            Device Name               Bird Name                 Bird Latin               
--------------------------------------------------------------------------------------------------
018594221744308360Sr56DmocjwP        2025-04-10 14:05:58  Birdies                   Eastern Bluebird          Sialia sialis            
018594221744307663G74FrEyoM9B        2025-04-10 13:54:21  Birdies                   Eastern Bluebird          Sialia sialis            
018594221744305940Q2JsDKoTGtm        2025-04-10 13:25:38  Birdies                   Unidentified                                       
018594221744305706UsJKzOrbVCH        2025-04-10 13:21:44  Birdies                   Eastern Bluebird          Sialia sialis            
018594221744305220fuWydkIg49g        2025-04-10 13:13:38  Birdies                   Unidentified                                       
018594221744305139l0SijUhR2Mc        2025-04-10 13:12:17  Birdies                   Unidentified                                       
018594221744304528n3UAE77kGPg        2025-04-10 13:02:06  Birdies                   Unidentified                                       
0185942217443044472MkMvzL2O0j        2025-04-10 13:00:45  Birdies                   Unidentified                                       
018594221744303696AAZDjT4vVQM        2025-04-10 12:48:14  Birdies                   Eastern Bluebird          Sialia sialis            
018594221744303514zXrR13t7xpl        2025-04-10 12:45:12  Birdies                   House Finch               Haemorhous mexicanus
```

### `./vico-cli events search --field birdName <> --startTime <> --endTime <>` [PASS]

```bash
➜  vicohome git:(main) ✗ ./vico-cli events search --field birdName "Eastern Phoebe" --startTime "2025-05-17 00:00:00" --endTime "2025-05-18 00:00:00"
Trace ID                             Timestamp            Device Name               Bird Name                 Bird Latin               
--------------------------------------------------------------------------------------------------
018594221744243886k4jua3TyFQq        2025-04-09 20:11:24  Birdies                   Eastern Phoebe            Sayornis phoebe
```

### `./vico-cli events search --field birdName <> --startTime <> --endTime <> --format json` [PASS]

```bash
➜  vicohome git:(main) ✗ ./vico-cli events search --field birdName "Eastern Phoebe" --startTime "2025-05-17 00:00:00" --endTime "2025-05-18 00:00:00" --format json
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
    "imageUrl": "https://a4x-prod-us.s3.amazonaws.com/ai-saas-out-storage/018594221744243886k4jua3TyFQq_gallery.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Date=20250410T183250Z\u0026X-Amz-SignedHeaders=host\u0026X-Amz-Expires=172800\u0026X-Amz-Credential=AKIAQBFG53LBAA5AEUVF%2F20250410%2Fus-east-1%2Fs3%2Faws4_request\u0026X-Amz-Signature=41455127a261dd15c9db24149edda5f7ead0a27725783dced738f84b913640ce",
    "videoUrl": "https://api-us.vicohome.io/video/download/m3u8/018594221744243886k4jua3TyFQq.m3u8?token=eyJhbGciOiJIUzUxMiJ9.eyJ1c2VySWQiOjE4NTk0MjIsInRyYWNlSWQiOiIwMTg1OTQyMjE3NDQyNDM4ODZrNGp1YTNUeUZRcSIsImV4cCI6MTc0NDQ4Mjc3MH0.oQEWEhLPTgARph1MdOuaT3uyj89RbK5kLQ4NN0viyN5jRaeUta0SonIu8-4swXE4AeMjpAISblmvWwf9B1kqnw"
  }
]
```

### `./vico-cli events search --field deviceName <> --startTime <> --endTime <>` [PASS]

```bash
➜  vicohome git:(main) ✗ ./vico-cli events search --field deviceName "Birdies" --startTime "2025-05-17 12:00:00" --endTime "2025-05-17 15:00:00"
Trace ID                             Timestamp            Device Name               Bird Name                 Bird Latin               
--------------------------------------------------------------------------------------------------
018594221744308360Sr56DmocjwP        2025-04-10 14:05:58  Birdies                   Eastern Bluebird          Sialia sialis            
018594221744307663G74FrEyoM9B        2025-04-10 13:54:21  Birdies                   Eastern Bluebird          Sialia sialis            
018594221744305940Q2JsDKoTGtm        2025-04-10 13:25:38  Birdies                   Unidentified                                       
018594221744305706UsJKzOrbVCH        2025-04-10 13:21:44  Birdies                   Eastern Bluebird          Sialia sialis            
018594221744305220fuWydkIg49g        2025-04-10 13:13:38  Birdies                   Unidentified                                       
018594221744305139l0SijUhR2Mc        2025-04-10 13:12:17  Birdies                   Unidentified                                       
018594221744304528n3UAE77kGPg        2025-04-10 13:02:06  Birdies                   Unidentified                                       
0185942217443044472MkMvzL2O0j        2025-04-10 13:00:45  Birdies                   Unidentified                                       
018594221744303696AAZDjT4vVQM        2025-04-10 12:48:14  Birdies                   Eastern Bluebird          Sialia sialis            
018594221744303514zXrR13t7xpl        2025-04-10 12:45:12  Birdies                   House Finch               Haemorhous mexicanus .
```