---
subcategory: "Compute Engine"
description: |-
  Get a network within GCE.
---

# google_compute_network

Get a network within GCE from its name.

## Example Usage

```tf
data "google_compute_network" "my-network" {
  name = "default-us-east1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the network.


- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - an identifier for the resource with format projects/{{project}}/global/networks/{{name}}

* `description` - Description of this network.

* `network_id` - The numeric unique identifier for the resource.

* `numeric_id` - (Deprecated) The numeric unique identifier for the resource. `numeric_id` is deprecated and will be removed in a future major release. Use `network_id` instead.

* `gateway_ipv4` - The IP address of the gateway.

* `internal_ipv6_range` - The ula internal ipv6 range assigned to this network.

* `network_profile` - A full or partial URL of the network profile to apply to this network.

* `subnetworks_self_links` - the list of subnetworks which belong to the network

* `peerings` - The network peerings configured on this network, including peerings not managed by Terraform. Structure is [documented below](#nested_peerings).

* `self_link` - The URI of the resource.

<a name="nested_peerings"></a>The `peerings` block contains:

* `name` - Name of the peering.

* `network` - The URI of the peer network.

* `state` - State for the peering, either `ACTIVE` or `INACTIVE`. The peering is `ACTIVE` when there's a matching configuration in the peer network.

* `state_details` - Details about the current state of the peering.

* `export_custom_routes` - Whether custom routes are exported to the peer network.

* `import_custom_routes` - Whether custom routes are imported from the peer network.

* `export_subnet_routes_with_public_ip` - Whether subnet routes with public IP ranges are exported to the peer network.

* `import_subnet_routes_with_public_ip` - Whether subnet routes with public IP ranges are imported from the peer network.

* `stack_type` - Which IP version(s) of traffic and routes are exchanged with the peer network, either `IPV4_ONLY` or `IPV4_IPV6`.

* `update_strategy` - The update strategy of the peering, either `INDEPENDENT` or `CONSENSUS`.
