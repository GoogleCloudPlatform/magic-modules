#!/bin/bash
set -e

OUTPUT_PATH="/Users/zhenhuali/go/src/github.com/GoogleCloudPlatform/container_node_pool"

echo "Running make clean-tgc..."
make clean-tgc OUTPUT_PATH="$OUTPUT_PATH"

echo "Running make tgc..."
make tgc OUTPUT_PATH="$OUTPUT_PATH"

echo "Switching to output directory..."
cd "$OUTPUT_PATH"
# exec $SHELL # Uncomment to stay in directory if running as executable
