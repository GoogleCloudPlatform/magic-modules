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
- **Determine Target Provider Version (`ga`, `beta`, or `both`)**:
  - Inspect the issue `Failure rates` section (e.g. `ga: 100%`, `beta: 0%` -> `target_provider: ga`; `ga: 0%`, `beta: 100%` -> `target_provider: beta`; `ga: 100%`, `beta: 100%` -> `target_provider: both`).
  - Match provider-specific GCS error message links in the issue body:
    - `ga error message` (e.g. `.../test-errors/ga/.../*.txt`)
    - `beta error message` (e.g. `.../test-errors/beta/.../*.txt`)
  - **Fetch the complete content of each failing provider's error text file using `gcloud storage cat`** to populate `error_message`.
- Distinguish between **Error Message Links** and **Debug Log Links** in the issue body:
  - **Error Message Links**: Contain the exact `go test` output, backtraces, and `stdout` plan diffs for GA and/or Beta runs.
  - **Debug Log Links**: Contain the full `TF_LOG=DEBUG` provider trace (`ga debug log` / `beta debug log`). Fetch and process via `tf_debug_parser.py` for `parsed_logs_dir`.
- Search the issue body or fetched error log file for:
  - Impacted acceptance test name (e.g., `TestAcc<Resource>_<Scenario>`).
  - Full error text, backtrace, and `stdout` plan diff.

#### Path B: Direct Prompt / Text Entry
- Extract `test_name`, `target_provider` (`ga`, `beta`, or `both`), and full `error_message` directly from user input.
- Extract any GCS or local log paths provided.

#### Path C: Remote / GCS Log URLs
- For Error Log or Debug Log links (`https://storage.cloud.google.com/<bucket>/<path>` or `https://storage.googleapis.com/<bucket>/<path>`), convert the URL to `gs://<bucket>/<path>` and fetch using `gcloud storage cat`:
  ```bash
  gcloud storage cat gs://<bucket>/<path>
  ```
- Store the error message text in full in `error_message`.
- Save debug logs locally to `<workspace_root>/debug_output/raw_test.log` for parsing.

---

### Step 2: Log Parsing (If Debug Log Available)
If a debug log (local or GCS) is present:
- Run the `parse-debug-logs` skill:
  ```bash
  python3 .agents/scripts/tf_debug_parser.py <path_to_log> --extract-dir debug_output/<test_name>
  ```
- Inspect `debug_output/<test_name>/outline.txt` to capture relevant API HTTP request/response JSON file paths.

---

### Step 3: Produce Complete Normalized Failure Payload

Assemble and present the normalized payload in the following format. **CRITICAL: Include `target_provider` (`ga`, `beta`, or `both`) and do NOT truncate multi-line error output, assertion backtraces, or `stdout` plan diffs. Use a multi-line YAML block scalar (`|`) to preserve full error context:**

```yaml
normalized_failure_payload:
  test_name: "<ExactTestFunctionName>"
  target_provider: "ga"  # "ga", "beta", or "both"
  error_message: |
    <Full error output, go test backtrace, and stdout plan diff for GA and/or Beta>
  parsed_logs_dir: "debug_output/<test_name>/"  # Optional
```

---

## Next Step & Handoff

Pass the **Normalized Failure Payload** to the `test-fixer` subagent (`.agents/agents/test-fixer/`) to initiate automated diagnosis and remediation.
