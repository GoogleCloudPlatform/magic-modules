---
subcategory: "Compute Engine"
description: |-
  Get a list of the interconnect attachments (VLAN attachments) in a project.
---

# google_compute_interconnect_attachments

Get a list of the interconnect attachments (VLAN attachments) in a project, in one region or in every region. For more information see
[the official documentation](https://cloud.google.com/network-connectivity/docs/interconnect/how-to/dedicated/creating-vlan-attachments)
and the API [list](https://cloud.google.com/compute/docs/reference/rest/v1/interconnectAttachments/list) and
[aggregated list](https://cloud.google.com/compute/docs/reference/rest/v1/interconnectAttachments/aggregatedList) methods.

## Example Usage

```hcl
data "google_compute_interconnect_attachments" "all" {
  project = "my-project"
}

data "google_compute_interconnect_attachments" "europe_west1" {
  project = "my-project"
  region  = "europe-west1"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `region` - (Optional) The region to list, as a name or a self link. If it is not
    provided, the attachments of every region are listed. The provider region is not used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/regions/{{region}}/interconnectAttachments`, or `projects/{{project}}/aggregated/interconnectAttachments` when no region is given

* `interconnect_attachments` - The list of interconnect attachments, sorted by region.

See [google_compute_interconnect_attachment](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_interconnect_attachment) resource for details of the available attributes on each attachment, with these differences:

* `region` is the name of the region, not a link to it.
* `labels` holds every label of the attachment.
* Links to Compute resources, such as `router` and `interconnect`, are returned in the `/compute/v1/` form in both the `google` and `google-beta` providers.
* The `project`, `deletion_policy`, `params`, `candidate_subnets`, `subnet_length`, `terraform_labels` and `effective_labels` arguments and attributes of that resource are not part of each attachment.
