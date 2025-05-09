---
title: "Run tests"
weight: 20
aliases:
  - /docs/getting-started/run-provider-tests
  - /docs/getting-started/use-built-provider
  - /get-started/run-provider-tests
  - /get-started/use-built-provider
  - /getting-started/run-provider-tests
  - /getting-started/use-built-provider
  - /develop/run-tests
  - /develop/test/run-tests
---

# Run tests

## Before you begin

[Generate the modified provider(s)]({{< ref "/develop/generate-providers" >}})


1. Set up application default credentials for Terraform

    ```bash
    gcloud auth application-default login
    export GOOGLE_USE_DEFAULT_CREDENTIALS=true
    ```

1. Set the following environment variables:

    ```bash
    export GOOGLE_PROJECT=PROJECT_ID
    export GOOGLE_REGION=us-central1
    export GOOGLE_ZONE=us-central1-a
    ```
    Replace `PROJECT_ID` with the ID of the Google Cloud project you are using for testing.

1. Optional: Some tests may require additional variables to be set, such as:

    ```
    GOOGLE_ORG
    GOOGLE_BILLING_ACCOUNT
    ```

## Run automated tests

{{< tabs "version" >}}

{{< tab "GA Provider" >}}

1. Run unit tests and linters

    ```bash
    cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
    make test
    make lint
    ```


