---
title: "make commands"
weight: 10
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
go run . --version ga --output $GOPATH/src/github.com/hashicorp/terraform-provider-google
go run . --version beta --output $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta 

# Only generate a specific product (plus all common files)
go run . --version ga --product pubsub --output $GOPATH/src/github.com/hashicorp/terraform-provider-google

# Only generate common files, including all third_party code
go run . --version ga --product doesnotexist --output $GOPATH/src/github.com/hashicorp/terraform-provider-google
```

#### Arguments

- `output`: Required. The location you are generating provider code into.
- `version`: Required. The version of the provider you are building into. Valid values are `ga` and `beta`.
- `product`: Limits generations to the specified folder within `mmv1/products` or `tpgtools/api`. Handwritten files from `mmv1/third_party/terraform` are always generated into the downstream regardless of this setting, so you can provide a non-existant product name to generate only handwritten code.

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
