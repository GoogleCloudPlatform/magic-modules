#!/bin/bash

# This script configures the git submodule under magic-modules so that it is
# ready to create a new pull request.  It is cloned in a detached-head state,
# but its branch is relevant to the PR creation process, so we want to make
# sure that it's on a branch, and most importantly that that branch tracks
# a branch upstream.

set -e
set -x

shopt -s dotglob
cp -r magic-modules/* magic-modules-out

cd magic-modules-out

# This says "check out the branch which contains HEAD, and set it up to track its upstream."
git checkout -t "$(git branch -a --contains "$(git rev-parse HEAD)" | grep -v "detached")"

cd build/terraform

git checkout -t "$(git branch -a --contains "$(git rev-parse HEAD)" | grep -v "detached")"
# This special string 'new' tells the PR resource to create a new PR.
git config pullrequest.id new
