#!/usr/bin/env bash

set -e
set -x

pushd "magic-modules-gcp/build/ansible"
git remote add origin git@github.com:modular-magician/ansible.git

../../tools/ansible-pr/run.sh
