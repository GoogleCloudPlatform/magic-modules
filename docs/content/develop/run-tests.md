---
title: "Run tests"
weight: 70
aliases:
  - /docs/getting-started/run-provider-tests
  - /docs/getting-started/use-built-provider
  - /get-started/run-provider-tests
  - /get-started/use-built-provider
  - /getting-started/run-provider-tests
  - /getting-started/use-built-provider
---

# Run tests

## Before you begin

1. [Generate the modified provider(s)](/magic-modules/docs/getting-started/generate-providers/)


1. Set up application default credentials for Terraform

    ```bash
    gcloud auth application-default login
    export GOOGLE_USE_DEFAULT_CREDENTIALS=TRUE
    ```

1. Set the following environment variables:

    ```bash
    export GOOGLE_USE_DEFAULT_CREDENTIALS=TRUE
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
    make testacc TEST=./google TESTARGS='-run=TestAccContainerNodePool'
    ```



1. Optional: Save verbose test output to a file for analysis.

    ```bash
    TF_LOG=TRACE make testacc TEST=./google TESTARGS='-run=TestAccContainerNodePool_basic' > output.log
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
    make testacc TEST=./google-beta TESTARGS='-run=TestAccContainerNodePool'
    ```



1. Optional: Save verbose test output to a file for analysis.

    ```bash
    TF_LOG=TRACE make testacc TEST=./google-beta TESTARGS='-run=TestAccContainerNodePool_basic' > output.log
    ```

1. Optional: Debug tests with [Delve](https://github.com/go-delve/delve). See [`dlv test` documentation](https://github.com/go-delve/delve/blob/master/Documentation/usage/dlv_test.md) for information about available flags.

    ```bash
    cd google-beta
    TF_ACC=1 dlv test -- --test.v --test.run TestAccComputeRegionBackendService_withCdnPolicy
    ```

{{< /tab >}}

{{< /tabs >}}

## Optional: Test with different `terraform` versions

Tests will use whatever version of the `terraform` binary is found on your `PATH`. If you are testing a change that you know only impacts certain `terraform` versions, follow these steps:

1. Install [`tfenv`](https://github.com/tfutils/tfenv).

1. Install the version of `terraform` you want to test.

    ```bash
    tfenv install VERSION
    ```

    Replace `VERSION` with the version you want to test.

1. Run automated tests following the [earlier section]({{< ref "/develop/run-tests#run-automated-tests" >}}).

## Optional: Test manually

Sometimes, for example for manual testing, you may want to build the provider from source and use it with `terraform`.

### Build provider

```bash
## ga provider
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
make build

## beta provider
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
make build
```

{{< tabs "built-provider" >}}

{{< tab "Developer overrides" >}}

In the sections below we describe how to create a Terraform CLI configuration file, and how to make the CLI use the file via an environment variable.

### Create developer overrides file

Choose your architecture below.

{{< tabs "setup" >}}

{{< tab "Mac (ARM64 and AMD64), Linux AMD64" >}}

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

{{< tab "Windows (Vista and above)" >}}

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

This CLI configuration file you created in the steps above will allow Terraform to use the binaries generated by the `make build` commands in the `terraform-provider-google` and `terraform-provider-google-beta` repositories instead of downloading the latest versions.

### Using Terraform CLI developer overrides

1. Create a new directory and a `main.tf` file with your resource and its dependencies.

1. To make Terraform use the configuration file you created, you need to set the `TF_CLI_CONFIG_FILE` environment variable to the path to the configuration file ([see the documentation here](https://developer.hashicorp.com/terraform/cli/config/environment-variables#tf_cli_config_file)).

    Assuming that a configuration file was created at `~/tf-dev-override.tfrc`, you can either export the environment variable or set it explicitly for each `terraform` command. Note that you need to use the full path:

    ```bash
    # either export the environment variable for your session
    export TF_CLI_CONFIG_FILE="/Users/UserName/tf-dev-override.tfrc"
    terraform plan

    # OR, set the environment variable value per command
    TF_CLI_CONFIG_FILE="/Users/UserName/tf-dev-override.tfrc" terraform plan
    ```

1. To check that the developer override is working, look for a warning near the start of the terminal output that looks like the example below. It is not necessary to run the terraform init command to use development overrides.

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

1. Apply your configuration.

    ```bash
    TF_LOG=DEBUG TF_LOG_PATH=output.log terraform apply
    ```

{{< hint info >}}
**Note:** Developer overrides work without you needing to alter your Terraform configuration in any way.
{{< /hint >}}

### Download production providers

To stop using developer overrides, unset the `TF_CLI_CONFIG_FILE` environment variable or stop setting it in the commands you are executing.

This will then let Terraform resume normal behaviour of pulling published provider versions from the public Registry. Any version constraints in your Terraform configuration will come back into effect. Also, you may need to run `terraform init` to download the required version of the provider into your project directory if you haven't already.

### Alternative: using a global CLI configuration file

If you do not want to use the `TF_CLI_CONFIG_FILE` environment variable, as described above, you can instead create a global version of the CLI configuration file. This configuration will be used automatically by Terraform. To do this, follow [the guidance in the official documentation](https://developer.hashicorp.com/terraform/cli/config/config-file#locations). For Windows, the file will be named `terraform.rc` and go in the `%APPDATA` directory. For all other systems, it will be named `.terraformrc` and go in the user's home directory.

In this scenario you will need to remember to edit this file to swap between using developer overrides and using the production provider versions.

{{< /tab >}}

{{< tab "Filesystem mirrors (Terraform 0.13 and earlier)" >}}

Filesystem mirrors can used explicitly or implicitly by Terraform. Explicit filesystem mirrors can be [defined via the CLI configuration file](https://developer.hashicorp.com/terraform/cli/config/config-file#filesystem_mirror). In contrast, once [implicit filesystem mirrors](https://developer.hashicorp.com/terraform/cli/config/config-file#implied-local-mirror-directories) are created by a user they are discovered and used by Terraform automatically.

Filesystem mirrors require providers' files to be saved with specific paths for them to work correctly. To help with this, you can use the [`terraform providers mirror` command](https://developer.hashicorp.com/terraform/cli/commands/providers/mirror) to download a published provider to your local filesystem with the required file path.

Implied filesystem mirrors require manual cleanup if you want to revert back to using providers downloaded from the public Registry, and if an implied filesystem mirror is in place that a user is unaware of it can lead to confusing behaviour that is hard to debug. Other disadvantages compared to developer overrides include:
- No warning in terminal output letting you know when your local files are in use.
- Need to set a version number for your local version of the provider which is compatible with your Terraform configuration's version constraints
- Setup and cleanup required each time you swap to and from using them

### Possible problems

Filesystem mirrors (particularly "[_implicit_ filesystem mirrors](https://developer.hashicorp.com/terraform/cli/config/config-file#implied-local-mirror-directories)") are used automatically by Terraform, so can interfere with the expected behaviour of Terraform if you're not aware they're present.

To stop using the filesystem mirror, you can run:

```bash
rm -rf ~/.terraform.d/plugins/registry.terraform.io/hashicorp/
```

Another way to debug this is to run a Terraform command with the `TF_LOG` environment variable set to `TRACE` . Then, look for a log line similar to the below:

```bash
[TRACE] getproviders.SearchLocalDirectory: found registry.terraform.io/hashicorp/google vX.X.X for darwin_arm64 at /Users/UserName/.terraform.d/plugins/registry.terraform.io/hashicorp/google/xxx
```

{{< /tab >}}

{{< /tabs >}}
