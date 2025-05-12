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

In addition to the arguments listed above, the following computed attributes are exported:

* `labels` - GCP labels for this Feature.
* `resource_state` - State of the Feature resource itself.
  * `state` - The current state of the Feature resource in the Hub API.
* `spec` - The Hub-wide configuration of the feature.
  * `multiclusteringress` - Multicluster Ingress-specific spec.
  * `configmanagement` - Config Management-specific spec.
  * `servicemesh` - Service Mesh-specific spec.
  * `policycontroller` - Policy Controller-specific spec.
  * `clusterupgrade` - Cluster Upgrade-specific spec.
* `state` - The Hub-wide Feature state.
  * `state` - The current state of the Feature in the Hub API.
  * `servicemesh` - Service Mesh-specific state.