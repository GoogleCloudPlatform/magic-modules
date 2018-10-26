#!/usr/bin/env bash

set -e
set -x

pushd "magic-modules/build/inspec/test/unit"

bundle install
rspec -I ../../libraries *

popd