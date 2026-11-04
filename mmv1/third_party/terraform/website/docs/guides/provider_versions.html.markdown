---
page_title: "Google Provider Versions"
description: |-
  How to use the different Google Provider versions
---

# Google Provider Versions

There are three versions of the Terraform Provider for Google Cloud:

* `google`

* `google-beta`

* `google-nightly`

Most users will regularly interact with `google` and `google-beta`, which cover generally available (GA) and public preview ("beta") products and APIs respectively.
* The `google` provider [documentation](https://registry.terraform.io/providers/hashicorp/google/latest/docs) covers both `google` and `google-beta`. `google-beta`-only features are clearly marked.
* `google-beta` is a superset of `google` and contains all GA products and features
* Both `google` and `google-beta` are stable providers and will only release breaking changes as part of major releases
* `google-beta` sends all requests to a beta endpoint (such as `v1beta`) for the product when one exists, even if not using any public preview features
* Many public modules freely mix `google` and `google-beta`, or use exclusively `google-beta`

-> In short, using `google-beta` over `google` is similar to using `gcloud beta`
over `gcloud`. Features that are exclusively available in `google-beta` are
Google Cloud features that are not yet GA, and they will be made available in
`google` after they go to GA.


## The `google-nightly` provider

The `google-nightly` provider is an unstable provider where fields may be added and removed at will. `google-nightly` is built off of the providers' shared source code outside of the regular release cadence and approximates the closest release version. `google-nightly` is documented at https://registry.terraform.io/providers/hashicorp/google-nightly/latest/docs

`google-nightly` supports early previews, unstable APIs, and other similar cases. Pre-GA products and features might have limited support, and changes to pre-GA products and features might not be compatible with other pre-GA versions. For more information, see the [launch stage descriptions](https://cloud.google.com/products#product-launch-stages).

We recommend pinning exact `google-nightly` versions, such as `= v2026.10.26-8.8.1`, and only using it where absolutely required, as underlying APIs or provider schemas may change without prior warning.


## Using the `google-beta` or `google-nightly` provider

To use the `google-beta` or `google-nightly` provider, set the `provider` field on each
resource where you want to use it:

```hcl
resource "google_compute_instance" "beta-instance" {
  provider = google-beta

  # ...
}
```

To customize the behavior of the `google-beta` or `google-nightly` providers, you can define their respective
provider blocks, which accept the same arguments as the `google` provider block:

```hcl
provider "google-beta" {
  project     = "my-project-id"
  region      = "us-central1"
}
```

~> If the `provider` field is omitted, Terraform will implicitly use the `google`
 provider by default even if you have only defined a `google-beta` or
 `google-nightly` provider block.

## Using multiple provider versions together

Multiple provider versions can be used together.

```hcl
resource "google_compute_instance" "ga-instance" {
  # Implicitly treated as `provider = google`

  # ...
}

resource "google_compute_instance" "beta-instance" {
  provider = google-beta

  # ...
}

resource "google_compute_instance" "nightly-instance" {
  provider = google-nightly

  # ...
}
```

Provider settings are managed per-provider and must be duplicated when using them:

```hcl
provider "google" {
  project     = "my-project-id"
  region      = "us-central1"
}

provider "google-beta" {
  project     = "my-project-id"
  region      = "us-central1"
}
```

-> Environment variables are shared between all providers.

## Importing resources with `google-beta` or `google-nightly`

Terraform will respect the defined provider automatically, similar to importing with aliased providers:


```bash
terraform import google_compute_instance.beta-instance my-instance
```

## Converting resources between versions

Resources can safely be converted from one version to the other by changing the `provider` block. For example:

To go from GA to beta, change the `provider` field from `"google"` to `"google-beta"`.

To go from beta to GA, do the reverse. If you were previously using beta features that you no longer wish to use:

1. (Optional) Explicitly set the fields back to their default values in your Terraform config file, and run `terraform apply`.
1. Change the `provider` field to `"google"`.
1. Remove any beta fields from your Terraform config.
1. Run `terraform plan` or `terraform refresh`+`terraform show` to see that the beta fields are no longer in state.
