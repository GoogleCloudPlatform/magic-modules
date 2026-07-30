---
subcategory: "BigQuery"
description: |-
  List IAM member bindings for a BigQuery dataset for use with terraform query
  and .tfquery.hcl files.
---

# google_bigquery_dataset_iam_member (list)

Lists IAM **member bindings** for a BigQuery dataset for use with
[`terraform query`](https://developer.hashicorp.com/terraform/cli/commands/query) and
**`.tfquery.hcl`** files. Results correspond to existing
[`google_bigquery_dataset_iam_member`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/bigquery_dataset_iam)
managed resources.

For how list resources work in this provider, file layout, Terraform version requirements, and
shared `list` block arguments, refer to the guide
[Use list resources with terraform query (Google Cloud provider)](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/using_list_resources_with_terraform_query).

## Example

```hcl
list "google_bigquery_dataset_iam_member" "all" {
  provider = google

  config {
    dataset_id = "my_dataset"
    # Optional. Defaults to the provider project when omitted.
    # project = "other-project"
  }
}
```

Run `terraform query` from the directory that contains the `.tfquery.hcl` file.

## Configuration (`config` block)

* `dataset_id` - (Required) BigQuery dataset ID to list IAM members from.

* `project` - (Optional) Project ID that owns the dataset. If unset, the provider's
  configured default project is used.

* `role` - (Optional) If set, only bindings with this exact role are returned.
  For example, `roles/editor`.

* `member` - (Optional) If set, only bindings where this principal is a member
  are returned. For example, `user:jane@example.com`.

## Results

By default each result includes **resource identity** for `google_bigquery_dataset_iam_member` (see
[Resource identity](https://developer.hashicorp.com/terraform/language/resources/identities)):

* `project` - Project ID that owns the dataset.
* `dataset_id` - Dataset ID the binding belongs to.
* `role` - The IAM role, for example `roles/editor`.
* `member` - The principal, for example `user:jane@example.com`.

With `include_resource = true` on the `list` block, results also include full resource-style
attributes documented for the managed
[`google_bigquery_dataset_iam_member` resource](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/bigquery_dataset_iam#google_bigquery_dataset_iam_member)
(for example `etag` and `condition` where present in state).
