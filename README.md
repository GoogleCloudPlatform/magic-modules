
# Magic Modules

<img src="mmv1/images/magic-modules.svg" alt="Magic Modules Logo" width="300" align="right" />

## Overview

Magic Modules is a tool used to autogenerate support in a variety of open source DevOps
tools for Google Cloud Platform. [GCP "resource"](https://cloud.google.com/docs/overview/#gcp_resources)
definitions are encoded in a shared data file, and that data is used to fill in
tool-specific templates across each of the tools Magic Modules
generates.

Magic Modules generates GCP support for:

* Terraform

In addition, Magic Modules generates support for several companion
features/tools:

* Terraform Validator
* Terraform in Cloud Shell

Importantly, Magic Modules *isn't* full code generation. Every change is made
manually; more than a code generator, Magic Modules is a force multiplier for
development. While many Magic Modules resources are defined exactly based on the
GCP API, we use Magic Modules to preemptively solve issues across each tool by
encoding our field-tested learnings from other tools in those definitions. In
effect, an issue solved in one tool will be solved for each other tool.

## Getting Started with Magic Modules

---

<details><summary><b>Step 1: Preparing your environment</b></summary>

* Go
  * If you're using a Mac with Homebrew installed, you can follow these
    instructions to set up Go: [YouTube video](https://www.youtube.com/watch?v=VQVyvulNnzs).
  * If you're using Cloud Shell, Go is already installed.
* Ruby 2.6.0
  * You can use `rbenv` to manage your Ruby version(s).
  * To install `rbenv`, run `sudo apt install rbenv`.
  * Then run `rbenv install 2.6.0`. 
    * For M1 Mac users, run `RUBY_CFLAGS="-Wno-error=implicit-function-declaration" rbenv install 2.6.0`
* [`Bundler`](https://github.com/bundler/bundler)
  * This can be installed with `gem install bundler`
* If you are getting "Too many open files" ulimit needs to be raised.
  * Mac OSX: `ulimit -n 1000`
</details>

<details><summary><b>Step 2: Preparing Magic Modules / One-time setup</b></summary>

**Important:**
Compiling Magic Modules can be done directly from the `mmv1` directory within this repository.
In the future we will add hybrid generation with multiple generators. All the information below
pertains only to the contents of the `mmv1` directory, and commands should be executed from
that directory.

To get started right away, use the bootstrap script with:

```bash
cd mmv1
./tools/bootstrap
```

---

Otherwise, follow the manual steps below:

If you're generating the Terraform providers (`google` and `google-beta`),
you'll need to check out the repo(s) you're generating in your GOPATH. For
example:

```bash
git clone https://github.com/hashicorp/terraform-provider-google.git $GOPATH/src/github.com/hashicorp/terraform-provider-google
git clone https://github.com/hashicorp/terraform-provider-google-beta.git $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
```

Magic Modules won't work with old versions of the Terraform provider repos. If
you're encountering issues with vendoring and paths, make sure both MM and the
Terraform provider are running on up to date copies of `main`.

Once you've prepared the target folders for the tools, run the following to
finish getting Magic Modules set up by installing the Ruby gems it needs to run:

```bash
cd mmv1
bundle install
```

Now, you can verify you're ready with:

```bash
./tools/doctor
```

Expected output:

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

<details><summary><b>Step 3: Generating the Terraform Providers</b></summary>

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
make terraform VERSION=ga OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google" PRODUCT=dataproc
```

It's worth noting that Magic Modules will only generate new files when run
locally. The "Magician"- the Magic Modules CI system- handles deletion of old
files when creating PRs.

</details>

## How to Contribute

---

### General Contribution Steps


1. Fork `Magic Modules` repository into your GitHub account if you haven't done before
1. Check the [issue tracker](https://github.com/hashicorp/terraform-provider-google/issues) to see whether your feature has already been requested
   * if there's an issue and it's already has a dedicated assignee, it indicates that someone may already work on a solution.
   * otherwise, you're welcome to work on the issue
1. Check whether the resource you would like to work on already exists in the providers ([Tpg](https://github.com/hashicorp/terraform-provider-google) / [Tpgb](https://github.com/hashicorp/terraform-provider-google-beta) or [check the website](https://registry.terraform.io/providers/hashicorp/google/latest/docs))
   * If it exists, check the header of the downstream file to identify the type of tools used to generate the resource.
  For some resources, the code file, the test file and the documentation file may not be generated via the same tools. (attach photo of the header)
   * If not, decide which tool you would like to use to implement the resource. (In most cases, generat)
1. Make the actual code change
   * The [Contribution Guide](#contribution-guide) below will lead you to the detailed instructions on how to make your change, based on the type of the change + the tool used to generate the code
1. build the providers and test the feature locally. Check [Testing Guidance]() for details on how to run provider test.
1. refer to the schema table to make sure all correct schema are implemented for the feature
1. Open a Pull Request refer to the PR creation checklist
1. wait until the the modules magician to generate downstream diff (which should takes about 15 mins after creating the PR) to make sure all changes are generated correctly in downstream repos
1. wait for the VCR result
  If any tests failed during VCR recording mode
    Check if the failed tests are related (attach photo – build logs for general errors debug logs for each test http calls responses). 
If they are related, make changes to your PR (strongly suggest testing the PR locally and pushing the commit to the PR only after the tests pass locally, otherwise it may increase review time)
After you push the new commits and think the PR is ready to review again, re-request the review your assigned reviewer (attach photo to show the re-request button) (if we decide to automatically re-request the reviewer after any new pushes, this step should be removed)
If you are not able to debug the failed test, leave a comment on Github to let your assigned reviewers know where you are blocked.
    If you are not able to access the logs, gently ping your assigned reviewers on github to let them know.
1. If your assigned reviewers does not reply/ review within 48 hours, gently ping them on github


### Contribution Guide

Task          |        Section         ||
--------------|-----------------|---------
Resource      | [handwritten](./mmv1/third_party/terraform/README.md#resource) | [mmv1](./mmv1/README.md#resource) |
Datasource    | [handwritten](./mmv1/third_party/terraform/README.md#datasource) |          |
IAM resource  | [handwritten](./mmv1/third_party/terraform/README.md#iam-resource) | [mmv1](./mmv1/README.md#iam-resource) |
Testing       | [handwritten](./mmv1/third_party/terraform/README.md#testing) | [mmv1](./mmv1/README.md#testing) |
Documentation | [handwritten](./mmv1/third_party/terraform/README.md#documentation) | [mmv1](./mmv1/README.md#documentation) |
Beta feature  | [handwritten](./mmv1/third_party/terraform/README.md#beta-feature) | [mmv1](./mmv1/README.md#beta-feature) |

### Testing your changes

Once you've made changes to resource definition, you can run Magic Modules
to generate changes to your tool; see
["Generating the Terraform Providers"](#generating-the-terraform-providers)
above if you need a refresher. Once it's generated, you should run the
tool-specific tests as if you were submitting a PR against that tool.

You can run tests in the `{{output_folder}}` you generated the tool in.
See the following tool-specific documentation for more details on testing that
tool;

Tool             | Testing Guide
-----------------|--------------
terraform        | [`google` provider testing guide](https://github.com/hashicorp/terraform-provider-google/blob/main/.github/CONTRIBUTING.md#tests)
terraform (beta) | [`google-beta` provider testing guide](https://github.com/hashicorp/terraform-provider-google-beta/blob/main/.github/CONTRIBUTING.md#tests)

Don't worry about testing every tool, only the primary tool you're making
changes against. The Magic Modules maintainers will ensure your changes work
against each tool.

If your changes have unintended consequences in another tool, a reviewer will
instruct you to mark the field excluded or provide specific feedback on what
changes to make to the tool-specific overrides in order for them to work
correctly.

### Submitting a PR

Before creating a commit, if you've modified any .rb files, make sure you run
`rake test`! That will run rubocop to ensure that the code you've written will
pass Travis.

To run rubocop automatically before committing, add a Git pre-commit hook with:

```bash
cp .github/pre-commit .git/hooks
```

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