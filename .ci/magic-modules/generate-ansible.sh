#!/bin/bash

# This script takes in 'magic-modules-branched', a git repo tracking the head of a PR against magic-modules.
# It outputs "ansible-generated", a non-submodule git repo containing the generated ansible code.

set -x
set -e
source "$(dirname "$0")/helpers.sh"
PATCH_DIR="$(pwd)/patches"

pushd magic-modules-branched

# Choose the author of the most recent commit as the downstream author
COMMIT_AUTHOR="$(git log --pretty="%an <%ae>" -n1 HEAD)"

for i in $(find products/ -name 'ansible.yaml' -printf '%h\n');
do
  bundle exec compiler -p $i -e ansible -o "build/ansible/"
done

ANSIBLE_COMMIT_MSG="$(cat .git/title)"

pushd "build/ansible"

# These config entries will set the "committer".
git config --global user.email "magic-modules@google.com"
git config --global user.name "Modular Magician"

git add -A

git commit -m "$ANSIBLE_COMMIT_MSG" --author="$COMMIT_AUTHOR" || true  # don't crash if no changes
git checkout -B "$(cat ../../branchname)"

apply_patches "$PATCH_DIR/modular-magician/ansible" "$ANSIBLE_COMMIT_MSG" "$COMMIT_AUTHOR" "devel"

popd
popd

git clone magic-modules-branched/build/ansible ./ansible-generated
