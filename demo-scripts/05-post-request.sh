#!/bin/bash
# Demo 5: POST Request with JSON

clear
echo "📤 Example 5: POST Request with JSON"
echo "====================================="
echo ""
sleep 2

echo "$ mozzy POST https://jsonplaceholder.typicode.com/posts --json '{\"title\":\"Hello World\",\"body\":\"This is mozzy\",\"userId\":1}' --color"
sleep 1
mozzy POST https://jsonplaceholder.typicode.com/posts --json '{"title":"Hello World","body":"This is mozzy","userId":1}' --color
echo ""
sleep 2

echo "✅ Done! Created a new post"
