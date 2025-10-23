---
subcategory: "Compute Engine"
description: |-
  Get information about a Google Compute Engine Reservation Block.
---

# google_compute_reservation_block

Get information about a Google Compute Engine Reservation Block. Reservation blocks are automatically created by Google Cloud within reservations and represent a physical grouping of resources.

For more information see the [official documentation](https://cloud.google.com/compute/docs/instances/reserving-zonal-resources)
and the [API](https://cloud.google.com/compute/docs/reference/rest/v1/reservationBlocks).

## Example Usage

```hcl
data "google_compute_reservation_block" "block" {
  name        = "my-reservation-block"
  reservation = "my-reservation"
  zone        = "us-central1-a"
}

output "block_status" {
  value = data.google_compute_reservation_block.block.status
}

output "block_in_use_count" {
  value = data.google_compute_reservation_block.block.in_use_count
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the reservation block.

* `reservation` - (Required) The name of the parent reservation.

* `zone` - (Required) The zone where the reservation block resides.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - The identifier for the resource with format `projects/{{project}}/zones/{{zone}}/reservations/{{reservation}}/reservationBlocks/{{name}}`

* `kind` - Type of the resource. Always `compute#reservationBlock` for reservation blocks.

* `resource_id` - The unique identifier for the resource.

* `creation_timestamp` - Creation timestamp in RFC3339 text format.

* `self_link` - Server-defined fully-qualified URL for this resource.

* `self_link_with_id` - Server-defined URL for this resource with the resource id.

* `count` - The number of resources that are allocated in this reservation block.

* `in_use_count` - The number of instances that are currently in use on this reservation block.

* `status` - Status of the reservation block.

* `reservation_sub_block_count` - The number of reservation sub-blocks associated with this reservation block.

* `reservation_sub_block_in_use_count` - The number of in-use reservation sub-blocks associated with this reservation block.

* `reservation_maintenance` - Maintenance information for this reservation block. Structure is [documented below](#nested_reservation_maintenance).

* `physical_topology` - The physical topology of the reservation block. Structure is [documented below](#nested_physical_topology).

* `health_info` - Health information for the reservation block. Structure is [documented below](#nested_health_info).

<a name="nested_reservation_maintenance"></a>The `reservation_maintenance` block contains:

* `maintenance_ongoing_count` - Number of hosts in the block that have ongoing maintenance.

* `maintenance_pending_count` - Number of hosts in the block that have pending maintenance.

* `scheduling_type` - The type of maintenance for the reservation.

* `subblock_infra_maintenance_ongoing_count` - Number of sub-block infrastructure that has ongoing maintenance.

* `subblock_infra_maintenance_pending_count` - Number of sub-block infrastructure that has pending maintenance.

* `instance_maintenance_ongoing_count` - Number of instances that have ongoing maintenance.

* `instance_maintenance_pending_count` - Number of instances that have pending maintenance.

<a name="nested_physical_topology"></a>The `physical_topology` block contains:

* `cluster` - The cluster name of the reservation block.

* `block` - The hash of the capacity block within the cluster.

<a name="nested_health_info"></a>The `health_info` block contains:

* `health_status` - The health status of the reservation block.

* `healthy_sub_block_count` - The number of sub-blocks that are healthy.

* `degraded_sub_block_count` - The number of sub-blocks that are degraded.
