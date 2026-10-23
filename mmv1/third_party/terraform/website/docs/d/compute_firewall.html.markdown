---
subcategory: "Compute Engine"
description: |-
  Fetches the details of a Compute Firewall.
---

# google_compute_firewall

Use this data source to get information about a Compute Firewall. For more details, see the [API documentation](https://cloud.google.com/compute/docs/reference/rest/v1/firewalls).

## Example Usage

```hcl
data "google_compute_firewall" "default" {
  name = "my-firewall-rule"
}
```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  The name of the firewall rule.

- - -

* `project` -
  (Optional)
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_compute_firewall](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_firewall) resource for details of all the available attributes.
