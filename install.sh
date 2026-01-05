#!/bin/sh
set -e

REPO="sanjay-subramanya/drift"
BINARY_NAME="drift"

# 1. Detect Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then ARCH="amd64"; elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then ARCH="arm64"; fi

# 2. Get latest version tag
LATEST_TAG=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

# 3. Download
URL="https://github.com/$REPO/releases/download/$LATEST_TAG/${BINARY_NAME}_${OS}_${ARCH}.tar.gz"
echo "Downloading $BINARY_NAME $LATEST_TAG..."
curl -L "$URL" -o "${BINARY_NAME}.tar.gz"

# 4. Install
tar -xzf "${BINARY_NAME}.tar.gz" "$BINARY_NAME"
chmod +x "$BINARY_NAME"
sudo mv "$BINARY_NAME" /usr/local/bin/
rm "${BINARY_NAME}.tar.gz"

echo "Successfully installed $BINARY_NAME to /usr/local/bin/"