# queryFreeTier Endpoint

## Overview
The queryFreeTier endpoint retrieves information about free tier services available on the Vicohome platform. Unlike device-specific subscription endpoints, this endpoint provides global information about what features and services are available without requiring a paid subscription.

## API Details
- **Path**: `/vip/freetier/info`
- **Method**: POST
- **Authentication**: Required
- **Request Content-Type**: application/json
- **Description**: Retrieves information about free tier services and features.

## Request Parameters
The request body only requires standard BaseEntry properties:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| app | Object | Yes | Application information |
| countryNo | String | Yes | Country code (e.g., "US") |
| language | String | Yes | Language code (e.g., "en") |
| tenantId | String | Yes | Tenant identifier |

## Request Example
```json
{
  "app": "vicohome",
  "countryNo": "US",
  "language": "en",
  "tenantId": "AAA12345678"
}
```

## Response Structure
The response contains information about available free tier services:

| Property | Type | Description |
|----------|------|-------------|
| result | Integer | Status code (0 for success, negative values for errors) |
| msg | String | Status message or error description |
| data | Object | Contains free tier information |

### Data Object Structure
| Property | Type | Description |
|----------|------|-------------|
| freeTierEnabled | Boolean | Whether free tier services are enabled |
| cloudStorageDays | Integer | Number of days of cloud storage available in free tier |
| freeTierFeatures | Array | List of features available in the free tier |
| trialEnabled | Boolean | Whether a trial of premium features is available |
| trialDurationDays | Integer | Duration of trial period in days |
| trialFeatures | Array | List of features available during the trial period |

### Feature Object Structure
| Property | Type | Description |
|----------|------|-------------|
| featureKey | String | Unique identifier for the feature |
| featureName | String | Display name of the feature |
| featureDescription | String | Description of what the feature provides |
| enabled | Boolean | Whether the feature is enabled |

## Response Example
```json
{
  "result": 0,
  "msg": "success",
  "data": {
    "freeTierEnabled": true,
    "cloudStorageDays": 1,
    "freeTierFeatures": [
      {
        "featureKey": "basic_notifications",
        "featureName": "Basic Notifications",
        "featureDescription": "Receive basic motion detection notifications",
        "enabled": true
      },
      {
        "featureKey": "live_view",
        "featureName": "Live View",
        "featureDescription": "View camera feed in real-time",
        "enabled": true
      }
    ],
    "trialEnabled": true,
    "trialDurationDays": 30,
    "trialFeatures": [
      {
        "featureKey": "advanced_ai",
        "featureName": "Advanced AI Detection",
        "featureDescription": "AI-powered object and person detection",
        "enabled": true
      },
      {
        "featureKey": "extended_storage",
        "featureName": "Extended Cloud Storage",
        "featureDescription": "30 days of cloud video storage",
        "enabled": true
      }
    ]
  }
}
```

## Error Codes
| Error Code | Description |
|------------|-------------|
| -1 | General error |
| -2 | Authentication error |
| -5 | Network error |

## Usage Context
This endpoint is typically used in the following scenarios:
- During app initialization to determine which free features to enable
- When displaying feature availability in the app interface
- Before attempting to use a feature to check if it's available in the free tier
- When showing subscription upgrade prompts to highlight premium features
- During new user onboarding to explain available free features

## Related Endpoints
- `queryDeviceCloudVipInfo` - Queries premium/VIP services for a specific device
- `getVipUserServiceInfo` - Gets VIP service information for a user
- `getTierListV4` - Gets detailed subscription tier options

## Implementation Notes
Unlike device-specific endpoints that require a serial number, this endpoint provides global information about what services are available in the free tier of the Vicohome platform. The response helps the application determine what features to display or enable without requiring a subscription.

The endpoint is called by the DeviceManageCore.queryFreeTier() method and doesn't require any device-specific parameters. This suggests it provides platform-wide information rather than device-specific capabilities.