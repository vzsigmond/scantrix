#!/bin/bash
set -e

APP_NAME="scantrix"
OUTPUT_DIR="build"
CMD_PATH="./cmd/scantrix"
PHP_AST_DIR="./pkg/php-ast"

PLATFORMS=(
  "linux/amd64"
  # Add more platforms like: "darwin/arm64"
)

# Ensure output dir
mkdir -p "$OUTPUT_DIR"

echo "🔨 Building Scantrix in Docker..."

for platform in "${PLATFORMS[@]}"; do
  IFS="/" read -r GOOS GOARCH <<< "$platform"
  TARGET_DIR="${APP_NAME}-${GOOS}-${GOARCH}"
  OUT_DIR="${OUTPUT_DIR}/${TARGET_DIR}"
  mkdir -p "$OUT_DIR"

  echo "➡️  Building for $GOOS/$GOARCH"

  docker run --rm -v "$PWD":/app -w /app \
    -e GOOS="$GOOS" -e GOARCH="$GOARCH" -e CGO_ENABLED=0 \
    golang:1.24-alpine \
    sh -c "go build -buildvcs=false -o ${OUT_DIR}/${APP_NAME} $CMD_PATH"

  # Bundle PHP AST if available for platform
  case "$platform" in
    "linux/amd64")
      echo "📦 Bundling php-ast for $platform"
      mkdir -p "$OUT_DIR/php-ast"
      cp -r "$PHP_AST_DIR/linux-x64/"* "$OUT_DIR/php-ast/"
      ;;
    *)
      echo "⚠️  No php-ast bundle for $platform"
      ;;
  esac
done

echo "✅ Done! Output in ./$OUTPUT_DIR/"
