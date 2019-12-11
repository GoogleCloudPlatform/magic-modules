#!/bin/bash

# This script takes in 'magic-modules-branched', a git repo tracking the head of a PR against magic-modules.
# It outputs "inspec-generated", a non-submodule git repo containing the generated inspec code.

set -x
set -e
source "$(dirname "$0")/helpers.sh"
PATCH_DIR="$(pwd)/patches"

pushd magic-modules-branched

# Choose the author of the most recent commit as the downstream author
# Note that we don't use the last submitted commit, we use the primary GH email
# of the GH PR submitted. If they've enabled a private email, we'll actually
# use their GH noreply email which isn't compatible with CLAs.
COMMIT_AUTHOR="$(git log --pretty="%an <%ae>" -n1 HEAD)"

bundle exec compiler -a -e inspec -o "build/inspec/"

INSPEC_COMMIT_MSG="$(cat .git/title)"

pushd "build/inspec"

# These config entries will set the "committer".
git config --global user.email "magic-modules@google.com"
git config --global user.name "Modular Magician"

git add -A

git commit -m "$INSPEC_COMMIT_MSG" --author="$COMMIT_AUTHOR" || true  # don't crash if no changes
git checkout -B "$(cat ../../branchname)"

apply_patches "$PATCH_DIR/modular-magician/inspec-gcp" "$INSPEC_COMMIT_MSG" "$COMMIT_AUTHOR" "master"

popd
popd

git clone magic-modules-branched/build/inspec ./inspec-generated
