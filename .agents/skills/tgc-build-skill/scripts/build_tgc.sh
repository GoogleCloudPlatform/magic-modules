#!/bin/bash
set -e

# Use current directory for magic-modules, and default GOPATH for TGC downstream
MAGIC_MODULES_PATH="${MAGIC_MODULES_PATH:-$(pwd)}"

if [ -n "$1" ]; then
  TGC_PATH="$1"
else
  TGC_PATH="${TGC_PATH:-$GOPATH/src/github.com/GoogleCloudPlatform/terraform-google-conversion}"
fi

echo "Starting TGC Build Process..."

echo "[Phase 1] Generating TGC Code from Magic Modules..."
cd "$MAGIC_MODULES_PATH"
make tgc OUTPUT_PATH="$TGC_PATH"

echo "[Phase 2] Compiling the TGC Binary downstream..."
cd "$TGC_PATH"
make mod-clean
make build

echo "TGC build compiled successfully!"
