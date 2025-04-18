---
subcategory: "Memorystore (Redis)"
description: |-
  Fetches the details of a Redis Cluster.
---

# google_redis_cluster

Use this data source to get information about a Redis Cluster. For more details, see the [API documentation](https://cloud.google.com/memorystore/docs/cluster/reference/rest/v1/projects.locations.clusters).

## Example Usage

```hcl
data "google_redis_cluster" "default" {
  name   = "my-redis-cluster"
  region = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  The name of the Redis cluster.

* `region` -
  (Required)
  The region of the Redis cluster.

* `project` - 
  (optional) 
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_redis_cluster](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/redis_cluster) resource for details of all the available attributes.
