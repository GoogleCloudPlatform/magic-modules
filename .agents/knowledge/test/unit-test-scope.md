---
name: unit-test-scope
description: When to write unit tests vs acceptance tests, and avoiding test bloat.
topics: [test]
task_types: [bug-fix, test-fix, field-add, new-resource]
source: authored
status: approved
last_verified: 2026-11-06
---

# Unit Test Scope & Bloat Prevention

Unit tests in `mmv1/third_party/terraform/services/*/*_test.go` and `*_internal_test.go` verify isolated Go logic. They should target reusable logic and systemic patterns, not one-off declarative schema definitions.

## When to Write Unit Tests

Write unit tests for custom helpers with branching, parsing, or non-trivial state logic:
- **Complex diff suppress functions** (`*DiffSuppress`): custom normalization, regex/string parsing, or multi-case equivalence (not shared `tpgresource.*` helpers or single-comparison wrappers).
- **Custom validation functions** (`ValidateFunc`) & **compiled regexes**: boundary inputs and domain-specific rules.
- **State upgraders & ID/URL parsers**: string/URL normalization or ID deconstruction.
- **Complex custom expanders/flatteners**: multi-branch data transformations with edge cases hard to hit in `TestAcc*`.

**Why:** Custom logic with edge cases can regress silently. See [Add unit tests (`docs/content/test/test.md`)](../../../docs/content/test/test.md#add-unit-tests).

## Do NOT Use Unit Tests For

- **Shared Helpers, Single-Condition Suppressors & CRUD Hooks**: Attaching `tpgresource.*` / `customdiff.*` helpers, single-condition `DiffSuppressFunc`s, or straightforward `encoder`/`decoder`/`update_encoder` hooks already covered by `TestAcc*`.
- **Concrete `*schema.ResourceDiff` Hooks**: `CustomizeDiff` functions typed directly on `*schema.ResourceDiff` (whose SDKv2 fields are unexported) when covered by `TestAcc*`, unless complex branching warrants extracting a `tpgresource.TerraformResourceDiff` helper for `ResourceDiffMock`.
- **Declarative Schema Properties**: Mocking `r.Validate(c)` to test `Required`, `Optional`, `AtLeastOneOf`, `ConflictsWith`, or `Default` (already tested by the Terraform SDK).
- **One-Off Typo Fixes & Ephemeral Repros**: Covered by `TestAcc*`; delete ad-hoc debug tests before opening a PR.
