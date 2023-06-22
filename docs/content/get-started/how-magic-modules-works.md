---
title: "How Magic Modules works"
weight: 20
aliases:
  - /docs/how-to/types-of-resources
  - /how-to/types-of-resources
---

# How Magic Modules works

Magic Modules can be thought of as a source of truth for how to map a GCP API resource representation to a Terraform resource representation. Magic Modules uses that mapping (and additional handwritten code where necessary) to generate "downstream" repositories - in particular, the Terraform providers for Google Platform: [`google`](https://github.com/hashicorp/terraform-provider-google) (or TPG) and [`google-beta`](https://github.com/hashicorp/terraform-provider-google-beta) (or TPGB).

Generation of the downstream repositories happens for every new commit in a PR (to a temporary branch owned by the [`modular-magician`](https://github.com/modular-magician/) robot user) and on every merge into the main branch (to the main branch of downstreams). Generation for PR commits allows contributors to manually examine the changes, as well as allowing automatic running of unit tests, acceptance tests, and automated checks such as breaking change detection.

## Resource types

There are three types of resources supported by Magic Modules: MMv1, Handwritten, and DCL/tpgtools. These are described in more detail in the following sections.

### MMv1 (preferred)

MMv1 is a Ruby-based code generator that implements Terraform Provider Google (TPG) resources from YAML specification files.

MMv1-generated resources like [google_compute_address](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address) can be identified by looking in their [Go source](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_compute_address.go) for an `AUTO GENERATED CODE` header as well as a Type `MMv1`. MMv1-generated resources should have source code present under their product folders, like [mmv1/products/compute](./products/compute) for the [google_compute_address](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address) resource.

### Handwritten

Handwritten resources like [google_container_cluster](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/container_cluster) can be identified if they have source code present under the [mmv1/third_party/terraform/resources](./resources) folder or by the absence of the `AUTO GENERATED CODE` header in their [Go source](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_container_cluster.go) in the downstream repositories. Handwritten datasources should be under the [mmv1/third_party/terraform/data_sources](./data_sources) folder, tests under the [mmv1/third_party/terraform/tests](./tests) folder and web documentation under the [mmv1/third_party/terraform/website](./website) folder.

### DCL aka tpgtools (maintenance mode)

DCL is a Go-based code generator that implements Terraform Provider Google (TPG) resources from YAML specification files.

DCL-generated resources like [google_bigquery_reservation_assignment](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/bigquery_reservation_assignment) can be identified by looking in their [Go source](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_bigquery_reservation_assignment.go) for an `AUTO GENERATED CODE` header as well as a Type `DCL`.

DCL is in maintenance mode, which means that new resources using the DCL are not being added.
