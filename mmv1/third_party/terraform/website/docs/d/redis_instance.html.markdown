---
subcategory: "Memorystore (Redis)"
description: |-
  Get information about a Google Cloud Redis instance.
---

# google_redis_instance

Get info about a Google Cloud Redis instance.
~> It may take a while for the attached tag bindings to be deleted after the instance is scheduled to be deleted.

## Example Usage

```tf
data "google_redis_instance" "my_instance" {
  name = "my-redis-instance"
}

output "instance_memory_size_gb" {
  value = data.google_redis_instance.my_instance.memory_size_gb
}

output "instance_connect_mode" {
  value = data.google_redis_instance.my_instance.connect_mode
}

output "instance_authorized_network" {
  value = data.google_redis_instance.my_instance.authorized_network
}
```
To create an instance with a tag

```tf
resource "fgoogle_redis_instance" "my_instance" {
  name       = "My instance"
  tags = {"tagKeys/281478409127147" : "tagValues/281479442205542"}
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of a Redis instance.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `region` - (Optional) The region in which the resource belongs. If it
    is not provided, the provider region is used.
    
* `tags` - (Optional) A map of resource manager tags. Resource manager tag keys and values have the same definition as resource manager tags. Keys must be in the format tagKeys/{tag_key_id}, and values are in the format tagValues/456. The field is ignored when empty. The field is immutable and causes resource replacement when mutated.

## Attributes Reference

See [google_redis_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/redis_instance) resource for details of the available attributes.
