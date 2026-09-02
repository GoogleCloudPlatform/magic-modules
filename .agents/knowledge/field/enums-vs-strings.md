---
name: enums-vs-strings
description: Model an API enum as Enum (strict, plan-time) or String (forward-compatible) - the deliberate trade-off.
topics: [field]
task_types: [field-add, new-resource]
source: docs/content/best-practices/validation.md + the schema-review skill's original trade-off note
status: draft
last_verified: 2026-07-09
---

# Enum vs. String

**The trade-off (decide deliberately, per field):**
- **`Enum`** when the value set is exhaustive for a clearly-defined domain and new values are extremely
  unlikely. Buys strict plan-time validation — users fail fast instead of mid-apply.
- **`String`** (with allowed values named in the description and a link to API docs) when the API is
  likely to add values. **Why:** an Enum's value list lives in the provider; when the API adds a value,
  every user needs a provider upgrade before they can use it — and imports of resources using the new
  value break. A String stays forward-compatible.

```yaml
- name: 'severity'
  type: Enum
  enum_values:
    - 'LOW'
    - 'MEDIUM'
    - 'HIGH'
```

**Rules for Enum:**
- Omit the `FIELD_NAME_UNSPECIFIED` value from `enum_values` — it is the proto convention's first enum
  value, representing 0/unset, not a real choice.
- Enums validate against the list automatically; a custom `validation` **overrides** that default — if you
  add one, it must re-verify the enum values itself.

**Rule for String-instead-of-enum:** list the allowed values in the field description with a link to the
API docs — otherwise users only discover an invalid value when the API rejects it mid-apply.
