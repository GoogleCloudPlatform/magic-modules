---
name: intake-test-failure
description: "Process raw failure info (GitHub issue URL, direct prompt text, GCS debug log links, or local logs) and normalize it into a clean failure payload for test fixing."
---

# `intake-test-failure`

This skill converts raw, unstructured, or varied failure reports into a standardized **Normalized Failure Payload** ready for consumption by the `test-fixer` subagent or developer remediation workflow.

## Input Formats Handled

1. **GitHub Issue URL** (e.g., `https://github.com/hashicorp/terraform-provider-google/issues/28244`)
2. **Direct Text Prompt** (Test name, error message, and optional log snippet/URL)
3. **GCS / Remote Log URL** (e.g., `https://storage.googleapis.com/...` or `gs://...`)
4. **Local Log File** (e.g., `test_output.log` or debug log file)

---

## Execution Steps

### Step 1: Extract Core Failure Information

#### Path A: GitHub Issue URL
- Use `read_url_content` (or `gh issue view` if CLI available) to read the issue content and inspect its GitHub labels.
- Check if labels contain `test-failure`, `test-failure-100`, `test-failure-50`, or any similar `test-failure*` labels to confirm this is an acceptance test failure issue.
- Search the issue body or metadata for:
  - Impacted acceptance test name (e.g., `TestAcc<Resource>_<Scenario>`).
  - Failure error backtrace or state assertion error (e.g., `Step 1/1 error: Check failed...`).
  - GCS debug log URL links (often ending in `.log` or `.txt` hosted on Google Cloud Storage).

#### Path B: Direct Prompt / Text Entry
- Extract `test_name` and `error_message` directly from user input.
- Extract any GCS or local log paths provided.

#### Path C: Remote / GCS Debug Log URL
- For GCS log links (`https://storage.cloud.google.com/<bucket>/<path>` or `https://storage.googleapis.com/<bucket>/<path>`), convert the URL to `gs://<bucket>/<path>` and fetch the exact log output using `gcloud storage cat`:
  ```bash
  gcloud storage cat gs://<bucket>/<path>
  ```
- Save the log locally to `<workspace_root>/debug_output/raw_test.log` if parsing is required.

---

### Step 2: Log Parsing (If Debug Log Available)
If a debug log (local or GCS) is present:
- Run the `parse-debug-logs` skill:
  ```bash
  python3 .agents/scripts/tf_debug_parser.py <path_to_log> --extract-dir debug_output/<test_name>
  ```
- Inspect `debug_output/<test_name>/outline.txt` to capture relevant API HTTP request/response JSON file paths.

---

### Step 3: Produce Minimal Normalized Failure Payload

Assemble and present the normalized payload in the following clean format:

```yaml
normalized_failure_payload:
  test_name: "<ExactTestFunctionName>"
  error_message: "<Concise error string or failed assertion>"
  parsed_logs_dir: "debug_output/<test_name>/"  # Optional
```

---

## Next Step & Handoff

Pass the **Normalized Failure Payload** to the `test-fixer` subagent (`.agents/agents/test-fixer/`) to initiate automated diagnosis and remediation.
