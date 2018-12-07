#!/bin/bash

set -e
set -x

function cleanup {
	cd $TF_PATH
	terraform destroy -force -var-file=inspec-gcp.tfvars -auto-approve
}

# Service account credentials for GCP to allow terraform to work
export GOOGLE_CLOUD_KEYFILE_JSON="/tmp/google-account.json"
export GOOGLE_APPLICATION_CREDENTIALS="/tmp/google-account.json"
# Setup GOPATH
export GOPATH=${PWD}/go

# CI sets the contents of our json account secret in our environment; dump it
# to disk for use in tests.
set +x
echo "${TERRAFORM_KEY}" > /tmp/google-account.json
set -x

gcloud auth activate-service-account terraform@graphite-test-sam-chef.iam.gserviceaccount.com --key-file=$GOOGLE_CLOUD_KEYFILE_JSON

pushd magic-modules-new-prs

# Compile inspec because we are running off of new-prs
bundle install
for i in $(find products/ -name 'inspec.yaml' -printf '%h\n');
do
  bundle exec compiler -p $i -e inspec -o "build/inspec/"
done
pushd build/inspec

# Setup for using current GCP resources
export GCP_PROJECT_NUMBER=542134042613
export GCP_PROJECT_ID=graphite-test-sam-chef
export GCP_PROJECT_NAME=graphite-test-sam-chef
export GCP_ZONE=europe-west2-a
export GCP_LOCATION=europe-west2

bundle
export TF_PATH=${PWD}/test/integration/build

trap cleanup EXIT
bundle exec rake test:integration

gsutil cp inspec-cassettes/* gs://magic-modules-inspec-bucket/inspec-cassettes