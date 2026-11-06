---
name: import-safe-flatteners
description: Write import-safe flatteners and handle server-populated defaults on ForceNew fields without breaking terraform import.
topics: [field, custom-code]
task_types: [field-add, new-resource, bug-fix, code-review]
source: docs/content/develop/custom-code.md + docs/content/develop/diffs.md
status: approved
last_verified: 2026-11-06
---

# Import-Safe Flatteners & `ForceNew` Server Defaults

## 1. Populate Readable Fields Unconditionally from `v`

Flatteners for fields returned by `GET` must read unconditionally from the API response value `v` (or response map `res`)—never guard the flattener on `!d.GetOk("field")` or `d.Get(...)`.

- **Why `!d.GetOk(...)` in flatteners breaks import:** During `terraform import` (`ImportStatePassthrough`) and `ImportStateVerify`, `ResourceData` `d` starts with only `id` populated (`d.GetOk("field")` is `false` for unpopulated keys, zero values, and empty slices). Gating a flattener on `if _, ok := d.GetOk("field"); !ok { return nil }` passes normal `Create`/`Update` steps where state is already set, then silently drops server state on import.
- **Sole exception:** Reading prior state via `d.Get(...)` for write-only secrets intentionally omitted by `GET` (such as `initial_user.0.password` in `mmv1/templates/terraform/custom_flatten/alloydb_cluster_input_user_flatten.go.tmpl`), paired with `ImportStateVerifyIgnore` in acceptance tests.

## 2. Never Gate `ForceNew` Server Defaults on `d.GetOk` in `Read`

When `GET` populates a server default on an unconfigured `ForceNew` field, **never** suppress diffs with `if _, ok := d.GetOk("field"); !ok { return nil }` in `Read` or flatteners:

- `terraform import` writes the server value into prior state; on the next `terraform plan`, `d.GetOk` is `true` while HCL is unset, triggering a destructive `ForceNew` recreation.
- **Rule:** Keep `Read` and flatteners faithful to the API response, and suppress plan diffs via `DiffSuppressFunc` (such as `tpgresource.EmptyOrDefaultStringSuppress(...)` in `mmv1/products/compute/Address.yaml`) or a separate computed `effective_*` attribute.
