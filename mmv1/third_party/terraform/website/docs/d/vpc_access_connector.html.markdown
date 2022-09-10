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

* `name` -
  (Required)
  The name of the resource (Max 25 characters).

* `network` -
  (Optional)
  Name or self_link of the VPC network. Required if `ip_cidr_range` is set.

* `ip_cidr_range` -
  (Optional)
  The range of internal addresses that follows RFC 4632 notation. Example: `10.132.0.0/28`.

* `min_throughput` -
  (Optional)
  Minimum throughput of the connector in Mbps. Default and min is 200.

* `max_throughput` -
  (Optional)
  Maximum throughput of the connector in Mbps, must be greater than `min_throughput`. Default is 300.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{region}}/connectors/{{name}}`

* `state` -
  State of the VPC access connector.
