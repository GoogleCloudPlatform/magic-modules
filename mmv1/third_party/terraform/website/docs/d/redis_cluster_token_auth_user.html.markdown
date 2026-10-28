---
subcategory: "Redis"
description: |-
  Get information about a Google Cloud Redis Cluster Token Auth User.
---

# google_redis_cluster_token_auth_user

Get information about a Google Cloud Redis Cluster Token Auth User. For more details refer to the [API documentation](https://cloud.google.com/memorystore/docs/cluster/reference/rest/v1/projects.locations.clusters).

## Example Usage

```hcl
data "google_redis_cluster_token_auth_user" "qa" {
  cluster = google_redis_cluster.cluster.id
  user_id = "my-user"
}
```

## Argument Reference

The following arguments are supported:

* `cluster` -
  (Required)
  The full name or ID of the Redis Cluster.

* `user_id` -
  (Required)
  The unique ID of the token auth user. Can also be set to `"default"` to inspect the cluster's built-in default user.

## Attributes Reference

See [`google_redis_cluster_token_auth_user`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/redis_cluster_token_auth_user) resource for details of all available attributes.
