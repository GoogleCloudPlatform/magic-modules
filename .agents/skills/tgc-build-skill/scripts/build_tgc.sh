#!/bin/bash
set -e

# Use current directory for magic-modules, and default GOPATH for TGC downstream
MAGIC_MODULES_PATH="${MAGIC_MODULES_PATH:-$(pwd)}"

TGC_PATH="${1:-$TGC_PATH}"
if [ -z "$TGC_PATH" ]; then
  echo "Error: TGC_PATH is not set and no path provided as argument."
  exit 1
fi

echo "Starting TGC Build Process..."

echo "[Phase 1] Generating TGC Code from Magic Modules..."
cd "$MAGIC_MODULES_PATH"
make clean-tgc OUTPUT_PATH="$TGC_PATH"
make tgc OUTPUT_PATH="$TGC_PATH"

echo "[Phase 2] Compiling the TGC Binary downstream..."
cd "$TGC_PATH"
make mod-clean
make build

echo "TGC build compiled successfully!"
