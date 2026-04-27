---
subcategory: "Compute Engine"
description: |-
  List Google Compute Engine VM instances in a project and zone for use with terraform query
  and .tfquery.hcl files.
---

# google_compute_instance (list)

Lists **Compute Engine VM instances** in a Google Cloud project and zone for use with
[`terraform query`](https://developer.hashicorp.com/terraform/cli/commands/query) and
**`.tfquery.hcl`** files. Results correspond to existing
[`google_compute_instance`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_instance)
managed resources.

For how list resources work in this provider, file layout, Terraform version requirements, and
shared `list` block arguments, refer to the guide
[Use list resources with terraform query (Google Cloud provider)](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/using_list_resources_with_terraform_query).

## Example

```hcl
list "google_compute_instance" "all" {
  provider = google

  config {
    # Optional. Defaults to the provider project when omitted.
    # project = "other-project"

    # Optional. Defaults to the provider zone when omitted.
    # zone = "us-central1-a"
  }
}
```

Run `terraform query` from the directory that contains the `.tfquery.hcl` file.

## Configuration (`config` block)

* `project` - (Optional) Project ID to list instances from. If unset, the provider's
  configured default project is used (same idea as the managed resource).

* `zone` - (Optional) Zone to list instances in. If unset, the provider's configured
  default zone is used (same idea as the managed resource).

## Results

By default each result includes **resource identity** for `google_compute_instance` (see
[Resource identity](https://developer.hashicorp.com/terraform/language/resources/identities)):

* `name` - Instance name (required for identity).
* `project` - Project ID when applicable.
* `zone` - Zone the instance resides in.

With `include_resource = true` on the `list` block, results also include resource attributes
populated from the API response. These include `machine_type`, `self_link`, `current_status`,
`description`, `tags`, `labels`, `metadata`, `network_interface`, `boot_disk`, `scratch_disk`,
`attached_disk`, `service_account`, `scheduling`, `guest_accelerator`, `shielded_instance_config`,
`confidential_instance_config`, `advanced_machine_features`, `deletion_protection`, `hostname`,
`cpu_platform`, `instance_id`, `creation_timestamp`, and `reservation_affinity`.

Note: `metadata_startup_script` is not populated (it is state-dependent). Attached disk ordering
and raw disk encryption key fields are also omitted as they depend on prior configuration state.
