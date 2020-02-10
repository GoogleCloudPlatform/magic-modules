#!/bin/bash

set -e

REPO=$1
REFERENCE=$2
SCRATCH_OWNER=modular-magician
if [ "$REPO" == "ga" ]; then
    GH_REPO=terraform-provider-google
    LOCAL_PATH=$GOPATH/src/github.com/terraform-providers/terraform-provider-google
elif [ "$REPO" == "beta" ]; then
    GH_REPO=terraform-provider-google-beta
    LOCAL_PATH=$GOPATH/src/github.com/terraform-providers/terraform-provider-google-beta
else
    echo "no repo, dying."
    exit 1
fi

SCRATCH_PATH=https://modular-magician:$GITHUB_TOKEN@github.com/$SCRATCH_OWNER/$GH_REPO
mkdir -p "$(dirname $LOCAL_PATH)"
git clone $SCRATCH_PATH $LOCAL_PATH --single-branch --branch "auto-pr-$REFERENCE" --depth 1
pushd $LOCAL_PATH

make tools
make docscheck
make lint
make
make test
