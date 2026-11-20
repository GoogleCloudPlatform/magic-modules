---
name: prepare-pr-link
description: "Prepares a feature branch, stages and commits changes, pushes to the user's personal fork, and generates a validated pre-filled GitHub comparison link for the user to review and submit a Pull Request."
---

# `prepare-pr-link`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

This skill provides step-by-step instructions for preparing, committing, pushing a feature branch to a personal fork, and generating a validated pre-filled GitHub Pull Request comparison link for the user to review and submit.

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
2. **Strict PR Title & Commit Subject Length Limit:**
   - Target **`<= 60` characters**; **hard limit `63` characters** (enforced by `.agents/scripts/build_pr_link.py`). Appending ` (#12345)` on squash-merge must stay within GitHub's 72-char commit-subject ceiling and avoid 70-char auto-truncation.
   - Keep git commit subjects **`<= 60` characters** as well.
   - Format: `<product/area>: <imperative verb phrase>`.
   - **Shorten titles by dropping redundant `google_<service>_` prefixes** after `<service>:` and **summarizing the theme** instead of listing every field or check:
     - Bad (104): `networkservices: add regex matching and URL rewrite fields to google_network_services_edge_cache_service`
     - Good (42): `networkservices: add regex matching fields`
     - Bad (88): `mmv1: validate default_value types, enum duplicates/membership, and required constraints`
     - Good (52): `mmv1: validate default_value, enum dupes, & required`
     - Bad (81): `agents: add endpoint, test result, and destroy error summaries to tf_debug_parser`
     - Good (31): `agents: enhance tf_debug_parser`
     - Good resource shortening: `compute: fix alias_ip_range netmask handling on instance` (56), `agentidentity: add REP support to auth_provider` (47).
3. **No Downstream Artifacts in Magic Modules:**
   - Do NOT commit generated downstream provider code into `magic-modules`.
4. **Workspace Cleanup:** Run `git status --porcelain` and ensure no untracked temporary test files exist before pushing.
5. **Prefer Adding Commits on Subsequent Updates:**
   - Once a branch has been pushed, always add new commits for follow-up changes rather than amending (`git commit --amend`), rebasing, or force-pushing (`-f`, `--force`, `--force-with-lease`) unless explicitly instructed otherwise.

---

## Execution Steps

### 1. Sync and Create Feature Branch

Ensure the branch has a clean, `<short-descriptive-branch-name>` matching the change (never push internal worktree or task branch names; rename with `git branch -m "$BRANCH"` if already committed on a worktree branch). If updating an already-pushed PR branch, check out or stay on that existing branch rather than recreating it from `$UPSTREAM_REMOTE/$BASE_BRANCH`:

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
git commit -m "<product>: <concise description <= 60 chars>"
```

### 3. Push to Personal Fork

```bash
FORK_REMOTE=$(git remote -v | grep -v -i "GoogleCloudPlatform/magic-modules" | head -n 1 | awk '{print $1}')
FORK_REMOTE="${FORK_REMOTE:-origin}"

git push -u "$FORK_REMOTE" "$BRANCH"
```

---

### 4. Write the PR Description & Release Notes

Write for a human reviewer, not like AI output. Lead with the problem and why the change is needed, orient the reviewer quickly, and keep it succinct.

#### Scale Structure to Complexity

| Complexity | Structure |
| :--- | :--- |
| **Simple** (single bugfix, field add, small doc/skill/test change) | **1–2 short prose paragraphs** (problem/why → what changed). No headings, no tables. |
| **Medium** (multi-check tool, CI, or validator change with edge cases) | **1–2 sentence problem statement** → optional compact **Before/After** or **Check/Effect** table (only when comparing discrete items) → `### Notable implementation details` (2–3 concise bullets on non-obvious choices or trade-offs) → 1 concise verification sentence if helpful. |
| **Complex** (multi-lifecycle or high-blast-radius fix) | Brief opening framing the problem and why the single-PR scope is warranted → logical cognitive-flow sections (e.g. why this approach is needed → what changed) → repro/verification context if relevant. |

