---
title: "Add MMv1 IAM resources"
summary: "For resources implemented through the MMv1 engine, the majority of configuration
for IAM support can be inferred based on the preexisting YAML specification file."
weight: 11
---

# Add MMv1 IAM resources

For resources implemented through the MMv1 engine, the majority of configuration
for IAM support can be inferred based on the preexisting YAML specification file.

To add support for IAM resources based on an existing resource, add an
`iam_policy` block to the resource's definition in `api.yaml`, such as the
following:

```yaml
    iam_policy: !ruby/object:Api::Resource::IamPolicy
      method_name_separator: ':'
      fetch_iam_policy_verb: :POST      
      parent_resource_attribute: 'registry'
      import_format: ["projects/{{project}}/locations/{{location}}/registries/{{name}}", "{{name}}"]         
```

The specification values can be determined based on a mixture of the resource
specification and the cloud.google.com `setIamPolicy`/`getIamPolicy` REST
documentation, such as
[this page](https://cloud.google.com/iot/docs/reference/cloudiot/rest/v1/projects.locations.registries/setIamPolicy)
for Cloud IOT Registries.

`parent_resource_attribute` - (Required) determines the field name of the parent
resource reference in the IAM resources. Generally, this should be the singular
form of the parent resource kind in snake case, i.e. `registries` -> `registry`
or `backendServices` -> `backend_service`.

`method_name_separator` - (Required) should be set to the character preceding
`setIamPolicy` in the "HTTP Request" section on the resource's `setIamPolicy`
page. This is almost always `:` for APIs other than Google Compute Engine (GCE),
MMv1's `compute` product.

`fetch_iam_policy_verb` - (Required) should be set to the HTTP verb listed in
the "HTTP Request" section on the resource's `getIamPolicy` page. This is
generally `POST` but is occasionally `GET`. Note: This is specified as a Ruby
symbol, prefixed with a `:`. For example, for `GET`, you would specify `:GET`.

`import_format` - (Optional) A list of templated strings used to determine the
Terraform import format. If the resource has a custom `import_format` or
`id_format` defined in `terraform.yaml`, this must be supplied.

  * If an `import_format` is set on the parent resource use that set of values exactly, substituting `parent_resource_attribute` for the field name of the **final** templated value.
  * If an `id_format` is set on the parent resource use that as the first entry (substituting the final templated value, as with `import_format`) and define a second format with **only** the templated values, `/`-separated. For example, `projects/{{project}}/locations/{{region}}/myResources/{{name}}` -> `["projects/{{project}}/locations/{{region}}/myResources/{{myResource}}", "{{project}}/{{region}}/{{myResource}}"]`. 
    * Optionally, you may provide a version of the shortened format that excludes entries called `{{project}}`, `{{region}}`, and `{{zone}}`. For example, given `{{project}}/{{region}}/{{myResource}}/{{entry}}`, `{{myResource}}/{{entry}}` is a valid format. When a user specifies this format, the provider's default values for `project`/`region`/`zone` will be used.

`allowed_iam_role` - (Optional) If the resource does not allow the
`roles/viewer` IAM role to be set, an alternate, valid role must be provided.

`iam_conditions_request_type` - (Optional) The method the IAM policy version is
set in `getIamPolicy`. If unset, IAM conditions are assumed to not be supported for the resource. One of `QUERY_PARAM`, `QUERY_PARAM_NESTED` or `REQUEST_BODY`. For resources where a query parameter is expected, `QUERY_PARAM` should be used if the key is `optionsRequestedPolicyVersion`, while `QUERY_PARAM_NESTED` should be used if it is `options.requestedPolicyVersion`.

`min_version` - (Optional) If the resource or IAM method is not generally
available, this should be set to `beta` or `alpha` as appropriate.

`set_iam_policy_verb` - (Optional, rare) Similar to `fetch_iam_policy_verb`, the
HTTP verb expected by `setIamPolicy`. Defaults to `:POST`, and should only be
specified if it differs (typically if `:PUT` is expected).

Several single-user settings are not documented on this page as they are not
expected to recur often. If you are unable to configure your API successfully,
you may want to consult https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/api/resource/iam_policy.rb
for additional configuration options.

Additionally, in order to generate IAM tests based on a preexisting resource
configuration, the first `examples` entry in `terraform.yaml` must be modified
to include a `primary_resource_name` entry:

```diff
      - !ruby/object:Provider::Terraform::Examples
        name: "disk_basic"
        primary_resource_id: "default"
+        primary_resource_name: "fmt.Sprintf(\"tf-test-test-disk%s\", context[\"random_suffix\"])"
        vars:
          disk_name: "test-disk"
```

`primary_resource_name` - Typically
`"fmt.Sprintf(\"tf-test-{{shortname}}%s\", context[\"random_suffix\"])"`,
substituting the parent resource's shortname from the example configuration for
`{{shortname}}`, such as `test-disk` above. This value is variable, as both the
key and value are user-defined parts of the example configuration. In some cases
the value must be customized further, albeit rarely.

Once an `iam_policy` block is added and filled out, and `primary_resource_name`
is set on the first example, you're finished, and you can run MMv1 to generate
the IAM resources you've added, alongside documentation, and tests.

## Adding IAM support to nonexistent resources

Some IAM targets don't exist as distinct resources, such as IAP, or their target
is supported through an engine other than MMv1 (i.e. through tpgtools/DCL or a
handwritten resource). For these resources, the `exclude_resource: true`
annotation can be used. To use it, partially define the resource in the
product's `api.yaml` file and apply the annotation. MMv1 won't attempt to
generate the resource itself and will only generate IAM resources targeting it.

The IAP product is a good reference for adding these: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/iap