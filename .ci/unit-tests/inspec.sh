#!/usr/bin/env bash

set -e
set -x
echo 'TODO(slevenick): reimplement the following'
exit 0

# Service account credentials for GCP to pull VCR cassettes
export GOOGLE_CLOUD_KEYFILE_JSON="/tmp/google-account.json"

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

pushd "magic-modules/build/inspec/test/integration"

# Generate a rsa private key to use in mocks
# Due to using gauth library InSpec + train expect to load a service account file from an env variable
# This service account file must contain a real RSA key, but this key is never used in unit tests.
rsatmp=$(mktemp /tmp/rsatmp.XXXXXX)
yes y | ssh-keygen -f "${rsatmp}" -t rsa -N ''


echo '{
  "type": "service_account",
  "project_id": "fake",
  "private_key_id": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "private_key": "<%= @fake_private_key %>",
  "client_email": "fake@fake.iam.gserviceaccount.com",
  "client_id": "123451234512345123451",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/fake%40fake.iam.gserviceaccount.com"
}' > inspec.json.erb

# Formatting a rsa key file for use is surprisingly difficult
echo -n "@fake_private_key = '$(echo -n "$(cat ${rsatmp})")'.gsub(\"\n\", '\n')" > var.rb
rm ${rsatmp}
erb -r './var' inspec.json.erb > inspec.json

pushd inspec-mm
cp -r ../../../libraries libraries
popd

export GOOGLE_APPLICATION_CREDENTIALS=${PWD}/inspec.json

bundle install
# TODO change this to use a github repo
gsutil cp -r gs://magic-modules-inspec-bucket/inspec-cassettes .

function cleanup {
  rm -rf inspec-cassettes
  rm -rf inspec-mm/libraries
  rm inspec.json
  rm inspec.json.erb
  rm var.rb
}
trap cleanup EXIT

inspec exec inspec-mm --attrs=attributes/attributes.yaml -t gcp2:// --no-distinct-exit
