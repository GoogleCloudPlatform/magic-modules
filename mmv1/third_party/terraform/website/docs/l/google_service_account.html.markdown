---
subcategory: "Cloud Platform"
description: |-
  List Google Cloud IAM service accounts in a project for use with terraform query
  and .tfquery.hcl files.
---

# google_service_account (list)

Use this **list** resource type with the Terraform CLI [`terraform query`](https://developer.hashicorp.com/terraform/cli/commands/query)
command and configuration in **`.tfquery.hcl`** files. It enumerates existing
[`google_service_account`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_service_account)
resources in a project.

For general list block arguments (`provider`, `include_resource`, `limit`, `count`, `for_each`, etc.),
see the Terraform language reference for the [`list` block](https://developer.hashicorp.com/terraform/language/block/tfquery/list).

## Example

Place the `list` block in a file named with the `.tfquery.hcl` suffix (for example `service_accounts.tfquery.hcl`).
You can keep provider configuration in ordinary `.tf` files in the same directory.

```hcl
terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
    }
  }
}

provider "google" {
  project = "my-project"
}

list "google_service_account" "all" {
  provider = google

  config {
    # Optional. If omitted, the provider default project is used.
    # project = "other-project"
  }
}
```

Run from the directory that contains your `.tfquery.hcl` file:

```shell
terraform query
```

## Configuration (`config` block)

The following arguments are supported inside the nested `config` block:

* `project` - (Optional) The Google Cloud project ID whose service accounts are listed.
  If unset, the provider's configured default project is used (same behavior as the managed resource).

## Results

Each matching service account is returned as a list result. By default, Terraform returns
**resource identity** for `google_service_account` (see
[Resource identity](https://developer.hashicorp.com/terraform/language/resources/identities)):

* `email` - Service account email address (required for identity).
* `project` - Project ID (optional in identity when it can be inferred).

Set `include_resource = true` on the `list` block to include the **full resource object**
in each result, with the same attributes as the managed
[`google_service_account` resource](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_service_account#attributes-reference)
(for example `unique_id`, `name`, `display_name`, `disabled`, `description`, `member`, and `account_id` where present in state).

## API

Listing uses the IAM API method
[`projects.serviceAccounts.list`](https://cloud.google.com/iam/docs/reference/rest/v1/projects.serviceAccounts/list).

The caller must have permission to list service accounts in the target project (for example
`iam.serviceAccounts.list` on the project).

## Timeouts and pagination

The provider requests pages from the API until all accounts are returned (subject to your
`list` block `limit`, if set). Standard provider HTTP timeouts and retry behavior apply.
