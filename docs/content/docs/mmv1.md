---
title: "MMv1"
weight: 1
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
---


# MMv1

## Overview

MMv1 is a Ruby-based code generator that implements Terraform Provider Google (TPG) resources from YAML specification files.

MMv1-generated resources like [google_compute_address](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address) can be identified by looking in their [Go source](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_compute_address.go) for an `AUTO GENERATED CODE` header as well as a Type `MMv1`. MMv1-generated resources should have source code present under their product folders, like [mmv1/products/compute](./products/compute) for the [google_compute_address](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address) resource.

## Table of Contents
- [Contributing](#contributing)
  - [Resource](#resource)
    - [Field Configuration](#field-configuration)
      - [`api.yaml`](#apiyaml)
      - [`terraform.yaml`](#terraformyaml)
    - [Field Configuration - Complex Types](#field-configuration---complex-types)
      - [Enum](#enum)
      - [ResourceRef](#resourceref)
      - [Array](#array)
      - [NestedObject](#nestedobject)
      - [KeyValuePairs (Labels / Annotations)](#keyvaluepairs-labels--annotations)
      - [Exactly One Of](#exactly-one-of)
    - [Advanced customization](#advanced-customization)
      - [DiffSuppressFunc](#diffsuppressfunc)
  - [IAM Resources](#iam-resource)
  - [Testing](#testing)
    - [Example Configuration File](#example-configuration-file)
    - [`terraform.yaml` metadata](#terraformyaml-metadata)
    - [Results](#results)
    - [Tests that use beta features](#tests-that-use-beta-features)
  - [Documentation](#documentation)
  - [Beta Feature](#beta-feature)
    - [Adding a beta resource](#adding-a-beta-resource)
    - [Adding beta field(s)](#adding-beta-fields)
    - [Tests that use a beta feature](#tests-that-use-a-beta-feature)
    - [Promote a beta feature](#promote-a-beta-feature)

## Contributing

We're glad to accept contributions to MMv1-generated resources. Tutorials and guidance on making changes are available below.

### Resource

Generated resources are created using the `mmv1` code generator, and are
configured by editing definition files under the
[`mmv1/products`](https://github.com/GoogleCloudPlatform/magic-modules/tree/master/mmv1/products)
path. Go to the service for your resource like
[`compute`](https://github.com/GoogleCloudPlatform/magic-modules/tree/master/mmv1/products/compute)
and open the `api.yaml` and `terraform.yaml` files. In each of those, find the
resource's `properties` field.

For example, for `google_spanner_database`:

*   [`api.yaml`](https://github.com/GoogleCloudPlatform/magic-modules/blob/8728bc89c37d5033b530c7d7157bb43865d9df58/mmv1/products/spanner/api.yaml#L123-L158)

*   [`terraform.yaml`](https://github.com/GoogleCloudPlatform/magic-modules/blob/8728bc89c37d5033b530c7d7157bb43865d9df58/mmv1/products/spanner/terraform.yaml#L16-L51)

In short,`properties` is an array of the resource's fields. `api.yaml` it
contains the fields of the resource based on how it behaves in the API, and
`terraform.yaml` contains Terraform-specific amendments to those fields'
behaviour. Not all fields will need to be added to `terraform.yaml`- only add an
entry for your field if you need to configure one of the available option(s).

#### Field Configuration

##### `api.yaml`

To add a field, you'll append an entry to `properties` within `api.yaml`, such
as the following adding support for a `fooBar` field in the API:

```yaml
      - !ruby/object:Api::Type::String
        name: 'fooBar'
        min_version: beta
        input: true
        description: |
          The cloud.google.com description of this field.
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
*   `input: true` indicates that a field can only be set when the API resource is
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

##### `terraform.yaml`

You can add additional values within `terraform.yaml`:

```yaml
      foobar: !ruby/object:Overrides::Terraform::PropertyOverride
        ignore_read: true
        default_from_api: true
        custom_expand: 'templates/terraform/custom_expand/shortname_to_url.go.erb'
```

Commonly configured values include the following:

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

#### Field Configuration - Complex Types

##### Enum

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

##### ResourceRef

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

##### Array

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

##### NestedObject

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

##### KeyValuePairs (Labels / Annotations)

```yaml
      - !ruby/object:Api::Type::KeyValuePairs
        name: 'labels'
        description: Labels to apply to this address.  A list of key->value pairs.
```

KeyValuePairs is a special type to handle string -> string maps, such as GCE
`labels` fields. No extra configuration is required.

##### Exactly One Of

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

#### Advanced customization

##### DiffSuppressFunc

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
`terraform.yaml`:

```yaml
--- !ruby/object:Provider::Terraform::Config
overrides: !ruby/object:Overrides::ResourceOverrides
  ResourceName:
    # various overrides go here
    custom_code: !ruby/object:Provider::Terraform::CustomCode
      constants: templates/terraform/constants/product_resource_name.go.erb
```

Once you have chosen a DiffSuppressFunc, you can declare it as an override on
your resource:

```yaml
--- !ruby/object:Provider::Terraform::Config
overrides: !ruby/object:Overrides::ResourceOverrides
  ResourceName:
    # various overrides go here
    properties:
      myField: !ruby/object:Overrides::Terraform::PropertyOverride
        diff_suppress_func: 'caseDiffSuppress'
```

The value of diff_suppress_func can be any valid DiffSuppressFunc, including the
result of a function call. For example:

```yaml
diff_suppress_func: 'optionalPrefixSuppress("folders/")'
```

Please make sure to add thorough unit tests (in addition to basic integration
tests) for your diff suppress func.

Example: DomainMapping (domainMappingLabelDiffSuppress)

-   [terraform.yaml resource overrides](https://github.com/GoogleCloudPlatform/magic-modules/blob/15fd46f60ed49ec1a6488d1b34394dcbd7cd3a41/mmv1/products/cloudrun/terraform.yaml#L16)
    -   [`custom_code`](https://github.com/GoogleCloudPlatform/magic-modules/blob/15fd46f60ed49ec1a6488d1b34394dcbd7cd3a41/mmv1/products/cloudrun/terraform.yaml#L31)
    -   [`diff_suppress_func: 'resourceBigQueryDatasetAccessRoleDiffSuppress'`](https://github.com/GoogleCloudPlatform/magic-modules/blob/15fd46f60ed49ec1a6488d1b34394dcbd7cd3a41/mmv1/products/cloudrun/terraform.yaml#L46)
-   [constants file](https://github.com/GoogleCloudPlatform/magic-modules/blob/15fd46f60ed49ec1a6488d1b34394dcbd7cd3a41/mmv1/templates/terraform/constants/cloud_run_domain_mapping.go.erb)
-   [unit tests](https://github.com/GoogleCloudPlatform/magic-modules/blob/15fd46f60ed49ec1a6488d1b34394dcbd7cd3a41/mmv1/third_party/terraform/tests/resource_cloud_run_domain_mapping_test.go#L9)

### IAM Resource

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

#### Adding IAM support to nonexistent resources

Some IAM targets don't exist as distinct resources, such as IAP, or their target
is supported through an engine other than MMv1 (i.e. through tpgtools/DCL or a
handwritten resource). For these resources, the `exclude_resource: true`
annotation can be used. To use it, partially define the resource in the
product's `api.yaml` file and apply the annotation. MMv1 won't attempt to
generate the resource itself and will only generate IAM resources targeting it.

The IAP product is a good reference for adding these: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/iap

### Testing

For generated resources, you can add an example to the
[`mmv1/templates/terraform/examples`](https://github.com/GoogleCloudPlatform/magic-modules/tree/master/mmv1/templates/terraform/examples)
directory, which contains a set of templated Terraform configurations.

After writing out the example and filling out some metadata, Magic Modules will
insert it into the resource documentation page, and generate a test case
stepping through the following stages:

1.  Run `terraform apply` on the configuration, waiting for it to succeed and
    recording the results in Terraform state
1.  Run `terraform plan`, and fail if Terraform detects any drift
1.  Clear the resource from state and run `terraform import` on it
1.  Deeply compare the original state from `terraform apply` and the `terraform
    import` results, returning an error if any values are not identical
1.  Destroy all resources in the configuration using `terraform destroy`,
    waiting for the destroy command to succeed
1.  Call `GET` on the resource, and fail the test if it is still present

#### Example Configuration File

First, you'll want to add the example file. It needs to end in the filename
`.tf.erb`, and is typically named `service_resource_descriptive_name`. For
example, `pubsub_topic_geo_restricted.tf.erb`. Inside, you'll write a complete
Terraform configuration that provisions the resource and all of the required
dependencies. For example, in
[`mmv1/templates/terraform/examples/pubsub_subscription_dead_letter.tf.erb`](https://github.com/GoogleCloudPlatform/magic-modules/blob/e7ef590f6007796f446b2d41875b3d26f4469ff4/mmv1/templates/terraform/examples/pubsub_subscription_dead_letter.tf.erb):

```tf
resource "google_pubsub_topic" "<%= ctx[:primary_resource_id] %>" {
  name = "<%= ctx[:vars]['topic_name'] %>"
}

resource "google_pubsub_topic" "<%= ctx[:primary_resource_id] %>_dead_letter" {
  name = "<%= ctx[:vars]['topic_name'] %>-dead-letter"
}

resource "google_pubsub_subscription" "<%= ctx[:primary_resource_id] %>" {
  name  = "<%= ctx[:vars]['subscription_name'] %>"
  topic = google_pubsub_topic.<%= ctx[:primary_resource_id] %>.name

  dead_letter_policy {
    dead_letter_topic = google_pubsub_topic.<%= ctx[:primary_resource_id] %>_dead_letter.id
    max_delivery_attempts = 10
  }
}
```

The `ctx` variable provides metadata at generation time, and should be used in
two ways:

*   The Terraform ID of a single instance of the primary resource should be
    supplied through `<%= ctx[:primary_resource_id] %>` (in this example
    multiple resources use the value, although only the first
    `google_pubsub_topic` requires it). The resource kind you are testing with
    an id equal to `<%= ctx[:primary_resource_id] %>` is the one that will be
    imported.
*   Unique values can be supplied through `<%= ctx[:vars]['{{var}}'] %>`, where
    `{{var}}` is an arbitrary key you define. These values are created by
    appending suffixes to them, and are typically only used for names- most
    values should be constant within the configuration.

#### `terraform.yaml` metadata

Once your configuration is written, go in `terraform.yaml` and find the
`examples` block for the resource. Generally it'll be above the `properties`
block. In there, append an entry such as the
[following](https://github.com/GoogleCloudPlatform/magic-modules/blob/e7ef590f6007796f446b2d41875b3d26f4469ff4/mmv1/products/pubsub/terraform.yaml#L108-L113):

```yaml
      - !ruby/object:Provider::Terraform::Examples
        name: "pubsub_subscription_dead_letter"
        primary_resource_id: "example"
        vars:
          topic_name: "example-topic"
          subscription_name: "example-subscription"
```

The `name` should match the base name of your example file,
`primary_resource_id` is an arbitrary snake_cased string that describes the
resource, and the `vars` map should contain each key you defined previously.

**Important**: Any vars that are part of the resource's id should include at
least one hyphen or underscore; this
[triggers addition of a `tf-test` or `tf_test` prefix](https://github.com/GoogleCloudPlatform/magic-modules/blob/6858338f013f5dc57729ec037883a3594441ea62/mmv1/provider/terraform/examples.rb#L244),
which is what we use to detect and delete stray resources that are sometimes
left over during test runs.

#### Results

Your configuration will ultimately generate a Go test case similar to the
[following](https://github.com/hashicorp/terraform-provider-google/blob/38e2913cb102225f9f9bda9f04b5498d3386a79c/google/resource_pubsub_subscription_generated_test.go#L135-L180)
based on the snippets above:

```go
func TestAccPubsubSubscription_pubsubSubscriptionDeadLetterExample(t *testing.T) {
    t.Parallel()

    context := map[string]interface{}{
        "random_suffix": randString(t, 10),
    }

    vcrTest(t, resource.TestCase{
        PreCheck:     func() { testAccPreCheck(t) },
        Providers:    testAccProviders,
        CheckDestroy: testAccCheckPubsubSubscriptionDestroyProducer(t),
        Steps: []resource.TestStep{
            {
                Config: testAccPubsubSubscription_pubsubSubscriptionDeadLetterExample(context),
            },
            {
                ResourceName:            "google_pubsub_subscription.example",
                ImportState:             true,
                ImportStateVerify:       true,
                ImportStateVerifyIgnore: []string{"topic"},
            },
        },
    })
}

func testAccPubsubSubscription_pubsubSubscriptionDeadLetterExample(context map[string]interface{}) string {
    return Nprintf(`
resource "google_pubsub_topic" "example" {
  name = "tf-test-example-topic%{random_suffix}"
}
resource "google_pubsub_topic" "example_dead_letter" {
  name = "tf-test-example-topic%{random_suffix}-dead-letter"
}
resource "google_pubsub_subscription" "example" {
  name  = "tf-test-example-subscription%{random_suffix}"
  topic = google_pubsub_topic.example.name
  dead_letter_policy {
    dead_letter_topic = google_pubsub_topic.example_dead_letter.id
    max_delivery_attempts = 10
  }
}
`, context)
}
```

#### Tests that use beta features

See [Tests that use a beta feature](#tests-that-use-a-beta-feature)

### Documentation

### Beta Feature

When the underlying API of a feature is not final (i.e. a `vN` version like
`v1` or `v2`), is in preview, or the API has no SLO we add it to the
`google-beta` provider rather than the `google `provider, allowing users to
self-select for the stability level they are comfortable with.

In MMv1, a "version tag" can be annotated on resources, fields, resource iam
metadata and examples to control the stability level that a feature is available
at. Version tags are a specification of the minimum version a feature is
available at, written as `min_version: {{version}}`. This is only
specified when a feature is available at `beta`, and omitting a tag indicates
the target is generally available, or available at `ga`.

#### Adding a beta resource

To add support for a beta resource in a preexisting product, ensure that a
`beta` level exists in the `versions` map in the `api.yaml` file for the product.
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

Next, annotate the resource (part of `resources` in `api.yaml`) i.e.:

```diff
  - !ruby/object:Api::Resource
    name: 'Config'
    base_url: projects/{{project}}/configs
    self_link: projects/{{project}}/configs/{{name}}
+    min_version: beta
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

#### Adding beta field(s)

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

#### Tests that use a beta feature

For tests that use beta features, you'll need to perform two additional steps:

1.  Add `provider = google-beta` to every resource in the test (even resources
    that aren't being tested and/or are also in the GA provider)
1.  Add `min_version: beta` to the `Provider::Terraform::Examples` block

For example, modifying the snippets above:

```tf
resource "google_pubsub_topic" "<%= ctx[:primary_resource_id] %>" {
  provider = google-beta

  name = "<%= ctx[:vars]['topic_name'] %>"
}

resource "google_pubsub_topic" "<%= ctx[:primary_resource_id] %>_dead_letter" {
  provider = google-beta

  name = "<%= ctx[:vars]['topic_name'] %>-dead-letter"
}

resource "google_pubsub_subscription" "<%= ctx[:primary_resource_id] %>" {
  provider = google-beta

  name  = "<%= ctx[:vars]['subscription_name'] %>"
  topic = google_pubsub_topic.<%= ctx[:primary_resource_id] %>.name

  dead_letter_policy {
    dead_letter_topic = google_pubsub_topic.<%= ctx[:primary_resource_id] %>_dead_letter.id
    max_delivery_attempts = 10
  }
}
```

```yaml
      - !ruby/object:Provider::Terraform::Examples
        name: "pubsub_subscription_dead_letter"
        min_version: beta
        primary_resource_id: "example"
        vars:
          topic_name: "example-topic"
          subscription_name: "example-subscription"
```

#### Promote a beta feature

In order to promote a beta feature to GA, remove the version tags previously set
on the feature or its tests. This will automatically make it available in the
`google` provider and remove the note that the feature is in beta in the provider
documentation.

For a resource, this typically means ensuring their removal:

* At the resource level, `resources` in `api.yaml`
* On all resource examples, `examples` in `terraform.yaml` (unless some examples
use other beta resources or fields)
  * Additionally, for any modified examples, all `provider = google-beta`
annotations must be cleared

For a field, this typically means ensuring their removal:

* At the field level, `properties` in `api.yaml`
* On any resource examples where this was the last beta feature, `examples` in
`terraform.yaml`
  * Additionally, for any modified examples, all `provider = google-beta`
annotations must be cleared

If the feature was tested using handwritten tests, the version guards must be
removed, as described in the
[guidance for handwritten resources](third_party/terraform/README.md#promote-a-beta-feature).

When writing a changelog entry for a promotion, write it as if it was a new
field or resource, and suffix it with `(ga only)`. For example, if the
`google_container_cluster` resource was promoted to GA in your change:

```
\`\`\`release-note:new-resource
`google_container_cluster` (ga only)
\`\`\`
```

Alternatively, for field promotions, you may use "{{service}}: promoted
{{field}} in {{resource}} to GA", i.e.

```
\`\`\`release-note:enhancement
container: promoted `node_locations` field in google_container_cluster` to GA
\`\`\`
```

