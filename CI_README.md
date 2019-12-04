GitHub CI tools for MagicModules and Google Providers
===

These tools manage the downstream repositories of [magic-modules](https://github.com/GoogleCloudPlatform/magic-modules), and are collectively referred to as "The Magician".

# CI For Downstream/Magic Modules Developers
If you're interested in developing the repositories that Magic Modules manages, here are the things you'll want to know.

## What The Magician Does
The Magician takes the PR you write against MagicModules and creates the downstream (generated) commits that are a result of your changes.  It posts the diffs created by your changes as a comment on your PR, to aid in review.  When your PR is merged, it updates the downstream repositories with your changes.

## Your Workflow

You'll write some code and open a PR.  The Magician will run the generator and the downstream tests.  If either the generator or the tests fail, you'll see "check failed" status on your PR, which will link to the step that failed.  Once all the generation steps succeed, the Magician will comment on your PR with links to review the generated diffs.  Your reviewer will review those diffs (as well as your code).  Once your PR is approved, simply merge it in - the Magician will ensure that the downstream repositories are updated within about 20 minutes.

# CI For CI Developers
If you're working on enhancing these CI tools, you'll need to know about the structure of the pipeline and about how to develop and test.

## How the pipeline works
The pipeline is written in Github Actions, and is defined in the workflow .yml files in the .github/workflows directory of this repository.  The documentation on those files is located [here](https://help.github.com/en/actions/automating-your-workflow-with-github-actions/building-actions).  Documentation on the substeps is located below.

### Generation & Diffing
The generation / diff pipeline has one Action per downstream.  It generates the downstream at the PR's merge commit and at the left-side parent (`HEAD~`), which guarantees isolation of exclusively the changes made by the PR in question.  It creates commits for the downstreams as before-and-after commits.  Those commits don't have relevant git history - they're not meant to be applied - but they are used in a [two-dot diff](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/about-comparing-branches-in-pull-requests#three-dot-and-two-dot-git-diff-comparisons).  This means that if there are no changes to the downstream, you'll see an empty diff.  These downstream branches are named `pr-$number-old` and `pr-$number-new` respectively.

### Downstream Pushing
The downstream pushing pipeline runs on a cron job - every 20 minutes.  It checks a branch called `downstream-master` to find the commit which was most recently pushed to downstreams, then collects all commits since then.  One at a time, it generates each downstream at each commit, building an in-order history.  It pushes all the downstreams directly to master branches, but it does not use the `--force` flag - if the branch goes out of sync, the push will fail without damaging the downstream repository history.

## Deploying the pipeline
The code on the `master` branch is used in all cases - if you are making changes to the Actions, your changes will not apply until they are merged in.

## Design choices & tradeoffs
* The downstreams share some setup code in common - especially TPG and TPGB.  We violated the DRY principle by writing separate workflows for each repo.  In practice, this has substantially reduced the amount of code - the coordination layer above the two repos was larger than the code saved by combining them.  We also increase speed, since each Action runs separately.
* The downstream push doesn't happen immediately.  We caused some delay with the cron-based approach that could have been avoided.  This ensures that two commits which are merged at about the same time will never conflict, since only one copy of the downstream push will be running at a time.
* The downstream push doesn't open PRs against downstreams.  This may inconvenience some existing workflows which rely on the downstream PRs.  This ensures that merge conflicts never come into play, since the downstreams never have dangling PRs.
* The downstream push is totally disconnected from the differ.  This means that the diff which is approved isn't guaranteed to be applied *exactly*, if for instance magic modules' behavior changes on master between diff generation and downstream push.  This is also intended to avoid merge conflicts by, effectively, rebasing each commit on top of master before final generation is done.
* We deleted the submodules.  They weren't useful to us and they were annoying to update - no one ever used them for anything as far as we know - but it's possible there's some long-forgotten workflow that someone is using which will be damaged.
