#!/bin/bash

echo "🚀 Welcome to mozzy - Your Postman for the Terminal!"
echo ""
echo "=================================================="
echo ""

echo "1️⃣  Test GET request with beautiful colors:"
echo "   $ mozzy GET https://jsonplaceholder.typicode.com/users/1"
echo ""
./mozzy GET https://jsonplaceholder.typicode.com/users/1
echo ""
echo "Press Enter to continue..."
read

echo "2️⃣  Test --jq filtering (extract just the name):"
echo "   $ mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .name"
echo ""
./mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .name
echo ""
echo "Press Enter to continue..."
read

echo "3️⃣  Test nested --jq filtering:"
echo "   $ mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .address.city"
echo ""
./mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .address.city
echo ""
echo "Press Enter to continue..."
read

echo "4️⃣  Test array filtering:"
echo "   $ mozzy GET https://jsonplaceholder.typicode.com/users --jq .[0].email"
echo ""
./mozzy GET https://jsonplaceholder.typicode.com/users --jq .[0].email
echo ""
echo "Press Enter to continue..."
read

echo "5️⃣  Test 404 error with helpful tips:"
echo "   $ mozzy GET https://jsonplaceholder.typicode.com/users/999"
echo ""
./mozzy GET https://jsonplaceholder.typicode.com/users/999 2>&1
echo ""
echo "Press Enter to continue..."
read

echo "6️⃣  Test 401 error with authentication tips:"
echo "   $ mozzy POST https://httpbin.org/status/401"
echo ""
./mozzy POST https://httpbin.org/status/401 2>&1
echo ""
echo "Press Enter to continue..."
read

echo "7️⃣  Save a request to your collection:"
echo "   $ mozzy save github-user GET https://api.github.com/users/torvalds --desc 'Get Linus Torvalds profile'"
echo ""
./mozzy save github-user GET https://api.github.com/users/torvalds --desc "Get Linus Torvalds profile"
echo ""
echo "Press Enter to continue..."
read

echo "8️⃣  List your saved requests:"
echo "   $ mozzy list"
echo ""
./mozzy list
echo ""
echo "Press Enter to continue..."
read

echo "9️⃣  Execute a saved request:"
echo "   $ mozzy exec github-user"
echo ""
./mozzy exec github-user 2>&1 | head -20
echo ""
echo "Press Enter to continue..."
read

echo "🔟 Test POST with JSON:"
echo "   $ mozzy POST https://jsonplaceholder.typicode.com/posts --json '{\"title\":\"Hello\",\"body\":\"World\",\"userId\":1}'"
echo ""
./mozzy POST https://jsonplaceholder.typicode.com/posts --json '{"title":"Hello from mozzy","body":"Testing POST requests","userId":1}'
echo ""
echo "Press Enter to continue..."
read

echo "1️⃣1️⃣  Create and run a workflow:"
cat > /tmp/demo-workflow.yaml << 'YAML'
name: GitHub API Demo
steps:
  - name: Get GitHub user
    method: GET
    url: https://api.github.com/users/octocat
    capture:
      login: .login
      repos: .public_repos

  - name: Get user repos
    method: GET
    url: https://api.github.com/users/{{login}}/repos
YAML

echo "   Created workflow file:"
cat /tmp/demo-workflow.yaml
echo ""
echo "   $ mozzy run /tmp/demo-workflow.yaml"
echo ""
./mozzy run /tmp/demo-workflow.yaml 2>&1 | head -30
echo ""

echo "=================================================="
echo ""
echo "✅ Demo complete! You've seen:"
echo "   • Colorized JSON output"
echo "   • --jq filtering (nested paths, arrays)"
echo "   • Helpful error messages with tips"
echo "   • Collections (save/list/exec)"
echo "   • POST requests with JSON"
echo "   • YAML workflows with captures"
echo ""
echo "📚 More commands to try:"
echo "   mozzy --help"
echo "   mozzy jwt decode <token>"
echo "   mozzy history"
echo "   mozzy --env prod GET /api/users"
echo ""
echo "🎉 Enjoy using mozzy!"
