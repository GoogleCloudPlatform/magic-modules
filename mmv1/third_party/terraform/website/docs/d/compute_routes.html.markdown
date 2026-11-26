---
subcategory: "Compute Engine"
description: |-
  Get a list of routes in a project.
---

# google_compute_routes

Get a list of routes in a project. For more information see
the official [API](https://cloud.google.com/compute/docs/reference/rest/v1/routes/list) documentation.

~> **Note:** This data source returns the routes that the `routes.list` API returns. According to [List routes for a VPC network](https://cloud.google.com/vpc/docs/using-routes#listingroutes), that covers system-generated default routes, subnet routes, peering subnet routes, NCC subnet routes and static routes, and does not include policy-based routes, dynamic routes, or peering static and dynamic routes.

## Example Usage

```hcl
data "google_compute_routes" "all" {
  project = "my-project"
  filter  = "network eq .*/my-vpc"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A string filter as defined in the [REST API](https://cloud.google.com/compute/docs/reference/rest/v1/routes/list#query-parameters).

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/global/routes`

* `routes` - The list of routes in the project. Structure is [documented below](#nested_routes).

<a name="nested_routes"></a>The `routes` block supports:

Links to Compute resources are returned in the `/compute/v1/` form in both the `google` and `google-beta` providers.

* `name` - Name of the route.
* `description` - Description of the route.
* `network` - URL of the network that the route applies to.
* `dest_range` - The destination range of outgoing packets that the route applies to.
* `priority` - The priority of the route. Where routes have an equal prefix length, the one with the lowest priority value wins.
* `tags` - A list of instance tags to which the route applies.
* `next_hop_gateway` - URL to a gateway that should handle matching packets.
* `next_hop_instance` - URL to an instance that should handle matching packets.
* `next_hop_ip` - Network IP address of an instance that should handle matching packets.
* `next_hop_vpn_tunnel` - URL to a VpnTunnel that should handle matching packets.
* `next_hop_ilb` - The URL or IP address of a forwarding rule of type `loadBalancingScheme=INTERNAL` that should handle matching packets.
* `next_hop_network` - URL of the local network if it should handle matching packets.
* `next_hop_peering` - The network peering name that should handle matching packets.
* `next_hop_hub` - The full resource name of the Network Connectivity Center hub that should handle matching packets.
* `next_hop_interconnect_attachment` - URL to an InterconnectAttachment which is the next hop for the route.
* `next_hop_origin` - Indicates the origin of the route. Can be `IGP` (Interior Gateway Protocol), `EGP` (Exterior Gateway Protocol), or `INCOMPLETE`.
* `next_hop_med` - Multi-Exit Discriminator, a BGP route metric that indicates the desirability of a particular route in a network.
* `next_hop_inter_region_cost` - Internal fixed region-to-region cost that Google Cloud calculates based on factors such as network performance, distance, and available bandwidth between regions.
* `route_type` - The type of the route. One of `TRANSIT`, `SUBNET`, `STATIC` or `BGP`.
* `route_status` - The status of the route. One of `ACTIVE`, `DROPPED`, `INACTIVE` or `PENDING`.
* `self_link` - The URI of the route.
* `creation_timestamp` - Creation timestamp in RFC3339 text format.
