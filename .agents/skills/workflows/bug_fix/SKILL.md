---
name: bug-fix-workflow
description: "Workflow for triaging, planning, fixing, and verifying reported provider bugs."
---

# `bug-fix-workflow`

This document outlines the structured lifecycle for investigating, reproducing, fixing, and verifying reported provider bugs in Magic Modules.

## Execution Steps

### 1. Triage & Context Gathering via `bug-triager` Subagent
Delegate initial issue intake, deep codebase research, external documentation lookup, and Git history tracing to the `bug-triager` subagent (`.agents/agents/bug-triager/`). This keeps the primary agent's context window clean for reproduction and remediation.

* **Action:** Invoke the `bug-triager` subagent using the `invoke_subagent` tool:
  ```python
  invoke_subagent(
      TypeName="bug-triager",
      Role="Bug Triage & Context Gatherer",
      Prompt="Triage reported bug: <ISSUE_URL_OR_DESCRIPTION>. Inspect issue comments and any linked Buganizer tickets (b/XXXX in description), gather external API docs, inspect Magic Modules schemas and code in mmv1/, check git history for introducing changes, consult .agents/knowledge/index.md, and return a comprehensive Triage & Context Report with root cause hypothesis and concrete reproduction strategy (recommending a new test or modifying an existing test)."
  )
  ```
* **Subagent Scope & Responsibilities:**
  * **External & Issue context:** Reads the target issue description, full issue comments thread (e.g., `gh issue view --comments`), linked Buganizer issues (`b/XXXX` or `b/<id>` links in description rendered via `/google/bin/releases/issues-cli/issues readonly render <issue_id>`), related bug reports, and external API documentation (e.g., Google Cloud REST API references).
  * **Internal context:** Consults the Knowledge Index (`.agents/knowledge/index.md`) for relevant topics/patterns, searches the codebase for affected schemas, fields, expanders, flatteners, or custom code, and inspects existing tests/samples to identify reproduction candidates.
  * **Historical context:** Traces Git history (`git log`, PRs, blame) in `magic-modules` and downstream providers to identify how the defect was introduced or how similar resources behave.
  * **Synthesis:** Formulates root cause hypothesis and recommends how to recreate the bug with a new test or by modifying an existing test.
* **Handoff:** Review the returned **Triage & Context Report**. Use the identified components, reproduction strategy, and root cause hypothesis to proceed directly to Step 2.

### 2. Empirical Issue Reproduction & Remediation Plan

#### Recreating the Defect with a Test (RED Check)
Recreate the reported defect on the unfixed baseline using either a new test or a modified existing test based on `bug-triager`'s analysis:

* **Option A: Add a New Acceptance/Regression Test (Preferred for distinct scenarios or missing coverage)**:
  * **MMv1 Generated Resources:**
    1. Create a new sample template file: `mmv1/templates/terraform/samples/services/<product>/<sample_name>.tf.tmpl` containing the minimal HCL configuration to trigger the bug.
    2. Register the sample in the resource YAML (`mmv1/products/<product>/<Resource>.yaml`) under `samples:`. Follow sample conventions from [`docs/content/test/test.md`](../../../docs/content/test/test.md) (use `resource_id_vars` for identifiers needing `tf-test` prefixes and random suffixes; use `vars` for values that vary between test steps; hardcode constants).
  * **Handwritten Resources:**
    - Add a new test function `TestAcc<Resource>_<BugScenario>` in `mmv1/third_party/terraform/services/<product>/resource_<name>_test.go` or `data_source_<name>_test.go`.
  * **Pure Go Functions (Unit Tests):**
    - If the bug is isolated to an algorithmic Go helper (`DiffSuppress`, `ValidateFunc`, state parsers) per [`.agents/knowledge/test/unit-test-scope.md`](../../../knowledge/test/unit-test-scope.md), add a focused unit test in `*_test.go`.

* **Option B: Modify an Existing Test**:
  * Identify an existing acceptance test or sample covering the affected resource.
  * Modify the test configuration (e.g., adding the problematic field, setting an edge-case value combination, or adding an update step) to exercise the reported bug path.

