---
name: schema-review
description: "Skill for reviewing Magic Modules schemas across Enum/String modeling and zero-value clearing."
---

# `schema-review`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your required task.

## Prerequisites
* You must be in the `magic-modules` root directory.
* You must be reviewing or adding a field or resource schema.

## Execution Steps

1. **Enum vs. String & Descriptions:** For fields backed by an API enum, apply [.agents/knowledge/field/enums-vs-strings.md](../../../knowledge/field/enums-vs-strings.md) — verify configurable `Enum` fields omit `*_UNSPECIFIED` (`output: true` fields are exempt) and `type: Enum` descriptions do not manually duplicate `Possible values:`.
2. **Zero-Value Clearing & Empty Blocks:** For mutable primitives, `default_from_api` + zero values, or empty `{}` blocks, apply [.agents/knowledge/field/zero-value-clearing.md](../../../knowledge/field/zero-value-clearing.md).

> [!NOTE]
> The judgment content lives in [.agents/knowledge/](../../../knowledge/index.md) — this skill references those entries and does not hold its own copy.

## Verification & Handoff
* Proceed with your field addition or PR review based on those entries.
