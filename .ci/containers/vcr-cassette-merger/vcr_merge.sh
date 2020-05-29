#!/bin/bash

set -e

REFERENCE=$1

PR_NUMBER=$(curl -s -H "Authorization: token ${GITHUB_TOKEN}" \
    "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls?state=closed&base=master&sort=updated&direction=desc" | \
    jq -r ".[] | if .merge_commit_sha == \"$REFERENCE\" then .number else empty end")

set +e
gsutil ls gs://vcr-$GOOGLE_PROJECT/auto-pr-$PR_NUMBER/fixtures/
if [ $? -eq 0 ]; then
	# We have recorded new cassettes for this branch
  gsutil -m cp gs://vcr-$GOOGLE_PROJECT/refs/heads/auto-pr-$PR_NUMBER/fixtures/* gs://vcr-$GOOGLE_PROJECT/fixtures/
  gsutil -m rm -r gs://vcr-$GOOGLE_PROJECT/refs/heads/auto-pr-$PR_NUMBER/
fi
set -e