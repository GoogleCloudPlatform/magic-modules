#!/bin/bash

# CI sets the contents of our json account secret in our environment; dump it
# to disk for use in tests.
echo "${SERVICE_ACCOUNT_KEY}" > /tmp/google-account.json
echo "${ANSIBLE_TEMPLATE}" > /tmp/ansible-template.ini

set -e
set -x

# Get the newest version of Ansible from the PR
pushd magic-modules-gcp
bundle install
for i in $(find products/ -name 'ansible.yaml' -printf '%h\n');
do
  bundle exec compiler -p $i -e ansible -o "build/ansible/"
done
popd

# Go to the newly-compiled version of Ansible
pushd magic-modules-gcp/build/ansible

# Setup Cloud configuration template with variables
cp /tmp/ansible-template.ini test/integration/cloud-config-gcp.ini

# Install dependencies for ansible
pip install -r requirements.txt

# Setup ansible
source hacking/env-setup

# Run ansible
ansible-test integration -v --allow-unsupported --continue-on-error $(find test/integration/targets -name "gcp*" -type d -printf "%P ")
