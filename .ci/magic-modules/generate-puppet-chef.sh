#!/bin/bash

# This script takes in 'magic-modules-branched', a git repo tracking the head of a PR against magic-modules.
# It outputs "$PROVIDER-generated", a folder containing at least one non-submodule git repo containing the
# generated puppet/chef code.

set -x
set -e
source "$(dirname "$0")/helpers.sh"
PATCH_DIR="$(pwd)/patches"

IFS="," read -ra PRODUCT_ARRAY <<< "$PRODUCTS"
for PRD in "${PRODUCT_ARRAY[@]}"; do
  pushd magic-modules-branched
    LAST_COMMIT_AUTHOR="$(git log --pretty="%an <%ae>" -n1 HEAD)"
    find build/"${PROVIDER}/${PRD}"/ -type f -not -name '.git*' -not -name '.last_run.json' -print0 | xargs -0 rm -rf --
    bundle install
    # Running with the --debug flag will cause Concourse to crash
    bundle exec compiler -p "products/$PRD" -e "$PROVIDER" -o "build/$PROVIDER/$PRD"

    # This command can crash - if that happens, the script should not fail.
    set +e
    COMMIT_MSG="$(python .ci/magic-modules/extract_from_pr_description.py --tag "$PROVIDER-$PRD" < .git/body)"
    set -e
    if [ -z "$COMMIT_MSG" ]; then
      COMMIT_MSG="Magic Modules changes."
    fi

    pushd "build/$PROVIDER/$PRD"
      # These config entries will set the "committer".
      git config --global user.email "magic-modules@google.com"
      git config --global user.name "Modular Magician"

      git add -A
      # Set the "author" to the commit's real author.
      git commit -m "$COMMIT_MSG" --author="$LAST_COMMIT_AUTHOR" || true  # don't crash if no changes
      git checkout -B "$(cat ../../../branchname)"
      apply_patches "$PATCH_DIR/GoogleCloudPlatform/$PROVIDER-google-$PRD" "$COMMIT_MSG" "$LAST_COMMIT_AUTHOR" "master"
    popd
  popd
  git clone "magic-modules-branched/build/$PROVIDER/$PRD" "$PROVIDER-generated/$PRD"

done
