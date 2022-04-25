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
TFV_SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/modular-magician/terraform-validator
TFV_LOCAL_PATH=$PWD/../tfv
TFOICS_SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/modular-magician/docs-examples
TFOICS_LOCAL_PATH=$PWD/../tfoics
TFCD_SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/modular-magician/terraform-docs-samples
TFCD_LOCAL_PATH=$PWD/../tfcd

# For backwards compatibility until at least Nov 15 2021
TFC_SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/modular-magician/terraform-google-conversion
TFC_LOCAL_PATH=$PWD/../tfc

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

# TF Conversion - for compatibility until at least Nov 15 2021
mkdir -p $TFC_LOCAL_PATH
# allow this to fail for compatibility during tfv/tgc transition phase
if git clone -b $NEW_BRANCH $TFC_SCRATCH_PATH $TFC_LOCAL_PATH; then
    pushd $TFC_LOCAL_PATH
    git fetch origin $OLD_BRANCH
    if ! git diff --exit-code origin/$OLD_BRANCH origin/$NEW_BRANCH; then
        SUMMARY=`git diff origin/$OLD_BRANCH origin/$NEW_BRANCH --shortstat`
        DIFFS="${DIFFS}${NEWLINE}TF Conversion: [Diff](https://github.com/modular-magician/terraform-google-conversion/compare/$OLD_BRANCH..$NEW_BRANCH) ($SUMMARY)"
    fi
    popd
fi

# TF Validator
mkdir -p $TFV_LOCAL_PATH
# allow this to fail for compatibility during tfv/tgc transition phase
if git clone -b $NEW_BRANCH $TFV_SCRATCH_PATH $TFV_LOCAL_PATH; then
    pushd $TFV_LOCAL_PATH
    git fetch origin $OLD_BRANCH
    if ! git diff --exit-code origin/$OLD_BRANCH origin/$NEW_BRANCH; then
        SUMMARY=`git diff origin/$OLD_BRANCH origin/$NEW_BRANCH --shortstat`
        DIFFS="${DIFFS}${NEWLINE}TF Validator: [Diff](https://github.com/modular-magician/terraform-validator/compare/$OLD_BRANCH..$NEW_BRANCH) ($SUMMARY)"
    fi
    popd
fi

# TF OICS
mkdir -p $TFOICS_LOCAL_PATH
git clone -b $NEW_BRANCH $TFOICS_SCRATCH_PATH $TFOICS_LOCAL_PATH
pushd $TFOICS_LOCAL_PATH
git fetch origin $OLD_BRANCH
if ! git diff --exit-code --quiet origin/$OLD_BRANCH origin/$NEW_BRANCH; then
    SUMMARY="$(git diff origin/$OLD_BRANCH origin/$NEW_BRANCH --shortstat)"
    DIFFS="${DIFFS}${NEWLINE}TF OiCS: [Diff](https://github.com/modular-magician/docs-examples/compare/$OLD_BRANCH..$NEW_BRANCH) ($SUMMARY)"
fi
popd

# TF Cloud Docs
mkdir -p $TFCD_LOCAL_PATH
git clone -b $NEW_BRANCH $TFCD_SCRATCH_PATH $TFCD_LOCAL_PATH
pushd $TFCD_LOCAL_PATH
git fetch origin $OLD_BRANCH
if ! git diff --exit-code --quiet origin/$OLD_BRANCH origin/$NEW_BRANCH; then
    SUMMARY="$(git diff origin/$OLD_BRANCH origin/$NEW_BRANCH --shortstat)"
    DIFFS="${DIFFS}${NEWLINE}TF Cloud Doc Samples: [Diff](https://github.com/modular-magician/terraform-docs-samples/compare/$OLD_BRANCH..$NEW_BRANCH) ($SUMMARY)"
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
