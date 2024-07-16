---
title: "Review a pull request"
weight: 30
---

# Review a pull request

This page provides guidelines for reviewing a Magic Modules pull request (PR).

1. Read the PR description to understand the context and ensure the PR either
   * is linked to a GitHub issue or an internal bug
      * if not, check the [issue tracker](https://github.com/hashicorp/terraform-provider-google/issues) to see whether the feature has already been requested and add the issues in the description, if any.
   * establishes clear context itself via title or description.
2. If the PR adds any new resource, ensure that the resource does not already exist in the [GA provider](https://github.com/hashicorp/terraform-provider-google) or [beta provider](https://github.com/hashicorp/terraform-provider-google-beta)
1. Read through all the changes in the PR, generated code in the downstreams and the API documentation to ensure that:
   1. the resource schema added in the PR matches the API structure.
   1. the features are added in the correct version
      * features only available in beta are not included in the GA google provider.
      * features added to the GA provider are also included in the beta provider -- beta should be a strict superset of GA.
   1. no [breaking changes]({{< ref "/develop/breaking-changes/make-a-breaking-change" >}}) are introduced without a valid justification. Add the `override-breaking-change` label if there is a valid justification.
      * remember to check for changes in default behaviour like changing the flags on delete! 
   1. verify the change **fully** resolves the linked issues, if any. If it does not, change the "Fixes" message to "Part of".
1. Check the tests added/modified to ensure that:
   1. all fields added/updated in the PR appear in at least one test.
      * It is advisable to test updating from a non-zero value to a zero value if feasible.
   1. all mutable fields are tested in at least one update test.
   1. all resources in the acceptance tests have a `tf-test` or `tf_test` prefix in their primary id field.
   1. all handwritten test Config steps include import steps following them
   1. all related tests pass in GA for features promoted from beta to GA.
      {{< hint info >}}Note:
      Presubmit VCR tests do not run in GA. Manual testing is required for promoted GA features.
      {{< /hint >}}
   1. newly added or modified diff suppress functions are tested in at least one unit test.
   1. the linked issue (if any) is covered by at least one test that reproduces the issue
      * for example - a bugfix should test the bug (or explain why it's not feasible to do so in the description, including manual results when possible) and an enhancement should test the new behaviour(s).
   1. all related PR presubmit tests have been completed successfully, including:
      * terraform-provider-breaking-change-test
      * presubmit-rake-tests
      * terraform-provider-google-build-and-unit-tests
      * terraform-provider-google-beta-build-and-unit-tests
      * VCR-test
      {{< hint info >}}Note:
      Some acceptance tests may be skipped in VCR and manual testing is required.
      {{< /hint >}}
   1. a significant number of preexisting tests have not been modified. Changing old tests often indicates a change is backwards incompatible.
1. Check documentation to ensure
   1. resouce-level and field-level documentation are generated correctly for MMv1-based resource
   1. documentation is added manually for handwritten resources.   
1. Check if release notes capture all changes in the PR, and are correctly formatted following the guidance in [write release notes]({{< ref "release-notes" >}}) before merging the PR.
