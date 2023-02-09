#!/bin/bash

MM_LOCAL_PATH=$1
TPG_LOCAL_PATH=$2
mm_commit_sha=$3
build_id=$4
build_step=$5
project_id=$6

github_username=modular-magician

pushd $MM_LOCAL_PATH/tools/breaking-change-detector
go test
exit_code=$?
popd


if [ $exit_code -ne 0 ]; then
	state="failure"
else
	state="success"
fi

post_body=$( jq -n \
	--arg context "breaking-change-detector-test" \
	--arg target_url "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}" \
	--arg state "${state}" \
	'{context: $context, target_url: $target_url, state: $state}')

curl \
  -X POST \
  -u "$github_username:$GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/$mm_commit_sha" \
  -d "$post_body"

build_step=$((build_step + 1))

pushd $MM_LOCAL_PATH/tools/missing-test-detector
PROVIDER_DIR=$TPG_LOCAL_PATH go test
exit_code=$?
popd


if [ $exit_code -ne 0 ]; then
	state="failure"
else
	state="success"
fi

post_body=$( jq -n \
	--arg context "missing-test-detector-test" \
	--arg target_url "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}" \
	--arg state "${state}" \
	'{context: $context, target_url: $target_url, state: $state}')

curl \
  -X POST \
  -u "$github_username:$GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/$mm_commit_sha" \
  -d "$post_body"

