#!/bin/bash

# This script takes in 'magic-modules-branched', a git repo tracking the head of a PR against magic-modules.
# It outputs "terraform-generated", a non-submodule git repo containing the generated terraform code.

set -x
set -e

# Create $GOPATH structure - in order to successfully run Terraform codegen, we need to run
# it with a correctly-set-up $GOPATH.  It calls out to `goimports`, which means that
# we need to have all the dependencies correctly downloaded.
export GOPATH="${PWD}/go"
mkdir -p "${GOPATH}/src/github.com/terraform-providers"

pushd magic-modules-branched
git submodule update --init build/terraform
ln -s "${PWD}/build/terraform/" "${GOPATH}/src/github.com/terraform-providers/terraform-provider-google"
popd

pushd "${GOPATH}/src/github.com/terraform-providers/terraform-provider-google"
go get
popd

pushd magic-modules-branched
bundle install
bundle exec compiler -p products/compute -e terraform -o "${GOPATH}/src/github.com/terraform-providers/terraform-provider-google/"

TERRAFORM_COMMIT_MSG="$(gawk -v fetch=terraform,all -f ci/magic-modules/extract-from-pr-description.awk < .git/body)"
if [ -z "$TERRAFORM_COMMIT_MSG" ]; then
  TERRAFORM_COMMIT_MSG="Magic Modules changes"
fi

pushd "build/terraform"
git config --global user.email "magic-modules@google.com"
git config --global user.name "Modular Magician"

git add -A
git commit -m "$TERRAFORM_COMMIT_MSG" || true  # don't crash if no changes
git checkout -B "$(cat ../../branchname)"
popd

popd

git clone magic-modules-branched/build/terraform ./terraform-generated
