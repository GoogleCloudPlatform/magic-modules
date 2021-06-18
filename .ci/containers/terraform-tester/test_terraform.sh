#!/bin/bash

set -e

version=$1
pr_number=$2
mm_commit_sha=$3
build_id=$4
project_id=$5
build_step=$6
github_username=modular-magician
if [ "$version" == "ga" ]; then
    gh_repo=terraform-provider-google
elif [ "$version" == "beta" ]; then
    gh_repo=terraform-provider-google-beta
else
    echo "no repo, dying."
    exit 1
fi

scratch_path=https://$github_username:$GITHUB_TOKEN@github.com/$github_username/$gh_repo
local_path=$GOPATH/src/github.com/terraform-providers/$gh_repo
mkdir -p "$(dirname $local_path)"
git clone $scratch_path $local_path --single-branch --branch "auto-pr-$pr_number" --depth 1
pushd $local_path


post_body=$( jq -n \
    --arg context "${gh_repo}-test" \
    --arg target_url "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}" \
    --arg state "pending" \
    '{context: $context, target_url: $target_url, state: $state}')

curl \
  -X POST \
  -u "$github_username:$GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/$mm_commit_sha" \
  -d "$post_body"


post_body=$( jq -n \
    --arg context "${gh_repo}-lint" \
    --arg target_url "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}" \
    --arg state "pending" \
    '{context: $context, target_url: $target_url, state: $state}')

curl \
  -X POST \
  -u "$github_username:$GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/$mm_commit_sha" \
  -d "$post_body"

set +e

make
lint_exit_code=$?
test_exit_code=1

if [ $lint_exit_code -eq 0 ]; then
    # only run lint & tests if the code compiled
    make lint
    lint_exit_code=$lint_exit_code || $?
    make test
    test_exit_code=$?
fi

make tools
lint_exit_code=$lint_exit_code || $?
make docscheck
lint_exit_code=$lint_exit_code || $?

set -e

if [ $test_exit_code -ne 0 ]; then
    test_state="failure"
else
    test_state="success"
fi

if [ $lint_exit_code -ne 0 ]; then
    lint_state="failure"
else
    lint_state="success"
fi

post_body=$( jq -n \
    --arg context "${gh_repo}-test" \
    --arg target_url "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}" \
    --arg state "${test_state}" \
    '{context: $context, target_url: $target_url, state: $state}')

curl \
  -X POST \
  -u "$github_username:$GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/$mm_commit_sha" \
  -d "$post_body"


post_body=$( jq -n \
    --arg context "${gh_repo}-lint" \
    --arg target_url "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}" \
    --arg state "${lint_state}" \
    '{context: $context, target_url: $target_url, state: $state}')

curl \
  -X POST \
  -u "$github_username:$GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/$mm_commit_sha" \
  -d "$post_body"
