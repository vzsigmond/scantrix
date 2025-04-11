#!/bin/bash
set -e

# parser.sh
# Usage: ./parser.sh /path/to/phpfile.php
#
# This script calls php.sh with a custom AST parser script (ASTParser.php).
# The parser script receives the path to a PHP file, parses it, and outputs JSON.

# 1) Check if a file path is provided
if [ $# -lt 1 ]; then
  echo "Usage: $0 /path/to/phpfile.php" >&2
  exit 1
fi

PHP_FILE="$1"
if [ ! -f "$PHP_FILE" ]; then
  echo "Error: file '$PHP_FILE' does not exist." >&2
  exit 1
fi

# 2) Determine this script's directory so we can reference php.sh + ASTParser.php
BIN_DIR="$(cd "$(dirname "$0")" && pwd)"
SCRIPT_DIR="$(cd "$BIN_DIR/.." && pwd)"
echo $SCRIPT_DIR

# 3) Call php.sh, pointing it to "ASTParser.php" (or "ast_script.php") plus the user’s PHP file
"$BIN_DIR/php.sh" "$SCRIPT_DIR/ASTParser.php" "$PHP_FILE"
