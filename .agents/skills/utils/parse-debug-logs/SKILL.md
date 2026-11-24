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
Execute the `tf_debug_parser.py` script against the debug logs using the mode that fits your task:

#### Parse Debug Logs
```bash
# Print a concise human-readable summary (timeline, endpoint summary, test results, destroy errors, failures, latest diff)
python3 <PATH_TO_MAGIC_MODULES>/.agents/scripts/tf_debug_parser.py /path/to/test_output.log --summary

# Output the full structured timeline and summaries as JSON
python3 <PATH_TO_MAGIC_MODULES>/.agents/scripts/tf_debug_parser.py /path/to/test_output.log --json

# Extract the compressed outline and individual JSON request/response payloads to a directory
python3 <PATH_TO_MAGIC_MODULES>/.agents/scripts/tf_debug_parser.py /path/to/test_output.log --extract-dir <PATH_TO_MAGIC_MODULES>/debug_output
```

### 3. Review & Diagnostic Summary Fields (`--summary` and `--json`)
Both `--summary` and `--json` include three review-focused diagnostic sections designed to surface common test and API defects quickly:

* **`endpoint_summary`** (`====== ENDPOINT SUMMARY ======`):
  * Groups all HTTP requests by `<host>/<api_version>` (e.g., `ces.googleapis.com/v1beta: 14`, `compute.googleapis.com/compute/v1: 6`, `storage.googleapis.com/storage/v1: 3`), parsed from `Host:` headers (falling back to full request URLs when present) with query parameters stripped.
  * In `--json`, emitted as a deterministically sorted object mapping `"<host>/<api_version>"` to request counts.
  * **Reviewer check:** Verify GA changes actually exercise the GA endpoint (e.g., `/v1` or `/compute/v1`) rather than hitting only `/v1beta` / `/compute/beta`.
* **`test_results`** (`====== TEST RESULTS ======`):
  * Parses Go test result lines (`--- PASS:`, `--- FAIL:`, `--- SKIP:`), including indented subtests (e.g., `    --- SKIP: TestFoo/subtest (0.00s)`).
  * In `--json`, emitted as `{"counts": {"pass": N, "fail": N, "skip": N}, "tests": [{"name": ..., "status": "PASS"|"FAIL"|"SKIP", "duration": ..., "line": ...}]}` while preserving existing `failures` entries for backward compatibility.
  * **Reviewer check:** Confirm tests did not silently skip (`counts.skip > 0`) or fail in a subtest.
* **`destroy_errors`** (`====== DESTROY ERRORS ======`):
  * Captures HTTP `>= 400` responses occurring during the `destroy` / `CheckDestroy` phase while excluding expected `404` read-after-delete polling responses and non-destroy step errors.
  * In `--json`, emitted as a list of objects containing `line`, `status`, `status_code`, `phase`, and when available the paired request `method`, `host`, `url`, and extracted `error_message`.
  * **Reviewer check:** Catch broken `CheckDestroy` logic or cleanup failures (e.g., `400`, `403`, `500` during destroy) even if a test swallowed the error and reported `PASS`.

### 4. Verification & Handoff
* For quick PR/test verification, inspect `--summary` or `--json` (`endpoint_summary`, `test_results`, `destroy_errors`, `failures`).
* For deep payload debugging with `--extract-dir`, use `view_file` to read `outline.txt` inside the generated `<test_name>_<timestamp>/` directory, then inspect any specific request/response JSON file referenced in `outline.txt` (e.g., `01_REQUEST_POST.json`).
* Pivot to `.agents/skills/operations/troubleshooting_reference.md` to map failures to a known fix strategy before invoking the `fix` workflow.
