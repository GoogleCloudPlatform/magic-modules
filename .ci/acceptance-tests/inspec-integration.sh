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

pushd magic-modules-gcp/build/inspec

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
bundle exec rake test:integration
popd