* **Execute Baseline Reproduction (RED Check):**
  1. Generate downstream code: `make provider VERSION=<ga|beta>` so the new or modified test is compiled into the provider repository.
  2. Run the test against the unfixed baseline:
     ```bash
     cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta && make testacc TEST=./google-beta/services/<product> TESTARGS='-run=<TestName>'
     ```
  3. **Evaluate Result:**
     * ❌ **If reproduction FAILS to recreate the bug** (test passes or fails for an unrelated reason): STOP. The test configuration does not accurately capture the reported defect. Re-examine the issue description, comments, or Buganizer notes, adjust the test configuration, and re-run until the failure is faithfully recreated.
     * ✅ **If reproduction SUCCEEDS in recreating the bug**: Document the failure signature (e.g., non-empty plan diff, unexpected API error, assertion failure).

#### Remediation Plan
* Identify root cause and propose the exact MMv1 changes to make the reproduction test pass.
* **Unit Test Scope:** Consult [`.agents/knowledge/test/unit-test-scope.md`](../../../knowledge/test/unit-test-scope.md).
* **Backwards Compatibility:** Check `docs/content/breaking-changes/` for schema or validation changes.
* **Version Upgrade Safety (`RELEASE_DIFF=true`):** Recommend running acceptance tests with `RELEASE_DIFF=true make testacc ...` to verify zero plan diffs against the released provider baseline.
* Create an Investigation & Remediation Report artifact.

---

### 🛑 CRITICAL: MANDATORY HUMAN-IN-THE-LOOP CHECKPOINT
* **STOP ALL TOOL CALLS AND YIELD EXECUTION TO THE USER.**
* Present the investigation report, empirical reproduction results, and remediation plan.
* **DO NOT** write fix code, generate downstream providers, or edit files until the user explicitly approves the plan.
* *Re-steering triggers:* If reproduction failed, verification failed, or fix scope expands, you MUST return here and pause.

---

### 3. Implementation & Code Generation (Only after user approval)
* Apply approved changes in Magic Modules (`mmv1/`).
* **Template Modifications:** If modifying engine templates (`mmv1/templates/terraform/`), consult [Template Modifications & Blast Radius](../../../knowledge/template/template-modifications.md) and obtain explicit user approval before proceeding.
* Execute code generation via `generate-provider` (`.agents/skills/operations/generate-provider/`) to compile downstream providers.

### 4. Post-Fix Verification (GREEN Check)
* **Verify Reproduction:** Run the **exact same reproduction test** (the newly added test or modified existing test) from Step 2 against the patched build to confirm it now passes (RED $\rightarrow$ GREEN).
* **Retain Regression Test:** Retain the newly added test or test modifications as a permanent regression test in the repository to prevent future regressions. Do NOT revert or delete it. (Only remove throwaway local scratch scripts that violate [`.agents/knowledge/test/unit-test-scope.md`](../../../knowledge/test/unit-test-scope.md)).
* **Acceptance Tests:** Execute `qa-test-runner` (`.agents/skills/operations/qa-test-runner/`) for target acceptance tests. Verify payloads and `PASS` status.

### 5. Resolution & Issue Reporting
* **Plan Completeness:** Verify all changed and generated files are staged, including:
  1. MMv1 source fixes (`mmv1/products/...` or handwritten files).
  2. The new regression test sample (`.tf.tmpl` and YAML registration) or modified existing test.
  3. Generated downstream provider code and tests.
* **Pre-PR Quality Gate:**
  1. **Build Verification:** Run `make build` in downstream provider repository to ensure compilation passes without errors.
  2. **Acceptance Test Verification:** Confirm target acceptance tests pass (`PASS`).
  3. **Pre-Gen Static Checks:** Run `./.agents/skills/utils/run-pre-gen-checks/scripts/run_pre_gen_checks.sh`.
  4. **Breaking Change Validation:** Run `validate-provider-changes` if schemas changed.
* **Prepare PR Link:** Execute `prepare-pr-link` (`.agents/skills/operations/prepare-pr-link/`) to push to fork and generate the comparison link.
* **Workspace Cleanup:** Run `git status --porcelain` and remove any untracked `.log`, `.test`, or temporary test artifacts.
* **HIL Final Checkpoint:** Present the final response draft (2–3 sentences) and verification report to the user for sign-off.

---

## The Loop
If verification fails during Step 4, repeat steps 2-4 as needed.
* **Scope Expansion Guardrail:** If resolving the root cause requires expanding scope beyond the approved plan (such as modifying engine templates in `mmv1/templates/` or altering additional fields/resources), do NOT apply changes silently. Loop back to Step 2, update the report, and obtain explicit user approval at the HITL Checkpoint.
* Reset to Step 3 (Implementation & Code Generation) after applying any approved fix changes to compile and re-test.
