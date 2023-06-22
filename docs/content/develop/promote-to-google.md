---
title: "Promote to `google`"
weight: 50
---

# Promote to the `google` provider

This document describes how to promote an existing resource or field that uses MMv1 and/or handwritten code from the `google-beta` provider to the `google` (also known as "GA") provider.

Handwritten code (including `custom_code`) commonly uses "version guards" in the form of `<% unless version == 'ga' -%>...<% end -%>` to wrap code that is beta-specific, which need to be removed during promotion.

For more information about types of resources and the generation process overall, see [How Magic Modules works]({{< ref "/get-started/how-magic-modules-works.md" >}}).

## Before you begin

1. Complete the [Generate the providers]({{< ref "/get-started/generate-providers" >}}) quickstart to set up your environment and your Google Cloud project.
2. Ensure that your `magic-modules`, `terraform-provider-google`, and `terraform-provider-google-beta` repositories are up to date.
   ```
   cd ~/magic-modules
   git checkout main && git clean -f . && git checkout -- . && git pull
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
   git checkout main && git clean -f . && git checkout -- . && git pull
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
   git checkout main && git clean -f . && git checkout -- . && git pull
   ```

## Promote fields and resources

1. For MMv1 fields or resources:
   - Remove `min_version: beta` from the resource's or field's configuration in `ResourceName.yaml`.
   - Add `min_version: beta` on any fields or subfields that should not be promoted.
   - Remove version guards from `custom_code` for the field or resource if necessary.
2. For handwritten resources:
   - Remove version guards from the field's schema definition, the `Create` and `Update` methods, and any other code used by the field or resource.
   - Add `<% unless version == 'ga' -%>...<% end -%>` version guards to any parts of the resource or field implementation that should not be promoted.

## Promote tests

1. Remove `min_version: beta` from any examples in `ResourceName.yaml` which only test fields and resources that are present in the `google` provider.
2. Remove version guards from any handwritten code related to fields and resources that are present in the `google` provider.
3. Remove `provider = google-beta` from any test configurations (from MMv1 `examples` or handwritten) which have been promoted.
4. Ensure that there is at least one test that will run for the `google` provider that covers any promoted fields and resources.

## Promote documentation

For handwritten resources, modify the documentation as appropriate for your change:

1. If the entire resource has been promoted to `google`, remove the beta warning at the top of the documentation.
2. Remove the `Beta` annotation for any fields that have been promoted.
3. Add `Beta` as an annotation on any fields or subfields that remained beta-only. For example:

   ```
   * `FIELD_NAME` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) FIELD_DESCRIPTION
   ```

   Replace `FIELD_NAME` and `FIELD_DESCRIPTION` with the field's name and description.

## Generate & Test your changes
