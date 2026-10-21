---
subcategory: "AlloyDB"
description: |-
  Fetches the details of an available AlloyDB database user.
---

# google_alloydb_user

Use this data source to get information about an available AlloyDB database user. For more details refer to the [API docs](https://cloud.google.com/alloydb/docs/reference/rest/v1/projects.locations.clusters.users).

## Example Usage

```hcl
data "google_alloydb_user" "default" {
  cluster_id = "my-cluster"
  user_id    = "my-user"
  location   = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` -
  (Required)
  The ID of the alloydb cluster that the user belongs to.

* `user_id` -
  (Required)
  The database role name of the user.

* `location` -
  (Optional)
  The canonical ID of the location. If it is not provided, the provider location is used. For example: `us-central1`.

* `project` -
  (Optional)
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_alloydb_user](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/alloydb_user) resource for details of all the available attributes.
