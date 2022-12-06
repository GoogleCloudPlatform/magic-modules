---
title: "Types of resources"
summary: "Check the header in the Go source to determine what type of resource it is. If there is no header, it is likely handwritten."
weight: 1
---

# Types of resources

## MMv1

MMv1 is a Ruby-based code generator that implements Terraform Provider Google (TPG) resources from YAML specification files.

MMv1-generated resources like [google_compute_address](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address) can be identified by looking in their [Go source](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_compute_address.go) for an `AUTO GENERATED CODE` header as well as a Type `MMv1`. MMv1-generated resources should have source code present under their product folders, like [mmv1/products/compute](./products/compute) for the [google_compute_address](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address) resource.

## Handwritten

Handwritten resources like [google_container_cluster](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/container_cluster) can be identified if they have source code present under the [mmv1/third_party/terraform/resources](./resources) folder or by the absence of the `AUTO GENERATED CODE` header in their [Go source](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_container_cluster.go) in the downstream repositories. Handwritten datasources should be under the [mmv1/third_party/terraform/data_sources](./data_sources) folder, tests under the [mmv1/third_party/terraform/tests](./tests) folder and web documentation under the [mmv1/third_party/terraform/website](./website) folder.

## DCL aka tpgtools (maintenance mode)

DCL is a Go-based code generator that implements Terraform Provider Google (TPG) resources from YAML specification files.

DCL-generated resources like [google_bigquery_reservation_assignment](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/bigquery_reservation_assignment) can be identified by looking in their [Go source](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_bigquery_reservation_assignment.go) for an `AUTO GENERATED CODE` header as well as a Type `DCL`.

DCL is in maintenance mode, which means that new resources using the DCL are not being added.
