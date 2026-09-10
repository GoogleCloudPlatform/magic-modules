---
name: schema-review
description: "Skill for reviewing Magic Modules schemas. Currently covers the Enum vs. String decision via the knowledge base."
---

# `schema-review`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

## Prerequisites
* You must be in the `magic-modules` root directory.
* You must be reviewing or adding a field.

## Execution Steps

### 1. Enum vs. String Trade-off
When adding a field that is defined as an Enum in the API, you must decide between `Enum` and `String` in Magic Modules. Read the knowledge entry [.agents/knowledge/field/enums-vs-strings.md](../../../knowledge/field/enums-vs-strings.md) and apply its trade-off.

> [!NOTE]
> This skill will be expanded over time with more schema review checklist items. The judgment content lives in [.agents/knowledge/](../../../knowledge/index.md) — this skill references it and does not hold its own copy.

### 2. Verification & Handoff

* Proceed with your field addition or PR review based on the decision.
