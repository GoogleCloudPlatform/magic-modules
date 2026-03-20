#!/bin/bash

# Script to run integration tests for TGC
# Usage: ./run_integration_test.sh <test-path> <test-target>
#
# Examples:
# ./run_integration_test.sh ./test/services/alloydb TestAccAlloydbBackup

if [ $# -lt 2 ]; then
    echo "Usage: $0 <test-path> <test-target>"
    echo "Example: $0 ./test/services/alloydb TestAccAlloydbBackup"
    exit 1
fi

TEST_PATH=$1
TEST_TARGET=$2

# Assuming standard GOPATH unless explicitly provided in the environment
TGC_DIR="${TGC_DIR:-$GOPATH/src/github.com/GoogleCloudPlatform/terraform-google-conversion}"

if [ ! -d "$TGC_DIR" ]; then
    echo "Error: TGC directory not found at $TGC_DIR"
    exit 1
fi

# Ensure the log file is unique
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
LOG_FILE="${TEST_TARGET}_${TIMESTAMP}.log"
LOG_DIR="$TGC_DIR/debug_output/raw_logs"

echo "Creating log directory $LOG_DIR..."
mkdir -p "$LOG_DIR"

echo "Running integration test $TEST_TARGET in $TEST_PATH..."
echo "Logs will be saved to: $LOG_DIR/$LOG_FILE"

cd "$TGC_DIR"
export WRITE_FILES=true
make test-integration-local TESTPATH="$TEST_PATH" TESTARGS="-run=$TEST_TARGET" > "$LOG_DIR/$LOG_FILE" 2>&1

echo "Test execution complete (or failed). Check the log file for outputs."
