#!/bin/bash

set -e
set -x

function cleanup {
	cd $TF_PATH
	./terraform destroy -auto-approve
}

# Service account credentials for GCP to allow terraform to work
export GOOGLE_CLOUD_KEYFILE_JSON="/tmp/google-account.json"
# Setup GOPATH
export GOPATH=${PWD}/go

# CI sets the contents of our json account secret in our environment; dump it
# to disk for use in tests.
echo "${TERRAFORM_KEY}" > /tmp/google-account.json

git clone https://github.com/slevenick/inspec-gcp.git

# new train plugin not published yet, install locally for now
pushd inspec-gcp
bundle
inspec plugin install train-gcp2/lib/train-gcp2.rb

popd

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

export TF_PATH=${PWD}
trap cleanup EXIT
./terraform apply -auto-approve
export GOOGLE_APPLICATION_CREDENTIALS="${PWD}/inspec.json"
inspec detect -t gcp2://
popd

# Copy inspec resources
pushd inspec-mm
cp -r ../../../libraries libraries
popd

# Run inspec
bundle

# Service accounts take several minutes to be authorized everywhere
set +e

for i in {1..50}
do
	inspec exec inspec-mm --attrs=attributes/attributes.yaml -t gcp2://
	if [ "$?" -eq "0" ]; then
		exit 0
	fi
done
set -e

popd
exit 100