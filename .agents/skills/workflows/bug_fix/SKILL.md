---
name: bug-fix-workflow
description: "Workflow for triaging, planning, fixing, and verifying reported provider bugs."
---

# `bug-fix-workflow`

This document outlines the structured lifecycle for investigating, reproducing, fixing, and verifying reported provider bugs in Magic Modules.

## Execution Steps

### 1. Triage & Context Gathering
* **External context:** Read the target issue description, related bug reports, and external API documentation (e.g., REST API references).
* **Internal context:** Consult the Knowledge Index (`.agents/knowledge/index.md`) for relevant topics/patterns and search the codebase for affected schemas, fields, expanders, or flatteners.
* **Historical context:** Trace Git history (`git log`, PRs, blame) in the repository and downstream providers to identify how the defect was introduced or how similar resources behave.

### 2. Empirical Issue Reproduction & Remediation Plan
* **Empirical Reproduction (RED Check):** Formulate a minimal reproduction (acceptance test, CLI plan with dev overrides, or targeted schema test) and run it against the *unfixed baseline*.
  * ❌ **If reproduction FAILS to recreate the bug** (or unexpectedly passes): STOP. The static analysis is incomplete. Re-investigate, update findings, and present a revised plan at the HITL Checkpoint.
  * ✅ **If reproduction SUCCEEDS in recreating the bug**: Document the failure signature.
* **Remediation Plan:**
  * Identify root cause and propose the exact MMv1 changes.
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
* **Verify Reproduction:** Run the **exact same reproduction** from Step 2 against the patched build to confirm it now passes (RED $\rightarrow$ GREEN).
* **Cleanup Ephemeral Tests:** Remove any ad-hoc reproduction tests before proceeding.
* **Acceptance Tests:** Execute `qa-test-runner` (`.agents/skills/operations/qa-test-runner/`) for target acceptance tests. Verify payloads and `PASS` status.

### 5. Resolution & Issue Reporting
* **Plan Completeness:** Verify all changed and generated files are staged.
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
