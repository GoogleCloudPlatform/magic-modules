#!/bin/bash

set -e
set -x

# Service account credentials for GCP to allow terraform to work
export GOOGLE_CLOUD_KEYFILE_JSON="/tmp/google-account.json"
export GOOGLE_APPLICATION_CREDENTIALS="/tmp/google-account.json"

# CI sets the contents of our json account secret in our environment; dump it
# to disk for use in tests.
set +x
echo "${TERRAFORM_KEY}" > /tmp/google-account.json
export GCP_PROJECT_NUMBER=${PROJECT_NUMBER}
export GCP_PROJECT_ID=${PROJECT_NAME}
export GCP_PROJECT_NAME=${PROJECT_NAME}
set -x

gcloud auth activate-service-account terraform@graphite-test-sam-chef.iam.gserviceaccount.com --key-file=$GOOGLE_CLOUD_KEYFILE_JSON
PR_ID="$(cat ./magic-modules-new-prs/.git/id)"


pushd magic-modules
rm build/inspec/test/integration/verify/controls/*
export VCR_MODE=none
bundle install
bundle exec compiler -a -e inspec -o "build/inspec/"

cp templates/inspec/vcr_config.rb build/inspec

pushd build/inspec
bundle
mkdir inspec-cassettes
# Check if PR_ID folder exists
set +e
gsutil ls -m gs://magic-modules-inspec-bucket/$PR_ID
if [ $? -eq 0 ]; then
  gsutil -m cp gs://magic-modules-inspec-bucket/$PR_ID/inspec-cassettes/* inspec-cassettes/
else
  gsutil -m cp gs://magic-modules-inspec-bucket/master/inspec-cassettes/* inspec-cassettes/
fi
set -e
bundle exec rake test:generate_integration_test_variables
bundle exec rake test:run_integration_tests

popd
popd