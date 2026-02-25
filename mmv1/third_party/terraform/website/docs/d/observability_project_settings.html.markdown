---
subcategory: "Cloud Observability"
description: |-
  Describes the Google Cloud Observability settings associated with a project.
---

# google_observability_project_settings

Describes the Google Cloud Observability Settings associated with a project.

To get more information about Observability Settings, see:

* [API documentation](https://docs.cloud.google.com/stackdriver/docs/reference/observability/api/rest)

## Example Usage - Observability Project Settings Basic

```hcl
data "google_observability_project_settings" "settings" {
  project  = "my-project-name"
  location = "global"
}
```

## Argument Reference

The following arguments are supported:

- - -

* `project` - (Required) The project for which to retrieve settings.

* `location` - (Required) The location of the settings.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/settings`

* `name` - The resource name of the settings.

* `default_storage_location` - The default storage location for new resources, e.g. buckets. Only valid for global location.

* `kms_key_name` - The default Cloud KMS key to use for new resources. Only valid for regional locations.

* `service_account_id` - The service account used by Cloud Observability for this project.
