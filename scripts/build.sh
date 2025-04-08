#!/bin/bash

# Build Scantrix for major platforms
set -e

OUTPUT_DIR="bin"
APP_NAME="scantrix"
CMD_PATH="./cmd/scantrix"

mkdir -p "$OUTPUT_DIR"

platforms=(
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
)

echo "🔨 Building Scantrix for multiple platforms..."

for platform in "${platforms[@]}"
do
  IFS="/" read -r GOOS GOARCH <<< "$platform"
  output_name="$OUTPUT_DIR/${APP_NAME}-${GOOS}-${GOARCH}"
  if [ "$GOOS" = "windows" ]; then
    output_name="${output_name}.exe"
  fi

  echo "➡️  Building for $GOOS/$GOARCH -> $output_name"

  env GOOS="$GOOS" GOARCH="$GOARCH" CGO_ENABLED=0 go build -buildvcs=false -o "$output_name" "$CMD_PATH"
done

echo "✅ All builds completed. Binaries saved in ./$OUTPUT_DIR"