1. Run acceptance tests for only modified resources. (Full test runs can take over 9 hours.) See [Go's documentation](https://pkg.go.dev/cmd/go#hdr-Testing_flags) for more information about `-run` and other flags.

    ```bash
    make testacc TEST=./google/services/container TESTARGS='-run=TestAccContainerNodePool_basic$$'
    ```

    To run all tests matching, e.g., `TestAccContainerNodePool*`, omit the trailing `$$`:

    ```bash
    make testacc TEST=./google/services/container TESTARGS='-run=TestAccContainerNodePool'
    ```

    > **Note:** Acceptance tests create actual infrastructure which can incur costs. Acceptance tests may not clean up after themselves if interrupted, so you may want to check for stray resources and / or billing charges.

1. Optional: Save verbose test output (including API requests and responses) to a file for analysis.

    ```bash
    TF_LOG=DEBUG make testacc TEST=./google/services/container TESTARGS='-run=TestAccContainerNodePool_basic$$' > output.log
    ```

1. Optional: Debug tests with [Delve](https://github.com/go-delve/delve). See [`dlv test` documentation](https://github.com/go-delve/delve/blob/master/Documentation/usage/dlv_test.md) for information about available flags.

    ```bash
    cd google
    TF_ACC=1 dlv test -- --test.v --test.run TestAccComputeRegionBackendService_withCdnPolicy
    ```

{{< /tab >}}

{{< tab "Beta Provider" >}}

1. Run unit tests and linters

    ```bash
    cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
    make test
    make lint
    ```

1. Run acceptance tests for only modified resources. (Full test runs can take over 9 hours.) See [Go's documentation](https://pkg.go.dev/cmd/go#hdr-Testing_flags) for more information about `-run` and other flags.

    ```bash
    make testacc TEST=./google-beta/services/container TESTARGS='-run=TestAccContainerNodePool'
    ```

    To run all tests matching, e.g., `TestAccContainerNodePool*`, omit the trailing `$$`:

    ```bash
    make testacc TEST=./google-beta/services/container TESTARGS='-run=TestAccContainerNodePool'
    ```

    > **Note:** Acceptance tests create actual infrastructure which can incur costs. Acceptance tests may not clean up after themselves if interrupted, so you may want to check for stray resources and / or billing charges.

1. Optional: Save verbose test output to a file for analysis.

    ```bash
    TF_LOG=DEBUG make testacc TEST=./google-beta/services/container TESTARGS='-run=TestAccContainerNodePool_basic$$' > output.log
    ```

1. Optional: Debug tests with [Delve](https://github.com/go-delve/delve). See [`dlv test` documentation](https://github.com/go-delve/delve/blob/master/Documentation/usage/dlv_test.md) for information about available flags.

    ```bash
    cd google-beta
    TF_ACC=1 dlv test -- --test.v --test.run TestAccComputeRegionBackendService_withCdnPolicy
    ```

{{< /tab >}}

{{< /tabs >}}

## Troubleshooting acceptance tests {#troubleshooting}

### After applying this test step, the plan was not empty.

This indicates that the returned values from the API (which will be displayed on the left) are different than what is in the configuration (displayed on the right). Fields which are listed as "known after apply" are not the cause of the diff and can be ignored in terms of resolving the issue.

See [Fix diffs]({{< ref "/develop/diffs" >}}) for more information on potential causes and fixes.

### Blocks of type "FIELD_NAME" are not expected here

The field called `FIELD_NAME` does not exist; this is either because it has not been implemented or because the test is running for the `google` provider and the field is only implemented in the `google-beta` provider. See [Add resource tests]({{< ref "/test/test" >}}) for information on using version guards to exclude beta-only fields from GA tests, or [Promote from beta to GA]({{< ref "/develop/promote-to-ga" >}}) for information on how to promote fields that were accidentally made beta-only.

### Provider produced inconsistent result after apply ... Root object was present, but now absent.

This indicates that after an apply to create or update a resource, the resource was not present in Terraform state. This generally means one of a few things:

- [API is eventually consistent or returns an Operation]({{< ref "/develop/diffs#eventually-consistent" >}})
- The URL for reads was built incorrectly. The exact fix will depend on why this is happening. Run the test with the `TF_LOG=DEBUG` environment variable and check whether the read URL matches what you expect.
- There is a call to unset the resource's id (`d.SetId("")`) somewhere it shouldn't be. The fix is to remove that extraneous call. This is rare.

### Error: Inconsistent dependency lock file

Tests require all of the providers they use (except the one actually being tested) to be explicitly stated. This error generally means one of a few things:

- This is a beta-only test and one of the `google_*` resources in the test doesn't have `provider = google-beta` set
  - ```hcl
    resource "google_compute_instance" "beta-instance" {
      provider = google-beta
      # ...
    }
    ```

- This is a GA+beta test and one of the `google_*` resources has `provider = google-beta` set
  - `provider = google-beta` can't be set unless the test is beta-only.
- The test relies on an external provider, such as `time`, and that is not explicitly declared
  - For MMv1 example-based tests, use [`examples.external_providers`](https://googlecloudplatform.github.io/magic-modules/reference/resource/#examples).
  - For Handwritten tests, use TestCase.ExternalProviders:
    ```go
    acctest.VcrTest(t, resource.TestCase{
      ExternalProviders: map[string]resource.ExternalProvider{
        "time": {},
      },
      // ...
    }
    ```

## Optional: Test with different `terraform` versions

Tests will use whatever version of the `terraform` binary is found on your `PATH`. If you are testing a change that you know only impacts certain `terraform` versions, follow these steps:

1. Install [`tfenv`](https://github.com/tfutils/tfenv).

1. Install the version of `terraform` you want to test.

    ```bash
    tfenv install VERSION
    ```

    Replace `VERSION` with the version you want to test.

1. Run automated tests following the [earlier section]({{< ref "/test/run-tests#run-automated-tests" >}}).

## Optional: Test manually

For manual testing, you can build the provider from source and run `terraform apply` to verify the behavior.

### Before you begin

Configure Terraform to use locally-built binaries for `google` and `google-beta` instead of downloading the latest versions.

{{< tabs "built-provider" >}}

{{< tab "Developer overrides (Mac / Linux)" >}}

1. Find the location where built provider binaries are created. To do this, run this command and make a note of the path value:

    ```bash
    go env GOBIN

    ## If the above returns nothing, then run the command below and add "/bin" to the end of the output path.
    go env GOPATH
    ```

1. Create an empty configuration file.

    ```bash
    ## create an empty file
    touch ~/tf-dev-override.tfrc

    ## open the file with a text editor of your choice, e.g:
    vi ~/tf-dev-override.tfrc
    ```

    Open the empty file with a text editor and paste in these contents:

    ```hcl
    provider_installation {

      # Developer overrides will stop Terraform from downloading the listed
      # providers their origin provider registries.
      dev_overrides {
          "hashicorp/google" = "GO_BIN_PATH/bin"
          "hashicorp/google-beta" = "GO_BIN_PATH/bin"
      }

      # For all other providers, install them directly from their origin provider
      # registries as normal. If you omit this, Terraform will _only_ use
      # the dev_overrides block, and so no other providers will be available.
      direct {}
    }
    ```

1. Edit the file to replace `GO_BIN_PATH` with the path you saved from the first step, making sure to keep `/bin` at the end of the path.

    **Please note**: the _full_ path is required and environment variables cannot be used. For example, `"/Users/UserName/go/bin"` is a valid path for a user called `UserName`, but `"~/go/bin"` or `"$HOME/go/bin"` will not work.

1. Save the file.

{{< /tab >}}

{{< tab "Developer overrides (Windows)" >}}

1. Find the location where built provider binaries are created. To do this, run this command and make a note of the path value:

    ```bash
    echo %GOPATH%
    ```

1. Create an empty configuration file in the `%APPDATA%` directory (use `$env:APPDATA` in PowerShell to find its location on your system).

    ```powershell
    ## create an empty file
    type nul > "$($env:APPDATA)\tf-dev-override.tfrc"

    ## open the file with a text editor of your choice, e.g:
    notepad "$($env:APPDATA)\tf-dev-override.tfrc"
    ```

    Open the empty file with a text editor and paste in these contents:

    ```hcl
    provider_installation {

      # Developer overrides will stop Terraform from downloading the listed
      # providers their origin provider registries.
      dev_overrides {
          "hashicorp/google" = "GO_BIN_PATH\bin"
          "hashicorp/google-beta" = "GO_BIN_PATH\bin"
      }

      # For all other providers, install them directly from their origin provider
      # registries as normal. If you omit this, Terraform will _only_ use
      # the dev_overrides block, and so no other providers will be available.
      direct {}
    }
    ```

1. Edit the file to replace `GO_BIN_PATH` with the output you saved from the first step, making sure to keep `\bin` at the end of the path.

    **Please note**: The _full_ path is required and environment variables cannot be used. For example, `C:\Users\UserName\go\bin` is a valid path for a user called `UserName`.

1. Save the file.

{{< /tab >}}

{{< /tabs >}}

### Run manual tests

1. [Generate the provider(s) you want to test]({{< ref "/develop/generate-providers" >}})
2. Build the provider(s) you want to test

    ```bash
    ## google provider
    cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
    make build

    ## google-beta provider
    cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
    make build
    ```

1. Create a new directory and a `main.tf` file with your resource and its dependencies.

1. In the new directory, run `terraform plan` as follows:
    ```bash
    TF_CLI_CONFIG_FILE="$HOME/tf-dev-override.tfrc" terraform plan
    ```

    Replace the TF_CLI_CONFIG_FILE value with the full path to your developer overrides file.
1. Optional: Verify that developer overrides are working by looking for output like the following near the start of the output:

    ```
    │ Warning: Provider development overrides are in effect
    │ 
    │ The following provider development overrides are set in the CLI configuration:
    │  - hashicorp/google in /Users/UserName/go/bin
    │  - hashicorp/google-beta in /Users/UserName/go/bin
    │ 
    │ The behavior may therefore not match any released version of the provider and applying
    │ changes may cause the state to become incompatible with published releases.
    ```

1. Run `terraform apply` with developer overrides.

    ```bash
    TF_CLI_CONFIG_FILE="$HOME/tf-dev-override.tfrc" terraform apply
    ```

1. Optional: Save verbose `terraform apply` output (including API requests and responses) to a file for analysis.

    ```bash
    TF_LOG=DEBUG TF_LOG_PATH=output.log TF_CLI_CONFIG_FILE="$HOME/tf-dev-override.tfrc" terraform apply
    ```

### Run Tests with VCR Locally

VCR tests record HTTP request/response interactions in cassettes and replay them in future runs without calling the real API.

Running tests in `REPLAYING` mode locally can sometimes be useful. In particular, it can allow you to test more quickly, cheaply, and without spinning up real infrastructure, once you've got an initial recording.

It can also be helpful for debugging tests that seem to work locally, but fail in CI in replaying mode.
   
VCR is controlled via two variables:
- `VCR_MODE`: `REPLAYING` or `RECORDING` mode 
- `VCR_PATH`: Path where recorded cassettes are stored. 

Ensure both variables are configured to properly trigger VCR tests locally.

If you don't already have an existing cassette that's up to date, first do a run in `RECORDING` mode:
   
```bash
VCR_PATH=$HOME/.vcr/ VCR_MODE=RECORDING  make testacc TEST=./google/services/alloydb TESTARGS='-run=TestAccContainerNodePool_basic$$'
```

Now run the same test again in `REPLAYING` mode:

```bash
VCR_PATH=$HOME/.vcr/ VCR_MODE=REPLAYING make testacc TEST=./google/services/alloydb TESTARGS='-run=TestAccContainerNodePool_basic$$'
```

### Cleanup

To stop using developer overrides, stop setting `TF_CLI_CONFIG_FILE` in the commands you are executing.

Terraform will resume its normal behaviour of pulling published provider versions from the public Registry. Any version constraints in your Terraform configuration will come back into effect. Also, you may need to run `terraform init` to download the required version of the provider into your project directory if you haven't already.

## What's next?

[Create a pull request]({{< ref "/code-review/create-pr" >}})
