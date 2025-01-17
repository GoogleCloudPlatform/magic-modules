---
title: "make commands"
weight: 30
---
# `make` commands reference

## `magic-modules`

### `make` / `make provider`

Generates the code for the downstream `google` and `google-beta` providers.

{{< hint info >}}
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
- `PRODUCT`: Limits generations to the specified folder within `mmv1/products` or `tpgtools/api`. Handwritten files from `mmv1/third_party/terraform` are always generated into the downstream regardless of this setting, so you can provide a non-existent product name to generate only handwritten code. Required if `RESOURCE` is specified.
- `RESOURCE`: Limits generation to the specified resource within a particular product. For `mmv1` resources, matches the resource's `name` field (set in its configuration file).For `tpgtools` resources, matches the terraform resource name.
- `ENGINE`: Modifies `make provider` to only generate code using the specified engine. Valid values are `mmv1` or `tpgtools`. (Providing `tpgtools` will still generate any prerequisite mmv1 files required for tpgtools.)

#### Cleaning up old files

Magic Modules will only generate on top of whatever is in the downstream repository. This means that, from time
to time, you may end up with stale files or changes in your downstream that cause compilation or tests to fail.

You can clean up by running the following command in your downstream repositories:

```bash
git checkout -- . && git clean -f google/ google-beta/ website/
```

### Container-based environment

{{< hint warning >}}This approach is in beta and still collecting feedback. Please [file an issue](https://github.com/hashicorp/terraform-provider-google/issues/new/choose) if you encounter challenges.{{< /hint >}}

[`./scripts/make-in-container.sh`](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/scripts/make-in-container.sh) runs `make` with the provided arguments inside a container with all necessary dependencies preinstalled. It uses [Docker](https://docker.io/) if available and [Podman](https://podman.io/) otherwise. Like `make`, this script must be run in the root of a `magic-modules` repository clone.

If you run into any problems, please [file an issue](https://github.com/hashicorp/terraform-provider-google/issues/new/choose).

#### Before you begin

1. Ensure that `GOPATH` is set on your host machine.

   ```bash
   printenv | grep GOPATH
   ```

   If not, add `export GOPATH=$HOME/go` to your terminal's startup script and restart your terminal.
1. Clone the `google` and `google-beta` provider repositories with the following commands:

   ```bash
   git clone https://github.com/hashicorp/terraform-provider-google.git $GOPATH/src/github.com/hashicorp/terraform-provider-google
   git clone https://github.com/hashicorp/terraform-provider-google-beta.git $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
   ```

#### Example

To build the `google` provider, run the following command in the root of a `magic-modules` repository clone:

```bash
./scripts/make-in-container.sh \
  terraform VERSION=ga \
  OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google"
```
