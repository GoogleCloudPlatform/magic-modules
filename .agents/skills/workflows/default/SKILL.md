---
name: default-workflow
description: "Fallback workflow for general implementation and debugging tasks that do not involve creating a new resource."
---

# `default-workflow`

This document outlines the structured 5-step lifecycle for formal implementation and debugging tasks in Magic Modules.

## Execution Steps

### 1. Repo Sync
*   Execute the `repo-sync` skill (located in `.agents/skills/operations/repo-sync/`). This skill handles checking the sync status and prompting for action if needed.

### 2. Triage
*   Gather context on the change or bug. Plan the change (New feature or bug fix) within schema or logic. 
*   Consult `.agents/knowledge/index.md` for the topics the change touches and open the relevant sources.
*   Execute the `triage` skill (located in `.agents/skills/operations/triage/`) to perform this work.
*   **Transfers to Step 3:** Approved implementation plan and file paths file.

### 3. Generate
*   Once the first pass at the change is made, execute the code generation using the `generate-provider` skill (located in `.agents/skills/operations/generate-provider/`) to compile machine code in the downstream.
*   **Transfers to Step 4:** Compiled provider binary.

### 4. Test and Debug
*   Invoke the specialized `qa-test-runner` subagent using the `invoke_subagent` tool to reproduce failures and interpret logs. The subagent evaluates if the test fails/passes and returns a human-readable interpretation of the results.
*   **Transfers to Step 5:** Human-readable Markdown report explaining whether the test succeeded or failed, and what discrepancy was found.

### 5. Fix
*   This is a Remediation Planning step (similar to Triage). Take the results from `qa-test-runner`, classify symptoms against `.agents/skills/utils/test-failure-decision-tree/SKILL.md` (all catalog scenarios), and compare against reference guides or user suggestions. Propose a specific fix code change to the user (or invoke `test-fixer` subagent). 
*   Execute the `fix` skill (located in `.agents/skills/operations/fix/`) to perform this planning.

---

## The Loop
Repeat steps 2-5 as needed during the session until the primary task is complete. Reset to Step 3 (Generate) after applying an approved fix to compile it!
