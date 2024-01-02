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
old_branch="auto-pr-$pr_number-old"
git_remote=https://github.com/$github_username/$gh_repo
local_path=$GOPATH/src/github.com/hashicorp/$gh_repo
mkdir -p "$(dirname $local_path)"
git clone $git_remote $local_path --branch $new_branch --depth 2
pushd $local_path

# Only skip tests if we can tell for sure that no go files were changed
echo "Checking for modified go files"
# Fetch the latest commit in the old branch, associating them locally
# This will let us compare the old and new branch by name on the next line
git fetch origin $old_branch:$old_branch --depth 1
# get the names of changed files and look for go files
# (ignoring "no matches found" errors from grep)
# If there was no code generated, this will always return nothing (because there's no diff)
gofiles=$(git diff $new_branch $old_branch --name-only | { grep -e "\.go$" -e "go.mod$" -e "go.sum$" || test $? = 1; })
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

# for backwards-compatibility
if [ -z "$BASE_BRANCH" ]; then
  BASE_BRANCH=main
else
  echo "BASE_BRANCH: $BASE_BRANCH"
fi

set +e
# cassette retrieval
mkdir fixtures
if [ "$BASE_BRANCH" != "FEATURE-BRANCH-major-release-5.0.0" ]; then
  # pull main cassettes (major release uses branch specific casssettes as primary ones)
  gsutil -m -q cp gs://ci-vcr-cassettes/beta/fixtures/* fixtures/
fi
if [ "$BASE_BRANCH" != "main" ]; then
  # copy feature branch specific cassettes over main. This might fail but that's ok if the folder doesnt exist
  gsutil -m -q cp gs://ci-vcr-cassettes/beta/refs/branches/$BASE_BRANCH/fixtures/* fixtures/
