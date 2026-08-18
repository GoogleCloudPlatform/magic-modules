---
name: prepare-pr-link
description: "Prepares a feature branch, stages and commits changes, pushes to the user's personal fork, and generates a pre-filled GitHub comparison link for the user to review and submit a Pull Request."
---

# `prepare-pr-link`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

This skill provides step-by-step instructions for preparing, committing, pushing a feature branch to a personal fork, and generating a pre-filled GitHub Pull Request comparison link for the user to review and submit.

## Prerequisites

* You are in the `magic-modules` root directory.
* Your git working directory is clean except for the files intended for the PR.
* Downstream provider changes are NOT staged or committed in `magic-modules`.
* Remote repositories are configured (e.g., `upstream` pointing to `GoogleCloudPlatform/magic-modules` and a personal fork remote such as `origin`).

---

## Pre-PR Verification & Guardrails

Before creating a branch or generating the link, verify all of the following rules:

1. **Single Self-Contained Change:** Each PR must contain only **one** logical change.
   - Adding multiple resources? Put **one resource per PR**.
   - Fixing a bug and adding new fields? Split into **two separate PRs**.
2. **Strict PR Title Length Limit:**
   - PR title must be concise and strictly **under 70 characters**.
   - Format: `<product>: <action> <target>` (e.g. `beyondcorp: deprecate google_beyondcorp_app_*` or `compute: add foo field to google_compute_instance`).
3. **No Downstream Artifacts in Magic Modules:**
   - Do NOT commit generated downstream provider code into `magic-modules`.
4. **Workspace Cleanup:** Run `git status --porcelain` and ensure no untracked temporary test files exist before pushing.

---

## Execution Steps

### 1. Sync and Create Feature Branch

```bash
UPSTREAM_REMOTE=$(git remote -v | grep -i "GoogleCloudPlatform/magic-modules" | head -n 1 | awk '{print $1}')
UPSTREAM_REMOTE="${UPSTREAM_REMOTE:-upstream}"

BASE_BRANCH="main" # or target major release feature branch, e.g. FEATURE-BRANCH-major-release-8.0.0
git fetch "$UPSTREAM_REMOTE" "$BASE_BRANCH"

BRANCH="<short-descriptive-branch-name>" # e.g. deprecate-beyondcorp-app
git checkout -b "$BRANCH" "$UPSTREAM_REMOTE/$BASE_BRANCH"
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

Refer to [docs/content/code-review/release-notes.md](../../../../docs/content/code-review/release-notes.md) for details on categories (`enhancement`, `bug`, `none`, `new-resource`, `new-datasource`, `new-list-resource`, `deprecation`, `breaking-change`).

#### Sample PR Body
```markdown
Summary of what changed and why in a few concise sentences.

Fixes https://github.com/hashicorp/terraform-provider-google/issues/12345

```release-note:enhancement
<product>: added `foo` field to `google_compute_instance`
```
```

---

### 5. Generate Pre-Filled Web URL (Always Provide as Hyperlink)

Run Python to generate the pre-filled comparison URL with title and body URL-encoded:

```python
import urllib.parse
import subprocess
import re

base_branch = "main" # or target feature branch, e.g. "FEATURE-BRANCH-major-release-8.0.0"
branch = "<BRANCH>"
title = "<PR_TITLE>"
body = """<PR_BODY>"""

# Auto-detect fork remote and username
fork_remote_cmd = "git remote -v | grep -v -i 'GoogleCloudPlatform/magic-modules' | head -n 1 | awk '{print $1}'"
fork_remote = subprocess.check_output(fork_remote_cmd, shell=True, text=True).strip() or "origin"
fork_url = subprocess.check_output(["git", "remote", "get-url", fork_remote], text=True).strip()

# Extract username from git URL (ssh or https)
m = re.search(r"github\.com[:/]([^/]+)/", fork_url)
username = m.group(1) if m else "<username>"

head_ref = f"{username}:{branch}"
url = f"https://github.com/GoogleCloudPlatform/magic-modules/compare/{base_branch}...{head_ref}?expand=1&title={urllib.parse.quote(title)}&body={urllib.parse.quote(body)}"
print(url)
```

Present the result in chat as:
👉 **[Create Pull Request on GitHub](<URL>)**

Inform the user that the branch has been pushed to their fork and they can click the link to review the diff and submit the Pull Request.
