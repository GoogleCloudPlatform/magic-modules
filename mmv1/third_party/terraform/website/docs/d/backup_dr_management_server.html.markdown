---
subcategory: "BackupDR"
description: |-
  Get information about a Backupdr Management server.
---

# google\_backup\_dr\_management\_server

Get information about a Google Backup DR Management server.

## Example Usage

```hcl
data "google_backup_dr_management_server" "my-backup-dr-management-server" {
   location =  "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required.) The region in which the management server resource belongs.

- - -

## Attributes Reference

See [google_backupdr_management_server](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/backup_dr_management_server) resource for details of the available attributes.
