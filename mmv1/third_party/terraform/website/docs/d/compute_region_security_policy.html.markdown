---
subcategory: "Compute Engine"
description: |-
  Fetches the details of a Compute Region Security Policy.
---

# google_compute_region_security_policy

Use this data source to get information about a Compute Region Security Policy. For more details, see the [API documentation](https://cloud.google.com/compute/docs/reference/rest/v1/regionSecurityPolicies).

## Example Usage

```hcl
data "google_compute_region_security_policy" "default" {
  name   = "my-region-security-policy"
  region = "us-west2"
}
```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  The name of the Region Security Policy.

* `region` -
  (Optional)
  The region in which the Region Security Policy resides. If not specified, the provider region is used.

* `project` -
  (Optional)
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_compute_region_security_policy](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_region_security_policy) resource for details of all the available attributes.
