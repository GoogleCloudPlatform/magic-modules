---
subcategory: "Oracle Database"
description: |-
  Get information about an ODB Network.
---

# google_oracle_database_odb_network

Get information about an ODB Network.

For more information see the
[API](https://cloud.google.com/oracle/database/docs/reference/rest/v1/projects.locations.odbNetworks).

## Example Usage

```hcl
data "google_oracle_database_odb_network" "my-network" {
  location = "us-east4"
  odb_network_id = "my-network-id"
}
```

## Argument Reference

The following arguments are supported:

* `odb_network_id` - (Required) The ID of the ODB Network.

* `location` - (Required) The location of the resource.

* `project` - (Optional) The project to which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes reference

See [google_oracle_database_odb_network](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/oracle_database_odb_network#argument-reference) resource for details of the available attributes.
