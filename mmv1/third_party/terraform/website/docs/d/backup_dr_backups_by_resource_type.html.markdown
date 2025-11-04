---
subcategory: "Backup and DR Service"
description: |-
  Get information about Backup and DR backups.
---

# google_backup_dr_fetch_backups

A list of Backup and DR backups.

## Example Usage

```hcl
data "google_backup_dr_backups_by_resource_type" "my_backups" {
  location = "us-central1"
  resource_type = "sqladmin.googleapis.com/Instance"
}

output "backup_name" {
  value = data.google_backup_dr_backups_by_resource_type.my_backups.backups[0].name
}
```

## Argument Reference
------------------

The following arguments are supported:

*   location- (Required) The location of the backups.
    
*   resource\_type- (Required) The resource type to get the backups for.
    
*   project - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.
    

## Attributes Reference
--------------------

In addition to the arguments listed above, the following attributes are exported:

*   backups - A list of the backups found. Each element of this list has the following attributes:
    
    *   name- The full name of the backup.
        
    *   backup\_type- The type of the backup.
        
    *   state- The state of the backup.
        
    *   create\_time- The time when the backup was created.
