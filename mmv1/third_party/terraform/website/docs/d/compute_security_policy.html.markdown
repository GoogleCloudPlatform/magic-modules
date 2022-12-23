---
subcategory: "Compute Engine"
page_title: "Google: google_compute_security_policy"
description: |-
  Get information about a Google Compute Security Policy.
---

# google\_compute\_security\_policy

To get more information about Google Compute Security Policy, see:

* [API documentation](https://cloud.google.com/compute/docs/reference/rest/beta/securityPolicies)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/armor/docs/configure-security-policies)

## Example Usage

```hcl
data "google_compute_security_policy" "foo" {
  name = "my-policy"
  project = "my-project"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the security policy.

* `project` - (Required) The project in which the resource belongs. If it is not provided, the provider project is used.

- - -

## Attributes Reference

See [google_compute_security_policy](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_security_policy) resource for details of the available attributes.
