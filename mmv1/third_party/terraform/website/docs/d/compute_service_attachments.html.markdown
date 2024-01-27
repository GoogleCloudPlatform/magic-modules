---
subcategory: "Compute Engine"
description: |-
  Get a service attachment in a Google Cloud project.
---

# google\_compute\_service_attachment

Get a service attachment  in a specified Google Cloud project

To get more information about service attachment, see:

* [API documentation](https://cloud.google.com/compute/docs/reference/rest/v1/serviceAttachments/)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/vpc/docs/configure-private-service-connect-producer)
## Example Usage

```tf
data "google_compute_networks" "my-networks" {
  project = "my-cloud-project"
}
data "google_compute_service_attachment" "my_attachment" {
  name = "psc-service-attachment-target"
  project = "my-cloud-project
  region = "us-west1"
}

```

## Argument Reference

The following arguments are supported:

* `project` - (required) The name of the project.
* `name` - (required) The name of the service attachment.
* `region` - (required) The region where the service attachment is created.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `name` -The name of the service attachment resource

* `region` - The region where the resource exists.

* `project` - The project name being queried.

* `self_link` - The URI of the resource.