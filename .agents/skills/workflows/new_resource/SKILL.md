---
name: new-resource-workflow
description: "Workflow specifically for creating a new resource, supporting both autogen and manual generation."
---

# `new-resource-workflow`

This document guides the agent through the task of creating a new Terraform resource in Magic Modules.

## Generation Options

There are two ways to generate the initial YAML definition for the resource:
1.  **Autogen Subagent (Default):** Use this method if an OpenAPI specification is available for the resource.
2.  **Manual Generation:** Create the YAML definition directly if no spec is available or if it requires custom crafting.

---

## Execution Steps

### 1. Prompt for OpenAPI Specification (Default Path)
*   Ask the user to provide the OpenAPI specification file (or its path) for the new resource.
*   **Note:** If no OpenAPI specification is available, state that you will fall back to the **Manual Path (Step 2B)** to draft the YAML directly, and ask for the user's confirmation to proceed.
*   If the user provides the spec, proceed to **Step 2A (Autogen Path)**.

### 2A. Autogen Path
*   **Action:** Invoke the `autogen` subagent to handle the generation and initial testing.
*   **Prompt:** "Use the OpenAPI spec to generate the new resource. Run tests and report back."
*   Wait for the subagent to return its report.
*   **Handoff:**
    - If the subagent reports that tests passed successfully, present the results to the user.
    - If the subagent reports test failures, enter the **Default Workflow** ([default/SKILL.md](../default/SKILL.md)) at **Step 5 (Fix)** to plan remediation.

### 2B. Manual Path
*   Consult `.agents/knowledge/index.md` for the topics the resource touches and open the relevant sources.
*   Follow the standard process to draft the YAML definition in `mmv1/products/...` based on API documentation and repository patterns.
*   **Handoff:** Once the YAML is drafted, enter the **Default Workflow** ([default/SKILL.md](../default/SKILL.md)) at **Step 3 (Generate)** to compile the provider and continue with testing.
