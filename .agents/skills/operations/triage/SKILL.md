---
name: triage
description: "Gather context on a change or bug and plan the implementation."
---

# `triage`

Use this skill when a user asks you to implement a feature or fix a bug. Gather context and plan before writing code.

## Execution Steps

### 1. Gather Context (Ask User First)
Propose multiple places to look (Resource YAML definitions and Cloud API documentation) and ask the user where to start. Plan the change (New feature or bug fix) within schema or logic.

### 2. Formulate a Proposal
- Write an `implementation_plan.md` in your temporary artifacts.
- State the intent of the change.
- List the files to modify.

### 3. Verification & Handoff
- Halt. Discuss the plan with the user.
- Once approved and the change is made, proceed back to the outer loop (Step 3: Generate).
