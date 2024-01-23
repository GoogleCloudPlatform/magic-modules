---
majorVersion: "5.0.0"
upgradeGuide: "version_5_upgrade.html.markdown"
title: "Make a breaking change"
summary: "Guidance on making a breaking changes"
weight: 20
aliases:
- /develop/make-a-breaking-change
---

# Make a breaking change

A "breaking change" is any change that requires an end user to modify any
previously-valid configuration after a provider upgrade. For more information,
see [Types of breaking changes]({{< ref "/develop/breaking-changes" >}}).

The `google` and `google-beta` providers are both considered "stable surfaces"
for the purpose of releases, which means that neither provider allows breaking
changes except during major releases, which are typically yearly.

Terraform users rely on the stability of Terraform providers (including the GCP
provider and other major providers.) Even as part of a major release, breaking
changes that are overly broad and/or have little benefit to users can cause
deeply negative reactions and significantly delay customers upgrading to the
new major version.

Breaking changes may cause significant churn for users by forcing them to
update their configurations. It also causes churn in tooling built on top of
the providers, such as:

* Terraform modules that use `google` or `google-beta` resources
* Policy tools like [`gcloud terraform vet`](https://cloud.google.com/docs/terraform/policy-validation/quickstart)
  * There may also be churn in customer policies
* [Config Connector](https://cloud.google.com/config-connector/docs/overview)
* [Pulumi GCP Classic](https://www.pulumi.com/registry/packages/gcp/)

This page covers the general process to make a breaking change. It does not
include exact, comprehensive details on how to make every potential breaking
change. Breaking changes are complicated; the exact process and implementation
may vary drastically depending on the implementation of the impacted resource
or field and the change being made.

## In minor releases

If a breaking change fixes a bug that impacts **all configurations** that
include a field or resource, it is generally allowed in a minor release. For
example:

* Removing update support from a field if that field is not actually updatable
  in the API.
* Marking a field required if omitting the field always causes an API error.

The following types of changes can be made if the default behavior stays the
same and new behavior can be enabled with a flag:

* Major resource-level or field-level behavioural changes

## In the {{% param "majorVersion" %}} major release

The general process for contributing a breaking change to the
`{{% param "majorVersion" %}}` major release is:

1. Make the `main` branch forwards-compatible with the major release
2. Add deprecations and warnings to the `main` branch of `magic-modules`
3. Add upgrade guide entries to the `main` branch of `magic-modules`
4. Make the breaking change on `FEATURE-BRANCH-major-release-{{% param "majorVersion" %}}`

These are covered in more detail in the following sections. The upgrade guide
and the actual breaking change will be merged only after both are completed.

### Make the `main` branch forwards-compatible with the major release

What forwards-compatibility means will vary depending on the breaking change. For example:

* If a required field is being removed, make the field optional
  on the `main` branch.
* If a field is being renamed, the new field must be added to the `main` branch


### Add deprecations and warnings to the `main` branch of `magic-modules`

Deprecations and warnings must be actionable at the time that they are added
to the `main` branch, and they must be added prior to the `{{% param "majorVersion" %}}`
major release. Every deprecation or warning should be surfaced to users of the
provider at runtime as well as in documentation.

#### Field deprecation (due to removal or rename)

{{< tabs "Field deprecations" >}}
{{< tab "MMv1" >}}
Set `deprecation_message` on the field. For example:

```yaml
- !ruby/object:Api::Type::String
  name: 'apiFieldName'
  description: |
    MULTILINE_FIELD_DESCRIPTION
  deprecation_message: "`api_field_name` is deprecated and will be removed in a future major release. Use `other_field_name` instead."
```

Replace the second sentence with an appropriate short description of the replacement path and/or the reason for
deprecation.

The deprecation message will automatically show up in the resource documentation.
{{< /tab >}}
{{< tab "Handwritten" >}}
1. Set `Deprecated` on the field. For example:

   ```go
   "api_field_name": {
      Type:       schema.String,
      Deprecated: "`api_field_name` is deprecated and will be removed in a future major release. Use `other_field_name` instead.",
      ...
   }
   ```
   Replace the second sentence with an appropriate short description of the replacement path and/or the reason for
   deprecation.
2. Update the [documentation for the field]({{< ref "/develop/resource#add-documentation" >}}) to include the deprecation notice. For example:

   ```markdown
   * `api_field_name` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html), Deprecated) FIELD_DESCRIPTION. `api_field_name` is deprecated and will be removed in a future major release. Use `other_field_name` instead.
   ```
{{< /tab >}}
{{< /tabs >}}

#### Resource deprecation (due to removal or rename)

{{< tabs "Resource deprecations" >}}
{{< tab "MMv1" >}}
Set `deprecation_message` on the resource. For example:

```yaml
deprecation_message: >-
  `google_RESOURCE_NAME` is deprecated and will be removed in a future major release.
  Use `google_OTHER_RESOURCE_NAME` instead.
```

Replace RESOURCE_NAME with the name of the resource (excluding the `google_` prefix). Replace the
second sentence with an appropriate short description of the replacement path and/or the reason for
deprecation.

The deprecation message will automatically show up in the resource documentation.
{{< /tab >}}
{{< tab "Handwritten" >}}
1. Set `DeprecationMessage` on the field. For example:

   ```go
   return &schema.Resource{
      ...
      DeprecationMessage: "`google_RESOURCE_NAME` is deprecated and will be removed in a future " +
                          "major release. Use `google_OTHER_RESOURCE_NAME` instead.",
      ...
   }
   ```

   Replace RESOURCE_NAME with the name of the resource (excluding the `google_` prefix). Replace the
   second sentence with an appropriate short description of the replacement path and/or the reason for
   deprecation.
2. Add a warning to the resource documentation stating that the resource is deprecated. For example:
   ```markdown
   ~> **Warning:** `google_RESOURCE_NAME` is deprecated and will be removed in a future
   major release. Use `google_OTHER_RESOURCE_NAME` instead.
   ```
{{< /tab >}}
{{< /tabs >}}

#### Other breaking changes

Other breaking changes should be called out in the docs for the impacted field
or resource. It is also great to log warnings at runtime if possible.

### Add upgrade guide entries to the `main` branch of `magic-modules`

Upgrade guide entries should be added to
[{{< param upgradeGuide >}}](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/website/docs/guides/{{< param upgradeGuide >}}).
Entries should focus on the changes that users need to make when upgrading
to `{{% param "majorVersion" %}}`, rather than how to write configurations
after upgrading.

See [Terraform Google Provider 4.0.0 Upgrade Guide](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/version_4_upgrade)
and other upgrade guides for examples.

The upgrade guide and the actual breaking change will be merged only after both are completed.

### Make the breaking change on `FEATURE-BRANCH-major-release-{{% param "majorVersion" %}}`

When working on your breaking change, make sure that your base branch
is `FEATURE-BRANCH-major-release-{{% param "majorVersion" %}}`. This
means that you will follow the standard
[contribution process]({{< ref "/get-started/contribution-process" >}})
with the following changes:

1. Before you start, check out and sync your local `magic-modules` and provider
   repositories with the upstream major release branches.
   ```bash
   cd ~/magic-modules
   git checkout FEATURE-BRANCH-major-release-{{% param "majorVersion" %}}
   git pull --ff-only origin FEATURE-BRANCH-major-release-{{% param "majorVersion" %}}
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
   git checkout FEATURE-BRANCH-major-release-{{% param "majorVersion" %}}
   git pull --ff-only origin FEATURE-BRANCH-major-release-{{% param "majorVersion" %}}
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
   git checkout FEATURE-BRANCH-major-release-{{% param "majorVersion" %}}
   git pull --ff-only origin FEATURE-BRANCH-major-release-{{% param "majorVersion" %}}
   ```
1. Make sure that any deprecation notices and warnings that you added in previous sections
   are present on the major release branch. Changes to the `main` branch will be
   merged into the major release branch every Monday.
1. Make the breaking change.
1. Remove any deprecation notices and warnings (including in documentation) not already removed by the breaking change.
1. When you create your pull request,
   [change the base branch](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/changing-the-base-branch-of-a-pull-request)
   to `FEATURE-BRANCH-major-release-{{% param "majorVersion" %}}`
1. To resolve merge conflicts with `git rebase` or `git merge`, use `FEATURE-BRANCH-major-release-{{% param "majorVersion" %}}` instead of `main`.

The upgrade guide and the actual breaking change will be merged only after both are completed.

## What's next?

- [Run tests]({{< ref "/develop/test/run-tests.md" >}})
