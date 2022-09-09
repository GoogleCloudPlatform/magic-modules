---
subcategory: "Compute Engine"
page_title: "Google: google_compute_snapshot"
description: |-
  Get information about a Google Compute Snapshot.
---

# google\_compute\_snapshot

Get information about a Google Compute Persistent disks.

[the official documentation](https://cloud.google.com/compute/docs/disks/create-snapshots) and its [API](https://cloud.google.com/compute/docs/reference/rest/v1/snapshots).

## Example Usage

```hcl
#by name 
data "google_compute_snapshot" "snapshot" {
  name    = "generic-tpl-20200107"
}

# using a filter
data "google_compute_snapshot" "latest-snapshot" {
  filter      = "name != generic-tpl-20200107"
  most_recent = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the instance template. One of `name` or `filter` must be provided.

- `filter` - (Optional) A filter to retrieve the instance templates.
    See [gcloud topic filters](https://cloud.google.com/sdk/gcloud/reference/topic/filters) for reference.
    If multiple instance templates match, either adjust the filter or specify `most_recent`. One of `name` or `filter` must be provided.

- `most_recent` - (Optional) If `filter` is provided, ensures the most recent template is returned when multiple instance templates match. One of `name` or `filter` must be provided.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `description` -
  
  An optional description of this resource.

* `storage_locations` -
  
  Cloud Storage bucket storage location of the snapshot (regional or multi-regional).

* `labels` -
  
  Labels to apply to this Snapshot.

* `zone` -
  
  A reference to the zone where the disk is hosted.

* `snapshot_encryption_key` -
  
  The customer-supplied encryption key of the snapshot. Required if the
  source snapshot is protected by a customer-supplied encryption key.
  Structure is [documented below](#nested_snapshot_encryption_key).

* `source_disk_encryption_key` -
  
  The customer-supplied encryption key of the source snapshot. Required
  if the source snapshot is protected by a customer-supplied encryption
  key.
  Structure is [documented below](#nested_source_disk_encryption_key).

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.


<a name="nested_snapshot_encryption_key"></a>The `snapshot_encryption_key` block supports:

* `raw_key` -
  
  Specifies a 256-bit customer-supplied encryption key, encoded in
  RFC 4648 base64 to either encrypt or decrypt this resource.
  **Note**: This property is sensitive and will not be displayed in the plan.

* `sha256` -
  The RFC 4648 base64 encoded SHA-256 hash of the customer-supplied
  encryption key that protects this resource.

* `kms_key_self_link` -
  
  The name of the encryption key that is stored in Google Cloud KMS.

* `kms_key_service_account` -
  
  The service account used for the encryption request for the given KMS key.
  If absent, the Compute Engine Service Agent service account is used.

<a name="nested_source_disk_encryption_key"></a>The `source_disk_encryption_key` block supports:

* `raw_key` -
  
  Specifies a 256-bit customer-supplied encryption key, encoded in
  RFC 4648 base64 to either encrypt or decrypt this resource.
  **Note**: This property is sensitive and will not be displayed in the plan.

* `kms_key_service_account` -
  
  The service account used for the encryption request for the given KMS key.
  If absent, the Compute Engine Service Agent service account is used.
---

* `id` - an identifier for the resource with format `projects/{{project}}/global/snapshots/{{name}}`

* `creation_timestamp` -
  Creation timestamp in RFC3339 text format.

* `snapshot_id` -
  The unique identifier for the resource.

* `disk_size_gb` -
  Size of the snapshot, specified in GB.

* `storage_bytes` -
  A size of the storage used by the snapshot. As snapshots share
  storage, this number is expected to change with snapshot
  creation/deletion.

* `licenses` -
  A list of public visible licenses that apply to this snapshot. This
  can be because the original image had licenses attached (such as a
  Windows image).  snapshotEncryptionKey nested object Encrypts the
  snapshot using a customer-supplied encryption key.

* `label_fingerprint` -
  The fingerprint used for optimistic locking of this resource. Used
  internally during updates.
* `self_link` - The URI of the created resource.
