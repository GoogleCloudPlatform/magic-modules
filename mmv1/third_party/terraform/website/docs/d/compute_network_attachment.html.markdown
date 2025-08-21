---
subcategory: "Compute Engine"
description: |-
  A data source to retrieve a network attachment
---

# `google_compute_network_attachment`

Get a specific network attachment within a region. For more information see
the [official documentation](https://cloud.google.com/vpc/docs/about-network-attachments)
and [API](https://cloud.google.com/compute/docs/reference/rest/v1/networkAttachments/get).

## Example Usage

```hcl
data "google_compute_network_attachment" "default" {
  project    = "my-project"
  name       = "my-network-attachment"
  region     = "europe-west1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the network attachment to retrieve.
  The name must be unique within the region.

* `region` - (Required) The region in which the network attachment resides.
  For example, `europe-west1`.

* `project` - (Optional) The ID of the project in which the resource belongs.
  If it is not provided, the provider project is used.

## Attributes Reference

See [google_compute_network_attachment](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/bigquery_table#attributes-reference) resource for details of the available attributes.