---
subcategory: "Secret Manager"
description: |-
  List all versions for a Secret Manager secret.
---

# google_secret_manager_secret_versions

List all versions of a Secret Manager secret.

## Example Usage

\`\`\`hcl
data "google_secret_manager_secret_versions" "versions" {
  secret = "my-secret"
}
\`\`\`

## Argument Reference

* `secret` - (Required) The secret to list versions for. Can be the secret_id or the full resource name.
* `project` - (Optional) The project in which the resource belongs. If not set, the provider project is used.
* `filter` - (Optional) Filter string for listing versions.

## Attributes Reference

* `versions` - A list of versions. Each version contains:
  * `name` - The full resource name of the version.
  * `version` - The version number.
  * `enabled` - Whether the version is enabled.
  * `create_time` - The time the version was created.
  * `destroy_time` - The time the version was destroyed (if applicable).
