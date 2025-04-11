#!/bin/bash
set -e

TARGET=linux-x64
OUTPUT_DIR="../pkg/php-ast/$TARGET"

# Move to the script directory
cd "$(dirname "$0")"

echo "🛠️ Building PHP bundle for $TARGET..."
docker build -f ../Dockerfile.php83-linux -t php-ast-bundle-linux ..

echo "📦 Creating container..."
ID=$(docker create php-ast-bundle-linux)

mkdir -p "$OUTPUT_DIR/lib"

docker cp "$ID:/export/php" "$OUTPUT_DIR/bin/php"
docker cp "$ID:/export/php.ini" "$OUTPUT_DIR/php.ini"
docker cp "$ID:/export/lib" "$OUTPUT_DIR/"

echo "🧹 Cleaning up..."
docker rm "$ID"

echo "✅ Bundle available at $OUTPUT_DIR"
