---
title: "MMv1 field reference"
weight: 31
aliases:
  - /reference/field-reference
---

# MMv1 field reference

This page documents commonly-used properties for fields. For a full list of
available properties, see [type.rb ↗](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/api/type.rb).

## Shared properties

### `min_version: beta`
Marks the field (and any subfields) as beta-only. Ensure a beta version block
is present in provider.yaml. Do not use if an ancestor field (or the overall
resource) is already marked as beta-only.

### `immutable`
If true, the field (and any subfields) are considered immutable - that is,
only settable on create. If unset or false, the field is still considered
immutable if any ancestor field (or the overall resource) is immutable,
unless `update_url` is set.

Example:

```yaml
immutable: true
```

### `update_url`
If set, changes to the field's value trigger a separate call to a specific
API method for updating the field's value. The field is not considered
immutable even if an ancestor field (or the overall resource) is immutable.
Terraform field names enclosed in double curly braces are replaced with the
field values from the resource at runtime.

Example:

```yaml
update_url: 'projects/{{project}}/locations/{{location}}/resourcenames/{{name}}/setFieldName'
```

### `update_verb`
If update_url is also set, overrides the verb used to update this specific
field. Allowed values: :POST, :PUT, :PATCH. Default: Resource's update_verb
(which defaults to :PUT if unset).

Example:

```yaml
update_verb: :POST
```

### `required`
If true, the field is required. If unset or false, the field is optional.

Example:

```yaml
required: true
```

### `output`
If true, the field is output-only - that is, it cannot be configured by the
user. If unset or false, the field is configurable.

Example:

```yaml
output: true
```

