---
title: "Add resource tests"
weight: 10
aliases:
  - /docs/how-to/add-mmv1-test
  - /how-to/add-mmv1-test
  - /develop/add-mmv1-test
  - /docs/how-to/add-handwritten-test
  - /how-to/add-handwritten-test
  - /develop/add-handwritten-test
  - /develop/test
---

# Add resource tests

This page describes how to add tests to a new resource in the `google` or `google-beta` Terraform provider.

For more information about testing, see the [official Terraform documentation](https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests).

## Before you begin

1. Determine whether your resources is using [MMv1 generation or handwritten]({{<ref "/get-started/how-magic-modules-works.md" >}}).
2. If you are not adding tests to an in-progress PR, ensure that your `magic-modules`, `terraform-provider-google`, and `terraform-provider-google-beta` repositories are up to date.
   ```bash
   cd ~/magic-modules
   git checkout main && git clean -f . && git checkout -- . && git pull
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
   git checkout main && git clean -f . && git checkout -- . && git pull
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
   git checkout main && git clean -f . && git checkout -- . && git pull
   ```

## Add a create test

A create test is a test that creates the target resource and immediately destroys it.

> **Note:** All resources should have a "basic" create test, which uses the smallest possible number of fields. Additional create tests can be used to ensure all fields on the resource are used in at least one test.

{{< tabs "create" >}}
{{< tab "MMv1" >}}
1. Using an editor of your choice, create a `*.tf.tmpl` file in [`mmv1/templates/terraform/examples/`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/examples). The name of the file should include the service name, resource name, and a descriptor. For example, `compute_subnetwork_basic.tf.tmpl`.
2. Write the Terraform configuration for your test. This should include all of the required dependencies. For example, `google_compute_subnetwork` has a dependency on `google_compute_network`:
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
3. If beta-only fields are being tested:
   - Add `provider = google-beta` to every resource in the file.
4. Modify the configuration to use templated values.
   - Replace the id of the primary resource you are testing with `{{$.PrimaryResourceId}}`.
   - Replace fields that are identifiers, like `id` or `name`, with an appropriately named variable. For example, `{{index $.Vars "subnetwork_name"}}`.
   - The resulting configuration for the above example would look like this:
   ```tf
   resource "google_compute_subnetwork" "{{$.PrimaryResourceId}}" {
     name          = "{{index $.Vars "subnetwork_name"}}"
     ip_cidr_range = "10.1.0.0/16"
     region        = "us-central1"
     network       = google_compute_network.network.name
   }

   resource "google_compute_network" "network" {
     name                    = "{{index $.Vars "network_name"}}"
     auto_create_subnetworks = false
   }
   ```
