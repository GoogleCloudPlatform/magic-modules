---
title: "Add or modify a resource"
weight: 30
aliases:
  - /docs/how-to/add-mmv1-resource
  - /how-to/add-mmv1-resource
  - /develop/add-mmv1-resource
  - /docs/how-to/mmv1-resource-documentation
  - /how-to/mmv1-resource-documentation
  - /develop/mmv1-resource-documentation
  - /docs/how-to/add-mmv1-iam
  - /how-to/add-mmv1-iam
  - /develop/add-mmv1-iam
  - /docs/how-to/update-handwritten-resource
  - /how-to/update-handwritten-resource
  - /develop/update-handwritten-resource
  - /docs/how-to/update-handwritten-documentation
  - /how-to/update-handwritten-documentation
  - /develop/update-handwritten-documentation
  - /docs/how-to
  - /how-to
  - /docs/getting-started/provider-documentation
  - /getting-started/provider-documentation
---

# Add or modify a resource

This page describes how to add a new resource to the `google` or `google-beta` Terraform provider using MMv1 and/or handwritten code.

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

## Add a resource

{{< tabs "resource" >}}
{{< tab "MMv1" >}}
1. Using an editor of your choice, in the appropriate [product folder]({{<ref "/get-started/how-magic-modules-works.md#mmv1" >}}), create a file called `RESOURCE_NAME.yaml`. Replace `RESOURCE_NAME` with the name of the API resource you are adding support for. For example, a configuration file for [NatAddress](https://cloud.google.com/apigee/docs/reference/apis/apigee/rest/v1/organizations.instances.natAddresses) would be called `NatAddress.yaml`.
2. Copy the following template into the new file:
   ```yaml
   # Copyright 2024 Google Inc.
   # Licensed under the Apache License, Version 2.0 (the "License");
   # you may not use this file except in compliance with the License.
   # You may obtain a copy of the License at
   #
   #     http://www.apache.org/licenses/LICENSE-2.0
   #
   # Unless required by applicable law or agreed to in writing, software
   # distributed under the License is distributed on an "AS IS" BASIS,
   # WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   # See the License for the specific language governing permissions and
   # limitations under the License.

   --- !ruby/object:Api::Resource
   # API resource name
   name: 'ResourceName'
   # Resource description for the provider documentation.
   description: |
     RESOURCE_DESCRIPTION
   references: !ruby/object:Api::Resource::ReferenceLinks
     guides:
      # Link to quickstart in the API's Guides section. For example:
      # 'Create and connect to a database': 'https://cloud.google.com/alloydb/docs/quickstart/create-and-connect'
       'QUICKSTART_TITLE': 'QUICKSTART_URL'
     # Link to the REST API reference for the resource. For example,
     # https://cloud.google.com/alloydb/docs/reference/rest/v1/projects.locations.backups
     api: 'API_REFERENCE_URL'
   # Marks the resource as beta-only. Ensure a beta version block is present in
   # provider.yaml.
   # min_version: beta

   # Inserts styled markdown into the header of the resource's page in the
   # provider documentation.
   # docs: !ruby/object:Provider::Terraform::Docs
   #   warning: |
   #     MULTILINE_WARNING_MARKDOWN
   #   note: |
   #     MULTILINE_NOTE_MARKDOWN

   # URL for the resource's standard List method. https://google.aip.dev/132
   # Terraform field names enclosed in double curly braces are replaced with
   # the field values from the resource at runtime.
   base_url: 'projects/{{project}}/locations/{{location}}/resourcenames'
   # URL for the resource's standard Get method. https://google.aip.dev/131
   # Terraform field names enclosed in double curly braces are replaced with
   # the field values from the resource at runtime.
   self_link: 'projects/{{project}}/locations/{{location}}/resourcenames/{{name}}'

   # If true, the resource and all its fields are considered immutable - that is,
   # only creatable, not updatable. Individual fields can override this if they
   # have a custom update method in the API.
   # immutable: true

   # Overrides one or more timeouts, in minutes. All timeouts default to 20.
   # timeouts: !ruby/object:Api::Timeouts
   #   insert_minutes: 20 
   #   update_minutes: 20 
   #   delete_minutes: 20 

   # URL for the resource's standard Create method, including query parameters.
   # https://google.aip.dev/133
   # Terraform field names enclosed in double curly braces are replaced with
   # the field values from the resource at runtime.
   create_url: 'projects/{{project}}/locations/{{location}}/resourcenames?resourceId={{name}}'
   # Overrides the HTTP verb used to create a new resource.
   # Allowed values: :POST, :PUT, :PATCH. Default: :POST
   # create_verb: :POST

   # Overrides the URL for the resource's standard Update method. (If unset, the
   # self_link URL is used by default.) https://google.aip.dev/134
   # Terraform field names enclosed in double curly braces are replaced with
   # the field values from the resource at runtime.
   # update_url: 'projects/{{project}}/locations/{{location}}/resourcenames/{{name}}'
   # The HTTP verb used to update a resource. Allowed values: :POST, :PUT, :PATCH. Default: :PUT.
   update_verb: :PATCH
   # If true, the resource sets an `updateMask` query parameter listing modified
   # fields when updating the resource. If false, it does not.
   update_mask: true

   # Overrides the URL for the resource's standard Delete method. (If unset, the
   # self_link URL is used by default.) https://google.aip.dev/135
   # Terraform field names enclosed in double curly braces are replaced with
   # the field values from the resource at runtime.
   # delete_url: 'projects/{{project}}/locations/{{location}}/resourcenames/{{name}}'
   # Overrides the HTTP verb used to delete a resource.
   # Allowed values: :POST, :PUT, :PATCH, :DELETE. Default: :DELETE
   # delete_verb: :DELETE

   # If true, code for handling long-running operations is generated along with
   # the resource. If false, that code is not generated.
   autogen_async: true
   # Sets parameters for handling operations returned by the API.
   async: !ruby/object:Api::OpAsync
     # Overrides which API calls return operations. Default: ['create',
     # 'update', 'delete']
     # actions: ['create', 'update', 'delete']
     operation: !ruby/object:Api::OpAsync::Operation
       base_url: '{{op_id}}'

     # If true, the provider sets the resource's Terraform ID after the resource is created,
     # taking into account values that are set by the API at create time. This is only possible
     # when the completed operation's JSON includes the created resource in the "response" field.
     # If false (or unset), the provider sets the resource's Terraform ID before the resource is
     # created, based only on the resource configuration.
     # result: !ruby/object:Api::OpAsync::Result
     #   resource_inside_response: true

   # All resources (of all kinds) that share a mutex value block rather than
   # executing concurrent API requests.
   # Terraform field names enclosed in double curly braces are replaced with
   # the field values from the resource at runtime.
   # mutex: RESOURCE_NAME/{{name}}

   parameters:
     - !ruby/object:Api::Type::String
       name: 'location'
       required: true
       immutable: true
       url_param_only: true
       description: |
         LOCATION_DESCRIPTION
     - !ruby/object:Api::Type::String
       name: 'name'
       required: true
       immutable: true
       url_param_only: true
       description: |
         NAME_DESCRIPTION

   properties:
     # Fields go here
   ```

3. Modify the template as needed to match the API resource's documented behavior.
4. Delete all remaining comments in the resource configuration (including attribute descriptions) that were copied from the above template.

> **Note:** The template includes the most commonly-used fields. For a comprehensive reference, see [ResourceName.yaml reference ↗]({{<ref "/reference/resource-reference.md" >}}).
{{< /tab >}}
{{< tab "Handwritten" >}}
> **Warning:** Handwritten resources are more difficult to develop and maintain. New handwritten resources will only be accepted if implementing the resource in MMv1 would require entirely overriding two or more CRUD methods.

1. Add the resource in MMv1.
2. [Generate the beta provider]({{< ref "/get-started/generate-providers.md" >}})
3. From the beta provider, copy the files generated for the resource to the following locations:
   - Resource: Copy to the appropriate service folder inside [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services)
   - Documentation: [`magic-modules/mmv1/third_party/terraform/website/docs/r`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/website/docs/r)
   - Tests: Copy to the appropriate service folder inside [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services), and remove `_generated` from the filename
   - Sweepers: [`magic-modules/mmv1/third_party/terraform/utils`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/utils)
4. Modify the Go code as needed.
   - Replace all occurrences of `github.com/hashicorp/terraform-provider-google-beta/google-beta` with `github.com/hashicorp/terraform-provider-google/google`
   - Remove the `Example` suffix from all test function names.
   - Remove the comments at the top of the file.
   - If beta-only fields are being tested, do the following:
     - Change the file suffix to `.go.erb`
     - Add `<% autogen_exception -%>` to the top of the file
     - Wrap each beta-only test in a separate version guard: `<% unless version == 'ga' -%>...<% else -%>...<% end -%>`
5. Register the resource `handwrittenResources` in [`magic-modules/mmv1/third_party/terraform/provider/provider_mmv1_resources.go.erb`](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/provider/provider_mmv1_resources.go.erb)
   - Add a version guard for any beta-only resources.
6. Optional: Complete other handwritten tasks that require the MMv1 configuration file.
    - [Add resource tests]({{< ref "/develop/test/test.md" >}})
    - [Add IAM support]({{<ref "#add-iam-support" >}})
7. Delete the MMv1 configuration file.
{{< /tab >}}
{{< /tabs >}}

## Add fields

In general, Terraform resources should implement all configurable fields and all read-only fields.
Even fields that seem like they would not be useful in Terraform (like update time or etag) often
end up being requested by users, so it's usually easier to just add them all at once. However,
optional or read-only fields can be omitted when adding a resource if they would require significant
additional work to implement.

{{< tabs "fields" >}}
{{< tab "MMv1" >}}
1. For each API field, copy the following template into the resource's `properties` attribute. Be sure to indent appropriately.

{{< tabs "MMv1 types" >}}
{{< tab "Simple" >}}
```yaml
- !ruby/object:Api::Type::String
  name: 'API_FIELD_NAME'
  description: |
    MULTILINE_FIELD_DESCRIPTION
  min_version: beta
  immutable: true
  required: true
  output: true
  conflicts:
    - field_one
    - nested_object.0.nested_field
  exactly_one_of:
    - field_one
    - nested_object.0.nested_field

```

Replace `String` in the field type with one of the following options:

- `String`
- `Integer`
- `Boolean`
- `Double`
- `KeyValuePairs` (string -> string map)
- `KeyValueLabels` (for standard resource 'labels' field)
- `KeyValueAnnotations` (for standard resource 'annotations' field)
{{< /tab >}}
{{< tab "Enum" >}}
```yaml
- !ruby/object:Api::Type::Enum
  name: 'API_FIELD_NAME'
  description: |
    MULTILINE_FIELD_DESCRIPTION
  min_version: beta
  immutable: true
  required: true
  output: true
  conflicts:
    - field_one
    - nested_object.0.nested_field
  exactly_one_of:
    - field_one
    - nested_object.0.nested_field
  values:
    - :VALUE_ONE
    - :VALUE_TWO
```
{{< /tab >}}
{{< tab "ResourceRef" >}}
```yaml
- !ruby/object:Api::Type::ResourceRef
  name: 'API_FIELD_NAME'
  description: |
    MULTILINE_FIELD_DESCRIPTION
  min_version: beta
  immutable: true
  required: true
  output: true
  conflicts:
    - field_one
    - nested_object.0.nested_field
  exactly_one_of:
    - field_one
    - nested_object.0.nested_field
  resource: 'ResourceName'
  imports: 'name'
```
{{< /tab >}}
{{< tab "Array" >}}
```yaml
- !ruby/object:Api::Type::Array
  name: 'API_FIELD_NAME'
  description: |
    MULTILINE_FIELD_DESCRIPTION
  min_version: beta
  immutable: true
  required: true
  output: true
  conflicts:
    - field_one
    - nested_object.0.nested_field
  exactly_one_of:
    - field_one
    - nested_object.0.nested_field
  # Array of primitives
  item_type: Api::Type::String

  # Array of nested objects
  item_type: !ruby/object:Api::Type::NestedObject
    properties:
      - !ruby/object:Api::Type::String
        name: 'FIELD_NAME'
        description: |
          MULTI_LINE_FIELD_DESCRIPTION
```
{{< /tab >}}
{{< tab "NestedObject" >}}
```yaml
- !ruby/object:Api::Type::Array
  name: 'API_FIELD_NAME'
  description: |
    MULTILINE_FIELD_DESCRIPTION
  min_version: beta
  immutable: true
  required: true
  output: true
  conflicts:
    - field_one
    - nested_object.0.nested_field
  exactly_one_of:
    - field_one
    - nested_object.0.nested_field
  properties:
    - !ruby/object:Api::Type::String
      name: 'FIELD_NAME'
      description: |
        MULTI_LINE_FIELD_DESCRIPTION
```
{{< /tab >}}
{{< tab "Map" >}}
```yaml
  - !ruby/object:Api::Type::Map
    name: 'API_FIELD_NAME'
    description: |
      MULTILINE_FIELD_DESCRIPTION
    key_name: KEY_NAME
    key_description: |
      MULTILINE_KEY_FIELD_DESCRIPTION
    value_type: !ruby/object:Api::Type::NestedObject
      properties:
        - !ruby/object:Api::Type::String
          name: 'FIELD_NAME'
          description: |
            MULTI_LINE_FIELD_DESCRIPTION
```

This type is only used for string -> complex type mappings, use "KeyValuePairs" for simple mappings. Complex maps can't be represented natively in Terraform, and this type is transformed into an associative array (TypeSet) with the key merged into the object alongside other top-level fields.

For `key_name` and `key_description`, provide a domain-appropriate name and description. For example, a map that references a specific type of resource would generally use the singular resource kind as the key name (such as "topic" for PubSub Topic) and a descriptor of the expected format depending on the context (such as resourceId vs full resource name).

{{< /tab >}}
{{< /tabs >}}

2. Modify the field configuration according to the API documentation and behavior.

> **Note:** The templates in this section only include the most commonly-used fields. For a comprehensive reference, see [MMv1 field reference]({{<ref "/develop/field-reference.md" >}}). For information about modifying the values sent and received for a field, see [Modify the API request or response]({{<ref "/develop/custom-code#modify-the-api-request-or-response" >}}).
{{< /tab >}}
{{< tab "Handwritten" >}}
1. Add the field to the handwritten resource's schema.
   - The new field(s) should mirror the API's structure to ease predictability and maintenance. However, if there is an existing related / similar field in the resource that uses a different convention, follow that convention instead.
   - Enum fields in the API should be represented as `TypeString` in Terraform for forwards-compatibility. Link to the API documentation of allowed values in the field description.
   - Terraform field names should always use [snake case ↗](https://en.wikipedia.org/wiki/Snake_case).
   - See [Schema Types ↗](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-types) and [Schema Behaviors ↗](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors) for more information about field schemas.
2. Add handling for the new field in the resource's Create method and Update methods.
   - "Expanders" convert Terraform resource data to API request data.
   - For top level fields, add an expander. If the field is set or has changed, call the expander and add the resulting value to the API request.
   - For other fields, add logic to the parent field's expander to add the field to the API request. Use a nested expander for complex logic.
3. Add handling for the new field in the resource's Read method.
   - "Flatteners" convert API response data to Terraform resource data.
   - For top level fields, add a flattener. Call `d.Set()` on the flattened API response value to store it in Terraform state.
   - For other fields, add logic to the parent field's flattener to convert the value from the API response to the Terraform state value. Use a nested flattener for complex logic.
4. If any of the added Go code (including any imports) is beta-only, change the file suffix to `.go.erb` and wrap the beta-only code in a version guard: `<% unless version == 'ga' -%>...<% else -%>...<% end -%>`.
   - Add a new guard rather than adding the field to an existing guard; it is easier to read.
{{< /tab >}}
{{< /tabs >}}


## Add IAM support

This section covers how to add IAM resources in Terraform if they are supported by a particular API resource (indicated by `setIamPolicy` and `getIamPolicy` methods in the API documentation for the resource).

{{< tabs "IAM" >}}
{{< tab "MMv1" >}}
1. Add the following top-level block to `ResourceName.yaml` directly above `parameters`.

```yaml
iam_policy: !ruby/object:Api::Resource::IamPolicy
  # Name of the field on the terraform IAM resources which references
  # the parent resource. Update to match the parent resource's name.
  parent_resource_attribute: 'resource_name'
  # Character preceding setIamPolicy in the full URL for the API method.
  # Usually `:`
  method_name_separator: ':'
  # HTTP method for getIamPolicy. Usually :POST.
  # Allowed values: :GET, :POST. Default: :GET
  fetch_iam_policy_verb: :POST
  # Overrides the HTTP method for setIamPolicy.
  # Allowed values: :POST, :PUT. Default: :POST
  # set_iam_policy_verb: :POST

  # Must match the parent resource's `import_format` (or `self_link` if
  # `import_format` is unset), but with the `parent_resource_attribute`
  # value substituted for the final field.
  import_format: [
    'projects/{{project}}/locations/{{location}}/resourcenames/{{resource_name}}'
  ]
  # Valid IAM role that can be set by generated tests. Default: 'roles/viewer'
  # allowed_iam_role: 'roles/viewer'

  # If IAM conditions are supported, set this attribute to indicate how the
  # conditions should be passed to the API. Allowed values: :QUERY_PARAM,
  # :REQUEST_BODY, :QUERY_PARAM_NESTED. Note: :QUERY_PARAM_NESTED should
  # only be used if the query param field contains a `.`
  # iam_conditions_request_type: :REQUEST_BODY

  # Marks IAM support as beta-only
  # min_version: beta
```

2. Modify the template as needed to match the API resource's documented behavior. These are the most commonly-used fields. For a comprehensive reference, see [IAM policy YAML reference ↗]({{<ref "/reference/iam-policy-reference.md" >}}).
3. Delete all remaining comments in the IAM configuration (including attribute descriptions) that were copied from the above template.
{{< /tab >}}
{{< tab "Handwritten" >}}
> **Warning:** IAM support for handwritten resources should be implemented using MMv1. New handwritten IAM resources will only be accepted if they cannot be implemented using MMv1.

### Add support in MMv1

1. Follow the MMv1 directions in [Add the resource]({{<ref "#add-the-resource" >}}) to create a skeleton ResourceName.yaml file for the handwritten resource, but set only the following top-level fields:
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

1. [Generate the beta provider]({{< ref "/get-started/generate-providers.md" >}})
2. From the beta provider, copy the files generated for the IAM resources to the following locations:
   - Resource: Copy to the appropriate service folder inside [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services)
   - Documentation: [`magic-modules/mmv1/third_party/terraform/website/docs/r`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/website/docs/r)
   - Tests: In the appropriate service folder inside [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services)
3. Modify the Go code as needed.
   - Replace all occurrences of `github.com/hashicorp/terraform-provider-google-beta/google-beta` with `github.com/hashicorp/terraform-provider-google/google`
   - Remove the comments at the top of the file.
   - If any of the added Go code is beta-only:
     - Change the file suffix to `.go.erb`
     - Add `<% autogen_exception -%>` to the top of the file
     - Wrap each beta-only code block (including any imports) in a separate version guard: `<% unless version == 'ga' -%>...<% else -%>...<% end -%>`
4. Register the binding, member, and policy resources `handwrittenIAMResources` in [`magic-modules/mmv1/third_party/terraform/provider/provider_mmv1_resources.go.erb`](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/provider/provider_mmv1_resources.go.erb)
   - Add a version guard for any beta-only resources.
{{< /tab >}}
{{< /tabs >}}

## Add documentation

{{< tabs "docs" >}}
{{< tab "MMv1" >}}
Documentation is autogenerated based on the resource and field configurations. To preview the documentation:

1. [Generate the providers]({{< ref "/get-started/generate-providers.md" >}})
2. Copy and paste the generated documentation into the Hashicorp Registry's [Doc Preview Tool](https://registry.terraform.io/tools/doc-preview) to see how it is rendered.
{{< /tab >}}
{{< tab "Handwritten" >}}
### Add or modify documentation files

1. Open the resource documentation in [`magic-modules/third_party/terraform/website/docs/r/`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/website/docs/r) using an editor of your choice.
   - The name of the file is the name of the resource without a `google_` prefix. For example, for `google_compute_instance`, the file is called `compute_instance.html.markdown`
2. Modify the documentation as needed according to [Handwritten documentation style guide]({{< ref "/develop/handwritten-docs-style-guide" >}}).
3. [Generate the providers]({{< ref "/get-started/generate-providers.md" >}})
4. Copy and paste the generated documentation into the Hashicorp Registry's [Doc Preview Tool](https://registry.terraform.io/tools/doc-preview) to see how it is rendered.
{{< /tab >}}
{{< /tabs >}}

## What's next?

- [Add custom resource code]({{< ref "/develop/custom-code.md" >}})
- [Add tests]({{< ref "/develop/test/test.md" >}})
- [Run tests]({{< ref "/develop/test/run-tests.md" >}})
