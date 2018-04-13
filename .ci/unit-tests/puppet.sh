#!/bin/bash
set -e
set -x

pushd "magic-modules/build/puppet/$PRODUCT"
bundle install
bundle exec rspec
popd
