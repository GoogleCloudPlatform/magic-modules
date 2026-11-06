---
name: permadiff-decision-path
description: Choose the right schema attribute or custom hook to eliminate post-apply diffs without masking real drift.
topics: [field, custom-code]
task_types: [field-add, new-resource, bug-fix, test-fix, code-review]
source: docs/content/develop/diffs.md + docs/content/reference/field.md
status: approved
last_verified: 2026-11-06
---

# Permadiff Decision Path

When `terraform plan` is non-empty after `apply` or during `ImportStateVerify`, match the cause to the narrowest mechanism rather than reaching for `ignore_read: true`:

## 1. Choose the Right Mechanism

1. **Static, predictable API default** -> **`default_value: <val>`** (plan-time default without marking `Computed`).
2. **Dynamic server-populated default returned on `GET`** -> **`default_from_api: true`** (`Optional` + `Computed`).
   - Set **only** when omitting the field on `Create` causes `GET` to return a non-empty value on the wire. If API prose mentions a conceptual default while `GET` omits unset values, keep the field pure `Optional` so users can clear it in HCL.
3. **API normalizes formatting (case, trailing slash, self-link vs. short name, duration/JSON)** -> **`diff_suppress_func` or `custom_flatten`**.
   - Prefer `tpgresource` helpers (`CaseDiffSuppress`, `CompareSelfLinkOrResourceName`, `DurationDiffSuppress`, `EmptyOrDefaultStringSuppress`). Use `custom_flatten` when canonicalizing on read also keeps `ImportStateVerify` clean.
4. **Secret, password, or one-time token omitted by `GET`** -> **`ignore_read: true` + `sensitive: true`** (plus `write_only: true` when exposing ephemeral `<field>_wo` attributes) and `ImportStateVerifyIgnore`.
   - Keep `ignore_read: true` and MMv1 `write_only: true` distinct: `addWriteOnlyFields` in `mmv1/api/resource.go` generates `<field>_wo` + `<field>_wo_version` while setting `Sensitive: true` and `IgnoreRead: false` on the base field, so a base field omitted by `GET` still requires `ignore_read: true` or a state-preserving flattener.
   - When every child inside a `NestedObject` sets `ignore_read: true` and `GET` omits the enclosing object, set `ignore_read: true` on the parent `NestedObject` too (such as `apiKey` in `mmv1/products/agentidentity/AuthProvider.yaml`).
5. **Action/trigger inputs or mutually-exclusive `oneof` update transitions** -> **`custom_code.update_encoder`, `custom_diff`, or targeted `DiffSuppressFunc`**. Keep working transition hooks rather than replacing them with `ignore_read: true`.
6. **Purely server-generated metadata (`etag`, `uid`, `create_time`)** -> **`output: true`**.

## 2. Do NOT Use For

- **Never** set `ignore_read: true` or `output: true` on readable user-configurable fields to hide a broken expander/flattener, missing `updateMask` entry, or zero-value clearing bug.
- **Never** branch on `d.GetOk` in `Read`/flatteners on `ForceNew` fields to hide server defaults—`terraform import` populates state and triggers `ForceNew` recreation on the next plan. Use `DiffSuppressFunc` (such as `tpgresource.EmptyOrDefaultStringSuppress(...)`) or a computed `effective_*` attribute.
- **Never** return `true` unconditionally in `DiffSuppressFunc`; always compare normalized `old` and `new`.
