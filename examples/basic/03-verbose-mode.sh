#!/bin/bash
# Verbose mode example - shows headers and timing

echo "=== GET with Verbose Output ==="
mozzy GET https://jsonplaceholder.typicode.com/users/1 --verbose

echo -e "\n=== This shows: ==="
echo "  - Request headers"
echo "  - Response headers"
echo "  - Timing breakdown (DNS, TLS, server, transfer)"
