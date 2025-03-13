---
subcategory: "Compute Engine"
description: |-
  Get information about a Regional Backend Service.
---

# google_compute_region_backend_service

Get information about a Regional Backend Service. For more information see
[the official documentation](https://cloud.google.com/compute/docs/load-balancing/internal/backend-service) and
[API](https://cloud.google.com/compute/docs/reference/rest/beta/regionBackendServices).

## Example Usage

```hcl
data "google_compute_region_backend_service" "my_backend" {
  name   = "my-backend-service"
  region = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the regional backend service.

* `region` - (Required) The region where the backend service resides.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - An identifier for the resource with format `projects/{{project}}/regions/{{region}}/backendServices/{{name}}`

* `self_link` - The URI of the created resource.

* `description` - An optional description of this resource.

* `protocol` - The protocol this backend service uses to communicate with backends.

* `session_affinity` - The session affinity setting for this backend service.

* `timeout_sec` - How many seconds to wait for the backend before considering it a failed request.

* `load_balancing_scheme` - The load balancing scheme. Possible values are: `INTERNAL`, `EXTERNAL`, `INTERNAL_MANAGED`, or `INTERNAL_SELF_MANAGED`.

* `health_checks` - The list of URLs to the health checks associated with this backend service.

* `backend` - The list of backends that serve this backend service.
