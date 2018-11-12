#!/usr/bin/env bash

set -e
set -x

pushd "magic-modules/build/inspec/test/integration"

# Generate a rsa private key to use in mocks
# Due to using gauth library InSpec + train expect to load a service account file from an env variable
# This service account file must contain a real RSA key, but this key is never used in unit tests.
ssh-keygen -f id_rsa -t rsa -N ''

echo '{
  "type": "service_account",
  "project_id": "fake",
  "private_key_id": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "private_key": "<%= @private_key %>",
  "client_email": "fake@fake.iam.gserviceaccount.com",
  "client_id": "123451234512345123451",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/fake%40fake.iam.gserviceaccount.com"
}' > inspec.json.erb

# Formatting a private key file for use is surprisingly difficult
echo -n "@private_key = '$(echo -n "$(cat id_rsa)")'.gsub(\"\n\", '\n')" > var.rb
erb -r './var' inspec.json.erb > inspec.json

export GOOGLE_APPLICATION_CREDENTIALS=${PWD}/inspec.json

inspec exec inspec-mm --attrs=attributes/attributes.yaml -t gcp2://
