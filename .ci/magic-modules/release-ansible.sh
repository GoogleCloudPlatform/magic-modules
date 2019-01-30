#!/bin/bash

# This script takes in 'magic-modules-branched', a git repo tracking the head of a PR against magic-modules.
# It outputs "ansible-generated", a non-submodule git repo containing the generated ansible code.

set -x
set -e

pushd magic-modules-gcp
./tools/ansible-pr/run
