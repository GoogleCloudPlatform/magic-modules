#!/bin/bash

set -e
pr_number=$1
mm_commit_sha=$2
echo "PR number: ${pr_number}"
echo "Commit SHA: ${mm_commit_sha}"
github_username=modular-magician
gh_repo=terraform-provider-google-beta

new_branch="auto-pr-$pr_number"
git_remote=https://github.com/$github_username/$gh_repo
local_path=$GOPATH/src/github.com/hashicorp/$gh_repo
mkdir -p "$(dirname $local_path)"
git clone $git_remote $local_path --branch $new_branch --depth 2
pushd $local_path

# Only skip tests if we can tell for sure that no go files were changed
echo "Checking for modified go files"
# get the names of changed files and look for go files
# (ignoring "no matches found" errors from grep)
gofiles=$(git diff --name-only HEAD~1 | { grep -e "\.go$" -e "go.mod$" -e "go.sum$" || test $? = 1; })
if [[ -z $gofiles ]]; then
  echo "Skipping tests: No go files changed"
  exit 0
else
  echo "Running tests: Go files changed"
fi

function add_comment {
  curl -H "Authorization: token ${GITHUB_TOKEN}" \
    -d "$(jq -r --arg comment "${1}" -n "{body: \$comment}")" \
    "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/${2}/comments"
}


# The following section is to cancel jobs that may be running for the given PR, before starting to run new tests
function cancel_queued {
  add_comment "Canceling build id ${1}"
  # Cancel tests that are queued up to run in TeamCity
  curl \
    --header "Accept: application/json" \
    --header "Authorization: Bearer $TEAMCITY_TOKEN" \
    --header "Content-Type:application/xml" \
    --data-binary @/teamcitycancelparams.xml \
    --request POST \
    -o canceled-queue.json \
    "https://ci-oss.hashicorp.engineering/app/rest/buildQueue/id:${1}"
}

function cancel_running {
  add_comment "Canceling build for running pr number ${1}"
  # Cancel tests that are currently running in TeamCity
  curl \
    --header "Accept: application/json" \
    --header "Authorization: Bearer $TEAMCITY_TOKEN" \
    --header "Content-Type:application/xml" \
    --data-binary @/teamcitycancelparams.xml \
    --request POST \
    -o canceled-running.json \
    "https://ci-oss.hashicorp.engineering/app/rest/builds/multiple/branch:name:(\$base64:${1}),buildType:(id:GoogleCloudBeta_ProviderGoogleCloudBetaMmUpstreamVcr),property:(name:env.VCR_MODE,value:REPLAYING),running:true,defaultFilter:false"
}

# We need to get the ids of the tests that are queued up to run before we can cancel them
curl \
  --header "Accept: application/json" \
  --header "Authorization: Bearer $TEAMCITY_TOKEN" \
  --header "Content-Type:application/json" \
  -o existing-queue.json \
  "https://ci-oss.hashicorp.engineering/app/rest/buildQueue?locator=buildType:(id:GoogleCloudBeta_ProviderGoogleCloudBetaMmUpstreamVcr)&fields=build(id,comment)"

queued_ids=$(cat existing-queue.json | jq -r ".build[] | select(.comment.text | endswith(\"${pr_number}\")) | .id")
for id in $queued_ids; do
  cancel_queued ${id}
done

# We can just cancel running jobs by a filter
branchname=$(echo -n "/auto-pr-${pr_number}"| base64)
cancel_running ${branchname}
canceled_errs=$(cat canceled-running.json | jq .errorCount -r)
if [ "${canceled_errs}" -gt "0" ]; then
  for row in $(cat canceled-running.json | jq -r '.operationResult[] | @base64'); do
    build_url=$(${row} | base64 --decode | jq -r '.related.build | select(.state == "running") | .webUrl')
    add_comment "Error trying to cancel build (${build_url})" ${pr_number}
  done
fi

# Old jobs should have been canceled, create the new jobs

