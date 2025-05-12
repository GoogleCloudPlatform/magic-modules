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
  - /develop/test/test
---

# Add resource tests

This page describes how to add tests to a new resource in the `google` or `google-beta` Terraform provider.

The providers have two basic types of tests:

- Unit tests: test specific functions thoroughly. Unit tests do not interact with GCP APIs.
- Acceptance tests (aka VCR tests, or create and update tests): test that resources interact as expected with the APIs. Acceptance tests interact with GCP APIs, but should only test the provider's behavior in constructing the API requests and parsing the responses.

Acceptance tests are also called "VCR tests" because they use [`go-vcr`](https://github.com/dnaeon/go-vcr) to record and play back HTTP requests. This allows tests to run more quickly on PRs because the resources don't actually need to be created, updated, or destroyed by the live API.

For more information about testing, see the [official Terraform documentation](https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests).

## Before you begin

1. Determine whether your resources is using [MMv1 generation or handwritten]({{<ref "/" >}}).
2. If you are not adding tests to an in-progress PR, ensure that your `magic-modules`, `terraform-provider-google`, and `terraform-provider-google-beta` repositories are up to date.
   ```bash
   cd ~/magic-modules
   git checkout main && git clean -f . && git checkout -- . && git pull
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
   git checkout main && git clean -f . && git checkout -- . && git pull
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
   git checkout main && git clean -f . && git checkout -- . && git pull
   ```

## Add unit tests

A unit test verifies functionality that is not related to interactions with the API, such as
[diff suppress functions]({{<ref "/reference/field#diff_suppress_func" >}}),
[validation functions]({{<ref "/reference/field#validation" >}}),
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

## Add a create test

A create test is an **acceptance test** that creates the target resource and immediately destroys it.

> **Note:** All resources should have a "basic" create test, which uses the smallest possible number of fields. Additional create tests can be used to ensure all fields on the resource are used in at least one test.

{{< tabs "create" >}}
{{< tab "MMv1" >}}
1. Add an entry to your `RESOURCE_NAME.yaml` file's `examples`. The fields listed here are the most commonly-used. For a comprehensive reference, see [MMv1 resource reference: `examples` ↗]({{<ref "/reference/resource#examples" >}}).
   ```yaml
   examples:
     # name must correspond to a configuration file that you'll create in the next step.
     # The name should include the product name, resource name, and a basic description
     # of the test. This will be used to generate the test name and the documentation
     # header.
     - name: "PRODUCT_RESOURCE_basic"
       # primary_resource_id will be used for the Terraform resource id in the configuration file.
       primary_resource_id: "example"
       # vars contains key/value pairs of variables to inject into the configuration file.
       # These can be referenced in the configuration file as a key inside `{{$.Vars}}`.
       # All resource IDs (even for resources not under test) should be declared
       # with variables that contain a `-` or `_`; this will ensure that, in tests,
       # the resources are created with a `tf-test` prefix to allow automatic cleanup
       # of dangling resources and a random suffix to avoid name collisions.
       vars:
         network_name: "example-network"
       # test_vars_overrides contains key/value pairs of literal overrides for
       # variables used in tests. This can be used to call functions to
       # generate or determine a variable's value – for example, bootstrapping
       # a shared network for your product to avoid test failures due to limits
       # on the default network.
       test_vars_overrides:
         network_name: 'acctest.BootstrapSharedServiceNetworkingConnection(t, "PRODUCT-RESOURCE-network-config")'
       # Set min_version: beta if the resource is not beta-only and any beta-only fields are being tested.
       min_version: beta
   ```

2. Create a `.tf.tmpl` file in [`mmv1/templates/terraform/examples/`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/examples). The name of the file should match the name of the example created in the previous step. For example, `PRODUCT_RESOURCE_basic.tf.tmpl`.
3. In that file, write the Terraform configuration for your test. This should include all of the required dependencies. For example, `google_compute_subnetwork` has a dependency on `google_compute_network`:
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
4. If the resource or the example is beta-only:
   - Add `provider = google-beta` to every resource in the file.
{{< /tab >}}
{{< tab "Handwritten" >}}
This section assumes you've used the [Add a resource]({{< ref "/develop/add-resource" >}}) guide to create your handwritten resource, and you have a working MMv1 config.

> **Note:** If not, you can create one now, or skip this guide and construct the test by hand. Writing tests by hand can sometimes be a better option if there is a similar test you can copy from.

1. Add the test in MMv1. Repeat for all the create tests you will need.
2. [Generate the beta provider]({{< ref "/develop/generate-providers" >}}).
3. From the beta provider, copy and paste the generated `*_generated_test.go` file into the appropriate service folder inside [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services/) as a new file call `*_test.go`.
4. Modify the tests as needed.
   - Replace all occurrences of `github.com/hashicorp/terraform-provider-google-beta/google-beta` with `github.com/hashicorp/terraform-provider-google/google`
   - Remove the comments at the top of the file.
   - Remove the `Example` suffix from all function names.
   - If beta-only fields are being tested, do the following:
     - Change the file suffix to `.go.tmpl`
     - Wrap each beta-only test in a separate version guard: `{{- if ne $.TargetVersionName "ga" -}}...{{- else }}...{{- end }}`
     - In each beta-only test, ensure that the TestCase sets `ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t)`
     - In each beta-only test, ensure that all Terraform resources in all configs have `provider = google-beta` set
{{< /tab >}}
{{< /tabs >}}

## Add an update test

An update test is an **acceptance test** that creates the target resource and then makes updates to fields that are updatable. Updatable fields are fields that can be updated without recreating the entire resource; that is, they are not marked `immutable` in MMv1 or `ForceNew` in handwritten code.

> **Note:** All updatable fields must be covered by at least one update test. In most cases, only a single update test is needed to test all fields at once.

{{< tabs "update" >}}
{{< tab "MMv1" >}}
1. [Generate the beta provider]({{< ref "/develop/generate-providers" >}}).
2. From the beta provider, copy and paste the generated `*_generated_test.go` file into the appropriate service folder inside [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services) as a new file call `*_test.go`.
3. Using an editor of your choice, delete the `*DestroyProducer` function, and all but one test. The remaining test should be the "full" test, or if there is no "full" test, the "basic" test. This will be the starting point for your new update test.
4. Modify the `TestAcc*` *test function* to support updates.
   - Change the suffix of the test function to `_update`.
   - Copy the 2 `TestStep` blocks and paste them immediately after, so that there are 4 total test steps.
   - Change the suffix of the first `Config` value to `_full` (or `_basic`).
   - Change the suffix of the second `Config` value to `_update`.
   - Add `ConfigPlanChecks` to the update step of the test to ensure the resource is updated in-place.
   - The resulting test function would look similar to this:
   ```go
   import "github.com/hashicorp/terraform-plugin-testing/plancheck"

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
   import "github.com/hashicorp/terraform-plugin-testing/plancheck"

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

## Bootstrapping API resources {#bootstrapping}

Most acceptance tests run in a the default org and default test project, which means that they can conflict for quota, resource namespaces, and control over shared resources. You can work around these limitations with "bootstrapped" resources.

### CryptoKeys

There are a few functions provided for bootstrapping CryptoKeys, depending on your needs.

- `BootstrapKMSKeyWithPurposeInLocationAndName(t *testing.T, purpose, locationID, keyShortName string)`
- `BootstrapKMSKeyWithPurposeInLocation(t *testing.T, purpose, locationID string)`
  - Uses a default key name based on the purpose.
- `BootstrapKMSKeyWithPurpose(t *testing.T, purpose string)`
  - Uses `global` location and a key name based on the purpose.
- `BootstrapKMSKeyInLocation(t *testing.T, locationID string)`
  - Uses `ENCRYPT_DECRYPT` for the purpose and the corresponding key name.
- `BootstrapKMSKey(t *testing.T)`
  - Uses `global` location, `ENCRYPT_DECRYPT` for the purpose, and the corresponding key name for that purpose. 

Example usage:

{{< tabs "bootstrap-cryptokeys" >}}
{{< tab "MMv1" >}}
```yaml
examples:
  - name: service_resource_basic
    primary_resource_id: example
    vars:
      kms_key_name: 'kms-key'
    test_vars_overrides:
      kms_key_name: 'acctest.BootstrapKMSKey(t).CryptoKey.Name'
```
{{< /tab >}}
{{< tab "Handwritten" >}}
```go
func TestAccProductResource_update(t *testing.T) {
   t.Parallel()

   context := map[string]interface{}{
      "kms": acctest.BootstrapKMSKey(t).CryptoKey.Name,
      // other variables
   }
   // rest of test
}
```
{{< /tab >}}
{{< /tabs >}}

### IAM resources

Specify member/role pairs that should always exist. `{project_number}` will be replaced with the default project's project number. `{organization_id}` will be replaced with the "target" test organization's ID – we don't modify IAM in the main test org to avoid accidentally locking ourselves out.

Permissions attached to resources created _in_ a test should instead be provisioned with standard terraform resources.

Example usage:

{{< tabs "bootstrap-iam" >}}
{{< tab "MMv1" >}}
```yaml
# Project-level IAM
examples:
  - name: service_resource_basic
    primary_resource_id: example
    bootstrap_iam:
      - member: "serviceAccount:service-{project_number}@gcp-sa-healthcare.iam.gserviceaccount.com"
        role: "roles/bigquery.dataEditor"
```

```yaml
# Org-level IAM
examples:
  - name: service_resource_basic
    primary_resource_id: example
    bootstrap_iam:
      - member: "serviceAccount:service-org-{organization_id}@gcp-sa-osconfig.iam.gserviceaccount.com"
        role: "roles/osconfig.serviceAgent"
    test_env_vars:
      org_id: ORG_TARGET
```
{{< /tab >}}
{{< tab "Handwritten" >}}
```go
// Project-level IAM
import (
  "github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
)

func TestAccProductResource_update(t *testing.T) {
    t.Parallel()

    acctest.BootstrapIamMembers(t, []acctest.IamMember{
        {
            Member: "serviceAccount:service-{project_number}@gcp-sa-pubsub.iam.gserviceaccount.com",
            Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
        },
    })
    // rest of test
}
```
```go
// Org-level IAM
import (
  "github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
  "github.com/hashicorp/terraform-provider-google-beta/google-beta/envvar"
)

func TestAccProductResource_update(t *testing.T) {
    t.Parallel()

    acctest.BootstrapIamMembers(t, []acctest.IamMember{
        {
            Member: "serviceAccount:service-org-{organization_id}@gcp-sa-osconfig.iam.gserviceaccount.com",
            Role:   "roles/osconfig.serviceAgent",
        },
    })
    context := map[string]string{
        "org_id": envvar.GetTestOrgTargetFromEnv(t),
    }
    // rest of test
}
```
{{< /tab >}}
{{< /tabs >}}

### Networks

Bootstrapping networks can be useful for two reasons:

1. Resources like `google_service_networking_connection` use a consumer network and create a complementing tenant network which we don't control. These tenant networks never get cleaned up and they can accumulate to the point where a limit is reached for the organization. By reusing a consumer network across test runs, we can reduce the number of tenant networks that are needed. (Googlers: See b/146351146 for more context.)
2. Bootstrap networks used in tests (gke clusters, dataproc clusters...) to limit traffic to the default network (preventing conflicts).

When creating a bootstrapped network in a test, you can specify an identifier. Note that if the network is being used for a `google_service_networking_connection`, you should use an identifier unique to the test to avoid race conditions where multiple tests attempt to modify the connection at once.

You can also bootstrap one or more subnetworks within a bootstrapped network if necessary, to avoid subnetwork-level quotas and race conditions.

Example usage:

{{< tabs "bootstrap-networks" >}}
{{< tab "MMv1" >}}
```yaml
examples:
  - name: service_resource_basic
    primary_resource_id: example
    vars:
      network_name: 'default'
      subnetwork_name: 'default'
    test_vars_overrides:
      network_name: 'acctest.BootstrapSharedTestNetwork(t, "network-identifier")'
      subnetwork_name: 'acctest.BootstrapSubnet(t, "subnet-identifier", acctest.BootstrapSharedTestNetwork(t, "network-identifier"))'
```
{{< /tab >}}
{{< tab "Handwritten" >}}
```go
func TestAccProductResource_update(t *testing.T) {
   t.Parallel()

   networkName := 
   subnetName := 
   context := map[string]interface{}{
      "network_name": acctest.BootstrapSharedTestNetwork(t, "network-identifier"),
      "subnetwork_name": acctest.BootstrapSubnet(t, "subnet-identifier", acctest.BootstrapSharedTestNetwork(t, "network-identifier")),
      // other variables
   }
   // rest of test
}
```
{{< /tab >}}
{{< /tabs >}}

## Skip tests in VCR replaying mode {#skip-vcr}

Acceptance tests are run in VCR replaying mode on PRs (using pre-recorded HTTP requests and responses) to reduce the time it takes to present results to contributors. However, not all resources or tests are possible to run in replaying mode. Incompatible tests should be skipped during VCR replaying mode. They will still run in our nightly test suite.

{{< tabs "skipping-tests-in-vcr-replaying" >}}

   {{< tab "Skip a generated test" >}}
   Skipping acceptance tests that are generated from example files can be achieved by adding `skip_vcr: true` in the example's YAML:

   ```yaml
   examples:
   - name: 'bigtable_app_profile_anycluster'
      ...

      # bigtable instance does not use the shared HTTP client, this test creates an instance
      skip_vcr: true
   ```

   If you skip a test in VCR mode, include a code comment explaining the reason for skipping (for example, a link to a GitHub issue.)

   {{< /tab >}}
   {{< tab "Skip a handwritten test" >}}
   Skipping acceptance tests that are handwritten can be achieved by adding `acctest.SkipIfVcr(t)` at the start of the test:

   ```go
   func TestAccPubsubTopic_update(t *testing.T) {
         acctest.SkipIfVcr(t) // See: https://github.com/hashicorp/terraform-provider-google/issues/9999
         acctest.VcrTest(t, resource.TestCase{ ... })
   }
   ```

   If you skip a test in VCR mode, include a code comment explaining the reason for skipping (for example, a link to a GitHub issue.)

   {{< /tab >}}
{{< /tabs >}}

### Reasons that tests are skipped in VCR replaying mode

| Problem                                          | How to fix/Other info  | Skip in VCR replaying? |
| ------------------------------------------------ | ---------------------- |------------- |
| *Incorrect or insufficient data is present in VCR recordings to replay tests*.  Tests will fail with `Requested interaction not found` errors during REPLAYING mode | Make sure that you're not introducing randomness into the test, e.g. by unnecessarily using the random provider to set a resource's name.| If you cannot avoid this issue you should skip the test, but try to ensure that it cannot be fixed first.|
*Bigtable acceptance tests aren't working in VCR mode*. `Requested interaction not found` errors are seen during Bigtable tests run in REPLAYING mode | Currently the provider uses a separate client than the rest of the provider to interact with the Bigtable API. As HTTP traffic to the Bigtable API doesn't go via the shared client it cannot be recorded in RECORDING mode.| Skip the test in VCR for Bigtable. |
| *Using multiple provider aliases doesn't work in VCR*. You may have two instances of the google provider in the test config but one of them doesn't seem to be using its provider arguments - for example, using the wrong default project. | See this GitHub issue: https://github.com/hashicorp/terraform-provider-google/issues/20019 . The problem is that, due to how the VCR system works, one provider instance will be configured and the other will be forced to reuse the first instance's configuration, despite them being given different provider arguments. |  Skip the test in VCR is using aliases is unavoidable. |
| *Using multiple versions of the google/google-beta provider in a single test isn't working in VCR*. Unexpected test failures may occur during tests in REPLAYING mode where `ExternalProviders` is used to pull in past versions of the google/google-beta provider. | When ExternalProviders is used to pulling in other versions of the provider, any HTTP traffic through the external provider will not be recorded. If the HTTP traffic produces an unexpected result or returns an API error then the test will fail in REPLAYING mode. | Skip the test in VCR when testing the current provider behaviour versus previous released versions. |

Some additional things to bear in mind are that VCR tests in REPLAYING mode will still interact with GCP APIs somewhat. For example:

- When the provider is configured it will use credentials to obtain access tokens from GCP
- Some acceptance tests use bootstrapping functions that ensure long-lived resources are present in a testing project before the provider is tested.

These tests can still run in VCR replaying mode; however, REPLAYING mode can't be used as a way to completely avoid HTTP traffic generally or with GCP APIs.


## What's next?

[Run your tests]({{< ref "/test/run-tests" >}})

## References

* [Official Terraform documentation on Acceptance Tests](https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests)
* [MMv1 resource reference: `examples` ↗]({{<ref "/reference/resource#examples" >}})

