---
subcategory: "Compute Engine"
description: |-
  Get a list of the interconnect groups in a project.
---

# google_compute_interconnect_groups

Get a list of the interconnect groups in a project. For more information see
[the official documentation](https://cloud.google.com/network-connectivity/docs/interconnect/concepts/dedicated-overview)
and [API](https://cloud.google.com/compute/docs/reference/rest/v1/interconnectGroups/list).

## Example Usage

```hcl
data "google_compute_interconnect_groups" "all" {
  project = "my-project"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The ID of the project whose interconnect groups are listed. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/global/interconnectGroups`

* `interconnect_groups` - The list of interconnect groups in the project. Each group has the attributes of the
    [google_compute_interconnect_group](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_interconnect_group)
    resource, including its `id`, except `project` and `deletion_policy`.
