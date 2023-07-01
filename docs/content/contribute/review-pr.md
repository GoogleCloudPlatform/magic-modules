---
title: "Review a PR"
weight: 11
---

# Review a PR

This page provides guidelines for reviewing Magic Modules pull requests

1. Read the PR description and the linked issues. Understand the context and make sure that the PR will actually resolve the issues.
2. If the PR adds any new resource, ensure the resource does not already exists in the [GA provider](https://github.com/hashicorp/terraform-provider-google) or [beta provider](https://github.com/hashicorp/terraform-provider-google-beta)
1. Read through all the changes and the API documentation to ensure
   1. the schema added in the PR matches the API structure
   1. features are added in the correct version 
      * features only available in beta are not appear in the GA google provider
      * features added to the GA provider also appear in the beta provider -- beta should be a strict superset of GA 
   1. no [breaking changes]({{< ref "/develop/breaking-changes" >}}) are introduced unless for major release PRs
1. Check the tests to ensure
   1. all fields added/updated in the PR appear in at least one test 
   1. mutable features are tested in at least one update test
   1. all related tests pass in GA for features promoted from beta to GA
   1. newly added or modified diff suppress functions are tested in at least one unit test
1. Check documentation to ensure
   1. documentation is added manually for handwritten resources
1. Check if release notes capture all changes in the PR following guidance [write release notes]({{< ref "release-notes" >}})




