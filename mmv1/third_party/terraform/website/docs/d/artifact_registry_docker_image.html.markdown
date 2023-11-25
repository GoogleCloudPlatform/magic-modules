---
subcategory: "Artifact Registry"
description: |-
  Get information about a Docker Image within a Google Artifact Registry Repository.
---

# google\_artifact\_registry\_docker\_image

This data source fetches information from a provided Artifact Registry repository, including the fully qualified name and URI for an image, based on a the latest version of image name and optional digest or tag.

~> **Note**
Requires one of the following OAuth scopes: `https://www.googleapis.com/auth/cloud-platform` or `https://www.googleapis.com/auth/cloud-platform.read-only`.

## Example Usage

```hcl
data "google_artifact_registry_docker_image" "my-image" {
  repository = "my-repository"
  location   = "my-location"
  image      = "my-image"
  tag        = "my-tag"
}

resource "google_cloud_run_v2_service" "default" {
 # ...
 
  template {
    containers {
      image = data.google_artifact_registry_docker_image.my-image.self_link
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `repository` - (Required) The repository name.

* `location` - (Required) The location of the artifact registry repository. For example, "us-west1".

* `image` - (Required) The image name to fetch.

- - -

* `project` - (Optional) The project ID in which the resource belongs. If it
    is not provided, the provider project is used.

* `digest` - (Optional) The image digest to fetch.  This cannot be used if `tag` is provided.

* `tag` - (Optional) The tag of the version of the image to fetch. This cannot be used if `digest` is provided.

If a `digest` or `tag` is not provided, then the last updated version of the image will be fetched.

## Attributes Reference

The following computed attributes are exported:

* `name` - The fully qualified name of the fetched image.  This name has the form: `projects/{{project}}/locations/{{location}}/repository/{{repository}}/dockerImages/{{docker_image}}`. For example, 
```
projects/test-project/locations/us-west4/repositories/test-repo/dockerImages/nginx@sha256:e9954c1fc875017be1c3e36eca16be2d9e9bccc4bf072163515467d6a823c7cf
```

* `self_link` - The URI to access the image.  For example, 
```
us-west4-docker.pkg.dev/test-project/test-repo/nginx@sha256:e9954c1fc875017be1c3e36eca16be2d9e9bccc4bf072163515467d6a823c7cf
```

* `tags` - A list of all tags associated with the image.

* `image_size_bytes` - Calculated size of the image in bytes.

* `media_type` - Media type of this image, e.g. `application/vnd.docker.distribution.manifest.v2+json`. 

* `upload_time` - The time, as a RFC 3339 string, the image was uploaded. For example, `2014-10-02T15:01:23.045123456Z`.

* `build_time` - The time, as a RFC 3339 string, this image was built. 

* `update_time` - The time, as a RFC 3339 string, this image was updated.
