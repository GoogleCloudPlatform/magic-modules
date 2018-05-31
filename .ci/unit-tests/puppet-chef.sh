#!/bin/bash
set -e
set -x

pushd "magic-modules/build/$PROVIDER/$PRODUCT"
if [ -n "$SPEC_DIR" ]; then
  cd "$SPEC_DIR"
fi
bundle install

DISABLE_COVERAGE=true bundle exec parallel_rspec -o '--exclude_pattern "$EXCLUDE_PATTERN"' spec/
popd
