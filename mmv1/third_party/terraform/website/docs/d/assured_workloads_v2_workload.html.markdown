---
subcategory: "Assured Workloads V2"
description: |-
  Get information about an Assured Workloads V2 Workload.
---

# google_assured_workloads_v2_workload

Get information about an existing Assured Workloads V2 Workload. For more information see the
[official documentation](https://cloud.google.com/assured-workloads/docs/overview) and
[API](https://cloud.google.com/assured-workloads/docs/reference/rest/v2/organizations.locations.workloads).

## Example Usage

```hcl
data "google_assured_workloads_v2_workload" "my_workload" {
  organization = "123456789"
  location     = "us-central1"
  workload_id  = "my-workload-id"
}
```

Or using the full canonical resource name:

```hcl
data "google_assured_workloads_v2_workload" "my_workload" {
  name = "organizations/123456789/locations/us-central1/workloads/my-workload-id"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The full resource name of the workload (e.g. `organizations/123456789/locations/us-central1/workloads/my-workload-id`). Either `name` or all of `organization`, `location`, and `workload_id` must be specified.

* `organization` - (Optional) The parent organization ID of the Assured Workloads environment.

* `location` - (Optional) The location of the Assured Workloads environment.

* `workload_id` - (Optional) User-assigned or server-assigned unique identifier for the workload.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `description` - User-friendly description for the workload.

* `resource_config` - Target GCP Folder or Project where controls are deployed.

* `computed_target_resource` - Output only. The resolved target resource (e.g. `projects/123` or `folders/456`).

* `target_resource_display_name` - Output only. User-facing display name of the target resource.

* `framework` - The compliance framework active on this workload.

* `cloud_control_configs` - Control configurations enforced on this workload.

* `cmek_config` - Customer Managed Encryption Keys (CMEK) configuration.

* `state` - Output only. Status of the workload (e.g. `ACTIVE`, `CREATING`, `FAILED`).

* `create_time` - Output only. Creation timestamp in RFC3339 format.

* `update_time` - Output only. Last update timestamp in RFC3339 format.

* `etag` - Concurrency control token for optimistic locking.
