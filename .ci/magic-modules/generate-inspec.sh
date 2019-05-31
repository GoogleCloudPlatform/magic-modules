#!/bin/bash

# This script takes in 'magic-modules-branched', a git repo tracking the head of a PR against magic-modules.
# It outputs "inspec-generated", a non-submodule git repo containing the generated inspec code.

set -x
set -e
source "$(dirname "$0")/helpers.sh"
PATCH_DIR="$(pwd)/patches"
pushd magic-modules-branched
LAST_COMMIT_AUTHOR="$(git log --pretty="%an <%ae>" -n1 HEAD)"
for i in $(find products/ -name 'inspec.yaml' -printf '%h\n');
do
  bundle exec compiler -p $i -e inspec -o "build/inspec/"
done

# This command can crash - if that happens, the script should not fail.
set +e
INSPEC_COMMIT_MSG="$(python .ci/magic-modules/extract_from_pr_description.py --tag inspec < .git/body)"
set -e
if [ -z "$INSPEC_COMMIT_MSG" ]; then
  INSPEC_COMMIT_MSG="$(cat .git/title)"
fi

pushd "build/inspec"
# These config entries will set the "committer".
git config --global user.email "magic-modules@google.com"
git config --global user.name "Modular Magician"

git add -A
# Set the "author" to the commit's real author.
git commit -m "$INSPEC_COMMIT_MSG" --author="$LAST_COMMIT_AUTHOR" || true  # don't crash if no changes
git checkout -B "$(cat ../../branchname)"

apply_patches "$PATCH_DIR/modular-magician/inspec-gcp" "$INSPEC_COMMIT_MSG" "$LAST_COMMIT_AUTHOR" "master"
popd
popd

git clone magic-modules-branched/build/inspec ./inspec-generated
