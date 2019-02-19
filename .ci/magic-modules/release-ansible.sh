#!/usr/bin/env bash

set -x

pushd "magic-modules-gcp/build/ansible"
../../tools/ansible-pr/run.sh
