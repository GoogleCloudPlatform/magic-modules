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
*   **External context:** Read the target issue description, related bug reports, and external API documentation (e.g., REST API references) to understand service behavior and parameters.
*   **Internal context:** Consult the Knowledge Index (`.agents/knowledge/index.md`) for any relevant topics, patterns, or repository-specific instructions. Then search the codebase to locate where the affected fields, schemas, expanders, or flatteners are defined.
*   **Historical context:** Trace Git logs, tags, and past PRs in the repository and downstream provider repositories to identify the lifecycle of the affected code, related fixes, or similar resource implementations (`Modeled after:` / `Based on:`).

### 3. Remediation Planning (Proposal)
*   Analyze the triage findings and identify the root cause.
*   **Backwards Compatibility Check:** When proposing code changes that alter schema validation (e.g., adding `ValidateFunc`), default values, or resource naming/ID formulas, check `docs/content/breaking-changes/` and explicitly note in the investigation report whether existing deployed resources or valid configurations are affected.
*   Create a final investigation report (an artifact) detailing the root cause, affected version matrix, and analysis.
*   Propose a resolution path to the user:
    *   If the bug can be resolved or closed without code changes, propose closing the issue (e.g. if the bug is already fixed in a released version, is a duplicate, is not a bug, or is otherwise closeable).
    *   If the bug requires code changes, propose the code modification plan (dipping into the `fix` planning skill located in `.agents/skills/operations/fix/`).
*   **Verification Proposal:** If running an acceptance test (live or VCR) is necessary to verify either the code change or the initial analysis (e.g., to prove it's already resolved in the latest version), explicitly include the test execution in the proposed plan. Do not run tests unilaterally unless approved.
    *   **Version Upgrade Safety (`RELEASE_DIFF=true`):** When testing bug fixes that modify schema validation, defaults, or ID/name formatting, recommend running the acceptance test with `RELEASE_DIFF=true make testacc ...` to verify zero plan diffs against the released provider baseline.
*   **HIL steering checkpoint:** Present the analysis, proposal, and testing plan to the user. **Do not proceed to implementation or testing until the user confirms the plan.**


### 4. Implementation & Code Generation (Only if code changes are required)
*   Apply the approved schema or logic changes in Magic Modules (`mmv1/`).
*   **Template Modifications:** If the fix requires modifying engine templates (`mmv1/templates/terraform/`), consult the Knowledge Index entry on [Template Modifications & Blast Radius](../../../knowledge/template/template-modifications.md). Obtain explicit user approval before modifying engine templates.
*   Execute code generation to compile the downstream provider (using the `generate-provider` skill located in `.agents/skills/operations/generate-provider/`).


### 5. Verification Testing (Only if approved in the plan)
*   Execute the `qa-test-runner` skill (located in `.agents/skills/operations/qa-test-runner/`) to trigger test runs and parse logs.
*   **Verification:** Verify the test finishes successfully (`PASS`). Inspect the API request/response payloads in the test logs or report to confirm the correct payload structure (e.g., verifying fields are correctly set or sent as `null`/omitted).


### 6. Resolution & Issue Reporting
*   **Plan Completeness:** Verify that every file listed in the remediation plan (including any necessary documentation) has been generated and staged.
*   **Pre-PR Quality & Verification Gate:** Before opening a PR or finalizing the branch, run the following verification pipeline:
    1. **Build Verification:** Run `make build` in downstream provider repository to ensure full compilation passes without syntax errors.
    2. **Acceptance Test Verification:** Confirm target acceptance tests pass (`PASS`).
    3. **Go Formatting:** Run `gofmt -s -w` on all modified or newly created `.go` files under `mmv1/third_party/terraform/`.
    4. **Pre-Gen Static Checks:** Run `./.agents/skills/utils/run-pre-gen-checks/scripts/run_pre_gen_checks.sh` to ensure Go formatting, YAML linting, template validation, and MMv1 unit tests pass.
    5. **Breaking Change Validation:** Run `./.agents/skills/utils/validate-provider-changes/scripts/validate_provider_changes.sh` if schemas or properties were modified.
*   **Workspace Cleanup:** Run `git status --porcelain` and remove any untracked `.log`, `.test`, or temporary test artifacts across both repositories before reporting resolution or creating a PR.
*   **Artifact Report:** If code changes or verification tests were performed, compile these results into a separate verification/test report artifact.
*   **PR Creation:** When opening a PR, execute the `create-pr` skill (`.agents/skills/operations/create-pr/`), which governs branch creation, PR title length, release notes, and reference linking (`Modeled after:` / `Based on:`).
*   **GitHub Response Draft:** Draft a final, succinct public response containing verified PR/commit links.
    *   **Succinct Public Communication:** Responses should be concise (2–3 sentences preferred): state what changed, why, and refer readers to the PR or documentation for technical deep-dives.
*   **HIL steering checkpoint:** Present the final response draft and any new verification reports to the user for sign-off and issue closure.


---

## The Loop
If verification fails during Step 5, repeat steps 3-5 as needed.
*   **Scope Expansion Guardrail:** If debugging reveals that resolving the root cause requires expanding scope beyond the approved plan (such as modifying engine templates in `mmv1/templates/` or altering additional fields/resources), do NOT apply changes silently. Loop back to Step 3, update the investigation report artifact, and obtain explicit user approval.
*   Reset to Step 4 (Implementation & Code Generation) after applying any approved fix changes to compile and re-test.
