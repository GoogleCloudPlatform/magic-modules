---
name: test-monitor-workflow
description: "Workflow for fetching, triaging, analyzing, and reporting on nightly acceptance test results across Beta and GA Google Cloud Terraform providers."
---

# `test-monitor-workflow`

This document guides the agent through the task of monitoring nightly acceptance test runs.

## Overview

The `test-monitor` workflow provides a non-destructive monitoring lifecycle that tracks nightly acceptance test execution, identifies persistent failures across a 7-day window, correlates failing tests with open GitHub issues, retrieves debug logs, and outputs a comprehensive Markdown monitoring report.

## Available Skills Used in This Workflow

*   **`fetch-nightly-results`** (`.agents/skills/fetch-nightly-results/SKILL.md`): Downloads nightly test status JSON files for Beta and GA providers from GCS (`gs://nightly-test-data/test-metadata/`).
*   **`automate-test-triage`** (`.agents/skills/automate-test-triage/SKILL.md`): Triages 7-day failure trends, matches open GitHub `test-failure` issues, sanitizes dynamic error tokens for provider comparison, and outputs reports.

## Execution Steps

### 1. Fetch Test Results
*   Execute the `fetch-nightly-results` skill to retrieve test status JSON files for both GA and Beta providers from Google Cloud Storage (`gs://nightly-test-data/test-metadata/`) for the latest available run and past 7 days using `gcloud storage cat` (if GCS authentication or permission fails, verify `gcloud auth login` and `roles/storage.objectViewer` access).

### 2. Triage & Aggregate Failure Trends
*   Execute the `automate-test-triage` skill by running `python3 .agents/skills/automate-test-triage/scripts/triage.py` from the root workspace (or pass `--date YYYY-MM-DD` to specify a target end date for the 7-day window).
*   Group the target date's latest run failures by error signature and flag High-Impact errors based on:
    1. 🚨 **Critical Severity**: Provider panic/crash (`panic:`, `runtime error:`, `SIGSEGV`) or API enablement errors in CI test environment projects, prioritized at the top of Section 1 regardless of test count.
    2. ⚠️ **High Volume**: Error signatures affecting $\ge 3$ tests in the latest run.
*   Separate non-actionable errors requiring human intervention (Quota / Rate Limit / Stockout, Internal Error / Error Code 13, Tenant Project Creation) into Section 2.
*   Filter for persistent actionable test failures (failing in latest run AND $\ge 4$ out of the past 7 days).

### 3. Cross-Reference GitHub Issues
*   Fetch open issues labelled `test-failure` from `hashicorp/terraform-provider-google` using `gh issue list`.
*   Match test names against issue titles to correlate persistent failures with existing open tracking tickets.

### 4. Categorize Root Causes & Analyze Errors
*   Extract `error_message` and `log_link` directly from `test-metadata` entries.
*   Analyze API error payloads and debug logs to categorize the failure domain (e.g., Provider Panic, Quota Exceeded, API Permission / IAM, Model Availability, State Mismatch, Flakiness).

### 5. Generate Monitoring Report
*   Save the final Markdown report to `tmp/test-status/test-report-<date>.md` (e.g., `tmp/test-status/test-report-2026-07-28.md`).
*   Present the executive summary and 4 report sections (where every section table includes clickable GCS debug log links `[Log](url)`, expanded 500-character error summaries, and links to Section 4) to the user: Section 1 (High-Impact Actionable Errors: Panics & API Enablement in Test Env, High Volume), Section 2 (Test Failures Requiring Human Action), Section 3 (Persistent Actionable Failures Grouped by Error Signature), and Section 4 (Detailed Test Failures Grouped by Service Package, showing all latest-run failures with failed out of total tests, percentages, and a Human Action Required? column).

---

## Boundaries & Guardrails
*   **Monitoring Only:** Do NOT modify source code or template files in `magic-modules`.
*   **No Code Fixes / PRs:** Do NOT create fix branches or attempt to open pull requests.
