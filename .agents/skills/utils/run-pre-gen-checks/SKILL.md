---
name: run-pre-gen-checks
description: "Run fast static checks that don't require generating downstreams, including git submodule verification, gofmt formatting, YAML linting, template validation (version-guard, unused-tmpl, and disallow-legacy-apiary), mmv1 core unit tests, internal tool unit tests, and .ci/magician unit tests. Use this skill as a fast initial check for any changes."
---

# `run-pre-gen-checks`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.
> This skill executes pre-generation static checks in parallel directly against `magic-modules` without generating downstream providers.

## Checks Included

Diff-based checks compare your committed, staged, unstaged, and untracked changes against the merge-base with `upstream/main` (falling back to `origin/main`, `main`, or `HEAD`).

- **Submodule Check:** Verifies no unauthorized git submodules are present (`git submodule status --recursive`).
- **Go Formatting Check:** Verifies all tracked and untracked (non-deleted) Go files are properly formatted (`gofmt -l`).
- **Template Validation (`tools/template-check`):**
  - `version-guard`: Checks modified and newly added (non-deleted, including untracked) `.tmpl` files for invalid version guards.
  - `unused-tmpl`: Checks newly added `.tmpl` files under `mmv1/` (committed, staged, or untracked) to ensure they are referenced by product YAMLs or generator code.
  - `disallow-legacy-apiary`: Checks tracked diffs and untracked files under `mmv1/templates/terraform/` and `mmv1/third_party/{terraform,cai2hcl,tgc*}/` for forbidden legacy Apiary Compute client imports or invocations, mirroring CI.
- **YAML Linting Check:** Lints modified and newly added (non-deleted) product YAML files under `mmv1/products/` against `.yamllint` (if `yamllint` is installed).
- **MMv1 Unit Tests:** Runs core unit tests in `mmv1` (`go test ./...`).
- **Tools & CI Unit Tests:** Diff-scoped unit tests for `tools/go-changelog`, `tools/issue-labeler`, `tools/template-check`, `tools/test-reader`, `.ci/magician`, and `.agents/scripts` (only runs suites for modified/untracked tool or CI files, or all suites when `--all` is passed).

## Prerequisites

- Operating inside the `magic-modules` repository or worktree.
- `go` (v1.26+) and `git` must be installed and configured. Optionally `yamllint` for YAML checks.

## Execution Steps

### 1. Run Pre-Generation & Static Checks Runner

Execute the script to run all pre-generation static checks and unit tests in parallel:

```bash
./.agents/skills/utils/run-pre-gen-checks/scripts/run_pre_gen_checks.sh
```

Pass `--all` if you want to force running unit tests across all internal tools and CI automation regardless of whether they were modified:

```bash
./.agents/skills/utils/run-pre-gen-checks/scripts/run_pre_gen_checks.sh --all
```

### 2. Output and Verification

- If all checks pass, the script exits with status `0`.
- If any check fails, the script outputs the relevant logs and exits immediately with status `1`.
