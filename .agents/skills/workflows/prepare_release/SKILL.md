---
name: prepare-release-workflow
description: "Prepare and cut weekly releases for both terraform-provider-google (TPG) and terraform-provider-google-beta (TPGB)."
---

# `prepare-release-workflow`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.
> This skill is designed to be completely self-contained and unambiguous for a fresh agent without prior context.

This workflow prepares and cuts the weekly releases for both the `terraform-provider-google` (TPG) and `terraform-provider-google-beta` (TPGB) repositories.

## Prerequisites

- You must be operating relative to the `magic-modules` workspace.
- You must locate the downstream target provider repositories.

## Step 0: Locate Repositories & Set Workspace Context

Before creating any branches or running scripts, you must establish the path of magic-modules and the downstream providers:
1. **Locate Magic Modules (MM_DIR):** The magic-modules repository is the repository you are running this agent from. Record its absolute path as `MM_DIR`.
2. **Locate TPG & TPGB Repositories:**
   - TPG (`terraform-provider-google`) and TPGB (`terraform-provider-google-beta`) are the downstream target provider repositories.
   - Try to find TPG and TPGB paths:
     - Check active workspaces or GOPATH (`go env GOPATH`/src/github.com/hashicorp/terraform-provider-google[-beta]).
     - Check sibling directories of `MM_DIR`.
   - Verify by checking if `git status` runs successfully inside these directories.
   - If not found or if there is any ambiguity, ask the user to verify or supply the absolute paths to both repositories before proceeding.
3. **Execution Mode:** You will prepare releases for TPG first, then TPGB. Throughout the run, always execute commands inside the correct provider repository by setting the command's current working directory (`Cwd`) to that repository's absolute path.
4. **Running Helper Scripts:** Always execute scripts located in magic-modules using their absolute paths, e.g., `bash $MM_DIR/.agents/scripts/get_nightly_sha.sh ...` with `Cwd` set to TPG or TPGB.

## Step 1: Determine Release Candidate, Remote, & Semver

For the current target repository (TPG or TPGB), determine exactly what commit to release, what versions to use, and where to push:

1. **Determine Upstream & Target Remotes:**
   - Run `git remote -v` in the target repository to inspect mappings.
   - Identify the remote pointing to the official upstream repository (`hashicorp/terraform-provider-google` or `hashicorp/terraform-provider-google-beta`) (this is the `UPSTREAM_REMOTE`).
   - The `TARGET_REMOTE` defaults to the `UPSTREAM_REMOTE`. If the user explicitly requested a different remote (e.g., their fork), use that as the `TARGET_REMOTE` instead.
2. **Check Candidate SHA & Nightly Status:**
   - By default, the weekly release uses the Thursday night cut. Check if the user explicitly specified a different day of the week to take the cut from.
   - Run the nightly SHA script:
     `bash $MM_DIR/.agents/scripts/get_nightly_sha.sh <UPSTREAM_REMOTE> [--day-of-week <day>]` (setting `Cwd` to the target provider repository).
   - If it outputs any error `AGENT_INSTRUCTION` (e.g., missing tokens), stop and resolve it with the user. If `TEAMCITY_TOKEN` is missing, present clear step-by-step instructions to the user: 1) Log into https://hashicorp.teamcity.com/, 2) Open Profile -> Access Tokens, 3) Click 'Create Access Token', name it, and copy it, 4) Export via `export TEAMCITY_TOKEN="<token>"`.
3. **Calculate Semver Versions:**
   - Run the release versions script:
     `bash $MM_DIR/.agents/scripts/get_release_versions.sh <UPSTREAM_REMOTE>` (setting `Cwd` to the target provider repository).
   - Extract the `previous_version` and `next_version` from the JSON output.
4. **Check GitHub Tokens:** Ensure `$JET_GITHUB_TOKEN` or `$GITHUB_TOKEN` is set and valid.

### Mandatory Confirmation Checkpoint
Present the standard confirmation checkpoint directly in your normal chat response. DO NOT use the `ask_question` tool for this. Instead, output the following text and wait for the user to reply in the chat:
```markdown
### 🛑 Confirm Release Candidate & Configuration for <Provider_Name>
* **Previous Release:** `v<previous_release_version>`
* **New Release:** `v<new_release_version>`
* **Target Commit SHA:** [`<short_sha>`](https://github.com/hashicorp/<provider_name>/commit/<long_sha>) ("<commit_message>") (`<long_sha>`)
* **TeamCity Build:** [View Execution](<web_url>) (Finished: `<local_time>`)
* **Target Remote:** `<remote_name>` (`<remote_url>`)

*(Note: Always display a prominent warning advising the user to check the TeamCity execution link and manually verify the build status before proceeding, as individual test failures do not necessarily indicate release quality but a manual check is generally good. Do not state that the build failed or mention build failures/status explicitly. The warning should read: "The TeamCity build finished. Please check the TeamCity Execution Link and manually verify the build status before proceeding.")*

Would you like to proceed with creating base release branch `release-<new_release_version>` off of commit `<short_sha>` and pushing it to `<remote_name>`?
```
Stop and wait for the user to reply to confirm and proceed with the release cut.

