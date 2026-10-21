---
subcategory: "Certificate Authority Service"
description: |-
  Get info about a Google CAS CaPool.
---

# google_privateca_ca_pool

Get info about a Google Cloud Certificate Authority Service CaPool.

To get more information about CaPool, see:

* [API documentation](https://cloud.google.com/certificate-authority-service/docs/reference/rest/v1/projects.locations.caPools)
* How-to Guides
    * [Certificate Authority Service Overview](https://cloud.google.com/certificate-authority-service/docs/overview)

## Example Usage

```hcl
data "google_privateca_ca_pool" "default" {
  name     = "my-pool"
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name for this CaPool.

- - -

* `location` - (Optional) Location of the CaPool. A full list of valid locations can be found by running `gcloud privateca locations list`. If it is not provided, the provider region or zone is used.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_privateca_ca_pool](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/privateca_ca_pool) resource for details of all the available attributes.
