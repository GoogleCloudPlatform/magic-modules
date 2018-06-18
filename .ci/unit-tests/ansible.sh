#!/usr/bin/env bash

set -e
set -x

pushd "magic-modules/build/ansible"
source hacking/env-setup
ansible-test sanity --tox --python 2.7 $(find test/integration/targets -name "gcp_*" -printf '%P ')
