---
subcategory: "Vertex AI Workbench"
description: |-
  Get information about a Vertex AI Workbench instance.
---

# google_workbench_instance

Use this data source to get information about a Vertex AI Workbench instance in a project and location.

To get more information about Instance, see:

* [API documentation](https://cloud.google.com/vertex-ai/docs/workbench/reference/rest/v2/projects.locations.instances)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/vertex-ai/docs/workbench/instances/introduction)

## Example Usage

```hcl
data "google_workbench_instance" "default" {
  name     = "my-workbench-instance"
  location = "us-east1-b"
}
```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  The name of the workbench instance.

* `location` -
  (Optional)
  The location (zone) where the workbench instance resides. If it is not provided, the provider location is used.

* `project` -
  (Optional)
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_workbench_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/workbench_instance) resource for details of all the available attributes.
