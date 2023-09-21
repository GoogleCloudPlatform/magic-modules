#!/bin/bash

# Check if there's at least one argument
if [ "$#" -eq 0 ]; then
    echo "No arguments provided"
    exit 1
fi

# Get the directory of the current script
DIR="$(dirname $(realpath $0))"

# Construct the path to the Go program
GO_PROGRAM="$DIR/../../../magician/"

pushd $GO_PROGRAM

set -x
# Pass all arguments to the child command
go run . "$@"
set +x
