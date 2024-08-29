---
subcategory: "BigQuery"
description: |-
  A datasource to retrieve a list of tables in a dataset.
---

# `google_bigquery_tables`

Get a list of tables in a BigQuery dataset. For more information see
the [official documentation](https://cloud.google.com/bigquery/docs)
and [API](https://cloud.google.com/bigquery/docs/reference/rest/v2/tables).

## Example Usage

```hcl
data "google_bigquery_tables" "tables" {
  dataset_id = "my-bq-dataset"
  project = "my-project"
}
```

## Argument Reference

The following arguments are supported:

* `dataset_id` - (Required) The dataset ID.

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

## Attributes Reference

* `tables` - Map of tables contained within the provided `dataset_id`.

