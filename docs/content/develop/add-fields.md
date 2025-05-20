---
title: "Add a field to an existing resource"
weight: 30
---

# Add a field to an existing resource

This page describes how to add a field to an existing resource in the `google` or `google-beta` Terraform
provider using MMv1 and/or handwritten code. In general, Terraform resources should implement all configurable
fields and all read-only fields. Even fields that seem like they would not be useful in Terraform
(like update time or etag) often end up being requested by users, so it's usually easier to just add them all at
once. However, optional or read-only fields can be omitted when adding a resource if they would require significant
additional work to implement.

For more information about types of resources and the generation process overall, see [How Magic Modules works]({{< ref "/" >}}).

## Before you begin

1. Complete the steps in [Set up your development environment]({{< ref "/develop/set-up-dev-environment" >}}) to set up your environment and your Google Cloud project.
1. [Ensure the resource to which you want to add the fields exists in the provider]({{< ref "/develop/add-resource" >}}).
1. Ensure that your `magic-modules`, `terraform-provider-google`, and `terraform-provider-google-beta` repositories are up to date.
   ```
   cd ~/magic-modules
   git checkout main && git clean -f . && git checkout -- . && git pull
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
   git checkout main && git clean -f . && git checkout -- . && git pull
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
   git checkout main && git clean -f . && git checkout -- . && git pull
   ```

## Add fields

{{< tabs "fields" >}}
{{< tab "MMv1" >}}
1. For each API field, copy the following template into the resource's `properties` attribute. Be sure to indent appropriately.

{{< tabs "MMv1 types" >}}
{{< tab "Simple" >}}
```yaml
- name: 'API_FIELD_NAME'
  type: String
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
- name: 'API_FIELD_NAME'
  type: Enum
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
  enum_values:
    - 'VALUE_ONE'
    - 'VALUE_TWO'
```
{{< /tab >}}
{{< tab "ResourceRef" >}}
```yaml
- name: 'API_FIELD_NAME'
  type: ResourceRef
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
- name: 'API_FIELD_NAME'
  type: Array
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
  item_type: 
    type: String

  # Array of nested objects
  item_type: 
    type: NestedObject
    properties:
      - name: 'FIELD_NAME'
        type: String
        description: |
          MULTI_LINE_FIELD_DESCRIPTION
```
{{< /tab >}}
{{< tab "NestedObject" >}}
```yaml
- name: 'API_FIELD_NAME'
  type: NestedObject
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
    - name: 'FIELD_NAME'
      type: String
      description: |
        MULTI_LINE_FIELD_DESCRIPTION
```
{{< /tab >}}
{{< tab "Map" >}}
```yaml
  - name: 'API_FIELD_NAME'
    type: Map
    description: |
      MULTILINE_FIELD_DESCRIPTION
    key_name: 'KEY_NAME'
    key_description: |
      MULTILINE_KEY_FIELD_DESCRIPTION
  # Map of primitive values
    value_type:
      name: mapIntegerName
      type: Integer

  # Map of complex values
    value_type:
      name: mapObjectName
      type: NestedObject
      properties:
      - name: 'FIELD_NAME'
        type: String
        description: |
          MULTI_LINE_FIELD_DESCRIPTION
```

This type is used for general-case string -> non-string type mappings, use "KeyValuePairs" for string -> string mappings. Complex maps can't be represented natively in Terraform, and this type is transformed into an associative array (TypeSet) with the key merged into the object alongside other top-level fields.

For `key_name` and `key_description`, provide a domain-appropriate name and description. For example, a map that references a specific type of resource would generally use the singular resource kind as the key name (such as "topic" for PubSub Topic) and a descriptor of the expected format depending on the context (such as resourceId vs full resource name).

{{< /tab >}}
{{< /tabs >}}

2. Modify the field configuration according to the API documentation and behavior.

> **Note:** The templates in this section only include the most commonly-used fields. For a comprehensive reference, see [MMv1 field reference]({{<ref "/reference/field" >}}). For information about modifying the values sent and received for a field, see [Modify the API request or response]({{<ref "/develop/custom-code#modify-the-api-request-or-response" >}}).
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
4. If any of the added Go code (including any imports) is beta-only, change the file suffix to `.go.tmpl` and wrap the beta-only code in a version guard: `{{- if ne $.TargetVersionName "ga" -}}...{{- else }}...{{- end }}`.
   - Add a new guard rather than adding the field to an existing guard; it is easier to read.
{{< /tab >}}
{{< /tabs >}}

## What's next?

+ [Add IAM support]({{< ref "/develop/add-iam-support" >}})
+ [Add documentation]({{< ref "/document/add-documentation" >}})
+ [Add custom resource code]({{< ref "/develop/custom-code" >}})
+ [Add tests]({{< ref "/test/test" >}})
+ [Run tests]({{< ref "/test/run-tests" >}})
