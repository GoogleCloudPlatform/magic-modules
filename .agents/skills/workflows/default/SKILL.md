---
name: default-workflow
description: "Fallback workflow for general implementation and debugging tasks that do not involve creating a new resource."
---

# `default-workflow`

This document outlines the structured 3-step lifecycle for formal implementation and debugging tasks in Magic Modules.

## Execution Steps

### 1. Triage
*   Gather context on the change or bug. Plan the change (New feature or bug fix) within schema or logic. 
*   Consult `.agents/knowledge/index.md` for the topics the change touches and open the relevant sources.
*   Execute the `triage` skill (located in `.agents/skills/operations/triage/`) to perform this work.
*   **Transfers to Step 2:** Approved implementation plan and file paths file.

### 2. Test and Debug
*   Invoke the specialized `qa-verification` subagent using the `invoke_subagent` tool to run verification and interpret logs. The subagent evaluates if the checks fail/pass and returns a human-readable interpretation of the results.
*   **Transfers to Step 3:** Human-readable Markdown report explaining whether the verification succeeded or failed, and what discrepancy was found.

### 3. Fix
*   This is a Remediation Planning step (similar to Triage). Take the results from `qa-verification` and compare against reference guides or user suggestions. Propose a specific fix code change to the user. 
*   Execute the `fix` skill (located in `.agents/skills/operations/fix/`) to perform this planning.

---

## The Loop
Repeat steps 1-3 as needed during the session until the primary task is complete. Reset to Step 2 (Test and Debug) after applying an approved fix!
