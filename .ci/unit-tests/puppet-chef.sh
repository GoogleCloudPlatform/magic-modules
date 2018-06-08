#!/bin/bash
set -e
set -x

pushd "magic-modules/build/$PROVIDER/$PRODUCT"
if [ -n "$SPEC_DIR" ]; then
  cd "$SPEC_DIR"
fi
bundle install

# parallel_rspec doesn't support --exclude_pattern
IFS="," read -ra excluded <<< "$EXCLUDE_PATTERN"
filtered=$(find spec/ -name '*_spec.rb' $(printf "! -wholename %s " ${excluded[@]}))

DISABLE_COVERAGE=true bundle exec parallel_rspec ${filtered[@]}
popd
