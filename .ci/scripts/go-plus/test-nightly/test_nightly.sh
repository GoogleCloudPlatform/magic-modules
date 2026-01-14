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

mkdir testlog
mkdir testlog/debug_log

export GOOGLE_REGION=us-central1
export GOOGLE_ZONE=us-central1-a
export ACCTEST_PARALLELISM=12
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

echo "running tests now"
TF_LOG=DEBUG TF_LOG_PATH_MASK=$local_path/testlog/debug_log/%s.log TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test ./google-beta/services/sql -parallel $ACCTEST_PARALLELISM -v -run=TestAccSqlDatabaseInstance_MysqlSwitchoverSuccess -timeout 1200m -ldflags="-X=github.com/hashicorp/terraform-provider-google-beta/version.ProviderVersion=acc" > test.log

# store build log
gsutil -h "Content-Type:text/plain" -q cp test.log gs://test-nightly/beta/$today/$build_id/logs/build-log/

# store test logs
gsutil -h "Content-Type:text/plain" -m -q cp testlog/debug_log/* gs://test-nightly/beta/$today/$build_id/logs/debug-log/

# handle provider crash
TESTS_PANIC=$(grep "^panic: " test.log)
if [[ -n $TESTS_PANIC ]]; then
  echo "#################################"
  echo "The provider crashed while running the tests"
  echo "#################################"
  exit 0
fi

FAILED_TESTS=$(grep "^--- FAIL: TestAcc" test.log)
PASSED_TESTS=$(grep "^--- PASS: TestAcc" test.log)
SKIPPED_TESTS=$(grep "^--- SKIP: TestAcc" test.log)

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

FAILED_TESTS_PATTERN=$(grep "^--- FAIL: TestAcc" test.log | awk '{print $3}' | awk -v d="|" '{s=(NR==1?s:s d)$0}END{print s}')

echo "#################################"
echo "Tests Analytics"
echo "Total tests: $(($FAILED_TESTS_COUNT+$PASSED_TESTS_COUNT+$SKIPPED_TESTS_COUNT))"
echo "Passed tests: $PASSED_TESTS_COUNT"
echo "Skipped tests: $SKIPPED_TESTS_COUNT"
echo "Failed tests: $FAILED_TESTS_COUNT"
echo "Failed tests list: $FAILED_TESTS_PATTERN"
echo "#################################"
echo ""

set -e