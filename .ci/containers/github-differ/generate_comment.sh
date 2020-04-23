#! /bin/bash

set -e

if [ $# -lt 1 ]; then
    echo "Usage: $0 pr-number"
    exit 1
fi
if [ -z "$GITHUB_TOKEN" ]; then
    echo "Did not provide GITHUB_TOKEN environment variable."
    exit 1
fi

PR_NUMBER=$1
NEW_BRANCH=auto-pr-$PR_NUMBER
OLD_BRANCH=auto-pr-$PR_NUMBER-old
TPG_SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/modular-magician/terraform-provider-google
TPG_LOCAL_PATH=$PWD/../tpg
TPGB_SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/modular-magician/terraform-provider-google-beta
TPGB_LOCAL_PATH=$PWD/../tpgb
TFC_SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/modular-magician/terraform-google-conversion
TFC_LOCAL_PATH=$PWD/../tfc
TFOICS_SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/modular-magician/docs-examples
TFOICS_LOCAL_PATH=$PWD/../tfoics
ANSIBLE_SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/modular-magician/google
ANSIBLE_LOCAL_PATH=$PWD/../ansible
INSPEC_SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/modular-magician/inspec-gcp
INSPEC_LOCAL_PATH=$PWD/../inspec

DIFFS=""
NEWLINE=$'\n'

# TPG
mkdir -p $TPG_LOCAL_PATH
git clone -b $NEW_BRANCH $TPG_SCRATCH_PATH $TPG_LOCAL_PATH
pushd $TPG_LOCAL_PATH
git fetch origin $OLD_BRANCH
if ! git diff --exit-code origin/$OLD_BRANCH origin/$NEW_BRANCH; then
    SUMMARY=`git diff origin/$OLD_BRANCH origin/$NEW_BRANCH --shortstat`
    DIFFS="${DIFFS}${NEWLINE}Terraform GA: [Diff](https://github.com/modular-magician/terraform-provider-google/compare/$OLD_BRANCH..$NEW_BRANCH) ($SUMMARY)"
fi
popd

# TPGB
mkdir -p $TPGB_LOCAL_PATH
git clone -b $NEW_BRANCH $TPGB_SCRATCH_PATH $TPGB_LOCAL_PATH
pushd $TPGB_LOCAL_PATH
git fetch origin $OLD_BRANCH
if ! git diff --exit-code origin/$OLD_BRANCH origin/$NEW_BRANCH; then
    SUMMARY=`git diff origin/$OLD_BRANCH origin/$NEW_BRANCH --shortstat`
    DIFFS="${DIFFS}${NEWLINE}Terraform Beta: [Diff](https://github.com/modular-magician/terraform-provider-google-beta/compare/$OLD_BRANCH..$NEW_BRANCH) ($SUMMARY)"
fi
popd

# Ansible
mkdir -p $ANSIBLE_LOCAL_PATH
git clone -b $NEW_BRANCH $ANSIBLE_SCRATCH_PATH $ANSIBLE_LOCAL_PATH
pushd $ANSIBLE_LOCAL_PATH
git fetch origin $OLD_BRANCH
if ! git diff --exit-code origin/$OLD_BRANCH origin/$NEW_BRANCH; then
    SUMMARY=`git diff origin/$OLD_BRANCH origin/$NEW_BRANCH --shortstat`
    DIFFS="${DIFFS}${NEWLINE}Ansible: [Diff](https://github.com/modular-magician/google/compare/$OLD_BRANCH..$NEW_BRANCH) ($SUMMARY)"
fi
popd

# TF Conversion
mkdir -p $TFC_LOCAL_PATH
git clone -b $NEW_BRANCH $TFC_SCRATCH_PATH $TFC_LOCAL_PATH
pushd $TFC_LOCAL_PATH
git fetch origin $OLD_BRANCH
if ! git diff --exit-code origin/$OLD_BRANCH origin/$NEW_BRANCH; then
    SUMMARY=`git diff origin/$OLD_BRANCH origin/$NEW_BRANCH --shortstat`
    DIFFS="${DIFFS}${NEWLINE}TF Conversion: [Diff](https://github.com/modular-magician/terraform-google-conversion/compare/$OLD_BRANCH..$NEW_BRANCH) ($SUMMARY)"
fi
popd

# Temporary guards for OiCS since this container image will be newer than YAML
# for existing PRs. We need to be backwards-compatible with old GCB runs that
# don't generate OiCS. This should be safe to remove after Feb 18th or so.
set +e
# Use a subshell to run commands serially and fail when one fails
OICSDIFFS=$(bash -e <<TRY
    # TF OICS
    mkdir -p $TFOICS_LOCAL_PATH
    git clone -b $NEW_BRANCH $TFOICS_SCRATCH_PATH $TFOICS_LOCAL_PATH
    pushd $TFOICS_LOCAL_PATH > /dev/null
    git fetch origin $OLD_BRANCH
    if ! git diff --exit-code --quiet origin/$OLD_BRANCH origin/$NEW_BRANCH; then
        SUMMARY=`git diff origin/$OLD_BRANCH origin/$NEW_BRANCH --shortstat`
        echo "TF OiCS: [Diff](https://github.com/modular-magician/docs-examples/compare/$OLD_BRANCH..$NEW_BRANCH) ($SUMMARY)"
    fi
    popd > /dev/null
TRY
)
if [ $? -ne 0 ]; then
  echo failed to generate OiCS
else
  DIFFS="${DIFFS}${NEWLINE}${OICSDIFFS}"
fi
set -e

# Inspec
mkdir -p $INSPEC_LOCAL_PATH
git clone -b $NEW_BRANCH $INSPEC_SCRATCH_PATH $INSPEC_LOCAL_PATH
pushd $INSPEC_LOCAL_PATH
git fetch origin $OLD_BRANCH
if ! git diff --exit-code origin/$OLD_BRANCH origin/$NEW_BRANCH; then
    SUMMARY=`git diff origin/$OLD_BRANCH origin/$NEW_BRANCH --shortstat`
    DIFFS="${DIFFS}${NEWLINE}Inspec: [Diff](https://github.com/modular-magician/inspec-gcp/compare/$OLD_BRANCH..$NEW_BRANCH) ($SUMMARY)"
fi
popd

if [ -z "$DIFFS" ]; then
    DIFFS="Hi!  I'm the modular magician.  Your PR hasn't generated any diffs, but I'll let you know if a future commit does."
else
    DIFFS="Hi!  I'm the modular magician.  Your PR generated some diffs in downstreams - here they are.$NEWLINE# Diff report:$NEWLINE$NEWLINE$DIFFS"
fi

curl -H "Authorization: token ${GITHUB_TOKEN}" \
      -d "$(jq -r --arg diffs "$DIFFS" -n "{body: \$diffs}")" \
      "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/${PR_NUMBER}/comments"
