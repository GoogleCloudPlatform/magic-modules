---
subcategory: "Cloud Workstations"
description: |-
  Get information about a Cloud Workstations workstation cluster.
---

# google_workstations_workstation_cluster

Use this data source to get information about a Cloud Workstations workstation cluster.

For more information see the
[official documentation](https://cloud.google.com/workstations/docs)
and [API](https://cloud.google.com/workstations/docs/reference/rest/v1/projects.locations.workstationClusters).

## Example Usage

```hcl
data "google_workstations_workstation_cluster" "default" {
  workstation_cluster_id = "my-workstation-cluster"
  location               = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `workstation_cluster_id` - (Required) The ID of the workstation cluster.

* `location` - (Optional) The location where the workstation cluster resides. If it is not provided, the provider region is used.

* `project` - (Optional) The project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_workstations_workstation_cluster](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/workstations_workstation_cluster#attributes-reference) Resource for details of the available attributes.
