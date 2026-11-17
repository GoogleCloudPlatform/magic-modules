---
subcategory: "Compute Engine"
description: |-
  Get SBOM software component versions for a Google Compute Engine Host.
---

# google_compute_host_version

Get SBOM software component versions for a Google Compute Engine Host (`hosts.getVersion`), optionally scoped to a reservation, reservation block, or reservation sub-block.

~> **Note:** Each read or refresh of this data source dispatches an asynchronous `POST .../getVersion` zonal operation against an active SBOM-supported host and waits for the operation to complete before reading `inlineSbomInfo` from the operation metadata.

For more information see the [official documentation](https://cloud.google.com/compute/docs/instances/reserving-zonal-resources)
and the [API](https://cloud.google.com/compute/docs/reference/rest/v1/hosts/getVersion).

## Example Usage

```hcl
data "google_compute_host_version" "host_version" {
  name            = "my-host"
  reservation     = "my-reservation"
  zone            = "us-west1-a"
  sbom_selections = ["SBOM_SELECTION_CURRENT"]
}

output "current_versions" {
  value = data.google_compute_host_version.host_version.current_component_versions
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name or resource ID of the host.

* `sbom_selections` - (Required) The SBOM selection to return. Possible values include `SBOM_SELECTION_CURRENT` and `SBOM_SELECTION_TARGET`. Exactly one value must be specified.

- - -

* `reservation` - (Optional) The name or self-link of the parent reservation. When omitted, calls `getVersion` on the top-level zonal host path.

* `reservation_block` - (Optional) The name of the parent reservation block. Requires `reservation` to be set.

* `reservation_sub_block` - (Optional) The name of the parent reservation sub-block. Requires `reservation` and `reservation_block` to be set.

* `zone` - (Optional) The zone where the host resides. If it is not provided, any zone in a `reservation` or `name` self-link is used before falling back to the provider zone.

* `project` - (Optional) The project in which the resource belongs. If it is not provided, any project in a `reservation` or `name` self-link is used before falling back to the provider project.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - The identifier for the resource with format `projects/{{project}}/zones/{{zone}}/hosts/{{name}}/getVersion/{{sbom_selection}}` or `projects/{{project}}/zones/{{zone}}/reservations/{{reservation}}[/reservationBlocks/{{reservation_block}}[/reservationSubBlocks/{{reservation_sub_block}}]]/hosts/{{name}}/getVersion/{{sbom_selection}}`.

* `current_component_versions` - A map of component names to their currently-applied versions or identifiers (populated when `SBOM_SELECTION_CURRENT` is requested).

* `target_component_versions` - A map of component names to their target versions or identifiers (populated when `SBOM_SELECTION_TARGET` is requested).

## Timeouts

This data source provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `read` - Default is 20 minutes.
