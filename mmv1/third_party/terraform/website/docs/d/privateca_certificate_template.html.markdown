---
subcategory: "Certificate Authority Service"
description: |-
  Get info about a Google CAS CertificateTemplate.
---

# google_privateca_certificate_template

Get info about a Google Cloud Certificate Authority Service CertificateTemplate.

To get more information about CertificateTemplate, see:

* [API documentation](https://cloud.google.com/certificate-authority-service/docs/reference/rest)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/certificate-authority-service)
    * [Understanding Certificate Templates](https://cloud.google.com/certificate-authority-service/docs/certificate-template)
    * [Common configurations and Certificate Profiles](https://cloud.google.com/certificate-authority-service/docs/certificate-profile)

## Example Usage

```hcl
data "google_privateca_certificate_template" "default" {
  name     = "my-template"
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The resource name for this CertificateTemplate.

- - -

* `location` - (Optional) The location for the resource. If it is not provided, the provider region or zone is used.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_privateca_certificate_template](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/privateca_certificate_template) resource for details of all the available attributes.
