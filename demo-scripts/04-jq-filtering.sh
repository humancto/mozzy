#!/bin/bash
# Demo 4: JQ Filtering

clear
echo "🔍 Example 4: Inline JQ Filtering"
echo "=================================="
echo ""
sleep 2

echo "$ mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .name --color"
sleep 1
mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .name --color
echo ""
sleep 2

echo "$ mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .address.city --color"
sleep 1
mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .address.city --color
echo ""
sleep 2

echo "$ mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .company.name --color"
sleep 1
mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .company.name --color
echo ""
sleep 2

echo "✅ Done! No need for external jq tool"
