---
page_title: "Use Terraform Actions in the Google Cloud provider"
description: |-
  Learn how Terraform Action work and how to use them in the Google Cloud Terraform provider
---

# Terraform Actions in the Google Cloud provider

Terraform Actions introduce an **operational execution layer** to Terraform that complements the traditional **declarative resource model**.

Actions allow Terraform configurations to perform **imperative operations**-such as starting a virtual machine, triguring build, or publishing a message-**without creating, updating or destroing  infrastructure**.

Terraform Actions are available in Terraform v1.14 and later. For more information, see the [official HashiCorp documentation for Terraform Actions](https://developer.hashicorp.com/terraform/language/invoke-actions)

---

## What Are Terraform Actions?

Terraform Actions are **Provider-defined operations** that are executed during `terraform apply`, but **do not participate in Terraform state**.

Terraform Actions allow you to trigger imperative operations—such as power cycling a VM—using the formal Terraform framework. Unlike data sources or managed resources, they don’t affect infrastructure state but execute operations on existing resources.

An Action:
- Executes an API call or operational task
- Produces no managed infrastructure
- Does not store state
- Can be safely re-run

---

## Why Use Actions?

Many real-worl workflows require **imperative operations** that do not map cleanly to Terraform resources.

Examples include:
- Starting or stopping an existing  Conpute engine
- Restarting a database for maintenence
- Publishing a Pub/Sub message to trigger downstream autotmation
- Invalidating a CDN cache
- Triggering a Cloud Build pipeline

Before Terrafiorm Actions, users typically handled these task by:
- Running 'gcloud' commands manually
- Writing custom scripts
- Using CI/CD glue outside ZTerraform

Actions bring these workflows **into Terraform**, while preserving Terraforms's core principles.

---

## How Terraform Core Executes Actions

Terraform Actions are a **Terraform Core feature** (Terraform CLI v1.14+)

During `terraform apply`:
1. Terraform evaluates the configuration
2. Terraform builds a dependency graph
3. When an `action` block is reached:
  - Terraform calls the provider's Action implementation
  - The provider executes the operation
  - Progress and diagnostics are steamedback to Terraform

Action:
- Are not planned
- Are not diffed
- Do not affect the plan output
- Are executed **only during the apply**

---

## Action Syntax

Actions are declared using a top-level `action` block:

```hcl
action "<action_type>" "<name>"{
  config {
    # action- specific configuration
  }
}
```

Key properties:
- action_type is defined by the provider
- name is local to the configuration
- config contains provider-defined arguments
- depends_on can be used for ordering

---

## Terraform Actions currently supported in the Google Cloud provider

### Compute Instance Power Action
The `google_compute_instance_power` action enables power operations (`stop`, `start`, `restart`) on a Compute Engine instance. This feature uses Terraform Actions introduced in Terraform CLI v1.14+.


```hcl
resource "google_compute_instance" "example" {
  name         = "instance-power-example"
  machine_type = "e2-micro"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"
  }

  lifecycle {
    action_trigger {
      events  = ["after_create", "after_update"]
      actions = [action.google_compute_instance_power.power]
    }
  }
}

action "google_compute_instance_power" "power" {
  config {
    instance  = google_compute_instance.example.name
    project   = google_compute_instance.example.project
    zone      = google_compute_instance.example.zone
    operation = "restart"
  }
}
```

Arguments Reference:
The following arguments are supported in the config block:
- `instance`: (Required) The name of the Compute Engine instance.
- `project`: (Required) The GCP project containing the instance.
- `zone`: (Required) The zone where the instance resides (e.g., us-central1-a).
- `operation`: (Required) The power operation to perform. Valid values:
    - "stop": Stops the instance.
    - "start": Starts the instance.
    - "restart": Restarts the instance.