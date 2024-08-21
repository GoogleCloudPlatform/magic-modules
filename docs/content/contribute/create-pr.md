---
title: "Create a pull request"
weight: 10
---

# Create a pull request (PR)

## Requirements

1. Make sure your [branch](https://docs.github.com/en/get-started/quickstart/github-flow#create-a-branch) contains a single self-contained change. For example:
	 - If you are adding multiple resources to the provider, only put one resource in each PR - even if the product requires all resources to be present before it can be meaningfully used.
	 - If you are adding a few fields and also fixing a bug, create one PR for adding the new fields and a separate PR for the bugs.
1. Follow the instructions at [Creating a pull request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request) to create a pull request to merge your branch into `GoogleCloudPlatform/magic-modules`.
   - Make sure the PR body includes the text `Fixes GITHUB_ISSUE_LINK.` once per issue resolved by your PR. Replace `GITHUB_ISSUE_LINK` with a link to a GitHub issue from the [provider issue tracker](https://github.com/hashicorp/terraform-provider-google/issues).
   - [Write release notes]({{< ref "/contribute/release-notes" >}})

## Code review

1. A reviewer will automatically be assigned to your PR.
1. Creating a new pull request or pushing a new commit automatically triggers our CI pipelines and workflows. After CI starts, downstream diff generation takes about 10 minutes; VCR tests can take up to 2 hours. If you are a community contributor, some tests will only run after approval from a reviewer.
   - While convenient, relying on CI to test iterative changes to PRs often adds extreme latency to reviews if there are errors in test configurations or at runtime. We **strongly** recommend you [test your changes locally before pushing]({{< ref "/develop/test/run-tests" >}}) even after the initial change.
1. If your assigned reviewer does not respond to changes on a pull request within two US business days, ping them on the pull request.

{{< hint info >}}
**TIP:** Speeding up review:
1. [Test your changes locally before pushing]({{< ref "/develop/test/run-tests" >}}) to iterate faster.
   - You can push them and test in parallel as well. New CI runs will preempt old ones where possible.
1. Resolve failed [status checks](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/collaborating-on-repositories-with-code-quality-features/about-status-checks) quickly
   - Directly ask your reviewer for help if you don't know how to proceed. If there are failed checks they may only check in if there's no progress after a couple days.
1. [Self-review your PR]({{< ref "/contribute/review-pr" >}}) or ask someone else familiar with Terraform to review
{{< /hint >}}


## Troubleshoot status check failures

### Provider unit tests or VCR tests {#provider-test-failures}

VCR test failures that do not immediately seem related to your PR are most likely safe to ignore unless your reviewer says otherwise.

1. Review the "diff generation" report to make sure the generated code looks as expected.
1. Check out the generated code for your PR to [run tests]({{< ref "/develop/test/run-tests" >}}) and iterate locally. For handwritten code or [custom code]({{< ref "/develop/custom-code" >}}), you can iterate directly in the provider and then copy the changes to your `magic-modules` branch once you have resolved the issue.
   {{< tabs "checkout-auto-pr" >}}
   {{< tab "terraform-provider-google" >}}
   ```bash
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
   git checkout -- . && git clean -f google/ google-beta/ website/
   git remote add modular-magician https://github.com/modular-magician/terraform-provider-google.git
   git fetch modular-magician
   git checkout modular-magician/auto-pr-PR_NUMBER
   make test
   make lint
   make testacc TEST=./google/services/container TESTARGS='-run=TestAccContainerNodePool'
   ```
   Replace PR_NUMBER with your PR's ID.
   {{< /tab >}}
   {{< tab "terraform-provider-google-beta" >}}
   ```bash
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
   git checkout -- . && git clean -f google/ google-beta/ website/
   git remote add modular-magician https://github.com/modular-magician/terraform-provider-google-beta.git
   git fetch modular-magician
   git checkout modular-magician/auto-pr-PR_NUMBER
   make test
   make lint
   make testacc TEST=./google/services/container TESTARGS='-run=TestAccContainerNodePool'
   ```
   Replace PR_NUMBER with your PR's ID.
   {{< /tab >}}
   {{< /tabs >}}
