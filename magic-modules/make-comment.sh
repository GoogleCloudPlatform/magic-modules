#!/bin/bash

# This creates a PR comment based on what PR was created downstream.

set -e
set -x

shopt -s dotglob

pushd mm-initial-pr
ID=$(git config --get pullrequest.id)
popd

pushd terraform-prs-out
TF_PR=$(git config --get pullrequest.id)
popd

cp -r magic-modules-out/* magic-modules-comment

pushd magic-modules-comment
git config pullrequest.id "$ID"
cat << EOF > ./pr_comment 
I am a robot that works on MagicModules PRs!

I built this PR into one or more PRs on other repositories, and when those are closed, this PR will also be merged and closed.
depends: https://github.com/$GH_USERNAME/terraform-provider-google/pull/$TF_PR
EOF
