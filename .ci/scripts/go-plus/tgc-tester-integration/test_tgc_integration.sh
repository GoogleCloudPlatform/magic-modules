#!/bin/bash

set -e

pr_number=$1
mm_commit_sha=$2
build_id=$3
project_id=$4
build_step=$5
gh_repo=$6
github_username=modular-magician


new_branch="auto-pr-$pr_number"
git_remote=https://$github_username:$GITHUB_TOKEN@github.com/$github_username/$gh_repo
local_path=$GOPATH/src/github.com/GoogleCloudPlatform/$gh_repo
mkdir -p "$(dirname $local_path)"
git clone $git_remote $local_path --branch $new_branch --depth 2
pushd $local_path

# Only skip tests if we can tell for sure that no go files were changed
echo "Checking for modified go files"
# get the names of changed files and look for go files
# (ignoring "no matches found" errors from grep)
gofiles=$(git diff --name-only HEAD~1 | { grep "\.go$" || test $? = 1; })
if [[ -z $gofiles ]]; then
    echo "Skipping tests: No go files changed"
    exit 0
else
    echo "Running tests: Go files changed"
fi

post_body=$( jq -n \
	--arg context "${gh_repo}-test-integration" \
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

go mod edit -replace github.com/hashicorp/terraform-provider-google-beta=github.com/$github_username/terraform-provider-google-beta@$new_branch
go mod tidy

make build
make test-integration
exit_code=$?

set -e

if [ $exit_code -ne 0 ]; then
	state="failure"
else
	state="success"
fi

post_body=$( jq -n \
	--arg context "${gh_repo}-test-integration" \
	--arg target_url "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}" \
	--arg state "${state}" \
	'{context: $context, target_url: $target_url, state: $state}')

curl \
  -X POST \
  -u "$github_username:$GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/$mm_commit_sha" \
  -d "$post_body"
