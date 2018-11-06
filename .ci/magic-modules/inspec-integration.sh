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

git clone https://github.com/slevenick/inspec-gcp.git
pushd inspec-gcp/test/integration

# Generate tfvars
pushd attributes
ruby compile_vars.rb > terraform.tfvars
mv terraform.tfvars ../terraform
popd

# Run terraform
pushd terraform
wget https://releases.hashicorp.com/terraform/0.11.10/terraform_0.11.10_linux_amd64.zip
apt-get install unzip
unzip terraform_0.11.10_linux_amd64.zip
./terraform init
./terraform plan
./terraform apply -auto-approve
export GOOGLE_APPLICATION_CREDENTIALS="${PWD}/inspec.json"
popd

# Copy inspec resources
pushd inspec-mm
cp -r ../../../libraries libraries
popd

# Run inspec
bundle
inspec exec inspec-mm --attrs=attributes/attributes.yaml -t gcp2://

pushd terraform
./terraform destroy -auto-approve
popd
popd