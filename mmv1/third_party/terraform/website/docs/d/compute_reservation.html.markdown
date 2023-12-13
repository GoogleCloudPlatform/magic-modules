---
subcategory: "Compute Engine"
description: |-
  Provide access to a Reservation's attributes
---

# google\_compute\_reservation

Provides access to available Google Compute Reservation Resources for a given project.
See more about [Reservations of Compute Engine resources](https://cloud.google.com/compute/docs/instances/reservations-overview) in the upstream docs.

```hcl
data "google_compute_reservation" "reservation" {
  name = "gce-reservation"
  zone = "us-central1-a"
}

resource "google_compute_reservation" "test" {
  name = "gce-reservation"
  zone = "us-central1-a"
  project = "PROJECT_ID"


  specific_reservation {
    count = 3
    instance_properties {
      min_cpu_platform = "Intel Cascade Lake"
      machine_type     = "n2-standard-2"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` (Required) - The name of the Compute Reservation.
* `zone` (Required) - Zone where the Compute Reservation resides.
* `project` (Optional) - Project from which to list the Compute Reservation. Defaults to project declared in the provider.

## Attributes Reference

See [google_compute_reservation](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_reservation) resource for details of the available attributes.
