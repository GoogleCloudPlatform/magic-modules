#!/usr/bin/env bash

set -e
set -x

pushd "magic-modules/build/ansible"
pip install tox

source hacking/env-setup
ansible-test sanity -v --tox --python 2.7 --base-branch origin/devel lib/ansible/modules/cloud/google/
