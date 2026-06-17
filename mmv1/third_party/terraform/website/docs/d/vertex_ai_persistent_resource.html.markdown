---
subcategory: "Vertex AI"
description: |-
  Get information about a Google Cloud Vertex AI Persistent Resource.
---

# data_source_google_vertex_ai_persistent_resource

Get information about a Google Cloud Vertex AI Persistent Resource.

## Example Usage

```terraform
data "google_vertex_ai_persistent_resource" "foo" {
  name     = "my-persistent-resource"
  location = "us-central1"
  project  = "my-project-id"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the PersistentResource.
* `location` - (Required) The location for the resource.
* `project` - (Optional) The project ID of the resource.

## Attributes Reference

See [google_vertex_ai_persistent_resource](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/vertex_ai_persistent_resource) resource for details of all the available attributes.
