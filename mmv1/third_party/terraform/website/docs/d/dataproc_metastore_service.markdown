---
subcategory: "Dataproc metastore"
description: |-
  Get a Dataproc Metastore Service from Google Cloud
---

# google_dataproc_metastore_service

Get a Dataproc Metastore service from Google Cloud by its id and location.
~> It may take a while for the attached tag bindings to be deleted after the service is scheduled to be deleted.

## Example Usage

```tf
data "google_dataproc_metastore_service" "foo" {
  service_id = "foo-bar"
  location   = "global"  
}
```

To create a service with a tag

```tf
resource "google_dataproc_metastore_service" "my-service"{
  service_id = "your-service-id"
  location = "global"
  tags = {"tagKeys/281478409127147" : "tagValues/281479442205542"}
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required) The ID of the metastore service.
* `location` - (Required) The location where the metastore service resides.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.
* `tags` - (Optional) A map of resource manager tags. Resource manager tag keys and values have the same definition as resource manager tags. Keys must be in the format tagKeys/{tag_key_id}, and values are in the format tagValues/{tag_value_id}. The field is ignored when empty. The field is immutable and causes resource replacement when mutated.

## Attributes Reference

See [google_dataproc_metastore_service](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/dataproc_metastore_service) resource for details of all the available attributes.
