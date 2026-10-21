---
subcategory: "Secure Source Manager"
description: |-
  Get information about a Secure Source Manager Instance.
---

# google_secure_source_manager_instance

Get information about a Secure Source Manager Instance.

To get more information about Secure Source Manager Instance, see:

* [API documentation](https://cloud.google.com/secure-source-manager/docs/reference/rest/v1/projects.locations.instances)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/secure-source-manager/docs/create-instance)

## Example Usage

```hcl
data "google_secure_source_manager_instance" "default" {
  instance_id = "my-instance"
  location    = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The name for the Instance.

- - -

* `location` - (Optional) The location for the Instance. If it is not provided, the provider region/zone is used.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_secure_source_manager_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/secure_source_manager_instance) resource for details of all the available attributes.
