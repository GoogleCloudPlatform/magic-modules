---
name: check-sync-provider
description: "Checks if the downstream provider repository is aligned with the correct Magic Modules base commit."
---

# `check-sync-provider`

Use this skill to verify if the downstream repository is at the correct baseline commit before generating code or running tests.

## Execution Steps

### 1. Identify Magic Modules Base Commit
*   Identify the official `GoogleCloudPlatform/magic-modules` remote (usually `upstream`).
*   Calculate the target base commit:
    *   If on `main`: Use current commit (`HEAD`).
    *   If on a feature branch: Run `git merge-base HEAD <canonical-remote>/main` to find where the branch diverged.

### 2. Find the Matching Downstream Commit
*   Search the downstream repository's log for a commit that references this Magic Modules base commit hash in its commit message.

### 3. Compare and Report
*   Verify if the `HEAD` of the downstream repository matches that found commit.
*   If they match, report **Synced**.
*   If they do not match, or the commit cannot be found, report **Out of Sync** and provide the details.
