---
subcategory: "Compute Engine"
description: |-
 Get a list of routers.
---

# google_compute_routers
Get a list of routers. For more information see
the official [API](https://cloud.google.com/compute/docs/reference/rest/v1/routers/list) documentation.

## Example Usage
```tf
data "google_compute_routers" "all" {
  project = google_compute_router.foobar.project
  region  = google_compute_router.foobar.region
}
```

## Argument Reference

The following arguments are supported:
* `project` - (Optional) The project in which the resource belongs. If it
 is not provided, the provider project is used.
* `region` - (Optional) The region to list routers in. If it is not provided,
 the provider region is used. Routers in other regions are not returned.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/regions/{{region}}`

* `routers` - The list of routers in the region. Structure is [documented below](#nested_routers).

<a name="nested_routers"></a>The `routers` block supports:

* `name` - Name of the router.
* `network` - URI of the network to which the router belongs.
* `description` - Description of the router.
* `creation_timestamp` - Creation timestamp in RFC3339 text format.
* `self_link` - The URI of the router.
* `bgp` - BGP information specific to the router. Structure is [documented below](#nested_bgp).
* `interfaces` - The interfaces of the router. Structure is [documented below](#nested_interfaces).
* `bgp_peers` - The BGP peers of the router. Structure is [documented below](#nested_bgp_peers).
* `nats` - The NAT services created in the router. Structure is [documented below](#nested_nats).
* `encrypted_interconnect_router` - Indicates if the router is dedicated for use with encrypted VLAN attachments (interconnectAttachments).
* `md5_authentication_keys` - Keys used for MD5 authentication. Structure is [documented below](#nested_md5_authentication_keys).

<a name="nested_bgp"></a>The `bgp` block supports:

* `asn` - Local BGP Autonomous System Number (ASN).
* `advertise_mode` - User-specified flag to indicate which mode to use for advertisement. One of `DEFAULT` or `CUSTOM`.
* `advertised_groups` - User-specified list of prefix groups to advertise in custom mode.
* `advertised_ip_ranges` - User-specified list of individual IP ranges to advertise in custom mode. Structure is [documented below](#nested_advertised_ip_ranges).
* `keepalive_interval` - The interval in seconds between BGP keepalive messages that are sent to the peer.

<a name="nested_advertised_ip_ranges"></a>The `advertised_ip_ranges` block supports:

* `range` - The IP range to advertise, as a CIDR-formatted string.
* `description` - User-specified description for the IP range.

<a name="nested_interfaces"></a>The `interfaces` block supports:

* `name` - Name of the interface.
* `linked_vpn_tunnel` - URI of the VPN tunnel linked to the interface.
* `linked_interconnect_attachment` - URI of the Interconnect attachment linked to the interface.
* `ip_range` - IP address and range of the interface.
* `private_ip_address` - The regional private internal IP address that is used to establish BGP sessions to a VM instance acting as a third-party Router Appliance.
* `redundant_interface` - Name of the interface that is redundant to this interface.
* `subnetwork` - The URI of the subnetwork resource that the interface belongs to.

<a name="nested_bgp_peers"></a>The `bgp_peers` block supports:

* `name` - Name of the BGP peer.
* `interface_name` - Name of the interface the BGP peer is associated with.
* `ip_address` - IP address of the interface inside Google Cloud Platform.
* `peer_ip_address` - IP address of the BGP interface outside Google Cloud Platform.
* `peer_asn` - Peer BGP Autonomous System Number (ASN).
* `advertised_route_priority` - The priority of routes advertised to the BGP peer.
* `advertise_mode` - User-specified flag to indicate which mode to use for advertisement. One of `DEFAULT` or `CUSTOM`.
* `management_type` - The resource that configures and manages the BGP peer. `MANAGED_BY_USER` for a peer managed by users, or `MANAGED_BY_ATTACHMENT` for a peer that is configured and managed by Cloud Interconnect.
* `enable` - The status of the BGP peer connection, as the string `TRUE` or `FALSE`.
* `enable_ipv6` - Whether IPv6 traffic is enabled over the BGP peer.

<a name="nested_nats"></a>The `nats` block supports:

* `name` - Name of the NAT service.
* `source_subnetwork_ip_ranges_to_nat` - How NAT is configured per subnetwork. One of `ALL_SUBNETWORKS_ALL_IP_RANGES`, `ALL_SUBNETWORKS_ALL_PRIMARY_IP_RANGES` or `LIST_OF_SUBNETWORKS`.
* `nat_ip_allocate_option` - How external IPs are allocated for the NAT. `AUTO_ONLY` for NAT IPs allocated by Google Cloud Platform, or `MANUAL_ONLY` for user-allocated NAT IP addresses.
* `nat_ips` - Self-links of the NAT IPs.
* `min_ports_per_vm` - Minimum number of ports allocated to a VM from the NAT.
* `udp_idle_timeout_sec` - Timeout (in seconds) for UDP connections.
* `icmp_idle_timeout_sec` - Timeout (in seconds) for ICMP connections.
* `tcp_established_idle_timeout_sec` - Timeout (in seconds) for TCP established connections.
* `tcp_transitory_idle_timeout_sec` - Timeout (in seconds) for TCP transitory connections.
* `enable_endpoint_independent_mapping` - Whether endpoint independent mapping is enabled.

<a name="nested_md5_authentication_keys"></a>The `md5_authentication_keys` block supports:

* `name` - Name used to identify the key.
* `key` - Value of the key used for MD5 authentication. The API treats this value as input only. **Note**: This property is sensitive and will not be displayed in the plan.
