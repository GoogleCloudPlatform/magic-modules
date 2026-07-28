---
name: bug-fix-workflow
description: "Workflow for triaging, planning, fixing, and verifying reported provider bugs."
---

# `bug-fix-workflow`

This document outlines the structured 6-step lifecycle for investigating, planning, fixing, and verifying reported provider bugs in Magic Modules.

## Prerequisites
*   You must be in the `magic-modules` root directory.
*   You must have access to the relevant downstream provider workspace(s) (GA and/or Beta) required by the bug. When in doubt, ask the user.


## Execution Steps

### 1. Repo Sync & Baseline Setup
*   Execute the `repo-sync` skill (located in `.agents/skills/operations/repo-sync/`). This skill handles checking the sync status and prompting for action if needed to establish a clean sync baseline.

### 2. Triage & Context Gathering
*   **External context:** Read the target GitHub issue description, related bug reports, and external API documentation (e.g., REST API references) to understand GCP service behavior and parameters.
*   **Internal context:** Search the codebase (grep) to locate where the affected fields, schemas, expanders, or flatteners are defined.
*   **Historical context:** Trace Git logs, tags, and PRs in downstream provider repositories to identify the lifecycle of the affected code (e.g., when it was introduced, deprecated, or modified).

### 3. Remediation Planning (Proposal)
*   Analyze the triage findings and identify the root cause.
*   Create a final investigation report (an artifact) detailing the root cause, affected version matrix, and analysis.
*   Propose a resolution path to the user:
    *   If the bug can be resolved or closed without code changes, propose closing the issue (e.g. if the bug is already fixed in a released version, is a duplicate, is not a bug, or is otherwise closeable).
    *   If the bug requires code changes, propose the code modification plan (dipping into the `fix` planning skill located in `.agents/skills/operations/fix/`).
*   **Verification Proposal:** If running an acceptance test (live or VCR) is necessary to verify either the code change or the initial analysis (e.g., to prove it's already resolved in the latest version), explicitly include the test execution in the proposed plan. Do not run tests unilaterally unless approved.
*   **HIL steering checkpoint:** Present the analysis, proposal, and testing plan to the user. **Do not proceed to implementation or testing until the user confirms the plan.**


### 4. Implementation & Code Generation (Only if code changes are required)
*   Apply the approved schema or logic changes in Magic Modules (`mmv1/`).
*   Execute code generation to compile the downstream provider (using the `generate-provider` skill located in `.agents/skills/operations/generate-provider/`).


### 5. Verification Testing (Only if approved in the plan)
*   Execute the `qa-test-runner` skill (located in `.agents/skills/operations/qa-test-runner/`) to trigger test runs and parse logs.
*   **Verification:** Verify the test finishes successfully (`PASS`). Inspect the API request/response payloads in the test logs or report to confirm the correct payload structure (e.g., verifying fields are correctly set or sent as `null`/omitted).


### 6. Resolution & Issue Reporting
*   If code changes or verification tests were performed, compile these results into a separate verification/test report artifact.
*   Draft a final, succinct GitHub response containing verified PR/commit links.
*   **HIL steering checkpoint:** Present the final response draft and any new verification reports to the user for sign-off and issue closure.


---

## The Loop
If verification fails during Step 5, repeat steps 3-5 as needed. Reset to Step 4 (Implementation & Code Generation) after applying any approved fix changes to compile and re-test.
