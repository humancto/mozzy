#!/bin/bash
# POST request example

echo "=== POST with JSON Body ==="
mozzy POST https://jsonplaceholder.typicode.com/posts \
  --json '{
    "title": "My First Post",
    "body": "This is the content of my post",
    "userId": 1
  }'

echo -e "\n=== POST with Custom Headers ==="
mozzy POST https://jsonplaceholder.typicode.com/posts \
  --header "X-Custom-Header: MyValue" \
  --json '{"title":"Test","body":"Body","userId":1}'
