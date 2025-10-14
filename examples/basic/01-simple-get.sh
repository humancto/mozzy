#!/bin/bash
# Simple GET request example

echo "=== Basic GET Request ==="
mozzy GET https://jsonplaceholder.typicode.com/users/1

echo -e "\n=== GET with JQ Filter ==="
mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .name

echo -e "\n=== GET with Nested JQ Filter ==="
mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .address.city
