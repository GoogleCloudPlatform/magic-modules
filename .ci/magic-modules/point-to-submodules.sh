#!/bin/bash

# This script takes in 'magic-modules-branched', a git repo tracking the head of a PR against magic-modules.
# It needs to output the same git repo, but with the code generation done and submodules updated, at 'magic-modules-submodules'.

set -e
set -x

pushd magic-modules-branched
BRANCH="$(cat ./branchfile)"
# Update this repo to track the submodules we just pushed:
git config -f .gitmodules submodule.build/terraform.branch "$BRANCH"
git config -f .gitmodules submodule.build/terraform.url "git@github.com:$GH_USERNAME/terraform-provider-google.git"

# Download those submodules so we can add them here.
git submodule update --remote --init build/terraform

# Commit those changes so that they can be tested in the next phase.
git add build/terraform
git add .gitmodules
git config --global user.email "magic-modules@google.com"
git config --global user.name "Modular Magician"
git commit -m "Automatic submodule update to generated code." || true  # don't crash if no changes
git checkout -B "$BRANCH"

cp -r ./ ../magic-modules-submodules
