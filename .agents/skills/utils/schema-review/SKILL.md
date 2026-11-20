---
name: schema-review
description: "Review Magic Modules field and resource schemas across Enum/String modeling and zero-value clearing."
---

# `schema-review`

When adding or reviewing Magic Modules schema fields, open and apply the matching knowledge entries (the rules live in `.agents/knowledge/` rather than in this skill):

1. **Enum vs. String & Descriptions:** Read [`.agents/knowledge/field/enums-vs-strings.md`](../../../knowledge/field/enums-vs-strings.md) whenever an API field is an enum.
2. **Zero-Value Clearing & Empty Blocks:** Read [`.agents/knowledge/field/zero-value-clearing.md`](../../../knowledge/field/zero-value-clearing.md) when handling `send_empty_value`, `default_from_api` with zero values, or `allow_empty_object`.
