# Workflows Menu

This document lists the available high-level workflows in this repository and the rules governing them.
Match the workflow to the task; the Default Workflow is the fallback for tasks that don't fit a more
specific one.

## Rules of Engagement

1.  **Source of Truth:** `magic-modules` contains YAML schema definitions in `mmv1/` and templates. Code is generated into downstream providers.
2.  **Use Skills:** Discover and leverage `.agents/skills/` to perform complex tasks.
3.  **Verify before a PR:** generate, build, and run the tests relevant to the change. A task is not done because it compiles.
4.  **Never weaken baseline test coverage:** no disabling or skipping tests, and no test-dodging behavior flags (`ignore_read`, `ImportStateVerifyIgnore`) without an adjacent comment justifying the API behavior that requires them.
5.  **PR descriptions are brief:** what changed and why, in a few sentences.
6.  **GitHub Issue Label Routing:** When triaging a GitHub issue URL, inspect the GitHub labels (`gh issue view --json labels` or issue payload). Match against the **Issue Label Routing Matrix** below. If no matching label is present, inspect the issue title/body content; if still unmatched, fall back to the **Default Workflow**.

## Issue Label Routing Matrix:

| GitHub Label Pattern | Issue Category | Target Workflow Skill |
| :--- | :--- | :--- |
| `test-failure`, `test-failure-*` (e.g. `test-failure-0`, `test-failure-10`, `test-failure-50`, `test-failure-100`) | Acceptance Test Failure | `.agents/skills/workflows/test_fix/SKILL.md` |
| `new-resource` | New Resource Creation | `.agents/skills/workflows/new_resource/SKILL.md` |
| `list-resource` | List Resource Addition | `.agents/skills/workflows/add_list_resource/SKILL.md` |
| `enhancement` (field additions) | New Field / Property Addition | `.agents/skills/workflows/add_fields/SKILL.md` |
| `deprecation` | Deprecation Notice Addition | `.agents/skills/workflows/deprecate_resource_or_field/SKILL.md` |
| `breaking-change` | Major Release Removal / Breaking Change | `.agents/skills/workflows/major_release_removal/SKILL.md` |
| *No matching label / un-labeled* | General Modification / Bug Fix | Fallback to `.agents/skills/workflows/default/SKILL.md` (or inspect issue body) |

## Available Workflows:

*   **Default Workflow** (`.agents/skills/workflows/default/SKILL.md`): For tasks that do not involve creating a new resource (fallback for general tasks).
*   **New Resource Workflow** (`.agents/skills/workflows/new_resource/SKILL.md`): Specifically for creating a new resource, supporting both autogen and manual generation.
*   **Add Fields Workflow** (`.agents/skills/workflows/add_fields/SKILL.md`): Specifically for adding new fields to existing MMv1 or handwritten resources.
*   **Test Fix Workflow** (`.agents/skills/workflows/test_fix/SKILL.md`): Specifically for resolving failing acceptance tests from GitHub issues, direct prompts, or debug logs.
*   **Add List Resource Workflow** (`.agents/skills/workflows/add_list_resource/SKILL.md`): Opts one product's eligible MMv1 resources into list-resource generation by setting `generate_list_resource: true`, and validates locally.
*   **Prepare Release Workflow** (`.agents/skills/workflows/prepare_release/SKILL.md`): Prepares and cuts weekly releases for both `terraform-provider-google` (TPG) and `terraform-provider-google-beta` (TPGB) providers.
*   **Test Monitor Workflow** (`.agents/skills/workflows/test_monitor/SKILL.md`): For fetching, triaging, analyzing, and reporting on nightly acceptance test results across Beta and GA providers.
*   **Add IAM Support Workflow** (`.agents/skills/workflows/add_iam_resources/SKILL.md`): Adds IAM support to an existing MMv1 resource.
*   **Bug Fix Workflow** (`.agents/skills/workflows/bug_fix/SKILL.md`): For triaging, planning, fixing, and verifying reported provider bugs.
*   **Deprecate Resource or Field Workflow** (`.agents/skills/workflows/deprecate_resource_or_field/SKILL.md`): Specifically for deprecating existing resources or fields on `main` prior to a major release.
*   **Major Release Removal Workflow** (`.agents/skills/workflows/major_release_removal/SKILL.md`): Specifically for executing resource and field removals on a major release feature branch.
*   *(Future workflows can be added here)*

## Subagents:

*   **`test-fixer`** (`.agents/agents/test-fixer/`): Automates diagnosis, remediation in Magic Modules, provider generation, and re-testing for failing acceptance tests.
*   **`removal-auditor`** (`.agents/agents/removal-auditor/`): Audits deprecation status on `main`, sync status on the major release branch, and scans repository dependencies for resource and field removals.
