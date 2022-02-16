#!/bin/bash

set -e
pr_number=$1
mm_commit_sha=$2
build_id=$3
project_id=$4
build_step=$5
echo "PR number: ${pr_number}"
echo "Commit SHA: ${mm_commit_sha}"
github_username=modular-magician
gh_repo=terraform-provider-google-beta

# For testing only.
echo "checking if this is the testing PR for vcr setup"
if [ "$pr_number" == "5703" ]; then
  echo "Running tests for new VCR setup"
else
  echo "Skipping new vcr tests: Not testing PR"
  exit 0
fi

echo "Keep running for testing PR"
echo "Project id: ${project_id}"

new_branch="auto-pr-$pr_number"
git_remote=https://github.com/$github_username/$gh_repo
local_path=$GOPATH/src/github.com/hashicorp/$gh_repo
mkdir -p "$(dirname $local_path)"
git clone $git_remote $local_path --branch $new_branch --depth 2
pushd $local_path

# # Only skip tests if we can tell for sure that no go files were changed
# echo "Checking for modified go files"
# # get the names of changed files and look for go files
# # (ignoring "no matches found" errors from grep)
# gofiles=$(git diff --name-only HEAD~1 | { grep -e "\.go$" -e "go.mod$" -e "go.sum$" || test $? = 1; })
# if [[ -z $gofiles ]]; then
#   echo "Skipping tests: No go files changed"
#   exit 0
# else
#   echo "Running tests: Go files changed"
# fi

post_body=$( jq -n \
    --arg context "VCR-test" \
    --arg target_url "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}" \
    --arg state "pending" \
    '{context: $context, target_url: $target_url, state: $state}')

curl \
  -X POST \
  -u "$github_username:$GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/$mm_commit_sha" \
  -d "$post_body"


# cassette retrieval
gsutil -m cp gs://vcr-$GOOGLE_PROJECT/beta/fixtures/* /fixtures/
# # copy branch specific cassettes over master. This might fail but that's ok if the folder doesnt exist
# gsutil -m cp gs://vcr-$GOOGLE_PROJECT/beta/%BRANCH_NAME%/fixtures/* fixtures/
# mkdir -p $VCR_PATH
# cp fixtures $VCR_PATH/../
# ls $VCR_PATH

export TF_LOG=DEBUG
export GOOGLE_REGION=us-central1
export GOOGLE_ZONE=us-central1-a
export GOOGLE_USE_DEFAULT_CREDENTIALS=TRUE
export VCR_PATH=/fixtures
export VCR_MODE=REPLAYING

ls $VCR_PATH

make testacc TEST=./google-beta TESTARGS='-run=TestAccComputeInstanceTemplate_basic'

test_exit_code=$?

if [ $test_exit_code -ne 0 ]; then
    test_state="failure"
else
    test_state="success"
fi

post_body=$( jq -n \
    --arg context "VCR-test" \
    --arg target_url "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}" \
    --arg state "${test_state}" \
    '{context: $context, target_url: $target_url, state: $state}')

curl \
  -X POST \
  -u "$github_username:$GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/$mm_commit_sha" \
  -d "$post_body"