## Step 2: Cut Release Branch & Generate Changelog

Once the user confirms the details in Step 1, run the release cutting script with `Cwd` set to the target provider repository:
```bash
bash $MM_DIR/.agents/scripts/cut_release_branch.sh <COMMIT_SHA> <PREVIOUS_RELEASE_VERSION> <RELEASE_VERSION> <TARGET_REMOTE> <UPSTREAM_REMOTE>
```
- **Crash Recovery & Unmerged Changelog Handling:** 
  - The `cut_release_branch.sh` script automatically detects if `CHANGELOG.md` on `main` is missing the previous release notes. If missing, it recovers and prepends the previous release section from the remote changelog branch.
  - If the script crashes midway, inspect the error, help the user resolve the issue, and manually pick up the missing steps rather than starting from scratch.
- If the script outputs an `AGENT_INSTRUCTION` (e.g., missing tokens or parameter validation failures), stop and ask the user to provide or fix them.
- **User Checkpoint:** Once the script succeeds, it will have automatically committed the raw changes, pushed the PR branch, and created the PR. 
  Provide the user with the PR link outputted by the script: "I have pushed the branch and created the PR. You can view the PR with the raw changelog diff here: <PR_LINK>". Then ask: "Shall I proceed to audit the changelog?"

## Step 3: The Release Notes Audit & PR

After the user tells you to proceed to the audit, rigorously audit the new notes in `CHANGELOG.md` against the following rules (operating inside the target provider repository):

1. **Cross-Provider Filtering (`changelog-gen` sync artifact removal):**
   - **For TPG (GA):** Delete any entry ending in `(beta)` or describing a Beta-only addition. Strip the `(ga)` suffix from valid GA entries.
   - **For TPGB (Beta):** Delete any entry ending in `(ga)` or describing a GA-only addition. Strip the `(beta)` suffix from valid beta entries (e.g. `added foo (beta)` -> `added foo`).
2. **Internal Tooling & Infra Exclusion:**
   - Remove or omit entries for commits/PRs that exclusively touch repository tooling, internal scripts, or CI infrastructure (such as `.agents/`, `.github/`, `.ci/`, unit test helpers, or build scripts).
3. **Net-Zero Change Check (Same-Cycle Introduced & Fixed Bugs):**
   - Check candidate PRs/commits merged during this release window (since `$PREVIOUS_RELEASE_VERSION`). If a PR fixes a bug or panic introduced in an earlier PR within *this same release cycle*, its net effect to released users is zero.
   - *Mandatory Safeguard*: NEVER silently delete potential net-zero entries without user confirmation. Delegate diff/PR inspection to compile candidate net-zero notes, present your findings with PR links and rationale to the user, and request explicit approval before removing them from `CHANGELOG.md`.
4. **One Line Per Block:** Ensure no single release note block spans multiple lines.
5. **Past Tense Verbs:** Every note must begin with a past-tense verb (`added`, `fixed`, `resolved`, `updated`, `rejected`, `migrated`). Never use present/future.
6. **User-Focused & Backticked HCL Names:**
   - Inspect exact git commit diffs (`git show <SHA>`) when in doubt about HCL attribute names.
   - Surround all Terraform resource names, data source names, and attribute/block names in backticks (`google_compute_instance`, `workload_identity_config`).
   - Replace internal struct names with exact backticked HCL blocks.
   - Replace informal descriptions with exact backticked resource names.
7. **Correct Product Prefixes:**
   - Prefix must be the product/service folder name or API subdomain followed by a colon and a lowercase letter (`compute: added...`). Do not backtick the prefix.
   - For provider-wide cross-product changes, use `provider:`.
8. **Correct Tab Sectioning:**
   - `new-resource` -> `* **New Resource:** \`RESOURCE_NAME\`` under `FEATURES:`.
   - `new-datasource` -> `* **New Data Source:** \`DATASOURCE_NAME\`` under `FEATURES:`.
   - `new-list-resource` -> `* **New List Resource:** \`LIST_RESOURCE_NAME\`` under `FEATURES:`.
   - `enhancement` -> `IMPROVEMENTS:`.
   - `bug` -> `BUG FIXES:`.
   - `none` / `changelog: no-release-note` -> Remove completely.
9. **Strict Alphabetical Sorting:** Order all notes strictly alphabetically within each section.

### User Checkpoint
Do not present the git diff directly in the chat. Instead, make the audited changes to `CHANGELOG.md` locally and ask the user to review the file and approve pushing the update to the PR.

### Final Push
- Once approved, commit and push the audited notes:
  ```bash
  git add CHANGELOG.md && git commit -m "changelog: apply style guide audit to $RELEASE_VERSION release notes"
  git push <TARGET_REMOTE> changelog-$RELEASE_VERSION
  ```
- Notify the user that the PR has been updated with the audited notes, and provide the PR link again for convenience.
- After completing the PR for TPG, check if the TPGB release has also been prepared. If not, proceed to Step 0/1 for TPGB.
