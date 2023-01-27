---
title: Generate the providers
weight: 20
---

# Generate the provider


You can compile the Terraform provider you're working on using `make` or
`make provider` in the magic-modules directory. Below you'll find a command reference,
sample commands, and instructions on how to clean up the repository.

## `make provider` reference
`make` or `make provider` will build the terraform provider. The following are variables
set either as environment or inline when calling `make`. They will inform the generator where
and how to generate.

Note: Generation is done by running definitions through templates, then unioning these generated files with the
actual provider downstream. Thus please ensure the provider you are generating into and magic-modules are in-sync
and/or up to date.

{{< build_table >}}

## Sample commands
Run the following commands from the root directory of the repository.
OUTPUT_PATH should be set to the location of your provider repository, which
is recommended to be inside your GOPATH.

```bash
cd magic-modules

make provider VERSION=ga OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google"
make provider VERSION=beta OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google-beta"

# Only generate a specific product (plus all common files)
make provider VERSION=ga OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google" PRODUCT=pubsub

# Only generate only a specific resources for a product
make provider VERSION=ga OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google" PRODUCT=pubsub RESOURCE=Topic

# Only generate common files, including all third_party code
make provider VERSION=ga OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google" PRODUCT=foo
```

## Cleaning up old files

Magic Modules will only generate on top of whatever is in the downstream repository. This means that, from time
to time, you may end up with stale files or changes in your downstream that cause compilation or tests to fail.

You can clean up by running the following commands in your downstream repository:

```bash
git checkout -- .
git clean -f google/ google-beta/ website/
```
