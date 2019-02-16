#! /bin/bash
set -e
set -x

PR_ID="$(cat ./mm-initial-pr/.git/id)"
ORIGINAL_PR_BRANCH="codegen-pr-$PR_ID"
pushd branchname
echo "$ORIGINAL_PR_BRANCH" > ./original_pr_branch_name
