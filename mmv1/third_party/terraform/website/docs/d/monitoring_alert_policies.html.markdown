---
subcategory: "Cloud (Stackdriver) Monitoring"
description: |-
  Get a list of the Cloud Monitoring alert policies in a project.
---

# google_monitoring_alert_policies

Get a list of the Cloud Monitoring alert policies in a project. For more information see
[the official documentation](https://cloud.google.com/monitoring/alerts/)
and [API](https://cloud.google.com/monitoring/api/ref_v3/rest/v3/projects.alertPolicies/list).

## Example Usage

```hcl
data "google_monitoring_alert_policies" "all" {
  project = "my-project"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The ID of the project whose alert policies are listed. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/alertPolicies`

* `alert_policies` - The list of alert policies in the project. Each policy has the attributes of the
    [google_monitoring_alert_policy](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/monitoring_alert_policy)
    resource, including its `id`, except `project` and `deletion_policy`.
