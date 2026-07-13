---
name: qa-test-runner
description: "Triggers the QA Test Runner subagent to reproduce failures and interpret logs."
---

# `qa-test-runner`

Use this skill to delegate heavy test reproduction and log parsing to a dedicated Subagent. This keeps your main context window clean.

## Execution Steps

### 1. Identify Target and Context
Gather context on *what* changed (the resource, the field, the PR) and which test you need to run.

### 2. Invoke Subagent
Use the `invoke_subagent` tool. Pass the context of the change as part of the trigger prompt so the subagent knows what it's looking for!

**Trigger Example:**
invoke_subagent("qa-test-runner", "Run <TEST_NAME>. We just made the following change: <CHANGE_DESCRIPTION>. Verify if it fails or passes. Use your active workspace to find the magic-modules path <PATH_TO_MAGIC_MODULES> and output debug logs to <PATH_TO_MAGIC_MODULES>/debug_output.")

### 3. Verify Report
Read the human-readable Markdown report returned by the subagent. Formulate a fix strategy using `.agents/skills/operations/troubleshooting_reference.md` based on their analysis.
