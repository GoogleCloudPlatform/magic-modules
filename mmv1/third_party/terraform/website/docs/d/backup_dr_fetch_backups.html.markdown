---
subcategory: "Backup and DR Service"
description: |-
  Get information about Backup and DR backups.
---

# google_backup_dr_fetch_backups

A list of Backup and DR backups.

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

## Example Usage

```hcl
data "google_backup_dr_fetch_backups" "my_backups" {
  location = "us-central1"
  resource_type = "sqladmin.googleapis.com/Instance"
}

output "backup_name" {
  value = data.google_backup_dr_fetch_backups.my_backups.backups[0].name
}
```
