#!/bin/bash

# CI sets the contents of our json account secret in our environment; dump it
# to disk for use in tests.
echo "${SERVICE_ACCOUNT_KEY}" > /tmp/google-account.json
echo "${ANSIBLE_TEMPLATE}" > /tmp/ansible-template.ini

set -e
set -x

# Install ansible from source
git clone https://github.com/ansible/ansible.git
pushd ansible
pip install -r requirements.txt
source hacking/env-setup
popd

# Build newest modules
pushd magic-modules-gcp
bundle install
bundle exec compiler -a -e ansible -o build/ansible

# Install collection
pushd magic-modules-gcp/build/ansible
ansible-galaxy build .
ansible-galaxy install ~/.ansible/collections
popd

# Setup Cloud configuration template with variables
pushd ~/.ansible/collections/ansible_collections/google/cloud
cp /tmp/ansible-template.ini tests/integration/cloud-config-gcp.ini

# Run ansible
ansible-test integration -v --allow-unsupported --continue-on-error $(find test/integration/targets -name "gcp*" -type d -printf "%P ")
