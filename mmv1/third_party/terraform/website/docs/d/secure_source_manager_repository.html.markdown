---
subcategory: "Secure Source Manager"
description: |-
  Get information about a Secure Source Manager Repository.
---

# google_secure_source_manager_repository

Get information about a Secure Source Manager Repository.

To get more information about Secure Source Manager Repository, see:

* [API documentation](https://cloud.google.com/secure-source-manager/docs/reference/rest/v1/projects.locations.repositories)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/secure-source-manager/docs/overview)

## Example Usage

```hcl
data "google_secure_source_manager_repository" "default" {
  repository_id = "my-repository"
  location      = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `repository_id` - (Required) The ID for the Repository.

- - -

* `location` - (Optional) The location for the Repository. If it is not provided, the provider region/zone is used.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_secure_source_manager_repository](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/secure_source_manager_repository) resource for details of all the available attributes.
