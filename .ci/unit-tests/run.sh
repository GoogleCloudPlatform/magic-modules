#!/usr/bin/env bash

set -e

# Setup GOPATH
export GOPATH=${PWD}/go

set -x


if [ -n "$VERSION" ]; then
  PROVIDER_NAME="terraform-provider-google-$VERSION"
  SUBMODULE_DIR="terraform-$VERSION"
else
  PROVIDER_NAME="terraform-provider-google"
  SUBMODULE_DIR="terraform"
fi

# Create GOPATH structure
mkdir -p "${GOPATH}/src/github.com/terraform-providers"
ln -s "${PWD}/magic-modules/build/$SUBMODULE_DIR" "${GOPATH}/src/github.com/terraform-providers/$PROVIDER_NAME"

cd "${GOPATH}/src/github.com/terraform-providers/$PROVIDER_NAME"

go test -v ./google -parallel 16 -run '^Test' -timeout 1m
