---
name: run-pre-gen-checks
description: "Run pre-generation static checks including gofmt formatting, template validation (version-guard and unused-tmpl), mmv1 core unit tests, and internal tool unit tests. Use this skill to run static checks directly against magic-modules without generating downstream providers."
---

# `run-pre-gen-checks`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.
> This skill executes pre-generation static checks in parallel directly against `magic-modules` without generating downstream providers.

## Prerequisites

- You must be operating in the `magic-modules` root directory.
- `go` (v1.26+) and `git` must be installed and configured.

## Execution Steps

### 1. Run Pre-Generation & Static Checks Runner

Execute the script to run gofmt, template validation checks, mmv1 unit tests, and internal tool unit tests in parallel:

```bash
./.agents/skills/utils/run-pre-gen-checks/scripts/run_pre_gen_checks.sh
```

### 2. Output and Verification

- If all formatting, template checks, and unit tests pass, the script exits with status `0`.
- If any check fails, the script exits immediately with status `1`.
