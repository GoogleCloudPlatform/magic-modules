#!/bin/bash

# This script takes in 'magic-modules-branched', a git repo tracking the head of a PR against magic-modules.
# It needs to output the same git repo, but with the code generation done and submodules updated, at 'magic-modules-submodules'.

set -e
set +x
# Don't show the credential in the output.
echo "$CREDS" > ~/github_private_key
set -x
chmod 400 ~/github_private_key

pushd magic-modules-branched
BRANCH="$(cat ./branchname)"
# Update this repo to track the submodules we just pushed:
IFS="," read -ra PRODUCT_ARRAY <<< "$PUPPET_MODULES"
for PRD in "${PRODUCT_ARRAY[@]}"; do
  git config -f .gitmodules "submodule.build/puppet/$PRD.branch" "$BRANCH"
  # Bundle repo does not use the same naming scheme as the others
  if [[ $PRD == *"_bundle"* ]]; then
    repo="puppet-google"
  else
    repo="puppet-google-$PRD"
  fi
  git config -f .gitmodules "submodule.build/puppet/$PRD.url" "git@github.com:$GH_USERNAME/$repo.git"
  git submodule sync "build/puppet/$PRD"
  ssh-agent bash -c "ssh-add ~/github_private_key; git submodule update --remote --init build/puppet/$PRD"
  git add "build/puppet/$PRD"
done

IFS="," read -ra PRODUCT_ARRAY <<< "$CHEF_MODULES"
for PRD in "${PRODUCT_ARRAY[@]}"; do
  git config -f .gitmodules "submodule.build/chef/$PRD.branch" "$BRANCH"
  # Bundle repo does not use the same naming scheme as the others
  if [[ $PRD == *"_bundle"* ]]; then
    repo="chef-google"
  else
    repo="chef-google-$PRD"
  fi
  git config -f .gitmodules "submodule.build/chef/$PRD.url" "git@github.com:$GH_USERNAME/$repo.git"
  git submodule sync "build/chef/$PRD"
  ssh-agent bash -c "ssh-add ~/github_private_key; git submodule update --remote --init build/chef/$PRD"
  git add "build/chef/$PRD"
done

if [ "$TERRAFORM_ENABLED" = "true" ]; then
  IFS="," read -ra TERRAFORM_VERSIONS <<< "$TERRAFORM_VERSIONS"
  for VERSION in "${TERRAFORM_VERSIONS[@]}"; do
    if [ -n "$VERSION" ]; then
      PROVIDER_NAME="terraform-provider-google-$VERSION"
      SUBMODULE_DIR="terraform-$VERSION"
    else
      PROVIDER_NAME="terraform-provider-google"
      SUBMODULE_DIR="terraform"
    fi
    git config -f .gitmodules "submodule.build/$SUBMODULE_DIR.branch" "$BRANCH"
    git config -f .gitmodules "submodule.build/$SUBMODULE_DIR.url" "git@github.com:$GH_USERNAME/$PROVIDER_NAME.git"
    git submodule sync "build/$SUBMODULE_DIR"
    ssh-agent bash -c "ssh-add ~/github_private_key; git submodule update --remote --init build/$SUBMODULE_DIR"
    git add "build/$SUBMODULE_DIR"
  done
fi

if [ "$ANSIBLE_ENABLED" = "true" ]; then
  git config -f .gitmodules submodule.build/ansible.branch "$BRANCH"
  git config -f .gitmodules submodule.build/ansible.url "git@github.com:$GH_USERNAME/ansible.git"
  git submodule sync build/ansible
  ssh-agent bash -c "ssh-add ~/github_private_key; git submodule update --remote --init build/ansible"
  git add build/ansible
fi

# Commit those changes so that they can be tested in the next phase.
git add .gitmodules
git config --global user.email "magic-modules@google.com"
git config --global user.name "Modular Magician"
git commit -m "Automatic submodule update to generated code." || true  # don't crash if no changes
git checkout -B "$BRANCH"

cp -r ./ ../magic-modules-submodules
