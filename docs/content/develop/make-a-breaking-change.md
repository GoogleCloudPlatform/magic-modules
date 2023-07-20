---
title: "Make a breaking change"
summary: "Guidance on making a breaking change during a major release"
weight: 95
---

# How to make a breaking change

{{< hint info >}}
**Note:** This page covers the general requirements to make a breaking change
and may not include exact, comprehensive details on making your change. Breaking
changes are generally complicated and resource or field specific traits can
drastically change how they need to be made.
{{< /hint >}}

## When to (or not to) make a breaking change

Within the Terraform ecosystem, there's a strong expectation of stability
**including across major versions**. Ecosystem-wide, the community has had
deeply negative reactions to provider changes with overly-broad impact, and
too-large changes can delay customers upgrading to adopt new major versions.

In general, we recommend avoiding breaking changes where possible. The
provider's schema is an API surface relied on by many GCP customers, and users
have responded negatively to instability in our surface and those of other
providers.

While the cost to make a change in the provider is relatively cheap, it has
knock-on effects through the whole ecosystem:

* Modules need to update to adapt to the breaking changes
* Customers need to update their provider version (and module versions)
* Terraform assistive tools like [`gcloud terraform vet`](https://cloud.google.com/docs/terraform/policy-validation/quickstart)
  and 3P tools with provider dependencies like
  [Config Connector](https://cloud.google.com/config-connector/docs/overview)
  and [Pulumi GCP Classic](https://www.pulumi.com/registry/packages/gcp/) need
  to be updated

When breaking changes are made, they must be made within a major release
(typically yearly) and meet the stringent ecosystem requirements for a breaking
change, detailed below.

## What counts as a breaking change?

While the general versioning model for providers is captured
in https://developer.hashicorp.com/terraform/plugin/best-practices/versioning,
that documentation serves as a baseline for new provider development and not
guidelines for mature providers like the major cloud providers.

The standard our joint Google & HashiCorp team has developed is that any change
that requires an end user to modify any previously-valid configuration after a
provider upgrade is a breaking change, and must happen in a major release. In
this context "valid" means that a configuration is syntactically correct
(passes `terraform validate`) and runs in `terraform apply` without returning an
error. This means that minor corrections are possible such as marking immutable
fields as immutable if they previously weren't.

There is a single exception to this rule- if a user has managed a resource
out-of-band from Terraform to enable a new field with a non-default value,
adding support for the field without handling that case is permissible.
Terraform will generate a plan that reconciles their state with their live
configuration.

For example, the following are (a non-exhaustive list of) breaking changes:

* Changing the output format of a field (i.e. from an integer to a string, or
  the pattern of a structured value)
* Adding a new required field or changing an existing non-required field to
  required
* Removing a field
* Major behaviour changes to a resource (if configurable, if they are done by
  default)

Meanwhile, the following are allowed in a minor version:

* Adding a new field
* Adding update support for a field
* Removing update support from a field that returned an error in **all cases** a
  user attempted to update it
* Marking a field required if *all configurations** that did not specify the
  field returned an error
* Major behavioural changes guarded by a flag where the **previous** behaviour
  is the default

## Making a breaking change

{{< hint info >}}
**Note:** This section of the document refers to several release versions of the
provider at once. Breaking changes are made in major releases such as 1.0.0,
4.0.0, or 4.0.0 (referred to as the "major release" or `N` here)  while typical
provider releases are minor versions, most notably the last minor release of the
previous major release series such as `2.20.3` for `3.0.0` or `3.90.1`
for `4.0.0` (referred to as the "last minor release" or `N-1.X` here).
{{</hint >}}

The Terraform provider ecosystem follows the standard that deprecations or
warnings must be resolvable by end users in the last minor release of the prior
release series before their removal or change in a major release. Additionally,
deprecation warnings must be actionable- at the time a deprecation is posted, a
user must be able to remove the field.

### Contributing to the next major release (`5.0.0`)

For the `5.0.0` major release, the major release branch that you'll contribute to is [`FEATURE-BRANCH-major-release-5.0.0`](https://github.com/GoogleCloudPlatform/magic-modules/tree/FEATURE-BRANCH-major-release-5.0.0).
All breaking changes targeting `5.0.0` must be committed to this branch.

A downstream branch with the same name `FEATURE-BRANCH-major-release-5.0.0` will be used to track the generated `5.0.0` changes in both [`google`](https://github.com/hashicorp/terraform-provider-google/tree/FEATURE-BRANCH-major-release-5.0.0) and [`google-beta`](https://github.com/hashicorp/terraform-provider-google-beta/tree/FEATURE-BRANCH-major-release-5.0.0) provider repos

The process of contributing to the major release `5.0.0` should follow most of the [General contributing steps]({{< ref "/get-started/contributing" >}}), with the following exceptions

1. Use `FEATURE-BRANCH-major-release-5.0.0` branch instead of the `main` branch as the base branch when you
   * checkout your working branch where you make your code changes.
   * sync your working branch using `git rebase` or `git merge`.
   * create a pull request in the magic-modules repo.
2. Make sure that you checkout to the `FEATURE-BRANCH-major-release-5.0.0` branch in your downstream `google` and `google-beta` repos before generating the providers locally.

### Renaming a field

The most common type of breaking change is a field rename, and most guidance is
tuned around that. To perform one, a provider contributor must:

1. Add support for the new field on or before the last minor release of the
   preceding release series (i.e. version `N-1.X.0`) by contributing to
   the `main` branch in the magic-modules repo
2. Mark the old field deprecated on or before the last minor release of the
   preceding release series (i.e. version `N-1.X.0`) by contributing to
   the `main` branch in the magic-modules repo
3. Write an upgrade guide entry in the major release's upgrade guide (such
   as [this one](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/version_5_upgrade)
   for `5.0.0`) by contributing to the `main` branch in the magic-modules repo
4. Remove the old field from the major release (i.e. version `N`) in the major
   release branch in the magic-modules repo

For example, if `google_foobar` has a field `baz` in version `3.80.0` that is
being replaced by `qux` in `4.0.0`, `qux` must be added on or before `3.90.1`,
and `baz` must be deprecated within the same window (either at the same time
as `qux` is added, or afterwards (if added earlier than `3.90.1`), but not
before).

Example  (`google_storage_bucket.bucket_policy_only` -> `google_storage_bucket.uniform_bucket_level_access`)
* https://github.com/GoogleCloudPlatform/magic-modules/pull/3916 introduced the new field and deprecated the original one in `3.38.0`
* https://github.com/GoogleCloudPlatform/magic-modules/pull/5340 added the upgrade guide entry and removed the field in `4.0.0`
* Note: In 4.0.0 the upgrade guide was split between main and the major release branch. Going forward, those changes will be made against main exclusively.

### Removing a field

Removing a field is similar to renaming one, except that a new field doesn't
need to be introduced:

1. Mark the field deprecated on or before the last minor release of the
   preceding release series (i.e. version `N-1.X.0`) by contributing to
   the `main` branch in the magic-modules repo
2. Write an upgrade guide entry in the major release's upgrade guide (such
   as [this one](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/version_5_upgrade)
   for `5.0.0`) by contributing to the `main` branch in the magic-modules repo
3. Remove the field from the major release (i.e. version `N`) in the major
   release branch in the magic-modules repo

Example (`google_container_cluster.instance_group_urls`)
* https://github.com/GoogleCloudPlatform/magic-modules/pull/5261 deprecated `instance_group_urls` and filled out the upgrade guide
* https://github.com/GoogleCloudPlatform/magic-modules/pull/5378 removed the field

### Marking optional fields required

There is no way to message a change to users in advance at the moment, other
than through the upgrade guide.

1. Write an upgrade guide entry in the major release's upgrade guide (such
   as [this one](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/version_5_upgrade)
   for `5.0.0`) by contributing to the `main` branch in the magic-modules repo
2. Mark the field as required in the major release (i.e. version `N`) in the major
   release branch in the magic-modules repo

Example (`google_app_engine_standard_app_version.entrypoint`)
* https://github.com/GoogleCloudPlatform/magic-modules/pull/5318 marked the field required and added an upgrade guide entry
* Note: In 4.0.0 the upgrade guide was split between main and the major release branch. Going forward, those changes will be made against main exclusively.

### Changing default values

Default values in Terraform are used to replace null values in configuration at
plan/apply time and **do not** respect previously-configured values by the user.
These changes are often undesirable, as their impact is extremely broad.

When a default is changed, every user that has not specified an explicit value in their configuration will see Terraform propose changing the value of the field **including** if the change will destroy and recreate the resource due to changing an immutable value. Default changes in the provider are comparable in impact to default changes in an API, and modifying examples and modules may achieve the intended effect with a smaller blast radius.

There is no way to message a change to users in advance at the moment, other
than through the upgrade guide.

1. Write an upgrade guide entry in the major release's upgrade guide (such
   as [this one](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/version_5_upgrade)
   for `5.0.0`) by contributing to the `main` branch in the magic-modules repo
2. Mark the new default in the major release (i.e. version `N`) in the major
   release branch in the magic-modules repo

Example (`google_container_cluster.enable_shielded_nodes`)
* https://github.com/GoogleCloudPlatform/magic-modules/pull/5263 changed the default value for the field and added an upgrade guide entry
* Note: In 4.0.0 the upgrade guide was split between main and the major release branch. Going forward, those changes will be made against main exclusively.

## References

* Terraform (Provider) Plugin Development
    * [Versioning and Changelog](https://developer.hashicorp.com/terraform/plugin/best-practices/versioning)
    * [SDKv2 Deprecations, Removals, and Renames](https://developer.hashicorp.com/terraform/plugin/sdkv2/best-practices/deprecations)
