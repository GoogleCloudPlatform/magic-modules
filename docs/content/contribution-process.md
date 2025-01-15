---
title: "Contribution process"
weight: 20
aliases:
  - /docs/getting-started/contributing
  - /getting-started/contributing
  - /get-started/contributing
  - /get-started/contribution-process
---

# Contribution process

This page explains how you can contribute code and documentation to the	`magic-modules` repository.

## Before you begin

1. Familiarize yourself with [GitHub flow](https://docs.github.com/en/get-started/quickstart/github-flow)
1. [Fork](https://docs.github.com/en/get-started/quickstart/fork-a-repo) the `Magic Modules` repository into your GitHub account
1. [Set up your development environment]({{< ref "/develop/set-up-dev-environment/" >}})
1. Check whether the feature you want to work on has already been [requested in the issue tracker](https://github.com/hashicorp/terraform-provider-google/issues).
   - If there's an issue and it already has a dedicated assignee, this indicates that someone might have already started to work on a solution. Otherwise, you're welcome to work on the issue.

## Contribute code

1. [Set up your development environment]({{< ref "/develop/set-up-dev-environment" >}})
1. [Create a new branch for your change](https://docs.github.com/en/get-started/quickstart/github-flow#create-a-branch)
1. Make the code change. For example:
   + [Add a resource]({{< ref "/develop/add-resource" >}})
   + [Add resource tests]({{< ref "/test/test" >}})
   + [Add a datasource]({{< ref "/develop/add-handwritten-datasource" >}})
   + [Promote to GA]({{< ref "/develop/promote-to-ga" >}})
   + [Make a breaking change]({{< ref "/breaking-changes/make-a-breaking-change" >}})
1. [Generate the providers]({{< ref "/develop/generate-providers" >}}) that include your change.
1. [Run provider tests locally]({{< ref "/test/run-tests" >}}) that are relevant to the change you made
1. [Create a pull request (PR)]({{< ref "/code-review/create-pr" >}})
1. Make changes in response to [code review]({{< ref "/code-review/create-pr#code-review" >}})

## After your change is merged

After your change is merged, it can take a week or longer to be released to customers.
