GCB CI tools for Magic Modules and Google Providers
===

These tools manage the downstream repositories of [magic-modules](https://github.com/GoogleCloudPlatform/magic-modules), and are collectively referred to as "The Magician".

# CI For Downstream/Magic Modules Developers
If you're interested in developing the repositories that Magic Modules manages, here are the things you'll want to know.

## What The Magician Does
The Magician takes the PR you write against Magic Modules and creates the downstream (generated) commits that are a result of your changes.  It posts the diffs created by your changes as a comment on your PR, to aid in review.  When your PR is merged, it updates the downstream repositories with your changes.

## Your Workflow
You'll write some code and open a PR.  The Magician will run the generator and the downstream tests.  If either the generator or the tests fail, you'll see "check failed" status on your PR, which will link to the step that failed.  Once all the generation steps succeed, the Magician will comment on your PR with links to review the generated diffs.  Your reviewer will review those diffs (as well as your code).  Once your PR is approved, simply merge it in - the Magician will ensure that the downstream repositories are updated within about 5 minutes.

# CI For CI Developers
If you're working on enhancing these CI tools, you'll need to know about the structure of the pipeline and about how to develop and test.

## How the pipeline works

### Generation & Diffing
The generation / diff pipeline generates the downstream at the PR's merge commit and at the left-side parent (`HEAD~`), which guarantees isolation of exclusively the changes made by the PR in question.  It creates commits for the downstreams as before-and-after commits.  Those commits don't have relevant git history - they're not meant to be applied - but they are used in a [two-dot diff](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/about-comparing-branches-in-pull-requests#three-dot-and-two-dot-git-diff-comparisons).  This means that if there are no changes to the downstream, you'll see an empty diff.  These downstream branches are named `auto-pr-$number-old` and `auto-pr-$number` respectively.

### Downstream Pushing
The Magician maintains a set of tags called `$REPO-sync` that tracks the Magic Modules commit the downstreams are up to date with.

In effect, this means that each downstream commit will correspond 1:1 to an MM commit. If an MM commit had no changes in a downstream, no commit will be created.  We are enforcing squash-merges-only in Magic Modules.

Downstream pushing can only take place if the sync tag points to the commit which precedes the commit that is being pushed.  It is possible for the pipeline to get stuck if a commit is merged into Magic Modules which cannot be generated - this happens most often when the Gemfile.lock is updated.  We expect the failure case to be somewhat uncommon, but when it happens, you need to:
1. Submit your PR changing Gemfile.lock.  The downstream builder will fail in one of the downstream generation steps, whichever one starts first.
2. Immediately after your downstream builder job is submitted, prevent submissions to the magic-modules repository.
3. Once the downstream-builder container is regenerated (about 15 minutes) re-enable submissions to magic-modules and press retry on your downstream builder job.

Disabling submissions to magic-modules can be done through the Admin console in GitHub.

It is safe to have more than one downstream-push running at the same time due to this property, in the event of overruns.  Each run will either
a) make no changes to any downstream and fail
or
b) atomically update every downstream to a fast-forward state that represents the appropriate HEAD as of the beginning of the run

#### Something went wrong!
Don't panic - this is all quite safe.  :)

It's possible for a job to be cancelled or fail in the middle of pushing downstreams in a transient way.  The sorts of failures that happen at scale - lightning strikes a datacenter or some other unlikely misfortune happens.  This has a chance to cause a hiccup in the downstream history, but isn't dangerous.  If that happens, the sync tags may need to be manually updated to sit at the same commit, just before the commit which needs to be generated.  Then, the downstream pusher workflow will need to be restarted.

Updating the sync tags is done like this:
First, check their state: `git fetch origin && git rev-parse origin/tpg-sync origin/tpgb-sync origin/ansible-sync origin/inspec-sync origin/tf-oics-sync origin/tf-conv-sync` will list the commits for each of the sync tags.
If you have changed the name of the `googlecloudplatform/magic-modules` remote from `origin`, substitute that name instead.
In normal, steady-state operation, these tags will all be identical.  When a failure occurs, some of them may be one commit ahead of the others.  It is rare for any of them to be 2 or more commits ahead of any other.  If they are not all equal, and there is no pusher task currently running, this means you need to reset them by hand.  If they are all equal, skip the next step.

