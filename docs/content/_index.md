---
title: "Home"
weight: 0
type: "docs"
date: 2022-11-14T09:50:49-08:00
---


# Magic Modules

Magic Modules is a code generator and CI system that's used to develop the Terraform providers
for Google Platform, [`google`](https://github.com/hashicorp/terraform-provider-google) (or TPG) and
[`google-beta`](https://github.com/hashicorp/terraform-provider-google-beta) (or TPGB).

Magic Modules allows contributors to make changes against a single codebase and develop both
provider versions simultaneously. After sending a pull request against this repository, the
`modular-magician` robot user will manage (most of) the heavy lifting from generating a
complete output, running presubmit tests, and updating the providers following your
change.

---

- [Magic Modules](#magic-modules)
  - [Getting Started with Magic Modules](#getting-started-with-magic-modules)
    - [Cloning Terraform providers](#cloning-terraform-providers)
    - [Setting up a container-based environment](#setting-up-a-container-based-environment)
    - [Preparing your environment manually](#preparing-your-environment-manually)
    - [Generating the Terraform Providers](#generating-the-terraform-providers)
    - [Testing](#testing)
      - [Using released terraform binary with local provider binary](#using-released-terraform-binary-with-local-provider-binary)
    - [Submitting a PR](#submitting-a-pr)
  - [Contributing](#contributing)
    - [General contributing steps](#general-contributing-steps)
    - [Detailed contributing guide](#detailed-contributing-guide)
  - [Glossary](#glossary)
  - [Other Resources](#other-resources)

---

## Getting Started with Magic Modules

### Cloning Terraform providers

If you're generating the Terraform providers (`google` and `google-beta`),
you'll need to check out the repo(s) you're generating in your GOPATH. For
example:

```bash
git clone https://github.com/hashicorp/terraform-provider-google.git $GOPATH/src/github.com/hashicorp/terraform-provider-google
git clone https://github.com/hashicorp/terraform-provider-google-beta.git $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
```

Or run the following to check them all out:

```bash
./tools/bootstrap
```

Magic Modules won't work with old versions of the Terraform provider repos. If
you're encountering issues with vendoring and paths, make sure both MM and the
Terraform provider are running on up to date copies of `main`.

### Setting up a container-based environment

*NOTE* this approach is in beta and still collecting feedback. Please file an issue if you encounter challenges, and try pulling the latest container (see command below) first to see if any recent changes may fix you.

You do not need to run these instructions if you are setting up your environment manually.

For ease of contribution, we provide containers with the required dependencies for building magic-modules, as well as the option to build them yourself.

You can work with containers with either [Podman](https://podman.io/) or [Docker](https://docker.io/).

[scripts/make-in-container.sh](scripts/make-in-container.sh) includes all the
instructions to build and run containers by hand. Refer to the script for
individual steps, but for most users you only have to run the script directly, as a replacement for make.

Here is an example of how to build Terraform (after [cloning the provider](#cloning-terraform-providers)):

```shell
./scripts/make-in-container.sh \
  terraform VERSION=ga \
  OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google"
```

Generally, you can replace any reference to `make` in this guide with `scripts/make-in-container.sh`.

### Preparing your environment manually

*NOTE*: you don't need to run these instructions if you are using a container-based environment.

You can also build magic-modules
within your local development environment.

To get started, you'll need:

* Go
  * If you're using a Mac with Homebrew installed, you can follow these
    instructions to set up Go: [YouTube video](https://www.youtube.com/watch?v=VQVyvulNnzs).
  * If you're using Cloud Shell, Go is already installed.
  * Currently it's recommended to use Go 1.18, Go 1.19 changed the gofmt rules which causes some gofmt issue and our CIs are all on 1.18.X
* Ruby 2.6.0
  * You can use [`rbenv`](https://github.com/rbenv/rbenv) to manage your Ruby version(s).
  * To install `rbenv`:
    * Homebrew: run `brew install rbenv ruby-build`
    * Debian, Ubuntu, and their derivatives: run `sudo apt install rbenv`
  * Then run `rbenv install 2.6.0`.
    * For M1 Mac users, run `RUBY_CFLAGS="-Wno-error=implicit-function-declaration" rbenv install 2.6.0`
* [`Bundler`](https://github.com/bundler/bundler)
  * This can be installed with `gem install bundler`
* Gems for magic-modules
  * This can be installed with `cd mmv1 && bundler install`
* Goimports
  * go install golang.org/x/tools/cmd/goimports / go install golang.org/x/tools/cmd/goimports@latest
* Terraform
  * [Install Terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli)
* If you are getting "Too many open files" ulimit needs to be raised.
  * Mac OSX: `ulimit -n 1000`

Now, you can verify you're ready with:

```bash
./tools/doctor
```

<details><summary>Expected output:</summary>

```
Check for rbenv in path...
   found!
Checking ruby version...
2.6.0 (set by [PATH]/magic-modules/mmv1/.ruby-version)
Check for bundler in path...
   found!
Check for go in path...
   found!
Check for goimports in path...
   found!
Check for git in path...
   found!
```
</details>

---

### Generating the Terraform Providers

Before making any changes, you can compile the Terraform provider you're working
on by running the following command. If Magic Modules has been installed
correctly, you'll get no errors.

The following commands should be run from the root directory of the repository.
OUTPUT_PATH should be set to the location of your provider repository, which is
recommended to be inside your GOPATH.

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

It's worth noting that Magic Modules will only generate new files when run
locally. The "Magician"- the Magic Modules CI system- handles deletion of old
files when creating PRs.

---

### Testing

Once you've made changes to resource definition, you can run Magic Modules
to generate changes to your provider; see
["Generating the Terraform Providers"](#generating-the-terraform-providers)
above if you need a refresher. Once it's generated, you should run the
tests to ensure your change work for the providers.

Tests generally assume the following environment variables must be set in order to run tests:

```
GOOGLE_PROJECT
GOOGLE_CREDENTIALS|GOOGLE_CLOUD_KEYFILE_JSON|GCLOUD_KEYFILE_JSON|GOOGLE_USE_DEFAULT_CREDENTIALS
GOOGLE_REGION
GOOGLE_ZONE
```

Note that the credentials you provide must be granted wide permissions on the specified project. These tests provision real resources, and require permission in order to do so. Most developers on the team grant their test service account `roles/editor` or `roles/owner` on their project. Additionally, to ensure that your tests are performed in a region and zone with wide support for GCP features, `GOOGLE_REGION` should be set to `us-central1` and `GOOGLE_ZONE` to `us-central1-a`.

Additional variable may be required for other tests, and should get flagged when running them by Go skipping the test and flagging in the output it was skipped, with a skip message explaining why. The most typical extra values required are those required for project creation:

```
GOOGLE_ORG
GOOGLE_BILLING_ACCOUNT
```

You can run tests against the provider you generated in the `OUTPUT_PATH` location. When running tests, specify which to run using `TESTARGS`, such as:

```bash
# for ga provider
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
TF_LOG=TRACE make testacc TEST=./google TESTARGS='-run=TestAccContainerNodePool_basic' > output.log

# for beta provider
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
TF_LOG=TRACE make testacc TEST=./google-beta TESTARGS='-run=TestAccContainerNodePool_basic' > output.log
```

The `TESTARGS` variable is regexp-like, so multiple tests can be run in parallel by specifying a common substring of those tests (for example, `TestAccContainerNodePool` to run all node pool tests). There are 2000+ tests, and running all of them takes over 9 hours and requires a lot of GCP quota.

Note: `TF_LOG=TRACE` is optional; it [enables verbose logging](https://www.terraform.io/docs/internals/debugging.html) during tests, including all API request/response cycles. `> output.log` redirects the test output to a file for analysis, which is useful because `TRACE` logging can be extremely verbose.

Tests will use whatever version of the `terraform` binary is found on your path. To test with multiple versions of `terraform` core, you must run the tests multiple times with different versions. You can use [`tfenv`](https://github.com/tfutils/tfenv) to manage your system `terraform` versions.

#### Using released terraform binary with local provider binary

After the provider code is generated into the location at `OUTPUT_PATH`, you need to build the code in your generated provider
to compile the provider binary:
```bash
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
make build
```

Setup:

```bash
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/hashicorp/google/5.0.0/darwin_arm64
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/hashicorp/google-beta/5.0.0/darwin_arm64
ln -s $GOPATH/bin/terraform-provider-google ~/.terraform.d/plugins/registry.terraform.io/hashicorp/google/5.0.0/darwin_arm64/terraform-provider-google_v5.0.0
ln -s $GOPATH/bin/terraform-provider-google-beta ~/.terraform.d/plugins/registry.terraform.io/hashicorp/google-beta/5.0.0/darwin_arm64/terraform-provider-google-beta_v5.0.0
```

NOTE: Substitute your machine architecture for `darwin_arm64`, i.e. `darwin_amd64` or `linux_amd64`

Once this setup is complete, terraform will automatically use the binaries generated by the `make build` commands in the `terraform-provider-google` and `terraform-provider-google-beta` repositories instead of downloading the latest versions. To undo this, you can run:

```bash
rm -rf ~/.terraform.d/plugins/registry.terraform.io/hashicorp/
```

For more information, check out Hashicorp's documentation on the [0.13+ filesystem layout](https://www.terraform.io/upgrade-guides/0-13.html#new-filesystem-layout-for-local-copies-of-providers).

If multiple versions are available in a plugin directory (for example after `terraform providers mirror` is used), Terraform will pick the most up-to-date provider version within version constraints. As such, we recommend using a version that is several major versions ahead for your local copy of the provider, such as `5.0.0`.

---

### Submitting a PR

Before creating a commit, if you've modified any .rb files, make sure you run
`rake test`! That will run rubocop to ensure that the code you've written will
pass Travis.

Once you've created your commit(s), you can submit the code normally as a PR in
the GitHub UI. The PR template includes some instructions to make sure we
generate good PR messages for the tools' repo histories.

Once your PR is submitted, one of the Magic Modules maintainers will review it.
They'll look over the code before running the "Magician", the Magic Modules CI
system that generates PRs against each tool. Each review pass, your reviewer
will run the Magician again to update the PRs against the tools.

If there are multiple tools affected, that first reviewer will be the "primary"
reviewer, and for each other affected tool a maintainer for that specific tool
will make a pass. The primary reviewer will make it clear which other
maintainers need to review, and prompt them to review your code; you will
communicate primarily with the first reviewer.

Even when multiple tools are affected, this will generally be a quick look by
that maintainer with no changes needing to be made.

Once you've gotten approvals from the primary reviewer and the reviewers for
any affected tools, the primary reviewer will merge your changes.

---

## Contributing

### General contributing steps

1. Fork `Magic Modules` repository into your GitHub account if you haven't done before.
1. Check the [issue tracker](https://github.com/hashicorp/terraform-provider-google/issues) to see whether your feature has already been requested.
   * if there's an issue and it's already has a dedicated assignee, it indicates that someone may have already started to work on a solution.
   * otherwise, you're welcome to work on the issue.
1. Check whether the resource you would like to work on already exists in the providers ([`google`](https://github.com/hashicorp/terraform-provider-google) / [`google-beta`](https://github.com/hashicorp/terraform-provider-google-beta) or [check the website](https://registry.terraform.io/providers/hashicorp/google/latest/docs)).
   * If it exists, check the header of the downstream file to identify the type of tools used to generate the resource. For some resources, the code file, the test file and the documentation file may not be generated via the same tools.
      * Generated resources like [`google_compute_address`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address) can be identified by looking in their [`Go source`](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_compute_address.go) for an `AUTO GENERATED CODE` header as well as a `Type`. "Generated resources" typically refers to just the `MMv1` type, and `DCL` type resources are considered "DCL-based". (Currently DCL-related contribution are not supported)
      * Handwritten resources like [`google_container_cluster`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/container_cluster) can be identified if they have source code present under the [`mmv1/third_party/terraform/resources`](./mmv1/third_party/terraform/resources) folder or by the absence of the `AUTO GENERATED CODE header` in their [`Go source`](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_container_cluster.go).
   * If not, decide which tool you would like to use to implement the resource.
      * MMv1 is strongly preferred over handwriting the resource unless the resource can not be generated.
      * Currently, only handwritten datasources are supported.
1. Make the actual code change.
   * The [Contribution Guide](#detailed-contributing-guide) below will guide you to the detailed instructions on how to make your change, based on the type of the change + the tool used to generate the code.
1. Build the providers that includes your change. Check [Generating the Terraform Providers](#generating-the-terraform-providers) section for details on how to generate the providers locally.
1. Test the feature against the providers you generated in the last step locally. Check [Testing Guidance](#testing) for details on how to run provider test locally. (Testing the PR locally and pushing the commit to the PR only after the tests pass locally may significantly reduce review cycles)
1. Push your changes to your `magic-modules` repo fork and send a pull request from that branch to the main branch on `magic-modules`. A reviewer will be assigned automatically to your PR. Check [Submitting a PR](#submitting-a-PR) section for details on how to submit a PR.
1. Wait until the the modules magician to generate downstream diff (which should takes about 15 mins after creating the PR) to make sure all changes are generated correctly in downstream repos.
1. Wait for the VCR test results.
   <details><summary>Get to know general workflow for VCR tests</summary>

      1. You submit your change.
      1. The recorded tests are ran against your changes by the `modular-magician`. Tests will fail if:
         1. Your PR has changed the HTTP request values sent by the provider
         1. Your PR does not change the HTTP request values, but fails on the values returned in an old recording
         1. The recordings are out of sync with the merge-base of your PR, and an unrelated contributor's change has caused a false positive
      1. The `modular-magician` will leave a message indicating the number of passing and failing VCR tests. If there is a failure, the `modular-magician` user will leave a message indicating the "`Triggering VCR tests in RECORDING mode for the following tests that failed during VCR:`" marking which tests failed.
         1. If a test does not appear related, it probably isn't!
      1. The `modular-magician` will kick off a second test run targeting only the failed tests, this time hitting the live GCP APIs. If there are tests that fail at this point, a message stating `Tests failed during RECORDING mode:` will be left indicating the tests.
         1. If a test that appears to be related to your change has failed here, it's likely your change has introduced an issue. You can view the debug logs for the test by clicking the "view" link beside the test case to attempt to debug what's going wrong.
         1. If a test appears to be completely unrelated has failed, it's possible that a GCP API has changed in a way that broke the provider or our environment capped on a quota.
   </details>

   Where possible, take a look at the logs and see if you can figure out what needs to be fixed related to your change.
   The false positive rate on these tests is extremely high between changes in the API, Cloud Build bugs, and eventual consistency issues in test recordings so we don't expect contributors to wholly interpret the results- that's the responsibility of your reviewer.
1. If your assigned reviewers does not reply/ review within a week, gently ping them on github.
1. After your PR is merged, it will be released to customers in around a week or two.

---

### Detailed contributing guide

Task          | Section
--------------|--------------
Resource      | [handwritten](./docs/handwritten/#resource) / [mmv1](./docs/mmv1/#resource)
Datasource    | [handwritten](./docs/handwritten/#datasource) / (only handwritten datasources are supported)
IAM resource  | [handwritten](./docs/handwritten/#iam-resource) / [mmv1](./docs/mmv1/#iam-resource)
Testing       | [handwritten](./docs/handwritten/#testing) / [mmv1](./docs/mmv1/#testing)
Documentation | [handwritten](./docs/handwritten/#documentation) / [mmv1](./docs/mmv1/#documentation)
Beta feature  | [handwritten](./docs/handwritten/#beta-feature) / [mmv1](./docs/mmv1/#beta-feature)

---

## Glossary

The maintainers of the repository will tend to use specific jargon to describe
concepts related to Magic Modules; here's a quick reference of what some of
those terms are.

Term          | Definition
--------------|--------------
tool          | One of the OSS DevOps projects Magic Modules generates GCP support in
provider      | Synonym for tool as referred to inside the codebase
downstream(s) | A PR created by the Magician against a tool
upstream      | A PR created against Magic Modules or the Magic Modules repo
The Magician  | The Magic Modules CI system that drives the GitHub robot `modular-magician`

---

## Other Resources

* [Extending Terraform](https://www.terraform.io/plugin)
   * [How Terraform Works](https://www.terraform.io/plugin/how-terraform-works)
   * [Writing Custom Providers / Calling APIs with Terraform Providers](https://learn.hashicorp.com/collections/terraform/providers?utm_source=WEBSITE&utm_medium=WEB_IO&utm_offer=ARTICLE_PAGE&utm_content=DOCS)
* [Terraform Glossary](https://www.terraform.io/docs/glossary)
