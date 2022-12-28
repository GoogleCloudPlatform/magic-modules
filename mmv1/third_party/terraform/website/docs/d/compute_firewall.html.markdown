---
subcategory: "Compute Engine"
page_title: "Google: google_compute_firewall"
description: |-
  Get information about a Google Compute Firewall.
---

# google\_compute\_firewall

To get more information about Google Compute Firewall, see:

* [API documentation](https://cloud.google.com/compute/docs/reference/rest/v1/firewalls/get)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/vpc/docs/firewalls)

## Example Usage

```hcl
data "google_compute_firewall" "foo" {
  name = "firewall-name"
  project = "my-project"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the compute firewall.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

- - -

## Attributes Reference

See [google_compute_firewall](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_firewall) resource for details of the available attributes.
