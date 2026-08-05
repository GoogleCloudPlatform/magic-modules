---
name: tgc-create-pr-skill
description: "Create a pull request for Magic Modules changes to support a TGC feature, adhering to upstream conventions, prefixing with tgc-revival, and formatting empty release notes."
---

# `tgc-create-pr-skill`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

## Prerequisites
* You must be working in the `magic-modules` root directory.
* You must have local changes inside `mmv1/` ready to be committed and pushed.
* You must have the `gh` CLI installed and authenticated on the host machine.

## Execution Steps

### 1. Verification of Staged Changes
Verify that only files inside the `mmv1/` folder are modified or staged. Ensure no scratch files, text files, or downstream TGC files are modified.

```bash
# Run in magic-modules root directory
git status
```

> [!IMPORTANT]
> If any files outside the `mmv1/` directory (e.g., downstream files, `.txt`, `.sh`, or `.py` scripts) are staged or tracked, you MUST reset them before committing.

### 2. Fetching Upstream Baseline
Ensure your local git repository is aware of the latest `upstream/main` changes before branching.

```bash
git fetch upstream main
```

### 3. Creating a New Branch Based on Upstream Main
Create and checkout a new local feature branch based on the fetched `upstream/main` (or `FETCH_HEAD`).

```bash
git checkout -b tgc-<feature-name> upstream/main
```

Verify that your uncommitted/staged `mmv1/` changes have successfully carried over to the new branch.

### 4. Committing the Changes
Stage and commit the `mmv1/` changes with a clear, descriptive commit message.

```bash
git add mmv1/...
git commit -m "[TGC] Support <feature_name> in <resource_name>"
```

### 5. Pushing the Branch
Push the local branch to your GitHub fork (`origin`).

```bash
git push -u origin tgc-<feature-name>
```

### 6. Verifying GitHub CLI Auth
Check that the GitHub CLI is logged in and active.

```bash
gh auth status
```

### 7. Creating the Pull Request
Create the pull request to the `GoogleCloudPlatform/magic-modules` upstream repository. Ensure the title is prefixed with `tgc-revival:` and the body includes the standard empty release note block.

```bash
gh pr create \
  --base main \
  --title "tgc-revival: Support <field/resource> in <ResourceName>" \
  --body "Support GKE <field/resource> in <ResourceName>

\`\`\`release-note:none
\`\`\`"
```

## Verification & Handoff
* Upon running the `gh pr create` command, copy the returned pull request URL and print it for the user.
* Update `task.md` to mark pull request creation and commit as fully completed.
