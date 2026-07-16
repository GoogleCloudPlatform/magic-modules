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

Identify the upstream remote (pointing to `GoogleCloudPlatform/magic-modules`), fetch `main`, and create a clean topic branch:

```bash
# Discover the remote for GoogleCloudPlatform/magic-modules (defaults to 'upstream' if unmatched)
UPSTREAM_REMOTE=$(git remote -v | grep -i "GoogleCloudPlatform/magic-modules" | head -n 1 | awk '{print $1}')
UPSTREAM_REMOTE="${UPSTREAM_REMOTE:-upstream}"

git fetch "$UPSTREAM_REMOTE" main
BRANCH="<short-descriptive-branch-name>" # e.g. add-compute-foo-field
git checkout -b "$BRANCH" "$UPSTREAM_REMOTE/main"
```

### 2. Stage and Commit Changes

Stage ONLY the files in `magic-modules` (such as YAML definitions under `mmv1/products/`, templates, or docs):

```bash
git add mmv1/products/<product>/
git commit -m "<product>: <concise description of change>"
```

### 3. Push to Fork

Identify your personal fork remote. If multiple non-upstream remotes exist or the target remote is ambiguous, stop and ask the user which remote to push to:

```bash
# Check available remotes
git remote -v

# Automatically select personal fork remote (remote not pointing to GoogleCloudPlatform).
# If multiple non-upstream remotes exist, ask the user to confirm the remote name.
FORK_REMOTE=$(git remote -v | grep -v -i "GoogleCloudPlatform/magic-modules" | head -n 1 | awk '{print $1}')
FORK_REMOTE="${FORK_REMOTE:-origin}"

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

#### Sample PR Body Content
```markdown
Summary of what changed and why in a few concise sentences.

Fixes https://github.com/hashicorp/terraform-provider-google/issues/12345

```release-note:enhancement
compute: added `foo` field to `google_compute_instance` resource
```
```

---

### 5. Create Pull Request with `gh` CLI Using `--body-file`

> [!CAUTION]
> **DO NOT pass inline double-quoted `--body "..."` strings containing backticks.**
> Enclosing triple backticks (` ```release-note:type ``` `) in double quotes causes `zsh`/`bash` to execute `` `release-note:type` `` as a live shell command substitution. The command fails, silently stripping the release note block from the published PR body!

Always write the body to a temporary file via a single-quoted HEREDOC (`cat <<'EOF'`) and invoke `gh pr create` with `--body-file`:

```bash
PR_TITLE="<product>: <short description>" # e.g. compute: add foo field to google_compute_instance

cat <<'EOF' > /tmp/pr_body.txt
<summary of what changed and why>

```release-note:<type>
<release note content>
```
EOF

gh pr create \
  --repo GoogleCloudPlatform/magic-modules \
  --base main \
  --head "$(gh api user -q .login):$BRANCH" \
  --title "$PR_TITLE" \
  --body-file /tmp/pr_body.txt
```

---

## Verification & Auto-Repair Handoff

1. Verify that `gh pr create` completed successfully and returned a valid Pull Request URL.
2. **Verify PR Body Integrity:** Run `gh pr view` to verify the published body rendered the release note block:
   ```bash
   gh pr view --repo GoogleCloudPlatform/magic-modules
   ```
3. **Auto-Repair Missing Release Note:** If the output body does NOT contain `release-note:`, repair it immediately:
   ```bash
   gh pr edit --repo GoogleCloudPlatform/magic-modules --body-file /tmp/pr_body.txt
   ```
4. Share the confirmed PR URL with the user and complete the execution.
