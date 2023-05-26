---
subcategory: "Cloud VMware Engine"
description: |-
  Get info about a Google VMwareEngine Network.
---

# google\_vmwareengine\_network

Use this data source to get details about a Google VMwareEngine network resource.

To get more information about VMwareEngine Network, see:
* [API documentation](https://cloud.google.com/vmware-engine/docs/reference/rest/v1/projects.locations.vmwareEngineNetworks)

## Example Usage

```hcl
data "google_vmwareengine_network" "my_nw" {
  provider = google-beta
  name     = "us-central1-default"
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the resource.
* `location` - (Required) Location of the resource.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.