### `sensitive`
If true, the field is considered "sensitive", which means that its value will
be obscured in Terraform output such as plans. If false, the value will not be
obscured. Either way, the value will still be stored in plaintext in Terraform
state. See
[Handling Sensitive Values in State](https://developer.hashicorp.com/terraform/plugin/best-practices/sensitive-state)
for more information.

Sensitive fields are often not returned by the API (because they are sensitive).
In this case, the field will also need to use [`ignore_read` or a `custom_flatten` function]({{< ref "/develop/permadiff#ignore_read" >}}).

Example:

```yaml
sensitive: true
```

### `ignore_read`
If true, the provider sets the field's value in the resource state based only
on the user's configuration. If false or unset, the provider sets the field's
value in the resource state based on the API response. Only use this attribute
if the field cannot be read from GCP due to either API or provider constraints.

Nested fields currently
[do not support `ignore_read`](https://github.com/hashicorp/terraform-provider-google/issues/12410)
but can replicate the behavior by implementing a
[`custom_flatten`]({{< ref "/develop/custom-code#custom_flatten" >}})
that always ignores the value returned by the API.
](https://github.com/GoogleCloudPlatform/magic-modules/blob/5923d4cb878396a04bed9beaf22a8478e8b1e6a5/mmv1/templates/terraform/custom_flatten/source_representation_instance_configuration_password.go.erb).
Any fields using a custom flatten also need to be added to `ignore_read_extra`
for any examples where the field is set.

Example: YAML

```yaml
ignore_read: true
```

Example: Custom flatten

```go
func flatten<%= prefix -%><%= titlelize_property(property) -%>(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
  return d.Get("password")
}
```

### `default_value`
Sets a client-side default value for the field. This should be used if the
API has a default value that applies in all cases and is stable. Removing
or changing a default value is a breaking change. If unset, the field defaults
to an "empty" value (such as zero, false, or an empty string).

Example:

```yaml
default_value: DEFAULT_VALUE
```

### `default_from_api`
If true, and the field is either not set or set to an "empty" value (such as
zero, false, or empty strings), the provider accepts any value returned from
the API as the value for the field. If false, and the field is either not set
or set to an "empty" value, the provider treats the field's `default_value`
as the value for the field and shows a diff if the API returns any other
value for the field. This attribute is useful for complex or
frequently-changed API-side defaults, but provides less useful information at
plan time than `default_value` and causes the provider to ignore user
configurations that explicitly set the field to an "empty" value.
`default_from_api` and `send_empty_value` cannot both be true on the same field.

Example:

```yaml
default_from_api: true
```

### `send_empty_value`
If true, the provider sends "empty" values (such as zero, false, or empty
strings) to the API if set explicitly in the user's configuration. If false,
"empty" values cause the field to be omitted entirely from the API request.
This attribute is useful for fields where the API would behave differently
for an "empty" value vs no value for a particular field - for example,
boolean fields that have an API-side default of true.
`send_empty_value` and `default_from_api` cannot both be true on the same field.

Example:

```yaml
send_empty_value: true
```

### `conflicts`
Specifies a list of fields (excluding the current field) that cannot be
specified at the same time as the current field. Must be set separately on
all listed fields. Not supported within
[lists of nested objects](https://github.com/hashicorp/terraform-plugin-sdk/issues/470#issue-630928923).

Example:

```yaml
- !ruby/object:Api::Type::String
  name: 'fieldOne'
  conflicts:
    - field_two
    - nested_object.0.nested_field
```

### `exactly_one_of`
Specifies a list of fields (including the current field) of which exactly one
must be set. Must be set separately on all listed fields. Not supported within
[lists of nested objects](https://github.com/hashicorp/terraform-plugin-sdk/issues/470#issue-630928923).

Example:

```yaml
- !ruby/object:Api::Type::String
  name: 'fieldOne'
  exactly_one_of:
    - field_one
    - field_two
    - nested_object.0.nested_field
```

### `at_least_one_of`
Specifies a list of fields (including the current field) that cannot be
specified at the same time (but at least one of which must be set). Must be
set separately on all listed fields. Not supported within
[lists of nested objects](https://github.com/hashicorp/terraform-plugin-sdk/issues/470#issue-630928923).

Example:

```yaml
- !ruby/object:Api::Type::String
  name: 'fieldOne'
  at_least_one_of:
    - field_one
    - field_two
    - nested_object.0.nested_field
```

### `diff_suppress_func`
Specifies the name of a [diff suppress function](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#diffsuppressfunc)
to use for this field. In many cases, a [custom flattener](https://googlecloudplatform.github.io/magic-modules/develop/custom-code/#custom_flatten)
is preferred because it will allow the user to see a clearer diff when the field actually is being changed. See
[Fix a permadiff]({{< ref "/develop/permadiff.md" >}}) for more information and best practices.

Example:

```yaml
- !ruby/object:Api::Type::String
  name: 'fieldOne'
  diff_suppress_func: 'tpgresource.CaseDiffSuppress'
```

## `Enum` properties

### `values`
Enum only. Sets allowed values as ruby "literal constants" (prefixed with a
colon). If the allowed values change frequently, use a String field instead
to allow better forwards-compatibility, and link to API documentation
stating the current allowed values in the String field's description. Do not
include UNSPECIFIED values in this list.

Example:

```yaml
values:
  - :VALUE_ONE
  - :VALUE_TWO
```

## `Array` properties

### `item_type`
Array only. Sets the expected type of the items in the array. Primitives
should use the name of the primitive class as a string; other types should
define the attributes of the nested type.

Example: Primitive value

```yaml
item_type: Api::Type::String
```

Example: Enum value

```yaml
item_type: !ruby/object:Api::Type::Enum
  name: 'required but unused'
  description: 'required but unused'
  values:
    - :VALUE_ONE
    - :VALUE_TWO
```

Example: Nested object

```yaml
item_type: !ruby/object:Api::Type::NestedObject
  properties:
    - !ruby/object:Api::Type::String
      name: 'FIELD_NAME'
      description: |
        MULTI_LINE_FIELD_DESCRIPTION
```

## `NestedObject` properties

### `properties`
NestedObject only. Defines fields nested inside the current field.

Example:

```yaml
properties:
  - !ruby/object:Api::Type::String
    name: 'FIELD_NAME'
    description: |
      MULTI_LINE_FIELD_DESCRIPTION
```
