---
subcategory: "AlloyDB"
description: |-
  Fetches the details of available locations.
---

# google_alloydb_instance

Use this data source to get information about the available instance. For more details refer the [API docs](https://cloud.google.com/alloydb/docs/reference/rest/v1/projects.locations.clusters.instances).

## Example Usage


```hcl
data "google_alloydb_instance" "qa" {
}
```

## Argument Reference

The following arguments are supported:

* `cluster` -
  (Required)
  Identifies the alloydb cluster. Must be in the format
  'projects/{project}/locations/{location}/clusters/{cluster_id}'

* `instance_id` -
  (Required)
  The ID of the alloydb instance.

## Attributes Reference

See [google_alloydb_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/alloydb_instance) resource for details of all the available attributes.
