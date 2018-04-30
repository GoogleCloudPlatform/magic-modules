#!/bin/bash
set -e
set -x

pushd "magic-modules/build/$PROVIDER/$PRODUCT"
bundle install
bundle exec rspec --exclude_pattern "$EXCLUDE_PATTERN"
popd
