#!/bin/bash
# Demo 2: Request Collections

clear
echo "📚 Example 2: Request Collections"
echo "=================================="
echo ""
sleep 2

echo "$ mozzy save gh-torvalds GET https://api.github.com/users/torvalds --desc 'Get Linus Torvalds profile'"
sleep 1
mozzy save gh-torvalds GET https://api.github.com/users/torvalds --desc "Get Linus Torvalds profile"
echo ""
sleep 2

echo "$ mozzy list"
sleep 1
mozzy list
echo ""
sleep 2

echo "$ mozzy exec gh-torvalds"
sleep 1
mozzy exec gh-torvalds
echo ""
sleep 2

echo "✅ Done!"
