---
subcategory: "Backup and DR Service"
description: |-
  Get information about a Backup and DR BackupPlan.
---

# google_backup_dr_backup_plan

Use this data source to get information about an existing Backup and DR BackupPlan.

## Example Usage

```hcl
data "google_backup_dr_backup_plan" "my-backup-plan" {
  location       = "us-central1"
  backup_plan_id = "bp-id"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The location in which the `BackupPlan` resource belongs.
* `backup_plan_id` - (Required) The ID of the `BackupPlan` resource.
* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

- - -

## Attributes Reference

See [google_backup_dr_backup_plan](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/backup_dr_backup_plan) resource for details of the available attributes.
