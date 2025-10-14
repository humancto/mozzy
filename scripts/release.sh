#!/bin/bash
set -e

echo "🚀 mozzy Release Script"
echo ""

# Check if version provided
if [ -z "$1" ]; then
    echo "Usage: ./scripts/release.sh <version>"
    echo "Example: ./scripts/release.sh 1.0.0"
    exit 1
fi

VERSION=$1
TAG="v${VERSION}"

echo "📋 Releasing mozzy ${TAG}"
echo ""

# Check if goreleaser is installed
if ! command -v goreleaser &> /dev/null; then
    echo "❌ goreleaser not found. Installing..."
    brew install goreleaser
fi

# Check if tag exists
if git rev-parse "$TAG" >/dev/null 2>&1; then
    echo "❌ Tag $TAG already exists"
    exit 1
fi

# Run tests
echo "🧪 Running tests..."
go test ./... || { echo "❌ Tests failed"; exit 1; }

# Check if we have GITHUB_TOKEN
if [ -z "$GITHUB_TOKEN" ]; then
    echo "❌ GITHUB_TOKEN not set"
    echo "Set it with: export GITHUB_TOKEN=your_token_here"
    exit 1
fi

# Tag the release
echo "🏷  Creating tag ${TAG}..."
git tag -a "$TAG" -m "Release ${TAG}"
git push origin "$TAG"

echo "📦 Building and releasing with goreleaser..."
goreleaser release --clean

echo ""
echo "✅ Release ${TAG} complete!"
echo ""
echo "📚 What's next:"
echo "  1. Check release: https://github.com/humancto/mozzy/releases/tag/${TAG}"
echo "  2. Users can install with:"
echo "     brew tap humancto/mozzy"
echo "     brew install mozzy"
echo ""
echo "  3. Once you hit 75+ stars, submit to homebrew-core:"
echo "     Then users can: brew install mozzy"
