#!/bin/bash

# This script takes in 'magic-modules', a git repo tracking the head of a PR against magic-modules.
# It needs to output the same git repo, but with the code generation done, at 'mm-output'.

# Setup GOPATH
export GOPATH="${PWD}/go"
WORKDIR="${PWD}"

set -x
set -e

# Create $GOPATH structure - in order to successfully run Terraform codegen, we need to run
# it with a correctly-set-up $GOPATH.  It calls out to `goimports`, which means that
# we need to have all the dependencies correctly downloaded.
mkdir -p "${GOPATH}/src/github.com/terraform-providers"

pushd magic-modules
git submodule update --init build/terraform
ln -s "${PWD}/build/terraform/" "${GOPATH}/src/github.com/terraform-providers/terraform-provider-google"
popd

pushd "${GOPATH}/src/github.com/terraform-providers/terraform-provider-google"

go get

popd

pushd magic-modules
# We're going to use the short commit sha of the git repo's head as the branch name for the generated code.
BRANCH=$(git rev-parse --short HEAD)

bundle install
bundle exec compiler -p products/compute -e terraform -o build/terraform

pushd "build/terraform"
git add -A
git config --global user.email "magic-modules@google.com"
git config --global user.name "Modular Magician"
git commit -m "magic modules change happened here" || true  # don't crash if no changes
# TODO(@ndmckinley): A better message that comes from the body of the magic-modules PR.
git checkout -B "$BRANCH"
popd

git config -f .gitmodules submodule.build/terraform.branch "$BRANCH"
git config -f .gitmodules submodule.build/terraform.url "git@github.com:$GH_USERNAME/terraform-provider-google.git"
git submodule sync build/terraform

# ./branchname is intentionally not committed - but run *before* the commit, because it should contain the hash of
# the commit which kicked off this process, *not* the resulting commit.
echo "$BRANCH" > ./branchname

git add build/terraform
git add .gitmodules
git commit -m "update terraform." || true  # don't crash if no changes
# TODO(@ndmckinley): A better message that comes from the body of the magic-modules PR.
git checkout -B "$BRANCH"

cp -r ./ "${WORKDIR}/mm-output/"
popd
