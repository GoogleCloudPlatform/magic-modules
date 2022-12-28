---
subcategory: "Cloud SQL"
page_title: "Google: google_sql_databases"
description: |-
  Get a list of databases in a Cloud SQL database instance.
---

# google\_sql\_databases

Use this data source to get information about a list of databases in a Cloud SQL instance, you can also apply some filter over this list of databases.

## Example Usage


```hcl
data "google_sql_databases" "qa" {
  instance = google_sql_database_instance.main.name
  filter{
    "name" = "name"
    "values" = ["db-*"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance` - (required) The name of the Cloud SQL database instance in which the database belongs.

* `project` - (optional) The ID of the project in which the instance belongs.

The optional `filters` sublist supports:

* `name` - (Required) Name of the filter. supported values include `name`, `charset`, `collation`.

* `values` - (Optional) To include databases which matches at least on of the regex provided in the list for the chosen filter.

* `exclude_values` - (Optional) To exclude databases which matches at least on of the regex provided in the list for the chosen filter.

This filtering is done locally on what GCP returns, and could have a performance impact if the count of databases is large.

## Attributes Reference
See [google_sql_database](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/sql_database) resource for details of all the available attributes.
