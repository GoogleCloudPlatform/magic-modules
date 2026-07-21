---
name: test-fix-workflow
description: "Workflow for diagnosing, fixing, and verifying failing Terraform acceptance tests from GitHub issue URLs (detecting test-failure labels), direct prompts, or log files using Failure Scenario Decision Trees."
---

# `test-fix-workflow`

This document guides the agent through the structured 3-step lifecycle for resolving Terraform acceptance test failures in Magic Modules using scenario classification.

Consult `.agents/knowledge/index.md` for the topics the failure touches and open the relevant sources.

---

## Execution Steps

### 1. Failure Information Intake
* Execute the `intake-test-failure` skill (`.agents/skills/utils/intake-test-failure/SKILL.md`) on the input provided by the user (GitHub issue URL, direct prompt text, GCS log link, or local log file).
* Check for GitHub issue labels such as `test-failure`, `test-failure-100`, `test-failure-50`, or any `test-failure-*` label to confirm test failure classification.
* Produce the **Normalized Failure Payload**:
  ```yaml
  normalized_failure_payload:
    test_name: "<ExactTestFunctionName>"
    error_message: "<Concise error string or failed assertion>"
    parsed_logs_dir: "debug_output/<test_name>/" # (Optional)
  ```

---

### 2. Failure Scenario Classification & Remediation (Choose Path)

Match failure symptoms against standard decision tree scenarios:

| Scenario ID | Symptom / Error Pattern | Root Cause Category | Primary File Location & Remedy |
| :--- | :--- | :--- | :--- |
| **Scenario 1** | State drift, `ImportStateVerify` mismatch, or non-empty plan after Apply | State normalization / Default mismatch | Edit `mmv1/products/<product>/<Resource>.yaml` to add `diff_suppress_func`, adjust state reader, or handle API defaults |
| **Scenario 2** | HTTP 400 `InvalidArgument`, `Unknown Field`, or field serialization error | API Request Payload schema mismatch | Compare `01_POST_request.json` against API schema; fix field camelCase/snake_case mapping or `send_empty_value` |
| **Scenario 3** | HTTP 404 immediately post-Create or concurrent operation conflict | LRO / Eventual Consistency timing | Configure `autogen_async` in YAML or add polling waiter logic in handwritten code overrides |
| **Scenario 4** | Pre-requisite resource failure, name collision, or test setup error | Sample HCL test template bug | Update `..._test.go.tmpl` to use dynamic random string suffixes or fix test dependencies |
| **Scenario 5** | gRPC Code 13, HTTP 500/502/503, `An internal error has occurred`, or API backend crash | Non-Remediable GCP Backend Error | **Early Exit & Handoff:** Do NOT edit `magic-modules`. Report failure as a GCP service-side internal error requiring API team investigation. |

#### Path A: Automated Subagent (Mandatory Default)
* **Action:** Invoke the `test-fixer` subagent (`.agents/agents/test-fixer/`) using the `invoke_subagent` tool.
* **Prompt:** Pass the **Normalized Failure Payload** to `test-fixer`.
* **Wait:** The subagent will classify the failure scenario, consult `.agents/knowledge/index.md` for relevant design rules, edit `magic-modules` source files, run `make provider` and `make build`, and execute acceptance tests to verify `PASS`.
* **Handoff:**
  - If `test-fixer` reports success, present the fix summary to the user.
  - If `test-fixer` reports unresolved issues, switch to **Path B (Interactive Debugging)**.

#### Path B: Interactive Debugging (Fallback)
* Use the `qa-test-runner` subagent to isolate logs and inspect request/response JSONs.
* Match symptoms against Scenario 1-4 decision trees and apply source modifications using `triage` and `fix` skills.
* Execute `generate-provider` (`make provider`) and re-run acceptance tests manually to verify.

---

### 3. Verification & PR Handoff
* Verify that the test output reports `PASS`.
* Verify that no test-dodging flags (`ignore_read`, `default_from_api`, etc.) were introduced without justification.
* If requested by the user (e.g. "create a PR for the fix"), automatically invoke the `create-pr` skill (`.agents/skills/operations/create-pr/`) to open a pull request.
