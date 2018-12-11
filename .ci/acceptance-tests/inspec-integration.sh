#!/bin/bash

set -e
set -x

function cleanup {
	cd $TF_PATH
	terraform destroy -auto-approve
}

# Service account credentials for GCP to allow terraform to work
export GOOGLE_CLOUD_KEYFILE_JSON="/tmp/google-account.json"
# Setup GOPATH
export GOPATH=${PWD}/go

# CI sets the contents of our json account secret in our environment; dump it
# to disk for use in tests.
set +x
echo "${TERRAFORM_KEY}" > /tmp/google-account.json
set -x

pushd magic-modules-new-prs

# Compile inspec because we are running off of new-prs
bundle install
for i in $(find products/ -name 'inspec.yaml' -printf '%h\n');
do
  bundle exec compiler -p $i -e inspec -o "build/inspec/"
done
pushd build/inspec/test/integration

# Generate tfvars
pushd attributes
ruby compile_vars.rb > terraform.tfvars
mv terraform.tfvars ../terraform
popd

# Run terraform
pushd terraform
terraform init
terraform plan

export TF_PATH=${PWD}
trap cleanup EXIT
terraform apply -auto-approve
export GOOGLE_APPLICATION_CREDENTIALS="${PWD}/inspec.json"
popd

# Run inspec
bundle

# Service accounts take several minutes to be authorized everywhere
set +e

for i in {1..30}
do
	# Cleanup cassettes folder each time, we don't want to use a recorded cassette if it records an unauthorized response
	rm -r inspec-cassettes
	
	if inspec exec verify-mm --attrs=attributes/attributes.yaml -t gcp:// --no-distinct-exit; then
		exit 0
	fi
done
set -e

popd
popd
exit 100