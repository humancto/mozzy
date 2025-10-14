#!/bin/bash

# Example: File downloads with mozzy

echo "=== File Download Examples ==="
echo ""

echo "1. Simple download (auto-detect filename)"
mozzy download https://jsonplaceholder.typicode.com/users/1

echo ""
echo "2. Download with custom filename"
mozzy download https://jsonplaceholder.typicode.com/users/1 -o user-data.json

echo ""
echo "3. Download without progress bar"
mozzy download https://jsonplaceholder.typicode.com/posts --no-progress

echo ""
echo "4. Overwrite existing file"
mozzy download https://jsonplaceholder.typicode.com/users/1 -o existing.json --overwrite

echo ""
echo "5. Download large files with progress tracking"
# This would show a real progress bar with ETA and speed
mozzy download https://github.com/golang/go/archive/refs/tags/go1.21.0.zip -o golang.zip

echo ""
echo "6. Download binary files"
mozzy download https://github.com/humancto/mozzy/releases/latest/download/mozzy_darwin_arm64.tar.gz

echo ""
echo "Note: Progress bar automatically detects file size from Content-Length header"
echo "      For unknown sizes, shows downloaded bytes and speed"
