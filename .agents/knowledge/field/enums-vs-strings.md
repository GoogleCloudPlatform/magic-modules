---
name: enums-vs-strings
description: Model an API enum as Enum (strict, plan-time) or String (forward-compatible) - the deliberate trade-off.
topics: [field]
task_types: [field-add, new-resource]
source: docs/content/best-practices/validation.md + the schema-review skill's original trade-off note
status: approved
last_verified: 2026-11-06
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

**Rules for `Enum` (MMv1 YAML):**
- Omit `*_UNSPECIFIED` on configurable (`Optional` / `Required`) fields — it represents proto `0`/unset, not a user choice. **Exception:** `output: true` (read-only) fields may retain `*_UNSPECIFIED` because MMv1 omits `ValidateFunc` on `output: true` fields and APIs often return `*_UNSPECIFIED` while provisioning.
- Do **not** list `Possible values:` in `description:` on `Enum` (or `Array` of `Enum`) — MMv1 auto-appends `enum_values` to schema and website docs.
- Custom `validation` **overrides** the generated enum check and must re-validate allowed values itself.

**Rule for `String`-instead-of-enum (MMv1 YAML):** Name allowed values and link the API docs in `description:` (MMv1 does not auto-append values for `String`).

**Rule for Handwritten Go (`validation.StringInSlice`):** On `Optional: true, Computed: true` `TypeString` fields without diff suppression, omit `*_UNSPECIFIED` unless the API persists it as a distinct user-settable value—otherwise unset vs. `"FOO_UNSPECIFIED"` causes permadiffs. Read-only `Computed` fields omit `ValidateFunc`.
