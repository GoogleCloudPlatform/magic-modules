---
title: "Generate the providers"
weight: 110
aliases:
  - /docs/getting-started/setup
  - /getting-started/setup
  - /docs/getting-started/generate-providers
  - /getting-started/generate-providers
  - /get-started/generate-providers
---

# Generate `google` and `google-beta` providers

After making a change to the Terraform providers for Google Cloud, you must
integrate your changes with the providers. This page explains how to generate
provider changes to the `google` and `google-beta` Terraform providers.

## Before you begin

1. [Set up your development environment]({{< ref "/develop/set-up-dev-environment" >}}).
1. Update `magic-modules` as needed. These updates could be any of the following changes:
  + [Adding a resource]({{< ref "/develop/add-resource" >}}).
  + [Adding a datasource]({{< ref "/develop/add-handwritten-datasource" >}}).
  + [Adding custom resource code]({{< ref "/develop/custom-code" >}}).
  + [Promoting a resource to GA]({{< ref "/develop/promote-to-ga" >}}).

## Generate a provider change

1. Clone the `google` and `google-beta` provider repositories with the following commands:

   ```bash
   git clone https://github.com/hashicorp/terraform-provider-google.git $GOPATH/src/github.com/hashicorp/terraform-provider-google
   git clone https://github.com/hashicorp/terraform-provider-google-beta.git $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
   ```
1. Generate changes for the `google` provider:
   ```bash
   make provider VERSION=ga OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google" PRODUCT=[PRODUCT_NAME]
   ```
    Where `[PRODUCT_NAME]` is one of the folder names in
    https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products.
  
    For example, if your product is `bigqueryanalyticshub`, the command would be
    the following:

     ```bash
     make provider VERSION=ga OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google" PRODUCT=bigqueryanalyticshub
     ```

1. Generate changes for the `google-beta` provider:
   ```bash
   make provider VERSION=beta OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google-beta" PRODUCT=[PRODUCT_NAME]
   ```

    Where `[PRODUCT_NAME]` is one of the folder names in https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products.
   
    For example, if your product name is `bigqueryanalyticshub`, the command would be the following:

     ```bash
     make provider VERSION=beta OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google-beta" PRODUCT=bigqueryanalyticshub
     ```
 
1. Confirm that the expected changes were generated:
   ```bash
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
   git diff -U0
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
   git diff -U0
   ```


   {{< hint info >}}
   **Note**: There may be additional changes present due to specifying a
   `PRODUCT=` value or due to the `magic-modules` repository being out of sync
   with the provider repositories.
   {{< /hint >}}


## Troubleshoot

### Too many open files {#too-many-open-files}

If you are getting “Too many open files” ulimit needs to be raised.

{{< tabs "ulimit" >}}
{{< tab "Mac OS" >}}
```bash
ulimit -n 8192
```
{{< /tab >}}
{{< /tabs >}}

## What's next

+ [Learn how to add resource tests]({{< ref "/test/test" >}})
+ [Learn how to run tests]({{< ref "/test/run-tests" >}})
+ [Learn about `make` commands]({{< ref "/reference/make-commands" >}})