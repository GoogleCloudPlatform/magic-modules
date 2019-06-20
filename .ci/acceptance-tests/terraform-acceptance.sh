#!/usr/bin/env bash

set -e
set -x

export GOOGLE_CREDENTIALS_FILE="/tmp/google-account.json"
export GOOGLE_REGION="us-central1"
export GOOGLE_ZONE="us-central1-a"
# Setup GOPATH
export GOPATH=${PWD}/go

# CI sets the contents of our json account secret in our environment; dump it
# to disk for use in tests.
set +x
echo "${GOOGLE_JSON_ACCOUNT}" > $GOOGLE_CREDENTIALS_FILE
set -x

# Create GOPATH structure
mkdir -p "${GOPATH}/src/github.com/terraform-providers"
ln -s "${PWD}/magic-modules/build/$SHORT_NAME" "${GOPATH}/src/github.com/terraform-providers/$PROVIDER_NAME"

cd "${GOPATH}/src/github.com/terraform-providers/$PROVIDER_NAME"

git diff HEAD~ > tmp.diff
OUTPUT=( $(go run scripts/affectedtests/affectedtests.go -diff tmp.diff) )
rm tmp.diff

if [ ${#OUTPUT[@]} -eq 0 ]; then
	echo "No tests to run"
else
	make testacc TEST=./$TEST_DIR TESTARGS="-run=\"$( IFS=$'|'; echo "${OUTPUT[*]}" )\""
fi
