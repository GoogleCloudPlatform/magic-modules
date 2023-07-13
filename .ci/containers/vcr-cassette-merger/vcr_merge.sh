#!/bin/bash

set -e

REFERENCE=$1

PR_NUMBER=$(curl -s -H "Authorization: token ${GITHUB_TOKEN}" \
    "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls?state=closed&base=main&sort=updated&direction=desc" | \
    jq -r ".[] | if .merge_commit_sha == \"$REFERENCE\" then .number else empty end")

# # for backwards-compatibility
# if [ -z "$BRANCH_NAME" ]; then
#   BRANCH_NAME=main
# else
#   echo "BRANCH_NAME: $BRANCH_NAME"
# fi

set +e
gsutil ls gs://ci-vcr-cassettes/refs/heads/auto-pr-$PR_NUMBER/fixtures/

# for testing only
echo "BRANCH_NAME: $BRANCH_NAME"

# if [ $? -eq 0 ]; then
#   # We have recorded new cassettes for this branch
#     if [ "$BRANCH_NAME" == "main" ]; then
#       gsutil -m cp gs://ci-vcr-cassettes/refs/heads/auto-pr-$PR_NUMBER/fixtures/* gs://ci-vcr-cassettes/fixtures/
#     else
#       echo "BRANCH_NAME: $BRANCH_NAME"
#       gsutil -m cp gs://ci-vcr-cassettes/refs/heads/auto-pr-$PR_NUMBER/fixtures/* gs://ci-vcr-cassettes/refs/branches/$BRANCH_NAME/fixtures/
#     fi
#   gsutil -m rm -r gs://ci-vcr-cassettes/refs/heads/auto-pr-$PR_NUMBER/
# fi

# # Beta cassettes
# gsutil ls gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-$PR_NUMBER/fixtures/
# if [ $? -eq 0 ]; then
#   # We have recorded new cassettes for this branch
#     if [ "$BRANCH_NAME" == "main" ]; then
#       gsutil -m cp gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-$PR_NUMBER/fixtures/* gs://ci-vcr-cassettes/beta/fixtures/
#     else
#       echo "BRANCH_NAME: $BRANCH_NAME"
#       gsutil -m cp gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-$PR_NUMBER/fixtures/* gs://ci-vcr-cassettes/beta/refs/branches/$BRANCH_NAME/fixtures/
#     fi
#   gsutil -m rm -r gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-$PR_NUMBER/
# fi


set -e
