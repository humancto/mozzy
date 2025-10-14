#!/bin/bash

# Example: Exporting saved requests and workflows

echo "=== Export Examples ==="
echo ""

echo "1. First, save a request to your collection"
mozzy GET https://jsonplaceholder.typicode.com/users/1 --save get-user

echo ""
echo "2. Export saved request as curl command"
mozzy export get-user --format curl

echo ""
echo "3. Export saved request as Postman collection"
mozzy export get-user --format postman > postman-collection.json

echo ""
echo "4. Export a workflow to curl commands"
mozzy export examples/workflows/api-workflow.yaml --format curl
