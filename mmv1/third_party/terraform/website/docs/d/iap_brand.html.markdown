---
subcategory: "Identity-Aware Proxy"
description: |-
  Provides the information of the Identity Aware Proxy brand.
---
# google_iap_brand

Get info about a Google Cloud IAP Brand.

## Example Usage

```tf
data "google_project" "project" {
  project_id = "foobar"
}

data "google_iap_brand" "project" {
  project =  data.google_project.project.id
}

```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The name of the brand.

* `brand` - (Optional) The brand id which is the project number. By default, the data source retrieve the project number automatically if `brand` is not specified.

## Attributes Reference

See [google_iap_client](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/iap_brand) resource for details of the available attributes.
