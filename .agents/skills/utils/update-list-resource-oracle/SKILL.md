---
name: update-list-resource-oracle
description: "Update the list-resource oracle scripts (run-pre-gen-checks and validate-provider-changes) when the eligibility rules, static checks, or schema-diff logic need to change. Commits the updates to a dedicated oracle branch — separate from any list-resource PR — so a human can review and merge them into upstream magic-modules before the agent system picks them up on the next run."
---

# `update-list-resource-oracle`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

Oracle changes are **never** mixed into a list-resource PR. This skill produces a
standalone commit on `update-list-resource-oracle` that a human merges separately
into `GoogleCloudPlatform/magic-modules:main`. Future agent runs pick up the
updated scripts automatically once that PR lands.

## What the oracle consists of

| Script | Path | What it checks |
|---|---|---|
| `run_pre_gen_checks.sh` | `.agents/skills/utils/run-pre-gen-checks/scripts/run_pre_gen_checks.sh` | Submodule guard · `gofmt` · `version-guard` + `unused-tmpl` template checks · YAML lint against `.yamllint` · `mmv1` unit tests · `tools/` and `.ci/magician` unit tests |
| `validate_provider_changes.sh` | `.agents/skills/utils/validate-provider-changes/scripts/validate_provider_changes.sh` | Breaking schema changes · missing acceptance tests · missing resource/datasource documentation (GA + Beta) |

The **eligibility scan** embedded in `.agents/skills/workflows/add_list_resource/SKILL.md`
is also part of the oracle surface: it defines which resources are candidates for
`generate_list_resource: true`. If the eligibility rules change (e.g. a new
auto-scope parameter is added, or an exclusion flag is introduced), update the
Python scan in that SKILL as well and include it in the same oracle commit.

## Prerequisites

* You are in the `magic-modules` root directory.
* `git` is configured with the fork remote (e.g. `origin` → `BBBmau/magic-modules`).
* You have identified **exactly what needs to change** in the oracle and **why**
  (e.g. a new YAML field that breaks the YAML lint rule, a new scope parameter
  the eligibility scan doesn't recognise, a `diff-processor` flag that changed).
  Do not make speculative edits.

## Execution Steps

### 1. Verify you are on the right branch

The oracle branch is fixed and dedicated. Never commit oracle changes to a
list-resource branch or to `main` directly.

```bash
git fetch upstream main
git fetch origin update-list-resource-oracle 2>/dev/null || true

# If the oracle branch already exists in the fork, resume it.
# Otherwise create it fresh from upstream/main.
if git ls-remote --exit-code origin update-list-resource-oracle >/dev/null 2>&1; then
  git checkout update-list-resource-oracle
else
  git checkout -b update-list-resource-oracle upstream/main
fi
```

### 2. Make the targeted edit

Edit only the specific lines that need to change. The scope of a valid oracle
update is one of:

**a) `run_pre_gen_checks.sh` — add, remove, or adjust a static check**

```bash
# Example: add a new parallel check block inside the script.
# Follow the existing pattern: background the check, capture its PID,
# wait for it, check its exit status, print its log at the end.
$EDITOR .agents/skills/utils/run-pre-gen-checks/scripts/run_pre_gen_checks.sh
```

**b) `validate_provider_changes.sh` — adjust schema-diff logic**

```bash
# Example: pass a new flag to diff-processor, or adjust the missing-docs
# detection logic for a changed JSON schema.
$EDITOR .agents/skills/utils/validate-provider-changes/scripts/validate_provider_changes.sh
```

**c) Eligibility scan in `add_list_resource/SKILL.md` — update candidate rules**

```bash
# Example: add a new auto-scope parameter so resources using it become eligible.
# Locate the AUTO_SCOPES set in the embedded Python script and extend it.
$EDITOR .agents/skills/workflows/add_list_resource/SKILL.md
```

### 3. Self-verify the oracle change

Run the updated script against the current `main` to confirm it exits 0 and
produces sensible output before committing. Use `HEAD` so you are checking the
uncommitted edit:

```bash
# For run-pre-gen-checks:
bash .agents/skills/utils/run-pre-gen-checks/scripts/run_pre_gen_checks.sh

# For validate-provider-changes:
bash .agents/skills/utils/validate-provider-changes/scripts/validate_provider_changes.sh HEAD
```

If either script exits non-zero, fix the script before proceeding. Do not
commit a broken oracle.

### 4. Commit only the oracle files

Stage only the files you changed. Do NOT stage anything under `mmv1/products/`.

```bash
git add .agents/skills/utils/run-pre-gen-checks/scripts/run_pre_gen_checks.sh
git add .agents/skills/utils/validate-provider-changes/scripts/validate_provider_changes.sh
git add .agents/skills/workflows/add_list_resource/SKILL.md
# Only add the files you actually changed — remove any that were not touched.

git status --porcelain   # confirm only oracle files are staged

git commit -m "oracle: <one-line description of what changed and why>"
```

### 5. Push to the fork

```bash
git push --force-with-lease origin update-list-resource-oracle
```

### 6. Report back

Output a summary of:
- Which file(s) were changed and what lines were modified.
- Why the change was necessary (the specific failure or gap it addresses).
- The commit SHA.
- A reminder that a human must open and merge the PR from
  `origin:update-list-resource-oracle` → `GoogleCloudPlatform/magic-modules:main`
  before the agent system picks up the updated oracle.

**Do NOT open the PR yourself.** This is a deliberate human gate. The oracle
defines what "correct" means for all future agent runs — it must be reviewed
before merging.

## Handoff & Guardrails

* **One concern per commit.** If two unrelated oracle scripts need changes,
  use two separate commits on the same branch (not two separate branches).
* **Never edit the generator** (`mmv1/api/`, `mmv1/provider/`,
  `mmv1/templates/terraform/list_resource*`). If the generator is the problem,
  escalate to the user.
* **Never weaken a check to make it pass.** If a resource genuinely cannot pass
  the oracle, remove it from the list-resource candidate list — do not relax
  the oracle rule.
* **The oracle branch is rebased, not merged.** If `update-list-resource-oracle`
  has fallen behind `upstream/main`, rebase it:
  ```bash
  git fetch upstream main
  git rebase upstream/main update-list-resource-oracle
  git push --force-with-lease origin update-list-resource-oracle
  ```
