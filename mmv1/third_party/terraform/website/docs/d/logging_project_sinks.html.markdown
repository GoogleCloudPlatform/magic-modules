---
subcategory: "Cloud (Stackdriver) Logging"
description: |-
  Get a list of the logging sinks in a project.
---

# google_logging_project_sinks

Get a list of the logging sinks owned by a project, including the built-in `_Required` and `_Default` sinks. Intercepting sinks owned by a parent folder or organization are not listed. For more information see
[the official documentation](https://cloud.google.com/logging/docs/export)
and
[API](https://cloud.google.com/logging/docs/reference/v2/rest/v2/projects.sinks/list).

## Example Usage

```hcl
data "google_logging_project_sinks" "all" {
  project = "my-project"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/sinks`

* `sinks` - The list of sinks in the project. Structure is [documented below](#nested_sinks).

<a name="nested_sinks"></a>The `sinks` block supports:

* `name` - The name of the logging sink.

* `destination` - The destination of the sink (or, in other words, where logs are written to).

* `filter` - The filter which is applied when exporting logs. Only log entries that match the filter are exported.

* `description` - The description of this sink.

* `disabled` - Whether this sink is disabled and it does not export any log entries.

* `writer_identity` - The identity associated with this sink. This identity must be granted write access to the configured `destination`. It is empty for sinks that route to a log bucket in the same project, which is the case for the built-in sinks by default.

* `bigquery_options` - Options that affect sinks exporting data to BigQuery. Structure is [documented below](#nested_bigquery_options).

* `exclusions` - Log entries that match any of the exclusion filters are not exported. Structure is [documented below](#nested_exclusions).

<a name="nested_bigquery_options"></a>The `bigquery_options` block supports:

* `use_partitioned_tables` - Whether [BigQuery's partition tables](https://cloud.google.com/bigquery/docs/partitioned-tables) are used.

<a name="nested_exclusions"></a>The `exclusions` block supports:

* `name` - A client-assigned identifier, such as `load-balancer-exclusion`.
* `description` - A description of this exclusion.
* `filter` - An advanced logs filter that matches the log entries to be excluded.
* `disabled` - Whether this exclusion is disabled and it does not exclude any log entries.
