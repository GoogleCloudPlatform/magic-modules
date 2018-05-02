#!/bin/bash

# This script updates the submodule to track terraform master.
set -e
set -x
shopt -s dotglob

# Since these creds are going to be managed externally, we need to pass
# them into the container as an environment variable.  We'll use
# ssh-agent to ensure that these are the credentials used to update.
echo "$CREDS" > ~/github_private_key
chmod 400 ~/github_private_key

pushd mm-approved-prs
ID=$(git config --get pullrequest.id)
# We need to know what branch to check out for the update.
BRANCH=$(git config --get pullrequest.branch)
popd

cp -r mm-approved-prs/* mm-output

pushd mm-output
git config pullrequest.id "$ID"
git branch -f "$BRANCH"
git checkout "$BRANCH"
git config --global user.email "magic-modules@google.com"
git config --global user.name "Modular Magician"
ssh-agent bash -c "ssh-add ~/github_private_key; git submodule update --remote --init $ALL_SUBMODULES"

# Word-splitting here is intentional.
git add $ALL_SUBMODULES 

# It's okay for the commit to fail if there's no changes.
set +e
git commit -m "Update tracked submodules -> HEAD on $(date)

Tracked submodules are $ALL_SUBMODULES."
echo "Merged PR #$ID." > ./commit_message
