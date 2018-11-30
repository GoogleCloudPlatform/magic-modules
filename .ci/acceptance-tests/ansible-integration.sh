#!/bin/bash

# CI sets the contents of our json account secret in our environment; dump it
# to disk for use in tests.
echo "${SERVICE_ACCOUNT_KEY}" > /tmp/google-account.json

set -e
set -x

pushd magic-modules-new-prs/build/ansible

# Setup Cloud configuration template with variables
cp test/integration/cloud-config-gcp.yml.template test/integration/cloud-config-gcp.yml
sed -i 's/@PROJECT/graphite-test-ansible/g' test/integration/cloud-config-gcp.yml
sed -i 's/@CRED_KIND/serviceaccount/g' test/integration/cloud-config-gcp.yml
sed -i 's/@CRED_FILE/\/tmp\/google-account.json/g' test/integration/cloud-config-gcp.yml

# Setup ansible
source hacking/env-setup

# Run ansible
ansible-test integration -v --allow-unsupported --continue-on-error $(find test/integration/targets -name "gcp*" -type d -printf "%P ")

popd
exit 100
