---
subcategory: "Certificate Manager"
description: |-
  Fetches the details of a Certificate Manager DNS Authorization.
---

# google_certificate_manager_dns_authorization

Use this data source to get information about a Certificate Manager DNS Authorization. For more details, see the [API documentation](https://cloud.google.com/certificate-manager/docs/reference/certificate-manager/rest/v1/projects.locations.dnsAuthorizations).

## Example Usage

```hcl
data "google_certificate_manager_dns_authorization" "default" {
  name     = "my-dns-auth"
  location = "global"
}
```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  The name of the DNS Authorization.

* `domain` -
  (Required)
  The name of the DNS Authorization.

* `location` -
  (Optional)
  The Certificate Manager location. If not specified, "global" is used.

* `project` -
  (Optional)
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_certificate_manager_dns_authorization](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/certificate_manager_dns_authorization) resource for details of all the available attributes.
