#!/bin/bash

set -e

PR_NUMBER=$1
SCRATCH_OWNER=modular-magician
GH_REPO=terraform-google-conversion

SCRATCH_PATH=https://$SCRATCH_OWNER:$GITHUB_TOKEN@github.com/$SCRATCH_OWNER/$GH_REPO
LOCAL_PATH=$GOPATH/src/github.com/GoogleCloudPlatform/$GH_REPO
mkdir -p "$(dirname $LOCAL_PATH)"
git clone $SCRATCH_PATH $LOCAL_PATH --single-branch --branch "auto-pr-$PR_NUMBER" --depth 1
pushd $LOCAL_PATH

make test
