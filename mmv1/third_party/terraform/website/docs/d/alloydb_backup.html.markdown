---
subcategory: "AlloyDB"
description: |-
  Fetches the details of an available AlloyDB backup.
---

# google_alloydb_backup

Use this data source to get information about an available AlloyDB backup. For more details refer to the [API docs](https://cloud.google.com/alloydb/docs/reference/rest/v1/projects.locations.backups).

## Example Usage

```hcl
data "google_alloydb_backup" "default" {
  backup_id = "my-backup"
  location  = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `backup_id` -
  (Required)
  The ID of the alloydb backup.

* `location` -
  (Optional)
  The canonical ID of the location. If it is not provided, the provider location is used. For example: `us-central1`.

* `project` -
  (Optional)
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_alloydb_backup](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/alloydb_backup) resource for details of all the available attributes.
