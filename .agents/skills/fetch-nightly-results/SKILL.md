---
name: fetch-nightly-results
description: "Skill to fetch latest nightly test results and metadata for Beta or GA providers from GCS test-metadata."
---

# `fetch-nightly-results`

> **Note to AI Agents:** Use this skill to retrieve test statuses, duration, error messages, and log links from `test-metadata`. Direct `gcloud storage cat` is the primary method; an optional helper script is provided if multi-day or complex filtering is required.

## Prerequisites
* `gcloud` CLI installed and authenticated with access to `gs://nightly-test-data` (requires `roles/storage.objectViewer` permission). If permission or authentication fails, verify your login via `gcloud auth login`.
* Must be executed from within the `magic-modules` workspace.

## Primary Execution Steps (Direct GCS Command)

### 1. Construct Target Metadata File URI
Construct the GCS URI based on provider version (`ga` or `beta`) and target date (`YYYY-MM-DD`).
*   **Metadata URI Format**: `gs://nightly-test-data/test-metadata/{version}/{date}-{version}.json`
*   **Example**: `gs://nightly-test-data/test-metadata/ga/2026-07-23-ga.json`
*   *Note*: If date is omitted, construct dates starting from yesterday and check backwards for the latest available run.

### 2. Fetch Content
Execute `gcloud storage cat` to read the test metadata:

```bash
gcloud storage cat gs://nightly-test-data/test-metadata/{version}/{date}-{version}.json
```
*   *Error Handling*: If `gcloud` returns a permission or authentication error (`permission denied`, `401`, `403`), check `gcloud auth login` and ensure read access (`roles/storage.objectViewer`) to `gs://nightly-test-data`.

### 3. JSON Payload Fields
Each element in `test-metadata` contains all test details:
*   `name`: Test function name (e.g., `TestAccComputeRegionInstanceTemplate_basic`).
*   `status`: Outcome (`SUCCESS`, `FAILURE`, etc.).
*   `service`: Affected GCP service.
*   `error_message`: Full Go test error message and backtrace snippet.
*   `log_link`: HTTPS link to the full execution debug log.
*   `duration`: Test duration in milliseconds.

---

## Optional Helper Script (If Script Is Needed)
If multi-day querying or complex status filtering is needed, use `.agents/skills/fetch-nightly-results/scripts/fetch_results.py`:

```bash
# Query specific test across past 7 days
python3 .agents/skills/fetch-nightly-results/scripts/fetch_results.py --test TestAccComputeRegionInstanceTemplate_basic --days 7 --status ALL
```
