#!/bin/bash

# This script takes in 'magic-modules-branched', a git repo tracking the head of a PR against magic-modules.
# It outputs "ansible-generated", a non-submodule git repo containing the generated ansible code.

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

# Remove all modules so that old files are removed in process.
rm build/ansible/plugins/modules/gcp_*

bundle exec compiler -a -e ansible -o "build/ansible/"

ANSIBLE_COMMIT_MSG="$(cat .git/title)"

pushd "build/ansible"
# This module is handwritten. It's the only one.
# It was deleted earlier, so it needs to be undeleted.
git checkout HEAD -- plugins/modules/gcp_storage_object.py

# These config entries will set the "committer".
git config --global user.email "magic-modules@google.com"
git config --global user.name "Modular Magician"

git add -A

git commit -m "$ANSIBLE_COMMIT_MSG" --author="$COMMIT_AUTHOR" || true  # don't crash if no changes
git checkout -B "$(cat ../../branchname)"

apply_patches "$PATCH_DIR/modular-magician/ansible" "$ANSIBLE_COMMIT_MSG" "$COMMIT_AUTHOR" "master"

popd
popd

git clone magic-modules-branched/build/ansible ./ansible-generated
