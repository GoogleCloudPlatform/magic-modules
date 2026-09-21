---
subcategory: "Cloud Build"
description: |-
  Get information about a Google Cloud Build worker pool.
---

# google_cloudbuild_worker_pool

To get more information about Cloudbuild worker pool, see:

* [API documentation](https://docs.cloud.google.com/build/docs/api/reference/rest/v1/projects.locations.workerPools)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/build/docs/automating-builds/create-manage-triggers)

## Example Usage

```hcl
data "google_cloudbuild_worker_pool" "name" {
  project  = "your-project-id"
  name     = "your-pool"
  location = "europe-west1"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The Cloud Build location for the worker pool.

* `name` - (Required) User-defined name of the worker pool.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

- - -

## Attributes Reference

See [google_cloudbuild_worker_pool](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloudbuild_worker_pool) resource for details of the available attributes.
