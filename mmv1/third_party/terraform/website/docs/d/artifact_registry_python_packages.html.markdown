---
subcategory: "Artifact Registry"
description: |-
  Get information about Python packages within a Google Artifact Registry repository.
---

# google_artifact_registry_python_packages

Get information about Artifact Registry Python packages.
See [the official documentation](https://cloud.google.com/artifact-registry/docs/python)
and [API](https://cloud.google.com/artifact-registry/docs/reference/rest/v1/projects.locations.repositories.pythonPackages/list).

## Example Usage

```hcl
data "google_artifact_registry_python_packages" "my_packages" {
  location      = "us-central1"
  repository_id = "example-repo"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The location of the Artifact Registry repository.

* `repository_id` - (Required) The last part of the repository name to fetch from.

* `project` - (Optional) The project ID in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

The following attributes are exported:

* `python_packages` - A list of all retrieved Artifact Registry Python packages. Structure is [defined below](#nested_python_packages).

<a name="nested_python_packages"></a>The `python_packages` block supports:

* `name` - The fully qualified name of the fetched package.  This name has the form: `projects/{{project}}/locations/{{location}}/repository/{{repository_id}}/pythonPackages/{{pythonPackage}}`. For example, `projects/example-project/locations/us-central1/repository/example-repo/pythonPackages/my-test-package:0.0.1`

* `package_name` - Extracted short name of the package (last part of `name`, without version). For example, from `.../my-test-package:0.0.1` → `my-test-package`.

* `version` - Version of this package.

* `create_time` - The time, as a RFC 3339 string, this package was created. 

* `update_time` - The time, as a RFC 3339 string, this package was updated.
