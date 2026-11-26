---
subcategory: "Service Directory"
description: |-
  Get a list of Service Directory namespaces in a location.
---

# google_service_directory_namespaces

Get a list of Service Directory namespaces in a location. For more information see
[the official documentation](https://cloud.google.com/service-directory/docs/configuring-service-directory#configuring_a_namespace)
and
[API](https://cloud.google.com/service-directory/docs/reference/rest/v1/projects.locations.namespaces/list).

## Example Usage

```hcl
data "google_service_directory_namespaces" "all" {
  project  = "my-project"
  location = "us-central1"
}
```

## Example Usage - listing namespaces in several locations

```hcl
data "google_service_directory_namespaces" "by_location" {
  for_each = toset(["us-central1", "europe-west1"])

  project  = "my-project"
  location = each.value
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The location to list namespaces in. Each data source reads one location; use `for_each` to read several.

* `filter` - (Optional) A filter on the `name` or the `labels` of the namespaces, as defined in the [REST API](https://cloud.google.com/service-directory/docs/reference/rest/v1/projects.locations.namespaces/list#query-parameters). For example, `labels.owner=sd` returns the namespaces that have the label `owner` set to `sd`. A filter on a field that does not exist returns no namespaces instead of an error.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/namespaces`

* `namespaces` - The list of namespaces in the location. Structure is [documented below](#nested_namespaces).

<a name="nested_namespaces"></a>The `namespaces` block supports:

* `name` - The resource name of the namespace, in the format `projects/*/locations/*/namespaces/*`.
* `namespace_id` - The ID of the namespace, which is the last part of `name`.
* `labels` - The labels associated with the namespace.
