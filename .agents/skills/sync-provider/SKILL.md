---
name: sync-provider
description: "Synchronize a downstream Terraform provider repository with Magic Modules by aligning commit history and verifying parity."
---

# `sync-provider`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.
> This skill is designed to be completely self-contained and unambiguous for a fresh agent without prior context.

## Prerequisites

- You must be operating relative to the `magic-modules` and downstream provider repositories.
- You must have the absolute path to the downstream repository.
- You must have verified there are no unsaved or uncommitted changes in the downstream provider directory that would be overwritten by code generation.

## How to Choose Your Path

Before proceeding, determine which synchronization method you need:

1. **Aligning to a Specific Base Commit**: Use this if you need to synchronize the repositories to a specific state in the past (e.g., to match a specific pull request or isolate a failure). Follow the numbered **Execution Steps** below.
2. **Synchronizing to Latest (Fast-Forward)**: Use this if you simply want to bring both repositories up to date with their latest remote commits. Follow the steps in the **Synchronizing to Latest** section.

## Execution Steps

### 1. Identify Magic Modules Base Commit

Check if the `magic-modules` repository is a fork or has multiple remotes:
```bash
git remote -v
```

Identify the official `googleapis/magic-modules` remote (e.g., `origin` or `upstream`). 
If it is ambiguous, **ask the user** which remote is the canonical upstream.

Then calculate the target base commit:

- **If on `main` (of the upstream):** Use the current commit (`HEAD`).
- **If on a feature branch:** Calculate the merge base where the branch diverged from the canonical `main`:
    ```bash
    git merge-base HEAD <canonical-remote>/main
    ```

Record the base commit hash and date. This hash will be used to find the matching commit in the downstream repository.

### 2. Check for Uncommitted Downstream Changes

Before checking out older commits, check if there are local modifications in the downstream repository:
```bash
cd <downstream-provider-path>
git status
```
If changes exist, always clean the uncommitted and untracked files in the downstream repository using `git reset --hard` and `git clean -fd`.


### 3. Find the Matching Downstream Commit

Search for a commit in the downstream repository that corresponds to the Magic Modules base commit hash:
```bash
cd <downstream-provider-path>
git log -n 50 --grep="[upstream:<MM-hash>]"
```

If not found, search `origin/main` (in case your local `main` is behind):
```bash
git log -n 50 --grep="[upstream:<MM-hash>]" origin/main
```
If found, pull your `main` or checkout the specific commit. If not found, try history traversal before falling back to date-based matching.

#### Fallback: History Traversal (Recommended)

If the exact base hash is not found, it is likely because the specific `magic-modules` commit did not generate code changes for this specific provider (resulting in no downstream generation commit). 

Walk backwards through the `magic-modules` commit history from your base commit. For each parent commit:

1. View previous commit history: `git log -n 10 --format="%H %s"`
2. Search for its hash in the downstream: `git log --grep="[upstream:<parent-hash>]"`
3. The first commit found is the true synchronization point where the downstream repository diverges or sits.

#### Fallback: Date-Based Matching

If no commit message contains the hash, find the commit in the downstream repository that is closest to the date/timestamp of the MM base commit:
```bash
git log --since="<date-of-base-minus-2-days>" --until="<date-of-base-plus-2-days>"
```
Select the commit that appears to be the matching nightly or generation commit. If ambiguous, **ask the user for clarification**.

### 4. Align Downstream

Checkout the matching commit in the downstream repository:
```bash
git checkout <matching-commit-hash>
```

### 5. Verify Parity

Return to `magic-modules` and run the automation script from `tgc-build-skill` to generate code and verify parity:
```bash
./.agents/skills/tgc-build-skill/scripts/build_tgc.sh <downstream-provider-path>
```
Verify the output of `git status` in the downstream repository. It should be clean or only contain changes from our specific branch.

### 6. Verification & Handoff

Verify that the `make tgc` command succeeded without unexpected diffs. Return to the workflow that invoked this skill to proceed.

---

## Synchronizing to Latest (Fast-Forward)

If your goal is to bring both repositories to their most recent remote commits rather than aligning to an older base:

### 1. Update Magic Modules to Latest
Fetch the latest changes from the canonical remote (e.g., `upstream` or `origin`) and rebase your feature branch on it. This avoids failures if `main` is already checked out in another worktree:
```bash
git fetch <canonical-remote> main
git rebase <canonical-remote>/main
```

### 2. Clean and Update Downstream to Latest
Downstream is generated, so clean local changes first to avoid conflicts. Then fetch the latest commit from the remote main branch:
```bash
cd <downstream-provider-path>
git reset --hard
git clean -fd
git fetch origin main
```

### 3. Project Latest Changes
Return to `magic-modules` and run the build script to project all changes to the latest downstream state:
```bash
./.agents/skills/tgc-build-skill/scripts/build_tgc.sh <downstream-provider-path>
```