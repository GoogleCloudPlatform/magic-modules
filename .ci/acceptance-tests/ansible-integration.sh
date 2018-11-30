#!/bin/bash

# CI sets the contents of our json account secret in our environment; dump it
# to disk for use in tests.
echo "${SERVICE_ACCOUNT_KEY}" > /tmp/google-account.json
echo "${ANSIBLE_TEMPLATE}" > /tmp/ansible-template.yml

set -e
set -x

pushd magic-modules-new-prs/build/ansible

# Setup Cloud configuration template with variables
cp /tmp/ansible-template.yml test/integration/cloud-config-gcp.yml

# Install dependencies for ansible
pip install -r requirements.txt

# Setup ansible
source hacking/env-setup

# Run ansible
ansible-test integration -v --allow-unsupported --continue-on-error $(find test/integration/targets -name "gcp*" -type d -printf "%P ")
