#!/bin/bash
# run_tgc_batch.sh
# Automates the entire TGC enablement, build, and test pipeline for a batch of resources.
# Usage: ./run_tgc_batch.sh <path_to_resources_txt>

if [ -z "$1" ]; then
    echo "Usage: ./run_tgc_batch.sh <path_to_resources_txt>"
    exit 1
fi

RESOURCE_FILE="$1"
if [ ! -f "$RESOURCE_FILE" ]; then
    echo "Error: File $RESOURCE_FILE not found!"
    exit 1
fi

echo "=== Processing New TGC Resource Batch ==="

# 1. Copy the input file to standard location for the Python script
cp "$RESOURCE_FILE" resources_to_add.txt

# 2. Run the mapping script
echo "=== 1. Mapping Resources to Products ==="
python3 map.py
if [ $? -ne 0 ]; then
    echo "Failed to map resources."
    exit 1
fi

# 3. Generate the workflow script
echo "=== 2. Generating Workflow Script ==="
python3 fix_master_workflow.py
if [ $? -ne 0 ]; then
    echo "Failed to generate workflow script."
    exit 1
fi

# Find the latest generated script (assuming it's named tgc_master_workflow*.sh)
SCRIPT_NAME=$(ls -t tgc_master_workflow*.sh | head -n 1)

if [ -z "$SCRIPT_NAME" ]; then
    echo "Failed to find generated workflow script."
    exit 1
fi

echo "=== 3. Executing Workflow Script: $SCRIPT_NAME ==="
chmod +x "$SCRIPT_NAME"

# Prompt before running the massive suite
echo "Generated $SCRIPT_NAME based on $RESOURCE_FILE."
echo "Press ENTER to start the build and test pipeline (or Ctrl+C to abort)..."
read -r

echo "Starting $SCRIPT_NAME in background..."
./"$SCRIPT_NAME" > tgc_batch_run.log 2>&1 &
PID=$!
echo "Pipeline running with PID $PID. Logs are going to tgc_batch_run.log"
echo "You can monitor progress with: tail -f tgc_batch_run.log"
