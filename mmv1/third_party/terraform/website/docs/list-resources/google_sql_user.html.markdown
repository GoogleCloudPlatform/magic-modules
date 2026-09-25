---
subcategory: "Cloud SQL"
description: |-
  List Cloud SQL users in a project and instance for use with terraform query
  and .tfquery.hcl files.
---

# google_sql_user (list)

Lists [google_sql_user](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/sql_user) resources for use with [terraform query](https://developer.hashicorp.com/terraform/cli/commands/query) and **.tfquery.hcl** files.

For how list resources work in this provider, file layout, Terraform version requirements, and shared list block arguments, refer to the guide [Use list resources with terraform query (Google Cloud provider)](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/using_list_resources_with_terraform_query).

## Example

```hcl
list "google_sql_user" "all" {
  provider = google

  config {
    # Optional. Defaults to the provider project when omitted.
    # project = "my-project"

    # Required. Cloud SQL instance name to list users from.
    instance = "my-instance"
  }
}
```

Run `terraform query` from the directory that contains the `.tfquery.hcl` file.

## Configuration (`config` block)

* `project` - (Optional) Project ID to list Cloud SQL users from. If unset, the provider's configured default project is used, matching managed resource behavior.

* `instance` - (Required) Cloud SQL instance name to list users from.

## Results

By default each result includes **resource identity** for `google_sql_user` (see [Resource identity](https://developer.hashicorp.com/terraform/language/block/import#identity)):

* `name` - User name.
* `instance` - Cloud SQL instance name.
* `project` - Project ID.
* `host` - Host from which the user can connect. MySQL-only; empty for Postgres.

With `include_resource = true` on the `list` block, results also include the full resource-style attributes documented for the managed [google_sql_user resource](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/sql_user#attributes-reference).
