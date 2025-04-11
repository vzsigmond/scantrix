#!/bin/bash
set -e

# php.sh - A wrapper that auto-detects the OS/arch and picks the right php-ast/<platform> folder.

# 1) Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m | tr '[:upper:]' '[:lower:]')"

PLATFORM=""
case "$OS" in
  linux*)
    # e.g., linux-x64, linux-arm64, etc.
    if [[ "$ARCH" == "x86_64" ]]; then
      PLATFORM="linux-x64"
    elif [[ "$ARCH" == "aarch64" ]] || [[ "$ARCH" == "arm64" ]]; then
      PLATFORM="linux-arm64"
    else
      echo "Unsupported Linux arch: $ARCH" >&2
      exit 1
    fi
    ;;
  darwin*)
    # e.g., darwin-x64, darwin-arm64
    if [[ "$ARCH" == "arm64" ]]; then
      PLATFORM="darwin-arm64"
    else
      PLATFORM="darwin-x64"
    fi
    ;;
  msys*|cygwin*|mingw*|nt|win*)
    # For Windows environment
    # you might do more checks, but typically "windows-x64" is the main target
    PLATFORM="windows-x64"
    ;;
  *)
    echo "Unsupported OS: $OS" >&2
    exit 1
    ;;
esac

# 2) Determine script directory
DIR="$(cd "$(dirname "$0")/.." && pwd)"
# 3) Build the paths to bin/ and lib/ inside php-ast/<platform>
BIN_DIR="$DIR/$PLATFORM/bin"
LIB_DIR="$DIR/$PLATFORM/lib"
# 4) Check if the bin/ directory exists

if [ ! -d "$BIN_DIR" ]; then
  echo "Error: PHP binary directory '$BIN_DIR' does not exist." >&2
  exit 1
fi

# 5) Set library path so the embedded PHP can find ast.so + other libs
export LD_LIBRARY_PATH="$LIB_DIR:$LD_LIBRARY_PATH"

# 6) Actually run "php" from the bin/ directory
#    - extension_dir points to LIB_DIR
#    - extension=... ensures ast.so is loaded
#    - -c "$DIR/php.ini" uses your local php.ini
"$BIN_DIR/php" \
  -d extension_dir="$LIB_DIR" \
  -d extension="$LIB_DIR/ast.so" \
  -c "$DIR/php.ini" \
  "$@"