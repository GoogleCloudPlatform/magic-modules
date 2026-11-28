---
subcategory: "Cloud (Stackdriver) Monitoring"
description: |-
  Get a list of the Cloud Monitoring notification channels in a project.
---

# google_monitoring_notification_channels

Get a list of the Cloud Monitoring notification channels in a project. For more information see
[the official documentation](https://cloud.google.com/monitoring/support/notification-options)
and [API](https://cloud.google.com/monitoring/api/ref_v3/rest/v3/projects.notificationChannels/list).

## Example Usage

```hcl
data "google_monitoring_notification_channels" "all" {
  project = "my-project"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The ID of the project whose notification channels are listed. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/notificationChannels`

* `notification_channels` - The list of notification channels in the project. Each channel has the attributes of the
    [google_monitoring_notification_channel](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/monitoring_notification_channel)
    resource, including its `id`, except `project`, `deletion_policy`, `force_delete` and `sensitive_labels`.
    Labels with sensitive data, such as `auth_token`, `password` and `service_key`, are obfuscated or left out by the API.
    All other labels, such as the `url` of a webhook, are stored in the Terraform state as the API returns them.
