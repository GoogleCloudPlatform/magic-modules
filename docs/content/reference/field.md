---
title: "MMv1 field reference"
weight: 20
aliases:
  - /reference/field-reference
  - /develop/field-reference
---

# MMv1 field reference

This page documents commonly-used properties for fields. For a full list of
available properties, see [type.go ↗](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/api/type.go).

## Shared properties

### `name`
Specifies the name of the field within Terraform. By default this will also 
be the key for the field in the API request message, if a separate `api_name`
is not declared using the corresponding property.

### `type`
Sets the expected data type of the field. All valid types are declared [here](https://github.com/GoogleCloudPlatform/magic-modules/blob/d7777055cb7618648725abd16d3b05e5c138fc56/mmv1/api/type.go#L673).

### `min_version: beta`
Marks the field (and any subfields) as beta-only. Ensure a beta version block
is present in provider.yaml. Do not use if an ancestor field (or the overall
resource) is already marked as beta-only.

### `immutable`
If true, the field is considered immutable - that is, only settable on create. If
unset or false, the field is considered to support update-in-place.

Immutability is not inherited from field to field: subfields are still considered to
be updatable in place by default. However, if the overall resource has
[`immutable`]({{< ref "/reference/resource#immutable" >}}) set to true, all its
fields are considered immutable.  Individual fields can override this for themselves
and their subfields with [`update_url`]({{< ref "/reference/field#update_url" >}})
if they have a custom update method in the API.

See [Best practices: Immutable fields]({{< ref "/best-practices/immutable-fields/" >}}) for more information.

Example:

```yaml
immutable: true
```

### `update_url`
If set, changes to the field's value trigger a separate call to a specific
API method for updating the field's value. Even if the overall resource is marked
immutable, the field and its subfields are not considered immutable unless explicitly
marked as such.

Terraform field names enclosed in double curly braces are replaced with the
field values from the resource at runtime.

Example:

```yaml
update_url: 'projects/{{project}}/locations/{{location}}/resourcenames/{{name}}/setFieldName'
```

### `update_verb`
If update_url is also set, overrides the verb used to update this specific
field. Allowed values: 'POST', 'PUT', 'PATCH'. Default: Resource's update_verb
(which defaults to 'PUT' if unset).

Example:

```yaml
update_verb: 'POST'
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
In this case, the field will also need to use [`ignore_read` or a `custom_flatten` function]({{< ref "/develop/diffs#ignore_read" >}}).

Example:

```yaml
sensitive: true
```

### `write_only_legacy` (deprecated)
If true, the field is considered "write-only", which means that its value will
be obscured in Terraform output as well as not be stored in state. This field is meant to replace `sensitive` as it doesn't store the value in state.
See [Ephemerality in Resources - Use Write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral/write-only)
for more information.

Write-only fields are only supported in Terraform v1.11+. Because the provider supports earlier Terraform versions, write only fields must be paired with (mutually exclusive) `sensitive` fields covering the same functionality for compatibility with those older versions.
This field cannot be used in conjuction with `immutable` or `sensitive`.

**Note**: Due to write-only not being read from the API, it is not possible to update the field directly unless a sidecar field is used. (e.g. `password` as a write-only field and `password_wo_version` as an immutable field meant for updating).

Example:

```yaml
write_only_legacy: true
```

**Deprecated**: This field is deprecated and will be removed in a future release.

### `ignore_read`
If true, the provider sets the field's value in the resource state based only
on the user's configuration. If false or unset, the provider sets the field's
value in the resource state based on the API response. Only use this attribute
if the field cannot be read from GCP due to either API or provider constraints.

`ignore_read` is current not supported inside arrays of nested objects. See [tpg#23630](https://github.com/hashicorp/terraform-provider-google/issues/23630)
for details and workarounds.

Example: YAML

```yaml
ignore_read: true
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

If true simulataneously with `default_from_api`, the provider will send empty values
explicitly set in configuration. If the field is unset, the provider will
accept API values as the default as usual with `default_from_api`.

Due to a [bug](https://github.com/hashicorp/terraform-provider-google/issues/13201),
NestedObject fields will currently be sent as `null` if unset (rather than being
omitted.)

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
- name: 'fieldOne'
  type: String
  conflicts:
    - field_two
    - nested_object.0.nested_field
```

### `required_with`
Specifies a list of fields (excluding the current field) that must all be specified
if at least one is specified. Must be set separately on all listed fields. Not supported within
[lists of nested objects](https://github.com/hashicorp/terraform-plugin-sdk/issues/470#issue-630928923).

Example:

```yaml
- name: 'fieldOne'
  type: String
  required_with:
    - field_two
    - nested_object.0.nested_field
```

### `exactly_one_of`
Specifies a list of fields (including the current field) of which exactly one
must be set. Must be set separately on all listed fields. Not supported within
[lists of nested objects](https://github.com/hashicorp/terraform-plugin-sdk/issues/470#issue-630928923).

Example:

```yaml
- name: 'fieldOne'
  type: String
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
- name: 'fieldOne'
  type: String
  at_least_one_of:
    - field_one
    - field_two
    - nested_object.0.nested_field
```

### `diff_suppress_func`
Specifies the name of a [diff suppress function](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#diffsuppressfunc)
to use for this field. In many cases, a [custom flattener]({{< ref "/develop/custom-code/#custom_flatten" >}})
is preferred because it will allow the user to see a clearer diff when the field actually is being changed. See
[Fix diffs]({{< ref "/develop/diffs" >}}) for more information and best practices.

The function specified can be a
[provider-specific function](https://github.com/hashicorp/terraform-provider-google-beta/blob/main/google-beta/tpgresource/common_diff_suppress.go)
(for example, `tgpresource.CaseDiffSuppress`) or a function defined in resource-specific
[custom code]({{<ref "/develop/custom-code#add-reusable-variables-and-functions" >}}).

Example:

```yaml
- name: 'fieldOne'
  type: String
  diff_suppress_func: 'tpgresource.CaseDiffSuppress'
```

### `validation`
In many cases, it is better to avoid client-side validation. See [Best practices: Validation]({{< ref "/best-practices/validation" >}}) for more information.

Controls the value set for the field's [`ValidateFunc`](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#validatefunc).

For Enum fields, this will override the default validation (that the provided value is one of the enum [`values`](#values)).
If you need additional validation on top of an enum, ensure that the supplied validation func also verifies the enum
values are correct.

This property has two mutually exclusive child properties:

- `function`: The name of a
  [validation function](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#validatefunc)
  to use for validation. The function can be a
  [Terraform-provided function](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/helper/validation)
  (for example, `validation.IntAtLeast(0)`), a
  [provider-specific function](https://github.com/hashicorp/terraform-provider-google-beta/blob/main/google-beta/verify/validation.go)
  (for example, `verify.ValidateBase64String`), or a function defined in
  resource-specific
  [custom code]({{<ref "/develop/custom-code#add-reusable-variables-and-functions" >}}).
- `regex`: A regex string to check values against. This can only be used on simple
  String fields. It is equivalent to
  [`function: verify.ValidateRegexp(REGEX_STRING)`](https://github.com/hashicorp/terraform-provider-google-beta/blob/0ef51142a4dd1c1a4fc308c1eb09dce307ebe5f5/google-beta/verify/validation.go#L425).

`validation` is not supported for Array fields (including sets); however, individual
elements in the array can be validated using [`item_validation`]({{<ref "#item_validation" >}}).

Example: Provider-specific function

```yaml
- name: 'fieldOne'
  type: String
  validation:
    function: 'verify.ValidateBase64String'
```

Example: Regex

```yaml
- name: 'fieldOne'
  type: String
  validation:
    regex: '^[a-zA-Z][a-zA-Z0-9_]*$'
```

### `is_set`
If true, the field is a Set rather than an Array. Set fields represent an
unordered set of unique elements. `set_hash_func` may be used to customize the
hash function used to index elements in the set, otherwise the schema default
function will be used. Adding this property to an existing field is usually a
breaking change.

```yaml
- name: 'fieldOne'
  type: Array
  is_set: true
```

### `set_hash_func`
Specifies a function for hashing elements in a Set field. If unspecified,
`schema.HashString` will be used if the elements are strings, otherwise
`schema.HashSchema`. The hash function should be defined in
`custom_code.constants`.

```yaml
set_hash_func: functionName
```

### `api_name`
Specifies a name to use for communication with the API that is different than
the name of the field in Terraform. In general, setting an `api_name` is not
recommended, because it makes it more difficult for users and maintainers to
understand how the resource maps to the underlying API.

```yaml
- name: 'fieldOne'
  type: String
  api_name: otherFieldName
```

### `url_param_only`
If true, the field is not sent in the resource body, and the provider does
not read the field value from the API response. If unset or false, the field
is sent in the resource body, and the provider reads the field value from the
API response.

```yaml
url_param_only: true
```

## `Enum` properties

### `enum_values`
Enum only. If the allowed values may change in the future, use a String field instead and link to API documentation
stating the current allowed values in the String field's description. 
See [Best practices: Validation]({{< ref "/best-practices/validation" >}}) for more information.

Do not include UNSPECIFIED values in this list.

Enums will validate that the provided field is in the allowed list unless a
custom [`validation`]({{<ref "#validation" >}}) is provided.

Example:

```yaml
enum_values:
  - 'VALUE_ONE'
  - 'VALUE_TWO'
```

## `Array` properties

### `item_type`
Array only. Sets the expected type of the items in the array. Primitives
should use the name of the primitive class as a string; other types should
define the attributes of the nested type.

Example: Primitive value

```yaml
item_type:
  type: String
```

Example: Enum value

```yaml
item_type:
  type: Enum
  description: 'required but unused'
  values:
    - 'VALUE_ONE'
    - 'VALUE_TWO'
```

Example: Nested object

```yaml
item_type:
  type: NestedObject
  properties:
    - name: 'FIELD_NAME'
      type: String
      description: |
        MULTI_LINE_FIELD_DESCRIPTION
```

### `item_validation`
Array only. Controls the [`ValidateFunc`](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#validatefunc)
used to validate individual items in the array. Behaves like [`validation`]({{<ref "#validation" >}}).

In many cases, it is better to avoid client-side validation. See [Best practices: Validation]({{< ref "/best-practices/validation" >}}) for more information.

For arrays of enums, this will override the default validation (that the provided value is one of the enum [`values`](#values)).
If you need additional validation on top of an enum, ensure that the supplied validation func also verifies the enum
values are correct.

Example: Provider-specific function

```yaml
- name: 'fieldOne'
  type: Array
  item_type:
    type: String
  item_validation:
    function: 'verify.ValidateBase64String'
```

Example: Regex

```yaml
- name: 'fieldOne'
  type: Array
  item_type:
    type: String
  item_validation:
    regex: '^[a-zA-Z][a-zA-Z0-9_]*$'
```

Example: Enum

```yaml
- name: 'fieldOne'
  type: Array
  item_type:
    type: Enum
    description: 'required but unused'
    values:
      - 'VALUE_ONE'
      - 'VALUE_TWO'
  item_validation: 
    function: 'customFunction'
```


## `NestedObject` properties

### `properties`
NestedObject only. Defines fields nested inside the current field.

Example:

```yaml
properties:
  - name: 'FIELD_NAME'
    type: String
    description: |
      MULTI_LINE_FIELD_DESCRIPTION
```
