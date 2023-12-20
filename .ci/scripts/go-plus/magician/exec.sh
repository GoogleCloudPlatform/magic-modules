#!/bin/bash

# Get the directory of the current script
DIR="$(dirname $(realpath $0))"

# Construct the path to the Go program directory and binary
GO_PROGRAM_DIR="$DIR/../../../magician"
GO_BINARY="$GO_PROGRAM_DIR/magician_binary"

pushd $GO_PROGRAM_DIR

set -x
# Check if the binary exists
if [ ! -f "$GO_BINARY" ]; then
    # If it doesn't exist, compile the binary
    echo "Building the magician binary at $GO_BINARY"
    go build -o "$GO_BINARY"
fi

# If there are no arguments only compile the binary
if [ "$#" -eq 0 ]; then
    echo "No arguments provided"
    exit 0
fi

# Run the binary and pass all arguments
$GO_BINARY "$@"
set +x
