---
title: "How Magic Modules works"
weight: 20
aliases:
  - /docs/how-to/types-of-resources
  - /how-to/types-of-resources
---

# How Magic Modules works

Magic Modules can be thought of as a source of truth for how to map a GCP API resource representation to a Terraform resource (or datasource) representation. Magic Modules uses that mapping (and additional handwritten code where necessary) to generate "downstream" repositories - in particular, the Terraform providers for Google Cloud: [`google`](https://github.com/hashicorp/terraform-provider-google) (or TPG) and [`google-beta`](https://github.com/hashicorp/terraform-provider-google-beta) (or TPGB).

Generation of the downstream repositories happens for every new commit in a PR (to a temporary branch owned by the [`modular-magician`](https://github.com/modular-magician/) robot user) and on every merge into the main branch (to the main branch of downstreams). Generation for PR commits allows contributors to manually examine the changes, as well as allowing automatic running of unit tests, acceptance tests, and automated checks such as breaking change detection.

## Resource types

There are three types of resources supported by Magic Modules: MMv1, Handwritten, and DCL/tpgtools. These are described in more detail in the following sections.

### MMv1

MMv1 consists of a set of "products"; each product contains one or more "resources".

Each product has a folder in [`magic-modules/mmv1/products`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products). The name of the folder is the "product name", which usually corresponds to the API subdomain covered by the product (such as `compute.googleapis.com`). Each product folder contains a product configuration file (`product.yaml`) and one or more resource configuration files (`ResourceName.yaml`). The actual name of a `ResourceName.yaml` file usually matches the name of a GCP API resource in the product's subdomain.

MMv1 resource configurations may reference handwritten code stored in [`magic-modules/mmv1/templates/terraform`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform), which will be injected into the generated resource file. Many MMv1 resources also have one or more handwritten tests, which are stored in the appropriate service folder inside [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services)

In the providers, MMv1-based resources are stored in `PROVIDER/services/PRODUCT/resource_PRODUCT_RESOURCE.go`, where `PROVIDER` is `google` or `google-beta`, `PRODUCT` is the product name, and RESOURCE is the GCP API resource's name converted to [snake case ↗](https://en.wikipedia.org/wiki/Snake_case).

MMv1-based files start with the following header:

```
***     AUTO GENERATED CODE    ***    Type: MMv1     ***
```

### Handwritten

Handwritten resources and datasources are technically part of MMv1; however, they are not generated from YAML configurations. Instead, they are written as Go code with minimal go template "version guards" to exclude beta-only features from the `google` provider.

Handwritten resources and datasources can be grouped by "service", which generally corresponds to the API subdomain the resource or datasource interacts with.

In addition to the core implementation, handwritten resources and datasources will also have documentation, tests, and sweepers (which clean up stray resources left behind by tests). Each type of code is stored in the following locations:

- Resource & datasource implementation: In the appropriate service folder inside [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services)
- Resource documentation: [`magic-modules/mmv1/third_party/terraform/website/docs/r`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/website/docs/r)
- Datasource documentation: [`magic-modules/mmv1/third_party/terraform/website/docs/d`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/website/docs/d)
- Tests: In the appropriate service folder inside [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services)
- Sweepers: [`magic-modules/mmv1/third_party/terraform/utils`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/utils)

In the providers, handwritten resources and datasources are stored in `PROVIDER/services/SERVICE/FILENAME.go`, where `PROVIDER` is `google` or `google-beta`, `SERVICE` is the service name, and `FILENAME` is the name of the handwritten file in magic-modules. Handwritten files do not have an `AUTO GENERATED CODE` header.

### DCL aka tpgtools (maintenance mode)

DCL / tpgtools is similar to MMv1; however, it is in maintenance mode, which means that new resources using the DCL are not being added.

DCL-based files start with the following header:

```
***     AUTO GENERATED CODE    ***    Type: DCL     ***
```
