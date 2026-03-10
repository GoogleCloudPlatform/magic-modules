---
subcategory: "Compute Engine"
description: |-
 Get a list of routers.
---

# google_compute_routers
Get a list of routers. For more information see
the official [API](https://cloud.google.com/compute/docs/reference/rest/v1/routers/list) documentation.

## Example Usage
```tf
data "google_compute_routers" "all" {
  project = google_compute_router.foobar.project
  region  = google_compute_router.foobar.region
}
```

## Argument Reference

The following arguments are supported:
* `project` - (Optional) The project in which the resource belongs. If it
 is not provided, the provider project is used.
* `region` - (Optional) If provided, only resources from the given regions are queried.