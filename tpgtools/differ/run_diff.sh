#!/bin/bash
set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
pushd $SCRIPT_DIR
TPG_ROOT=$(mktemp -d -t tpg-XXXXXX)
git clone --depth 1 https://github.com/terraform-providers/terraform-provider-google-beta $TPG_ROOT/oldtpg
git clone --depth 1 https://github.com/terraform-providers/terraform-provider-google-beta $TPG_ROOT/newtpg
cp -r ../../ $TPG_ROOT/mm
popd
pushd $TPG_ROOT/mm/tpgtools/differ
pushd ../../
make OUTPUT_PATH=$TPG_ROOT/newtpg VERSION=beta
popd
go mod edit -require=newtpg@v0.0.0
go mod edit -replace=newtpg=$TPG_ROOT/newtpg
pushd $TPG_ROOT/newtpg
go mod edit -module newtpg
popd
go mod edit -require=oldtpg@v0.0.0
go mod edit -replace=oldtpg=$TPG_ROOT/oldtpg
pushd $TPG_ROOT/oldtpg
go mod edit -module oldtpg
popd

go clean -modcache -cache
go mod tidy
go run ./ --resource=google_$1_$2

popd
