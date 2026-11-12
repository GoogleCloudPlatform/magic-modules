---
subcategory: "Compute Engine"
description: |-
  Get information about a Google Compute Engine Host.
---

# google_compute_reservation_host

Get information about a Google Compute Engine Host, optionally scoped to a reservation, reservation block, or reservation sub-block.

For more information see the [official documentation](https://cloud.google.com/compute/docs/instances/reserving-zonal-resources)
and the [API](https://cloud.google.com/compute/docs/reference/rest/v1/hosts/get).

## Example Usage

```hcl
data "google_compute_reservation_host" "host" {
  name                  = "my-host"
  reservation           = "my-reservation"
  reservation_block     = "my-reservation-block"
  reservation_sub_block = "my-reservation-sub-block"
  zone                  = "us-central1-a"
}

output "host_state" {
  value = data.google_compute_reservation_host.host.state
}

output "host_status" {
  value = data.google_compute_reservation_host.host.status
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name or resource ID of the host.

- - -

* `reservation` - (Optional) The name of the parent reservation. When omitted, queries the top-level zonal host collection.

* `reservation_block` - (Optional) The name of the parent reservation block. Requires `reservation` to be set.

* `reservation_sub_block` - (Optional) The name of the parent reservation sub-block. Requires `reservation` and `reservation_block` to be set.

* `zone` - (Optional) The zone where the host resides. If it is not provided, the provider zone is used.

* `project` - (Optional) The project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - The identifier for the resource with format `projects/{{project}}/zones/{{zone}}/hosts/{{name}}` or `projects/{{project}}/zones/{{zone}}/{{association}}/hosts/{{name}}`.

* `kind` - Type of the resource. Always `compute#host` for hosts.

* `resource_id` - The unique server-generated identifier for the host.

* `creation_timestamp` - Creation timestamp in RFC3339 text format.

* `self_link` - Server-defined fully-qualified URL for this resource.

* `self_link_with_id` - Server-defined URL for this resource with the resource ID.

* `alias_links` - All alias resource paths for this host.

* `state` - The state of the host (for example, `CREATING`, `ACTIVE`, `DELETING`, or `UNAVAILABLE`).

* `status` - The status of the host. Structure is [documented below](#nested_status).

<a name="nested_status"></a>The `status` block contains:

* `physical_topology` - The physical topology of the host, if present. Structure is [documented below](#nested_physical_topology).

* `running_instances` - The URIs of the instances currently running on this host.

<a name="nested_physical_topology"></a>The `physical_topology` block contains:

* `cluster` - The cluster name of the reservation sub-block.

* `block` - The unique identifier of the capacity block within the cluster.

* `sub_block` - The unique identifier of the capacity sub-block within the capacity block.

* `host` - The unique identifier of the capacity host within the capacity sub-block.
