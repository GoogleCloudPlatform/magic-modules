#!/usr/bin/env bash

set -x

# Install dependencies for Template Generator
pushd "magic-modules-gcp"
bundle install
popd

pushd "magic-modules-gcp/build/ansible"
# Setup Git config.
git config --global user.email "alexstephen@google.com"
git config --global user.name "Alex Stephen"

# Run creation script.
../../tools/ansible-pr/run.sh
