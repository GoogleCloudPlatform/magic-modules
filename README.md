[![Build Status](https://travis-ci.org/GoogleCloudPlatform/magic-modules.svg?branch=master)](https://travis-ci.org/GoogleCloudPlatform/magic-modules)


# Magic Modules

<img src="images/magic-modules.svg" width="300" align="right" />

## Overview

Magic Modules is a tool we use to autogenerate infrastructure-as-code tools for
Google Cloud Platform. GCP ["resources"](https://cloud.google.com/docs/overview/#gcp_resources)
are encoded in a shared data file, and that data is used to fill in
"Mad Libs"-style templates across each of the tools Magic Modules generates.

They include:

* Terraform
* Ansible
* InSpec

Not only is Magic Modules a force multiplier for our developers, Magic Modules
allows us to preemptively solve issues by encoding field-tested learnings about
the GCP API into each tool; issues solved for one tool will be fixed in each
other tool.

## Getting Started with Magic Modules

We've prepared an interactive tutorial that you can try out with Open in Cloud
Shell below; if you're getting set up on a local workstation, this guide serves
as a reference.

[![Open in Cloud Shell](http://gstatic.com/cloudssh/images/open-btn.svg)](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/magic-modules&tutorial=TUTORIAL.md)

### Requirements

To get started, you'll need:

* Ruby 2.5.0
  * You can use `rbenv` to manage your Ruby version(s)
* [`Bundler`](https://github.com/bundler/bundler)
  * This can be installed with `gem install bundler`

### Preparing Magic Modules / One-time setup

If you're developing Ansible or Inspec, we use submodules to manage the Magic
Modules generated outputs:

```
git submodule update --init
```

If you're generating the Terraform providers (`google` and `google-beta`),
you'll need to check out the repo(s) you're generating in your GOPATH. For
example:

```
git clone https://github.com/terraform-providers/terraform-provider-google.git $GOPATH/src/github.com/terraform-providers/terraform-provider-google
git clone https://github.com/terraform-providers/terraform-provider-google-beta.git $GOPATH/src/github.com/terraform-providers/terraform-provider-google-beta
```

Once you've prepared the target folders for the tools, run the following to
finish getting Magic Modules set up by installing the Ruby gems it needs to run:

```
bundle install
```


### Generating downstream tools

Before making any changes, you can compile the "downstream" tool you're working
on by running the following command. If Magic Modules has been installed
correctly, you'll get no errors when you run a command:

```
bundle exec compiler -a -v "ga" -e {{tool}} -o "{{output_folder}}"
```

Generally, you'll want to generate into the same output; here's a reference of
common commands

{{tool}}  | {{output_folder}}
----------|----------
terraform | $GOPATH/src/github.com/terraform-providers/terraform-provider-google
ansible   | build/ansible
inspec    | build/inspec

For example, to generate Terraform:

```
bundle exec compiler -a -v "ga" -e terraform -o "$GOPATH/src/github.com/terraform-providers/terraform-provider-google"
```

It's worth noting that Magic Modules will only generate new files when ran
locally. The "Magician"- the Magic Modules CI system- handles deletion of old
files when creating PRs.

#### Terraform's `google-beta` provider

Terraform is the only tool to handle Beta features right now; you can generate
`google-beta` by running the following, substitution `"beta"` for the version
and using the repository for the `google-beta` provider.

```
bundle exec compiler -a -v "beta" -e terraform -o "$GOPATH/src/github.com/terraform-providers/terraform-provider-google-beta"
```

### Making changes to resources

TODO

You can modify resources in the `api.yaml` of the GCP product it is a part of as
well as in the tool-specific yaml file such as `terraform.yaml`.

Most fields in .yaml files are documented in the [`api` package](https://github.com/GoogleCloudPlatform/magic-modules/tree/master/api).

#### Making changes to handwritten files

TODO

### Testing your changes

Once you've generated your changes for the tool, you can test them by running the
tool-specific tests as if you were submitting a PR against that tool.

You can run tests in the `{{output_folder}}` from above. See the following for
more details;

Tool             | Testing Guide
-----------------|--------------
ansible          | [instructions](https://docs.ansible.com/ansible/devel/dev_guide/testing.html)
inspec           | TODO(slevenick): Add this
terraform        | [`google` provider testing guide](https://github.com/terraform-providers/terraform-provider-google/blob/master/.github/CONTRIBUTING.md#tests)
terraform (beta) | [`google-beta` provider testing guide](https://github.com/terraform-providers/terraform-provider-google-beta/blob/master/.github/CONTRIBUTING.md#tests)

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

#### Contributor License Agreements

We have a few legal steps to accept most community changes. Please fill out
either the individual or corporate Contributor License Agreement (CLA).

  * If you are an individual writing original source code and you're sure you
    own the intellectual property, then you'll need to sign an [individual CLA](http://code.google.com/legal/individual-cla-v1.0.html).
  * If you work for a company that wants to allow you to contribute your work,
    then you'll need to sign a [corporate CLA](http://code.google.com/legal/corporate-cla-v1.0.html).

Follow either of the two links above to access the appropriate CLA and
instructions for how to sign and return it. Once we receive it, we'll be able to
review/accept your PR.
