---
subcategory: "Redis"
description: |-
  Get information about an Auth Token for a Google Cloud Redis Cluster Token Auth User.
---

# google_redis_cluster_auth_token

Get information about an Auth Token associated with a Google Cloud Redis Cluster Token Auth User. For more details refer to the [API documentation](https://cloud.google.com/memorystore/docs/cluster/reference/rest/v1/projects.locations.clusters.tokenAuthUsers.authTokens).

## Example Usage

```hcl
data "google_redis_cluster_auth_token" "qa" {
  token_auth_user = google_redis_cluster_token_auth_user.user.name
  token_id        = "version-1"
}
```

## Argument Reference

The following arguments are supported:

* `token_auth_user` -
  (Required)
  The full resource name of the parent Token Auth User.

* `token_id` -
  (Required)
  The version ID of the generated auth token (e.g. `version-1`).

* `project` -
  (Optional)
  The project in which the resource belongs. If not provided, provider default is used.

## Attributes Reference

See [`google_redis_cluster_auth_token`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/redis_cluster_auth_token) resource for details of all available attributes.
