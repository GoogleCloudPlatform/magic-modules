#!/usr/bin/env bash

set -e

# Setup GOPATH
export GOPATH=${PWD}/go

set -x

# Create GOPATH structure
mkdir -p "${GOPATH}/src/github.com/terraform-providers"
ln -s "${PWD}/terraform-provider-google" "${GOPATH}/src/github.com/terraform-providers/terraform-provider-google"

cd "${GOPATH}/src/github.com/terraform-providers/terraform-provider-google"

make build vet