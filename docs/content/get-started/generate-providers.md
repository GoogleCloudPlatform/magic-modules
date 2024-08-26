---
title: "Generate the providers"
weight: 10
aliases:
  - /docs/getting-started/setup
  - /getting-started/setup
  - /docs/getting-started/generate-providers
  - /getting-started/generate-providers
---


# Generate `google` and `google-beta` providers

This quickstart guides you through setting up your development environment, making a change to `magic-modules`, generating provider changes to the `google` and `google-beta` Terraform providers, and running tests related to the change.

## Before you begin

1. [Install the gcloud CLI.](https://cloud.google.com/sdk/docs/install)
1. In the Google Cloud console, on the project selector page, select or [create a Google Cloud project](https://cloud.google.com/resource-manager/docs/creating-managing-projects).
   {{< hint info >}}
   **Note:** If you don't already have a project to use for testing changes to the Terraform providers, create a project instead of selecting an existing poject. After you finish these steps, you can delete the project, removing all resources associated with the project.
   {{< /hint >}}
   {{< button href="https://console.cloud.google.com/projectselector2/home/dashboard" >}}Go to project selector{{< /button >}}
1. Make sure that billing is enabled for your Google Cloud project. Learn how to [check if billing is enabled on a project](https://cloud.google.com/billing/docs/how-to/verify-billing-enabled).

## Set up your development environment

{{< hint warning >}}
If you are familiar with Docker or Podman, you may want to use the experimental [container-based environment]({{< ref "/reference/make-commands.md#container-based-environment" >}}) instead of this section.
{{< /hint >}}

1. [Install git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
1. [Install rbenv](https://github.com/rbenv/rbenv#installation), ensuring you follow **both** steps 1 and 2. 
1. Use rbenv to install ruby 3.1.0
   ```bash
   rbenv install 3.1.0
   ```
1. [Install go](https://go.dev/doc/install)
1. Add the following values to your environment settings such as `.bashrc`:
   ```bash
   # Add GOPATH variable for convenience
   export GOPATH=$(go env GOPATH)
   # Add Go binaries to PATH
   export PATH=$PATH:$(go env GOPATH)/bin
   ```
1. Install goimports
   ```bash
   go install golang.org/x/tools/cmd/goimports@latest
   ```
1. [Install terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)
1. Clone the `magic-modules` repository
   ```bash
   cd ~
   git clone https://github.com/GoogleCloudPlatform/magic-modules.git
   ```
1. Run the following command from the root of your cloned `magic-modules` repository.
  
   ```bash
   cd magic-modules
   ./scripts/doctor
   ```
 
   Expected output if everything is installed properly:
 
   ```
   Check for ruby in path...
      found!
   Check for go in path...
      found!
   Check for goimports in path...
      found!
   Check for git in path...
      found!
   Check for terraform in path...
      found!
   Check for make in path...
      found!
   ```

## Generate a provider change

1. In your cloned magic-modules repository, edit `mmv1/products/bigqueryanalyticshub/DataExchange.yaml` to change the description for the `displayName` field:
   ```yaml
   - !ruby/object:Api::Type::NestedObject
     name: 'displayName'
     description: |
       UPDATED_DESCRIPTION
   ```
1. Clone the `google` and `google-beta` provider repositories with the following commands:

   ```bash
   git clone https://github.com/hashicorp/terraform-provider-google.git $GOPATH/src/github.com/hashicorp/terraform-provider-google
   git clone https://github.com/hashicorp/terraform-provider-google-beta.git $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
   ```
1. Generate changes for the `google` provider
   ```bash
   make provider VERSION=ga OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google" PRODUCT=bigqueryanalyticshub
   ```
1. Generate changes for the `google-beta` provider
   ```bash
   make provider VERSION=beta OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google-beta" PRODUCT=bigqueryanalyticshub
   ```
1. Confirm that the expected changes were generated
   ```bash
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
   git diff -U0
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
   git diff -U0
   ```

   In both cases, the changes should include:

   ```diff
   diff --git a/google/services/bigqueryanalyticshub/resource_bigquery_analytics_hudiff --git a/google/services/bigqueryanalyticshub/resource_bigquery_analytics_hub_data_exchange.go b/google/services/bigqueryanalyticshub/resource_bigquery_analytics_hub_data_exchange.go
   --- a/google/services/bigqueryanalyticshub/resource_bigquery_analytics_hub_data_exchange.go
   +++ b/google/services/bigqueryanalyticshub/resource_bigquery_analytics_hub_data_exchange.go
   @@ -66 +66 @@ func ResourceBigqueryAnalyticsHubDataExchange() *schema.Resource {
   -                               Description: `Human-readable display name of the data exchange. The display name must contain only Unicode letters, numbers (0-9), underscores (_), dashes (-), spaces ( ), and must not start or end with spaces.`,
   +                               Description: `UPDATED_DESCRIPTION`,
   diff --git a/website/docs/r/bigquery_analytics_hub_data_exchange.html.markdown b/website/docs/r/bigquery_analytics_hub_data_exchange.html.markdown
   --- a/website/docs/r/bigquery_analytics_hub_data_exchange.html.markdown
   +++ b/website/docs/r/bigquery_analytics_hub_data_exchange.html.markdown
   @@ -63 +63 @@ The following arguments are supported:
   -  Human-readable display name of the data exchange. The display name must contain only Unicode letters, numbers (0-9), underscores (_), dashes (-), spaces ( ), and must not start or end with spaces.
   +  UPDATED_DESCRIPTION
   ```

   {{< hint info >}}
   **Note**: There may be additional changes present due to specifying a `PRODUCT=` value or due to the `magic-modules` repository being out of sync with the provider repositories. This is okay as long as tests in the following section pass.
   {{< /hint >}}


## Test changes

1. Set up application default credentials for Terraform
   ```bash
   gcloud auth application-default login
   export GOOGLE_USE_DEFAULT_CREDENTIALS=true
   ```
1. Set required environment variables
   ```bash
   export GOOGLE_PROJECT=PROJECT_ID
   export GOOGLE_REGION=us-central1
   export GOOGLE_ZONE=us-central1-a
   ```
   Replace `PROJECT_ID` with the ID of your Google Cloud project.

1. Enable required APIs
   ```bash
   gcloud config set project $GOOGLE_PROJECT
   gcloud services enable analyticshub.googleapis.com
   ```
1. Run all linters
   ```bash
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
   make lint
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
   make lint
   ```
1. Run all unit tests
   ```bash
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
   make test
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
   make test
   ```
1. Run acceptance tests for BigqueryAnalyticsHub DataExchange

   ```bash
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
   make testacc TEST=./google/services/bigqueryanalyticshub TESTARGS='-run=TestAccBigqueryAnalyticsHubDataExchange_'
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
   make testacc TEST=./google-beta/services/bigqueryanalyticshub TESTARGS='-run=TestAccBigqueryAnalyticsHubDataExchange_'
   ```

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

## Cleanup

1. Optional: Revoke credentials from the gcloud CLI.

```bash
gcloud auth revoke
```

## What's next

- [Learn about Magic Modules]({{< ref "/get-started/how-magic-modules-works.md" >}})
- [Learn about the contribution process]({{< ref "/get-started/contribution-process.md" >}})
- [Learn about make commands]({{< ref "/reference/make-commands.md" >}})
