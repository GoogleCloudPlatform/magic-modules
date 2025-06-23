---
subcategory: "Parallelstore"
description: |-
  Fetches the details of available instance.
---


# google_parallelstore_instance


Use this data source to get information about the available instance. For more details refer the [API docs](https://cloud.google.com/parallelstore/docs/reference/rest/v1/projects.locations.instances).


## Example Usage


```hcl
data "google_parallelstore_instance" "default" {
  name     = "instance-name"
  location = "us-central1-a"
}
```


## Argument Reference


The following arguments are supported:


* `name` -
  (Required)
  The ID of the parallelstore instance.
  'parallelstore_instance_id'


* `project` - 
  (optional) 
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.


* `location` -
  (optional)
  The canonical id of the location. If it is not provided, the provider project is used. For example: us-central1-a.


## Attributes Reference


See [google_parallelstore_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/parallelstore_instance) resource for details of all the available attributes.