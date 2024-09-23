---
subcategory: "Cloud Spanner"
description: |-
  Get a spanner instance from Google Cloud
---

# google_spanner_instance

Get a spanner instance from Google Cloud by its name.
~> It may take a while for the attached tag bindings to be deleted after the instance is scheduled to be deleted.

## Example Usage

```tf
data "google_spanner_instance" "foo" {
  name = "bar"
}
```

To create an instance with a tag

```tf
resource "spanner_instance" "my_instance" {
  name       = "My instance"
  instance_id = "your-instance-id"
  tags = {"tagKeys/281478409127147" : "tagValues/281479442205542"}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the spanner instance.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.
    
* `tags` - (Optional) A map of resource manager tags. Resource manager tag keys and values have the same definition as resource manager tags. Keys must be in the format tagKeys/{tag_key_id}, and values are in the format tagValues/{tag_value_id}. The field is ignored when empty. The field is immutable and causes resource replacement when mutated.

## Attributes Reference
See [google_spanner_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/spanner_instance) resource for details of all the available attributes.
