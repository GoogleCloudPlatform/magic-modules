---
subcategory: "Serverless VPC Access"
page_title: "Google: google_vpc_access_connector"
description: |-
  Get a Serverless VPC Access connector.
---

# google\_vpc\_access\_connector

Get a Serverless VPC Access connector.

To get more information about Connector, see:

* [API documentation](https://cloud.google.com/vpc/docs/reference/vpcaccess/rest/v1/projects.locations.connectors)
* How-to Guides
    * [Configuring Serverless VPC Access](https://cloud.google.com/vpc/docs/configure-serverless-vpc-access)

## Example Usage

```hcl
data "google_vpc_access_connector" "sample" {
  name = "vpc-con"
}

resource "google_vpc_access_connector" "connector" {
  name          = "vpc-con"
  ip_cidr_range = "10.8.0.0/28"
  network       = "default"
  region        = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the resource.

## Attributes Reference

See [google_vpc_access_connector](https://www.terraform.io/docs/providers/google/r/vpc_access_connector.html) resource for details of the available attributes.
