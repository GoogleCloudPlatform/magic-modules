---
title: "Add resource tests"
weight: 40
aliases:
  - /docs/how-to/add-mmv1-test
  - /how-to/add-mmv1-test
  - /develop/add-mmv1-test
  - /docs/how-to/add-handwritten-test
  - /how-to/add-handwritten-test
  - /develop/add-handwritten-test
---

# Add resource tests

This page describes how to add tests to a new resource in the `google` or `google-beta` Terraform provider.

For more information about testing, see the [official Terraform documentation](https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests).

## Before you begin

1. Determine whether your resources is using MMv1 generation or handwritten.
2. Ensure that your `magic-modules`, `terraform-provider-google`, and `terraform-provider-google-beta` repositories are up to date.
   ```
   cd ~/magic-modules
   git checkout main && git clean -f . && git checkout -- . && git pull
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
   git checkout main && git clean -f . && git checkout -- . && git pull
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
   git checkout main && git clean -f . && git checkout -- . && git pull
   ```

## Add a standard test

In this section, you will start by creating a "basic" test, which means it tests the simplest possible configuration of your resource (often only required fields). After the basic test is created, you will then optionally create a "full" test, which means it tests a configuration with all possible fields. As an alternative to a single full test, you can create multiple tests based on advanced use-cases.

> **Note:** All fields must be covered by at least one test.

{{< tabs "standard" >}}
{{< tab "MMv1" >}}
1. Using an editor of your choice, create a `*.tf.erb` file in [`mmv1/templates/terraform/examples/`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/examples). The name of the file should include the service name, resource name, and a descriptor. For example, `compute_subnetwork_basic.tf.erb`.
2. Write the Terraform configuration for your test. This will need to include all of the required dependencies. For example:
   ```tf
   resource "google_compute_subnetwork" "primary" {
     name          = "my-subnet"
     ip_cidr_range = "10.1.0.0/16"
     region        = "us-central1"
     network       = google_compute_network.network.name
   }

   resource "google_compute_network" "network" {
     name                    = "my-network"
     auto_create_subnetworks = false
   }
   ```
3. Replace the id of the primary resource you are testing with `<%= ctx[:primary_resource_id] %>`.
4. Replace fields that are identifiers, like `id` or `name`, with an appropriately named variable. For example, `<%= ctx[:vars]['subnetwork_name'] %>`. The resulting configuration for the above example would look like this:
   ```tf
   resource "google_compute_subnetwork" "<%= ctx[:primary_resource_id] %>" {
     name          = "<%= ctx[:vars]['subnetwork_name'] %>"
     ip_cidr_range = "10.1.0.0/16"
     region        = "us-central1"
     network       = google_compute_network.network.name
   }

   resource "google_compute_network" "network" {
     name                    = "<%= ctx[:vars]['network_name'] %>"
     auto_create_subnetworks = false
   }
   ```
5. Modify the relevant `RESOURCE_NAME.yaml` file under [magic-modules/mmv1/products](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products) to include an [`examples`](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/provider/terraform/examples.rb) block with your test. For example:
   ```yaml
   examples:
     - !ruby/object:Provider::Terraform::Examples
       name: "compute_subnetwork_basic"
       primary_resource_id: "example"
       vars:
         subnetwork_name: "example-subnet"
         network_name: "example-network"
   ```

