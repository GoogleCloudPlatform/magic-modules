---
name: check-breaking-changes
description: "Detect breaking schema changes between the current branch and a base branch (usually main) using the local diff-processor."
---

# `check-breaking-changes`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.
> This skill allows you to run local breaking change checks on Magic Modules schema changes.

## Prerequisites

- You must be operating in the `magic-modules` root directory.
- You must have `go` installed and configured.

## Execution Steps

### 1. Make the script executable (if not already)
Ensure the script is executable:
```bash
chmod +x .agents/skills/utils/check-breaking-changes/scripts/check_breaking_changes.sh
```

### 2. Run the Breaking Changes Checker
Run the script to compare your current local changes (committed or uncommitted) against the base branch (defaults to the merge base with `origin/main` or `main`):
```bash
./.agents/skills/utils/check-breaking-changes/scripts/check_breaking_changes.sh
```

To compare against a specific branch or commit, pass it as an argument:
```bash
./.agents/skills/utils/check-breaking-changes/scripts/check_breaking_changes.sh <base_branch_or_commit>
```

### 3. Analyze the Output
- **No breaking changes detected:** If the tool outputs `No breaking changes detected!`, your schema changes are backwards-compatible.
- **Breaking changes detected:** The tool will print a JSON list of breaking changes, including the affected fields, resources, and the specific rule violated.

### 4. Verification & Handoff
If breaking changes are detected, present them clearly to the user. Explain why each change is considered breaking and discuss potential mitigations (e.g., postponing to a major release, using a default value, or adding deprecation warnings first).
