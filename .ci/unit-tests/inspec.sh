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

pushd magic-modules-gcp
rm "build/inspec/test/integration/verify/controls/*"
export VCR_MODE=none
bundle exec compiler -p $i -e inspec -o "build/inspec/"

cp templates/inspec/vcr_config.rb build/inspec

pushd build/inspec
bundle
bundle exec rake test:run_integration_tests

popd
popd