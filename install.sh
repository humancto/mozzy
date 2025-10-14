#!/bin/bash

echo "🚀 Installing mozzy..."
echo ""

# Detect OS and architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

if [ "$OS" = "Darwin" ]; then
    if [ "$ARCH" = "arm64" ]; then
        PLATFORM="darwin_arm64"
    else
        PLATFORM="darwin_amd64"
    fi
elif [ "$OS" = "Linux" ]; then
    if [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
        PLATFORM="linux_arm64"
    else
        PLATFORM="linux_amd64"
    fi
else
    echo "❌ Unsupported OS: $OS"
    exit 1
fi

echo "Detected platform: $PLATFORM"
echo ""

# Download URL
VERSION="1.1.0"
URL="https://github.com/humancto/mozzy/releases/download/v${VERSION}/mozzy_${VERSION}_${PLATFORM}.tar.gz"

echo "Downloading mozzy v${VERSION}..."
curl -L "$URL" -o /tmp/mozzy.tar.gz

if [ $? -ne 0 ]; then
    echo "❌ Download failed"
    exit 1
fi

echo "Extracting..."
cd /tmp
tar -xzf mozzy.tar.gz

echo "Installing to /usr/local/bin..."
sudo mv mozzy /usr/local/bin/
sudo chmod +x /usr/local/bin/mozzy

# Cleanup
rm -f /tmp/mozzy.tar.gz

echo ""
echo "✅ mozzy installed successfully!"
echo ""
echo "Try it:"
echo "  mozzy --version"
echo "  mozzy GET https://jsonplaceholder.typicode.com/users/1 --color"
echo ""
