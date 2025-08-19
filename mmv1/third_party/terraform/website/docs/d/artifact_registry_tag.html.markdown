---
subcategory: "Artifact Registry"
description: |-
  Get information about a tag within a Google Artifact Registry repository.
---

# google_artifact_registry_tag
This data source fetches information of a tag from a provided Artifact Registry repository.

## Example Usage

```hcl
data "google_artifact_registry_tags" "my_tags" {
  location      = "us-central1"
  repository_id = "example-repo"
  package_name  = "example-package"
  tag_name      = "latest"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The location of the artifact registry.

* `repository_id` - (Required) The last part of the repository name to fetch from.

* `package_name` - (Required) The name of the package.

* `tag_name` - (Required) The name of the tag.

* `project` - (Optional) The project ID in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

The following computed attributes are exported:

* `name` - The name of the tag, for example: `projects/p1/locations/us-central1/repositories/repo1/packages/pkg1/tags/tag1`. If the package part contains slashes, the slashes are escaped.

* `version` - The version of the tag.
