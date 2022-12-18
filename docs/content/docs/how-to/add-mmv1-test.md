---
title: "Add an MMv1 test"
summary: "An example terraform configuration can be used to generate docs and tests for a resource."
weight: 12
---

# Add an MMv1 test

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

## Example Configuration File

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

## `terraform.yaml` metadata

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

## Results

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

## Update tests

Update tests can only be [added as handwritten tests](/magic-modules/docs/how-to/add-handwritten-test/#update-tests).

## Tests that use beta features

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