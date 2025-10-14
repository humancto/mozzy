#!/bin/bash

# Example: File uploads with mozzy

echo "=== File Upload Examples ==="
echo ""

echo "1. Upload single file"
# mozzy upload https://httpbin.org/post -f test.jpg

echo ""
echo "2. Upload with custom field name"
# mozzy upload https://httpbin.org/post -f avatar.jpg --field-name profileImage

echo ""
echo "3. Upload multiple files"
# mozzy upload https://httpbin.org/post -f file1.jpg -f file2.png

echo ""
echo "4. Upload with form data"
# mozzy upload https://httpbin.org/post -f resume.pdf --data "name=John Doe" --data "email=john@example.com"

echo ""
echo "5. Upload with authentication"
# mozzy upload https://api.example.com/upload -f document.pdf --auth "your-token-here"

echo ""
echo "6. Upload with custom headers"
# mozzy upload https://api.example.com/upload -f file.zip --header "X-API-Key: secret123"

echo ""
echo "7. Upload without progress (for scripting)"
# mozzy upload https://api.example.com/upload -f data.json --no-progress

echo ""
echo "8. Upload large files (shows progress bar)"
# mozzy upload https://api.example.com/upload -f large-video.mp4

echo ""
echo "Real example with httpbin.org:"
echo "Creating test file..."
echo "test content" > /tmp/test.txt

echo "Uploading..."
mozzy upload https://httpbin.org/post -f /tmp/test.txt --data "description=Test upload" --jq .files

rm /tmp/test.txt
