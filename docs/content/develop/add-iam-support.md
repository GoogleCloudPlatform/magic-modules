---
title: "Add IAM support"
weight: 40
---

# Add IAM support

This page covers how to add IAM resources in Terraform if they are supported by a particular API resource (indicated by
`setIamPolicy` and `getIamPolicy` methods in the API documentation for the resource).

For more information about types of resources and the generation process overall, see [How Magic Modules works]({{< ref "/" >}}).

## Before you begin

1. Complete the steps in [Set up your development environment]({{< ref "/develop/set-up-dev-environment" >}}) to set up your environment and your Google Cloud project.
1. Ensure that your `magic-modules`, `terraform-provider-google`, and `terraform-provider-google-beta` repositories are up to date.
   ```
   cd ~/magic-modules
   git checkout main && git clean -f . && git checkout -- . && git pull
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
   git checkout main && git clean -f . && git checkout -- . && git pull
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
   git checkout main && git clean -f . && git checkout -- . && git pull
   ```

## Add IAM support

{{< tabs "IAM" >}}
{{< tab "MMv1" >}}
IAM support for MMv1-generated resources is configured within the `ResourceName.yaml` file, and will create the `google_product_resource_iam_policy`, `google_product_resource_iam_binding`, `google_product_resource_iam_member` resource, website, and test files for that resource target when an `iam_policy` block is present.

1. Add the following top-level block to `ResourceName.yaml` directly above `parameters`.

```yaml
iam_policy:
  # Name of the field on the terraform IAM resources which references
  # the parent resource. Update to match the parent resource's name.
  parent_resource_attribute: 'resource_name'
  # Character preceding setIamPolicy in the full URL for the API method.
  # Usually `:`
  method_name_separator: ':'
  # HTTP method for getIamPolicy. Usually 'POST'.
  fetch_iam_policy_verb: 'POST'
  # Overrides the HTTP method for setIamPolicy. Default: 'POST'
  # set_iam_policy_verb: 'POST'

  # Must match the parent resource's `import_format` (or `self_link` if
  # `import_format` is unset), but with the `parent_resource_attribute`
  # value substituted for the final field.
  import_format:
    - 'projects/{{project}}/locations/{{location}}/resourcenames/{{resource_name}}'

  # If IAM conditions are supported, set this attribute to indicate how the
  # conditions should be passed to the API. Allowed values: 'QUERY_PARAM',
  # 'REQUEST_BODY', 'QUERY_PARAM_NESTED'. Note: 'QUERY_PARAM_NESTED' should
  # only be used if the query param field contains a `.`
  # iam_conditions_request_type: 'REQUEST_BODY'

  # Marks IAM support as beta-only
  # min_version: beta
```

2. Modify the template as needed to match the API resource's documented behavior. These are the most commonly-used fields. For a comprehensive reference, see [MMv1 resource reference: `iam_policy` ↗]({{<ref "/reference/resource#iam_policy" >}}).
3. Delete all remaining comments in the IAM configuration (including attribute descriptions) that were copied from the above template.
{{< /tab >}}
{{< tab "Handwritten" >}}
> **Warning:** IAM support for handwritten resources should be implemented using MMv1. New handwritten IAM resources will only be accepted if they cannot be implemented using MMv1.

### Add support in MMv1

1. Follow the MMv1 directions in [Add a resource]({{<ref "/develop/add-resource" >}}) to create a skeleton ResourceName.yaml file for the handwritten resource, but set only the following top-level fields:
   - `name`
   - `description` (required but unused)
   - `base_url` (set to URL of IAM parent resource)
   - `self_link` (set to same value as `base_url`)
   - `id_format` (set to same value as `base_url`)
   - `import_format` (including `base_url` value)
   - `exclude_resource` (set to `true`)
   - `properties`
2. Follow the MMv1 directions in [Add fields]({{<ref "#add-fields" >}}) to add only the fields used by base_url.
3. Follow the MMv1 directions in this section to add IAM support.

### Convert to handwritten (not usually necessary)

1. [Generate the beta provider]({{< ref "/develop/generate-providers" >}})
2. From the beta provider, copy the files generated for the IAM resources to the following locations:
   - Resource: Copy to the appropriate service folder inside [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services)
   - Documentation: [`magic-modules/mmv1/third_party/terraform/website/docs/r`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/website/docs/r)
   - Tests: In the appropriate service folder inside [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services)
3. Modify the Go code as needed.
   - Replace all occurrences of `github.com/hashicorp/terraform-provider-google-beta/google-beta` with `github.com/hashicorp/terraform-provider-google/google`
   - Remove the comments at the top of the file.
   - If any of the added Go code is beta-only:
     - Change the file suffix to `.go.tmpl`
     - Wrap each beta-only code block (including any imports) in a separate version guard: `{{- if ne $.TargetVersionName "ga" -}}...{{- else }}...{{- end }}`
4. Register the binding, member, and policy resources `handwrittenIAMResources` in [`magic-modules/mmv1/third_party/terraform/provider/provider_mmv1_resources.go.tmpl`](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/provider/provider_mmv1_resources.go.tmpl)
   - Add a version guard for any beta-only resources.
{{< /tab >}}
{{< /tabs >}}

## Add documentation

{{< tabs "docs" >}}
{{< tab "MMv1" >}}
Documentation is autogenerated based on the resource and field configurations. To preview the documentation:

1. [Generate the providers]({{< ref "/develop/generate-providers" >}})
2. Copy and paste the generated documentation into the Hashicorp Registry's [Doc Preview Tool](https://registry.terraform.io/tools/doc-preview) to see how it is rendered.
{{< /tab >}}
{{< tab "Handwritten" >}}
### Add or modify documentation files

1. Open the resource documentation in [`magic-modules/third_party/terraform/website/docs/r/`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/website/docs/r) using an editor of your choice.
   - The name of the file is the name of the resource without a `google_` prefix. For example, for `google_compute_instance`, the file is called `compute_instance.html.markdown`
2. Modify the documentation as needed according to [Handwritten documentation style guide]({{< ref "/document/handwritten-docs-style-guide" >}}).
3. [Generate the providers]({{< ref "/develop/generate-providers" >}})
4. Copy and paste the generated documentation into the Hashicorp Registry's [Doc Preview Tool](https://registry.terraform.io/tools/doc-preview) to see how it is rendered.
{{< /tab >}}
{{< /tabs >}}

## What's next?

+ [Add documentation]({{< ref "/document/add-documentation" >}})
+ [Add custom resource code]({{< ref "/develop/custom-code" >}})
+ [Add tests]({{< ref "/test/test" >}})
+ [Run tests]({{< ref "/test/run-tests" >}})
