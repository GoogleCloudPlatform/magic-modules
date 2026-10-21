---
subcategory: "Network Connectivity"
description: |-
  Fetches the details of a Network Connectivity Hub.
---

# google_network_connectivity_hub

Use this data source to get information about a Network Connectivity Hub. For more details, see the [API documentation](https://cloud.google.com/network-connectivity/docs/reference/networkconnectivity/rest/v1/projects.locations.global.hubs) and [Official Documentation](https://cloud.google.com/network-connectivity/docs/network-connectivity-center/concepts/overview).

## Example Usage

```hcl
data "google_network_connectivity_hub" "default" {
  name = "my-hub"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the hub.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_network_connectivity_hub](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/network_connectivity_hub) resource for details of all the available attributes.
