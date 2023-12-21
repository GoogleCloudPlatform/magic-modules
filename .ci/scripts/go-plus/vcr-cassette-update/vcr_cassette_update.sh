#!/bin/bash

set -e

build_id=$1
github_username=hashicorp
gh_repo=terraform-provider-google-beta
git_remote=https://github.com/$github_username/$gh_repo
local_path=$GOPATH/src/github.com/hashicorp/$gh_repo
mkdir -p "$(dirname $local_path)"
git clone $git_remote $local_path --depth 1
pushd $local_path

# get today's date in YYYY-MM-DD format
today=$(date +%F)

set +e

echo $SA_KEY > sa_key.json
gcloud auth activate-service-account $GOOGLE_SERVICE_ACCOUNT --key-file=$local_path/sa_key.json --project=$GOOGLE_PROJECT

# cassette retrieval
mkdir fixtures

# pull main cassettes
gsutil -m -q cp gs://ci-vcr-cassettes/beta/fixtures/* fixtures/

# main cassettes backup
# incase nightly run goes wrong. this will be used to restore the cassettes
gsutil -m -q cp fixtures/* gs://vcr-nightly/beta/$today/$build_id/main_cassettes_backup/fixtures/

mkdir testlog
mkdir testlog/replaying
mkdir testlog/recording
mkdir testlog/recording_build

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

echo "running tests in REPLAYING mode now"
TF_LOG=DEBUG TF_LOG_PATH_MASK=$local_path/testlog/replaying/%s.log TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test $GOOGLE_TEST_DIRECTORY -parallel $ACCTEST_PARALLELISM -v -run=TestAcc -timeout 240m -ldflags="-X=github.com/hashicorp/terraform-provider-google-beta/version.ProviderVersion=acc" > replaying_test.log

# store replaying build log
gsutil -h "Content-Type:text/plain" -q cp replaying_test.log gs://vcr-nightly/beta/$today/$build_id/logs/build-log/

# store replaying test logs
gsutil -h "Content-Type:text/plain" -m -q cp testlog/replaying/* gs://vcr-nightly/beta/$today/$build_id/logs/replaying/

# handle provider crash
TESTS_PANIC=$(grep "^panic: " replaying_test.log)
if [[ -n $TESTS_PANIC ]]; then
  echo "#################################"
  echo "The provider crashed while running the VCR tests in REPLAYING mode"
  echo "#################################"
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

echo "#################################"
echo "Tests Analytics"
echo "Total tests: $(($FAILED_TESTS_COUNT+$PASSED_TESTS_COUNT+$SKIPPED_TESTS_COUNT))"
echo "Passed tests: $PASSED_TESTS_COUNT"
echo "Skipped tests: $SKIPPED_TESTS_COUNT"
echo "Affected tests: $FAILED_TESTS_COUNT"
echo "Affected tests list: $FAILED_TESTS_PATTERN"
echo "#################################"
echo ""

if [[ -n $FAILED_TESTS_PATTERN ]]; then
  echo "running affected tests in RECORDING mode now"

  # Clear fixtures folder
  rm $VCR_PATH/*

  # set RECORDING mode
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
  gsutil -m -q cp fixtures/* gs://ci-vcr-cassettes/beta/fixtures/

  # store recording build log
  gsutil -h "Content-Type:text/plain" -q cp recording_test.log gs://vcr-nightly/beta/$today/$build_id/logs/build-log/

  # store recording individual build logs
  gsutil -h "Content-Type:text/plain" -m -q cp testlog/recording_build/* gs://vcr-nightly/beta/$today/$build_id/logs/build-log/recording_build/

  # store recording test logs
  gsutil -h "Content-Type:text/plain" -m -q cp testlog/recording/* gs://vcr-nightly/beta/$today/$build_id/logs/recording/

  # handle provider crash
  RECORDING_TESTS_PANIC=$(grep "^panic: " recording_test.log)

  if [[ -n $RECORDING_TESTS_PANIC ]]; then
    echo "#################################"
    echo "The provider crashed while running the VCR tests in RECORDING mode"
    echo "#################################"
    exit 0
  fi

  RECORDING_FAILED_TESTS=$(grep "^--- FAIL: TestAcc" recording_test.log | awk '{print $3}')
  RECORDING_PASSED_TESTS=$(grep "^--- PASS: TestAcc" recording_test.log | awk '{print $3}')

  RECORDING_PASSED_TESTS_COUNT=0
  RECORDING_FAILED_TESTS_COUNT=0

  echo "#################################"
  echo "RECORDING Tests Report"
  if [[ -n $RECORDING_PASSED_TESTS ]]; then
    RECORDING_PASSED_TESTS_COUNT=$(echo "$RECORDING_PASSED_TESTS" | wc -l)
    echo "Tests passed during RECORDING mode:"
    echo $RECORDING_PASSED_TESTS
    echo ""
  fi

  if [[ -n $RECORDING_FAILED_TESTS ]]; then
    RECORDING_FAILED_TESTS_COUNT=$(echo "$RECORDING_FAILED_TESTS" | wc -l)
    echo "Tests failed during RECORDING mode:"
    echo $RECORDING_FAILED_TESTS
    echo ""
    if [[ $RECORDING_PASSED_TESTS_COUNT+$RECORDING_FAILED_TESTS_COUNT -lt $FAILED_TESTS_COUNT ]]; then
      echo "Several tests got terminated during RECORDING mode"
    fi
  else
    if [[ $RECORDING_PASSED_TESTS_COUNT+$RECORDING_FAILED_TESTS_COUNT -lt $FAILED_TESTS_COUNT ]]; then
      echo "Several tests got terminated during RECORDING mode"
    elif [[ $test_exit_code -ne 0 ]]; then
      # check for any uncaught errors in RECORDING mode
      echo "Errors occurred during RECORDING mode."
    else
      echo "All tests passed!"
    fi
  fi
  echo "#################################"
else
  if [[ $test_exit_code -ne 0 ]]; then
    # check for any uncaught errors errors in REPLAYING mode
    echo "#################################"
    echo "Errors occurred during REPLAYING mode."
    echo "#################################"
  else
    echo "#################################"
    echo "All tests passed in REPLAYING mode."
    echo "#################################"
  fi
fi

set -e