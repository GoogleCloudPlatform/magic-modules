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

function add_comment {
  curl -H "Authorization: token ${GITHUB_TOKEN}" \
    -d "$(jq -r --arg comment "${1}" -n "{body: \$comment}")" \
    "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/${2}/comments"
}

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

set +e
# cassette retrieval
mkdir fixtures
gsutil -m -q cp gs://vcr-$GOOGLE_PROJECT/beta/fixtures/* fixtures/
# copy branch specific cassettes over master. This might fail but that's ok if the folder doesnt exist
gsutil -m -q cp gs://vcr-$GOOGLE_PROJECT/beta/refs/heads/auto-pr-$pr_number/fixtures/* fixtures/

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
export GOOGLE_APPLICATION_CREDENTIALS=$local_path/sa_key.json

echo "cassette copied"
echo "VCR_PATH: $VCR_PATH"
echo "ACCTEST_PARALLELISM: $ACCTEST_PARALLELISM" 
echo "GOOGLE_APPLICATION_CREDENTIALS: $GOOGLE_APPLICATION_CREDENTIALS"

TF_LOG=DEBUG TF_LOG_PATH_MASK=$local_path/testlog/replaying/%s.log TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test ./google-beta -parallel $ACCTEST_PARALLELISM -v '-run=TestAcc' -timeout 240m -ldflags="-X=github.com/hashicorp/terraform-provider-google/vers0ion.ProviderVersion=acc" > test.txt

test_exit_code=$?

FAILED_TESTS=$(grep "FAIL:" test.txt)
PASSED_TESTS=$(grep "PASS:" test.txt)
SKIPPED_TESTS=$(grep "SKIP:" test.txt)

FAILED_TESTS_COUNT=$(echo "$FAILED_TESTS" | wc -l)
PASSED_TESTS_COUNT=$(echo "$PASSED_TESTS" | wc -l)
SKIPPED_TESTS_COUNT=$(echo "$SKIPPED_TESTS" | wc -l)

FAILED_TESTS_PATTERN=$(grep FAIL: test.txt | awk '{print $3}' | awk -v d="|" '{s=(NR==1?s:s d)$0}END{print s}')

comment="Tests count: ${NEWLINE}"
comment+="Total tests: $(($FAILED_TESTS_COUNT+$PASSED_TESTS_COUNT+$SKIPPED_TESTS_COUNT)) ${NEWLINE}"
comment+="Passed tests $PASSED_TESTS_COUNT ${NEWLINE}"
comment+="Skipped tests: $SKIPPED_TESTS_COUNT ${NEWLINE}"
comment+="Failed tests: $FAILED_TESTS_COUNT"

add_comment "${comment}" "${pr_number}"

# store replaying build logs
# gsutil -m -q cp $local_path/testlog/replaying/* gs://replaying/test/log/path/for/each/pr #modify to correct GCS path

if [[ -n $FAILED_TESTS_PATTERN ]]; then
  
  comment="I have triggered VCR tests in RECORDING mode for the following tests that failed during VCR: $FAILED_TESTS_PATTERN"
  add_comment "${comment}" "${pr_number}"

  # RECORDING mode
  export VCR_MODE=RECORDING
  TF_LOG=DEBUG TF_LOG_PATH_MASK=$local_path/testlog/recording/%s.log TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test ./google-beta -parallel $ACCTEST_PARALLELISM -v '-run='$FAILED_TESTS_PATTERN -timeout 240m -ldflags="-X=github.com/hashicorp/terraform-provider-google/version.ProviderVersion=acc" > recording_test.txt
  test_exit_code=$?

  RECORDING_FAILED_TESTS=$(grep "FAIL:" recording_test.txt | awk '{print $3}')
  RECORDING_PASSED_TESTS=$(grep "PASS:" recording_test.txt | awk '{print $3}')

  comment=""  
  if [[ -n $RECORDING_PASSED_TESTS ]]; then
    comment="Tests passed during RECORDING mode:${NEWLINE} $RECORDING_PASSED_TESTS ${NEWLINE}${NEWLINE}"
  fi

  if [[ -n $RECORDING_FAILED_TESTS ]]; then
    comment+="Tests failed during RECORDING mode:${NEWLINE} $RECORDING_FAILED_TESTS ${NEWLINE}"
    comment+="Please fix these to complete your PR"
  else
    comment+="All tests passed"
  fi

  add_comment "${comment}" ${pr_number}

  # store cassettes
  gsutil -m -q cp fixtures/* gs://vcr-$GOOGLE_PROJECT/beta/refs/heads/auto-pr-$pr_number/fixtures/

  # store recording build logs
  # gsutil -m -q cp $local_path/testlog/recording/* gs://recording/test/log/path/for/each/pr #modify to correct GCS path

fi


if [ $test_exit_code -ne 0 ]; then
    test_state="failure"
else
    test_state="success"
fi


set -e

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



