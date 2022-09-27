---
subcategory: "Compute Engine"
page_title: "Google: google_compute_region_network_endpoint_group"
description: |-
  Get access to a Network Endpoint Group's Attributes.
---

# google\_compute\_region\_network\_endpoint\_group

Use this data source to access a Regional Network Endpoint Group's attributes.

The NEG may be found by providing either a `self_link`, or a `name` and a `region`.

## Example Usage

```hcl
// Cloud Run Example
resource "google_compute_region_network_endpoint_group" "cloudrun_neg" {
  name                  = "cloudrun-neg"
  network_endpoint_type = "SERVERLESS"
  region                = "us-central1"
  cloud_run {
    service = google_cloud_run_service.cloudrun_neg.name
  }
}

resource "google_cloud_run_service" "cloudrun_neg" {
  name     = "cloudrun-neg"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/hello"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

data "google_compute_region_network_endpoint_group" "data_source" {
  self_link = google_compute_region_network_endpoint_group.cloudrun_neg.self_link
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the instance group. One of `name` or `self_link` must be provided.
* `self_link` - (Optional) The link to the instance group. One of `name` or `self_link` must be provided.
* `project` - (Optional) The ID of the project in which the resource belongs. If `self_link` is provided, this value is ignored. If neither `self_link` nor `project` are provided, the provider project is used.
* `region` - (Optional) The region in which the resource belongs. If `self_link` is provided, this value is ignored. If neither `self_link` nor `region` are provided, the provider region is used.

## Attributes Reference

In addition the arguments listed above, the following attributes are exported:

* `network` - The network to which all network endpoints in the NEG belong.
* `subnetwork` - subnetwork to which all network endpoints in the NEG belong.
* `description` - The NEG description.
* `network_endpoint_type` - Type of network endpoints in this network endpoint group.
* `default_port` - The NEG default port.
* `size` - Number of network endpoints in the network endpoint group.
