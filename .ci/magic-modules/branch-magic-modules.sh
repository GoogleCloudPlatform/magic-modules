#! /bin/bash
set -e
set -x

pushd "magic-modules"
# "codegen-pr" vs "codegen-sha" is *LOAD-BEARING*.  Don't change
# them (or introduce other options) unless you also change the
# logic in create-or-update-pr - because we decide whether to
# create or to update by which one of these we're prefixed by.
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
ssh-agent bash -c "ssh-add ~/github_private_key; git submodule update --remote --init build/terraform build/puppet/sql build/puppet/compute"

cp -r ./ ../magic-modules-branched/
