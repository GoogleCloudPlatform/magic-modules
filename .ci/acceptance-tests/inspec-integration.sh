#!/bin/bash

set -e
set -x

# TODO make this work
# Service account credentials for GCP to allow terraform to work
export GOOGLE_CREDENTIALS_FILE="/tmp/google-account.json"
# Setup GOPATH
export GOPATH=${PWD}/go

# CI sets the contents of our json account secret in our environment; dump it
# to disk for use in tests.
echo "${GOOGLE_JSON_ACCOUNT}" > /tmp/google-account.json

git clone https://github.com/modular-magician/inspec-gcp.git
pushd inspec/test/integration

# Generate tfvars
pushd attributes
ruby compile_vars.rb > terraform.tfvars
mv terraform.tfvars ../terraform
popd

# Run terraform
pushd terraform
terraform plan
terraform apply -auto-approve
export GOOGLE_APPLICATION_CREDENTIALS="${PWD}/inspec.json"
popd

# Copy inspec resources
pushd inspec-mm
cp -r ../../../libraries libraries
popd

# Run inspec
rbenv exec inspec exec inspec-mm --attrs=attributes/attributes.yaml -t gcp2://

terraform destroy -auto-approve
popd