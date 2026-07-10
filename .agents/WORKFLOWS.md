# Workflows Menu

This document lists the available high-level workflows in this repository and the rules governing them.
Match the workflow to the task; the Default Workflow is the fallback for tasks that don't fit a more
specific one.

## Rules of Engagement

1.  **Source of Truth:** `magic-modules` contains YAML schema definitions in `mmv1/` and templates. Code is generated into downstream providers.
2.  **Use Skills:** Discover and leverage `.agents/skills/` to perform complex tasks.
3.  **Verify before a PR:** generate, build, and run the tests relevant to the change. A task is not done because it compiles.
4.  **Never weaken the baseline:** no disabling or skipping tests, and no test-dodging behavior flags (`ignore_read`, `default_from_api`, `ImportStateVerifyIgnore`) without an adjacent comment justifying the API behavior that requires them.
5.  **PR descriptions are brief:** what changed and why, in a few sentences.

## Available Workflows:

*   **Default Workflow** (`.agents/skills/workflows/default/SKILL.md`): For tasks that do not involve creating a new resource (fallback for general tasks).
*   **New Resource Workflow** (`.agents/skills/workflows/new_resource/SKILL.md`): Specifically for creating a new resource, supporting both autogen and manual generation.
*   **Add List Resource Workflow** (`.agents/skills/workflows/add_list_resource/SKILL.md`): Opts one product's eligible MMv1 resources into list-resource generation by setting `generate_list_resource: true`, validates locally, and opens a PR.
*   *(Future workflows can be added here)*
