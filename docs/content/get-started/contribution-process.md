---
title: "Contribution process"
weight: 50
aliases:
  - /docs/getting-started/contributing
  - /getting-started/contributing
  - /get-started/contributing
---

# Contribution process

## Before you begin

1. Familiarize yourself with [GitHub flow](https://docs.github.com/en/get-started/quickstart/github-flow)
1. [Fork](https://docs.github.com/en/get-started/quickstart/fork-a-repo) the `Magic Modules` repository into your GitHub account
1. [Set up your development environment](https://googlecloudplatform.github.io/magic-modules/get-started/generate-providers/)
1. Check whether the feature you want to work on has already been [requested in the issue tracker](https://github.com/hashicorp/terraform-provider-google/issues).
   - If there's an issue and it already has a dedicated assignee, this indicates that someone might have already started to work on a solution. Otherwise, you're welcome to work on the issue.

## Contribute code

1. [Create a new branch for your change](https://docs.github.com/en/get-started/quickstart/github-flow#create-a-branch)
1. Make the code change. For example:
   - [Add or modify a resource]({{< ref "/develop/resource" >}})
   - [Add resource tests]({{< ref "/develop/test" >}})
   - [Add a datasource]({{< ref "/develop/add-handwritten-datasource" >}})
   - [Promote to GA]({{< ref "/develop/promote-to-ga" >}})
   - [Make a breaking change]({{< ref "/develop/make-a-breaking-change" >}})
1. [Generate the providers]({{< ref "/get-started/generate-providers" >}}) that include your change.
1. [Run provider tests locally]({{< ref "/develop/run-tests" >}}) that are relevant to the change you made
1. [Create a pull request (PR)](https://docs.github.com/en/get-started/quickstart/github-flow#create-a-pull-request)
   - A reviewer will be automatically assigned for your PR.
   - Tests for community contributors will only run after approval from a reviewer
   - After tests start, downstream diff generation takes about 10 minutes; VCR tests can take up to 2 hours.
   - Make sure your PR body includes the test `Fixes GITHUB_ISSUE_LINK.` once per issue resolved by your PR. Replace `GITHUB_ISSUE_LINK` with a link to a GitHub issue from the [provider issue tracker](https://github.com/hashicorp/terraform-provider-google/issues).

## Code review

{{< hint info >}}
**TIP:** Speeding up review:
1. Make sure your PR only includes one self-contained change. For example, if you are adding two resources, create one PR for each resource.
1. [Run provider tests locally]({{< ref "/develop/run-tests" >}}) that are relevant to the change you made
1. [Self-review your PR]({{< ref "/contribute/review-pr" >}}") or ask someone you know to review
   - Try to resolve test failures where possible, and ask for help if you get stuck.
{{< /hint >}}

If your assigned reviewer does not respond to changes on a pull request within two US business days, ping them on the pull request.

## After your change is merged

After your change is merged, it can take a week or longer to be released to customers.
