---
name: repo-sync
description: "Checks sync status and invokes the subagent to align repositories if needed."
---

# `repo-sync`

Use this skill to check if the downstream workspace is synchronized with Magic Modules and to perform the sync if needed.

## Execution Steps

### 1. Check Sync Status
*   Execute the `check-sync-provider` skill (located in `.agents/skills/utils/check-sync-provider/SKILL.md`) to determine if the downstream repository is aligned with the correct Magic Modules base commit.

### 2. Prompt User if Needed
*   If the workspace is clearly out of sync or if you are unsure, **ask the user** if they would like to run a full repository sync to establish a clean baseline.
*   If the user says no, skip the sync and report back to the workflow.

### 3. Execute Sync
*   If the user agrees or if you are confident a sync is required and safe, use the `invoke_subagent` tool to call the `repo-sync` subagent.
*   **Prompt to send:** "Please initialize and sync the workspace located at the following path: `<insert_absolute_path_here>`"
*   Wait for the subagent to return its final completion message.