fi
# copy PR branch specific cassettes over main. This might fail but that's ok if the folder doesnt exist
gsutil -m -q cp gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-$pr_number/fixtures/* fixtures/

echo $SA_KEY > sa_key.json
gcloud auth activate-service-account $GOOGLE_SERVICE_ACCOUNT --key-file=$local_path/sa_key.json --project=$GOOGLE_PROJECT

mkdir testlog
mkdir testlog/replaying
mkdir testlog/recording
mkdir testlog/recording_build
mkdir testlog/replaying_after_recording
mkdir testlog/replaying_build_after_recording

export GOOGLE_REGION=us-central1
export GOOGLE_ZONE=us-central1-a
export VCR_PATH=$local_path/fixtures
export VCR_MODE=REPLAYING
export ACCTEST_PARALLELISM=32
export GOOGLE_CREDENTIALS=$SA_KEY
export GOOGLE_APPLICATION_CREDENTIALS=$local_path/sa_key.json
export GOOGLE_TEST_DIRECTORY=$(go list ./... | grep -v github.com/hashicorp/terraform-provider-google-beta/scripts)

echo "checking terraform version"
terraform version

go build $GOOGLE_TEST_DIRECTORY
if [ $? != 0 ]; then
  echo "Skipping tests: Build failure detected"
  exit 1
fi

update_status "pending"

run_full_VCR=false

# declare an associative array ("hashmap") to track affected service packages
declare -A affected_services

for file in $gofiles
do
  if [[ $file = google-beta/services* ]]; then
    # $file should be in format 'google-beta/service/SERVICE_NAME'
    # $(echo "$file" | awk -F / '{ print $3 }') is to get the service package name
    # separate the string with '/' and get the third part
    affected_services[$(echo "$file" | awk -F / '{ print $3 }')]=1
  elif [[ $file = google-beta/provider/provider_mmv1_resources.go ]] || [[ $file = google-beta/provider/provider_dcl_resources.go ]]; then
    echo "ignore changes in $file"
  else
    run_full_VCR=true
    echo "run full tests $file"
    break
  fi

done

test_exit_code=0

affected_services_comment="None"

if [[ "$run_full_VCR" = true ]]; then
  echo "run full VCR tests"
  affected_services_comment="all service packages are affected"
  TF_LOG=DEBUG TF_LOG_PATH_MASK=$local_path/testlog/replaying/%s.log TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test $GOOGLE_TEST_DIRECTORY -parallel $ACCTEST_PARALLELISM -v -run=TestAcc -timeout 240m -ldflags="-X=github.com/hashicorp/terraform-provider-google-beta/version.ProviderVersion=acc" > replaying_test.log # write log into file

  test_exit_code=$?
else
  # clear GOOGLE_TEST_DIRECTORY
  GOOGLE_TEST_DIRECTORY=""
  affected_services_comment="<ul>"
  for service in "${!affected_services[@]}"
  do
    # append affected service package path
    GOOGLE_TEST_DIRECTORY+=" ./google-beta/services/$service"
    echo "run VCR tests in $service"
    TF_LOG=DEBUG TF_LOG_PATH_MASK=$local_path/testlog/replaying/%s.log TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test ./google-beta/services/$service -parallel $ACCTEST_PARALLELISM -v -run=TestAcc -timeout 240m -ldflags="-X=github.com/hashicorp/terraform-provider-google-beta/version.ProviderVersion=acc" >> replaying_test.log # append logs into file

    test_exit_code=$(($test_exit_code || $?))
    affected_services_comment+="<li>$service</li>"
  done
  affected_services_comment+="</ul>"
fi

# store replaying build log
gsutil -h "Content-Type:text/plain" -q cp replaying_test.log gs://ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/build-log/

# store replaying test logs
gsutil -h "Content-Type:text/plain" -m -q cp testlog/replaying/* gs://ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/replaying/

# handle provider crash
TESTS_PANIC=$(grep "^panic: " replaying_test.log)

if [[ -n $TESTS_PANIC ]]; then
  comment="$\textcolor{red}{\textsf{The provider crashed while running the VCR tests in REPLAYING mode}}$ ${NEWLINE}"
  comment+="$\textcolor{red}{\textsf{Please fix it to complete your PR}}$ ${NEWLINE}"
  comment+="View the [build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/build-log/replaying_test.log)"
  add_comment "${comment}"
  update_status "failure"
  exit 0
fi

FAILED_TESTS=$(grep "^--- FAIL: TestAcc" replaying_test.log)
PASSED_TESTS=$(grep "^--- PASS: TestAcc" replaying_test.log)
SKIPPED_TESTS=$(grep "^--- SKIP: TestAcc" replaying_test.log)

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

FAILED_TESTS_PATTERN=$(grep "^--- FAIL: TestAcc" replaying_test.log | awk '{print $3}' | awk -v d="|" '{s=(NR==1?s:s d)$0}END{print s}')

comment="#### Tests analytics ${NEWLINE}"
comment+="Total tests: \`$(($FAILED_TESTS_COUNT+$PASSED_TESTS_COUNT+$SKIPPED_TESTS_COUNT))\` ${NEWLINE}"
comment+="Passed tests \`$PASSED_TESTS_COUNT\` ${NEWLINE}"
comment+="Skipped tests: \`$SKIPPED_TESTS_COUNT\` ${NEWLINE}"
comment+="Affected tests: \`$FAILED_TESTS_COUNT\` ${NEWLINE}${NEWLINE}"
comment+="<details><summary>Click here to see the affected service packages</summary><blockquote>$affected_services_comment</blockquote></details> ${NEWLINE}${NEWLINE}"

if [[ -n $FAILED_TESTS_PATTERN ]]; then
  comment+="#### Action taken ${NEWLINE}"
  comment+="<details> <summary>Found $FAILED_TESTS_COUNT affected test(s) by replaying old test recordings. Starting RECORDING based on the most recent commit. Click here to see the affected tests</summary><blockquote>$FAILED_TESTS_PATTERN </blockquote></details> ${NEWLINE}${NEWLINE}"
  comment+="[Get to know how VCR tests work](https://googlecloudplatform.github.io/magic-modules/docs/getting-started/contributing/#general-contributing-steps)"
  add_comment "${comment}"
  # Clear fixtures folder
  rm $VCR_PATH/*

  # Clear replaying-log folder
  rm testlog/replaying/*
  
  # RECORDING mode
  export VCR_MODE=RECORDING
  FAILED_TESTS=$(grep "^--- FAIL: TestAcc" replaying_test.log | awk '{print $3}')
  # test_exit_code=0
  parallel --jobs 16 TF_LOG=DEBUG TF_LOG_PATH_MASK=$local_path/testlog/recording/%s.log TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test {1} -parallel 1 -v -run="{2}$" -timeout 240m -ldflags="-X=github.com/hashicorp/terraform-provider-google-beta/version.ProviderVersion=acc" ">>" testlog/recording_build/{2}_recording_test.log ::: $GOOGLE_TEST_DIRECTORY ::: $FAILED_TESTS

  test_exit_code=$?

  # Concatenate recording build logs to one file
  # Note: build logs are different from debug logs
  for failed_test in $FAILED_TESTS
  do
    cat testlog/recording_build/${failed_test}_recording_test.log >> recording_test.log
  done

  # store cassettes
  gsutil -m -q cp fixtures/* gs://ci-vcr-cassettes/beta/refs/heads/auto-pr-$pr_number/fixtures/

  # store recording build log
  gsutil -h "Content-Type:text/plain" -q cp recording_test.log gs://ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/build-log/

  # store recording individual build logs
  gsutil -h "Content-Type:text/plain" -m -q cp testlog/recording_build/* gs://ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/build-log/recording_build/

  # store recording test logs
  gsutil -h "Content-Type:text/plain" -m -q cp testlog/recording/* gs://ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/recording/

  # handle provider crash
  RECORDING_TESTS_PANIC=$(grep "^panic: " recording_test.log)

  if [[ -n $RECORDING_TESTS_PANIC ]]; then
  
    comment="$\textcolor{red}{\textsf{The provider crashed while running the VCR tests in RECORDING mode}}$ ${NEWLINE}"
    comment+="$\textcolor{red}{\textsf{Please fix it to complete your PR}}$ ${NEWLINE}"
    comment+="View the [build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/build-log/recording_test.log)"
    add_comment "${comment}"
    update_status "failure"
    exit 0
  fi


  RECORDING_FAILED_TESTS=$(grep "^--- FAIL: TestAcc" recording_test.log | awk -v pr_number=$pr_number -v build_id=$build_id '{print "`"$3"`[[Error message](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-"pr_number"/artifacts/"build_id"/build-log/recording_build/"$3"_recording_test.log)] [[Debug log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-"pr_number"/artifacts/"build_id"/recording/"$3".log)]"}')
  RECORDING_PASSED_TESTS=$(grep "^--- PASS: TestAcc" recording_test.log | awk -v pr_number=$pr_number -v build_id=$build_id '{print "`"$3"`[[Debug log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-"pr_number"/artifacts/"build_id"/recording/"$3".log)]"}')
  RECORDING_PASSED_TEST_LIST=$(grep "^--- PASS: TestAcc" recording_test.log | awk '{print $3}')

  comment=""
  RECORDING_PASSED_TESTS_COUNT=0
  RECORDING_FAILED_TESTS_COUNT=0
  if [[ -n $RECORDING_PASSED_TESTS ]]; then
    comment+="$\textcolor{green}{\textsf{Tests passed during RECORDING mode:}}$ ${NEWLINE} $RECORDING_PASSED_TESTS ${NEWLINE}${NEWLINE}"
    RECORDING_PASSED_TESTS_COUNT=$(echo "$RECORDING_PASSED_TESTS" | wc -l)
    comment+="##### Rerun these tests in REPLAYING mode to catch issues ${NEWLINE}${NEWLINE}"

    # Rerun passed tests in REPLAYING mode 3 times to catch issues
    export VCR_MODE=REPLAYING
    count=3
    parallel --jobs 16 TF_LOG=DEBUG TF_LOG_PATH_MASK=$local_path/testlog/replaying_after_recording/%s.log TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test {1} -parallel 1 -count=$count -v -run="{2}$" -timeout 120m -ldflags="-X=github.com/hashicorp/terraform-provider-google-beta/version.ProviderVersion=acc" ">>" testlog/replaying_build_after_recording/{2}_replaying_test.log ::: $GOOGLE_TEST_DIRECTORY ::: $RECORDING_PASSED_TEST_LIST

    test_exit_code=$(($test_exit_code || $?))

    # Concatenate recording build logs to one file
    for test in $RECORDING_PASSED_TEST_LIST
    do
      cat testlog/replaying_build_after_recording/${test}_replaying_test.log >> replaying_build_after_recording.log
    done

    # store replaying individual build logs
    gsutil -h "Content-Type:text/plain" -m -q cp testlog/replaying_build_after_recording/* gs://ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/build-log/replaying_build_after_recording/

    # store replaying test logs
    gsutil -h "Content-Type:text/plain" -m -q cp testlog/replaying_after_recording/* gs://ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/replaying_after_recording/

    REPLAYING_FAILED_TESTS=$(grep "^--- FAIL: TestAcc" replaying_build_after_recording.log | sort -u -t' ' -k3,3 | awk -v pr_number=$pr_number -v build_id=$build_id '{print "`"$3"`[[Error message](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-"pr_number"/artifacts/"build_id"/build-log/replaying_build_after_recording/"$3"_replaying_test.log)] [[Debug log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-"pr_number"/artifacts/"build_id"/replaying_after_recording/"$3".log)]"}')
    if [[ -n $REPLAYING_FAILED_TESTS ]]; then
      comment+="$\textcolor{red}{\textsf{Tests failed when rerunning REPLAYING mode:}}$ ${NEWLINE} $REPLAYING_FAILED_TESTS ${NEWLINE}${NEWLINE}"
      comment+="Tests failed due to non-determinism or randomness when the VCR replayed the response after the HTTP request was made.${NEWLINE}${NEWLINE}"
      comment+="Please fix these to complete your PR. If you believe these test failures to be incorrect or unrelated to your change, or if you have any questions, please raise the concern with your reviewer.${NEWLINE}"
    else
      comment+="$\textcolor{green}{\textsf{No issues found for passed tests after REPLAYING rerun.}}$ ${NEWLINE}"
    fi
    comment+="${NEWLINE}---${NEWLINE}"

    # Clear replaying-log folder
    rm testlog/replaying_after_recording/*
    rm testlog/replaying_build_after_recording/*
  fi

  if [[ -n $RECORDING_FAILED_TESTS ]]; then
    comment+="$\textcolor{red}{\textsf{Tests failed during RECORDING mode:}}$ ${NEWLINE} $RECORDING_FAILED_TESTS ${NEWLINE}${NEWLINE}"
    RECORDING_FAILED_TESTS_COUNT=$(echo "$RECORDING_FAILED_TESTS" | wc -l)
    if [[ $RECORDING_PASSED_TESTS_COUNT+$RECORDING_FAILED_TESTS_COUNT -lt $FAILED_TESTS_COUNT ]]; then
      test_exit_code=1
      comment+="$\textcolor{red}{\textsf{Several tests got terminated during RECORDING mode.}}$ ${NEWLINE}"
    fi
    comment+="$\textcolor{red}{\textsf{Please fix these to complete your PR.}}$ ${NEWLINE}"
  else
    if [[ $RECORDING_PASSED_TESTS_COUNT+$RECORDING_FAILED_TESTS_COUNT -lt $FAILED_TESTS_COUNT ]]; then
      test_exit_code=1
      comment+="$\textcolor{red}{\textsf{Several tests got terminated during RECORDING mode.}}$ ${NEWLINE}"
    elif [[ $test_exit_code -ne 0 ]]; then
      # check for any uncaught errors in RECORDING mode
      comment+="$\textcolor{red}{\textsf{Errors occurred during RECORDING mode. Please fix them to complete your PR.}}$ ${NEWLINE}"
    else
      comment+="$\textcolor{green}{\textsf{All tests passed!}}$ ${NEWLINE}"
    fi
  fi

  comment+="View the [build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/build-log/recording_test.log) or the [debug log](https://console.cloud.google.com/storage/browser/ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/recording) for each test"
  add_comment "${comment}"

else
  if [[ $test_exit_code -ne 0 ]]; then
    # check for any uncaught errors errors in REPLAYING mode
    comment+="$\textcolor{red}{\textsf{Errors occurred during REPLAYING mode. Please fix them to complete your PR}}$ ${NEWLINE}"
  else
    comment+="$\textcolor{green}{\textsf{All tests passed in REPLAYING mode.}}$ ${NEWLINE}"
  fi
  comment+="View the [build log](https://storage.cloud.google.com/ci-vcr-logs/beta/refs/heads/auto-pr-$pr_number/artifacts/$build_id/build-log/replaying_test.log)"
  add_comment "${comment}"
fi


if [[ $test_exit_code -ne 0 ]]; then
  test_state="failure"
else
  test_state="success"
fi

set -e

update_status ${test_state}