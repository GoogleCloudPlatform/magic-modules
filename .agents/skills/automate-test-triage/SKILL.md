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
*   **Output File**: `tmp/test-status/test-report-<date>.md` (e.g., `tmp/test-status/test-report-2026-07-28.md`)

### 2. Triage & Aggregation Rules
The script performs the following steps:
1. **Fetch Metadata**: Reads `gs://nightly-test-data/test-metadata/{version}/{date}-{version}.json` for GA and Beta providers across the past 7 days using `gcloud storage cat`.
2. **Extract Error Output**: Reads `error_message` and `log_link` directly from each test entry in `test-metadata`.
3. **High-Impact Classification (Section 1)**: Identifies actionable failures in today's latest run flagged by:
   * 🚨 **Critical Severity**: Provider panic/crash (`panic:`, `runtime error:`, `SIGSEGV`) or API enablement errors in CI test environment projects, regardless of test count.
   * ⚠️ **High Volume**: Error signatures affecting $\ge 3$ tests in the latest run.
4. **Human-Action Separation (Section 2)**: Identifies non-actionable failures requiring human intervention (e.g., Quota / Rate Limit / Stockout, Internal Error / Error Code 13, Tenant Project Creation failure).
5. **Multi-Day Persistence Filtering (Section 3)**: Identifies actionable tests that failed in the latest run AND failed at least 4 out of the past 7 days (filtering out one-off flaky runs and non-actionable errors).
6. **Error Sanitization & Dedup**: Replaces dynamic tokens (project IDs, project/org numbers, timestamps, UUIDs, SA emails) to compare error signatures across Beta and GA runs.
7. **GitHub Issue Matching**: Queries open `test-failure` issues from `hashicorp/terraform-provider-google` via `gh issue list` and matches test names against issue titles.
8. **Report Generation**: Outputs a 4-section Markdown report with a Table of Contents where every table includes clickable GCS debug log links (`[Log](url)`), error summaries, and cross-links to Section 4. **Sections 1, 2, 3, and 4 are open by default (`<details open>`)**:
   * **Section 1 (`<details open>`)**: High-Impact Actionable Errors in Latest Run (Critical Severity & High Volume; includes note indicating 500-char truncation and linking to Section 4).
   * **Section 2 (`<details open>`)**: Test Failures Requiring Human Action (Non-Actionable by Agents; includes note indicating 500-char truncation and linking to Section 4).
   * **Section 3 (`<details open>`)**: Persistent Actionable Failures Grouped by Error Signature (Past 7 Days; includes note indicating 500-char truncation and linking to Section 4).
   * **Section 4 (`<details open>`)**: Detailed Test Failures Grouped by Service Package (sorted alphabetically A–Z; the **first service package is open by default (`<details open>`)** and remaining service packages are collapsed by default (`<details>`); displays **full untruncated error messages**, all latest-run test failures per package, and a `Human Action Required?` column).
