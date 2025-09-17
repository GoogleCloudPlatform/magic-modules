---
subcategory: "Tags"
description: |-
  A TagBindingCollection represents a collection of tag bindings directly bound to a cloud resource.
---

# google_tags_tag_binding_collection

A TagBindingCollection represents a collection of tag bindings directly bound to a cloud resource.

To get more information about TagBindingCollection, see:

* [API documentation](https://cloud.google.com/resource-manager/reference/rest/v3/locations.tagBindingCollections)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/resource-manager/docs/tags/tags-creating-and-managing)

## Example Usage - Cloud Run Service

To bind tags to a Cloud Run service:

```hcl
resource "google_project" "project" {
  project_id = "project_id"
  name       = "project_id"
  org_id     = "123456789"
}

resource "google_tags_tag_key" "key1" {
  parent      = "organizations/123456789"
  short_name  = "keyname1"
  description = "For keyname1 resources."
}

resource "google_tags_tag_value" "value1" {
  parent      = google_tags_tag_key.key1.id
  short_name  = "valuename1"
  description = "For valuename1 resources."
}

resource "google_tags_tag_key" "key2" {
  parent      = "organizations/123456789"
  short_name  = "keyname2"
  description = "For keyname2 resources."
}

resource "google_tags_tag_value" "value2" {
  parent      = google_tags_tag_key.key2.id
  short_name  = "valuename2"
  description = "For valuename2 resources."
}

resource "google_tags_tag_binding_collection" "bindingcollection" {
  full_resource_name    = "//run.googleapis.com/projects/${data.google_project.project.number}/locations/${google_cloud_run_service.default.location}/services/${google_cloud_run_service.default.name}"
  location  = "us-central1"
  tags      = {
    # Format: "{TagKey.namespaced_name}" = "{TagValue.short_name}"
    "${google_tags_tag_key.key1.namespaced_name}" = google_tags_tag_value.value1.short_name
    "${google_tags_tag_key.key2.namespaced_name}" = google_tags_tag_value.value2.short_name
  }
}
```

## Argument Reference

The following arguments are supported:


* `full_resource_name` -
  (Required)
  The full resource name of the resource to which the tags are bound. E.g. //cloudresourcemanager.googleapis.com/projects/123

* `tags` -
  (Required)
  A map of tag keys to values directly bound to this resource, specified in namespaced name format. E.g. "123/environment": "production"

* `location` -
  (Required)
  The location of the target resource. E.g. "global", "us-central1", "us-east2-c".

- - -



## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - Identifier of the TagBindingCollection resource, in the format `locations/{location}/tagBindingCollections/{encoded_full_resource_name}`

* `name` -
  The name of the TagBindingCollection, in the format: `locations/{location}/tagBindingCollections/{encoded_full_resource_name}`

* `active_tags` -
  The most recent state of all direct tags on the resource, as reported by the API. 
  This includes the tags configured through Terraform, the system, and other clients.


## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.

## Import


TagBindingCollection can be imported using any of these accepted formats:

* `locations/{{location}}/tagBindingCollections/{{encoded_full_resource_name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import TagBindingCollection using one of the formats above. For example:

```tf
import {
  id = "locations/{{location}}/tagBindingCollections/{{encoded_full_resource_name}}"
  to = google_tags_tag_binding_collection.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), TagBindingCollection can be imported using one of the formats above. For example:

```
$ terraform import google_tags_tag_binding_collection.default locations/{{location}}/tagBindingCollections/{{encoded_full_resource_name}}
```
