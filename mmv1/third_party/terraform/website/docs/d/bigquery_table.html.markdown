---
subcategory: "BigQuery"
description: |-
  A datasource to retrieve a specific table in a dataset.
---

# `google_bigquery_table`

Get a specific table in a BigQuery dataset. For more information see
the [official documentation](https://cloud.google.com/bigquery/docs)
and [API](https://cloud.google.com/bigquery/docs/reference/rest/v2/tables/get).

## Example Usage

```hcl
data "google_bigquery_table" "table" {
  project    = "my-project"
  dataset_id = "my-bq-dataset"
  table_id   = "my-table"
}
```

## Argument Reference

The following arguments are supported:

* `dataset_id` - (Required) The dataset ID.

* `table_id` - (Required) The table ID.

* `project` - (Optional) The ID of the project in which the resource belongs.
  If it is not provided, the provider project is used.

## Attributes Reference

See [google_bigquery_table](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/bigquery_table#attributes-reference) resource for details of the available attributes.