> **Warning:** Values in `vars` must include a `-` (or `_`). They [trigger the addition of a `tf-test` prefix](https://github.com/GoogleCloudPlatform/magic-modules/blob/6858338f013f5dc57729ec037883a3594441ea62/mmv1/provider/terraform/examples.rb#L244), which the sweeper uses to clean them up after tests run.
6. If beta-only fields are being tested:
   - Add `provider = google-beta` to every resource in the `*.tf.erb` config file.
   - Add `min_version: beta` to the `examples` block in `RESOURCE_NAME.yaml`.
{{< /tab >}}
{{< tab "Handwritten" >}}
This section assumes you've used the [Add a resource]({{< ref "/develop/resource.md" >}}) guide to create your handwritten resource, and you have a working MMv1 config.

> **Note:** If not, you can create one now, or skip this guide and construct the test by hand. Writing tests by hand can sometimes be a better option if there is a similar test you can copy from.

1. Add the test in MMv1. Repeat for all the standard tests you will need.
2. [Generate the providers]({{< ref "/get-started/generate-providers.md" >}}).
3. From the provider, copy and paste the generated `*_generated_test.go` file into [`magic-modules/mmv1/third_party/terraform/tests`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/tests) as a new file call `*_test.go`.
4. Using an editor of your choice, remove the `Example` suffix from all function names.
5. Modify the tests as needed.
   - Remove the comments at the top of the file.
   - If beta-only fields are being tested, do the following:
     - Change the file suffix to `.go.erb`.
     - Add `<% autogen_exception -%>` to the top of the file.
     - Wrap the beta-only tests in a version guard: `<% unless version = 'ga' -%>...<% else -%>...<% end -%>`.
{{< /tab >}}
{{< /tabs >}}

## Add an update test

In this section, you will create an update test, which will make sure that updatable fields can be properly updated. Updatable fields are fields that can be updated without recreating the entire resource, ie. they are not marked `immutable`. In most cases, only a single update test is needed to test all fields at once.

> **Note:** All updatable fields must be covered by at least one update test.

{{< tabs "update" >}}
{{< tab "MMv1" >}}
1. [Generate the providers]({{< ref "/get-started/generate-providers.md" >}}).
2. From the provider, copy and paste the generated `*_generated_test.go` file into [`magic-modules/mmv1/third_party/terraform/tests`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/tests) as a new file call `*_test.go`.
3. Using an editor of your choice, delete the `*DestroyProducer` function, and all but one test. The remaining test should be the "full" test, or if there is no "full" test, the "basic" test. This will be the starting point for your new update test.
4. Modify the `TestAcc*` test function to support updates.
   - Change the suffix of `TestAcc*` to `_update`.
   - Copy the 2 `TestStep` blocks and paste them immediately after, so that there are 4 total test steps.
   - Change the suffix of the first `Config` value to `_full` (or `_basic`).
   - Change the suffix of the second `Config` value to `_update`.
5. Modify the `testAcc*` Terraform template function to support updates.
   - Copy the `testAcc*` template function and paste it immediately after so that there are 2 template functions.
   - Change the suffix of the first `testAcc*` function to `_full` (or `_basic`).
   - Change the suffix of the second `testAcc*` function to `_update`.
6. Modify the test as needed.
   - Modify the `testAcc*_update` Terraform template function so that updatable fields are changed or removed. This may require additions to the `context` map in `TestAcc*_update`.
   - Remove the comments at the top of the file.
   - If beta-only fields are being tested, do the following:
     - Change the file suffix to `.go.erb`.
     - Add `<% autogen_exception -%>` to the top of the file.
     - Wrap the beta-only tests in a version guard: `<% unless version = 'ga' -%>...<% else -%>...<% end -%>`.
{{< /tab >}}
{{< tab "Handwritten" >}}
1. Using an editor of your choice, open the existing `*_test.go` or `*_test.go.erb` file in [`magic-modules/mmv1/third_party/terraform/tests`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/tests) which contains your standard tests.
2. Copy the `TestAcc*` test function for the existing "full" test. If there is no "full" test, use the "basic" test. This will be the starting point for your new update test.
3. Modify the `TestAcc*` test function to support updates.
   - Change the suffix of `TestAcc*` to `_update`.
   - Copy the 2 `TestStep` blocks and paste them immediately after, so that there are 4 total test steps.
   - Change the suffix of the second `Config` value to `_update`.
4. Add a `testAcc*` Terraform template function to support updates.
   - Copy the full (or basic) `testAcc*` template function.
   - Change the suffix of the new `testAcc*` function to `_update`.
5. Modify the test as needed.
   - Modify the new `testAcc*` Terraform template function so that updatable fields are changed or removed. This may require additions to the `context` map in `TestAcc*_update`.
   - Remove the comments at the top of the file.
   - If beta-only fields are being tested, do the following:
     - Change the file suffix to `.go.erb`.
     - Add `<% autogen_exception -%>` to the top of the file.
     - Wrap the beta-only tests in a version guard: `<% unless version = 'ga' -%>...<% else -%>...<% end -%>`.
{{< /tab >}}
{{< /tabs >}}

## What's next?

- [Test your changes]({{< ref "/get-started/run-provider-tests.md" >}})
