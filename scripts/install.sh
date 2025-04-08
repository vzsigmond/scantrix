#!/bin/bash

set -e

APP_NAME="scantrix"
GITHUB_REPO="vzsigmond/scantrix"
BIN_NAME="$APP_NAME"
DEST_DIR="/usr/local/bin"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
elif [[ "$ARCH" == arm* || "$ARCH" == aarch64 ]]; then
  ARCH="arm64"
else
  echo "❌ Unsupported architecture: $ARCH"
  exit 1
fi

EXT=""
if [ "$OS" = "windows" ]; then
  EXT=".exe"
  DEST_DIR="/c/ProgramData/scantrix"
fi

TARGET_NAME="${APP_NAME}-${OS}-${ARCH}${EXT}"
DOWNLOAD_URL="https://github.com/$GITHUB_REPO/releases/latest/download/$TARGET_NAME"

echo "⬇️ Downloading $DOWNLOAD_URL..."
curl -L "$DOWNLOAD_URL" -o "$BIN_NAME$EXT"

echo "🔧 Installing to $DEST_DIR/$BIN_NAME$EXT"
sudo mkdir -p "$DEST_DIR"
sudo mv "$BIN_NAME$EXT" "$DEST_DIR/$BIN_NAME$EXT"
sudo chmod +x "$DEST_DIR/$BIN_NAME$EXT"

if [ "$OS" != "windows" ]; then
  sudo ln -sf "$DEST_DIR/$BIN_NAME$EXT" "/usr/local/bin/$BIN_NAME"
fi

echo "✅ Installed: $BIN_NAME"
echo "👉 Run '$BIN_NAME --help' to get started."
