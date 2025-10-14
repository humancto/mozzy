#!/bin/bash

# Example: Retry with conditions

echo "=== Retry Conditions Examples ==="
echo ""

echo "1. Retry only on 5xx errors (default behavior)"
mozzy GET https://httpstat.us/503 --retry 3 --verbose

echo ""
echo "2. Retry on specific status code (503 Service Unavailable)"
mozzy GET https://httpstat.us/503 --retry 3 --retry-on "503" --verbose

echo ""
echo "3. Retry on rate limit (429) and server errors (5xx)"
mozzy GET https://api.example.com/data --retry 5 --retry-on "429,5xx"

echo ""
echo "4. Retry on any status >= 500"
mozzy GET https://httpstat.us/500 --retry 3 --retry-on ">=500"

echo ""
echo "5. Retry on 4xx and 5xx errors"
mozzy GET https://httpstat.us/404 --retry 2 --retry-on "4xx,5xx"

echo ""
echo "6. Retry only on network errors"
mozzy GET https://nonexistent.example.com --retry 3 --retry-on "network_error"

echo ""
echo "7. Always retry (including successful responses)"
mozzy GET https://api.example.com/flaky --retry 3 --retry-on "always"

echo ""
echo "8. Never retry (disable default retry behavior)"
mozzy GET https://httpstat.us/500 --retry 3 --retry-on "never"
