Fixes https://github.com/hashicorp/terraform-provider-google/issues/7659

`google_bigquery_dataset_access` has no unique ID of its own in the API — an access grant is identified within a dataset's `access[]` list by an aggregate of `role` plus exactly one member field (`user_by_email`, `group_by_email`, `domain`, `special_group`, `iam_member`, or the `view`/`dataset`/`routine` nested blocks). This adds composite-ID import support covering all of them.

## Changes

- `mmv1/products/bigquery/DatasetAccess.yaml`: remove `exclude_import: true`, add `import_format` (all accepted id shapes, used for generated docs) and `custom_code.custom_import`.
- `mmv1/templates/terraform/custom_import/bigquery_dataset_access.go.tmpl`: new. Parses the flat member-field id shapes via `tpgresource.ParseImportId`, delegating to the nested-block helper below first.
- `mmv1/templates/terraform/constants/bigquery_dataset_access.go.tmpl`: add `resourceBigQueryDatasetAccessImportNested`, since `view`/`dataset`/`routine` can't be set through `ParseImportId` (it only assigns flat string fields, not the nested list-of-object schema these use).
- `mmv1/api/resource.go`: fix `IdentityProperties()`, which was pulling an arbitrary field out of `import_format` into the unrelated Resource Identity schema whenever a resource also sets `exclude_identity_from_identity_import` — a combination this resource is the first to use. No other resource combines both flags, so this is a no-op everywhere else (verified via a full `go run . --output=...` regeneration diff).
- `mmv1/third_party/terraform/services/bigquery/resource_bigquery_dataset_access_test.go`: add an import step after every handwritten `Config` step across all 12 `TestAccBigQueryDatasetAccess_*` tests, reading the expected import id from post-apply state (rather than hardcoding it) since the API normalizes some fields, e.g. a predefined role like `roles/bigquery.dataEditor` is stored as `WRITER`.

## Local verification

- `go run . --output=<provider> --product=bigquery --resource=DatasetAccess --version=ga` generates cleanly.
- `go build` and `go vet` pass on the regenerated `google/services/bigquery` package.
- Full unfiltered regeneration diffed against baseline to confirm the `api/resource.go` change only affects `DatasetAccess`.
- Accepted import id formats cross-checked against the actual regexes with a standalone Go script (35 shapes, all match; adversarial near-misses like a missing `roles/` literal or wrong keyword correctly don't).

```release-note:enhancement
bigquery: added import support to `google_bigquery_dataset_access` resource
```
