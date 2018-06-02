#!/bin/bash
set -e
set -x

pushd "magic-modules/build/$PROVIDER/$PRODUCT"
if [ -n "$SPEC_DIR" ]; then
  cd "$SPEC_DIR"
fi
bundle install

# parallel_rspec doesn't support --exclude_pattern
test_files=(spec/**/*_spec.rb)
IFS="," read -ra excluded <<< "$EXCLUDE_PATTERN"
filtered=$(echo ${test_files[@]} ${excluded[@]} | tr " " "\n" | sort | uniq -u | tr "\n" " ")

DISABLE_COVERAGE=true bundle exec parallel_rspec ${filtered[@]}
popd
