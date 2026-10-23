---
subcategory: "Compute Engine"
description: |-
  A data source to retrieve a service attachment.
---

# `google_compute_service_attachment`

Get a specific [service attachment](https://cloud.google.com/vpc/docs/configure-private-service-connect-services) within a region. For more information see the
[official documentation](https://cloud.google.com/vpc/docs/configure-private-service-connect-services)
and [API](https://cloud.google.com/compute/docs/reference/rest/v1/serviceAttachments/get).

## Example Usage

```hcl
data "google_compute_service_attachment" "default" {
  project = "my-project"
  name    = "my-service-attachment"
  region  = "us-west2"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the service attachment to retrieve.

---

* `region` - (Optional) The region in which the service attachment resides.
  If it is not provided, the provider region is used.

* `project` - (Optional) The ID of the project in which the resource belongs.
  If it is not provided, the provider project is used.

## Attributes Reference

See [google_compute_service_attachment](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_service_attachment#attributes-reference) resource for details of the available attributes.
