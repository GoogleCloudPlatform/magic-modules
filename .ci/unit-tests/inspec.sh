#!/usr/bin/env bash

set -e
set -x

pushd "magic-modules/build/inspec/libraries"

bundle install
rspec -I . ../test/unit/*

popd