#!/bin/bash

set -e

APP_NAME="scantrix"
GITHUB_USER="vzsigmond"
GITHUB_REPO="scantrix"
INSTALL_DIR="/usr/local/bin"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Normalize architecture
if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
elif [[ "$ARCH" == arm64* || "$ARCH" == aarch64 ]]; then
  ARCH="arm64"
else
  echo "❌ Unsupported architecture: $ARCH"
  exit 1
fi

EXT=""
if [ "$OS" = "windows" ]; then
  EXT=".exe"
fi

BINARY_NAME="${APP_NAME}-${OS}-${ARCH}${EXT}"
DOWNLOAD_URL="https://github.com/${GITHUB_USER}/${GITHUB_REPO}/releases/latest/download/${BINARY_NAME}"

echo "⬇️ Downloading latest Scantrix binary for $OS/$ARCH..."
curl -L "$DOWNLOAD_URL" -o "$APP_NAME$EXT"

echo "🔧 Installing to $INSTALL_DIR/$APP_NAME$EXT"
chmod +x "$APP_NAME$EXT"
sudo mv "$APP_NAME$EXT" "$INSTALL_DIR/$APP_NAME$EXT"

echo "✅ Scantrix is now installed/upgraded at: $INSTALL_DIR/$APP_NAME$EXT"
echo "👉 Run '$APP_NAME --help' to get started."
