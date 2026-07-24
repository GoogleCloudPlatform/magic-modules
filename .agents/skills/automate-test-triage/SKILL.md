---
name: automate-test-triage
description: "Skill to automatically triage failing tests by reading test-metadata status files, matching GitHub issues, and generating a multi-section failure report."
---

# `automate-test-triage`

> **Note to AI Agents:** Use this skill to automate triaging failing tests from GCS `test-metadata` JSON files over a multi-day window for monitoring purposes. Do NOT modify source files or create code fixes.

## Prerequisites
* Access to `gs://nightly-test-data/test-metadata/` via `gcloud storage cat`.
* `gh` CLI installed for querying GitHub issues.

## Execution Steps

### 1. Execute Triage Script
Run the automated triage script from the workspace root:
*   **Script**: `.agents/skills/automate-test-triage/scripts/triage.py`
*   **Command**: `python3 .agents/skills/automate-test-triage/scripts/triage.py`
*   **Output File**: `tmp/test-status/persistent_failures.md`

### 2. Triage & Aggregation Rules
The script performs the following steps:
1. **Fetch Metadata**: Reads `gs://nightly-test-data/test-metadata/{version}/{date}-{version}.json` for GA and Beta providers across the past 7 days using `gcloud storage cat`.
2. **Extract Error Output**: Reads `error_message` and `log_link` directly from each test entry in `test-metadata`.
3. **High-Impact Classification (Section 1)**: Identifies failures in today's latest run flagged by:
   * 🚨 **Critical Severity**: Provider panic/crash (`panic:`, `runtime error:`, `SIGSEGV`), regardless of test count.
   * ⚠️ **High Volume**: Error signatures affecting $\ge 3$ tests in the latest run.
4. **Multi-Day Persistence Filtering (Section 2 & 3)**: Identifies tests that failed in the latest run AND failed at least 4 out of the past 7 days (filtering out one-off flaky runs and generic non-actionable errors).
5. **Error Sanitization & Dedup**: Replaces dynamic tokens (project IDs, project/org numbers, timestamps, UUIDs, SA emails) to compare error signatures across Beta and GA runs.
6. **GitHub Issue Matching**: Queries open `test-failure` issues from `hashicorp/terraform-provider-google` via `gh issue list` and matches test names against issue titles.
7. **Report Generation**: Outputs a 3-section Markdown report:
   * **Section 1**: High-Impact Errors in Latest Run (Panic/Crash Severity & High Volume).
   * **Section 2**: Persistent Failures Grouped by Error Signature (Past 7 Days).
   * **Section 3**: Detailed Per-Test Persistent Failures Table.
