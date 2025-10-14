#!/bin/bash
# JWT Decode Example

# Sample JWT token (this is a public example token)
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

echo "=== Decoding JWT Token ==="
echo "Token: $TOKEN"
echo ""

mozzy jwt decode $TOKEN

echo -e "\n=== What this shows: ==="
echo "  - Header (algorithm, type)"
echo "  - Payload (claims)"
echo "  - Expiration time (if present)"
