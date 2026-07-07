---
name: validate-provider-changes
description: "Validate that changes to the generated providers don't introduce breaking changes or fields with missing tests or missing documentation. Use this skill when the user wants to check for breaking changes, missing tests, or missing documentation, or other problems in the downstream providers."
---

# `validate-provider-changes`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.
> This skill allows you to run local breaking changes, missing tests, and missing documentation checks on Magic Modules schema and provider changes.

## Prerequisites

- You must be operating in the `magic-modules` root directory.
- You must have `go` installed and configured.

## Execution Steps

### 1. Run the Provider Changes Validator

Run the script to compare your current local changes (committed or uncommitted) against the base branch (defaults to the merge base with `origin/main` or `main`):

```bash
./.agents/skills/utils/validate-provider-changes/scripts/validate_provider_changes.sh
```

To compare against a specific branch or commit (e.g. `HEAD` to check uncommitted changes against the current commit):

```bash
./.agents/skills/utils/validate-provider-changes/scripts/validate_provider_changes.sh HEAD
```

Or to compare against a base branch:

```bash
./.agents/skills/utils/validate-provider-changes/scripts/validate_provider_changes.sh <base_branch_or_commit>
```

### 2. Analyze the Output

The tool runs three checks across both GA and Beta provider versions:

- **Breaking Changes:** Detects backwards-incompatible schema changes.
- **Missing Tests:** Identifies new or changed fields that are not covered by any acceptance tests.
- **Missing Documentation:** Identifies new fields that are not documented in the resource/datasource markdown files.

### 3. Verification & Handoff

If any issues (breaking changes, missing tests, or missing docs) are detected, present them clearly to the user. Discuss potential mitigations and fixes before proceeding.
