---
name: parse-debug-logs
description: "Parses massive Terraform debug logs into an easy-to-read outline of API delays, failures, and request payloads."
---

# `parse-debug-logs`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

## Prerequisites
* You must have a debug log file resulting from a failed `run-acctests` execution (e.g., `test_output.log`).

## Execution Steps

### 1. Verification
Verify the log file exists.

#### Verify Log Exists
```bash
# Replace /path/to/test_output.log with the actual path
ls -lh /path/to/test_output.log
```

### 2. The Core Commands
Execute the `tf_debug_parser.py` script against the debug logs.

#### Parse Debug Logs
```bash
# Execute the parser to extract the outline
python3 <PATH_TO_MAGIC_MODULES>/.agents/scripts/tf_debug_parser.py /path/to/test_output.log --extract-dir <PATH_TO_MAGIC_MODULES>/debug_output
```

### 3. Verification & Handoff
* Use the `view_file` tool to read the generated `outline.txt` INSIDE the output directory (or its subdirectory if the script creates one). It contains a compressed timeline of API requests/responses and Terraform state transitions.
* If you find an API error (HTTP 4xx or 5xx) or an unexpected default field, use `view_file` on the specific JSON payload file referenced in `outline.txt` (e.g., `01_REQUEST_POST.json`) found in that same directory.
* Pivot to the `.agents/skills/operations/troubleshooting_reference.md` document to map the failure to a known fix strategy before invoking the `fix` workflow.
