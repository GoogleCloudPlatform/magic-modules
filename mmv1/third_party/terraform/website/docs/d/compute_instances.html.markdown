---
subcategory: "Compute Engine"
description: |-
  Get a list of instances within GCE.
---

# google_compute_instances

Get a list of instances within GCE, either in a specific zone or across every
zone in a project. This is useful for inventory and automation use cases,
such as identifying all instances that need to be backed up.

For information about instances and how to use them, see
[the official documentation](https://cloud.google.com/compute/docs/instances)
and [API](https://cloud.google.com/compute/docs/reference/latest/instances).

## Example Usage

```hcl
# All instances in a single zone.
data "google_compute_instances" "zone" {
  zone = "us-central1-a"
}

# Every instance in the project, across all zones.
data "google_compute_instances" "project" {
  filter = "status = RUNNING"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The ID of the project in which the resources belong. If it
    is not provided, the provider project is used.

* `zone` - (Optional) The zone to list instances in, for example `us-central1-a`. If it
    is not provided, instances from every zone in the project are returned. Note that,
    unlike most other resources, this does not fall back to a `zone` configured on the
    provider: omitting it here always means "every zone".

* `filter` - (Optional) A filter expression as described in the
    [REST API](https://cloud.google.com/compute/docs/reference/rest/v1/instances/list#query-parameters),
    used to restrict the instances returned.

## Attributes Reference

The following attributes are exported:

* `id` - Identifier

* `instances` - A list of instances matching the given filter, zone and project. Structure is [defined below](#nested_instances).

<a name="nested_instances"></a>The `instances` block supports the fields also found on the
[google_compute_instance](/docs/providers/google/d/compute_instance.html) datasource, with
the following differences:

* Fields that require an additional API call per instance to populate -- such as
  `boot_disk.0.initialize_params` -- are omitted, so that this datasource stays fast and
  usable against projects with many instances. Use the `google_compute_instance` datasource
  to get full detail on a single, known instance.

* `current_status` reflects the instance's status (e.g. `RUNNING`, `TERMINATED`) at read time.
