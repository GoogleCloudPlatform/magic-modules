#! /bin/bash
set -e
set -x

pushd "magic-modules"
export GH_TOKEN
if PR_ID=$(git config --get pullrequest.id) &&
  [ -z "$USE_SHA" ] &&
  DEPS=$(python ./.ci/magic-modules/get_downstream_prs.py "$PR_ID") &&
  [ -z "$DEPS" ]; then
  BRANCH="codegen-pr-$(git config --get pullrequest.id)"
else
  BRANCH="codegen-sha-$(git rev-parse --short HEAD)"
fi
git checkout -B "$BRANCH"
# ./branchname is intentionally never committed - it isn't necessary once
# this output is no longer available.
echo "$BRANCH" > ./branchname

set +x
# Don't show the credential in the output.
echo "$CREDS" > ~/github_private_key
set -x
chmod 400 ~/github_private_key

# Update to head on master on all submodules, so we avoid spurious diffs.
# Note: $ALL_SUBMODULES will be re-split by the ssh-agent's "bash".
ssh-agent bash -c "ssh-add ~/github_private_key; git submodule update --remote --init $ALL_SUBMODULES"

cp -r ./ ../magic-modules-branched/

if [ "true" == "$INCLUDE_PREVIOUS" ] ; then
    # Since this is fetched after a merge commit, HEAD~ is
    # the newest commit on the branch being merged into.
    git reset --hard HEAD~
    BRANCH="$BRANCH-previous"
    git checkout -B "$BRANCH"
    # ./branchname is intentionally never committed - it isn't necessary once
    # this output is no longer available.
    echo "$BRANCH" > ./branchname
    ssh-agent bash -c "ssh-add ~/github_private_key; git submodule update --remote --init $ALL_SUBMODULES"
    cp -r ./ ../magic-modules-previous/
fi
