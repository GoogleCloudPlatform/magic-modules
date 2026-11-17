---
name: test-adequacy-traps
description: Avoid silent no-op update steps, ExpectError-only configs, and custom_check_destroy loop bugs in acceptance tests.
topics: [test]
task_types: [field-add, new-resource, bug-fix, test-fix]
source: docs/content/code-review/review-pr.md + docs/content/reference/sample.md
status: draft
last_verified: 2026-11-17
---

# Acceptance Test Adequacy Traps

For general acceptance and unit test review expectations, see [`docs/content/code-review/review-pr.md`](../../../docs/content/code-review/review-pr.md) and [`unit-test-scope.md`](unit-test-scope.md).

## 1. Silent No-Op Update & `ExpectError` Steps

- **Identical HCL across steps never tests `Update`:** Setting a mutable field to the **same** value in Step 1 (`_basic`) and Step 2 (`_update`) produces an empty diff—neither `d.HasChange` nor `updateMask` runs. Mutate the value (`valA -> valB`) or introduce the field in Step 2 (`unset -> valA`).
- **`ExpectError` / `PlanOnly` steps are negative-only:** Because they fail before apply, they never verify `Create`, `Read`, `Update`, or `ImportStateVerify`. Every configurable field still needs a positive applied `Config` step followed by `ImportStateVerify: true`.

## 2. Custom `test_check_destroy` Loop & Retry Pitfalls

`custom_code.test_check_destroy` replaces the branch inside `for name, rs := range s.RootModule().Resources` in `mmv1/templates/terraform/samples/base_configs/test_file.go.tmpl`:
- **Only customize when default `GET` is insufficient:** MMv1's generated `CheckDestroy` already handles standard deletion; reserve `custom_code.test_check_destroy` for `200 OK` soft-delete states, singleton resets, or custom query parameters (see `mmv1/templates/terraform/custom_check_destroy/`).
- **Use `continue`, verify URL substitution, and keep retries:** Once a resource is confirmed deleted, `continue` (not `return nil`) so remaining resources in state are checked, propagate `ReplaceVarsForTest` errors so malformed URLs cannot false-pass, and keep `transport_tpg.Retry` loops for asynchronous teardown.
