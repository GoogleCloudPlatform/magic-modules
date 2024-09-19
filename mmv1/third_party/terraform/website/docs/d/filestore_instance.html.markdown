---
subcategory: "Filestore"
description: |-
  Get information about a Google Cloud Filestore instance.
---

# google_filestore_instance

Get info about a Google Cloud Filestore instance.

~> It may take a while for the attached tag bindings to be deleted after the instance is scheduled to be deleted.

## Example Usage

```tf
data "google_filestore_instance" "my_instance" {
  name = "my-filestore-instance"
}

output "instance_ip_addresses" {
  value = data.google_filestore_instance.my_instance.networks.ip_addresses
}

output "instance_connect_mode" {
  value = data.google_filestore_instance.my_instance.networks.connect_mode
}

output "instance_file_share_name" {
  value = data.google_filestore_instance.my_instance.file_shares.name
}
```
To create an instance with a tag

```tf
resource "filestore_instance" "my_instance" {
  name       = "My instance"
  instance_id = "your-instance-id"
  tags = {"tagKeys/281478409127147" : "tagValues/281479442205542"}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of a Filestore instance.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `location` - (Optional) The name of the location of the instance. This 
    can be a region for ENTERPRISE tier instances. If it is not provided, 
    the provider region or zone is used.
    
* `tags` - (Optional) A map of resource manager tags. Resource manager tag keys and values have the same definition as resource manager tags. Keys must be in the format tagKeys/{tag_key_id}, and values are in the format tagValues/456. The field is ignored when empty. The field is immutable and causes resource replacement when mutated.

## Attributes Reference

See [google_filestore_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/filestore_instance) resource for details of the available attributes.
