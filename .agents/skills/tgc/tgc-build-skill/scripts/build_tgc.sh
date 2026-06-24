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

# Verify active task list GEMINI entrypoints before building
ACTIVE_TASK_MD=$(find "$HOME/.gemini/jetski/brain" -name task.md -type f -print0 2>/dev/null | xargs -0 stat -f "%m %N" 2>/dev/null | sort -rn | head -n 1 | cut -d' ' -f2-)

if [ ! -z "$ACTIVE_TASK_MD" ]; then
  echo "Verifying task list at $ACTIVE_TASK_MD..."
  "$MAGIC_MODULES_PATH/.agents/scripts/tgc/verify_task_list.py" "$ACTIVE_TASK_MD"
else
  echo "Warning: Active task.md could not be detected automatically. Skipping task list verification."
fi

# Automatically verify field ordering and test configuration for all modified YAML product configurations
echo "Checking for modified YAML product configurations..."
CHANGED_YAMLS=$(git diff --name-only 2>/dev/null | grep -E "^mmv1/products/.*\.ya?ml$" || true)
if [ ! -z "$CHANGED_YAMLS" ]; then
  for yaml in $CHANGED_YAMLS; do
    if [ -f "$yaml" ]; then
      echo "Verifying field ordering compliance for $yaml..."
      "$MAGIC_MODULES_PATH/.agents/scripts/tgc/verify_yaml_field_order.py" mmv1/api/resource.go "$yaml"
      
      echo "Verifying test configurations for $yaml..."
      "$MAGIC_MODULES_PATH/.agents/scripts/tgc/verify_test_configs.py" "$yaml"
    fi
  done
fi


echo "[Phase 1] Generating TGC Code from Magic Modules..."
cd "$MAGIC_MODULES_PATH"
make tgc OUTPUT_PATH="$TGC_PATH"

echo "[Phase 2] Compiling the TGC Binary downstream..."
cd "$TGC_PATH"
make mod-clean
make build

echo "TGC build compiled successfully!"

echo "[Phase 3] Executing selective unit tests..."
TGC_DIR="$TGC_PATH" "$MAGIC_MODULES_PATH/.agents/skills/tgc/tgc-build-skill/scripts/run_changed_folders_tests.sh"
