---
name: test-adequacy-traps
description: Prevent false-pass CheckDestroy templates, no-op update steps, ExpectError-only configs, and samples[0] GA breaks.
topics: [test, custom-code]
task_types: [field-add, new-resource, bug-fix, test-fix]
source: docs/content/test/test.md + docs/content/reference/sample.md + mmv1/api/resource.go
status: approved
last_verified: 2026-11-06
---

# Acceptance Test Adequacy Traps

For general test review criteria (2-step update sufficiency, `ImportStateVerify`, unit-test scope, and GA test verification), see [`docs/content/code-review/review-pr.md`](../../../docs/content/code-review/review-pr.md) and [`docs/content/test/test.md#sweepers`](../../../docs/content/test/test.md#sweepers).

## 1. Custom `test_check_destroy` (`custom_check_destroy`) Error Handling

`custom_code.test_check_destroy` replaces the `else` branch inside `for name, rs := range s.RootModule().Resources` in `mmv1/templates/terraform/samples/base_configs/test_file.go.tmpl` (`config`, `url`, and `billingProject` are not pre-declared in scope).

- **Do not treat arbitrary `err != nil` as deleted:** Checking only `if err == nil { return fmt.Errorf("still exists") }` lets a `400 Bad Request` (from a malformed URL or missing query parameter) pass `CheckDestroy` while leaking the cloud resource. Confirm expected deletion responses (such as `transport_tpg.IsGoogleApiErrorWithCode(err, 404)`, `403` on IAM-wrapped endpoints that mask non-existence, or the API's documented deleted state on `200`), `continue` so remaining state resources are checked, and `return err` on unexpected errors.
- **Keep eventual-consistency retries:** Do not strip `transport_tpg.Retry` loops when an API tears down asynchronously.

```go
config := acctest.GoogleProviderConfig(t)
url, err := tpgresource.ReplaceVarsForTest(config, rs, transport_tpg.BaseUrl(myproduct.Product, config)+"{{$.SelfLinkUri}}")
if err != nil {
	return err
}
billingProject := ""
if config.BillingProject != "" {
	billingProject = config.BillingProject
}
res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
	Config:    config,
	Method:    "GET",
	Project:   billingProject,
	RawURL:    url,
	UserAgent: config.UserAgent,
})
if err != nil {
	if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
		continue
	}
	return err
}
return fmt.Errorf("{{ $.ResourceName }} still exists at %s: %v", url, res)
```

## 2. Silent No-Op Update & `ExpectError` Steps

- **Identical HCL across steps never tests `Update`:** Setting a mutable field to the **same** value in Step 1 (`_basic`) and Step 2 (`_update`) produces an empty diff—neither `d.HasChange` nor `updateMask` runs. Mutate the value (`valA -> valB`) or introduce the field in Step 2 (`unset -> valA`).
- **`ExpectError` / `PlanOnly` steps are negative-only:** They never create, read, update, or import state. Every configurable field still needs a positive applied `Config` step followed by `ImportStateVerify: true`.

## 3. `FirstTestConfig()` Generator Nuance

`FirstTestConfig()` (`mmv1/api/resource.go`) picks the first sample in `samples:` without `exclude_test: true` or a higher sample-level `min_version` to scaffold generated IAM and datasource tests—it does **not** filter out `skip_test` or step-level `min_version: beta`. Keep a runnable, GA-compatible minimal sample first (see [`docs/content/reference/sample.md`](../../../docs/content/reference/sample.md)).

## Do NOT Use For

- Do **not** add `custom_check_destroy` when MMv1's default `GET` on `self_link` already returns `404` on deletion—customize only for `200` soft-delete states, singleton sub-object resets, or custom query parameters.
