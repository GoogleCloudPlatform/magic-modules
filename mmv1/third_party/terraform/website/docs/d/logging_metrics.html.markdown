---
subcategory: "Cloud (Stackdriver) Logging"
description: |-
  Get a list of the logs-based metrics in a project.
---

# google_logging_metrics

Get a list of the logs-based metrics in a project. For more information see
[the official documentation](https://cloud.google.com/logging/docs/logs-based-metrics)
and
[API](https://cloud.google.com/logging/docs/reference/v2/rest/v2/projects.metrics/list).

## Example Usage

```hcl
data "google_logging_metrics" "all" {
  project = "my-project"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/metrics`

* `metrics` - The list of logs-based metrics in the project.

See [google_logging_metric](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/logging_metric) resource for details of the available attributes on each metric. The `project` and `deletion_policy` arguments of that resource are not part of each metric.
