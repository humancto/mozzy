#!/bin/bash
# Demo 1: GitHub API Exploration

clear
echo "🐙 Example 1: Exploring GitHub's API"
echo "======================================"
echo ""
sleep 2

echo "$ mozzy GET https://api.github.com/users/torvalds --color"
sleep 1
mozzy GET https://api.github.com/users/torvalds --color
echo ""
sleep 2

echo "$ mozzy GET https://api.github.com/users/torvalds --jq .name"
sleep 1
mozzy GET https://api.github.com/users/torvalds --jq .name
echo ""
sleep 2

echo "$ mozzy GET https://api.github.com/users/torvalds/repos --jq '.[0].name'"
sleep 1
mozzy GET https://api.github.com/users/torvalds/repos --jq '.[0].name'
echo ""
sleep 2

echo "✅ Done!"
