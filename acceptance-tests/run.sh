#!/usr/bin/env bash

set -e

export GOOGLE_CREDENTIALS_FILE="/tmp/google-account.json"
export TF_ACC=1
export GOOGLE_REGION="us-central1"
# Setup GOPATH
export GOPATH=${PWD}/go

# CI sets the contents of our json account secret in our environment; dump it
# to disk for use in tests.
echo "${GOOGLE_JSON_ACCOUNT}" > /tmp/google-account.json

set -x

# Create GOPATH structure
mkdir -p "${GOPATH}/src/github.com/terraform-providers"
ln -s "${PWD}/terraform-provider-google" "${GOPATH}/src/github.com/terraform-providers/terraform-provider-google"

cd "${GOPATH}/src/github.com/terraform-providers/terraform-provider-google"

# TODO: Run all acceptance tests. We need to run our acceptance tests in an account outside the google org for the following reasons:
# - Enforcers automatically creates firewall rules which caused all tests creating network to fail to destroy properly
# - We can't run tests creating projects, folders or managing org policies in the google org.
go test -v ./google -parallel 16 -run '^TestAccComputeAddress_' -timeout 120m
