---
subcategory: "GKEHub"
description: |-
  Retrieves the details of a GKE Hub Feature.
---

# `google_gke_hub_feature`
Retrieves the details of a specific GKE Hub Feature. Use this data source to retrieve the feature's configuration and state.

## Example Usage

```hcl
data "google_gke_hub_feature" "example" {
  location = "global"
  name     = "servicemesh"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the feature you want to know the status of.
* `location` - (Required) The location for the GKE Hub Feature.
* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

## Attributes Reference

See [google_gke_hub_feature](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/gke_hub_feature) resource for details of the available attributes.