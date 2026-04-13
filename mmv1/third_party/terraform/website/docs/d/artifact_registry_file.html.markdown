---
subcategory: "Artifact Registry"
description: |-
  Downloads a file from a Google Artifact Registry repository.
---

# google_artifact_registry_file

Downloads a single file from a Google Artifact Registry repository to a local
path and exposes its metadata and content hashes. Applies to file-based
Artifact Registry formats (Generic, Maven, npm, Python, Apt, Yum, Go). For
Docker/OCI images, use
[`google_artifact_registry_docker_image`](./artifact_registry_docker_image.html.markdown).

To get more information about Artifact Registry files, see:

* [API documentation](https://cloud.google.com/artifact-registry/docs/reference/rest/v1/projects.locations.repositories.files)

## Example Usage

```hcl
data "google_artifact_registry_file" "example" {
  location      = "us-central1"
  repository_id = "my-generic-repo"
  file_id       = "my-package:1.0.0:my-artifact.tar.gz"
  output_path   = "${path.module}/tmp/my-artifact.tar.gz"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The location of the repository.
* `repository_id` - (Required) The ID of the repository.
* `file_id` - (Required) The Artifact Registry file ID. For Generic repositories this is `<package>:<version>:<filename>`; for other formats refer to the file listing in the API. Slashes and other reserved characters are URL-encoded by the provider.
* `output_path` - (Required) Local filesystem path where the downloaded bytes are written. Parent directories are created if missing.
* `project` - (Optional) The project in which the repository lives. Defaults to the provider project.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `name` - The fully-qualified file resource name (`projects/.../files/...`).
* `size_bytes` - Size of the file in bytes, as reported by Artifact Registry.
* `hashes` - Map of hash type (e.g. `SHA256`, `MD5`) to the corresponding hash value reported by Artifact Registry.
* `create_time` - Creation time (RFC 3339).
* `update_time` - Last update time (RFC 3339).
* `output_sha256` - Hex-encoded SHA-256 of the downloaded file contents.
* `output_base64sha256` - Base64-encoded SHA-256 of the downloaded file contents.
