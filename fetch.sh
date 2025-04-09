#!/usr/bin/env bash

EMAIL=""
PASSWORD=""
LOGIN_TYPE=0

# Login and obtain Bearer Token

echo "Authenticating with Vicohome API..."
RESPONSE=$(curl -s -X POST "https://api-us.vicohome.io/account/login" \
  -H "accept: application/json" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\",\"loginType\":$LOGIN_TYPE}")

BEARER_TOKEN=$(echo "$RESPONSE" | jq -r '.data.token.token')

# Filter for last 24 hours of birds

# API endpoint
START_TIMESTAMP=$( date -v-24H +%s)
END_TIMESTAMP=$(date +%s)

# Create the request body with the calculated timestamps
REQUEST_BODY='{
  "startTimestamp": "'$START_TIMESTAMP'",
  "endTimestamp": "'$END_TIMESTAMP'",
  "language": "en",
  "countryNo": "US"
}'

# Make the API request
echo "Fetching events from $START_DATE to $END_DATE from Vicohome API..."
RESPONSE=$(curl -s -X POST https://api-us.vicohome.io/library/newselectlibrary \
  -H "accept: application/json" \
  -H "Content-Type: application/json" \
  -H "Authorization: $BEARER_TOKEN" \
  -d "$REQUEST_BODY")

echo "----------------------------------------"
echo "$RESPONSE" | jq

echo "$RESPONSE" | jq > last-24-hours.json
echo "Response saved to last-24-hours.json"

# View Single Event by Trace
echo "Fetching event data for $TRACE_ID from Vicohome API..."
TRACE_ID="0185942217441562400ny2nw9yIw0"

REQUEST_BODY='{
  "traceId": "0185942217441562400ny2nw9yIw0",
  "language": "en",
  "countryNo": "US"
}'

RESPONSE=$(curl -s -X POST https://api-us.vicohome.io/library/newselectsinglelibrary \
  -H "accept: application/json" \
  -H "Content-Type: application/json" \
  -H "Authorization: $BEARER_TOKEN" \
  -d "$REQUEST_BODY")

echo "----------------------------------------"
echo "$RESPONSE" | jq

echo "$RESPONSE" | jq > 0185942217441562400ny2nw9yIw0.json
echo "Response saved to 0185942217441562400ny2nw9yIw0.json"