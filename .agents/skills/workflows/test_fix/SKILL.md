---
name: test-fix-workflow
description: "Workflow for diagnosing, fixing, and verifying failing Terraform acceptance tests from GitHub issue URLs (detecting test-failure labels), direct prompts, or log files using Failure Scenario Decision Trees."
---

# `test-fix-workflow`

This document guides the agent through the structured 3-step lifecycle for resolving Terraform acceptance test failures in Magic Modules using scenario classification.

Consult `.agents/knowledge/index.md` for the topics the failure touches and open the relevant sources.

---

## Essential Reference Documentation

When diagnosing and fixing test failures, consult these core documentation guides in `docs/content/` (each contains tabbed examples for **MMv1 Generated** vs **Handwritten** code):
* `docs/content/develop/diffs.md` — Authoritative guide for fixing plan diffs, permadiffs, API server defaults, and `ignore_read` fields (Scenario 2).
* `docs/content/test/test.md` — Authoritative guide for acceptance test conventions, dynamic randomization, dependency bootstrapping, and VCR skipping rules (Scenario 5).
* `docs/content/test/run-tests.md` — Authoritative guide for compiling downstream providers and running `make testacc` verification loops.

---

## Execution Steps

### 1. Failure Information Intake
* Execute the `intake-test-failure` skill (`.agents/skills/utils/intake-test-failure/SKILL.md`) on the input provided by the user (GitHub issue URL, direct prompt text, GCS log link, or local log file).
* Check for GitHub issue labels such as `test-failure`, `test-failure-100`, `test-failure-50`, or any `test-failure-*` label (e.g., `test-failure-0`, `test-failure-10`) to confirm test failure classification.
* Inspect issue failure rates to determine `target_provider` (`ga`, `beta`, or `both`).
* Produce the **Normalized Failure Payload**:
  ```yaml
  normalized_failure_payload:
    test_name: "<ExactTestFunctionName>"
    target_provider: "ga" # "ga", "beta", or "both"
    error_message: |
      <Full error output, go test backtrace, and stdout plan diff for GA and/or Beta>
    parsed_logs_dir: "debug_output/<test_name>/" # (Optional)
  ```

---

### 2. Failure Scenario Classification & Remediation (Choose Path)

Match failure symptoms against the central decision tree catalog in `.agents/skills/utils/test-failure-decision-tree/SKILL.md` (all catalog scenarios).

Consult `.agents/skills/utils/test-failure-decision-tree/SKILL.md` for full symptom patterns, root causes, and remediation recipes.

#### Path A: Automated Subagent (Mandatory Default)
* **Action:** Invoke the `test-fixer` subagent (`.agents/agents/test-fixer/`) using the `invoke_subagent` tool.
* **Prompt:** Pass the **Normalized Failure Payload** to `test-fixer`.
* **Wait:** The subagent will classify the failure scenario, consult `.agents/knowledge/index.md` for relevant design rules, perform remediation (including automatically running `gcloud services enable` for shared CI projects or editing `magic-modules` test configs for test-created secondary projects), run `make provider VERSION=<ga|beta>` and `make build`, and execute target acceptance tests for `ga`, `beta`, or `both` to verify `PASS`.
* **Handoff:**
  - If `test-fixer` reports success, present the fix summary to the user.
  - If `test-fixer` reports unresolved issues, switch to **Path B (Interactive Debugging)** and consult `.agents/skills/utils/test-failure-decision-tree/SKILL.md`.

#### Path B: Interactive Debugging (Fallback)
* Use the `qa-test-runner` subagent to isolate logs and inspect request/response JSONs.
* Match symptoms against the decision tree catalog in `.agents/skills/utils/test-failure-decision-tree/SKILL.md` (all catalog scenarios) and apply source modifications using `triage` (`.agents/skills/operations/triage/`) and `fix` (`.agents/skills/operations/fix/`) skills.
* Execute provider generation (`make provider VERSION=<ga|beta>`) and re-run acceptance tests for each target provider version (`ga`, `beta`, or `both`) to verify `PASS`.

---

### 3. Verification & PR Handoff
* Verify that test output reports `PASS` for all failing targets (`ga`, `beta`, or `both`).
* **Verify Remediation Scope**: Verify that code modifications are strictly scoped to the specific field(s) or resource(s) proven by the test failure log or plan diff to be causing the failure. Do NOT include assumption-based edits to adjacent/sibling fields without empirical evidence.
* **Verify Breaking Change Compliance & User Notification**: Consult all files in `docs/content/breaking-changes/` (`breaking-changes.md` and `make-a-breaking-change.md`) before making any schema or behavioral modifications. If a breaking change needs to be made, explicitly state it to the user and list out the reason why, referencing the applicable policy in `docs/content/breaking-changes/`.
* Verify that no test-dodging flags (`ignore_read`, `default_from_api`, etc.) were introduced without justification.
* If requested by the user (e.g. "create a PR for the fix"), automatically invoke the `create-pr` skill (`.agents/skills/operations/create-pr/`) to open a pull request.
