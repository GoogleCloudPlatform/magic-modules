#!/usr/bin/env bash

set -x

# Install dependencies for Template Generator
pushd "magic-modules-new-prs"
bundle install

# Setup SSH keys.

# Since these creds are going to be managed externally, we need to pass
# them into the container as an environment variable.  We'll use
# ssh-agent to ensure that these are the credentials used to update.
set +x
echo "$CREDS" > ~/github_private_key
set -x
chmod 400 ~/github_private_key

ssh-agent bash -c "ssh-add ~/github_private_key; git submodule update --remote --init build/ansible"
popd

pushd "magic-modules-new-prs/build/ansible"
# Setup Git config.
git config --global user.email "magic-modules@google.com"
git config --global user.name "Modular Magician"

# Run creation script.
ssh-agent bash -c "ssh-add ~/github_private_key; ../../tools/ansible-pr/run.sh"
