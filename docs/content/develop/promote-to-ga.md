---
title: "Promote to GA"
weight: 50
---

# Promote from beta to GA

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

{{< tabs "resources" >}}
{{< tab "MMv1" >}}
1. Remove `min_version: beta` from the resource's or field's configuration in `ResourceName.yaml`.
2. If necessary, remove version guards from resource-level `custom_code`.
3. Add `min_version: beta` on any fields or subfields that should not be promoted.
4. If necessary, add `<% unless version == 'ga' -%>...<% end -%>` version guards to resource-level `custom_code` that should not be promoted.
{{< /tab >}}
{{< tab "Handwritten" >}}
1. Remove version guards from the resource's implementation for any functionality being promoted. Be sure to check:
   - The overall resource (if the entire resource was beta-only)
   - The resource schema
   - For top-level fields, the resource's `Create`, `Update`, and `Read` methods
   - For other fields, expanders and flatteners
   - Any other resource-specific code
2. Add `<% unless version == 'ga' -%>...<% end -%>` version guards to any parts of the resource or field implementation that should not be promoted. Be sure to check:
   - The resource schema
   - For top-level fields, the resource's `Create`, `Update`, and `Read` methods
   - For other fields, expanders and flatteners
   - Any other resource-specific code
{{< /tab >}}
{{< /tabs >}}

## Promote tests

1. Remove `min_version: beta` from any examples in a `ResourceName.yaml` which only test fields and resources that are present in the `google` provider.
2. Remove version guards from any handwritten code related to fields and resources that are present in the `google` provider.
3. Remove `provider = google-beta` from any test configurations (from MMv1 `examples` or handwritten) which have been promoted.
4. Ensure that there is at least one test that will run for the `google` provider that covers any promoted fields and resources.

## Promote documentation

For handwritten resources, modify the documentation as appropriate for your change:

1. If the entire resource has been promoted to `google`, remove the beta warning at the top of the documentation.
2. Remove the `Beta` annotation for any fields that have been promoted.
3. Add `Beta` as an annotation on any fields or subfields that remained beta-only. For example:

   ```markdown
   * `FIELD_NAME` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) FIELD_DESCRIPTION
   ```

   Replace `FIELD_NAME` and `FIELD_DESCRIPTION` with the field's name and description.

## What's next?

- [Test your changes]({{< ref "/develop/test/run-tests.md" >}})
