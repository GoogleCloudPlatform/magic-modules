---
title: Upgrade provider dependencies
weight: 65
---

## Before you begin

[Set up your development environment]({{< ref "/get-started/generate-providers/" >}}).
You should have the git repositories `magic-modules`,
[`terraform-provider-google`](http://github.com/hashicorp/terraform-provider-google), and
[`terraform-provider-google-beta`](http://github.com/hashicorp/terraform-provider-google-beta)
cloned onto your filesystem.

## Upgrade provider dependencies

Follow this procedure to update package dependencies for the Google Terrafrom
providers. This is typically done if a new feature has been introduced to
`google.golang.org/api`, if there is some security update, or other useful
or necessary functionality in a dependency we would like to take advantage of.

1. Change directories to `terraform-provider-google-beta`: `cd terraform-provider-google-beta`.
2. Upgrade your desired package.
   1. To update the `api` package, run `go get -u google.golang.org/api`.
   2. To update all pacakges, run `go get -u ./..`.
   3. [More Info on Go Package Management](https://golang.cafe/blog/how-to-upgrade-golang-dependencies.html)
3. Validate the new package specification.
   1. Run `make lint` to check the configuration.
   2. If your receive an error saying `go get`, repeat steps 2 and 3 until `make lint` passes.
4. Install the new configuration to `magic-modules`
   1. Copy `go.mod`
      1. Copy the contents of `terraform-provider-google-beta/go.mod` to `magic-modules/mmv1/third_party/terraform/go.mod.erb`
      2. `go.mod.erb` should start with `module github.com/hashicorp/terraform-provider-google`.
   2. Copy `go.sum`
      1. Copy the contents of `terraform-provider-google-beta/go.sum` to `magic-modules/mmv1/third_party/terraform/go.sum`.
5. Validatate Generated Modules build
   1. Reset each provider.
      1. Run `git clean -f -d && git reset --hard` in `terraform-provider-google`
      2. Run `git clean -f -d && git reset --hard` in `terraform-provider-google-beta`
   2. Validate Build
      1. Run `make lint` in `terraform-provider-google`
      2. Run `make lint` in `terraform-provider-google-beta`
6. Follow [Run tests]({{< ref "/develop/run-tests" >}} to further validate the new configuration.
7. Submit the new configuration by [creating a pull request](https://docs.github.com/en/get-started/quickstart/github-flow#create-a-pull-request).

## What's next?

- [Run tests]({{< ref "/develop/run-tests" >}}
- [Create a pull request](https://docs.github.com/en/get-started/quickstart/github-flow#create-a-pull-request)
