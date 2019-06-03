#!/bin/bash

set -e
set -x

# Service account credentials for GCP to allow terraform to work
export GOOGLE_CLOUD_KEYFILE_JSON="/tmp/google-account.json"
export GOOGLE_APPLICATION_CREDENTIALS="/tmp/google-account.json"

# CI sets the contents of our json account secret in our environment; dump it
# to disk for use in tests.
set +x
echo "${SERVICE_ACCOUNT}" > /tmp/google-account.json
set -x

gcloud auth activate-service-account  magic-modules-spreadsheet@magic-modules.iam.gserviceaccount.com --key-file=$GOOGLE_CLOUD_KEYFILE_JSON

pushd magic-modules-gcp
bundle install
gem install rspec

# || true will suppress errors, but it's necessary for this to run. If unset,
# Concourse will fail on *any* rspec step failing (eg: any API mismatch)
bundle exec rspec tools/linter/spreadsheet.rb  || true

echo "File created"
date=$(date +'%m%d%Y')
echo "Date established"

gsutil cp output.csv gs://magic-modules-coverage/$date.csv
popd
