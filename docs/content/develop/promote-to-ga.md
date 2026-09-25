---
title: "Promote to GA"
weight: 80
---

# Promote from beta to GA

This document describes how to promote an existing resource or field that uses MMv1 and/or handwritten code from the `google-beta` provider to the `google` (also known as "GA") provider.

Handwritten code (including `custom_code`) commonly uses "version guards" in the form of `{{- if ne $.TargetVersionName "ga" }}...{{- end }}` to wrap code that is beta-specific, which need to be removed during promotion.

For more information about types of resources and the generation process overall, see [How Magic Modules works]({{< ref "/" >}}).

## Before you begin

1. Complete the steps in [Set up your development environment]({{< ref "/develop/set-up-dev-environment" >}}) to set up your environment and your Google Cloud project.
1. Ensure that your `magic-modules`, `terraform-provider-google`, and `terraform-provider-google-beta` repositories are up to date.
    ```bash
    cd ~/magic-modules
    git checkout main && git clean -f . && git checkout -- . && git pull
    cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
    git checkout main && git clean -f . && git checkout -- . && git pull
    cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
    git checkout main && git clean -f . && git checkout -- . && git pull
    ```

## Promote fields and resources

{{% tabs "resources" %}}
{{< tab "MMv1" >}}
1. In `product.yaml` (located in `mmv1/products/<product_name>/product.yaml`), ensure that the `versions` list includes the `ga` version and its corresponding `base_url`. If the product was previously beta-only, this entry will be missing and must be added before any resources can be promoted. Keep the existing `beta` entry: if it is removed, the `google-beta` provider silently falls back to the GA `base_url`.
2. Remove `min_version: 'beta'` from the resource's or field's configuration in `ResourceName.yaml`.
3. If necessary, remove version guards from resource-level `custom_code`.
4. Add `min_version: 'beta'` on any fields or subfields that should not be promoted.
5. If necessary, add `{{- if ne $.TargetVersionName "ga" }}...{{- end }} ` version guards to resource-level `custom_code` that should not be promoted.
6. Update API versions hardcoded outside `product.yaml`, for example in a resource-level `base_url`/`self_link`, `custom_code` URLs (use `transport_tpg.BaseUrl`), or `references` links.
{{< /tab >}}
{{< tab "Handwritten" >}}
1. Remove version guards from the resource's implementation for any functionality being promoted. Be sure to check:
   - The overall resource (if the entire resource was beta-only)
   - The resource schema
   - For top-level fields, the resource's `Create`, `Update`, and `Read` methods
   - For other fields, expanders and flatteners
   - Any other resource-specific code
   - Related files without the resource's name: handwritten data sources, sweepers, `bootstrap_test_utils.go`, and schema helpers shared with sibling resources (e.g. `google_compute_instance` and `google_compute_instance_template`)
2. Add `{{- if ne $.TargetVersionName "ga" }}...{{- end }}` version guards to any parts of the resource or field implementation that should not be promoted. Be sure to check:
   - The resource schema
   - For top-level fields, the resource's `Create`, `Update`, and `Read` methods
   - For other fields, expanders and flatteners
   - Any other resource-specific code
3. If a `.go.tmpl` file no longer contains any version guards, rename it to `.go` and format it with `gofmt`.
{{< /tab >}}
{{% /tabs %}}

## Promote tests

1. Remove `min_version: beta` from any samples in a `ResourceName.yaml` which only test fields and resources that are present in the `google` provider. This includes fields on other resources in the configuration: a test with a beta-only dependency can't run in `google`.
2. Remove version guards from any handwritten code related to fields and resources that are present in the `google` provider.
3. Delete `provider = google-beta` from any test configurations (from MMv1 samples or handwritten) which have been promoted. Don't replace it with `provider = google`.
4. Replace `ProtoV5ProviderBetaFactories` with `ProtoV5ProviderFactories` in all promoted handwritten tests.
5. Ensure that there is at least one test that will run for the `google` provider that covers any promoted fields and resources.
6. Run the promoted tests against the `google` provider (see [Run tests]({{< ref "/test/run-tests" >}})) and include the results in your pull request. Presubmit VCR tests only run against `google-beta`.

## Promote documentation

For handwritten resources, modify the documentation as appropriate for your change:

1. If the entire resource has been promoted to `google`, remove the beta warning at the top of the documentation.
2. Remove the `Beta` annotation for any fields that have been promoted.
3. Add `Beta` as an annotation on any fields or subfields that remained beta-only. For example:

   ```markdown
   * `FIELD_NAME` - (Optional, [Beta](../guides/provider_versions.html.markdown)) FIELD_DESCRIPTION
   ```

   Replace `FIELD_NAME` and `FIELD_DESCRIPTION` with the field's name and description.

## What's next?

- [Test your changes]({{< ref "/test/run-tests" >}})
