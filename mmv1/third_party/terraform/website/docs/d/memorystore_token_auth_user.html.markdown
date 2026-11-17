---
subcategory: "Memorystore"
description: |-
  Get information about a Google Cloud Memorystore Token Auth User.
---

# google_memorystore_token_auth_user

Get information about a Google Cloud Memorystore Token Auth User. For more details refer to the [API documentation](https://cloud.google.com/memorystore/docs/valkey/reference/rest/v1/projects.locations.instances.tokenAuthUsers).

## Example Usage

```hcl
data "google_memorystore_token_auth_user" "qa" {
  instance = google_memorystore_instance.instance.instance_id
  user_id  = "my-user"
  location = "europe-west4"
}
```

## Argument Reference

The following arguments are supported:

* `instance` -
  (Required)
  The ID or resource name of the Memorystore instance.

* `user_id` -
  (Required)
  The unique ID of the token auth user. Can also be set to `"default"` to inspect the instance's built-in default user.

* `location` -
  (Optional)
  The location of the Memorystore instance. When `instance` is provided as a full resource name (`projects/{project}/locations/{location}/instances/{instance}`), `location` is inferred from `instance` if not explicitly set, otherwise falling back to the provider region.

* `project` -
  (Optional)
  The ID of the project in which the resource belongs. When `instance` is provided as a full resource name (`projects/{project}/locations/{location}/instances/{instance}`), `project` is inferred from `instance` if not explicitly set, otherwise falling back to the provider project.

## Attributes Reference

See [`google_memorystore_token_auth_user`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/memorystore_token_auth_user) resource for details of all available attributes.
