---
name: prepare-pr-link
description: "Prepares a feature branch, stages and commits changes, pushes to the user's personal fork, and generates a validated pre-filled GitHub comparison link for the user to review and submit a Pull Request."
---

# `prepare-pr-link`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

## Pre-PR Guardrails

Run from the `magic-modules` root with `upstream` (`GoogleCloudPlatform/magic-modules`) and a personal fork remote (`origin`) configured. Verify before pushing:

1. **Single self-contained change:** One logical change per PR (one resource per PR; split bugfixes and new fields into separate PRs).
2. **PR title & commit subject length:** Target `<= 60` chars; hard limit `63` chars (enforced by `.agents/scripts/build_pr_link.py`, keeping squash-merge ` (#12345)` under GitHub's 72-char ceiling). Keep commit subjects `<= 60` chars too. Format: `<product/area>: <imperative verb phrase>`. Drop redundant `google_<service>_` prefixes after `<service>:` and summarize the theme instead of listing every field or check:
   - Bad (71): `compute: fix alias_ip_range netmask handling on google_compute_instance` → Good (56): `compute: fix alias_ip_range netmask handling on instance`
   - Bad (88): `mmv1: validate default_value types, enum duplicates/membership, and required constraints` → Good (52): `mmv1: validate default_value, enum dupes, & required`
3. **No downstream artifacts & clean workspace:** Do not commit generated downstream provider code into `magic-modules`. Run `git status --porcelain` before pushing to confirm no untracked temporary files remain.
4. **Add commits on updates:** Once a branch is pushed, always add new commits for follow-up changes rather than amending (`git commit --amend`), rebasing, or force-pushing (`-f`, `--force`, `--force-with-lease`) unless instructed otherwise.

## Execution Steps

### 1. Sync and Create Feature Branch

Use a clean `<short-descriptive-branch-name>` matching the change (never push internal worktree or task branch names; rename with `git branch -m "$BRANCH"` if needed). When updating an already-pushed PR branch, stay on that branch instead of recreating it:

```bash
UPSTREAM_REMOTE=$(git remote -v | grep -i "GoogleCloudPlatform/magic-modules" | head -n 1 | awk '{print $1}')
BASE_BRANCH="main" # or target major release branch, e.g. FEATURE-BRANCH-major-release-8.0.0
git fetch "${UPSTREAM_REMOTE:-upstream}" "$BASE_BRANCH"

BRANCH="<short-descriptive-branch-name>" # e.g. deprecate-beyondcorp-app
git checkout -b "$BRANCH" "${UPSTREAM_REMOTE:-upstream}/$BASE_BRANCH"
```

### 2. Stage, Commit, and Push

```bash
git add mmv1/products/<product>/ # or other relevant files
git commit -m "<product/area>: <concise description <= 60 chars>"

FORK_REMOTE=$(git remote -v | grep -v -i "GoogleCloudPlatform/magic-modules" | head -n 1 | awk '{print $1}')
git push -u "${FORK_REMOTE:-origin}" "$BRANCH"
```

### 3. Write the PR Description & Release Notes

- **Default to plain, brief prose:** Write for a human reviewer—1–2 short paragraphs (or a short flat bullet list when listing a few distinct changes) explaining the problem/why and what changed. Skip headings, tables, and boilerplate verification checklists by default.
- **Scale up only when warranted:** For genuinely complex or multi-part changes, add short sections for non-obvious design choices, edge cases, or trade-offs. Use tables **sparingly**—only when actually comparing discrete cases (e.g., before/after behavior across multiple rules), never as a default summary of files changed.
- **Avoid AI scaffolding:** Skip boilerplate headings (`## Summary`, `### Verification`, etc.), file-by-file diff narration, internal eval/agent jargon, and leftover PR template HTML comments (`<!-- ... -->`). Mention local testing only when it adds signal beyond CI (e.g., `TF_ACC` runs or repro steps).
- **Links & release notes:** Put `Needs <PR URL>` at the top if blocked. Add standalone `Fixes <github_issue_url>` / `Fixes b/<bug_id>` (or `Part of ...`) right before the release-note block—when working from a `b/` bug, link any matching GitHub issue too *(exception: keep internal tracker IDs, internal links, and unannounced dates out of GA promotions)*. End with a release-note block per [docs/content/code-review/release-notes.md](../../../../docs/content/code-review/release-notes.md) (syntax validated by `build_pr_link.py`):
  - Valid types (10): `enhancement`, `bug`, `note`, `none`, `new-resource`, `new-list-resource`, `new-datasource`, `new-action`, `deprecation`, `breaking-change`. Single line per block (`none` empty; `enhancement`/`bug` use `<product>: <past-tense description>`; `new-resource`/`new-datasource` use `` `google_<name>` ``).

Sample PR body:

~~~markdown
`RegionalEndpoint` GET responses return an empty `subnetwork` field on legacy
endpoints, causing a perpetual diff in shared customer modules. Add a diff
suppress function on `subnetwork` when the API returns empty.

Fixes https://github.com/hashicorp/terraform-provider-google/issues/12345

```release-note:bug
compute: fixed perpetual diff on `subnetwork` in `google_compute_regional_endpoint`
```
~~~

### 4. Generate Pre-Filled PR Link (Mandatory Script)

Always run `.agents/scripts/build_pr_link.py` to validate the title, release-note block, branch, and URL length before emitting the link. Do not improvise ad-hoc `urllib.parse` one-liners.

```bash
python3 .agents/scripts/build_pr_link.py \
  --title "<product/area>: <concise title <= 60 chars>" \
  --markdown \
  --body-file - <<'EOF'
<PR_BODY>
EOF
```

- Pass `--base <branch>` for non-`main` targets (e.g., `FEATURE-BRANCH-major-release-8.0.0`). `--branch` and `--user` default to the current git branch and fork remote.
- Fix any reported errors and review stderr warnings before sharing. Present the `--markdown` output directly in chat (`👉 **[Create Pull Request on GitHub](<URL>)**`) and let the user know the branch is pushed and ready to review.
