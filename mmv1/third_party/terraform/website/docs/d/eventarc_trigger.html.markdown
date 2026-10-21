---
subcategory: "Eventarc"
description: |-
  Get information about a Google Cloud Eventarc Trigger.
---

# google_eventarc_trigger

Get information about a Google Cloud Eventarc Trigger. For more information see
the [official documentation](https://cloud.google.com/eventarc/standard/docs/overview#eventarc-triggers)
and [API](https://cloud.google.com/eventarc/docs/reference/rest/v1/projects.locations.triggers).

## Example Usage

```hcl
data "google_eventarc_trigger" "my-trigger" {
  name     = "my-trigger"
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Eventarc Trigger.

- - -

* `location` - (Optional) The location for the resource. If it is not provided,
    the provider region or zone is used.

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_eventarc_trigger](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/eventarc_trigger) resource for details of the available attributes.
