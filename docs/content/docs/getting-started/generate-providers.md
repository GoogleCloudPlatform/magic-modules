---
title: Generate the providers
weight: 20
---

# Generate the providers

You can compile the Terraform provider you're working on by running the following
commands from the root directory of the repository. OUTPUT_PATH should be set to
the location of your provider repository, which is recommended to be inside your GOPATH.

```bash
cd magic-modules

make terraform VERSION=ga OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google"
make terraform VERSION=beta OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google-beta"

# Only generate a specific product (plus all common files)
make terraform VERSION=ga OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google" PRODUCT=pubsub

# Only generate only a specific resources for a product
make terraform VERSION=ga OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google" PRODUCT=pubsub RESOURCE=Topic
```

The `PRODUCT` variable values correspond to folder names in `mmv1/products` for generated resources. The `RESOURCE` variable value needs to match the name of the resource inside the `api.yaml` file for that product.

Handwritten files in `mmv1/third_party` are always compiled. If you are only working on common files or third_party code, you can pass a non-existent `PRODUCT`
to reduce the generation time.

```bash
# Only generate common files, including all third_party code
make terraform VERSION=ga OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google" PRODUCT=foo
```

## Cleaning up old files

Magic Modules will only generate on top of whatever is in the downstream repository. This means that, from time
to time, you may end up with stale files or changes in your downstream that cause compilation or tests to fail.

You can clean up by running the following commands in your downstream repository:

```bash
git checkout -- .
git clean -f google/ google-beta/ website/
```