Concourse CI tools for MagicModules and Google Providers
===

These tools manage the downstream repositories of [magic-modules](https://github.com/GoogleCloudPlatform/magic-modules).

# Jobs
The concourse pipeline defined here runs through four stages when a Github pull request is opened against `magic-modules`.
## `mm-generate`
The first job's goal is to generate the downstream repositories of MagicModules.  This is done almost entirely by the `generate` task.

### `generate`
`generate` takes two inputs:
* The magic-modules repository, with the pull request's `head` checked out and the repo in "detached `HEAD`" state.
* This CI repository.

It then runs `magic-modules/generate.yml`, which specifies one output:

* The magic-modules repository, after generation has been accomplished.

After that, the generated repositories are uploaded to GitHub, in the concourse process runner's forks.  The pull request is updated to point to those forks.

## Individual repo tests

The second stage's goal is to confirm that the individual repos still pass tests.  It runs only after the first stage finishes.

### `terraform-test`

`terraform-test` takes two inputs:
* The updated `magic-modules` repo, after the robot adds its commit to point the submodules to the new generated version.
* This CI repository.

It then runs `unit-tests/task.yml`, which has no outputs because it makes no changes to any code.  It just runs and succeeds if the unit tests pass, and fails if the unit tests don't pass.  The failure and detailed logs are available in concourse.

## `create-prs`

The third stage's goal is to create the downstream PRs.

It takes three inputs:
* The CI repository
* A copy of `magic-modules` which has passed the test stage
* The initial PR

It simply runs `magic-modules/create-pr.yml`, which creates the downstream pull requests using the `hub` CLI.  When complete, it updates `magic-modules` to point to the new submodule commits (for added clarity on the pull request), comments on that PR with a list of downstream PRs, and marks it as successfully generated.  If it fails, it marks generation as failed.

## `merge-prs`

The fourth stage's goal is to merge `magic-modules` PRs after they have been approved and all the downstream PRs have been merged.

### `merge-and-update`

`merge-and-update` takes two inputs:
* The approved PR, after all downstream PRs have been merged.
* This CI repository.

It then runs `magic-modules/merge.yml`, which declares one output:

* The `magic-modules` PR repo after it has been updated to be ready to merge.

#### A note on tracking submodules
This job sets the submodules back to tracking their downstream repositories on `master`.  It updates the submodules to point to the most recent commit on `master`.  This may not be ideal - other commits may have been made to `master` since the downstream PR was merged - however, this is the best way to ensure that we do not go backwards.  Imagine the following situation:

`magic-modules` PR #1 is created.  `downstream-repo-a` PR #7 and `downstream-repo-b` PR #8 are created from that PR.  `magic-modules` PR #2 is created.  `downstream-repo-a` PR #9 is created from that PR.  PR #7 is merged, then PR #9.  Since all of PR #2's downstream PRs are merged, PR #2 is merged, and it includes the changes from PR #1.  PR #8 is finally merged, and so PR #1 is ready to be merged.

When merging PR #1, if we update `downstream-repo-a` to the merge commit for PR #7, we will go backwards, erasing PR #2.  If we update to `master`, we will definitely include both the changes of both PR #1 and PR #2.
