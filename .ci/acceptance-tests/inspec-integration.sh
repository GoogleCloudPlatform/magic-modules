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

export CLOUD_SDK_REPO="cloud-sdk-stretch"
echo "deb http://packages.cloud.google.com/apt $CLOUD_SDK_REPO main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list
curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
apt-get update && apt-get install google-cloud-sdk -y

gcloud auth activate-service-account terraform@graphite-test-sam-chef.iam.gserviceaccount.com --key-file=$GOOGLE_CLOUD_KEYFILE_JSON

# Download train plugin (it's not published yet)
gsutil cp -r gs://magic-modules-inspec-bucket/train-gcp2 .
gem install inspec
inspec plugin install train-gcp2/lib/train-gcp2.rb

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
curl https://releases.hashicorp.com/terraform/0.11.10/terraform_0.11.10_linux_amd64.zip > terraform_0.11.10_linux_amd64.zip
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

for i in {1..30}
do
	# Cleanup cassettes folder each time, we don't want to use a recorded cassette if it records an unauthorized response
	rm -r inspec-cassettes
	inspec exec inspec-mm --attrs=attributes/attributes.yaml -t gcp2://
	if [ "$?" -eq "0" ]; then
		# Upload cassettes to storage bucket for unit test use
		gsutil cp inspec-cassettes/* gs://magic-modules-inspec-bucket/inspec-cassettes
		exit 0
	fi
done
set -e

popd
popd
exit 100