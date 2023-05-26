---
subcategory: "Cloud VMware Engine"
description: |-
  Get info about a Google VMwareEngine Cluster.
---

# google\_vmwareengine\_cluster

Use this data source to get details about a Google VMwareEngine Cluster resource.

To get more information about VMwareEngine Cluster, see:
* [API documentation](https://cloud.google.com/vmware-engine/docs/reference/rest/v1/projects.locations.privateClouds.clusters)

## Example Usage

```hcl
data "google_vmwareengine_cluster" "my_cluster" {
  provider = google-beta
  name     = "my-pc"
  parent   = "project/locations/us-west1-a/privateClouds/my-cloud"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the resource.
* `parent` - (Required) The resource name of the private cloud that this cluster belongs.
