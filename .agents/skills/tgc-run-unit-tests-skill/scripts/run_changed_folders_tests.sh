#!/bin/bash
set -e

# Default to current directory if TGC_DIR not set
TGC_DIR="${TGC_DIR:-$(pwd)}"
cd "$TGC_DIR"

echo "Checking for changes in $TGC_DIR..."

# Get changed files (staged and unstaged)
# Also include untracked files to be thorough
CHANGED_FILES=$(git diff --name-only HEAD; git ls-files --others --exclude-standard)

if [ -z "$CHANGED_FILES" ]; then
    echo "No changes detected."
    exit 0
fi

echo "Changed files:"
echo "$CHANGED_FILES"
echo "----------------"

# Extract directory paths
DIRS=$(echo "$CHANGED_FILES" | while read f; do dirname "$f"; done | sort -u)

for DIR in $DIRS; do
    if [ "$DIR" = "." ]; then
        continue
    fi
    if [[ "$DIR" == test/services/* ]]; then
        echo "ℹ️ Skipping service tests in $DIR."
        continue
    fi
    # Check if directory exists
    if [ -d "$DIR" ]; then
        # Check if it has tests
        if find "$DIR" -name "*_test.go" -print -quit | grep -q .; then
            echo "🚀 Running unit tests for $DIR..."
            make test-local TEST="./$DIR/..."
        else
            echo "ℹ️ No tests found in $DIR tree, skipping."
        fi
    else
        echo "ℹ️ $DIR is not a directory, skipping."
    fi
done
