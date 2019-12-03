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

if [ "$TERRAFORM_ENABLED" = "true" ]; then
  IFS="," read -ra TERRAFORM_VERSIONS <<< "$TERRAFORM_VERSIONS"
  for VERSION in "${TERRAFORM_VERSIONS[@]}"; do
    IFS=":" read -ra TERRAFORM_DATA <<< "$VERSION"
    PROVIDER_NAME="${TERRAFORM_DATA[0]}"
    SUBMODULE_DIR="${TERRAFORM_DATA[1]}"

    git config -f .gitmodules "submodule.build/$SUBMODULE_DIR.branch" "$BRANCH"
    git config -f .gitmodules "submodule.build/$SUBMODULE_DIR.url" "https://github.com/$GH_USERNAME/$PROVIDER_NAME.git"
    git submodule sync "build/$SUBMODULE_DIR"
    ssh-agent bash -c "ssh-add ~/github_private_key; git submodule update --remote --init build/$SUBMODULE_DIR"
    git add "build/$SUBMODULE_DIR"
  done
fi

if [ "$ANSIBLE_ENABLED" = "true" ]; then
  git config -f .gitmodules submodule.build/ansible.branch "$BRANCH"
  git config -f .gitmodules submodule.build/ansible.url "https://github.com/$GH_USERNAME/ansible_collections_google.git"
  git submodule sync build/ansible
  ssh-agent bash -c "ssh-add ~/github_private_key; git submodule update --remote --init build/ansible"
  git add build/ansible
fi

if [ "$INSPEC_ENABLED" = "true" ]; then
  git config -f .gitmodules submodule.build/inspec.branch "$BRANCH"
  git config -f .gitmodules submodule.build/inspec.url "https://github.com/$GH_USERNAME/inspec-gcp.git"
  git submodule sync build/inspec
  ssh-agent bash -c "ssh-add ~/github_private_key; git submodule update --remote --init build/inspec"
  git add build/inspec
fi

# Commit those changes so that they can be tested in the next phase.
git add .gitmodules
git config --global user.email "magic-modules@google.com"
git config --global user.name "Modular Magician"
git commit -m "Automatic submodule update to generated code." || true  # don't crash if no changes
git checkout -B "$BRANCH"

cp -r ./ ../magic-modules-submodules