5. Modify the relevant `RESOURCE_NAME.yaml` file under [magic-modules/mmv1/products](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products) to include an [`examples`](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/api/resource/examples.go) block with your test. The `name` must match the name of your `*.tf.tmpl` file. For example:
   ```yaml
   examples:
     - name: "compute_subnetwork_basic"
       primary_resource_id: "example"
       vars:
         subnetwork_name: "example-subnet"
         network_name: "example-network"
   ```
{{< hint warning >}}
**Warning:** Values in `vars` must include a `-` (or `_`). They [trigger the addition of a `tf-test` prefix](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/api/resource/examples.go#L224), which the sweeper uses to clean them up after tests run.
{{< /hint >}}
6. If beta-only fields are being tested:
   - Add `min_version: 'beta'` to the `examples` block in `RESOURCE_NAME.yaml`.
{{< /tab >}}
{{< tab "Handwritten" >}}
This section assumes you've used the [Add a resource]({{< ref "/develop/resource.md" >}}) guide to create your handwritten resource, and you have a working MMv1 config.

> **Note:** If not, you can create one now, or skip this guide and construct the test by hand. Writing tests by hand can sometimes be a better option if there is a similar test you can copy from.

1. Add the test in MMv1. Repeat for all the create tests you will need.
2. [Generate the beta provider]({{< ref "/get-started/generate-providers.md" >}}).
3. From the beta provider, copy and paste the generated `*_generated_test.go` file into the appropriate service folder inside [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services/) as a new file call `*_test.go`.
4. Modify the tests as needed.
   - Replace all occurrences of `github.com/hashicorp/terraform-provider-google-beta/google-beta` with `github.com/hashicorp/terraform-provider-google/google`
   - Remove the comments at the top of the file.
   - Remove the `Example` suffix from all function names.
   - If beta-only fields are being tested, do the following:
     - Change the file suffix to `.go.tmpl`
     - Wrap each beta-only test in a separate version guard: `{{- if ne $.TargetVersionName "ga" -}}...{{- else }}...{{- end }}`
{{< /tab >}}
{{< /tabs >}}

## Add an update test

An update test is a test that creates the target resource and then makes updates to fields that are updatable. Updatable fields are fields that can be updated without recreating the entire resource; that is, they are not marked `immutable` in MMv1 or `ForceNew` in handwritten code.

> **Note:** All updatable fields must be covered by at least one update test. In most cases, only a single update test is needed to test all fields at once.

{{< tabs "update" >}}
{{< tab "MMv1" >}}
1. [Generate the beta provider]({{< ref "/get-started/generate-providers.md" >}}).
2. From the beta provider, copy and paste the generated `*_generated_test.go` file into the appropriate service folder inside [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services) as a new file call `*_test.go`.
3. Using an editor of your choice, delete the `*DestroyProducer` function, and all but one test. The remaining test should be the "full" test, or if there is no "full" test, the "basic" test. This will be the starting point for your new update test.
4. Modify the `TestAcc*` *test function* to support updates.
   - Change the suffix of the test function to `_update`.
   - Copy the 2 `TestStep` blocks and paste them immediately after, so that there are 4 total test steps.
   - Change the suffix of the first `Config` value to `_full` (or `_basic`).
   - Change the suffix of the second `Config` value to `_update`.
   - The resulting test function would look similar to this:
   ```go
   func TestAccPubsubTopic_update(t *testing.T) {
      ...
      acctest.VcrTest(t, resource.TestCase{
         ...
         Steps: []resource.TestStep{
            {
               Config: testAccPubsubTopic_full(...),
            },
            {
               ...
            },
            {
               Config: testAccPubsubTopic_update(...),
            },
            {
               ...
            },
         },
      })
   }
   ```
5. Modify the `testAcc*` Terraform *template function* to support updates.
   - Copy the template function and paste it immediately after so that there are 2 template functions.
   - Change the suffix of the first template function to `_full` (or `_basic`).
   - Change the suffix of the second template function to `_update`.
   - The resulting template functions would look similar to this:
   ```go
   func testAccPubsubTopic_full(...) string {
       ...
   }

   func testAccPubsubTopic_update(...) string {
       ...
   }
   ```
6. Modify the test as needed.
   - Replace all occurrences of `github.com/hashicorp/terraform-provider-google-beta/google-beta` with `github.com/hashicorp/terraform-provider-google/google`
   - Modify the template function ending in `_update` so that updatable fields are changed or removed. This may require additions to the `context` map in the test function.
   - Remove the comments at the top of the file.
   - If beta-only fields are being tested, do the following:
     - Change the file suffix to `.go.tmpl`
     - Wrap each beta-only test in a separate version guard: `{{- if ne $.TargetVersionName "ga" -}}...{{- else }}...{{- end }}`
     - In each beta-only test, ensure that the TestCase sets `ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t)`
     - In each beta-only test, ensure that all Terraform resources in all configs have `provider = google-beta` set
{{< /tab >}}
{{< tab "Handwritten" >}}
1. Using an editor of your choice, open the existing `*_test.go` or `*_test.go.tmpl` file in the appropriate service folder inside [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services) which contains your create tests.
2. Copy the `TestAcc*` *test function* for the existing "full" test. If there is no "full" test, use the "basic" test. This will be the starting point for your new update test.
3. Modify the test function to support updates.
   - Change the suffix of the test function to `_update`.
   - Copy the 2 `TestStep` blocks and paste them immediately after, so that there are 4 total test steps.
   - Change the suffix of the second `Config` value to `_update`.
   - Add `ConfigPlanChecks` to the update step of the test to ensure the resource is updated in-place.
   - The resulting test function would look similar to this:
   ```go
   func TestAccPubsubTopic_update(t *testing.T) {
      ...
      acctest.VcrTest(t, resource.TestCase{
         ...
         Steps: []resource.TestStep{
            {
               Config: testAccPubsubTopic_full(...),
            },
            {
               ...
            },
            {
               Config: testAccPubsubTopic_update(...),
               ConfigPlanChecks: resource.ConfigPlanChecks{
                  PreApply: []plancheck.PlanCheck{
                     plancheck.ExpectResourceAction("google_pubsub_topic.foo", plancheck.ResourceActionUpdate),
                  },
               },
            },
            {
               ...
            },
         },
      })
   }
   ```
4. Add a Terraform *template function* to support updates.
   - Copy the full (or basic) `testAcc*` template function.
   - Change the suffix of the new template function to `_update`.
   - The new template function would look similar to this:
   ```go
   func testAccPubsubTopic_update(...) string {
       ...
   }
   ```
5. Modify the test as needed.
   - Modify the new template function so that updatable fields are changed or removed. This may require additions to the `context` map in the test function.
   - Remove the comments at the top of the file.
   - If beta-only fields are being tested, do the following:
     - Change the file suffix to `.go.tmpl`
     - Wrap each beta-only test in a separate version guard: `{{- if ne $.TargetVersionName "ga" -}}...{{- else }}...{{- end }}`
     - In each beta-only test, ensure that the TestCase sets `ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t)`
     - In each beta-only test, ensure that all Terraform resources in all configs have `provider = google-beta` set
{{< /tab >}}
{{< /tabs >}}

## Add unit tests

A unit test verifies functionality that is not related to interactions with the API, such as
[diff suppress functions]({{<ref "/develop/field-reference#diff_suppress_func" >}}),
[validation functions]({{<ref "/develop/field-reference#validation" >}}),
CustomizeDiff functions, and so on.

Unit tests should be added to the appropriate folder in [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services) in the file called `resource_PRODUCT_RESOURCE_test.go`. (You may need to create this file if it does not already exist. Replace PRODUCT with the product name and RESOURCE with the resource name; it should match the name of the generated resource file.)

Unit tests should be named like `TestFunctionName` - for example, `TestDiskImageDiffSuppress` would contain tests for the `DiskImageDiffSuppress` function.

Example:

```go
func TestSignatureAlgorithmDiffSuppress(t *testing.T) {
   cases := map[string]struct {
      Old, New           string
      ExpectDiffSuppress bool
   }{
      "ECDSA_P256 equivalent": {
         Old:                "ECDSA_P256_SHA256",
         New:                "EC_SIGN_P256_SHA256",
         ExpectDiffSuppress: true,
      },
      // Additional cases excluded for brevity
   }

   for tn, tc := range cases {
      if signatureAlgorithmDiffSuppress("signature_algorithm", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
         t.Errorf("bad: %s, %q => %q expect DiffSuppress to return %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
      }
   }
}
```

## What's next?

- [Run your tests]({{< ref "/develop/test/run-tests.md" >}})
