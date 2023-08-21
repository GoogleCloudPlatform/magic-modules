# Upgrading Dependencies

## Procedure

1. Please have an instance of `magic-modules` available. It will be used later in this procedure.
2. Checkout & `cd` to [`terraform-provider-google-beta`](http://github.com/hashicorp/terraform-provider-google-beta)
3. Upgrade your desired package. EX: `go get -u google.golang.org/api` You may need to upgrade all packages at once with `go get -u ./..`. [More Info](https://golang.cafe/blog/how-to-upgrade-golang-dependencies.html).
4. Run `make lint`. If there is an error saying `go get` should be run to fill in missing dependencies, please do so. Repeat until `make lint` passes.
5. Copy `go.mod` to `magic-modules` at `mmv1/third_party/terraform/go.mod.erb`. Do not include the `module github.com/hashicorp/terraform-provider-google-beta` beginning line that should stay as `module github.com/hashicorp/terraform-provider-google`.
6. Copy `go.sum` to `magic-modules` at `mmv1/third_party/terraform/go.sum`.
7. `git clean -f -d && git reset --hard` in beta provider
8. `git clean -f -d && git reset --hard` in ga provider
9. Generate the GA Provider. `make lint` should pass.
10. Generate the Beta Provider. `make lint` should pass.
11. Next Unit and Acceptance test should pass. The easiest way to do so is to
    create a draft PR and run the build.
12. Send PR to update dependencies. This can be done as a seperate PR, or as
    part of another PR adding functionality.
