---
subcategory: "Network Services"
description: |-
  Get information about a Google Cloud Network Services Gateway.
---

# google_network_services_gateway

Get information about a Google Cloud Network Services Gateway. This includes
OPEN_MESH gateways and Secure Web Proxy (`SECURE_WEB_GATEWAY`) gateways.

To get more information about Gateway, see:

* [API documentation](https://cloud.google.com/traffic-director/docs/reference/network-services/rest/v1/projects.locations.gateways)
* How-to Guides
    * [Secure Web Proxy](https://cloud.google.com/secure-web-proxy/docs)

## Example Usage

```hcl
data "google_network_services_gateway" "default" {
  name     = "my-gateway"
  location = "europe-west2"
}

output "gateway_ip" {
  value = data.google_network_services_gateway.default.addresses
}

output "gateway_ports" {
  value = data.google_network_services_gateway.default.ports
}

output "gateway_type" {
  value = data.google_network_services_gateway.default.type
}

output "gateway_network" {
  value = data.google_network_services_gateway.default.network
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Gateway resource.

* `location` - (Required) The location of the gateway. Use a region such as
    `europe-west2` for Secure Web Proxy (`SECURE_WEB_GATEWAY`) gateways, or
    `global` for `OPEN_MESH` gateways.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - An identifier for the resource with format `projects/{{project}}/locations/{{location}}/gateways/{{name}}`.

* `self_link` - Server-defined URL of this resource.

* `type` - The type of the customer managed gateway. Possible values are `OPEN_MESH` and `SECURE_WEB_GATEWAY`.

* `addresses` - Zero or one IPv4 or IPv6 address on which the Gateway receives traffic.
    When no address is provided, an IP from the subnetwork is allocated.
    This field only applies to gateways of type `SECURE_WEB_GATEWAY`.
    Gateways of type `OPEN_MESH` listen on `0.0.0.0` for IPv4 and `::` for IPv6.

* `ports` - One or more port numbers (1-65535) on which the Gateway receives traffic.

* `all_ports` - If true, the gateway listens on all ports (1-65535). Configurable only for `SECURE_WEB_GATEWAY`.

* `network` - The relative resource name of the VPC network using this configuration, for example `projects/*/global/networks/network-1`. Specific to `SECURE_WEB_GATEWAY`.

* `subnetwork` - The relative resource name of the subnetwork in which this Secure Web Proxy is allocated, for example `projects/*/regions/us-central1/subnetworks/network-1`. Specific to `SECURE_WEB_GATEWAY`.

* `gateway_security_policy` - A fully-qualified GatewaySecurityPolicy URL. Specific to `SECURE_WEB_GATEWAY`.

* `certificate_urls` - Fully-qualified Certificate URLs presented by the proxy when establishing a TLS connection. Specific to `SECURE_WEB_GATEWAY`.

* `server_tls_policy` - A fully-qualified ServerTLSPolicy URL. If empty, TLS termination is disabled.

* `routing_mode` - Routing mode of the Gateway. Possible values are `NEXT_HOP_ROUTING_MODE` and `EXPLICIT_ROUTING_MODE`. Configurable only for `SECURE_WEB_GATEWAY`.

* `ip_version` - IP version used by this gateway. Possible values are `IPV4` and `IPV6`.

* `scope` - Scope used to merge configuration across multiple Gateway instances.

* `description` - A free-text description of the resource.

* `envoy_headers` - Whether Envoy inserts internal debug headers into upstream requests. Possible values are `NONE` and `DEBUG_HEADERS`.

* `allow_global_access` - If true, the gateway allows traffic from clients outside the region where the gateway is located. Configurable only for `SECURE_WEB_GATEWAY`.

* `create_time` - The timestamp when the resource was created.

* `update_time` - The timestamp when the resource was updated.

* `labels` - Set of label tags associated with the Gateway resource. This field is non-authoritative; see `effective_labels` for all labels present on the resource.

* `terraform_labels` - The combination of labels configured directly on the resource and default labels configured on the provider.

* `effective_labels` - All labels present on the resource in GCP, including labels configured through Terraform, other clients, and services.
