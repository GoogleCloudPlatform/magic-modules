#!/bin/bash
set -e
set -x

pushd "magic-modules/build/$PROVIDER/$PRODUCT"
bundle install

# parallel_rspec doesn't support --exclude_pattern
if [ -z "$EXCLUDE_PATTERN" ]; then
    echo "No EXCLUDE_PATTERN"
    ls -la
    DISABLE_COVERAGE=true bundle exec parallel_rspec spec/
else
    echo "Excluding $EXCLUDE_PATTERN"
    IFS="," read -ra excluded <<< "$EXCLUDE_PATTERN"
    filtered=$(find spec -name '*_spec.rb' $(printf "! -wholename %s " ${excluded[@]}))
    DISABLE_COVERAGE=true bundle exec parallel_rspec ${filtered[@]}
fi

popd
