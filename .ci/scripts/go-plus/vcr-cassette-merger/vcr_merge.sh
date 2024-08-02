#!/bin/bash

set -e

REFERENCE=$1

# for backwards-compatibility
if [ -z "$BASE_BRANCH" ]; then
  BASE_BRANCH=main
else
  echo "BASE_BRANCH: $BASE_BRANCH"
fi

PR_NUMBER=$(curl -s -H "Authorization: token ${GITHUB_TOKEN_CLASSIC}" \
    "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls?state=closed&base=$BASE_BRANCH&sort=updated&direction=desc" | \
    jq -r ".[] | if .merge_commit_sha == \"$REFERENCE\" then .number else empty end")

set +e
gsutil ls gs://ci-vcr-cassettes/refs/heads/auto-pr-$PR_NUMBER/fixtures/
if [ $? -eq 0 ]; then
  # We have recorded new cassettes for this branch
    if [ "$BASE_BRANCH" == "main" ]; then
      gsutil -m cp gs://ci-vcr-cassettes/refs/heads/auto-pr-$PR_NUMBER/fixtures/* gs://ci-vcr-cassettes/fixtures/
    else
      gsutil -m cp gs://ci-vcr-cassettes/refs/heads/auto-pr-$PR_NUMBER/fixtures/* gs://ci-vcr-cassettes/refs/branches/$BASE_BRANCH/fixtures/
    fi
  gsutil -m rm -r gs://ci-vcr-cassettes/refs/heads/auto-pr-$PR_NUMBER/
fi

# Beta cassettes
gsutil ls gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-$PR_NUMBER/fixtures/
if [ $? -eq 0 ]; then
  # We have recorded new cassettes for this branch
    if [ "$BASE_BRANCH" == "main" ]; then
      gsutil -m cp gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-$PR_NUMBER/fixtures/* gs://ci-vcr-cassettes/beta/fixtures/
    else
      gsutil -m cp gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-$PR_NUMBER/fixtures/* gs://ci-vcr-cassettes/beta/refs/branches/$BASE_BRANCH/fixtures/
    fi
  gsutil -m rm -r gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-$PR_NUMBER/
fi


set -e
