---
subcategory: "Data Lineage"
description: |-
  Configuration for Data Lineage.
---

# google_data_lineage_config

Get a configuration for Data Lineage.

To get more information about Config, see [Official Documentation](https://docs.cloud.google.com/dataplex/docs/about-data-lineage#control-lineage-ingestion)

## Example Usage

```hcl
data "google_data_lineage_config" "default" {
  parent = "projects/my-project-name"
  location = "global"
}
```

## Argument Reference

The following arguments are supported:

* `parent` -
  (Required)
  Parent scope for the config.
  Format: projects/{project-id|project-number} or folders/{folder-number} or organizations/{organization-number}.

* `location` -
  (Required)
  The region of the data lineage configuration for integration.

## Attributes Reference

See [google_data_lineage_config](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/data_lineage_config) resource for details of all the available attributes.
