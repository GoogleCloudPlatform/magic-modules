---
title: "Contributing"
weight: 50
aliases:
  - /docs/getting-started/contributing
  - /getting-started/contributing
---

# General contributing steps

1. If you haven't done so already, fork the `Magic Modules` repository into your GitHub account.
1. Check the [issue tracker](https://github.com/hashicorp/terraform-provider-google/issues) to see whether your feature has already been requested.
   * If there's an issue and it already has a dedicated assignee, this indicates that someone might have already started to work on a solution.
   * Otherwise, you're welcome to work on the issue.
1. Check whether the resource you would like to work on already exists in the following places:
   * [`google` provider](https://github.com/hashicorp/terraform-provider-google) providers
   * [`google-beta` provider](https://github.com/hashicorp/terraform-provider-google-beta)
   * [Hashicorp documentation](https://registry.terraform.io/providers/hashicorp/google/latest/docs)
   
   If it exists, check the header of the downstream file to identify the type of tools used to generate the resource.
   For some resources, the code file, the test file and the documentation file might not be generated via the same tools.
      * Generated resources like [`google_compute_address`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address) can be identified by looking in their [`Go source`](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_compute_address.go) for an `AUTO GENERATED CODE` header as well as a `Type`. "Generated resources" typically refers to just the `MMv1` type, and `DCL` type resources are considered "DCL-based". (Currently DCL-related contribution are not supported)
      * Handwritten resources like [`google_container_cluster`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/container_cluster) can be identified if they have source code present under the [`mmv1/third_party/terraform/resources`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/resources) folder or by the absence of the `AUTO GENERATED CODE header` in their [`Go source`](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_container_cluster.go).
   
   If not, decide which tool you would like to use to implement the resource:
      * MMv1 is strongly preferred over handwriting the resource unless the resource cannot be generated.
      * Currently, only handwritten datasources are supported.
1. Make the code change.
   * The Develop section provides detailed instructions on how to make your change.
1. [Generate the providers]({{< ref "/get-started/generate-providers" >}}) that include your change.
1. [Run provider tests locally]({{< ref "/develop/run-tests" >}}) that are relevant to the change you made. (Testing the PR locally and pushing the commit to the PR only after the tests pass locally may significantly reduce back-and-forth in review.)
1. Push your changes to your `magic-modules` repo fork and send a pull request from that branch to the main branch on `magic-modules`. A reviewer will be assigned automatically to your PR.
1. Get approval to start Clould Builder jobs from the reviewer if you're an community contributor
1. Wait for the the modules magician to generate downstream diff (which should take about 15 mins after creating the PR) to make sure all changes are generated correctly in downstream repos.
1. Wait for the VCR test results.
{{< details "Get to know general workflow for VCR tests" >}}
   1. Submit your change.
   1. The recorded tests are ran against your changes by the `modular-magician`. Tests will fail if:
      * Your PR has changed the HTTP request values sent by the provider
      * Your PR does not change the HTTP request values, but fails on the values returned in an old recording
      * The recordings are out of sync with the merge-base of your PR, and an unrelated contributor's change has caused a false positive
   1. The `modular-magician` will leave a message indicating the number of passing and failing VCR tests. If there is a failure, the `modular-magician` user will leave a message indicating the "`Triggering VCR tests in RECORDING mode for the following tests that failed during VCR:`" marking which tests failed.
      * If a test does not appear related to your PR, it probably isn't!
   1. The `modular-magician` will kick off a second test run targeting only the failed tests, this time hitting the live GCP APIs. If there are tests that fail at this point, a message stating `Tests failed during RECORDING mode:` will be left indicating the tests.
      * If a test that appears to be related to your change has failed here, it's likely your change has introduced an issue. You can view the debug logs for the test by clicking the "view" link beside the test case to attempt to debug what's going wrong.
      * If a test that appears to be completely unrelated has failed, it's possible that a GCP API has changed in a way that broke the provider or our environment capped on a quota.
{{< /details >}}

   Where possible, take a look at the logs and see if you can figure out what needs to be fixed related to your change.
   The false positive rate on these tests is extremely high between changes in the API, Cloud Build bugs, and eventual consistency issues in test recordings so we don't expect contributors to wholly interpret the results — that's the responsibility of your reviewer.
1. If your assigned reviewer does not respond to changes on a pull request within two US business days, ping them on the pull request.
1. After your PR is merged, it will be released to customers in around one to two weeks.
