---
title: "Contributing"
weight: 50
---

# General contributing steps

1. Fork `Magic Modules` repository into your GitHub account if you haven't done before.
1. Check the [issue tracker](https://github.com/hashicorp/terraform-provider-google/issues) to see whether your feature has already been requested.
   * if there's an issue and it's already has a dedicated assignee, it indicates that someone may have already started to work on a solution.
   * otherwise, you're welcome to work on the issue.
1. Check whether the resource you would like to work on already exists in the providers ([`google`](https://github.com/hashicorp/terraform-provider-google) / [`google-beta`](https://github.com/hashicorp/terraform-provider-google-beta) or [check the website](https://registry.terraform.io/providers/hashicorp/google/latest/docs)).
   * If it exists, check the header of the downstream file to identify the type of tools used to generate the resource. For some resources, the code file, the test file and the documentation file may not be generated via the same tools.
      * Generated resources like [`google_compute_address`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address) can be identified by looking in their [`Go source`](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_compute_address.go) for an `AUTO GENERATED CODE` header as well as a `Type`. "Generated resources" typically refers to just the `MMv1` type, and `DCL` type resources are considered "DCL-based". (Currently DCL-related contribution are not supported)
      * Handwritten resources like [`google_container_cluster`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/container_cluster) can be identified if they have source code present under the [`mmv1/third_party/terraform/resources`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/resources) folder or by the absence of the `AUTO GENERATED CODE header` in their [`Go source`](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_container_cluster.go).
   * If not, decide which tool you would like to use to implement the resource.
      * MMv1 is strongly preferred over handwriting the resource unless the resource can not be generated.
      * Currently, only handwritten datasources are supported.
1. Make the actual code change.
   * The [How To](/magic-modules/docs/how-to) section will guide you to the detailed instructions on how to make your change.
1. [Generate the providers](/magic-modules/docs/getting-started/generate-providers/) that include your change.
1. [Run provider tests locally](/magic-modules/docs/getting-started/run-provider-tests/) that are relevant to the change you made. (Testing the PR locally and pushing the commit to the PR only after the tests pass locally may significantly reduce back-and-forth in review.)
1. Push your changes to your `magic-modules` repo fork and send a pull request from that branch to the main branch on `magic-modules`. A reviewer will be assigned automatically to your PR.
1. Wait until the the modules magician to generate downstream diff (which should take about 15 mins after creating the PR) to make sure all changes are generated correctly in downstream repos.
1. Wait for the VCR test results.
   {{< details "Get to know general workflow for VCR tests" >}}

      1. You submit your change.
      1. The recorded tests are ran against your changes by the `modular-magician`. Tests will fail if:
         1. Your PR has changed the HTTP request values sent by the provider
         1. Your PR does not change the HTTP request values, but fails on the values returned in an old recording
         1. The recordings are out of sync with the merge-base of your PR, and an unrelated contributor's change has caused a false positive
      1. The `modular-magician` will leave a message indicating the number of passing and failing VCR tests. If there is a failure, the `modular-magician` user will leave a message indicating the "`Triggering VCR tests in RECORDING mode for the following tests that failed during VCR:`" marking which tests failed.
         1. If a test does not appear related, it probably isn't!
      1. The `modular-magician` will kick off a second test run targeting only the failed tests, this time hitting the live GCP APIs. If there are tests that fail at this point, a message stating `Tests failed during RECORDING mode:` will be left indicating the tests.
         1. If a test that appears to be related to your change has failed here, it's likely your change has introduced an issue. You can view the debug logs for the test by clicking the "view" link beside the test case to attempt to debug what's going wrong.
         1. If a test appears to be completely unrelated has failed, it's possible that a GCP API has changed in a way that broke the provider or our environment capped on a quota.
   {{< /details >}}

   Where possible, take a look at the logs and see if you can figure out what needs to be fixed related to your change.
   The false positive rate on these tests is extremely high between changes in the API, Cloud Build bugs, and eventual consistency issues in test recordings so we don't expect contributors to wholly interpret the results- that's the responsibility of your reviewer.
1. If your assigned reviewers does not reply / review within a week, gently ping them on github.
1. After your PR is merged, it will be released to customers in around a week or two.