#### Avoid

- **AI boilerplate headings:** `## Summary`, `### Summary of Changes`, `## Motivation`, `## Motivation & Proof`, or bullet-dump `### Verification` logs.
- **File-by-file diff narration:** Group by behavior and reason; reviewers already see the file list in GitHub's Files Changed tab.
- **Internal eval/session jargon:** Never include eval case IDs, agent/swarm roles, or local worktree paths in public PR descriptions.
- **Leftover template comments:** Strip HTML comments (`<!-- ... -->`) from `.github/PULL_REQUEST_TEMPLATE.md`.

#### Dependencies, Issue Links & Release Notes

1. **Blocking dependencies:** Put `Needs <PR URL>` at the very top if this PR depends on another PR merging first.
2. **Issue links:** Add standalone `Fixes <github_issue_url>` and/or `Fixes b/<bug_id>` (or `Part of ...`) right before the release-note block. When working from an internal `b/` bug, check whether a matching GitHub issue exists and link both. *(Exception: when promoting to GA, keep internal tracker IDs, internal links, and unannounced dates out of commits, the PR, and release notes.)*
3. **Release-note block:** Required at the bottom. See [docs/content/code-review/release-notes.md](../../../../docs/content/code-review/release-notes.md).
   - **Valid types (10):** `enhancement`, `bug`, `note`, `none`, `new-resource`, `new-list-resource`, `new-datasource`, `new-action`, `deprecation`, `breaking-change`.
   - **Single line per block:** Never put multiple lines inside one block; use multiple blocks if a PR has multiple changes.
   - `release-note:none` must be empty. `enhancement`/`bug` must start with `<product>: <past-tense description>`. `new-resource`/`new-datasource` must be `` `google_<name>` `` (with optional ` (beta)`/` (ga)`).

#### Sample PR Body (Simple)

~~~markdown
`RegionalEndpoint` GET responses return an empty `subnetwork` field on legacy
endpoints, causing a perpetual diff in shared customer modules. Add a diff
suppress function on `subnetwork` when the API returns empty.

Fixes https://github.com/hashicorp/terraform-provider-google/issues/12345

```release-note:bug
compute: fixed perpetual diff on `subnetwork` in `google_compute_regional_endpoint`
```
~~~

#### Sample PR Body (Medium)

~~~markdown
`mmv1` schema validation missed several YAML mistakes that only surfaced during
downstream provider compilation. Add static checks during compilation so bad
schemas fail fast:

| Check | Behavior |
| :--- | :--- |
| `default_value` type | Errors when the value does not match the property type |
| `enum_values` | Rejects duplicates and defaults missing from the enum |

### Notable implementation details

- Checks run inside `Property.Validate` after defaults are inherited so overrides are evaluated accurately.
- Existing legacy schemas with known exemptions are skipped via a minimal allowlist.

Verified with `go test ./...` under `mmv1/`.

```release-note:none
```
~~~

---

### 5. Generate Pre-Filled Web URL (Mandatory Script)

Always run `.agents/scripts/build_pr_link.py` to validate the title, release-note block, branch, and URL length before emitting the link. **Do not** improvise ad-hoc `python3 -c 'import urllib.parse...'` one-liners.

```bash
python3 .agents/scripts/build_pr_link.py \
  --title "<product>: <concise title <= 60 chars>" \
  --markdown \
  --body-file - <<'EOF'
<PR_BODY>
EOF
```

- Pass `--base <branch>` when targeting a non-`main` branch (for example, `FEATURE-BRANCH-major-release-8.0.0`).
- `--branch` and `--user` default to the current git branch and non-upstream GitHub fork remote; pass them explicitly if needed.
- If the script exits non-zero, fix every reported error (especially title length > 63 chars or release-note format errors) and re-run it. Review any soft warnings printed to stderr.

Present the script's `--markdown` output directly in chat:
👉 **[Create Pull Request on GitHub](<URL>)**

Inform the user that the branch has been pushed to their fork and they can click the link to review the diff and submit the Pull Request.
