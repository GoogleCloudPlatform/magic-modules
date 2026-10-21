---
subcategory: "Google Cloud NetApp Volumes"
description: |-
  Fetches the details of a Google Cloud NetApp Volumes storage pool.
---

# google_netapp_storage_pool

Use this data source to get information about a Google Cloud NetApp Volumes storage pool. For more details, see the [API documentation](https://cloud.google.com/netapp/volumes/docs/reference/rest/v1/projects.locations.storagePools) and [official documentation](https://cloud.google.com/netapp/volumes/docs/get-started/quickstarts/create-storage-pool).

## Example Usage

```hcl
data "google_netapp_storage_pool" "example" {
  name     = "my-storage-pool"
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The resource name of the storage pool.

- - -

* `location` - (Optional) Name of the location of the storage pool. If it is not provided, the provider region or zone is used.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_netapp_storage_pool](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/netapp_storage_pool) resource for details of all the available attributes.
