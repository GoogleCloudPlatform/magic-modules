---
title: "Run tests"
weight: 26
aliases:
  - /docs/getting-started/run-provider-tests
  - /docs/getting-started/use-built-provider
  - /docs/how-to/run-tests
  - /get-started/run-provider-tests
  - /get-started/use-built-provider
  - /getting-started/run-provider-tests
  - /getting-started/use-built-provider
  - /how-to/run-tests
---

# Run tests

## Run provider tests locally

{{< hint info >}}
**Note:** If you want to test changes you've made in Magic Modules, you need to first [generate](/magic-modules/docs/getting-started/generate-providers/) the provider you want to test.
{{< /hint >}}

### Setup

Authentication is described in more detail [here](https://github.com/hashicorp/terraform-provider-google/wiki/Developer-Best-Practices#authentication).

In order to run tests, set the following environment variables:

```
GOOGLE_PROJECT
GOOGLE_CREDENTIALS|GOOGLE_CLOUD_KEYFILE_JSON|GCLOUD_KEYFILE_JSON|GOOGLE_USE_DEFAULT_CREDENTIALS
GOOGLE_REGION
GOOGLE_ZONE
```

Note that the credentials you provide must be granted wide permissions on the specified project. These tests provision real resources, and require permission in order to do so. Most developers on the team grant their test service account `roles/editor` or `roles/owner` on their project. Additionally, to ensure that your tests are performed in a region and zone with wide support for GCP features, `GOOGLE_REGION` should be set to `us-central1` and `GOOGLE_ZONE` to `us-central1-a`.

Additional variables may be required for other tests, and should get flagged when running them by Go skipping the test and flagging in the output it was skipped, with a skip message explaining why. The most typical extra values required are those required for project creation:

```
GOOGLE_ORG
GOOGLE_BILLING_ACCOUNT
```

### Run unit tests

Unit tests (that is, tests that do not interact with the GCP API) are very fast and you can generally run them all if you have changed any of them:

```bash
## for ga provider
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
make test
make lint

## for beta provider
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
make test
make lint
```

### Run acceptance tests

You can run tests against the provider you generated in the `OUTPUT_PATH` location. When running tests, specify which to run using `TESTARGS`, such as:

```bash
## for ga provider
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
make testacc TEST=./google TESTARGS='-run=TestAccContainerNodePool_basic'

## for beta provider
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
make testacc TEST=./google-beta TESTARGS='-run=TestAccContainerNodePool_basic'
```

TESTARGS allows you to pass [testing flags](https://pkg.go.dev/cmd/go#hdr-Testing_flags) to `go test`. The most important is `-run`, which allows you to limit the tests that get run. There are 2000+ tests, and running all of them takes over 9 hours and requires a lot of GCP quota.

`-run` is regexp-like, so multiple tests can be run in parallel by specifying a common substring of those tests (for example, `TestAccContainerNodePool` to run all node pool tests).

### Debugging tests

You can [increase your test verbosity](https://www.terraform.io/docs/internals/debugging.html)  and redirect the output to a log file for analysis. This is often helpful in debugging issues.

```bash
## for ga provider
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
TF_LOG=TRACE make testacc TEST=./google TESTARGS='-run=TestAccContainerNodePool_basic' > output.log

## for beta provider
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
TF_LOG=TRACE make testacc TEST=./google-beta TESTARGS='-run=TestAccContainerNodePool_basic' > output.log
```

You can also debug tests with [Delve](https://github.com/go-delve/delve):

```bash
## Navigate to the google package within your local GCP Terraform provider Git clone.
cd $GOPATH/src/github.com/terraform-providers/terraform-provider-google/google

## Execute the dlv command to launch the test.
## Note that the --test.run flag uses the same regexp matching as go test --run.
TF_ACC=1 dlv test -- --test.v --test.run TestAccComputeRegionBackendService_withCdnPolicy
Type 'help' for list of commands.
(dlv) b google.TestAccComputeRegionBackendService_withCdnPolicy
Breakpoint 1 set at 0x1de072b for github.com/terraform-providers/terraform-provider-google/google.TestAccComputeRegionBackendService_withCdnPolicy() ./resource_compute_region_backend_service_test.go:540
(dlv) c
=== RUN   TestAccComputeRegionBackendService_withCdnPolicy
> github.com/terraform-providers/terraform-provider-google/google.TestAccComputeRegionBackendService_withCdnPolicy() ./resource_compute_region_backend_service_test.go:540 (hits goroutine(7):1 total:1) (PC: 0x1de072b)
   535:                         },
   536:                 },
   537:         })
   538: }
   539:
=> 540: func TestAccComputeRegionBackendService_withCdnPolicy(t *testing.T) {
   541:         t.Parallel()
   542:
   543:         var svc compute.BackendService
   544:         resource.Test(t, resource.TestCase{
   545:                 PreCheck:     func() { acctest.AccTestPreCheck(t) },
(dlv)
```

### Testing with different `terraform` versions

Tests will use whatever version of the `terraform` binary is found on your path. To test with multiple versions of `terraform` core, you must run the tests multiple times with different versions. You can use [`tfenv`](https://github.com/tfutils/tfenv) to manage your system `terraform` versions.

## Use built provider locally

Sometimes, for example for manual testing, you may want to build the provider from source and use it with `terraform`.

### Developer Overrides

{{< hint info >}}
**Note:** Developer overrides are only available in Terraform v0.14 and later.
{{< /hint >}}

In the sections below we describe how to create a Terraform CLI configuration file, and how to make the CLI use the file via an environment variable.

### Setup

Choose your architecture below.

{{< details "Mac (ARM64 and AMD64), Linux AMD64" >}}

First, you need to find the location where built provider binaries are created. To do this, run this command and make a note of the path value:

```bash
go env GOBIN

## If the above returns nothing, then run the command below and add "/bin" to the end of the output path.
go env GOPATH
```

Next, create an empty configuration file. This could be in your `$HOME` directory or in a project directory; location does not matter. The extension `.tfrc` is required but the file name can be whatever you choose.

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
      "hashicorp/google" = "<REPLACE-ME>/bin"
      "hashicorp/google-beta" = "<REPLACE-ME>/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Edit the file to replace `<REPLACE-ME>` with the path you saved from the first step, making sure to keep `/bin` at the end of the path.

**Please note**: the _full_ path is required and environment variables cannot be used. For example, `"/Users/MyUserName/go/bin"` is a valid path for a user called `MyUserName`, but `"~/go/bin"` or `"$HOME/go/bin"` will not work.

Finally, save the file.

{{< /details >}}

{{< details "Windows (Vista and above)" >}}

First, you need to find the location where built provider binaries are created. To do this, run this command and make a note of the path value:

```bash
echo %GOPATH%
```

Next, create an empty configuration file. The location does not matter and could be in your home directory or a specific project directory. The extension `.tfrc` is required but the file name can be whatever you choose.

If you are unsure where to put the file, put it in the `%APPDATA%` directory (use `$env:APPDATA` in PowerShell to find its location on your system).

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
      "hashicorp/google" = "<REPLACE-ME>\bin"
      "hashicorp/google-beta" = "<REPLACE-ME>\bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Edit the file to replace `<REPLACE-ME>` with the output you saved from the first step, making sure to keep `\bin` at the end of the path.

**Please note**: The _full_ path is required and environment variables cannot be used. For example, `C:\Users\MyUserName\go\bin` is a valid path for a user called `MyUserName`.

Finally, save the file.

{{< /details >}}

This CLI configuration file you created in the steps above will allow Terraform to automatically use the binaries generated by the `make build` commands in the `terraform-provider-google` and `terraform-provider-google-beta` repositories instead of downloading the latest versions. All other providers will continue to be downloaded from the public Registry as normal.

### Build provider

```bash
## ga provider
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
make build

## beta provider
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
make build
```

### Using Terraform CLI developer overrides

To make Terraform use the configuration file you created, you need to set the `TF_CLI_CONFIG_FILE` environment variable to be a string containing the path to the configuration file ([see the documentation here](https://developer.hashicorp.com/terraform/cli/config/environment-variables#tf_cli_config_file)). The path can be either a relative or absolute path.

Assuming that a configuration file was created at `~/tf-dev-override.tfrc`, you can either export the environment variable or set it explicitly for each `terraform` command. Note that you need to use the full path:

```bash
# either export the environment variable for your session
export TF_CLI_CONFIG_FILE="/Users/MyUserName/tf-dev-override.tfrc"

# OR, set the environment variable value per command
TF_CLI_CONFIG_FILE="/Users/MyUserName/tf-dev-override.tfrc" terraform plan
```

To check that the developer override is working, run a `terraform plan` command and look for a warning near the start of the terminal output that looks like the example below. It is not necessary to run the terraform init command to use development overrides.

```
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp/google in /Users/MyUserName/go/bin
│  - hashicorp/google-beta in /Users/MyUserName/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying
│ changes may cause the state to become incompatible with published releases.
```

{{< hint info >}}
**Note:** Developer overrides work without you needing to alter your Terraform configuration in any way.
{{< /hint >}}

### Download production providers

To stop using developer overrides, unset the `TF_CLI_CONFIG_FILE` environment variable or stop setting it in the commands you are executing.

This will then let Terraform resume normal behaviour of pulling published provider versions from the public Registry. Any version constraints in your Terraform configuration will come back into effect. Also, you may need to run `terraform init` to download the required version of the provider into your project directory if you haven't already.

### Alternative: using a global CLI configuration file

If you do not want to use the `TF_CLI_CONFIG_FILE` environment variable, as described above, you can instead create a global version of the CLI configuration file. This configuration will be used automatically by Terraform. To do this, follow [the guidance in the official documentation](https://developer.hashicorp.com/terraform/cli/config/config-file#locations).

In this scenario you will need to remember to edit this file to swap between using developer overrides and using the production provider versions.

### Possible problems

Filesystem mirrors (particularly "[_implicit_ filesystem mirrors](https://developer.hashicorp.com/terraform/cli/config/config-file#implied-local-mirror-directories)") are used automatically by Terraform, so can interfere with the expected behaviour of Terraform if you're not aware they're present.

To stop using the filesystem mirror, you can run:

```bash
rm -rf ~/.terraform.d/plugins/registry.terraform.io/hashicorp/
```

Another way to debug this is to run a Terraform command with the `TF_LOG` environment variable set to `TRACE` . Then, look for a log line similar to the below:

```bash
[TRACE] getproviders.SearchLocalDirectory: found registry.terraform.io/hashicorp/google vX.X.X for darwin_arm64 at /Users/MyUserName/.terraform.d/plugins/registry.terraform.io/hashicorp/google/xxx
```
