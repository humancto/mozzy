#!/bin/bash
# Demo 3: Verbose Mode with Timing

clear
echo "⏱️  Example 3: Verbose Mode & Timing"
echo "===================================="
echo ""
sleep 2

echo "$ mozzy GET https://jsonplaceholder.typicode.com/users/1 --verbose --color"
sleep 1
mozzy GET https://jsonplaceholder.typicode.com/users/1 --verbose --color
echo ""
sleep 2

echo "✅ Done! Notice the timing breakdown and headers"
