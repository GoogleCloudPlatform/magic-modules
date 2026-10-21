---
subcategory: "Cloud Pub/Sub"
description: |-
  Get information about a Google Cloud Pub/Sub Schema.
---

# google_pubsub_schema

Get information about a Google Cloud Pub/Sub Schema. For more information see
the [official documentation](https://cloud.google.com/pubsub/docs/schemas)
and [API](https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.schemas).

## Example Usage

```hcl
data "google_pubsub_schema" "my-pubsub-schema" {
  name = "my-pubsub-schema"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cloud Pub/Sub Schema.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_pubsub_schema](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/pubsub_schema) resource for details of the available attributes.
