---
subcategory: "Google Cloud NetApp Volumes"
description: |-
  Fetches the details of a Google Cloud NetApp Volumes volume.
---

# google_netapp_volume

Use this data source to get information about a Google Cloud NetApp Volumes volume. For more details, see the [API documentation](https://cloud.google.com/netapp/volumes/docs/reference/rest/v1/projects.locations.volumes) and [official documentation](https://cloud.google.com/netapp/volumes/docs/configure-and-use/volumes/overview).

## Example Usage

```hcl
data "google_netapp_volume" "example" {
  name     = "my-volume"
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the volume.

- - -

* `location` - (Optional) Name of the location of the volume. If it is not provided, the provider region or zone is used.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_netapp_volume](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/netapp_volume) resource for details of all the available attributes.
