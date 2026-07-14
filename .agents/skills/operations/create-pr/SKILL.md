---
name: create-pr
description: "Create a Pull Request (PR) against GoogleCloudPlatform/magic-modules following repository standards, including branch management, commit formatting, mandatory release notes, pre-PR verification checks, and gh CLI commands."
---

# `create-pr`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

This skill provides step-by-step instructions for preparing, formatting, and opening a Pull Request (PR) for `magic-modules` following the official contribution guidelines.

## Prerequisites

* You are in the `magic-modules` root directory.
* Your git working directory is clean except for the files intended for the PR.
* Downstream provider changes are NOT staged or committed to `magic-modules`.
* `gh` CLI is installed and authenticated (`gh auth status`).
* Remote repositories are configured (e.g., `upstream` pointing to `GoogleCloudPlatform/magic-modules` and a personal fork remote such as `origin`).

---

## Pre-PR Verification & Guardrails

Before creating a branch or opening a PR, verify all of the following rules:

1. **Single Self-Contained Change:** Each PR must contain only **one** logical change.
   - Adding multiple resources? Put **one resource per PR**.
   - Fixing a bug and adding new fields? Split into **two separate PRs**.
2. **No Downstream Artifacts in Magic Modules:**
   - Do NOT commit generated downstream provider code (e.g. `$GOPATH/src/github.com/hashicorp/terraform-provider-google`) into `magic-modules`.

---

## Execution Steps

### 1. Sync and Create Feature Branch

Fetch the latest `upstream/main` and create a clean topic branch:

```bash
git fetch upstream main
BRANCH="<short-descriptive-branch-name>" # e.g. add-compute-foo-field
git checkout -b "$BRANCH" upstream/main
```

### 2. Stage and Commit Changes

Stage ONLY the files in `magic-modules` (such as YAML definitions under `mmv1/products/`, templates, or docs):

```bash
git add mmv1/products/<product>/
git commit -m "<product>: <concise description of change>"
```

### 3. Push to Fork

Push the branch to your personal GitHub fork (`origin` or your configured fork remote):

```bash
FORK_REMOTE="origin" # Verify fork remote via `git remote -v`
git push -u "$FORK_REMOTE" "$BRANCH"
```

---

### 4. Format Release Notes & PR Body

Every PR must contain at least one release note block in the PR body.

Refer to the official guide for detailed release note rules, categories, and examples:
* [docs/content/code-review/release-notes.md](../../../docs/content/code-review/release-notes.md)

#### Release Note Block Format
```markdown
```release-note:TYPE
CONTENT
```
```

Common types include `new-resource`, `new-datasource`, `new-list-resource`, `enhancement`, `bug`, `deprecation`, `breaking-change`, `note`, and `none`.

#### Sample PR Body String
```markdown
Summary of what changed and why in a few concise sentences.

Fixes https://github.com/hashicorp/terraform-provider-google/issues/12345

```release-note:enhancement
compute: added `foo` field to `google_compute_instance` resource
```
```

---

### 5. Create Pull Request with `gh` CLI

Construct the PR body string directly and execute `gh pr create` with `--body`:

```bash
PR_TITLE="<product>: <short description>" # e.g. compute: add foo field to google_compute_instance

PR_BODY=$(cat <<'EOF'
<summary of what changed and why>

```release-note:<type>
<release note content>
```
EOF
)

gh pr create \
  --repo GoogleCloudPlatform/magic-modules \
  --base main \
  --head "$(gh api user -q .login):$BRANCH" \
  --title "$PR_TITLE" \
  --body "$PR_BODY"
```

---

## Verification & Handoff

1. Verify that `gh pr create` completed successfully and returned a valid Pull Request URL.
2. Share the PR URL with the user and complete the execution.

