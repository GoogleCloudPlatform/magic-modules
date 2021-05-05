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
* Ansible
* InSpec

In addition, Magic Modules generates support for several companion
features/tools:

* Terraform Google Inventory Mapper
* Terraform in Cloud Shell

Importantly, Magic Modules *isn't* full code generation. Every change is made
manually; more than a code generator, Magic Modules is a force multiplier for
development. While many Magic Modules resources are defined exactly based on the
GCP API, we use Magic Modules to preemptively solve issues across each tool by
encoding our field-tested learnings from other tools in those definitions. In
effect, an issue solved in one tool will be solved for each other tool.

## Getting Started with Magic Modules

You can try out Magic Modules immediately with Open in Cloud Shell below; if
you're getting set up on a local workstation, this guide serves as a reference
to help you get it set up.

[![Open in Cloud Shell](http://gstatic.com/cloudssh/images/open-btn.svg)](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/magic-modules&tutorial=mmv1/TUTORIAL.md)

### Requirements

To get started, you'll need:

* Ruby 2.6.0
  * You can use `rbenv` to manage your Ruby version(s)
* [`Bundler`](https://github.com/bundler/bundler)
  * This can be installed with `gem install bundler`
* If you are getting "Too many open files" ulimit needs to be raised.
  * Mac OSX: `ulimit -n 1000`

### Preparing Magic Modules / One-time setup

**Important:**
Compiling Magic Modules can be done directly from the `mmv1` directory within this repository.
In the future we will add hybrid generation with multiple generators. All the information below
pertains only to the contents of the `mmv1` directory, and commands should be executed from
that directory.


To get started right away, use the bootstrap script with:

```bash
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
Terraform provider are running on up to date copies of `master`.

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

### Generating downstream tools

Before making any changes, you can compile the "downstream" tool you're working
on by running the following command. If Magic Modules has been installed
correctly, you'll get no errors when you run a command:

```bash
bundle exec compiler -a -v "ga" -e {{tool}} -o "{{output_folder}}"
```

Generally, you'll want to generate into the same output.  For terraform, that
will be `$GOPATH/src/github.com/hashicorp/terraform-provider-google` (optionally `-beta`).
For Ansible and Inspec, wherever you have cloned those repositories.

For example, to generate Terraform:

```bash
bundle exec compiler -a -v "ga" -e terraform -o "$GOPATH/src/github.com/hashicorp/terraform-provider-google"
```

It's worth noting that Magic Modules will only generate new files when ran
locally. The "Magician"- the Magic Modules CI system- handles deletion of old
files when creating PRs.

#### Compiler options

`-e`, `-v`, and `-f` let you select which project should be generated.

Target                      | compiler options
----------------------------|-----------------
ansible                     | `-e ansible`
inspec                      | `-e inspec -v "beta"`
terraform                   | `-e terraform -v "ga"`
terraform (beta)            | `-e terraform -v "beta"`
terraform-google-conversion | `-e terraform -f validator`

Other important options are:

- `-a` Generate for all products
- `-p products/<folder_name>` Generate for a specific project, i.e. `-p products/appengine`

### Making changes to resources

Once again, see the Open in Cloud Shell example above for an interactive example
of making a Magic Modules change; this section will serve as a reference more
than a specific example.

Magic Modules mirrors the GCP REST API; there are [products](https://github.com/GoogleCloudPlatform/magic-modules/blob/master/mmv1/api/product.rb)
such as Compute or Container (GKE) that contains [resources](https://github.com/GoogleCloudPlatform/magic-modules/blob/master/mmv1/api/resource.rb),
[GCP resources](https://cloud.google.com/docs/overview/#gcp_resources) such as
Compute VM Instances or GKE Clusters.

Products are separate folders under [`products/`], and each folder contains a
file named `api.yaml` that contains the resources that make up the API
definition.

Resources are made up of some metadata like their `"name"` in the API such as
Address or Instance, some additional metadata (see the fields in [resource.rb](https://github.com/GoogleCloudPlatform/magic-modules/blob/master/mmv1/api/resource.rb)),
and the meat of a resource, its fields. They're represented by `properties` in
Magic Modules, an array of [types](https://github.com/GoogleCloudPlatform/magic-modules/blob/master/mmv1/api/type.rb).

Adding a new field to a resource in Magic Modules is often as easy as adding a
`type` to the `properties` array for the resource. See [this example](https://github.com/GoogleCloudPlatform/magic-modules/pull/1126/files#diff-fb4f76e7d870258668a3beac48bf164c)
where a field was added to all the tools (currently only Terraform) that support
beta fields.

#### Tool-specific overrides

While most small changes won't require fiddling with overrides, each tool has
"overrides" when it needs to deviate from the definition in `api.yaml`. This is
often minor differences- the naming of a field, or whether it's required or not.

You can find them under the folder for a product, with the name `{{tool}}.yaml`.
For example, Ansible's overrides for Cloud SQL are present at `products/sql/ansible.yaml`

You can find a full reference for each tool under `overrides/{{tool}}/resource_override.rb`
and `overrides/{{tool}}/property_override.rb`, as well as some other tool-specific
functionality.

#### Making changes to handwritten files

The Google providers for Terraform have a large number of handwritten files,
written before Magic Modules was used with them. While conversion is ongoing,
many resources are still managed by hand. You can modify handwritten files
under the `third_party/terraform` directory.

Features that are only present in certain versions need to be "guarded" by
wrapping those lines of code in version guards;

```erb
<% unless version == 'ga' -%>
  // beta-only code
<% end -%>
```

### Testing your changes

Once you've made changes to resource definition, you can run Magic Modules
to generate changes to your tool; see "Generating downstream tools" above if
you need a refresher. Once it's generated, you should run the tool-specific
tests as if you were submitting a PR against that tool.

You can run tests in the `{{output_folder}}` you generated the tool in.
See the following tool-specific documentation for more details on testing that
tool;

Tool             | Testing Guide
-----------------|--------------
ansible          | [instructions](https://docs.ansible.com/ansible/devel/dev_guide/testing.html)
inspec           | [testing inspec-gcp](https://github.com/inspec/inspec-gcp/#test-inspec-gcp-resources)
terraform        | [`google` provider testing guide](https://github.com/hashicorp/terraform-provider-google/blob/master/.github/CONTRIBUTING.md#tests)
terraform (beta) | [`google-beta` provider testing guide](https://github.com/hashicorp/terraform-provider-google-beta/blob/master/.github/CONTRIBUTING.md#tests)

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

## Compiling MMv1 + tpgtools

We are currently developing a new generation tool for Terraform called tpgtools.
This relies on a [declarative client library](https://github.com/GoogleCloudPlatform/declarative-resource-client-library)
that handles the actuation of GCP resources. We plan to gradually move resources
from being generated by the existing Magic Modules Ruby code (mmv1) to using
tpgtools. While we move resources over there will be a period of time when
both generators are in use. To assist with generation we have a series of `make`
targets that will run the compilers in tandem to generate the Terraform provider.

Sample Usage to compile at beta:
`make OUTPUT_PATH=/path/to/terraform-provider-google-beta VERSION=beta`

Target single product:
`make OUTPUT_PATH=/path/to/terraform-provider-google VERSION=ga PRODUCT=compute`

Target single resource
`make OUTPUT_PATH=/path/to/terraform-provider-google VERSION=ga PRODUCT=compute RESOURCE=image`

For more advanced usage of mmv1 compiler flags, please execute the compiler directly
from within the mmv1 directory.
