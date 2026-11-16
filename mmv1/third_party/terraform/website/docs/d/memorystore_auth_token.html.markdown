---
subcategory: "Memorystore"
description: |-
  Get information about an Auth Token for a Google Cloud Memorystore Token Auth User.
---

# google_memorystore_auth_token

Get information about an Auth Token associated with a Google Cloud Memorystore Token Auth User. For more details refer to the [API documentation](https://cloud.google.com/memorystore/docs/valkey/reference/rest/v1/projects.locations.instances.tokenAuthUsers.authTokens).

## Example Usage

```hcl
data "google_memorystore_auth_token" "qa" {
  instance = google_memorystore_instance.instance.instance_id
  user_id  = "default"
  token_id = "version-1"
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

* `token_id` -
  (Required)
  The ID of the auth token (e.g. `version-1`), or its full resource name.

* `location` -
  (Optional)
  The location of the Memorystore instance. If not provided, the provider region is used.

* `project` -
  (Optional)
  The ID of the project in which the resource belongs. If not provided, the provider project is used.

## Attributes Reference

See [`google_memorystore_auth_token`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/memorystore_auth_token) resource for details of all available attributes.
