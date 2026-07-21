---
name: qa-verification
description: "Triggers the QA Verification subagent to run the verify-changes workflow and interpret results."
---

# `qa-verification`

Use this skill to delegate running the verification workflow and log parsing to a dedicated Subagent. This keeps your main context window clean.

## Execution Steps

### 1. Identify Target and Context
Gather context on *what* changed (the resource, the field, the PR) and which services or tests need verification.

### 2. Invoke Subagent
Use the `invoke_subagent` tool. Pass the context of the change as part of the trigger prompt so the subagent knows what it's verifying!

**Trigger Example:**
`invoke_subagent("qa-verification", "Run the verify_changes workflow. We just made the following change: <CHANGE_DESCRIPTION>. Verify if all CI checks pass or fail.")`

### 3. Verify Report
Read the human-readable Markdown report returned by the subagent. If any phase failed, formulate a fix strategy using `.agents/skills/operations/troubleshooting_reference.md` based on their analysis.
