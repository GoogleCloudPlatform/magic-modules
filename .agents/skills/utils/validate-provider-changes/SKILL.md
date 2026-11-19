---
name: validate-provider-changes
description: "Validate that changes to the generated providers don't introduce breaking changes, missing tests, unreadable tests, missing documentation, or missing ResourceIdentity coverage. Use this skill when the user wants to check for breaking changes, missing or unreadable tests, missing documentation, missing identity coverage, or other diff-processor problems in the downstream providers."
---

# `validate-provider-changes`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.
> This skill runs all 5 `diff-processor` validation checks from CI (`.ci/magician/cmd/generate_comment.go`) locally across both GA and Beta provider changes.

## Prerequisites

- You must be operating in the `magic-modules` root directory.
- You must have `go` installed and configured.

## Execution Steps

### 1. Run the Provider Changes Validator

Run the script to compare your current local changes (committed or uncommitted) against the base branch (defaults to the merge base with `origin/main`, `upstream/main`, or `main`):

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

To enable exhaustive missing-test analysis (verifying `SET`, `MODIFY`, and `UNSET` transitions for newly added fields at every nesting level via `diff-processor detect-missing-tests --exhaustive`):

```bash
./.agents/skills/utils/validate-provider-changes/scripts/validate_provider_changes.sh --exhaustive [base_branch_or_commit]
```

Pass `--help` (or `-h`) to display usage and options.

### 2. Analyze the Output

The tool runs all 5 `diff-processor` checks from CI (`.ci/magician/cmd/generate_comment.go`) across **both GA (`terraform-provider-google`) and Beta (`terraform-provider-google-beta`)**. While CI runs checks 2–5 on Beta only, running all 5 on both GA and Beta locally also catches GA-only coverage gaps (e.g., Beta → GA promotions):

- **Breaking Changes (`breaking-changes`) — GA & Beta:** Detects backwards-incompatible schema changes across both `terraform-provider-google` (GA) and `terraform-provider-google-beta` (Beta).
- **Missing Tests (`detect-missing-tests`) — GA & Beta:** Identifies new or changed fields in `google/services` (GA) and `google-beta/services` (Beta) that are not covered by acceptance tests. When `--exhaustive` is passed, additionally checks that newly added fields have `SET`, `MODIFY`, and `UNSET` test coverage at every nesting level.
- **Unreadable Tests (`detect-unreadable-tests`) — GA & Beta:** Checks changed service directories across both GA and Beta for acceptance tests that the test reader cannot parse (and that are not listed in `KnownUnreadableTests`), as well as stale entries in `KnownUnreadableTests` that can now be read.
- **Missing Documentation (`detect-missing-docs`) — GA & Beta:** Identifies new resource or data source fields that are not documented in `website/docs` markdown files across both GA and Beta.
- **Missing ResourceIdentity Coverage (`detect-missing-identity`) — GA & Beta:** Checks changed resources declaring `ResourceIdentity` across both GA and Beta for missing `SetResourceIdentityAttributes` calls in CRUD functions or missing `ImportBlockWithResourceIdentity` import tests.

### 3. Verification & Handoff

If any issues (breaking changes, missing tests, unreadable tests, missing docs, or missing ResourceIdentity coverage) are detected, present them clearly to the user. Discuss potential mitigations and fixes before proceeding.
