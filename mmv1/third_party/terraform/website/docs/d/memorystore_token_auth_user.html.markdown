---
subcategory: "Memorystore"
description: |-
  Get information about a Google Cloud Memorystore (Valkey) Token Auth User.
---

# google_memorystore_token_auth_user

Get information about a Google Cloud Memorystore (Valkey) Token Auth User. For more details refer to the [API documentation](https://cloud.google.com/memorystore/docs/valkey/reference/rest/v1/projects.locations.instances.tokenAuthUsers).

## Example Usage

```hcl
data "google_memorystore_token_auth_user" "qa" {
  instance = "my-instance"
  user_id  = "my-user"
  location = "europe-west4"
}
```

## Argument Reference

The following arguments are supported:

* `instance` -
  (Required)
  The full name or ID of the Memorystore instance.

* `user_id` -
  (Required)
  The unique ID of the token auth user.

* `location` -
  (Optional)
  The location of the Memorystore instance. If not provided, provider default is used.

* `project` -
  (Optional)
  The project in which the resource belongs. If not provided, provider default is used.

## Attributes Reference

See [`google_memorystore_token_auth_user`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/memorystore_token_auth_user) resource for details of all available attributes.
