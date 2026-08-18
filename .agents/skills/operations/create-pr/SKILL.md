---
name: create-pr
description: "Create a Pull Request (PR) against GoogleCloudPlatform/magic-modules following repository standards, including branch management, commit formatting, mandatory release notes, pre-PR verification checks, and gh CLI commands."
---

# `create-pr`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

This skill provides step-by-step instructions for preparing, formatting, and opening a Pull Request (PR) for `magic-modules` following official contribution guidelines.

## Prerequisites

* You are in the `magic-modules` root directory.
* Your git working directory is clean except for the files intended for the PR.
* Downstream provider changes are NOT staged or committed to `magic-modules`.
* Remote repositories are configured (e.g., `upstream` pointing to `GoogleCloudPlatform/magic-modules` and a personal fork remote such as `origin`).

---

## Pre-PR Verification & Guardrails

Before creating a branch or opening a PR, verify all of the following rules:

1. **Single Self-Contained Change:** Each PR must contain only **one** logical change.
   - Adding multiple resources? Put **one resource per PR**.
   - Fixing a bug and adding new fields? Split into **two separate PRs**.
2. **Strict PR Title Length Limit:**
   - PR title must be concise and strictly **under 70 characters**.
   - Format: `<product>: <action> <target>` (e.g. `beyondcorp: deprecate google_beyondcorp_app_*` or `compute: add foo field to google_compute_instance`).
3. **No Downstream Artifacts in Magic Modules:**
   - Do NOT commit generated downstream provider code into `magic-modules`.
4. **Workspace Cleanup:** Run `git status --porcelain` and ensure no untracked temporary test files exist before opening the PR.

---

## Execution Steps

### 1. Sync and Create Feature Branch

```bash
UPSTREAM_REMOTE=$(git remote -v | grep -i "GoogleCloudPlatform/magic-modules" | head -n 1 | awk '{print $1}')
UPSTREAM_REMOTE="${UPSTREAM_REMOTE:-upstream}"

git fetch "$UPSTREAM_REMOTE" main
BRANCH="<short-descriptive-branch-name>" # e.g. deprecate-beyondcorp-app
git checkout -b "$BRANCH" "$UPSTREAM_REMOTE/main"
```

### 2. Stage and Commit Changes

```bash
git add mmv1/products/<product>/ # or other relevant files
git commit -m "<product>: <concise description under 70 chars>"
```

### 3. Push to Personal Fork

```bash
FORK_REMOTE=$(git remote -v | grep -v -i "GoogleCloudPlatform/magic-modules" | head -n 1 | awk '{print $1}')
FORK_REMOTE="${FORK_REMOTE:-origin}"

git push -u "$FORK_REMOTE" "$BRANCH"
```

---

### 4. Format PR Body & Release Notes

Every PR must contain a clear summary and a release note block in the PR body.

Refer to [docs/content/code-review/release-notes.md](../../../../docs/content/code-review/release-notes.md) for details on categories (`enhancement`, `bug`, `none`, `new-resource`, `deprecation`, `breaking-change`).

#### Sample PR Body
```markdown
Summary of what changed and why in a few concise sentences.

Fixes https://github.com/hashicorp/terraform-provider-google/issues/12345

```release-note:enhancement
<product>: added `foo` field to `google_compute_instance`
```
```

---

### 5. Create PR or Provide Pre-Filled Web Link

Write the PR body to a temporary file (`/tmp/pr_body.txt`) using a single-quoted HEREDOC:

```bash
PR_TITLE="<product>: <short description under 70 chars>"
BASE_BRANCH="main" # or FEATURE-BRANCH-major-release-8.0.0

cat <<'EOF' > /tmp/pr_body.txt
<summary of what changed and why>

Fixes <issue link if applicable>

```release-note:<type>
<product>: <release note description>
```
EOF
```

#### Attempt `gh pr create`:
```bash
gh pr create \
  --repo GoogleCloudPlatform/magic-modules \
  --base "$BASE_BRANCH" \
  --head "$(gh api user -q .login):$BRANCH" \
  --title "$PR_TITLE" \
  --body-file /tmp/pr_body.txt
```

#### Generate Pre-Filled Web URL (Always Provide as Hyperlink):
Always generate and present a clickable markdown link with the PR title and description pre-filled in query parameters for easy user review and submission:

```python
import urllib.parse

base_branch = "main" # or "FEATURE-BRANCH-major-release-8.0.0"
head_ref = f"{username}:{branch}" # e.g. "c2thorn:deprecate-beyondcorp-app"
title = "..."
body = "..."

url = f"https://github.com/GoogleCloudPlatform/magic-modules/compare/{base_branch}...{head_ref}?expand=1&title={urllib.parse.quote(title)}&body={urllib.parse.quote(body)}"
print(url)
```

Present the result in chat as:
👉 **[Create Pull Request on GitHub](<URL>)**

---

## Verification & Summary

1. If `gh pr create` succeeds, view the published PR: `gh pr view --repo GoogleCloudPlatform/magic-modules`.
2. Share the confirmed PR URL (or the pre-filled direct compare URL) with the user.
