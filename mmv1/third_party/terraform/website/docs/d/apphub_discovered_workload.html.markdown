---
subcategory: "Apphub"
description: |-
  Get information about a discovered workload.
---

# google\_apphub\_discovered_workload

Get information about a discovered workload from its uri.


## Example Usage


```hcl
data "google_apphub_discovered_workload" "my-workload" {
  location = "my-location"
  workload_uri = "my-workload-uri"
}
```

## Argument Reference

The following arguments are supported:

* `project` - The host project of the discovered workload.
* `workload_uri` - (Required) The uri of the workload.
* `location` - (Required) The location of the discovered workload.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `discovered_workload` - Represents a network/api interface that exposes some functionality to clients for consumption over the network. Structure is [documented below](#nested_discovered_workloads)

<a name="nested_discovered_workloads"></a>A `discovered_workload` object would contain the following fields:-

* `name` - Resource name of a Workload. Format: "projects/{host-project-id}/locations/{location}/applications/{application-id}/workloads/{workload-id}".

* `workload_reference` - Reference to an underlying networking resource that can comprise a Workload. Structure is [documented below](#nested_workload_reference)

<a name="nested_workload_reference"></a>A `workload_reference` object would contain the following fields:-

* uri - The underlying resource URI.

* path - Additional path under the resource URI.

* `workload_properties` - Properties of an underlying compute resource that can comprise a Workload. Structure is [documented below](#nested_workload_properties)

<a name="nested_workload_properties"></a>A `workload_properties` object would contain the following fields:-

* gcp_project - The service project identifier that the underlying cloud resource resides in.

* location - The location that the underlying resource resides in.

* zone - The location that the underlying resource resides in if it is zonal.
