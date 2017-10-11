#!/usr/bin/env bash

set -e

export GOOGLE_CREDENTIALS_FILE="/tmp/google-account.json"
export GCLOUD_PROJECT="terraform-ci-acc-tests"
export TF_ACC=1
export GOOGLE_REGION="us-central1"
# TODO actually use a separate project for xpn resources
export GOOGLE_XPN_HOST_PROJECT="man-i-wish-i-was-a-real-project"
# Setup GOPATH
export GOPATH=${PWD}/go

# CI sets the contents of our json account secret in our environment; dump it
# to disk for use in tests.
echo "${google_json_account}" > /tmp/google-account.json

set -x

# Create GOPATH structure
mkdir -p "${GOPATH}/src/github.com/terraform-providers"
ln -s "${PWD}/terraform-provider-google" "${GOPATH}/src/github.com/terraform-providers/terraform-provider-google"

cd "${GOPATH}/src/github.com/terraform-providers/terraform-provider-google"
#
go test -v ./google -parallel 16 -run '^TestAcc' -timeout 120m
