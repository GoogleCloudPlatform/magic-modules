---
subcategory: "Assured Workloads V2"
description: |-
  Get information about an Assured Workloads V2 workload.
---

# google_assured_workloads_v2_workload

Get information about an existing Assured Workloads V2 workload.

To get more information about Assured Workloads, see:

* [API documentation](https://cloud.google.com/assured-workloads/docs/reference/rest/v2/organizations.locations.workloads)
* How-to Guides
    * [Assured Workloads Overview](https://cloud.google.com/assured-workloads/docs/overview)

~> **Note:** Workloads have no project, so Assured Workloads attributes API quota and billing to the client project of your credentials unless a billing project is set. Set `billing_project` and `user_project_override = true` in the provider configuration to use a project where the Assured Workloads API is enabled.

## Example Usage

Using `organization`, `location`, and `workload_id`:

```hcl
data "google_assured_workloads_v2_workload" "my_workload" {
  organization = "123456789"
  location     = "us-east4"
  workload_id  = "my-workload-id"
}
```

Or using the full resource `name`:

```hcl
data "google_assured_workloads_v2_workload" "my_workload" {
  name = "organizations/123456789/locations/us-east4/workloads/my-workload-id"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The full resource name of the workload, in the format `organizations/{organization}/locations/{location}/workloads/{workload_id}`. Either `name` or all of `organization`, `location`, and `workload_id` must be specified.

* `organization` - (Optional) The ID of the parent organization of the workload, for example `123456789`. Required when `name` is not specified.

* `location` - (Optional) The location of the workload, for example `us-east4`. Required when `name` is not specified.

* `workload_id` - (Optional) The unique identifier of the workload within the parent organization and location. Required when `name` is not specified.

## Attributes Reference

See [google_assured_workloads_v2_workload](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/assured_workloads_v2_workload) resource for details of the available attributes. Note that the create-only `resource_config` block is not exported by this data source; use `computed_target_resource` to inspect the resolved target folder or project.
