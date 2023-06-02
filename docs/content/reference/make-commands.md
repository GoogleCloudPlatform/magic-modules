---
title: "make commands"
weight: 10
---
# `make` commands reference

## `magic-modules`

### `make` / `make provider`

Generates the code for the downstream `google` and `google-beta` providers.

{{< hint >}}
**Note:** Generation works best if the downstream provider has a commit checked out corresponding to the latest `main` branch commit that is present in your `magic-modules` working branch. This can generally be identified based on matching commit messages.
{{< /hint >}}

Examples:

```bash
make provider VERSION=ga OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google"
make provider VERSION=beta OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google-beta"

# Only generate a specific product (plus all common files)
make provider VERSION=ga OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google" PRODUCT=pubsub

# Only generate only a specific resources for a product
make provider VERSION=ga OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google" PRODUCT=pubsub RESOURCE=Topic

# Only generate common files, including all third_party code
make provider VERSION=ga OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google" PRODUCT=doesnotexist
```

#### Arguments

- `OUTPUT_PATH`: Required. The location you are generating provider code into.
- `VERSION`: Required. The version of the provider you are building into. Valid values are `ga` and `beta`.
- `PRODUCT`: Limits generations to the specified folder within `mmv1/products` or `tpgtools/api`. Handwritten files from `mmv1/third_party/terraform` are always generated into the downstream regardless of this setting, so you can provide a non-existant product name to generate only handwritten code. Required if `RESOURCE` is specified.
- `RESOURCE`: Limits generation to the specified resource within a particular product. For `mmv1` resources, matches the resource's `name` field (set in its configuration file).For `tpgtools` resources, matches the terraform resource name.
- `ENGINE`: Modifies make provider to only generate code using the specified engine. Valid values are `mmv1` or `tpgtools`. (Providing `tpgtools` will still generate any prerequisite mmv1 files required for tpgtools.)

#### Cleaning up old files

Magic Modules will only generate on top of whatever is in the downstream repository. This means that, from time
to time, you may end up with stale files or changes in your downstream that cause compilation or tests to fail.

You can clean up by running the following command in your downstream repositories:

```bash
git checkout -- . && git clean -f google/ google-beta/ website/
```

### Container-based environment

{{< hint warning >}}This approach is in beta and still collecting feedback. Please [file an issue](https://github.com/hashicorp/terraform-provider-google/issues/new/choose) if you encounter challenges.{{< /hint >}}

For ease of contribution, we provide containers with the required dependencies for building magic-modules, as well as the option to build them yourself.

[scripts/make-in-container.sh](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/scripts/make-in-container.sh) acts as a drop-in replacement for magic-modules `make` commands by setting up the containers and running `make` inside the container.

For example, to build the `google` provider:

```bash
./scripts/make-in-container.sh \
  terraform VERSION=ga \
  OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google"
```

`make-in-container.sh` will use [Docker](https://docker.io/) if available and otherwise attempt to fall back to [Podman](https://podman.io/). If you run into any problems, try pulling the latest version of `gcr.io/graphite-docker-images/downstream-builder`. If that doesn't resolve your problem, please [file an issue](https://github.com/hashicorp/terraform-provider-google/issues/new/choose).