sed -i 's/{{PR_NUMBER}}/'"$pr_number"'/g' /teamcityparams.xml
curl --header "Accept: application/json" --header "Authorization: Bearer $TEAMCITY_TOKEN" https://ci-oss.hashicorp.engineering/app/rest/buildQueue --request POST --header "Content-Type:application/xml" --data-binary @/teamcityparams.xml -o build.json

function update_status {
  local context="beta-provider-vcr-test"
  local post_body=$( jq -n \
    --arg context "${context}" \
    --arg target_url "${1}" \
    --arg state "${2}" \
    '{context: $context, target_url: $target_url, state: $state}')
  echo "Updating status ${context} to ${2} with target_url ${1} for sha ${mm_commit_sha}"
  curl \
    -X POST \
    -u "$github_username:$GITHUB_TOKEN" \
    -H "Accept: application/vnd.github.v3+json" \
    "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/${mm_commit_sha}" \
    -d "$post_body"
}

build_url=$(cat build.json | jq -r .webUrl)
update_status "${build_url}" "pending"

ID=$(cat build.json | jq .id -r)
curl --header "Authorization: Bearer $TEAMCITY_TOKEN" --header "Accept: application/json" https://ci-oss.hashicorp.engineering/app/rest/builds/id:$ID --output poll.json
STATUS=$(cat poll.json | jq .status -r)
STATE=$(cat poll.json | jq .state -r)
counter=0
while [[ "$STATE" != "finished" ]]; do
  if [ "$counter" -gt "500" ]; then
    echo "Failed to wait for job to finish, exiting"
    # Call this an error because we don't know if the tests failed or not
    update_status "${build_url}" "error"
    # exit 0 because this script didn't have an error; the failure
    # is reported via the Github Status API
    exit 0
  fi
  sleep 30
  curl --header "Authorization: Bearer $TEAMCITY_TOKEN" --header "Accept: application/json" https://ci-oss.hashicorp.engineering/app/rest/builds/id:$ID --output poll.json
  STATUS=$(cat poll.json | jq .status -r)
  STATE=$(cat poll.json | jq .state -r)
  echo "Trying again, State: $STATE Status: $STATUS"
  counter=$((counter + 1))
done

if [ "$STATUS" == "UNKNOWN" ]; then
    echo "Run canceled."
    exit 0
fi

if [ "$STATUS" == "SUCCESS" ]; then
  echo "Tests succeeded."
  update_status "${build_url}" "success"
  exit 0
fi

# This is an intentionally dumb list; if something is removed and re-added with
# the same name, we'll still catch it. If that ends up causing noise, we can do
# something more clever.
NEW_TESTS=$(git diff --unified=0 HEAD~1 | { grep -oP '(?<=^\+func )TestAcc\w+(?=\(t \*testing.T\) {)' || test $? = 1; } | tr '\n' ' ')
RUN_TESTS=""
FAILED_TESTS=""
NEWLINE=$'\n'
set +e
TEST_RESULTS_URL="http://ci-oss.hashicorp.engineering/app/rest/testOccurrences?locator=build:${ID}"
while [ -n "${TEST_RESULTS_URL}" ]; do
  echo $TEST_RESULTS_URL
  curl \
    --header "Accept: application/json" \
    --header "Authorization: Bearer $TEAMCITY_TOKEN" \
    --output tests.json \
    -L \
    "${TEST_RESULTS_URL}"

  # Alert on tests that failed without running anything
  if [[ $(cat tests.json | jq -r '.count') == "0" ]]; then
    echo "Job failed without failing tests"
    update_status "${build_url}" "failure"
    # exit 0 because this script didn't have an error; the failure
    # is reported via the Github Status API
    exit 0
  fi

  NEW_RUN_TESTS=$(cat tests.json | jq -r '.testOccurrence | map(select(.status != "UNKNOWN"))| map(.name) | .[]')
  NEW_FAILED_TESTS=$(cat tests.json | jq -r '.testOccurrence | map(select(.status == "FAILURE")) | map(.name) | join("|")')
  if [ -n "$NEW_RUN_TESTS" ] && [ -n "$RUN_TESTS" ]
  then
    RUN_TESTS+="${NEWLINE}"
  fi
  if [ -n "$NEW_FAILED_TESTS" ] && [ -n "$FAILED_TESTS" ]
  then
    FAILED_TESTS+="|"
  fi

  RUN_TESTS+=$NEW_RUN_TESTS
  FAILED_TESTS+=$NEW_FAILED_TESTS

  next_href=$(cat tests.json | jq -r '.nextHref')
  if [[ $next_href == "null" ]]; then
    break
  else
    TEST_RESULTS_URL="http://ci-oss.hashicorp.engineering$next_href"
  fi
