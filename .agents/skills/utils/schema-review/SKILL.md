---
name: schema-review
description: "Stub skill for reviewing Magic Modules schemas. Currently covers Enum vs. String trade-offs."
---

# `schema-review`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

## Prerequisites
* You must be in the `magic-modules` root directory.
* You must be reviewing or adding a field.

## Execution Steps

### 1. Enum vs. String Trade-off
When adding a field that is defined as an Enum in the API, you must decide between `Enum` and `String` in Magic Modules.

#### The Trade-off:
*   **Prefer `Enum`** when you want strict, plan-time validation of values to fail fast.
*   **Prefer `String`** (with allowed values documented) when the API is dynamic or likely to add values. This prevents Terraform from crashing with a validation error for users when the API updates *before* the provider catches up.

> [!NOTE]
> This skill is a **stub** and will be expanded over time with more schema review checklist items (such as `api_name` overrides, array flattening, output-only fields, etc.).

### 2. Verification & Handoff
Instructions on how the agent should verify the command succeeded, and what workflow or rule it should return to next.

* Proceed with your field addition or PR review based on this stylistic choice.
* Reference the `modernization_roadmap.md` for upcoming expansions to this checklist.
