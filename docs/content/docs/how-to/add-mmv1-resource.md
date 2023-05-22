---
title: "Add an MMv1 resource"
summary: "Generated resources are created using the `mmv1` code generator, and are
configured by editing definition files under the
[`mmv1/products`](https://github.com/GoogleCloudPlatform/magic-modules/tree/master/mmv1/products)
path."
weight: 10
---

# Add an MMv1 resource

Generated resources are created using the `mmv1` code generator, and are
configured by editing definition files under the
[`mmv1/products`](https://github.com/GoogleCloudPlatform/magic-modules/tree/master/mmv1/products)
path. Go to the service for your resource like
[`compute`](https://github.com/GoogleCloudPlatform/magic-modules/tree/master/mmv1/products/compute)
and open the `ResourceName.yaml` and `product.yaml` files.

For example, for `google_spanner_database`:

*   [`Database.yaml`](https://github.com/GoogleCloudPlatform/magic-modules/blob/67cef91ee76fc4871566f03e7caee1ef664f8aa0/mmv1/products/spanner/Database.yaml)

*   [`product.yaml`](https://github.com/GoogleCloudPlatform/magic-modules/blob/67cef91ee76fc4871566f03e7caee1ef664f8aa0/mmv1/products/spanner/product.yaml)

In short, `properties` is an array of the resource's fields. `Database.yaml` is named after the resource 
in PascalCase and contains the fields of the resource based on how it behaves in the API.
`product.yaml` is lower cased and contains the product-level information that applies to all resources
within that product.

## Field Configuration

### `ResourceName.yaml`

To add a field, you'll append an entry to `properties` within `ResourceName.yaml`, such
as the following adding support for a `fooBar` field in the API:

```yaml
      - !ruby/object:Api::Type::String
        name: 'fooBar'
        min_version: beta
        immutable: true
        description: |
          The cloud.google.com description of this field.
        default_from_api: true
        custom_expand: 'templates/terraform/custom_expand/shortname_to_url.go.erb'
```

The first line of that snippet is the type of the field in the API, including
primitives like `String`, `Integer`, `Boolean`, `Double`. Additional special
types are detailed below in "Complex Types".

You can configure settings on the field that describe it in the API. Avoid
setting values to `false`, and omit them instead.

*   `description` is the cloud.google.com description of the field, and must be
    filled out manually.
*   `min_version` can be set to `beta` if the field is only available at public
    preview or beta. If it is GA, do not set a `min_version` value.
*   `required: true` indicates that a field is required. New top-level fields
    should not be considered required, as that is a breaking change. Subfields
    of newly-added optional fields can be added as required.
*   `immutable: true` indicates that a field can only be set when the API resource is
    created. Changing the field will force the resource to be recreated.
*   `output: true` indicates that a field is output-only in the API and cannot
    be configured by the user.
*   `default_value: {{value}}` adds a default value for the field. This should
    only be used if the default value is fixed in the API.
*   `send_empty_value: true` indicates that an explicit zero value should be
    sent to the API. This is useful when a value has a nonzero default in the
    API but the zero value for the type can be set. This is extremely common for
    booleans that default to `true`.
*   `update_verb` and `update_url` configure a custom update function for a
    field.`update_verb`should be set to a literal symbol for the type (such
    as :POST for `POST`) and the URL to a templated URL such
    as`projects/{{project}}/global/backendServices/{{name}}/setSecurityPolicy`.
*   `default_from_api: true` indicates that Terraform needs to handle a field
    specially. This is common for fields with complex defaults from the API that
    can't be represented with `default_value`. If a `default_from_api: true`
    field is set in a user's config, Terraform will treat it as an optional
    field, detecting drift and correcting drift. If it is not set, it will be
    treated as an output-only field.
*   `ignore_read: true` indicates that a value is not returned from an API, and
    Terraform should not look for it in API responses.
*   `custom_expand` and `custom_flatten` are custom functions to read/write a
    value from state. They refer to files holding function contents under
    [`mmv1/template/terraform/custom_expand`](https://github.com/GoogleCloudPlatform/magic-modules/tree/8728bc89c37d5033b530c7d7157bb43865d9df58/mmv1/templates/terraform/custom_expand)
    and
    [`mmv1/template/terraform/custom_flatten`](https://github.com/GoogleCloudPlatform/magic-modules/tree/8728bc89c37d5033b530c7d7157bb43865d9df58/mmv1/templates/terraform/custom_flatten)
    respectively.

## Field Configuration - Complex Types

### Enum

```yaml
          - !ruby/object:Api::Type::Enum
            name: 'metadata'
            description: |
              Can only be specified if VPC flow logging for this subnetwork is enabled.
              Configures whether metadata fields should be added to the reported VPC
              flow logs.
            values:
              - :EXCLUDE_ALL_METADATA
              - :INCLUDE_ALL_METADATA
              - :CUSTOM_METADATA
            default_value: :INCLUDE_ALL_METADATA
```

Enum values represent enums in the underlying API where it is valuable to
restrict the range of inputs to a fixed set of values. They are strings that
support a `values` key to define the array of possible values specified as
literal constants, and `default_value` should be specified as a literal constant
as well.

Most API enums should be typed as `String` instead- if the value will not be
fixed for >1 year, use a `String`.

### ResourceRef

```yaml
      - !ruby/object:Api::Type::ResourceRef
        name: 'urlMap'
        resource: 'UrlMap'
        imports: 'selfLink'
        description: |
          A reference to the UrlMap resource that defines the mapping from URL
          to the BackendService.
```

ResourceRefs are fields that reference other resource. They're most typical in
GCE, and making a field a `ResourceRef` instead of a `String` will cause
Terraform to allow switching between reference formats and versions safely. If a
field can refer to multiple resource types, use a `String` instead.

In a `ResourceRef`, `resource` and `imports` must be defined but Terraform
ignores those values. `resource` should be set to the resource kind, and
`imports` to `selfLink` within GCE and `name` elsewhere.

### Array

```yaml
      - !ruby/object:Api::Type::Array
        name: scopes
        item_type: Api::Type::String
        description: |
          The list of scopes to be made available for this service
          account.
```

```yaml
      - !ruby/object:Api::Type::Array
        name: 'instances'
        description: |
          A list of virtual machine instances serving this pool.
          They must live in zones contained in the same region as this pool.
        item_type: !ruby/object:Api::Type::ResourceRef
          name: 'instance'
          description: 'The instance being served by this pool.'
          resource: 'Instance'
          imports: 'selfLink'
```

Arrays refer to arrays in the underlying API, with their item being specified
through an `item_type` field. `item_type` accepts any type, although primitives
(String / Integer / Boolean) must be specified differently than other types as
shown above.

### NestedObject

```yaml
      - !ruby/object:Api::Type::NestedObject
        name: 'imageEncryptionKey'
        description: |
          Encrypts the image using a customer-supplied encryption key.
          After you encrypt an image with a customer-supplied key, you must
          provide the same key if you use the image later (e.g. to create a
          disk from the image)
        properties:
          - !ruby/object:Api::Type::String
            name: 'rawKey'
            description: |
              Specifies a 256-bit customer-supplied encryption key, encoded in
              RFC 4648 base64 to either encrypt or decrypt this resource.
          - !ruby/object:Api::Type::String
            name: 'sha256'
            output: true
            description: |
              The RFC 4648 base64 encoded SHA-256 hash of the
              customer-supplied encryption key that protects this resource.
```

NestedObject is an object in the JSON API, and contains a `properties` subfield
where a sub-properties array can be defined (including additional NestedObjects)

### KeyValuePairs (Labels / Annotations)

```yaml
      - !ruby/object:Api::Type::KeyValuePairs
        name: 'labels'
        description: Labels to apply to this address.  A list of key->value pairs.
```

KeyValuePairs is a special type to handle string -> string maps, such as GCE
`labels` fields. No extra configuration is required.

### Exactly One Of

To restrain a parent object to contain exactly one of its nested objects,
use `exactly_one_of` in the affected child objects.

Yaml:

```yaml
objects:
  - !ruby/object:Api::Resource
    name: 'Connection'
    ...
    properties:
      - !ruby/object:Api::Type::NestedObject
        name: 'cloudSql'
        exactly_one_of:
          - cloud_sql
          - aws
        properties:
          ...
      - !ruby/object:Api::Type::NestedObject
        name: aws
        exactly_one_of:
          - cloud_sql
          - aws
        properties:
          ...
```

## Advanced customization

### DiffSuppressFunc

Terraform allows fields to specify a
[DiffSuppressFunc](https://www.terraform.io/plugin/sdkv2/schemas/schema-behaviors#diffsuppressfunc),
which allows you to ignore diffs in cases where the two values are
**functionally identical**. This is generally useful when the API returns a
normalized value - for example by standardizing the case.

Note: The *preferred* behavior for APIs is to always return the value that the
user sent. DiffSuppressFunc is a workaround for APIs that don't.

The Terraform provider comes with a set of
["common diff suppress functions"](https://github.com/hashicorp/terraform-provider-google/blob/main/google/common_diff_suppress.go).
These fit frequent needs like ignoring whitespace at the beginning and end of a
string, or ignoring case differences.

If you need to define a custom diff specifically for your resource, you can do
so in a "constants" file, which is a `.go.erb` file in
[mmv1/templates/terraform/constants](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/constants)
named `<product>_<resource>.erb`. You can then declare this custom code in
`ResourceName.yaml`:

```yaml
--- !ruby/object:Api::Resource
  name: ResourceName

  custom_code: !ruby/object:Provider::Terraform::CustomCode
    constants: templates/terraform/constants/product_resource_name.go.erb
```

Once you have chosen a DiffSuppressFunc, you can declare it as an override on
your resource:

```yaml
--- !ruby/object:Api::Resource
  name: ResourceName

  properties:
    - !ruby/object:Api::Type::String
      name: "myField"
      diff_suppress_func: 'tpgresource.CaseDiffSuppress'
```

The value of diff_suppress_func can be any valid DiffSuppressFunc, including the
result of a function call. For example:

```yaml
diff_suppress_func: 'tpgresource.OptionalPrefixSuppress("folders/")'
```

Please make sure to add thorough unit tests (in addition to basic integration
tests) for your diff suppress func.

Example: DomainMapping (DomainMappingLabelDiffSuppress)

-   [DomainMapping.yaml resource file](https://github.com/GoogleCloudPlatform/magic-modules/blob/67cef91ee76fc4871566f03e7caee1ef664f8aa0/mmv1/products/cloudrun/DomainMapping.yaml)
    -   [`custom_code`](https://github.com/GoogleCloudPlatform/magic-modules/blob/67cef91ee76fc4871566f03e7caee1ef664f8aa0/mmv1/products/cloudrun/DomainMapping.yaml#L40)
    -   [`diff_suppress_func: 'resourceBigQueryDatasetAccessRoleDiffSuppress'`](https://github.com/GoogleCloudPlatform/magic-modules/blob/67cef91ee76fc4871566f03e7caee1ef664f8aa0/mmv1/products/bigquery/DatasetAccess.yaml#L112)
-   [constants file](https://github.com/GoogleCloudPlatform/magic-modules/blob/67cef91ee76fc4871566f03e7caee1ef664f8aa0/mmv1/templates/terraform/constants/cloud_run_domain_mapping.go.erb)
-   [unit tests](https://github.com/GoogleCloudPlatform/magic-modules/blob/67cef91ee76fc4871566f03e7caee1ef664f8aa0/mmv1/third_party/terraform/tests/resource_cloud_run_domain_mapping_test.go#L9)

## Documentation

When adding a new MMv1 product or resource there are fields that you need to set within `product.yaml` and resource-level `ResourceName.yaml` that are specific to documentation for that resource. To learn more about MMv1 generated documentation and what YAML fields you need to pay attention to, see [Add and update MMv1 resource documentation](/magic-modules/docs/how-to/mmv1-resource-documentation).

# Beta features

When the underlying API of a feature is not final (i.e. a `vN` version like
`v1` or `v2`), is in preview, or the API has no SLO we add it to the
`google-beta` provider rather than the `google` provider, allowing users to
self-select for the stability level they are comfortable with.

In MMv1, a "version tag" can be annotated on resources, fields, resource iam
metadata and examples to control the stability level that a feature is available
at. Version tags are a specification of the minimum version a feature is
available at, written as `min_version: {{version}}`. This is only
specified when a feature is available at `beta`, and omitting a tag indicates
the target is generally available, or available at `ga`.

## Adding a beta resource

To add support for a beta resource in a preexisting product, ensure that a
`beta` level exists in the `versions` map in the `product.yaml` file for the product.
If one doesn't already exist, add it, setting the `base_url` to the appropriate
value. This is generally an API version including `beta`, such as `v1beta`, but
may be the same `base_url` as the `ga` entry for services that mix fields with
different stability levels within a single endpoint.

For example:

```diff
versions:
  - !ruby/object:Api::Product::Version
    name: ga
    base_url: https://compute.googleapis.com/compute/v1/
+  - !ruby/object:Api::Product::Version
+    name: beta
+    base_url: https://compute.googleapis.com/compute/beta/
```

If the product doesn't already exist, it's only necessary to add the `beta`
entry, i.e.:

```
versions:
  - !ruby/object:Api::Product::Version
    name: beta
    base_url: https://runtimeconfig.googleapis.com/v1beta1/
```

Next, annotate the resource (`Config.yaml`) i.e.:

```diff
---!ruby/object:Api::Resource
name: 'Config'
base_url: projects/{{project}}/configs
self_link: projects/{{project}}/configs/{{name}}
+min_version: beta
description: |
  A RuntimeConfig resource is the primary resource in the Cloud RuntimeConfig service.
  A RuntimeConfig resource consists of metadata and a hierarchy of variables.
iam_policy: !ruby/object:Api::Resource::IamPolicy
  parent_resource_attribute: 'config'
  method_name_separator: ':'
  exclude: false
properties:
...
```

You'll notice above that the `iam_policy` is not annotated with a version tag.
Due to the resource having a `min_version` tagged already, that's passed through
to the `iam_policy` (although the same is *not* true for `examples` entries used
to [create tests](#tests-that-use-a-beta-feature)). IAM-level tagging is only
necessary in the (rare) case that a resource is available at a higher stability
level than its `getIamPolicy`/`setIamPolicy` methods.

## Adding beta field(s)

NOTE: If a resource is already tagged as `min_version: beta`, follow the general
instructions for adding a field instead.

To add support for a beta field to a GA resource, ensure that the `beta` entry
already exists in the `versions` map for the product. See
[above](#adding-a-beta-resource) for details on doing so.

Next, add the field(s) as normal with a `min_version: beta` tag specified. In
the case of nested fields, only the highest-level field must be tagged, as
demonstrated below:

```diff
         - !ruby/object:Api::Type::NestedObject
            name: 'scaleDownControl'
+            min_version: beta
            description: |
              Defines scale down controls to reduce the risk of response latency
              and outages due to abrupt scale-in events
            properties:
            - !ruby/object:Api::Type::Integer
              name: 'timeWindowSec'
              description: |
                How long back autoscaling should look when computing recommendations
                to include directives regarding slower scale down, as described above.
...
```

## Promote a beta feature

In order to promote a beta feature to GA, remove the version tags previously set
on the feature or its tests. This will automatically make it available in the
`google` provider and remove the note that the feature is in beta in the provider
documentation.

For a resource, this typically means ensuring their removal:

* At the resource level, `ResourceName.yaml`
* On all resource examples, `examples` in `ResourceName.yaml` (unless some examples
use other beta resources or fields)
  * Additionally, for any modified examples, all `provider = google-beta`
annotations must be cleared

For a field, this typically means ensuring their removal:

* At the field level, `properties` in `ResourceName.yaml`
* On any resource examples where this was the last beta feature, `examples` in
`ResourceName.yaml`
  * Additionally, for any modified examples, all `provider = google-beta`
annotations must be cleared

If the feature was tested using handwritten tests, the version guards must be
removed, as described in the
[guidance for handwritten resources](third_party/terraform/README.md#promote-a-beta-feature).

When writing a changelog entry for a promotion, write it as if it was a new
field or resource, and suffix it with `(ga only)`. For example, if the
`google_container_cluster` resource was promoted to GA in your change:

~~~
```release-note:new-resource
`google_container_cluster` (ga only)
```
~~~

Alternatively, for field promotions, you may use `{{service}}: promoted
{{field}} in {{resource}} to GA`, i.e.

~~~
```release-note:enhancement
container: promoted `node_locations` field in google_container_cluster` to GA
```
~~~
