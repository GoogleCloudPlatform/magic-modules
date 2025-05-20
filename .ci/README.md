magic-modules CI
===

## Playbook: Resolving downstream build failures

There are four downstreams, each of which has a corresponding "sync branch" in `magic-modules` which tracks the _most recently generated commit_ from magic-modules to that downstream.

- Downstream: [hashicorp/terraform-provider-google](https://github.com/hashicorp/terraform-provider-google)  
  Sync branch: [`tpg-sync`](https://github.com/GoogleCloudPlatform/magic-modules/tree/tpg-sync)
- Downstream: [hashicorp/terraform-provider-google-beta](https://github.com/hashicorp/terraform-provider-google-beta)  
  Sync branch: [`tpgb-sync`](https://github.com/GoogleCloudPlatform/magic-modules/tree/tpgb-sync)
- Downstream: [GoogleCloudPlatform/terraform-google-conversion](https://github.com/GoogleCloudPlatform/terraform-google-conversion)  
  Sync branch: [`tgc-sync`](https://github.com/GoogleCloudPlatform/magic-modules/tree/tgc-sync)
- Downstream: [terraform-google-modules/docs-examples](https://github.com/terraform-google-modules/docs-examples)  
  Sync branch: [`tf-oics-sync`](https://github.com/GoogleCloudPlatform/magic-modules/tree/tf-oics-sync)  
  Note: `oics` refers to "Open in Cloud Shell".

The goal of this system is that each downstream commit will have exactly one MM commit that it corresponds to, and each MM commit will correspond to at most one commit in a downstream. If an MM commit had no changes in a downstream, no commit will be created.

The sync branches allow downstream generation for each downstream to wait until the previous commit for that downstream has finished generating. If downstream generation fails for one commit, the following commits will continue to wait for 24 hours.

Run the following command to verify what commits the sync branches are pointing to:

```
git fetch origin && git rev-parse origin/tpg-sync origin/tpgb-sync origin/tf-oics-sync origin/tgc-sync
```

### Transient GitHub failures
Most downstream build failures are transient GitHub failures. To resolve these, click the "Retry" button in Cloud Build. This is safe because downstream builds are idempotent; if a commit has already been generated, we will not make a new commit.

### Downstream build job is not triggered by commits.
This is rare but we've seen this happened before. In this case, we need to manually trigger a Cloud Build job by running 
```
gcloud builds triggers run build-downstreams --project=graphite-docker-images --substitutions=BRANCH_NAME=<BASE_BRANCH_NAME> --sha=<COMMIT_SHA>
```
You'll need to substitute `<COMMIT_SHA>` with the commit sha that you'd like to trigger the build against and `<BASE_BRANCH_NAME>` with base branch that this commit is pushed into, likely `main` but could be feature branches in some cases.

### Magician / generation bugs
Magician or generation bugs are extremely rare. They cause generation itself to fail loudly, and the commits that introduce them need to be skipped.

Be sure that you understand the mechanism for how the commit is causing a failure before proceeding. Skipping commits will cause multiple MM commits to be squashed into a single downstream commit, which breaks commit linking expectations and release notes.

1. Lock the magic-modules main branch (or ask the team to)
2. Create a PR to fix the bug
3. Unlock the main branch, submit the PR, and re-lock the main branch
4. Update the sync branches to point to the bad commit
   - This indicates that we've "already generated" the downstreams for this commit, so builds for the following commit will stop waiting
6. Verify that downstream generation succeeded - this should only take <15 minutes (even though the overall build-downstreams build takes ~1hr)
7. Unlock the main branch

### Manually pushing commits to downstreams
In general, this should not be necessary because this situation shouldn't come up. It is a historical process that may not work out of the box. The case where you might want to try this is:

- A commit broke the magician (or downstream generation)
- Multiple PRs were merged before the main branch was locked

In this case, skipping just the initial commit would not work (because the following commit doesn't contain a fix) and skipping all the commits is not desirable (because we don't want to squash them.) You would need to instead locally apply the CI fix and then "manually" replicate the work of downstream generation.

Legacy fix (may no longer work):
When this happened the first time, the team wrote this little shell snippet, which might do most of the work for you.  You will need to get the Magician's github token, either by generating a new one (be sure to clean up after yourself when done), by decrypting the value in .ci/gcb-push-downstream.yml as cloudbuild does, or by accessing the token in Google's internal secret store.

```
SYNC_TAG=tpgb-sync
REPO=terraform
VERSION=beta
git clone https://github.com/GoogleCloudPlatform/magic-modules fix-gcb-run
pushd fix-gcb-run
docker pull gcr.io/graphite-docker-images/downstream-builder;
for commit in $(git log $SYNC_TAG..main --pretty=%H | tac); do
  git checkout $commit && \
  docker run -v `pwd`:/workspace -w /workspace -e GITHUB_TOKEN=$MAGICIAN_GITHUB_TOKEN -it gcr.io/graphite-docker-images/downstream-builder downstream $REPO $VERSION $commit || \
  break;
done
```

In the event of a failure, this will stop running.  If it succeeds, update the sync tag with `git push origin HEAD:tpg-sync`.

## Deploying the pipeline
The code on the PR's branch is used to plan actions - no merge is performed.
If you are making changes to the workflows, your changes will not trigger a workflow run, because of the risk of an untrusted contributor introducing malicious code in this way.  You will need to test locally by using the [cloud build local builder](https://cloud.google.com/cloud-build/docs/build-debug-locally).
If you are making changes to the containers, your changes will not apply until they are merged in and built - this can take up to 15 minutes.  If you need to make a breaking change, you will need to pause the pipeline while the build happens.  If you are making changes to both the containers and the workflows and those changes need to be coordinated, you will need to pause the build while the containers build and enforce every open PR be rebased on top of our PR.  It is probably better to build in backwards-compatibility into your containers.  We recommend a 14 day window - 14 days after your change goes in, you can remove the backwards-compatibility.

Pausing the pipeline is done in the cloud console, by setting the downstream-builder trigger to disabled.  You can find that trigger [here](https://console.cloud.google.com/cloud-build/triggers/edit/f80a7496-b2f4-4980-a706-c5425a52045b?project=graphite-docker-images)

## Changes to cloud build yaml:
If changes are made to `gcb-contributor-membership-checker.yml` or `gcb-community-checker.yml` they will not be reflected in presubmit runs for existing PRs without a rebase. This is because these build triggers are linked to pull request creation and not pushes to the PR branch. If changes are needed to these build files they will need to be made in a backwards-compatible manner. Note that changes to other files used by these triggers will be immediately reflected in all PRs, leading to a possible disconnect between the yaml files and the rest of the CI code.
