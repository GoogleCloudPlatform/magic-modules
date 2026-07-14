---
name: check-schema-diff
description: "Run comprehensive checks (breaking changes, missing tests, and missing documentation) on Magic Modules schema changes between the current branch and a base branch (usually main) using the local diff-processor."
---

# `check-schema-diff`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.
> This skill allows you to run local breaking changes, missing tests, and missing documentation checks on Magic Modules schema changes.

## Prerequisites

- You must be operating in the `magic-modules` root directory.
- You must have `go` installed and configured.

## Execution Steps

### 1. Make the script executable (if not already)
Ensure the script is executable:
```bash
chmod +x .agents/skills/utils/check-schema-diff/scripts/check_schema_diff.sh
```

### 2. Run the Schema Diff Checker
Run the script to compare your current local changes (committed or uncommitted) against the base branch (defaults to the merge base with `origin/main` or `main`):
```bash
./.agents/skills/utils/check-schema-diff/scripts/check_schema_diff.sh
```

To compare against a specific branch or commit, pass it as an argument:
```bash
./.agents/skills/utils/check-schema-diff/scripts/check_schema_diff.sh <base_branch_or_commit>
```

### 3. Analyze the Output
The tool runs three checks across both GA and Beta provider versions:
- **Breaking Changes:** Detects backwards-incompatible schema changes.
- **Missing Tests:** Identifies new or changed fields that are not covered by any acceptance tests.
- **Missing Documentation:** Identifies new fields that are not documented in the resource/datasource markdown files.

### 4. Verification & Handoff
If any issues (breaking changes, missing tests, or missing docs) are detected, present them clearly to the user. Discuss potential mitigations and fixes before proceeding.
