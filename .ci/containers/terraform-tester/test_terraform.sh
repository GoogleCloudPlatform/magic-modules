#!/bin/bash

source /utils.sh

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

new_branch="auto-pr-$pr_number"
git_remote=https://$github_username:$GITHUB_TOKEN@github.com/$github_username/$gh_repo
local_path=$GOPATH/src/github.com/hashicorp/$gh_repo
mkdir -p "$(dirname $local_path)"
git clone $git_remote $local_path --branch $new_branch --depth 2
pushd $local_path

# Only skip tests if no go files were changed
echo "Checking for modified go files"
if grep_files_modified "\.go$"; then
    echo "Running tests: Go files changed"
else
    echo "Skipping tests: No go files changed"
    exit 0
fi

update_status "${gh_repo}-test" "pending" "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}"
update_status "${gh_repo}-lint" "pending" "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}"

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

update_status "${gh_repo}-test" "${test_state}" "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}"
update_status "${gh_repo}-lint" "${lint_state}" "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}"
