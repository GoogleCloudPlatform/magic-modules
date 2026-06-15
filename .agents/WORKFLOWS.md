# Workflows Menu

This document lists the available high-level workflows in this repository and the rules governing them. To use one, the user must explicitly invoke it.

## Rules of Engagement (Active in Workflow Mode)
1.  **No Blind Fixes:** Never modify YAML or Go files without presenting an implementation plan and receiving approval.
2.  **Source of Truth:** `magic-modules` contains YAML schema definitions in `mmv1/` and templates. Code is generated into downstream providers.
3.  **Use Skills:** Discover and leverage `.agents/skills/` to perform complex tasks.

## Available Workflows:

*   **Default Workflow** (`.agents/skills/workflows/default/SKILL.md`): For tasks that do not involve creating a new resource (fallback for general tasks).
*   **New Resource Workflow** (`.agents/skills/workflows/new_resource/SKILL.md`): Specifically for creating a new resource, supporting both autogen and manual generation.
*   **Add List Resource Workflow** (`.agents/skills/workflows/add_list_resource/SKILL.md`): Opts one product's eligible MMv1 resources into list-resource generation by setting `generate_list_resource: true`, validates locally, and opens a PR.
*   *(Future workflows can be added here)*
