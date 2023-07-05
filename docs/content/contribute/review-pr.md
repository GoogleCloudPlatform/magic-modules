---
title: "Review a PR"
weight: 11
---

# Review a PR

This page provides guidelines for reviewing Magic Modules pull requests

1. Read the pull request description and the linked issues to understand the context and check if the pull request actually resolves the issues.
2. If the PR adds any new resource, ensure that the resource does not already exists in the [GA provider](https://github.com/hashicorp/terraform-provider-google) or [beta provider](https://github.com/hashicorp/terraform-provider-google-beta)
1. Read through all the changes in the pull request, generated code in the downstreams and the API documentation to ensure that:
   1. the resource schema added in the pull request matches the API structure
   1. the features are added in the correct version 
      * features only available in beta are not included in the GA google provider.
      * features added to the GA provider are also included in the beta provider -- beta should be a strict superset of GA.
   1. no [breaking changes]({{< ref "/develop/breaking-changes" >}}) are introduced without a valid justification.
1. Check the tests to ensure that:
   1. all related tests, including acceptance tests and unit tests, have been completed successfully. 
      {{< hint info >}}Note:
      Some acceptance tests may be skipped in presubmit VCR tests and manual testing is required.
      {{< /hint >}}
   1. all fields added/updated in the pull request appear in at least one test.
   1. all mutable features are tested in at least one update test.
   1. all related tests pass in GA for features promoted from beta to GA.
   1. newly added or modified diff suppress functions are tested in at least one unit test.
1. Check documentation to ensure
   1. resouce-level and field-level documentation are generated correctly for MMv1-based resource
   1. documentation is added manually for handwritten resources.
1. Check if release notes capture all changes in the pull request, following the guidance in [write release notes]({{< ref "release-notes" >}}).