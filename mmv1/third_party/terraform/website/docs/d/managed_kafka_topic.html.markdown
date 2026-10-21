---
subcategory: "Managed Kafka"
description: |-
  A Managed Service for Apache Kafka topic.
---

# google_managed_kafka_topic

Use this data source to get information about a Managed Service for Apache Kafka topic.

To get more information about Managed Service for Apache Kafka topic, see:

* [API documentation](https://cloud.google.com/managed-service-for-apache-kafka/docs/reference/rest/v1/projects.locations.clusters.topics)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/managed-service-for-apache-kafka/docs)

## Example Usage

```hcl
data "google_managed_kafka_topic" "example" {
  topic_id = "my-topic"
  cluster  = "my-cluster"
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `topic_id` - (Required) The ID of the topic, which is the final component of the topic's name.

* `cluster` - (Required) The cluster name.

- - -

* `location` - (Optional) ID of the location of the Kafka resource. See https://cloud.google.com/managed-kafka/docs/locations for a list of supported locations. If it is not provided, the provider region or zone is used.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_managed_kafka_topic](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/managed_kafka_topic) resource for details of all the available attributes.
