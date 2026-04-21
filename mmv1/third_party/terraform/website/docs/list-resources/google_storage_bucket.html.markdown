---
subcategory: "Cloud Storage"
description: |-
  List Google Cloud Storage buckets in a project for use with terraform query
  and .tfquery.hcl files.
---

# google_storage_bucket (list)

Lists **buckets** in a Google Cloud Storage project for use with
[`terraform query`](https://developer.hashicorp.com/terraform/cli/commands/query) and
**`.tfquery.hcl`** files. Results correspond to existing
[`google_storage_bucket`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_storage_bucket)
managed resources.

For how list resources work in this provider, file layout, Terraform version requirements, and
shared `list` block arguments, refer to the guide
[Use list resources with terraform query (Google Cloud provider)](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/using_list_resources_with_terraform_query).

## Example

```hcl
list "google_storage_bucket" "all" {
  provider = google

  config {
    # Optional. Defaults to the provider project when omitted.
    # project = "other-project"

    # Optional. Prefix filter passed to the Storage buckets.list API.
    # prefix = "logs-"

    # Optional. Positive cap on buckets returned per API page chain (omit for API default).
    # max_results = 100

    # Optional. "full" or "noAcl" (omit for API default, typically noAcl).
    # projection = "noAcl"

    # Optional. Maps to returnPartialSuccess query parameter.
    # return_partial_success = false

    # Optional. Maps to softDeleted query parameter (list soft-deleted buckets).
    # soft_deleted = false
  }
}
```

Run `terraform query` from the directory that contains the `.tfquery.hcl` file.

## Configuration (`config` block)

* `project` - (Optional) Project ID whose buckets are listed. If unset, the provider's configured
  default project is used (same idea as the managed resource).

* `prefix` - (Optional) Filter to buckets whose names begin with this prefix, per the
  [buckets.list](https://cloud.google.com/storage/docs/json_api/v1/buckets/list) `prefix` parameter.

* `max_results` - (Optional) Non-negative integer. When greater than zero, passed as `maxResults`
  to limit how many buckets the API may return in a single list response (pagination still applies).
  Must not be negative.

* `projection` - (Optional) When set, must be `full` or `noAcl` (after trimming whitespace), per
  the API `projection` parameter. Leave unset to use the API default.

* `return_partial_success` - (Optional) When `true`, sets the API `returnPartialSuccess` query
  parameter.

* `soft_deleted` - (Optional) When `true`, sets the API `softDeleted` query parameter to include
  soft-deleted buckets where supported.

## Results

By default each result includes **resource identity** for `google_storage_bucket` (see
[Resource identity](https://developer.hashicorp.com/terraform/language/resources/identities)):

* `name` - Bucket name (required for identity).
* `project` - Project ID when applicable (optional for import).

With `include_resource = true` on the `list` block, results also include the full resource-style
attributes documented for the managed
[`google_storage_bucket` resource](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_storage_bucket#attributes-reference)
(for example `location`, `storage_class`, `labels`, `terraform_labels`, `effective_labels`, and
other attributes present after read).
