#!/usr/bin/env bash

set -e

# Setup GOPATH
export GOPATH=${PWD}/go

set -x

# Create GOPATH structure
mkdir -p "${GOPATH}/src/github.com/terraform-providers"
ln -s "${PWD}/terraform-diff/${SUBDIR}/new" "${GOPATH}/src/github.com/terraform-providers/$PROVIDER_NAME"

cd "${GOPATH}/src/github.com/terraform-providers/$PROVIDER_NAME"

go test -v ./$TEST_DIR -parallel 16 -run '^Test' -timeout 1m

