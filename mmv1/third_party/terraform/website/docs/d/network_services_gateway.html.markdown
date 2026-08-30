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
  location = "global"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Gateway resource.

* `location` - (Required) The location of the gateway. Use `global` for OPEN_MESH gateways.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_network_services_gateway](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/network_services_gateway) resource for details of all the available attributes.
