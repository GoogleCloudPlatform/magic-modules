---
subcategory: "Compute Engine"
description: |-
  Get information about a Google Compute Engine Region Backend Service.
---

# google\_compute\_region\_backend\_service

Get information about a Google Compute Engine Region Backend Service. For more information see
[the official documentation](https://cloud.google.com/compute/docs/load-balancing/internal/regional-backend-service).

## Example Usage

```hcl
data "google_compute_region_backend_service" "default" {
  name   = "my-backend-service"
  region = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the resource.

* `region` - (Required) The region where the regional backend service resides.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs.
  If it is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `description` - Description of the resource.

* `backend` - The list of backends that serve this backend service. Structure is documented below.

* `fingerprint` - The fingerprint of this resource.

* `health_checks` - The list of URLs to the health checks used to verify traffic to backends.

* `load_balancing_scheme` - The load balancing scheme.

* `protocol` - The protocol used to communicate with backends.

* `session_affinity` - The session affinity setting.

* `timeout_sec` - The backend service timeout.

* `connection_draining_timeout_sec` - Time for which instance will be drained.

* `port_name` - Name of backend port.

* `self_link` - The URI of the created resource.

* `creation_timestamp` - Creation timestamp in RFC3339 text format.

The `backend` block supports:

* `group` - The fully-qualified URL of an instance group or network endpoint group.

* `balancing_mode` - The balancing mode for this backend.

* `capacity_scaler` - A multiplier applied to the group's maximum servicing capacity.

* `description` - An optional description.

* `failover` - If this backend is a failover backend.

* `max_connections` - Maximum number of connections for the group.

* `max_connections_per_instance` - Maximum number of connections per instance.

* `max_connections_per_endpoint` - Maximum number of connections per endpoint.

* `max_rate` - Maximum requests per second (RPS) for the group.

* `max_rate_per_instance` - Maximum RPS per instance.

* `max_rate_per_endpoint` - Maximum RPS per endpoint.

* `max_utilization` - The target CPU utilization for the group. 