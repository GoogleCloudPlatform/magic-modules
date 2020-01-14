#!/bin/bash

set -e
set -x

# Service account credentials for GCP to allow terraform to work
export GOOGLE_CLOUD_KEYFILE_JSON="/tmp/google-account.json"
export GOOGLE_APPLICATION_CREDENTIALS="/tmp/google-account.json"
# Setup GOPATH
export GOPATH=${PWD}/go

# CI sets the contents of our json account secret in our environment; dump it
# to disk for use in tests.
set +x
echo "${TERRAFORM_KEY}" > /tmp/google-account.json
export GCP_PROJECT_NUMBER=${PROJECT_NUMBER}
export GCP_PROJECT_ID=${PROJECT_NAME}
export GCP_PROJECT_NAME=${PROJECT_NAME}
set -x

gcloud auth activate-service-account terraform@graphite-test-sam-chef.iam.gserviceaccount.com --key-file=$GOOGLE_CLOUD_KEYFILE_JSON
# TODO(slevenick): Check to see if we have already run this
PR_ID="$(cat ./mm-approved-prs/.git/id)"

# Check if PR_ID folder exists
set +e
gsutil ls gs://magic-modules-inspec-bucket/$PR_ID
if [ $? -ne 0 ]; then
	# Bucket does not exist, so we did not have to record new cassettes to pass the inspec-test step.
	# This means no new cassettes need to be generated after this PR is merged.
	exit 0
fi
set -e

pushd mm-approved-prs
export VCR_MODE=all
# Running other controls may cause caching issues due to underlying clients caching responses
rm build/inspec/test/integration/verify/controls/*
bundle install
bundle exec compiler -a -e inspec -o "build/inspec/"
cp templates/inspec/vcr_config.rb build/inspec

pushd build/inspec

# Setup for using current GCP resources
export GCP_ZONE=europe-west2-a
export GCP_LOCATION=europe-west2

bundle install

function cleanup {
	cd $INSPEC_DIR
	bundle exec rake test:cleanup_integration_tests
}

export INSPEC_DIR=${PWD}
trap cleanup EXIT

seed=$RANDOM
bundle exec rake test:init_workspace
# Seed plan_integration_tests so VCR cassettes work with random resource suffixes
bundle exec rake test:plan_integration_tests[$seed]
bundle exec rake test:setup_integration_tests
bundle exec rake test:run_integration_tests
bundle exec rake test:cleanup_integration_tests

echo $seed > inspec-cassettes/seed.txt

gsutil -m cp inspec-cassettes/* gs://magic-modules-inspec-bucket/$PR_ID/inspec-cassettes/approved/

popd