---
subcategory: "Assured Workloads V2"
description: |-
  Get information about Assured Workloads V2 organization settings.
---

# google_assured_workloads_v2_settings

Get information about Assured Workloads V2 organization settings, such as Organization-Level Personnel Controls (OLPC). For more information see the
[official documentation](https://cloud.google.com/assured-workloads/docs/settings) and
[API](https://cloud.google.com/assured-workloads/docs/reference/rest/v2/organizations/assuredWorkloadsSettings).

## Example Usage

```hcl
data "google_assured_workloads_v2_settings" "my_settings" {
  organization = "123456789"
}
```

## Argument Reference

The following arguments are supported:

* `organization` - (Required) The parent organization ID.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - The identifier of the settings singleton resource (`organizations/{organization}/assuredWorkloadsSettings`).

* `name` - The resource name of the settings singleton.

* `olpc_mode` - Organization-Level Personnel Control (OLPC) mode.

* `etag` - Concurrency token for optimistic locking.
