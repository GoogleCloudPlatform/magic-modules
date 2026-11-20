---
name: unit-test-scope
description: When to write unit tests vs acceptance tests, and avoiding test bloat.
topics: [test]
task_types: [bug-fix, test-fix, field-add, new-resource]
source: authored
status: draft
last_verified: 2026-11-17
---

# Unit Test Scope & Bloat Prevention

Unit tests in `mmv1/third_party/terraform/services/*/*_test.go` and `*_internal_test.go` verify isolated Go logic. They should target reusable logic and systemic patterns, not one-off declarative schema definitions.

## When to Write Unit Tests

Write unit tests for custom functions with algorithmic logic, branching, or parsing:
- **Complex diff suppress functions** (`*DiffSuppress`): custom normalization, regex/string parsing, or multi-case equivalence (not shared `tpgresource.*` helpers or single-comparison wrappers).
- **Custom validation functions** (`ValidateFunc`): testing regex matches or domain-specific input rules.
- **State upgraders & parsers**: string parsing, URL parsing, or resource ID deconstruction.
- **Custom expanders/flatteners**: complex data transformations with branching.

**Why:** These functions contain custom execution logic that could regress if modified. For a reference implementation, see [Add unit tests (`docs/content/test/test.md`)](../../../docs/content/test/test.md#add-unit-tests).

## Do NOT Use Unit Tests For

- **Shared Helpers & `*schema.ResourceDiff` Hooks**: Attaching `tpgresource.*` / `customdiff.*` helpers, single-condition `DiffSuppressFunc`s, or `CustomizeDiff` hooks typed directly on `*schema.ResourceDiff` (whose SDKv2 fields are unexported) when already covered by `TestAcc*`.
  - **Why:** Acceptance tests exercise these end-to-end; extract a `tpgresource.TerraformResourceDiff` helper for `ResourceDiffMock` only when complex conditional branching warrants isolated testing.
- **Declarative Schema Properties**: Do NOT write unit tests that mock `r.Validate(c)` simply to check declarative schema properties (e.g. `Required`, `Optional`, `AtLeastOneOf`, `ConflictsWith`, `Default`).
  - **Why:** The Terraform SDK already tests its schema engine. Testing single declarative field permutations adds high-maintenance test bloat.
- **One-Off Typo Fixes**: Do NOT add persistent unit tests for one-off typos or copy-paste slice errors.
  - **Why:** These are covered by standard acceptance tests (`TestAcc*`) or verified ephemerally during triage reproduction.
- **Ephemeral Debug Tests**: Remove any ad-hoc unit tests created solely to reproduce an issue locally before finalizing a PR.
