---
subcategory: "Oracle Database"
description: |-
  Get information about an ODB Subnet.
---

# google_oracle_database_odb_subnet

Get information about an ODB Subnet.

For more information see the
* [API documentation](https://cloud.google.com/oracle/database/docs/reference/rest/v1/projects.locations.odbNetworks.odbSubnets)


## Example Usage

```hcl
data "google_oracle_database_odb_subnet" "my-subnet" {
  location = "us-east4"
  odbnetwork = "my-network-id"
  odb_subnet_id = "my-subnet-id"
}
```

## Argument Reference

The following arguments are supported:

* `odb_subnet_id` - (Required) The ID of the ODB Subnet.

* `odbnetwork` - (Required) The ID of the parent ODB Network.

* `location` - (Required) The location of the resource.

* `project` - (Optional) The project to which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes reference

See [google_oracle_database_odb_subnet](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/oracle_database_odb_subnet#argument-reference) resource for details of the available attributes.
