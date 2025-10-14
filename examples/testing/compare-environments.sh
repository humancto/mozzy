#!/bin/bash
# Compare API responses between prod and staging

ENDPOINT="/users/1"

echo "🔍 Comparing API responses across environments..."
echo ""

# Fetch from both environments
echo "Fetching from production..."
mozzy GET "https://jsonplaceholder.typicode.com${ENDPOINT}" > /tmp/prod-response.json

echo "Fetching from staging (simulated with different user)..."
mozzy GET "https://jsonplaceholder.typicode.com/users/2" > /tmp/staging-response.json

echo ""
echo "📊 Running diff..."
echo ""

# Compare
mozzy diff /tmp/prod-response.json /tmp/staging-response.json

# Cleanup
rm /tmp/prod-response.json /tmp/staging-response.json
