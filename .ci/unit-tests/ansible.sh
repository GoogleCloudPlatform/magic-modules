#!/usr/bin/env bash

set -e
set -x

apt-get update
pushd "magic-modules/build/ansible"
apt-get install -y man
pip install tox


source hacking/env-setup
ansible-test sanity -v --tox --python 2.7 --base-branch origin/devel lib/ansible/modules/cloud/google/ lib/ansible/module_utils/gcp_utils.py
