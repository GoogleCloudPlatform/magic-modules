---
title: Set up your environment
weight: 10
---

# Set up your environment

## Cloning Terraform providers

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

## Setting up a container-based environment

*NOTE* this approach is in beta and still collecting feedback. Please file an issue if you encounter challenges, and try pulling the latest container (see command below) first to see if any recent changes may fix you.

You do not need to run these instructions if you are setting up your environment manually.

For ease of contribution, we provide containers with the required dependencies for building magic-modules, as well as the option to build them yourself.

You can work with containers with either [Podman](https://podman.io/) or [Docker](https://docker.io/).

[scripts/make-in-container.sh](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/scripts/make-in-container.sh) includes all the
instructions to build and run containers by hand. Refer to the script for
individual steps, but for most users you only have to run the script directly, as a replacement for make.

Here is an example of how to build Terraform (after [cloning the provider](#cloning-terraform-providers)):

```shell
./scripts/make-in-container.sh \
  terraform VERSION=ga \
  OUTPUT_PATH="$GOPATH/src/github.com/hashicorp/terraform-provider-google"
```

Generally, you can replace any reference to `make` in this guide with `scripts/make-in-container.sh`.

## Preparing your environment manually

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

## Test your setup

Try [generating the providers](/magic-modules/docs/getting-started/generate-providers/). If your environment is set up correctly, it should succeed with no errors!