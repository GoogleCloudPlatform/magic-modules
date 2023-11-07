---
subcategory: "Compute Engine"
description: |-
  Get information about a Google Compute Regional Persistent disks.
---

# google\_compute\_region\_disk

Get information about a Google Compute Regional Persistent disks.

[the official documentation](https://cloud.google.com/compute/docs/disks) and its [API](https://cloud.google.com/compute/docs/reference/rest/v1/regionDisks).

## Example Usage

```hcl
data "google_compute_region_disk" "disk" {
  name    = "persistent-regional-disk"
  project = "example"
  region  = "us-central1"
}

resource "google_compute_instance" "default" {
  # ...
    
  attached_disk {
    source = data.google_compute_disk.disk.self_link
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of a specific disk.

- - -

* `region` - (Optional) A reference to the region where the disk resides.

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/zones/{{zone}}/disks/{{name}}`

* `async_primary_disk` - A nested object resource Structure is [documented below](#nested_async_primary_disk).

* `replica_zones` - URLs of the zones where the disk should be replicated to

* `label_fingerprint` -
  The fingerprint used for optimistic locking of this resource.  Used
  internally during updates.

* `creation_timestamp` -
  Creation timestamp in RFC3339 text format.

* `last_attach_timestamp` -
  Last attach timestamp in RFC3339 text format.

* `last_detach_timestamp` -
  Last detach timestamp in RFC3339 text format.

* `users` -
  Links to the users of the disk (attached instances) in form:
  project/zones/zone/instances/instance

* `source_snapshot_id` -
  The unique ID of the snapshot used to create this disk. This value
  identifies the exact snapshot that was used to create this persistent
  disk. For example, if you created the persistent disk from a snapshot
  that was later deleted and recreated under the same name, the source
  snapshot ID would identify the exact version of the snapshot that was
  used.

* `description` -
  The optional description of this resource.

* `labels` - All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.

* `size` -
  Size of the persistent disk, specified in GB.

* `physical_block_size_bytes` -
  Physical block size of the persistent disk, in bytes.

* `type` -
  URL of the disk type resource describing which disk type to use to
  create the disk.

* `region` -
  A reference to the region where the disk resides.

* `snapshot` -
  The source snapshot used to create this disk.

* `source_snapshot_encryption_key` -
  (Optional)
  The customer-supplied encryption key of the source snapshot.

* `self_link` - The URI of the created resource.

<a name="nested_async_primary_disk"></a>The `async_primary_disk` block supports:

* `disk` - Primary disk for asynchronous disk replication.