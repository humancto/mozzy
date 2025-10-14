#!/bin/bash
# Collection Example: GitHub API

echo "=== Saving GitHub API Requests to Collection ==="

# Save multiple requests
mozzy save gh-user GET https://api.github.com/users/torvalds \
  --desc "Get Linus Torvalds profile"

mozzy save gh-repos GET https://api.github.com/users/torvalds/repos \
  --desc "Get Linus Torvalds repositories"

mozzy save gh-events GET https://api.github.com/users/torvalds/events \
  --desc "Get Linus Torvalds recent events"

echo -e "\n=== Listing Saved Requests ==="
mozzy list

echo -e "\n=== Executing Saved Request ==="
mozzy exec gh-user

echo -e "\n=== Using JQ with Saved Request ==="
mozzy exec gh-repos --jq '.[0].name'
