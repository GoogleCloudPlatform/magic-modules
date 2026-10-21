---
subcategory: "Managed Kafka"
description: |-
  A Managed Service for Kafka Connect cluster.
---

# google_managed_kafka_connect_cluster

Use this data source to get information about a Managed Service for Kafka Connect cluster.

To get more information about Managed Service for Kafka Connect cluster, see:

* [API documentation](https://cloud.google.com/managed-service-for-apache-kafka/docs/reference/rest/v1/projects.locations.connectClusters)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/managed-service-for-apache-kafka/docs)

## Example Usage

```hcl
data "google_managed_kafka_connect_cluster" "example" {
  connect_cluster_id = "my-connect-cluster"
  location           = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `connect_cluster_id` - (Required) The ID of the Connect Cluster, which is the final component of the connect cluster's name.

- - -

* `location` - (Optional) ID of the location of the Kafka Connect resource. See https://cloud.google.com/managed-kafka/docs/locations for a list of supported locations. If it is not provided, the provider region or zone is used.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_managed_kafka_connect_cluster](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/managed_kafka_connect_cluster) resource for details of all the available attributes.
