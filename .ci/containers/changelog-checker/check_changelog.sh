#! /bin/bash

set -e

pr_number=$1
mm_commit_sha=$2
build_id=$3
project_id=$4
build_step=$5
github_username=modular-magician

post_body=$( jq -n \
	--arg context "valid-changelog" \
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

/go-changelog/changelog-pr-body-check
exit_code=$?

set -e

if [ $exit_code -ne 0 ]; then
	state="failure"
else
	state="success"
fi

post_body=$( jq -n \
	--arg context "valid-changelog" \
	--arg target_url "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}" \
	--arg state "${state}" \
	'{context: $context, target_url: $target_url, state: $state}')

curl \
  -X POST \
  -u "$github_username:$GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/$mm_commit_sha" \
  -d "$post_body"
