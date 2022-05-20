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
NEWLINE=$'\n'

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
    "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/${pr_number}/comments"
}

function update_status {
  post_body=$( jq -n \
    --arg context "VCR-test" \
    --arg target_url "https://console.cloud.google.com/cloud-build/builds;region=global/${build_id};step=${build_step}?project=${project_id}" \
    --arg state "${1}" \
    '{context: $context, target_url: $target_url, state: $state}')

  curl \
    -X POST \
    -u "$github_username:$GITHUB_TOKEN" \
    -H "Accept: application/vnd.github.v3+json" \
    "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/$mm_commit_sha" \
    -d "$post_body"
}

update_status "pending"

set +e
# cassette retrieval
mkdir fixtures
gsutil -m -q cp gs://ci-vcr-cassettes/beta/fixtures/* fixtures/
# copy branch specific cassettes over master. This might fail but that's ok if the folder doesnt exist
gsutil -m -q cp gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-$pr_number/fixtures/* fixtures/

echo $SA_KEY > sa_key.json
gcloud auth activate-service-account $GOOGLE_SERVICE_ACCOUNT --key-file=$local_path/sa_key.json --project=$GOOGLE_PROJECT

mkdir testlog
mkdir testlog/replaying
mkdir testlog/recording

export GOOGLE_REGION=us-central1
export GOOGLE_ZONE=us-central1-a
export VCR_PATH=$local_path/fixtures
export VCR_MODE=REPLAYING
export ACCTEST_PARALLELISM=32
export GOOGLE_CREDENTIALS=$SA_KEY
export GOOGLE_APPLICATION_CREDENTIALS=$local_path/sa_key.json

echo "checking terraform version"
terraform version

TF_LOG=DEBUG TF_LOG_PATH_MASK=$local_path/testlog/replaying/%s.log TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test ./google-beta -parallel $ACCTEST_PARALLELISM -v -run=TestAcc -timeout 240m -ldflags="-X=github.com/hashicorp/terraform-provider-google-beta/version.ProviderVersion=acc" > replaying_test.log

test_exit_code=$?

TESTS_TERMINATED=$(grep "^cannot run Terraform provider tests" replaying_test.log)

counter=1
test_suffix=""

while [[ -n $TESTS_TERMINATED ]]; do
  # store the previous replaying build log
  gsutil -h "Content-Type:text/plain" -q cp replaying_test$test_suffix.log gs://ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/build-log/

  if [[ $counter -gt 3 ]]; then
    comment="Failed to run VCR tests in REPLAYING mode${NEWLINE}"
    comment+="View the [build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/build-log/replaying_test$test_suffix.log)${NEWLINE}"
    comment+="If you believe the error is unrelated to your PR, please rerun the tests"
    add_comment "${comment}"
    update_status "failure"
    exit 0
  fi

  comment="Rerun tests in REPLAYING mode"
  add_comment "${comment}"

  test_suffix="$counter"

  # rerun the test
  TF_LOG=DEBUG TF_LOG_PATH_MASK=$local_path/testlog/replaying/%s.log TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test ./google-beta -parallel $ACCTEST_PARALLELISM -v -run=TestAcc -timeout 240m -ldflags="-X=github.com/hashicorp/terraform-provider-google-beta/version.ProviderVersion=acc" > replaying_test$test_suffix.log
  test_exit_code=$?
  TESTS_TERMINATED=$(grep "^cannot run Terraform provider tests" replaying_test$test_suffix.log)
  counter=$((counter + 1))
done

# store replaying build log
gsutil -h "Content-Type:text/plain" -q cp replaying_test$test_suffix.log gs://ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/build-log/

# store replaying test logs
gsutil -h "Content-Type:text/plain" -m -q cp testlog/replaying/* gs://ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/replaying/

# handle provider crash
TESTS_PANIC=$(grep "^panic: " replaying_test$test_suffix.log)

if [[ -n $TESTS_PANIC ]]; then
  comment="The provider crashed while running the VCR tests in REPLAYING mode${NEWLINE}"
  comment+="Please fix it to complete your PR${NEWLINE}"
  comment+="View the [build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/build-log/replaying_test$test_suffix.log)"
  add_comment "${comment}"
  update_status "failure"
  exit 0
fi

FAILED_TESTS=$(grep "^--- FAIL: TestAcc" replaying_test$test_suffix.log)
PASSED_TESTS=$(grep "^--- PASS: TestAcc" replaying_test$test_suffix.log)
SKIPPED_TESTS=$(grep "^--- SKIP: TestAcc" replaying_test$test_suffix.log)

if [[ -n $FAILED_TESTS ]]; then
  FAILED_TESTS_COUNT=$(echo "$FAILED_TESTS" | wc -l)
else
  FAILED_TESTS_COUNT=0
fi

if [[ -n $PASSED_TESTS ]]; then
  PASSED_TESTS_COUNT=$(echo "$PASSED_TESTS" | wc -l)
else
  PASSED_TESTS_COUNT=0
fi

if [[ -n $SKIPPED_TESTS ]]; then
  SKIPPED_TESTS_COUNT=$(echo "$SKIPPED_TESTS" | wc -l)
else
  SKIPPED_TESTS_COUNT=0
fi

FAILED_TESTS_PATTERN=$(grep "^--- FAIL: TestAcc" replaying_test$test_suffix.log | awk '{print $3}' | awk -v d="|" '{s=(NR==1?s:s d)$0}END{print s}')

comment="#### Tests analytics ${NEWLINE}"
comment+="Total tests: \`$(($FAILED_TESTS_COUNT+$PASSED_TESTS_COUNT+$SKIPPED_TESTS_COUNT))\` ${NEWLINE}"
comment+="Passed tests \`$PASSED_TESTS_COUNT\` ${NEWLINE}"
comment+="Skipped tests: \`$SKIPPED_TESTS_COUNT\` ${NEWLINE}"
comment+="Failed tests: \`$FAILED_TESTS_COUNT\` ${NEWLINE}${NEWLINE}"

if [[ -n $FAILED_TESTS_PATTERN ]]; then
  comment+="#### Action taken ${NEWLINE}"
  comment+="Triggering VCR tests in RECORDING mode for the following tests that failed during VCR: $FAILED_TESTS_PATTERN"
  add_comment "${comment}"
  # RECORDING mode
  export VCR_MODE=RECORDING
  # Clear fixtures folder
  rm $VCR_PATH/*
  TF_LOG=DEBUG TF_LOG_PATH_MASK=$local_path/testlog/recording/%s.log TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test ./google-beta -parallel 1 -v -run=$FAILED_TESTS_PATTERN -timeout 240m -ldflags="-X=github.com/hashicorp/terraform-provider-google-beta/version.ProviderVersion=acc" > recording_test.log
  test_exit_code=$?

  # store cassettes
  gsutil -m -q cp fixtures/* gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-$pr_number/fixtures/

  # store recording build log
  gsutil -h "Content-Type:text/plain" -q cp recording_test.log gs://ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/build-log/

  # store recording test logs
  gsutil -h "Content-Type:text/plain" -m -q cp testlog/recording/* gs://ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/recording/

  # handle provider crash
  RECORDING_TESTS_PANIC=$(grep "^panic: " recording_test.log)

  if [[ -n $RECORDING_TESTS_PANIC ]]; then
    comment="The provider crashed while running the VCR tests in RECORDING mode${NEWLINE}"
    comment+="Please fix it to complete your PR${NEWLINE}"
    comment+="View the [build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/build-log/recording_test.log)"
    add_comment "${comment}"
    update_status "failure"
    exit 0
  fi


  RECORDING_FAILED_TESTS=$(grep "^--- FAIL: TestAcc" recording_test.log | awk -v pr_number=$pr_number -v build_id=$build_id '{print "`"$3"`[[view](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-"pr_number"/artifacts/"build_id"/recording/"$3".log)]"}')
  RECORDING_PASSED_TESTS=$(grep "^--- PASS: TestAcc" recording_test.log | awk -v pr_number=$pr_number -v build_id=$build_id '{print "`"$3"`[[view](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-"pr_number"/artifacts/"build_id"/recording/"$3".log)]"}')

  comment=""
  if [[ -n $RECORDING_PASSED_TESTS ]]; then
    comment+="Tests passed during RECORDING mode:${NEWLINE} $RECORDING_PASSED_TESTS ${NEWLINE}${NEWLINE}"
  fi

  if [[ -n $RECORDING_FAILED_TESTS ]]; then
    comment+="Tests failed during RECORDING mode:${NEWLINE} $RECORDING_FAILED_TESTS ${NEWLINE}${NEWLINE}"
    comment+="Please fix these to complete your PR${NEWLINE}"
  else
    if [[ $test_exit_code -ne 0 ]]; then
      # check for any uncaught errors in RECORDING mode
      comment+="Errors occurred during RECORDING mode. Please fix them to complete your PR${NEWLINE}"
    else
      comment+="All tests passed${NEWLINE}"
    fi
  fi

  comment+="View the [build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/build-log/recording_test.log) or the [debug log](https://console.cloud.google.com/storage/browser/ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/recording) for each test"
  add_comment "${comment}"

else
  if [[ $test_exit_code -ne 0 ]]; then
    # check for any uncaught errors errors in REPLAYING mode
    comment+="Errors occurred during REPLAYING mode. Please fix them to complete your PR${NEWLINE}"
  else
    comment+="All tests passed in REPLAYING mode${NEWLINE}"
  fi
  comment+="View the [build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/build-log/replaying_test$test_suffix.log)"
  add_comment "${comment}"
fi


if [[ $test_exit_code -ne 0 ]]; then
  test_state="failure"
else
  test_state="success"
fi

set -e

update_status ${test_state}