Second, find which commit caused the error.  This will usually be easy - cloud build lists the commit which triggered a build, so you can probably just use that one.  You need to set all the sync tags to the parent of that commit.  Say the commit which caused the error is `12345abc`.  You can find the parent of that commit with `git rev-parse 12345abc~` (note the `~` suffix).  Some of the sync tags are likely set to this value already.  For the remainder, simply perform a git push.  Assuming that the parent commit is `98765fed`, that would be `git push origin 98765fed:tf-conv-sync`.

If you are unlucky, there may be open PRs - this only happens if the failure occurred during the ~5 second period surrounding the merging of one of the downstreams.  Close those PRs before proceeding to the final step.

Click "retry" on the failed job in Cloud Build.  Watch the retried job and see if it succeeds - it should!  If it does not, the underlying problem may not have been fixed.

## Deploying the pipeline
The code on the PR's branch is used to plan actions - no merge is performed.
If you are making changes to the workflows, your changes will not trigger a workflow run, because of the risk of an untrusted contributor introducing malicious code in this way.  You will need to test locally by using the [cloud build local builder](https://cloud.google.com/cloud-build/docs/build-debug-locally).
If you are making changes to the containers, your changes will not apply until they are merged in and built - this can take up to 15 minutes.  If you need to make a breaking change, you will need to pause the pipeline while the build happens.  If you are making changes to both the containers and the workflows and those changes need to be coordinated, you will need to pause the build while the containers build and enforce every open PR be rebased on top of our PR.  It is probably better to build in backwards-compatibility into your containers.  We recommend a 14 day window - 14 days after your change goes in, you can remove the backwards-compatibility.

Pausing the pipeline is done in the cloud console, by setting the downstream-builder trigger to disabled.  You can find that trigger [here](https://console.cloud.google.com/cloud-build/triggers/edit/f80a7496-b2f4-4980-a706-c5425a52045b?project=graphite-docker-images)


## Dependency change handbook:
If someone (often a bot) creates a PR which updates Gemfile or Gemfile.lock, they will not be able to generate diffs.  This is because bundler doesn't allow you to run a binary unless your installed gems exactly match the Gemfile.lock, and since we have to run generation before and after the change, there is no possible container that will satisfy all requirements.

The best approach is
* Build the `downstream-generator` container locally, with the new Gemfile and Gemfile.lock.  This will involve hand-modifying the Dockerfile to use the local Gemfile/Gemfile.lock instead of wget from this repo's `master` branch.  You don't need to check in those changes.
* When that container is built, and while nothing else is running in GCB (wait, if you need to), push the container to GCR, and as soon as possible afterwards, merge the dependency-changing PR.

## Historical Note: Design choices & tradeoffs
* The downstream push doesn't wait for checks on its PRs against downstreams.  This may inconvenience some existing workflows which rely on the downstream PR checks.  This ensures that merge conflicts never come into play, since the downstreams never have dangling PRs, but it requires some up-front work to get those checks into the differ.  If a new check is introduced into the downstream Travis, we will need to introduce it into the terraform-tester container.
* The downstream push is disconnected from the output of the differ (but runs the same code).  This means that the diff which is approved isn't guaranteed to be applied *exactly*, if for instance magic modules' behavior changes on master between diff generation and downstream push.  This is also intended to avoid merge conflicts by, effectively, rebasing each commit on top of master before final generation is done.
    * Imagine the following situation: PR A and PR B are opened simultaneously. PR A changes the copyright date in each file to 2020. PR B adds a new resource. PR A is merged seconds before PR B, so they are picked up in the same push-downstream run.  The commit from PR B will produce a new file with the 2020 copyright date, even though the diff said 2019, since PR A was merged first.
* We deleted the submodules.  They weren't useful to us and they were annoying to update - they're not in use anymore as far as we know - but it's possible there's some long-forgotten workflow that someone is using which will be damaged.  So far, we haven't seen any such issue.
