#!/usr/bin/env bash

set -e

# Setup GOPATH
export GOPATH=${PWD}/go

set -x

# Create GOPATH structure
mkdir -p "${GOPATH}/src/github.com/terraform-providers"
ln -s "${PWD}/$1" "${GOPATH}/src/github.com/terraform-providers/terraform-provider-google"

cd "${GOPATH}/src/github.com/terraform-providers/terraform-provider-google"

go test -v ./google -parallel 16 -run '^Test' -timeout 1m
