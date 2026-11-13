---
subcategory: "Compute Engine"
description: |-
  Get a specific Compute Engine DHCP options configuration within a region.
---

# google\_compute\_dhcp\_options\_config

Get a specific Compute Engine DHCP options configuration within a region.

To get more information about DhcpOptionsConfig, see:

* [API documentation](https://cloud.google.com/compute/docs/reference/rest/beta/dhcpOptionsConfigs)
* How-to Guides
    * [VPC network overview](https://cloud.google.com/vpc/docs/vpc)

~> **Warning:** This data source is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](../guides/provider_versions.html.markdown) for more details on beta resources.

## Example Usage

```hcl
data "google_compute_dhcp_options_config" "default" {
  provider = google-beta

  name    = "my-dhcp-config"
  region  = "us-central1"
  project = "my-project"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DHCP options configuration to retrieve.

- - -

* `region` - (Optional) The region in which the DHCP options configuration resides.
  If it is not provided, the provider region is used.

* `project` - (Optional) The ID of the project in which the resource belongs.
  If it is not provided, the provider project is used.

## Attributes Reference

See [google_compute_dhcp_options_config](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_dhcp_options_config#attributes-reference) resource for details of all the available attributes.
