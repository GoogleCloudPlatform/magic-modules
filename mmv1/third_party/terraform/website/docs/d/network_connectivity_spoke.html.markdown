---
subcategory: "Network Connectivity"
description: |-
  Fetches the details of a Network Connectivity Spoke.
---

# google_network_connectivity_spoke

Use this data source to get information about a Network Connectivity Spoke. For more details, see the [API documentation](https://cloud.google.com/network-connectivity/docs/reference/networkconnectivity/rest/v1/projects.locations.spokes) and [Official Documentation](https://cloud.google.com/network-connectivity/docs/network-connectivity-center/concepts/overview).

## Example Usage

```hcl
data "google_network_connectivity_spoke" "default" {
  name     = "my-spoke"
  location = "global"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the spoke.

- - -

* `location` - (Optional) The location for the resource. If it is not provided, the provider region or zone is used.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_network_connectivity_spoke](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/network_connectivity_spoke) resource for details of all the available attributes.
