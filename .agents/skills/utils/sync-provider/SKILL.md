---
name: sync-provider
description: "Synchronize a downstream Terraform provider repository with Magic Modules by aligning commit history and verifying parity."
---

# `sync-provider`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.
> This skill is designed to be completely self-contained and unambiguous for a fresh agent without prior context.

## Prerequisites

- You must be operating relative to the `magic-modules` and downstream provider repositories.
- You must have the absolute path to the downstream repository (e.g., `terraform-provider-google-beta`).
- You must have verified there are no unsaved or uncommitted changes in the downstream provider directory that would be overwritten by code generation.

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
If changes exist, point them out to the user and **ask if they wish to save/stash them**. Do not proceed without consent. If consent is given, you can stash or reset them.

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

Return to `magic-modules` and run code generation to verify that there are no unexpected diffs:
```bash
cd <magic-modules-path>
make provider VERSION=<beta|ga> OUTPUT_PATH="<downstream-provider-path>"
```
Verify the output of `git status` in the downstream repository. It should be clean or only contain changes from our specific branch.

### 6. Verification & Handoff

Verify that the `make provider` command succeeded without unexpected diffs. Return to the workflow that invoked this skill to proceed.
