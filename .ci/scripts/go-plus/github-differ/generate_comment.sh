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

# Get the directory of the current script
script_dir=$(dirname "$0")

# make utility functions avalible
source "$script_dir/utils.sh"

PR_NUMBER=$1
NEW_BRANCH=auto-pr-$PR_NUMBER
OLD_BRANCH=auto-pr-$PR_NUMBER-old
MM_LOCAL_PATH=$PWD
TPG_SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/modular-magician/terraform-provider-google
TPG_LOCAL_PATH=$PWD/../tpg
TPGB_SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/modular-magician/terraform-provider-google-beta
TPGB_LOCAL_PATH=$PWD/../tpgb
TFOICS_SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/modular-magician/docs-examples
TFOICS_LOCAL_PATH=$PWD/../tfoics
BREAKING_CHANGE_BUILD_FAILURE=0

# For backwards compatibility until at least Nov 15 2021
TFC_SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/modular-magician/terraform-google-conversion
TFC_LOCAL_PATH=$PWD/../tfc

DIFFS=""
NEWLINE=$'\n'

# TPG/TPGB difference
mkdir -p $TPG_LOCAL_PATH
git clone -b $NEW_BRANCH $TPG_SCRATCH_PATH $TPG_LOCAL_PATH
pushd $TPG_LOCAL_PATH
git fetch origin $OLD_BRANCH
if ! git diff --exit-code origin/$OLD_BRANCH origin/$NEW_BRANCH; then
    SUMMARY=`git diff origin/$OLD_BRANCH origin/$NEW_BRANCH --shortstat`
    DIFFS="${DIFFS}${NEWLINE}Terraform GA: [Diff](https://github.com/modular-magician/terraform-provider-google/compare/$OLD_BRANCH..$NEW_BRANCH) ($SUMMARY)"
fi
popd

mkdir -p $TPGB_LOCAL_PATH
git clone -b $NEW_BRANCH $TPGB_SCRATCH_PATH $TPGB_LOCAL_PATH
pushd $TPGB_LOCAL_PATH
git fetch origin $OLD_BRANCH
if ! git diff --exit-code origin/$OLD_BRANCH origin/$NEW_BRANCH; then
    SUMMARY=`git diff origin/$OLD_BRANCH origin/$NEW_BRANCH --shortstat`
    DIFFS="${DIFFS}${NEWLINE}Terraform Beta: [Diff](https://github.com/modular-magician/terraform-provider-google-beta/compare/$OLD_BRANCH..$NEW_BRANCH) ($SUMMARY)"
fi
popd

## Breaking change setup and execution
set +e
pushd $MM_LOCAL_PATH/tools/diff-processor
cp -r $TPG_LOCAL_PATH old/
cp -r $TPG_LOCAL_PATH new/
make build OLD_REF=$OLD_BRANCH NEW_REF=$NEW_BRANCH
TPG_BREAKING="$(bin/diff-processor breaking-changes)"
retVal=$?
if [ $retVal -ne 0 ]; then
    TPG_BREAKING=""
    BREAKING_CHANGE_BUILD_FAILURE=1
fi

rm -rf ./old/ ./new/ ./bin/
cp -r $TPGB_LOCAL_PATH old/
cp -r $TPGB_LOCAL_PATH new/
make build OLD_REF=$OLD_BRANCH NEW_REF=$NEW_BRANCH
TPGB_BREAKING="$(bin/diff-processor breaking-changes)"
retVal=$?
if [ $retVal -ne 0 ]; then
    TPGB_BREAKING=""
    BREAKING_CHANGE_BUILD_FAILURE=1
fi
if [ $BREAKING_CHANGE_BUILD_FAILURE -eq 0 ]; then
    echo "Breaking changes succeeded"
    # Export variables here so that they can be used in compare_breaking_changes
    # exporting earlier would cause the retvals to not be calculated properly.
    export TPG_BREAKING=$TPG_BREAKING
    export TPGB_BREAKING=$TPGB_BREAKING
    BREAKINGCHANGES="$($script_dir/compare_breaking_changes.sh)"
