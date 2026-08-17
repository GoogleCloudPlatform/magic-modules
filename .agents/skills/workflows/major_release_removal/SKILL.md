---
name: major-release-removal-workflow
description: "Workflow for executing resource and field removals on a major release feature branch."
---

# `major-release-removal-workflow`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.

This workflow governs removing deprecated resources or fields on a major release feature branch (e.g., `FEATURE-BRANCH-major-release-8.0.0`).

---

## Prerequisites

- You must be operating in the `magic-modules` repository root directory.
- You must know the target major release branch (e.g., `FEATURE-BRANCH-major-release-8.0.0`).
- You must know the target product, resource, and field (or full resource) to remove.

---

## Execution Steps

### 1. Pre-Flight Audit via `removal-auditor` Subagent

Launch the `removal-auditor` subagent (`.agents/agents/removal-auditor/`) with the target product, resource, field, and major release branch.

The subagent will inspect `upstream/main` (deprecation & replacement checks), `upstream/<major_release_branch>` (sync status), and all repo dependencies, returning a **Removal Audit & Blast-Radius Report**.

#### Gating Checks:
- **If deprecation is missing on `main`**: Stop this workflow. Invoke the [`deprecate-resource-or-field-workflow`](../deprecate_resource_or_field/SKILL.md) on `main` first.
- **If the release branch is not synced with `main`**: Stop this workflow. Invoke [`sync-main-to-major-release-branch`](file:///usr/local/google/home/camthornton/.gemini/config/skills/sync-main-to-major-release-branch/SKILL.md) first.

---

### 2. Plan Presentation & Review Checkpoint

Present the subagent's audit report and proposed removal plan to the user as an artifact (schema deletions, template/sample/test cleanup, TGC converters, and proposed upgrade guide text).

**Do NOT proceed to code modifications until the user approves the plan.**

---

### 3. Checkout Major Release Branch

Create a local working branch based on `upstream/<major_release_branch>` (e.g. `upstream/FEATURE-BRANCH-major-release-8.0.0`). Ensure downstream provider workspaces are also checked out on the matching major release branch.

---

### 4. Code Removal & Cleanup

Execute removals across the exact files identified in the approved Removal Audit Report:
- **Schema**: Delete the property or resource YAML in `mmv1/products/<product>/...` (or Go schema / `ResourceMap` in `mmv1/third_party/terraform/services/<product>/...`).
- **Custom Templates & Hooks**: Remove associated templates in `mmv1/templates/terraform/` (expanders, flatteners, hooks, constants, state migrations).
- **Samples & Examples**: Remove obsolete sample `.tf.tmpl` files and test configs.
- **Acceptance Tests**: Remove deleted fields from test configs and `ImportStateVerifyIgnore` slices; delete test files for removed resources.
- **TGC Converters**: Remove converter mappings and IAM registrations from `mmv1/third_party/tgc/resource_converters.go.tmpl`.
- **Documentation**: Remove handwritten documentation markdown if applicable.

---

### 5. Update Version Upgrade Guide

Add an entry to `mmv1/third_party/terraform/website/docs/guides/version_<MAJOR>_upgrade.html.markdown` following existing entries in that file and guidance in `docs/content/breaking-changes/make-a-breaking-change.md`:
- **Resource Removal**: `## Resource: google_<resource> is now removed` with migration advice.
- **Field Removal**: Under `## Resource: google_<resource>`, add `### <field> is now removed` explaining the removal and replacement argument.

---

### 6. Pre-Gen Checks & Code Generation

1. Run fast static pre-gen checks: [run-pre-gen-checks](.agents/skills/utils/run-pre-gen-checks/SKILL.md).
2. Generate provider code: [generate-provider](.agents/skills/operations/generate-provider/SKILL.md).
3. Verify compilation: Run `make build` in downstream provider repositories.

---

### 7. Verification Testing

Invoke [qa-test-runner](.agents/skills/operations/qa-test-runner/SKILL.md) to run acceptance tests for remaining or adjacent resources to ensure no regressions.

---

### 8. PR Creation & Release Note

Open a PR targeting the major release branch (e.g., `FEATURE-BRANCH-major-release-8.0.0`):
- Title: `<product>: remove deprecated <field_name|resource_name> for <MAJOR_VERSION>`
- Body:
  ```markdown
  ```release-note:breaking-change
  <product>: removed deprecated `<field_name>` from `google_<resource_name>`
  ```
  *(or `<product>: removed deprecated `google_<resource_name>` resource`)*
  ```
