---
subcategory: "Compute Engine"
description: |-
  Get information about a Google Compute Engine Reservation Sub-Block.
---

# google_compute_reservation_sub_block

Get information about a Google Compute Engine Reservation Sub-Block. Reservation sub-blocks are automatically created by Google Cloud within reservation blocks and represent a finer-grained physical grouping of resources.

For more information see the [official documentation](https://cloud.google.com/compute/docs/instances/reserving-zonal-resources)
and the [API](https://cloud.google.com/compute/docs/reference/rest/v1/reservationSubBlocks).

## Example Usage

```hcl
data "google_compute_reservation_sub_block" "sub_block" {
  name              = "my-reservation-sub-block"
  reservation_block = "my-reservation-block"
  reservation       = "my-reservation"
  zone              = "us-central1-a"
}

output "sub_block_status" {
  value = data.google_compute_reservation_sub_block.sub_block.status
}

output "sub_block_health" {
  value = data.google_compute_reservation_sub_block.sub_block.health_info
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the reservation sub-block.

* `reservation_block` - (Required) The name of the parent reservation block.

* `reservation` - (Required) The name of the parent reservation.

* `zone` - (Required) The zone where the reservation sub-block resides.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - The identifier for the resource with format `projects/{{project}}/zones/{{zone}}/reservations/{{reservation}}/reservationBlocks/{{reservation_block}}/reservationSubBlocks/{{name}}`

* `kind` - Type of the resource. Always `compute#reservationSubBlock` for reservation sub-blocks.

* `resource_id` - The unique identifier for the resource.

* `creation_timestamp` - Creation timestamp in RFC3339 text format.

* `self_link` - Server-defined fully-qualified URL for this resource.

* `self_link_with_id` - Server-defined URL for this resource with the resource id.

* `count` - The number of hosts that are allocated in this reservation sub-block.

* `in_use_count` - The number of instances that are currently in use on this reservation sub-block.

* `status` - Status of the reservation sub-block.

* `reservation_sub_block_maintenance` - Maintenance information for this reservation sub-block. Structure is [documented below](#nested_reservation_sub_block_maintenance).

* `physical_topology` - The physical topology of the reservation sub-block. Structure is [documented below](#nested_physical_topology).

* `health_info` - Health information for the reservation sub-block. Structure is [documented below](#nested_health_info).

<a name="nested_reservation_sub_block_maintenance"></a>The `reservation_sub_block_maintenance` block contains:

* `maintenance_ongoing_count` - Number of hosts in the sub-block that have ongoing maintenance.

* `maintenance_pending_count` - Number of hosts in the sub-block that have pending maintenance.

* `scheduling_type` - The type of maintenance for the reservation.

* `subblock_infra_maintenance_ongoing_count` - Number of sub-block infrastructure that has ongoing maintenance.

* `subblock_infra_maintenance_pending_count` - Number of sub-block infrastructure that has pending maintenance.

* `instance_maintenance_ongoing_count` - Number of instances that have ongoing maintenance.

* `instance_maintenance_pending_count` - Number of instances that have pending maintenance.

<a name="nested_physical_topology"></a>The `physical_topology` block contains:

* `cluster` - The cluster name of the reservation sub-block.

* `block` - The hash of the capacity block within the cluster.

* `sub_block` - The hash of the capacity sub-block within the capacity block.

<a name="nested_health_info"></a>The `health_info` block contains:

* `health_status` - The health status of the reservation sub-block.

* `healthy_host_count` - The number of healthy hosts in the reservation sub-block.

* `degraded_host_count` - The number of degraded hosts in the reservation sub-block.

* `healthy_infra_count` - The number of healthy infrastructure (e.g. NVLink domain) in the reservation sub-block.

* `degraded_infra_count` - The number of degraded infrastructure (e.g. NVLink domain) in the reservation sub-block.