else
    echo "Breaking changes failed"
    BREAKINGCHANGES="## Breaking Change Detection Failed${NEWLINE}The breaking change detector crashed during execution. This is usually due to the downstream provider(s) failing to compile. Please investigate or follow up with your reviewer."
fi
popd
set -e

## Missing test setup and execution
TPGB_LOCAL_PATH_OLD="${TPGB_LOCAL_PATH}old"
mkdir -p $TPGB_LOCAL_PATH_OLD
cp -r $TPGB_LOCAL_PATH/. $TPGB_LOCAL_PATH_OLD

pushd $TPGB_LOCAL_PATH
git checkout origin/$NEW_BRANCH
update_package_name "github.com/hashicorp/terraform-provider-google-beta" "google/provider/new"
popd

pushd $TPGB_LOCAL_PATH_OLD
git checkout origin/$OLD_BRANCH
update_package_name "github.com/hashicorp/terraform-provider-google-beta" "google/provider/old"
popd

set +e
pushd $MM_LOCAL_PATH/tools/missing-test-detector
go mod edit -replace google/provider/new=$(realpath $TPGB_LOCAL_PATH)
go mod edit -replace google/provider/old=$(realpath $TPGB_LOCAL_PATH_OLD)
go mod tidy
MISSINGTESTS="$(go run . -services-dir=$TPGB_LOCAL_PATH/google-beta/services)"
retVal=$?
if [ $retVal -ne 0 ]; then
    MISSINGTESTS=""
fi
set -e
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

MESSAGE="Hi there, I'm the Modular magician. I've detected the following information about your changes:${NEWLINE}${NEWLINE}"

BREAKINGSTATE="success"
if [ -n "$BREAKINGCHANGES" ]; then
  MESSAGE="${MESSAGE}${BREAKINGCHANGES}${NEWLINE}${NEWLINE}"

  BREAKINGCHANGE_OVERRIDE=$(curl \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer $GITHUB_TOKEN" \
    "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/pulls/$PR_NUMBER"\
    | jq ".labels|any(.id==4598495472)")

  if [ "${BREAKINGCHANGE_OVERRIDE}" == "true" ]; then
    BREAKINGSTATE="success"
  else
    BREAKINGSTATE="failure"
  fi
fi


if [ -z "$DIFFS" ]; then
  MESSAGE="${MESSAGE}## Diff report ${NEWLINE}Your PR hasn't generated any diffs, but I'll let you know if a future commit does."
else
  MESSAGE="${MESSAGE}## Diff report ${NEWLINE}Your PR generated some diffs in downstreams - here they are.${NEWLINE}${DIFFS}"
  if [ -n "$MISSINGTESTS" ]; then
    MESSAGE="${MESSAGE}${NEWLINE}${MISSINGTESTS}${NEWLINE}"
  fi
fi



#;region=global/${BUILD_ID};step=19?project=${PROJECT_ID}"
BREAKINGSTATE_BODY=$( jq -n \
	--arg context "terraform-provider-breaking-change-test" \
	--arg target_url "https://console.cloud.google.com/cloud-build/builds;region=global/${BUILD_ID};step=${BUILD_STEP}?project=${PROJECT_ID}" \
	--arg breakingstate "${BREAKINGSTATE}" \
	'{context: $context, target_url: $target_url, state: $breakingstate}')

curl \
  -X POST \
  -u "$github_username:$GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/$COMMIT_SHA" \
  -d "$BREAKINGSTATE_BODY"

curl -H "Authorization: token ${GITHUB_TOKEN}" \
      -d "$(jq -r --arg diffs "$MESSAGE" -n "{body: \$diffs}")" \
      "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/${PR_NUMBER}/comments"

if ! git diff --exit-code origin/main tools; then
    ## Run unit tests for breaking change and missing test detector.
    "$script_dir/test_tools.sh" "$MM_LOCAL_PATH" "$TPGB_LOCAL_PATH" "$COMMIT_SHA" "$BUILD_ID" "$BUILD_STEP" "$PROJECT_ID"
fi
