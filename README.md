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

