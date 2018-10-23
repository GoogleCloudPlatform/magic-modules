#!/usr/bin/env bash

set -e
set -x

echo 'TODO slevenick write tests'
pushd "magic-modules/build/inspec/libraries"

rspec -I . ../test/unit

popd