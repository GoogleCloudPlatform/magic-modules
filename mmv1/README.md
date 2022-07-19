# MMv1

## Overview

MMv1 is a Ruby-based code generator that implements Terraform Provider Google (TPG) resources from YAML specification files.

MMv1-generated resources like [google_compute_address](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address) can be identified by looking in their [Go source](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_compute_address.go) for an `AUTO GENERATED CODE` header as well as a Type `MMv1`. MMv1-generated resources should have source code present under their product folders, like [mmv1/products/compute](./products/compute) for the [google_compute_address](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address) resource.

## Table of Contents
- [Contributing](#contributing)
  - [Resource](#resource)
  - [IAM Resources](#iam-resource)
  - [Testing](#testing)
    - [Example Configuration File](#example-configuration-file)
    - [`terraform.yaml` metadata](#terraformyaml-metadata)
    - [Results](#results)
    - [Tests that use beta features](#tests-that-use-beta-features)
  - [Documentation](#documentation)
  - [Beta Feature](#beta-feature)
    - [Add or update a beta feature](#add-or-update-a-beta-feature)
    - [Tests that use a beta feature](#tests-that-use-a-beta-feature)
    - [Promote a beta feature](#promote-a-beta-feature)

## Contributing

We're glad to accept contributions to MMv1-generated resources. Tutorials and guidance on making changes are available below.

### Resource

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

#### Add or update a beta feature

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
