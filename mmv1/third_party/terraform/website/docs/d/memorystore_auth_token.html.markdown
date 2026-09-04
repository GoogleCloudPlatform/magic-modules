---
subcategory: "Memorystore"
description: |-
  Get information about an Auth Token for a Google Cloud Memorystore (Valkey) Token Auth User.
---

# google_memorystore_auth_token

Get information about an Auth Token associated with a Google Cloud Memorystore (Valkey) Token Auth User. For more details refer to the [API documentation](https://cloud.google.com/memorystore/docs/valkey/reference/rest/v1/projects.locations.instances.tokenAuthUsers.authTokens).

## Example Usage

```hcl
data "google_memorystore_auth_token" "qa" {
  token_auth_user = google_memorystore_token_auth_user.user.name
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

See [`google_memorystore_auth_token`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/memorystore_auth_token) resource for details of all available attributes.