done

MISSING_TESTS=""
for new_test in $NEW_TESTS; do
  if ! echo "${RUN_TESTS}" | grep -qP "^${new_test}$"; then
    MISSING_TESTS+="- ${new_test}${NEWLINE}"
  fi
done

if [[ -n $MISSING_TESTS ]]; then
  comment="Tests were added that did not run in TeamCity:${NEWLINE}${NEWLINE}"
  comment+=${MISSING_TESTS}
  add_comment "${comment}" "${pr_number}"
fi

set -e

sed -i 's/{{PR_NUMBER}}/'"$pr_number"'/g' /teamcityparamsrecording.xml
sed -i 's/{{FAILED_TESTS}}/'"$FAILED_TESTS"'/g' /teamcityparamsrecording.xml
curl --header "Accept: application/json" --header "Authorization: Bearer $TEAMCITY_TOKEN" https://ci-oss.hashicorp.engineering/app/rest/buildQueue --request POST --header "Content-Type:application/xml" --data-binary @/teamcityparamsrecording.xml --output record.json
build_url=$(cat record.json | jq -r .webUrl)
comment="I have triggered VCR tests in RECORDING mode for the following tests that failed during VCR: $FAILED_TESTS You can view the result here: $build_url"

add_comment "${comment}" ${pr_number}
update_status "${build_url}" "pending"

# Reset for checking failed tests
rm poll.json
rm tests.json

ID=$(cat record.json | jq .id -r)
curl --header "Authorization: Bearer $TEAMCITY_TOKEN" --header "Accept: application/json" https://ci-oss.hashicorp.engineering/app/rest/builds/id:$ID --output poll.json
STATUS=$(cat poll.json | jq .status -r)
STATE=$(cat poll.json | jq .state -r)
counter=0
while [[ "$STATE" != "finished" ]]; do
  if [ "$counter" -gt "500" ]; then
    echo "Failed to wait for job to finish, exiting"
    # Call this an error because we don't know if the tests failed or not
    update_status "${build_url}" "error"
    # exit 0 because this script didn't have an error; the failure
    # is reported via the Github Status API
    exit 0
  fi
  sleep 30
  curl --header "Authorization: Bearer $TEAMCITY_TOKEN" --header "Accept: application/json" https://ci-oss.hashicorp.engineering/app/rest/builds/id:$ID --output poll.json
  STATUS=$(cat poll.json | jq .status -r)
  STATE=$(cat poll.json | jq .state -r)
  echo "Trying again, State: $STATE Status: $STATUS"
  counter=$((counter + 1))
done

if [ "$STATUS" == "SUCCESS" ]; then
  echo "Tests succeeded."
  update_status "${build_url}" "success"
  exit 0
fi

curl --header "Accept: application/json" --header "Authorization: Bearer $TEAMCITY_TOKEN" http://ci-oss.hashicorp.engineering/app/rest/testOccurrences?locator=build:$ID,status:FAILURE --output failed.json -L
set +e
FAILED_TESTS=$(cat failed.json | jq -r '.testOccurrence | map(.name) | join("|")')
ret=$?
if [ $ret -ne 0 ]; then
  echo "Job failed without failing tests"
  update_status "${build_url}" "failure"
  # exit 0 because this script didn't have an error; the failure
  # is reported via the Github Status API
  exit 0
fi
set -e

comment="Tests failed during RECORDING mode: $FAILED_TESTS Please fix these to complete your PR"

add_comment "${comment}" ${pr_number}
update_status "${build_url}" "failure"

# exit 0 because this script didn't have an error; the failure
# is reported via the Github Status API
exit 0
