#!/bin/bash

set -e

pr_number=$1
mm_commit_sha=$2
build_id=$3
project_id=$4
build_step=$5
github_username=modular-magician
gh_repo=terraform-google-conversion

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

update_status "terraform-google-conversion-test" "pending" "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}"

set +e

make test
exit_code=$?

set -e

if [ $exit_code -ne 0 ]; then
	state="failure"
else
	state="success"
fi

update_status "terraform-google-conversion-test" "${state}" "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}"
