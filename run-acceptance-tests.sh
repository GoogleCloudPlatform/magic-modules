#!/bin/sh

export GOOGLE_CREDENTIALS_FILE="/tmp/google-account.json"
export GCLOUD_PROJECT="terraform-ci-acc-tests"
export TF_ACC=1
export GOOGLE_REGION="us-central1"
# TODO actually use a separate project for xpn resources
export GOOGLE_XPN_HOST_PROJECT="man-i-wish-i-was-a-real-project"

# CI sets the contents of our json account secret in our environment; dump it
# to disk for use in tests.
echo "${GOOGLE_ACCOUNT_JSON}" > /tmp/google-account.json

cd terraform
make testacc TEST=./google
