#!/bin/bash

# Example: Load testing with mozzy

echo "=== Load Testing Examples ==="
echo ""

echo "1. Fixed number of requests with 10 concurrent workers"
mozzy load https://jsonplaceholder.typicode.com/users/1 --requests 100 --concurrent 10

echo ""
echo "2. Duration-based testing (run for 30 seconds)"
mozzy load https://jsonplaceholder.typicode.com/users/1 --duration 30s --concurrent 5

echo ""
echo "3. High concurrency test"
mozzy load https://jsonplaceholder.typicode.com/posts --requests 500 --concurrent 50

echo ""
echo "4. Load test with authentication"
mozzy load https://api.example.com/protected --requests 100 --concurrent 10 --auth "Bearer your-token-here"
