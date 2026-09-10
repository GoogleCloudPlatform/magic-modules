#!/bin/bash
set -e

RESOURCE_NAME=$1
GENERATED_TEST_FILE=$2

if [ $# -lt 2 ]; then
    echo "Usage: $0 <ResourceName> <GeneratedTestFilePath>"
    echo "Example: $0 DialogflowCXAgent test/services/dialogflowcx/dialogflowcx_agent_generated_test.go"
    exit 1
fi

TGC_DIR="${TGC_DIR:-$(pwd)}" # Assume running from downstream root or provided

# Verify generated test file exists
if [ ! -f "$GENERATED_TEST_FILE" ]; then
    echo "❌ Error: Generated test file not found at $GENERATED_TEST_FILE"
    exit 1
fi

echo "Checking tests for $RESOURCE_NAME..."

MISSING=0
FOUND_ANY=0

# Search all metadata files
for METADATA_FILE in $(find "$TGC_DIR/test" -name "tests_metadata_*.json" -maxdepth 1); do
    echo "Searching in $METADATA_FILE..."
    EXPECTED_TESTS=$(grep -E "\"test_name\": *\"TestAcc${RESOURCE_NAME}" "$METADATA_FILE" | sed -E 's/.*"test_name": *\"([^"]+)".*/\1/' | sort -u)

    if [ -z "$EXPECTED_TESTS" ]; then
        continue
    fi

    FOUND_ANY=1
    for TEST in $EXPECTED_TESTS; do
        if ! grep -q "$TEST" "$GENERATED_TEST_FILE"; then
            echo "❌ Missing test: $TEST (found in $METADATA_FILE)"
            MISSING=$((MISSING+1))
        fi
    done
done

if [ $FOUND_ANY -eq 0 ]; then
    echo "No tests found for $RESOURCE_NAME in any metadata file."
    exit 0
fi

if [ $MISSING -eq 0 ]; then
    echo "✅ All tests present in generated file."
else
    echo "Total missing tests: $MISSING"
    exit 1
fi
