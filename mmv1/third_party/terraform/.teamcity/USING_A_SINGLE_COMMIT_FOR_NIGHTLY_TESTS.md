# How TeamCity decides which commit to use for nightly tests

## Background and problem

Our nightly tests are implemented as a build per package in the provider. A cron schedule triggers builds from each build configuration at 4am UTC and we use locks to prevent builds conflicting with each other. When builds are triggered they enter the build queue and wait for an agent to be available to run the build. Overall the test suite takes approximately 8-10 hours to complete (excluding sweeper builds) and builds are leaving the queue at various times.

Builds in TeamCity use the latest commit from the branch they’re testing at the point that they leave the build queue and start to run. This can result in situations where the build queue contains multiple builds, early-running builds use the latest commit (A) on the main branch, a PR is merged and introduces a subsequent commit (B), and then builds that exit the queue later will run tests using commit B.


This is a problem as our release cut process assumes that all acceptance tests run on a given night are testing the same commit, and that commit is used to cut the release. If a night’s tests span multiple commits the Release Shepherd will need to analyze multiple builds and identify what commits were tested and determine whether tests pass equally for all those commits (and then decide on a single commit to use for the release cut!).

## Solution

To solve this problem we need to direct TeamCity to checkout a particular commit when running nightly tests.

We cannot identify a commit to use for all tests and instruct TeamCity to checkout that specific commit directly from main. There is an open Feature Request tracked on JetBrains’ website for this feature: [Allow VCS root to checkout specified revision instead of the most current version.](https://youtrack.jetbrains.com/issue/TW-11400)

Because of this limitation we need to label a commit with either a tag or a new branch and direct TeamCity to checkout the tag/branch, using them as an indirect way of checking out the commit we labeled.

The solution we've implemented includes:

* A GitHub action in the [google](https://github.com/hashicorp/terraform-provider-google/blob/main/.github/workflows/teamcity-nightly-workflow.yaml) and [google-beta](https://github.com/hashicorp/terraform-provider-google-beta/blob/main/.github/workflows/teamcity-nightly-workflow.yaml) repositories that:
    * Runs at **3am UTC**   
    * Renames the previous day's `nightly-test` branch to `UTC-nightly-test-YYYY-MM-DD`, where the date corresponds to when the base commit was made in UTC.
    * Creates a new `nightly-test` branch using the latest commit on the `main` branch
    * Sweeps up old `UTC-nightly-test-YYYY-MM-DD` branches [over 3 days old](https://github.com/hashicorp/terraform-provider-google/blob/5bce89216324fcf9165ef5fc8d1634e55465282b/.github/workflows/teamcity-nightly-workflow.yaml#L83)
* [Updates to TeamCity](https://github.com/GoogleCloudPlatform/magic-modules/pull/10785) so that any builds triggered by the nightly cron at **4am UTC** check out the `nightly-test` branch 

This diagram shows what happens at 3am UTC when the GHAs run, and how this system helps TeamCity use a consistent commit if there are PRs merged after 4am UTC when test builds are triggered and queued:

<p align="center">
  <img src="https://github.com/user-attachments/assets/76646535-01e9-4ea4-9963-89b58887a318">
</p>

