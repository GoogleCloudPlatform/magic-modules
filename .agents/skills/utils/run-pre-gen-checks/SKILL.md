---
name: run-pre-gen-checks
description: "Run fast static checks that don't require generating downstreams, including git submodule verification, gofmt formatting, YAML linting, template validation (version-guard and unused-tmpl), mmv1 core unit tests, internal tool unit tests, and .ci/magician unit tests. Use this skill as a fast initial check for any changes."
---

# `run-pre-gen-checks`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.
> This skill executes pre-generation static checks in parallel directly against `magic-modules` without generating downstream providers.

## Checks Included

- **Submodule Check:** Verifies no unauthorized git submodules are present (`git submodule status --recursive`).
- **Go Formatting Check:** Verifies all tracked Go files are properly formatted (`gofmt -l`).
- **Template Validation:** Runs `version-guard` and `unused-tmpl` checks via `tools/template-check`.
- **YAML Linting Check:** Lints modified product YAML files against `.yamllint` (if `yamllint` is installed).
- **MMv1 Unit Tests:** Runs core unit tests in `mmv1` (`go test ./...`).
- **Tools & CI Unit Tests:** Runs unit tests for `tools/go-changelog`, `tools/issue-labeler`, `tools/template-check`, `tools/test-reader`, and `.ci/magician`.

## Prerequisites

- You must be operating in the `magic-modules` root directory.
- `go` (v1.26+) and `git` must be installed and configured. Optionally `yamllint` for YAML checks.

## Execution Steps

### 1. Run Pre-Generation & Static Checks Runner

Execute the script to run all pre-generation static checks and unit tests in parallel:

```bash
./.agents/skills/utils/run-pre-gen-checks/scripts/run_pre_gen_checks.sh
```

### 2. Output and Verification

- If all checks pass, the script exits with status `0`.
- If any check fails, the script outputs the relevant logs and exits immediately with status `1`.
