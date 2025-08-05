---
subcategory: "Artifact Registry"
description: |-
  Get information about versions within a Google Artifact Registry package.
---

# google_artifact_registry_versions

Get information about Artifact Registry versions.
See [the official documentation](https://cloud.google.com/artifact-registry/docs/overview)
and [API](https://cloud.google.com/artifact-registry/docs/reference/rest/v1/projects.locations.repositories.packages.versions/list).

## Example Usage

```hcl
data "google_artifact_registry_versions" "my_versions" {
  location      = "us-central1"
  repository_id = "example-repo"
  package_name  = "example-package"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The location of the artifact registry.

* `repository_id` - (Required) The last part of the repository name to fetch from.

* `package_name` - (Required) The name of the package.

* `filter` - (Optional) An expression for filtering the results of the request. Filter rules are case insensitive. The fields eligible for filtering are `name` and `version`. Further information can be found in the [REST API](https://cloud.google.com/artifact-registry/docs/reference/rest/v1/projects.locations.repositories.packages.versions/list#query-parameters).

* `view` - (Optional) The view, which determines what version information is returned in a response. Possible values are `"BASIC"` and `"FULL"`. Defaults to `"BASIC"`.

* `project` - (Optional) The project ID in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

The following attributes are exported:

* `versions` - A list of all retrieved Artifact Registry versions. Structure is [defined below](#nested_versions).

<a name="nested_versions"></a>The `versions` block supports:

* `name` - The name of the version, for example: `projects/p1/locations/us-central1/repositories/repo1/packages/pkg1/versions/version1`. If the package part contains slashes, the slashes are escaped.

* `description` - Description of the version, as specified in its metadata.

* `related_tags` - A list of related tags. Will contain up to 100 tags that reference this version.

* `create_time` - The time, as a RFC 3339 string, this package was created. 

* `update_time` - The time, as a RFC 3339 string, this package was last updated. This includes publishing a new version of the package.

* `annotations` - Client specified annotations.
