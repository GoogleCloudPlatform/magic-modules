#!/bin/bash

# This script configures the git submodule under magic-modules so that it is
# ready to create a new pull request.  It is cloned in a detached-head state,
# but its branch is relevant to the PR creation process, so we want to make
# sure that it's on a branch, and most importantly that that branch tracks
# a branch upstream.

set -e
set -x

shopt -s dotglob
cp -r magic-modules/* magic-modules-with-comment

pushd magic-modules-with-comment

# This says "check out the branch which contains HEAD, and set it up to track its upstream."
# The 'checkout' command tells us to checkout the branch, the '-t' says "which tracks".
# This 'git branch -a --contains' command lists all branches which contain the SHA provided
# by 'git rev-parse HEAD' (which returns the rev at HEAD).  Then we remove the 'detached'
# branch which we're currently on, collapse whitespace, and check it out.
# This will resolve to something like 'git checkout -t remotes/origin/abcdef12'
git checkout -t "$(git branch -a --contains "$(git rev-parse HEAD)" | grep -v "detached" | xargs echo -n)"

pushd build/terraform

cat << EOF > ./downstream_body
$(git log -1 --pretty=%B)

<!-- This change is generated by MagicModules. -->
EOF
git checkout -t "$(git branch -a --contains "$(git rev-parse HEAD)" | grep -v "detached" | xargs echo -n)"
TF_PR=$(hub pull-request -b "$TERRAFORM_REPO:master" -F ./downstream_body)
popd

cat << EOF > ./pr_comment 
I am a robot that works on MagicModules PRs!

I built this PR into one or more PRs on other repositories, and when those are closed, this PR will also be merged and closed.
depends: $TF_PR
EOF
