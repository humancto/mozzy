#!/bin/bash
# JWT Sign and Verify Example

SECRET="my-super-secret-key"

# Create a payload file
cat > /tmp/payload.json << EOF
{
  "user": "alice",
  "role": "admin",
  "email": "alice@example.com"
}
EOF

echo "=== Signing JWT with Payload ==="
cat /tmp/payload.json
echo ""

TOKEN=$(mozzy jwt sign /tmp/payload.json --secret $SECRET)
echo "Generated Token: $TOKEN"
echo ""

echo "=== Verifying JWT Signature ==="
mozzy jwt verify $TOKEN --secret $SECRET
echo ""

echo "=== Decoding JWT ==="
mozzy jwt decode $TOKEN

# Cleanup
rm /tmp/payload.json
